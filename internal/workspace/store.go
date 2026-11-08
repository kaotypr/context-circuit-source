package workspace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
	"github.com/gofrs/flock"
)

// Store contains workspace-owned files. Call WithLock around a compound edit.
// Local locks serialize one directory, not independent Git clones.
type Store struct {
	Root string
	// Command is the PATH command the running CLI is installed as, named in a
	// remedy the workspace has no launcher for; empty means the public one.
	Command string
	// CLIVersion is the running CLI's version, which an extension manifest's
	// requires_cli is compared with; empty means a store nobody told.
	CLIVersion string

	// archivesMoved names the legacy archive folders and the _archive/ folders
	// they moved to in this run, so the commit recommendation of the command
	// that renamed them names the moved records too.
	archivesMoved []string
}

func Open(root string) (*Store, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("workspace is not a directory: %s", abs)
	}
	return &Store{Root: abs}, nil
}

func (s *Store) Path(rel string) (string, error) {
	if !fs.ValidPath(filepath.ToSlash(rel)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("invalid workspace path: %s", rel)
	}
	current := s.Root
	for part := range strings.SplitSeq(filepath.ToSlash(rel), "/") {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("workspace files cannot be symlinks: %s", rel)
		}
	}
	return filepath.Join(s.Root, filepath.FromSlash(rel)), nil
}

func (s *Store) Read(rel string) ([]byte, error) {
	p, err := s.Path(rel)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return nil, fmt.Errorf("expected a regular workspace file under 16 MiB: %s", rel)
	}
	return os.ReadFile(p)
}

func Decode(data []byte, value any) error {
	doc, err := parser.ParseBytes(data, 0)
	if err != nil {
		return err
	}
	if len(doc.Docs) != 1 {
		return errors.New("expected exactly one YAML document")
	}
	return yaml.NewDecoder(bytes.NewReader(data), yaml.DisallowUnknownField()).Decode(value)
}

func (s *Store) YAML(rel string, value any) error {
	data, err := s.Read(rel)
	if err != nil {
		return err
	}
	if err := Decode(data, value); err != nil {
		return fmt.Errorf("%s: %w", rel, err)
	}
	return nil
}

func (s *Store) Write(rel string, data []byte, mode fs.FileMode) error {
	p, err := s.Path(rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".cc-write-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// Check again after preparing the replacement; don't follow a swapped link.
	if _, err := s.Path(rel); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), p)
}

func (s *Store) WriteYAML(rel string, value any, mode fs.FileMode) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	return s.Write(rel, data, mode)
}

// Edit changes only the selected YAML node. Other nodes and comments stay in the
// parsed document instead of being discarded by a struct round trip.
func Edit(data []byte, keys []string, value any) ([]byte, error) {
	if len(keys) == 0 {
		return nil, errors.New("missing YAML field")
	}
	doc, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(doc.Docs) != 1 {
		return nil, errors.New("expected exactly one YAML document")
	}
	path := yamlPath(keys).Build()
	node, err := path.FilterFile(doc)
	if err == nil && node != nil {
		if err := replaceEditedNode(doc, path, node, value); err != nil {
			return nil, err
		}
		return serializeYAML(doc)
	}
	if err != nil && !yaml.IsNotFoundNodeError(err) {
		return nil, err
	}
	if err := addMissingNode(doc, keys, value); err != nil {
		return nil, err
	}
	return serializeYAML(doc)
}

// listIndex matches a key that addresses one item of a list, written "[N]", so
// an edit can change one entry of a list in place and keep the comments around
// the rest instead of rewriting the whole list.
var listIndex = regexp.MustCompile(`^\[([0-9]+)\]$`)

func isListIndex(key string) bool { return listIndex.MatchString(key) }

// ListIndex is the key Edit and Update read as the index-th item of a list.
func ListIndex(index int) string { return "[" + strconv.Itoa(index) + "]" }

// yamlPath builds the path to the node keys name: a map key, or "[N]" for the
// Nth item of a list.
func yamlPath(keys []string) *yaml.PathBuilder {
	builder := (&yaml.PathBuilder{}).Root()
	for _, key := range keys {
		if match := listIndex.FindStringSubmatch(key); match != nil {
			index, _ := strconv.Atoi(match[1])
			builder = builder.Index(uint(index))
			continue
		}
		builder = builder.Child(key)
	}
	return builder
}

