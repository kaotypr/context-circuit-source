package workspace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kaotypr/context-circuit-source/internal/cow"
)

type ReuseOptions struct {
	Mode  string
	Paths []string
}

func (options ReuseOptions) Validate() error {
	if options.Mode != "" {
		if err := cow.Mode(options.Mode); err != nil {
			return err
		}
	}
	for _, path := range options.Paths {
		path = filepath.ToSlash(path)
		if path == "." || !filepath.IsLocal(path) || strings.Contains(path, "\\") {
			return fmt.Errorf("reuse path must be repository-relative: %s", path)
		}
		for _, part := range strings.Split(path, "/") {
			if part == ".git" || part == ".context-circuit" {
				return fmt.Errorf("cannot reuse control files: %s", path)
			}
		}
	}
	return nil
}

type ReuseEntry struct {
	Path   string    `json:"path"`
	Status string    `json:"status"`
	Reason string    `json:"reason,omitempty"`
	Stats  cow.Stats `json:"files"`
}

// ReuseEnvironment selects only ignored runtime entries, never arbitrary dirty
// source files. Contents of .env files are copied opaquely and never reported.
func ReuseEnvironment(ctx context.Context, source, destination string, options ReuseOptions) ([]ReuseEntry, error) {
	if options.Mode == "" {
		options.Mode = "auto"
	}
	if err := options.Validate(); err != nil {
		return nil, err
	}
	if options.Mode == "off" {
		return nil, nil
	}
	ignored, err := Git(ctx, source, "ls-files", "--others", "--ignored", "--exclude-standard", "--directory", "-z")
	if err != nil {
		return nil, err
	}
	selected := map[string]bool{}
	for _, path := range strings.Split(ignored, "\x00") {
		path = strings.TrimSuffix(path, "/")
		base := filepath.Base(path)
		if base == "node_modules" || base == ".env" || strings.HasPrefix(base, ".env.") {
			selected[path] = true
		}
	}
	for _, path := range options.Paths {
		path = filepath.ToSlash(path)
		selected[path] = true
	}
	paths := make([]string, 0, len(selected))
	for path := range selected {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	srcStore, err := Open(source)
	if err != nil {
		return nil, err
	}
	dstStore, err := Open(destination)
	if err != nil {
		return nil, err
	}
	var results []ReuseEntry
	deps := &dependencyCheck{}
	for _, rel := range paths {
		if err := ctx.Err(); err != nil {
			return results, err
		}
		entry := ReuseEntry{Path: rel, Status: "skipped"}
		// Path checks reject symbolic-link parents and roots before any copying.
		src, err := srcStore.Path(rel)
		if err != nil {
			return results, err
		}
		dst, err := dstStore.Path(rel)
		if err != nil {
			return results, err
		}
		entry.Reason, err = reuseSkipReason(ctx, source, destination, rel, src, dst, deps)
		if err != nil {
			return results, err
		}
		if entry.Reason != "" {
			results = append(results, entry)
			continue
		}
		entry.Stats, err = cow.Copy(ctx, src, dst, source, destination, options.Mode)
		// Reuse is a convenience, so an entry that cannot be copied whole,
		// such as node_modules holding an `npm link`, is reported rather
		// than failing the prepare. Only an entry refused by design is
		// skipped: a full disk or a permission error stays fatal. Copy
		// stages and renames, so nothing partial is left behind.
		if err != nil && options.Mode != "required" && errors.Is(err, cow.ErrNotReusable) {
			entry.Stats, entry.Reason = cow.Stats{}, "not copied: "+err.Error()
			results = append(results, entry)
			continue
		}
		if err != nil {
			return results, fmt.Errorf("reuse %s failed; existing worktree retained: %w", rel, err)
		}
		if err := restrictEnvMode(dst, rel); err != nil {
			return results, err
		}
		entry.Status = "reused"
		results = append(results, entry)
	}
	return results, nil
}

// reuseSkipReason returns why an entry must not be copied, or "" when it may be.
func reuseSkipReason(ctx context.Context, source, destination, rel, src, dst string, deps *dependencyCheck) (string, error) {
	ignorePath := rel
	if info, err := os.Stat(src); err == nil && info.IsDir() {
		ignorePath += "/"
	}
	if _, err := os.Lstat(dst); !os.IsNotExist(err) {
		if err != nil {
			return "", err
		}
		return "destination already exists; preserved", nil
	}
	if _, err := os.Lstat(src); os.IsNotExist(err) {
		return "source entry is absent", nil
	} else if err != nil {
		return "", err
	}
	if _, err := Git(ctx, source, "check-ignore", "-q", "--", ignorePath); err != nil {
		return "source path is not ignored by Git", nil
	}
	if _, err := Git(ctx, destination, "check-ignore", "-q", "--", ignorePath); err != nil {
		return "target path is not ignored by Git", nil
	}
	if filepath.Base(rel) == "node_modules" {
		return deps.reason(ctx, source, destination)
	}
	return "", nil
}

// dependencyCheck compares dependency inputs once per reuse run.
type dependencyCheck struct {
	checked bool
	value   string
}

func (d *dependencyCheck) reason(ctx context.Context, source, destination string) (string, error) {
	if d.checked {
		return d.value, nil
	}
	value, err := compareDependencyInputs(ctx, source, destination)
	if err != nil {
		return "", err
	}
	d.checked, d.value = true, value
	return value, nil
}

// restrictEnvMode keeps copied .env files private to the owner.
func restrictEnvMode(dst, rel string) error {
	base := filepath.Base(rel)
	if base != ".env" && !strings.HasPrefix(base, ".env.") {
		return nil
	}
	if info, err := os.Stat(dst); err == nil && info.Mode().IsRegular() {
		return os.Chmod(dst, 0600)
	}
	return nil
}

func compareDependencyInputs(ctx context.Context, source, destination string) (string, error) {
	names := map[string]bool{}
	for _, root := range []string{source, destination} {
		list, err := Git(ctx, root, "ls-files", "--cached", "-z")
		if err != nil {
			return "", err
		}
		for _, name := range strings.Split(list, "\x00") {
			switch filepath.Base(name) {
			case "package.json", "package-lock.json", "npm-shrinkwrap.json", "yarn.lock", "pnpm-lock.yaml", "pnpm-workspace.yaml", "bun.lock", "bun.lockb", ".nvmrc", ".node-version":
				names[name] = true
			}
		}
	}
	leftStore, err := Open(source)
	if err != nil {
		return "", err
	}
	rightStore, err := Open(destination)
	if err != nil {
		return "", err
	}
	for name := range names {
		left, leftErr := leftStore.Read(name)
		right, rightErr := rightStore.Read(name)
		if os.IsNotExist(leftErr) || os.IsNotExist(rightErr) {
			return "dependency inputs differ; prepare dependencies for this branch", nil
		}
		if leftErr != nil {
			return "", leftErr
		}
		if rightErr != nil {
			return "", rightErr
		}
		if !bytes.Equal(left, right) {
			return "dependency inputs differ; prepare dependencies for this branch", nil
		}
	}
	return "", nil
}
