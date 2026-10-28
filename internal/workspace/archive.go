package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// Archive is what `record archive` moved, and what it did to the links that
// pointed into or out of the moved records.
type Archive struct {
	Records []ArchivedRecord `json:"records"`
	// AlreadyArchived names records asked for by ID that were archived
	// before; they are left where they are.
	AlreadyArchived []string `json:"already_archived,omitempty"`
	// Unfinished names the moved records whose work is not done, which only
	// naming them by ID archives.
	Unfinished []string `json:"unfinished,omitempty"`
	// LinksRewritten are inline relative links the move would have broken,
	// and LinksReported the ones it would break that it cannot rewrite.
	LinksRewritten []LinkChange `json:"links_rewritten,omitempty"`
	LinksReported  []LinkChange `json:"links_reported,omitempty"`

	// Changed names every record whose files the archive moved or rewrote,
	// which is what is left to commit.
	Changed []string `json:"-"`
}

// ArchivedRecord is one record moved under _archive/. From and To name the
// entry file; a plan folder moves whole.
type ArchivedRecord struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by,omitempty"`
	Status    string `json:"status"`
	From      string `json:"from"`
	To        string `json:"to"`
}

// LinkChange is one link in a file, at the file's path after the move. To is
// what it was rewritten to, or, for a reported link, what it should now be.
type LinkChange struct {
	File string `json:"file"`
	From string `json:"from"`
	To   string `json:"to"`
}

// ArchiveRecords moves records under the _archive/ beside them. With no IDs it
// takes the active member's finished records; with IDs it takes exactly those,
// whoever wrote them and whether or not they are finished, as archiving by
// hand always could. ID files are never touched, so every number stays
// reserved, and lookup by ID finds a record wherever it sits.
func (s *Store) ArchiveRecords(ids []string) (Archive, error) {
	result := Archive{Records: []ArchivedRecord{}}
	records, archived, err := s.allRecords()
	if err != nil {
		return result, err
	}
	byID := indexByID(records)
	var chosen []Record
	if len(ids) == 0 {
		member, err := s.ActiveMember()
		if err != nil {
			return result, err
		}
		chosen = unarchivedFinished(member, records, archived, byID)
	} else {
		ids = slices.Clone(ids)
		slices.Sort(ids)
		ids = slices.Compact(ids)
		var unknown []string
		for _, id := range ids {
			record, ok := byID[id]
			switch {
			case !recordPattern.MatchString(id):
				return result, fmt.Errorf("nothing was archived: %q is not an intent or plan ID, such as i100 or p1000", id)
			case !ok:
				unknown = append(unknown, id)
			case archived[id]:
				result.AlreadyArchived = append(result.AlreadyArchived, id)
			default:
				chosen = append(chosen, record)
			}
		}
		if len(unknown) > 0 {
			return result, fmt.Errorf("nothing was archived: no record matches %s", strings.Join(unknown, ", "))
		}
	}
	if len(chosen) == 0 {
		return result, nil
	}
	moves := map[string]string{} // entry file, or plan folder, -> its new path
	for _, record := range chosen {
		to := placement.ArchivedPath(record.Path)
		if to == record.Path {
			return result, fmt.Errorf("nothing was archived: %s cannot be archived from %s", record.ID, record.Path)
		}
		from, dst := record.Path, to
		if folderPlan(record) {
			from, dst = path.Dir(from), path.Dir(dst)
		}
		p, err := s.Path(dst)
		if err != nil {
			return result, err
		}
		if _, err := os.Lstat(p); err == nil {
			return result, fmt.Errorf("nothing was archived: %s already exists; move or remove it by hand first", dst)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return result, err
		}
		moves[from] = dst
		result.Changed = append(result.Changed, record.ID)
		result.Records = append(result.Records, ArchivedRecord{record.ID, record.CreatedBy, deletedStatus(record), record.Path, to})
		if !finished(record, byID) {
			result.Unfinished = append(result.Unfinished, record.ID)
		}
	}
	rewrites, err := s.planLinkRewrites(records, moves, &result)
	if err != nil {
		return result, err
	}
	for _, from := range slices.Sorted(maps.Keys(moves)) {
		src, err := s.Path(from)
		if err != nil {
			return result, err
		}
		dst, err := s.Path(moves[from])
		if err != nil {
			return result, err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return result, err
		}
		if err := os.Rename(src, dst); err != nil {
			return result, err
		}
	}
	for _, file := range slices.Sorted(maps.Keys(rewrites)) {
		p, err := s.Path(file)
		if err != nil {
			return result, err
		}
		info, err := os.Stat(p)
		if err != nil {
			return result, err
		}
		if err := s.Write(file, rewrites[file], info.Mode().Perm()); err != nil {
			return result, err
		}
	}
	slices.Sort(result.Changed)
	result.Changed = slices.Compact(result.Changed)
	return result, nil
}