// rebuildEmptyContainer replaces an empty flow container at keys[:len(keys)-1]
// with a block mapping holding the new entry, by merging into the node that
// holds it. It reports false when that holder is itself a flow mapping, because
// a block value merged into flow would not serialize.
func rebuildEmptyContainer(doc *ast.File, keys []string, value any) (bool, error) {
	if len(keys) < 2 || isListIndex(keys[len(keys)-2]) {
		return false, nil
	}
	holderPath := yamlPath(keys[:len(keys)-2]).Build()
	holderNode, err := holderPath.FilterFile(doc)
	if err != nil {
		return false, err
	}
	if outer, ok := holderNode.(*ast.MappingNode); ok && outer.IsFlowStyle {
		return false, nil
	}
	container := map[string]any{keys[len(keys)-2]: map[string]any{keys[len(keys)-1]: value}}
	encoded, err := yaml.Marshal(container)
	if err != nil {
		return false, err
	}
	if err := holderPath.MergeFromReader(doc, bytes.NewReader(encoded)); err != nil {
		return false, err
	}
	return true, nil
}

func serializeYAML(doc *ast.File) ([]byte, error) {
	result := []byte(doc.String())
	if _, err := parser.ParseBytes(result, 0); err != nil {
		return nil, fmt.Errorf("YAML edit could not be serialized: %w", err)
	}
	return result, nil
}

func replaceYAMLNode(doc *ast.File, path *yaml.Path, old ast.Node, encoded []byte) error {
	replacement, err := parser.ParseBytes(encoded, parser.ParseComments)
	if err != nil {
		return err
	}
	next := replacement.Docs[0].Body
	// Preserve a populated container's flow style when replacing it within a
	// flow map. An empty container has no entries whose style must be matched,
	// so it adopts block style and stays readable as entries accumulate.
	grew := false
	switch previous := old.(type) {
	case *ast.MappingNode:
		if mapping, ok := next.(*ast.MappingNode); ok {
			if len(previous.Values) > 0 {
				mapping.SetIsFlowStyle(previous.IsFlowStyle)
			} else {
				grew = true
			}
		}
	case *ast.SequenceNode:
		if sequence, ok := next.(*ast.SequenceNode); ok {
			if len(previous.Values) > 0 {
				sequence.SetIsFlowStyle(previous.IsFlowStyle)
			} else {
				grew = true
			}
		}
	}
	if err := next.SetComment(old.GetComment()); err != nil {
		return err
	}
	column := keyColumn(doc, old)
	if err := path.ReplaceWithNode(doc, next); err != nil {
		return err
	}
	// A block container prints itself from its own start column, and replacing
	// carries over the column the `[]` or `{}` occupied after the key. Left
	// alone the first entry lands under the value, so the indent depends on how
	// long the key happens to be. Align it with the key, after the replace that
	// would otherwise overwrite it.
	if grew && column > 0 {
		indentNode(next, column)
	}
	return nil
}

// keyColumn reports the column of the key whose value is target, or 0 when the
// node is not a mapping value (a document root has no key to align to).
func keyColumn(doc *ast.File, target ast.Node) int {
	found := 0
	var walk func(ast.Node)
	walk = func(node ast.Node) {
		if found > 0 || node == nil {
			return
		}
		switch typed := node.(type) {
		case *ast.MappingNode:
			for _, value := range typed.Values {
				walk(value)
			}
		case *ast.MappingValueNode:
			if typed.Value == target {
				if token := typed.Key.GetToken(); token != nil {
					found = token.Position.Column
				}
				return
			}
			walk(typed.Value)
		case *ast.SequenceNode:
			for _, value := range typed.Values {
				walk(value)
			}
		}
	}
	for _, document := range doc.Docs {
		walk(document.Body)
	}
	return found
}

// indentNode moves a block container's own start column, which is what its
// printer measures every entry's indentation from.
func indentNode(node ast.Node, column int) {
	set := func(t *token.Token) {
		if t != nil {
			t.Position.Column = column
		}
	}
	switch typed := node.(type) {
	case *ast.SequenceNode:
		set(typed.Start)
		for _, value := range typed.Values {
			set(value.GetToken())
		}
	case *ast.MappingNode:
		set(typed.Start)
		for _, value := range typed.Values {
			set(value.GetToken())
		}
	}
}

