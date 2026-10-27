package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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
	// PlanFolderParents names the folders whose children are plan folders.
	PlanFolderParents() []string
	// GitPathMatcher matches every path Git may report for one record ID: its
	// entry file and, for a plan folder, everything in it.
	GitPathMatcher(id string) func(rel string) bool
}

// placement is the layout every caller uses. A fork that places records
// differently replaces this value.
var placement Placement = flatPlacement{}

// flatPlacement is the shipped layout: intent/iNNN-slug.md and
// plans/pNNNN-slug/plan.md, each with an archive/ folder beside its records.
type flatPlacement struct{}

const archiveFolder = "archive"

func (flatPlacement) Roots() []string { return []string{"intent", "plans"} }

func (flatPlacement) PlanFolderParents() []string {
	return []string{"plans", "plans/" + archiveFolder}
}

func (flatPlacement) NewRecordPath(kind, id, slug string) string {
	if kind == "plan" {
		return fmt.Sprintf("plans/%s-%s/plan.md", id, slug)
	}
	return fmt.Sprintf("intent/%s-%s.md", id, slug)
}

func (flatPlacement) ArchivedPath(rel string) string {
	root, rest, ok := strings.Cut(rel, "/")
	if !ok || strings.HasPrefix(rest, archiveFolder+"/") {
		return rel
	}
	return root + "/" + archiveFolder + "/" + rest
}

func (flatPlacement) Classify(rel string) (string, bool, bool) {
	parts := strings.Split(rel, "/")
	if len(parts) < 2 || parts[0] != "intent" && parts[0] != "plans" {
		return "", false, false
	}
	archived := len(parts) > 2 && parts[1] == archiveFolder
	level := parts[1:]
	if archived {
		level = parts[2:]
	}
	switch {
	case len(level) == 1 && recordFilename.MatchString(level[0]):
	case len(level) == 2 && parts[0] == "plans" && level[1] == "plan.md" && recordFilename.MatchString(level[0]+".md"):
	default:
		return "", false, false
	}
	return recordPathID(rel), archived, true
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
			if entry.IsDir() {
				if path != base && !archived && entry.Name() == archiveFolder {
					return filepath.SkipDir
				}
				if folder == "intent" && path != base && entry.Name() != archiveFolder {
					return filepath.SkipDir
				}
				return nil
			}
			rel, _ := filepath.Rel(s.Root, path)
			rel = filepath.ToSlash(rel)
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

func (flatPlacement) GitPathMatcher(id string) func(string) bool {
	quoted := regexp.QuoteMeta(id)
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^intent/(` + archiveFolder + `/)?` + quoted + `-[a-z][a-z0-9-]*\.md$`),
		regexp.MustCompile(`^plans/(` + archiveFolder + `/)?` + quoted + `-[a-z][a-z0-9-]*(\.md|/.+)$`),
	}
	return func(rel string) bool {
		for _, pattern := range patterns {
			if pattern.MatchString(rel) {
				return true
			}
		}
		return false
	}
}