// movedTo is where a path is after the moves: its own new path, or the new
// path of the moved folder it sits in; ok is false for a path nothing moved.
func movedTo(rel string, moves map[string]string) (string, bool) {
	if to, ok := moves[rel]; ok {
		return to, true
	}
	for from, to := range moves {
		if rest, ok := strings.CutPrefix(rel, from+"/"); ok {
			return to + "/" + rest, true
		}
	}
	return rel, false
}

// inlineLink matches an inline Markdown link or image target, with an optional
// title: ](target) or ](target "title"). A target in angle brackets is matched
// separately, because it is reported rather than rewritten.
var inlineLink = regexp.MustCompile(`\]\(\s*([^()\s<>]+)((?:\s+(?:"[^"]*"|'[^']*'))?\s*)\)`)
var angleLink = regexp.MustCompile(`\]\(\s*<([^>]*)>`)
var referenceDefinition = regexp.MustCompile(`^\s{0,3}\[[^\]]+\]:\s*<?([^\s>]+)>?`)
var schemeLink = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

// markdownFiles maps every Markdown file of the records to its record's ID:
// an entry file, and for a plan folder every .md file in it.
func (s *Store) markdownFiles(records []Record) (map[string]string, error) {
	files := map[string]string{}
	for _, record := range records {
		if !folderPlan(record) {
			files[record.Path] = record.ID
			continue
		}
		dir := path.Dir(record.Path)
		base, err := s.Path(dir)
		if err != nil {
			return nil, err
		}
		err = filepath.WalkDir(base, func(p string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".md") {
				rel, _ := filepath.Rel(s.Root, p)
				files[filepath.ToSlash(rel)] = record.ID
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

// planLinkRewrites works out, before anything moves, the new contents of every
// record file whose inline relative links the moves would break: links out of
// a moved file and links into one. A link between two files that move
// together, as inside one plan folder, is left alone because it still holds.
// Anything it cannot rewrite safely is reported instead. The result is keyed
// by each file's path after the moves.
func (s *Store) planLinkRewrites(records []Record, moves map[string]string, result *Archive) (map[string][]byte, error) {
	files, err := s.markdownFiles(records)
	if err != nil {
		return nil, err
	}
	rewrites := map[string][]byte{}
	for _, file := range slices.Sorted(maps.Keys(files)) {
		data, err := s.Read(file)
		if err != nil {
			return nil, err
		}
		newFile, fileMoved := movedTo(file, moves)
		// fix returns the target a link should have after the moves, and
		// whether the moves change it.
		fix := func(target string) (string, bool) {
			if target == "" || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "/") || schemeLink.MatchString(target) {
				return target, false
			}
			pathPart, suffix := target, ""
			if i := strings.IndexAny(target, "#?"); i >= 0 {
				pathPart, suffix = target[:i], target[i:]
			}
			old := path.Clean(path.Join(path.Dir(file), pathPart))
			if old == ".." || strings.HasPrefix(old, "../") {
				return target, false
			}
			if p, err := s.Path(old); err != nil {
				return target, false
			} else if _, err := os.Lstat(p); err != nil {
				// A link already broken is not this move's to fix.
				return target, false
			}
			newTarget, targetMoved := movedTo(old, moves)
			if !fileMoved && !targetMoved {
				return target, false
			}
			rel, err := filepath.Rel(filepath.FromSlash(path.Dir(newFile)), filepath.FromSlash(newTarget))
			if err != nil {
				return target, false
			}
			rel = filepath.ToSlash(rel)
			if strings.HasSuffix(pathPart, "/") {
				rel += "/"
			}
			if rel+suffix == target {
				return target, false
			}
			return rel + suffix, true
		}
		lines := strings.SplitAfter(string(data), "\n")
		changed, fenced := false, false
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
				fenced = !fenced
				continue
			}
			if fenced {
				continue
			}
			if m := referenceDefinition.FindStringSubmatch(line); m != nil {
				if to, ok := fix(m[1]); ok {
					result.LinksReported = append(result.LinksReported, LinkChange{newFile, m[1], to})
				}
				continue
			}
			for _, m := range angleLink.FindAllStringSubmatch(line, -1) {
				if to, ok := fix(m[1]); ok {
					result.LinksReported = append(result.LinksReported, LinkChange{newFile, m[1], to})
				}
			}
			// Code spans are left alone: only the text between them is read.
			spans := strings.Split(line, "`")
			for k := 0; k < len(spans); k += 2 {
				spans[k] = inlineLink.ReplaceAllStringFunc(spans[k], func(match string) string {
					m := inlineLink.FindStringSubmatch(match)
					to, ok := fix(m[1])
					if !ok {
						return match
					}
					result.LinksRewritten = append(result.LinksRewritten, LinkChange{newFile, m[1], to})
					changed = true
					return "](" + to + m[2] + ")"
				})
			}
			lines[i] = strings.Join(spans, "`")
		}
		if changed {
			rewrites[newFile] = []byte(strings.Join(lines, ""))
			result.Changed = append(result.Changed, files[file])
		}
	}
	return rewrites, nil
}
