package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Placement answers where intent and plan records are found and where a new
// one is placed. Every record path the CLI computes goes through it, so it is
// the one seam that changes when records are laid out differently.
//
// Paths are relative to the workspace root and slash-separated.
type Placement interface {
	// Roots names the top-level folders records live under, in walk order.
	Roots() []string
	// RecordPaths lists every record entry file, sorted; archived records only
	// when archived is true.
	RecordPaths(s *Store, archived bool) ([]string, error)
	// Classify reports the record ID a path is the entry file of, and whether
	// it is archived; ok is false for a path that is not a record entry.
	Classify(rel string) (id string, archived bool, ok bool)
	// NewRecordPath is where a new record of kind "intent" or "plan" is written.
	NewRecordPath(kind, id, slug string) string
	// ArchivedPath is where an active record entry goes when it is archived.
	ArchivedPath(rel string) string
	// PlanFolderParents names the existing folders whose children are plan
	// folders.
	PlanFolderParents(s *Store) []string
	// GitPathMatcher matches every path Git may report for one record ID: its
	// entry file and, for a plan folder, everything in it.
	GitPathMatcher(id string) func(rel string) bool
	// LayoutIssues reports what in the record folders this placement cannot
	// read, or reads differently from how it looks.
	LayoutIssues(s *Store) []Finding
}

// placement is the layout every caller uses. A fork that places records
// differently replaces this value.
var placement Placement = flatPlacement{}

// flatPlacement is the shipped layout. A record R is intent/iNNN-slug.md or
// plans/pNNNN-slug/plan.md (or a legacy plans/pNNNN-slug.md), and under each
// root it may sit in one level of group folder, with _archive/ above or below
// that group:
//
//	root/R  root/G/R  root/_archive/R  root/G/_archive/R  root/_archive/G/R
//
// A group is ordinary file organization and means nothing to the CLI. Folder
// names starting with "_" are kept for the CLI, and _archive/ is the only one
// it gives meaning to. A top-level archive/ is the legacy archive folder: it
// always reads as archived, and is never a group, until it is renamed to
// _archive/. An archive/ inside a group is an ordinary folder.
type flatPlacement struct{}

const (
	archiveFolder       = "_archive"
	legacyArchiveFolder = "archive"
	reservedPrefix      = "_"
)

func (flatPlacement) Roots() []string { return []string{"intent", "plans"} }

// isArchiveFolder reports whether a folder directly under a root is an
// archive: _archive/, or the legacy archive/ that awaits its rename.
func isArchiveFolder(name string) bool {
	return name == archiveFolder || name == legacyArchiveFolder
}

// recordFolderName reports whether a folder is named like a record, as a plan
// folder is. Such a folder is never a group.
func recordFolderName(name string) bool {
	return recordFilename.MatchString(name + ".md")
}

// isGroup reports whether a folder name may group records: anything but a
// reserved name, the legacy archive name, or a record's name.
func isGroup(name string) bool {
	return name != "" && !strings.HasPrefix(name, reservedPrefix) && name != legacyArchiveFolder && !recordFolderName(name)
}

// folderChain reports whether records may sit directly in the folders dirs
// below a root, and whether they are archived there.
func folderChain(dirs []string) (archived, ok bool) {
	switch len(dirs) {
	case 0:
		return false, true
	case 1:
		if isArchiveFolder(dirs[0]) {
			return true, true
		}
		return false, isGroup(dirs[0])
	case 2:
		if isArchiveFolder(dirs[0]) && isGroup(dirs[1]) || isGroup(dirs[0]) && dirs[1] == archiveFolder {
			return true, true
		}
	}
	return false, false
}

func (flatPlacement) NewRecordPath(kind, id, slug string) string {
	if kind == "plan" {
		return fmt.Sprintf("plans/%s-%s/plan.md", id, slug)
	}
	return fmt.Sprintf("intent/%s-%s.md", id, slug)
}

// ArchivedPath moves an active record under the _archive/ beside it: its
// root's for a record at the top level, its group's for a record in a group.
// A path that is already archived, or is no record entry, is returned as is.
func (p flatPlacement) ArchivedPath(rel string) string {
	_, archived, ok := p.Classify(rel)
	if !ok || archived {
		return rel
	}
	parts := strings.Split(rel, "/")
	cut := len(parts) - 1
	if parts[0] == "plans" && parts[cut] == "plan.md" {
		cut--
	}
	return strings.Join(parts[:cut], "/") + "/" + archiveFolder + "/" + strings.Join(parts[cut:], "/")
}

// entryDirs splits a path with the shape of a record entry into the folders
// between its root and the record; ok is false for any other path.
func entryDirs(rel string) (dirs []string, ok bool) {
	parts := strings.Split(rel, "/")
	if len(parts) < 2 || parts[0] != "intent" && parts[0] != "plans" {
		return nil, false
	}
	last := len(parts) - 1
	switch {
	case recordFilename.MatchString(parts[last]):
		return parts[1:last], true
	case parts[0] == "plans" && parts[last] == "plan.md" && last >= 2 && recordFolderName(parts[last-1]):
		return parts[1 : last-1], true
	}
	return nil, false
}

func (flatPlacement) Classify(rel string) (string, bool, bool) {
	dirs, ok := entryDirs(rel)
	if !ok {
		return "", false, false
	}
	archived, ok := folderChain(dirs)
	if !ok {
		return "", false, false
	}
	return recordPathID(rel), archived, true
}