func (s *Store) Update(rel string, keys []string, value any, mode fs.FileMode) error {
	return s.change(rel, mode, func(before []byte) ([]byte, error) { return Edit(before, keys, value) })
}

// change rewrites one file through edit, refusing when the file moved under it.
func (s *Store) change(rel string, mode fs.FileMode, edit func([]byte) ([]byte, error)) error {
	before, err := s.Read(rel)
	if err != nil {
		return err
	}
	after, err := edit(before)
	if err != nil {
		return err
	}
	now, err := s.Read(rel)
	if err != nil {
		return err
	}
	if !bytes.Equal(now, before) {
		return fmt.Errorf("%s changed during this edit; retry", rel)
	}
	return s.Write(rel, after, mode)
}

// AppendListItem and RemoveListItem applied to a workspace file, refusing when
// the file moved under the edit.
func (s *Store) AppendItem(rel string, keys []string, value any, mode fs.FileMode) error {
	return s.change(rel, mode, func(before []byte) ([]byte, error) { return AppendListItem(before, keys, value) })
}

func (s *Store) RemoveItem(rel string, keys []string, index int, mode fs.FileMode) error {
	return s.change(rel, mode, func(before []byte) ([]byte, error) { return RemoveListItem(before, keys, index) })
}

func (s *Store) WithLock(ctx context.Context, fn func() (any, error)) (any, error) {
	p, err := s.Path(".context-circuit/local/write.lock")
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return nil, err
	}
	lock := flock.New(p, flock.SetPermissions(0600))
	defer lock.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ok, err := lock.TryLockContext(ctx, 25*time.Millisecond)
	if err != nil || !ok {
		return nil, fmt.Errorf("workspace is busy; retry after the current edit: %v", err)
	}
	return fn()
}

