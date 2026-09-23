package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

func folderPlan(record Record) bool { return filepath.Base(record.Path) == "plan.md" }

// requiredPlanFiles returns workspace-relative paths for one repository. It
// checks every assignment, including those destined for another worker.
func (s *Store) requiredPlanFiles(record Record, repo string) ([]string, error) {
	if repo != "" && !slices.Contains(record.Repositories, repo) {
		return nil, fmt.Errorf("repository %s is not in plan %s", repo, record.ID)
	}
	if record.RequiredFiles == nil {
		return nil, nil
	}
	if !folderPlan(record) {
		return nil, fmt.Errorf("required_files needs a folder plan")
	}
	for key := range record.RequiredFiles.ByRepository {
		if !slices.Contains(record.Repositories, key) {
			return nil, fmt.Errorf("unknown required_files repository: %s", key)
		}
	}
	var selected []string
	check := func(group string, files []string, include bool) error {
		seen := map[string]bool{}
		for _, file := range files {
			if file == "" || path.IsAbs(file) || !fs.ValidPath(file) || strings.Contains(file, "\\") || file == "plan.md" {
				return fmt.Errorf("unsafe required file %q in %s", file, group)
			}
			if seen[file] {
				return fmt.Errorf("duplicate required file %q in %s", file, group)
			}
			seen[file] = true
			rel := path.Join(path.Dir(record.Path), file)
			absolute, err := s.Path(rel)
			if err != nil {
				return fmt.Errorf("required file %s: %w", rel, err)
			}
			info, err := os.Stat(absolute)
			if err != nil {
				return fmt.Errorf("required file %s: %w", rel, err)
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("required file %s is not regular", rel)
			}
			opened, err := os.Open(absolute)
			if err != nil {
				return fmt.Errorf("required file %s: %w", rel, err)
			}
			if err := opened.Close(); err != nil {
				return fmt.Errorf("required file %s: %w", rel, err)
			}
			if include {
				selected = append(selected, rel)
			}
		}
		return nil
	}
	if err := check("shared", record.RequiredFiles.Shared, true); err != nil {
		return nil, err
	}
	for key, files := range record.RequiredFiles.ByRepository {
		if err := check(key, files, key == repo); err != nil {
			return nil, err
		}
	}
	return unique(selected), nil
}

func (s *Store) planFileIssues(record Record) []Finding {
	if !folderPlan(record) {
		return nil
	}
	var issues []Finding
	assigned := map[string]bool{}
	if record.RequiredFiles != nil {
		for _, files := range append([][]string{record.RequiredFiles.Shared}, values(record.RequiredFiles.ByRepository)...) {
			for _, file := range files {
				assigned[file] = true
			}
		}
	}
	if _, err := s.requiredPlanFiles(record, ""); err != nil {
		issues = append(issues, found(record.ID+": "+err.Error(), "correct required_files and restore its regular files inside the plan folder"))
	}
	_, body, _ := splitRecord([]byte(record.Content))
	section := ""
	if at := strings.Index(string(body), "## Details\n"); at >= 0 {
		section = string(body)[at+len("## Details\n"):]
		if end := strings.Index(section, "\n## "); end >= 0 {
			section = section[:end]
		}
	}
	for file := range assigned {
		if !strings.Contains(section, "]("+file+")") && !strings.Contains(section, "](<"+file+">)") {
			issues = append(issues, found(record.ID+": required file lacks a Details link: "+file, "link the relative file path from plan.md's Details section"))
		}
	}
	base, err := s.Path(path.Dir(record.Path))
	if err != nil {
		return append(issues, found(record.ID+": "+err.Error(), "repair the plan folder path"))
	}
	_ = filepath.WalkDir(base, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			issues = append(issues, found(record.ID+": "+walkErr.Error(), "repair the plan folder"))
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			issues = append(issues, found(record.ID+": symlink in plan folder: "+filePath, "remove the symlink and keep details inside the plan folder"))
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(base, filePath)
		rel = filepath.ToSlash(rel)
		if rel != "plan.md" && entry.Type().IsRegular() && !assigned[rel] {
			issues = append(issues, found(record.ID+": unassigned supporting file: "+rel, "assign it in required_files or remove it from the plan folder"))
		}
		return nil
	})
	return issues
}

func values(m map[string][]string) [][]string {
	var out [][]string
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