// descends reports whether a walk enters the folder at dirs below a root:
// never a reserved folder other than _archive/, and an archive only when
// archived records are wanted.
func descends(dirs []string, archived bool) bool {
	name := dirs[len(dirs)-1]
	if strings.HasPrefix(name, reservedPrefix) && name != archiveFolder {
		return false
	}
	if archived || len(dirs) > 2 {
		return true
	}
	inArchive, ok := folderChain(dirs)
	return !ok || !inArchive
}

func (p flatPlacement) RecordPaths(s *Store, archived bool) ([]string, error) {
	var paths []string
	for _, folder := range p.Roots() {
		base, err := s.Path(folder)
		if err != nil {
			return nil, err
		}
		err = filepath.WalkDir(base, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("record directories cannot contain symlinks: %s", path)
			}
			rel, _ := filepath.Rel(s.Root, path)
			rel = filepath.ToSlash(rel)
			if entry.IsDir() {
				if path != base && !descends(strings.Split(rel, "/")[1:], archived) {
					return filepath.SkipDir
				}
				return nil
			}
			if _, isArchived, ok := p.Classify(rel); ok && (archived || !isArchived) {
				paths = append(paths, rel)
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	slices.Sort(paths)
	return paths, nil
}

// PlanFolderParents lists plans/ and every group and archive folder under it
// that exists, which are the folders the classifier reads plan folders from.
func (flatPlacement) PlanFolderParents(s *Store) []string {
	parents := []string{"plans"}
	var dirs func(rel string) []string
	dirs = func(rel string) []string {
		p, err := s.Path(rel)
		if err != nil {
			return nil
		}
		entries, err := os.ReadDir(p)
		if err != nil {
			return nil
		}
		var names []string
		for _, entry := range entries {
			if entry.IsDir() {
				names = append(names, entry.Name())
			}
		}
		return names
	}
	for _, first := range dirs("plans") {
		if _, ok := folderChain([]string{first}); !ok {
			continue
		}
		parents = append(parents, "plans/"+first)
		for _, second := range dirs("plans/" + first) {
			if _, ok := folderChain([]string{first, second}); ok {
				parents = append(parents, "plans/"+first+"/"+second)
			}
		}
	}
	return parents
}

// GitPathMatcher matches the record's entry file wherever the classifier reads
// it and, for a plan folder, every file inside that folder.
func (p flatPlacement) GitPathMatcher(id string) func(string) bool {
	return func(rel string) bool {
		if got, _, ok := p.Classify(rel); ok {
			return got == id
		}
		parts := strings.Split(rel, "/")
		if len(parts) < 3 || parts[0] != "plans" {
			return false
		}
		for k := 1; k < len(parts)-1; k++ {
			if !strings.HasPrefix(parts[k], id+"-") || !recordFolderName(parts[k]) {
				continue
			}
			if got, _, ok := p.Classify(strings.Join(parts[:k+1], "/") + "/plan.md"); ok && got == id {
				return true
			}
		}
		return false
	}
}

// LayoutIssues walks intent/ and plans/ and reports what the classifier would
// silently pass over or read differently from how it looks: a reserved folder,
// a group named like a record, a record file at a depth no command reads, and
// a legacy archive/ that still awaits its rename to _archive/.
func (p flatPlacement) LayoutIssues(s *Store) []Finding {
	var issues []Finding
	for _, root := range p.Roots() {
		base, err := s.Path(root)
		if err != nil {
			continue
		}
		if info, err := os.Lstat(filepath.Join(base, legacyArchiveFolder)); err == nil && info.IsDir() {
			issues = append(issues, found("legacy archive folder awaits its rename to "+archiveFolder+"/: "+root+"/"+legacyArchiveFolder,
				"run any changing command with the pinned CLI, which renames it to "+root+"/"+archiveFolder+"/ and keeps every record in it archived; one member makes and pushes that change before the others change anything"))
		}
		_ = filepath.WalkDir(base, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil || entry.Type()&os.ModeSymlink != 0 {
				return nil
			}
			rel, _ := filepath.Rel(s.Root, path)
			rel = filepath.ToSlash(rel)
			dirs := strings.Split(rel, "/")[1:]
			if entry.IsDir() {
				if path == base {
					return nil
				}
				name := entry.Name()
				if strings.HasPrefix(name, reservedPrefix) && name != archiveFolder {
					issues = append(issues, found("reserved folder name: "+rel,
						"rename the folder without its leading "+reservedPrefix+"; folder names starting with "+reservedPrefix+" are kept for the CLI, and no record in this one is read"))
					return filepath.SkipDir
				}
				parent := dirs[:len(dirs)-1]
				if _, chainOK := folderChain(parent); chainOK && recordFolderName(name) {
					if _, err := os.Lstat(filepath.Join(path, "plan.md")); err != nil && holdsRecords(path) {
						issues = append(issues, found("folder named like a record holds records: "+rel,
							"rename the folder so it does not start with a record ID; a group folder cannot look like a plan folder"))
						return filepath.SkipDir
					}
					// A plan folder; its supporting files are no records.
					if root == "plans" {
						return filepath.SkipDir
					}
				}
				return nil
			}
			_, entryShape := entryDirs(rel)
			if entryShape {
				if _, _, ok := p.Classify(rel); !ok {
					issues = append(issues, found("record file outside the folders records are read from: "+rel,
						"move it to "+root+"/, one group folder below it, or an "+archiveFolder+"/ above or below that group; no command reads it where it is"))
				}
			}
			return nil
		})
	}
	return issues
}

// holdsRecords reports whether a folder holds a record entry file or a plan
// folder directly.
func holdsRecords(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && recordFilename.MatchString(entry.Name()) || entry.IsDir() && recordFolderName(entry.Name()) {
			return true
		}
	}
	return false
}