// Export only writes into paths that do not exist. A directory containing Git
// metadata or unrelated files is allowed; no generated destination is replaced.
func (s *Store) Export(files map[string][]byte) error {
	names := make([]string, 0, len(files))
	for name := range files {
		p, err := s.Path(name)
		if err != nil {
			return err
		}
		if _, err := os.Lstat(p); !os.IsNotExist(err) {
			return fmt.Errorf("initialization would overwrite %s", name)
		}
		parent := filepath.Dir(p)
		for parent != s.Root {
			if info, err := os.Stat(parent); err == nil && !info.IsDir() {
				return fmt.Errorf("not a directory: %s", parent)
			}
			parent = filepath.Dir(parent)
		}
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		p, _ := s.Path(name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return err
		}
		f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = f.Write(files[name])
		closeErr := f.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

// AppendListItem adds value as the last item of the list at keys, creating the
// list when the key is missing, and leaves every other node and comment alone.
func AppendListItem(data []byte, keys []string, value any) ([]byte, error) {
	doc, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(doc.Docs) != 1 {
		return nil, errors.New("expected exactly one YAML document")
	}
	path := yamlPath(keys).Build()
	node, err := path.FilterFile(doc)
	if err != nil && !yaml.IsNotFoundNodeError(err) {
		return nil, err
	}
	sequence, ok := node.(*ast.SequenceNode)
	if ok && sequence.IsFlowStyle && len(sequence.Values) > 0 {
		return nil, fmt.Errorf("%s is a one-line list; write it as a block list first", strings.Join(keys, "."))
	}
	if err != nil || !ok || len(sequence.Values) == 0 {
		// A missing or empty list carries no entries whose layout must be
		// kept, so it is written again in block style.
		return Edit(data, keys, []any{value})
	}
	encoded, err := yaml.Marshal([]any{value})
	if err != nil {
		return nil, err
	}
	if err := path.MergeFromReader(doc, bytes.NewReader(encoded)); err != nil {
		return nil, err
	}
	return serializeYAML(doc)
}

// RemoveListItem drops the index-th item of the list at keys. A list left
// empty is written as [] so the key itself stays.
func RemoveListItem(data []byte, keys []string, index int) ([]byte, error) {
	doc, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(doc.Docs) != 1 {
		return nil, errors.New("expected exactly one YAML document")
	}
	node, err := yamlPath(keys).Build().FilterFile(doc)
	if err != nil {
		return nil, err
	}
	sequence, ok := node.(*ast.SequenceNode)
	if !ok || index < 0 || index >= len(sequence.Values) {
		return nil, fmt.Errorf("no list item %d under %s", index, strings.Join(keys, "."))
	}
	if len(sequence.Values) == 1 {
		return Edit(data, keys, []any{})
	}
	drop := func(n int) bool { return index < n }
	sequence.Values = slices.Delete(sequence.Values, index, index+1)
	if drop(len(sequence.ValueHeadComments)) {
		sequence.ValueHeadComments = slices.Delete(sequence.ValueHeadComments, index, index+1)
	}
	if drop(len(sequence.Entries)) {
		sequence.Entries = slices.Delete(sequence.Entries, index, index+1)
	}
	return serializeYAML(doc)
}

// RemoveMapKey drops the last of keys from the mapping the rest of keys name,
// and leaves every other node and comment alone. The mapping must keep at
// least one other key: a caller removing the last one removes its parent's key
// instead.
func RemoveMapKey(data []byte, keys []string) ([]byte, error) {
	if len(keys) < 2 {
		return nil, errors.New("missing YAML field")
	}
	doc, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(doc.Docs) != 1 {
		return nil, errors.New("expected exactly one YAML document")
	}
	node, err := yamlPath(keys[:len(keys)-1]).Build().FilterFile(doc)
	if err != nil {
		return nil, err
	}
	last := keys[len(keys)-1]
	mapping, ok := node.(*ast.MappingNode)
	if !ok {
		return nil, fmt.Errorf("%s is not a mapping with %s and another key", strings.Join(keys[:len(keys)-1], "."), last)
	}
	index := slices.IndexFunc(mapping.Values, func(value *ast.MappingValueNode) bool {
		token := value.Key.GetToken()
		return token != nil && token.Value == last
	})
	if index < 0 {
		return nil, fmt.Errorf("no key %s under %s", last, strings.Join(keys[:len(keys)-1], "."))
	}
	if len(mapping.Values) == 1 {
		return nil, fmt.Errorf("%s is the only key under %s", last, strings.Join(keys[:len(keys)-1], "."))
	}
	mapping.Values = slices.Delete(mapping.Values, index, index+1)
	return serializeYAML(doc)
}

// RemoveKey applies RemoveMapKey to a workspace file, refusing when the file
// moved under the edit.
func (s *Store) RemoveKey(rel string, keys []string, mode fs.FileMode) error {
	return s.change(rel, mode, func(before []byte) ([]byte, error) { return RemoveMapKey(before, keys) })
}

// replaceEditedNode overwrites a node that already exists.
func replaceEditedNode(doc *ast.File, path *yaml.Path, node ast.Node, value any) error {
	encoded, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	return replaceYAMLNode(doc, path, node, encoded)
}

// addMissingNode adds the last key under its parent, matching the parent's
// block or flow style.
func addMissingNode(doc *ast.File, keys []string, value any) error {
	if isListIndex(keys[len(keys)-1]) {
		return fmt.Errorf("no list item %s to edit", keys[len(keys)-1])
	}
	parent := yamlPath(keys[:len(keys)-1])
	encoded, err := yaml.Marshal(map[string]any{keys[len(keys)-1]: value})
	if err != nil {
		return err
	}
	parentPath := parent.Build()
	parentNode, err := parentPath.FilterFile(doc)
	if err != nil {
		return err
	}
	mapping, isMapping := parentNode.(*ast.MappingNode)
	flow := isMapping && mapping.IsFlowStyle
	// A seed writes an empty container as `{}`. It carries no entries to
	// stay consistent with, so rebuild it in block style through the node
	// that holds it instead of letting every later entry extend one line.
	if flow && len(mapping.Values) == 0 {
		rebuilt, err := rebuildEmptyContainer(doc, keys, value)
		if err != nil {
			return err
		}
		if rebuilt {
			return nil
		}
	}
	// A populated flow container must receive a flow value. Mixing styles
	// through MergeFromReader produces malformed YAML in v1.19.2.
	if flow {
		encoded, err = yaml.MarshalWithOptions(map[string]any{keys[len(keys)-1]: value}, yaml.Flow(true))
		if err != nil {
			return err
		}
	}
	return parentPath.MergeFromReader(doc, bytes.NewReader(encoded))
}
