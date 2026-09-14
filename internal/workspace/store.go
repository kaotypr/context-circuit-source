package workspace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/gofrs/flock"
)

// Store contains workspace-owned files. Call WithLock around a compound edit.
// Local locks serialize one directory, not independent Git clones.
type Store struct{ Root string }

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
	for _, part := range strings.Split(filepath.ToSlash(rel), "/") {
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
	builder := (&yaml.PathBuilder{}).Root()
	for _, key := range keys {
		builder = builder.Child(key)
	}
	path := builder.Build()
	node, err := path.FilterFile(doc)
	if err == nil && node != nil {
		encoded, e := yaml.Marshal(value)
		if e != nil {
			return nil, e
		}
		if e = replaceYAMLNode(doc, path, node, encoded); e != nil {
			return nil, e
		}
	} else {
		if err != nil && !yaml.IsNotFoundNodeError(err) {
			return nil, err
		}
		parent := (&yaml.PathBuilder{}).Root()
		for _, key := range keys[:len(keys)-1] {
			parent = parent.Child(key)
		}
		encoded, e := yaml.Marshal(map[string]any{keys[len(keys)-1]: value})
		if e != nil {
			return nil, e
		}
		parentPath := parent.Build()
		parentNode, e := parentPath.FilterFile(doc)
		if e != nil {
			return nil, e
		}
		mapping, isMapping := parentNode.(*ast.MappingNode)
		flow := isMapping && mapping.IsFlowStyle
		// A seed writes an empty container as `{}`. It carries no entries to
		// stay consistent with, so rebuild it in block style through the node
		// that holds it instead of letting every later entry extend one line.
		if flow && len(mapping.Values) == 0 {
			rebuilt, e := rebuildEmptyContainer(doc, keys, value)
			if e != nil {
				return nil, e
			}
			if rebuilt {
				return serializeYAML(doc)
			}
		}
		// A populated flow container must receive a flow value. Mixing styles
		// through MergeFromReader produces malformed YAML in v1.19.2.
		if flow {
			encoded, e = yaml.MarshalWithOptions(map[string]any{keys[len(keys)-1]: value}, yaml.Flow(true))
			if e != nil {
				return nil, e
			}
		}
		e = parentPath.MergeFromReader(doc, bytes.NewReader(encoded))
		if e != nil {
			return nil, e
		}
	}
	return serializeYAML(doc)
}

// rebuildEmptyContainer replaces an empty flow container at keys[:len(keys)-1]
// with a block mapping holding the new entry, by merging into the node that
// holds it. It reports false when that holder is itself a flow mapping, because
// a block value merged into flow would not serialize.
func rebuildEmptyContainer(doc *ast.File, keys []string, value any) (bool, error) {
	if len(keys) < 2 {
		return false, nil
	}
	holder := (&yaml.PathBuilder{}).Root()
	for _, key := range keys[:len(keys)-2] {
		holder = holder.Child(key)
	}
	holderPath := holder.Build()
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
	switch previous := old.(type) {
	case *ast.MappingNode:
		if mapping, ok := next.(*ast.MappingNode); ok && len(previous.Values) > 0 {
			mapping.SetIsFlowStyle(previous.IsFlowStyle)
		}
	case *ast.SequenceNode:
		if sequence, ok := next.(*ast.SequenceNode); ok && len(previous.Values) > 0 {
			sequence.SetIsFlowStyle(previous.IsFlowStyle)
		}
	}
	if err := next.SetComment(old.GetComment()); err != nil {
		return err
	}
	return path.ReplaceWithNode(doc, next)
}

func (s *Store) Update(rel string, keys []string, value any, mode fs.FileMode) error {
	before, err := s.Read(rel)
	if err != nil {
		return err
	}
	after, err := Edit(before, keys, value)
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
	sort.Strings(names)
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
