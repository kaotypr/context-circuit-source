package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ArchiveRename reports what renaming the legacy archive/ folders to _archive/
// did. The schema does not change: a top-level archive/ always reads as
// archived, so nothing reads differently before or after, and ID files are
// untouched, so every number stays reserved.
type ArchiveRename struct {
	// Renamed names each legacy folder moved whole, as "intent/archive".
	Renamed []string `json:"renamed"`
	// Merged names each legacy folder whose contents joined an _archive/ that
	// already existed.
	Merged []string `json:"merged"`
}

// folders lists every legacy folder the rename moved, in root order.
func (r ArchiveRename) folders() []string {
	return append(append([]string{}, r.Renamed...), r.Merged...)
}

// Notice is the paragraph a person reads on the error stream when a command
// renamed the archive folders before doing what it was asked.
func (r ArchiveRename) Notice() string {
	var moves []string
	for _, from := range r.Renamed {
		moves = append(moves, fmt.Sprintf("renamed %s to %s", from, archiveTarget(from)))
	}
	for _, from := range r.Merged {
		moves = append(moves, fmt.Sprintf("moved the contents of %s into %s", from, archiveTarget(from)))
	}
	var targets []string
	for _, from := range r.folders() {
		targets = append(targets, archiveTarget(from)+"/")
	}
	return "notice: " + strings.Join(moves, "; ") + ", the one archive folder this CLI reads; every record in it stays archived and keeps its reserved number. " +
		"Commit and push the moved records in " + strings.Join(targets, " and ") + " before anyone else changes the workspace, so every clone reads them in the same place; committing is the person's decision.\n"
}

// archiveTarget is the _archive/ folder a legacy archive folder moves to.
func archiveTarget(legacy string) string {
	return strings.TrimSuffix(legacy, legacyArchiveFolder) + archiveFolder
}

// legacyArchives lists every top-level legacy archive/ folder still present.
func (s *Store) legacyArchives() ([]string, error) {
	var found []string
	for _, root := range placement.Roots() {
		rel := root + "/" + legacyArchiveFolder
		p, err := s.Path(rel)
		if err != nil {
			return nil, err
		}
		info, err := os.Lstat(p)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			found = append(found, rel)
		}
	}
	return found, nil
}

// ArchiveRenamePending reports whether a top-level intent/archive/ or
// plans/archive/ still awaits its rename to _archive/, including one a clone
// brought back after the rename.
func (s *Store) ArchiveRenamePending() (bool, error) {
	found, err := s.legacyArchives()
	return len(found) > 0, err
}

// archiveMove is one entry the rename moves.
type archiveMove struct{ from, to string }

// RenameLegacyArchives renames each top-level archive/ to _archive/. Where an
// _archive/ already exists the two are merged: a name only one side holds is
// moved, a folder both hold that is no record is merged one level further, and
// an identical file both hold is dropped from archive/. Any other name both
// hold is a clash, and the rename refuses before moving anything.
func (s *Store) RenameLegacyArchives() (*ArchiveRename, error) {
	legacy, err := s.legacyArchives()
	if err != nil || len(legacy) == 0 {
		return nil, err
	}
	rename := &ArchiveRename{Renamed: []string{}, Merged: []string{}}
	var moves []archiveMove
	var drops, emptied []string
	var clashes []string
	for _, rel := range legacy {
		from, err := s.Path(rel)
		if err != nil {
			return nil, err
		}
		to, err := s.Path(archiveTarget(rel))
		if err != nil {
			return nil, err
		}
		if _, err := os.Lstat(to); os.IsNotExist(err) {
			moves = append(moves, archiveMove{from, to})
			rename.Renamed = append(rename.Renamed, rel)
			continue
		} else if err != nil {
			return nil, err
		}
		rename.Merged = append(rename.Merged, rel)
		if err := planMerge(from, to, &moves, &drops, &emptied, &clashes); err != nil {
			return nil, err
		}
	}
	if len(clashes) > 0 {
		var names []string
		for _, clash := range clashes {
			rel, _ := filepath.Rel(s.Root, clash)
			names = append(names, filepath.ToSlash(rel))
		}
		return nil, fmt.Errorf("the legacy archive and %s/ both hold %s; nothing was moved. Keep one of each by hand, moving the other out of the legacy archive folder, then run the command again", archiveFolder, strings.Join(names, ", "))
	}
	for _, move := range moves {
		if err := os.Rename(move.from, move.to); err != nil {
			return nil, err
		}
	}
	for _, drop := range drops {
		if err := os.Remove(drop); err != nil {
			return nil, err
		}
	}
	// Folders emptied by a merge are removed deepest first; os.Remove leaves
	// a folder that still holds anything.
	for i := len(emptied) - 1; i >= 0; i-- {
		if err := os.Remove(emptied[i]); err != nil {
			return nil, err
		}
	}
	for _, rel := range rename.folders() {
		s.archivesMoved = append(s.archivesMoved, rel, archiveTarget(rel))
	}
	return rename, nil
}

// planMerge works out how the folder from merges into the folder to, which
// both exist, without changing anything. from is listed in emptied before any
// folder below it, so removing in reverse removes children first.
func planMerge(from, to string, moves *[]archiveMove, drops, emptied, clashes *[]string) error {
	entries, err := os.ReadDir(from)
	if err != nil {
		return err
	}
	*emptied = append(*emptied, from)
	for _, entry := range entries {
		src, dst := filepath.Join(from, entry.Name()), filepath.Join(to, entry.Name())
		target, err := os.Lstat(dst)
		if os.IsNotExist(err) {
			*moves = append(*moves, archiveMove{src, dst})
			continue
		}
		if err != nil {
			return err
		}
		switch {
		case entry.IsDir() && target.IsDir() && !recordFolderName(entry.Name()):
			if err := planMerge(src, dst, moves, drops, emptied, clashes); err != nil {
				return err
			}
		case entry.Type().IsRegular() && target.Mode().IsRegular() && sameContent(src, dst):
			*drops = append(*drops, src)
		default:
			*clashes = append(*clashes, src)
		}
	}
	return nil
}

func sameContent(a, b string) bool {
	left, err := os.ReadFile(a)
	if err != nil {
		return false
	}
	right, err := os.ReadFile(b)
	return err == nil && bytes.Equal(left, right)
}
