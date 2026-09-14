// Package cow copies local runtime files into independent working copies.
package cow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Stats struct {
	Cloned int `json:"cloned"`
	Copied int `json:"copied"`
	Links  int `json:"links"`
}

func Mode(mode string) error {
	if mode != "auto" && mode != "required" && mode != "copy" && mode != "off" {
		return fmt.Errorf("copy mode must be auto, required, copy, or off: %q", mode)
	}
	return nil
}

// Copy publishes one complete entry without overwriting a destination. Relative
// links stay inside the new checkout; absolute links are translated only when
// they pointed inside the source checkout. External links are never followed.
func Copy(ctx context.Context, source, destination, sourceRoot, targetRoot, mode string) (Stats, error) {
	var stats Stats
	if err := Mode(mode); err != nil {
		return stats, err
	}
	if mode == "off" {
		return stats, nil
	}
	realSourceRoot, err := filepath.EvalSymlinks(sourceRoot)
	if err != nil {
		return stats, err
	}
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return stats, fmt.Errorf("destination already exists: %s", destination)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0700); err != nil {
		return stats, err
	}
	stage, err := os.MkdirTemp(filepath.Dir(destination), ".cc-copy-*")
	if err != nil {
		return stats, err
	}
	defer os.RemoveAll(stage)
	staged := filepath.Join(stage, "entry")
	err = filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if entry.Name() == ".git" {
			return errors.New("cannot copy nested Git metadata")
		}
		dest := filepath.Join(staged, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			resolved := link
			if !filepath.IsAbs(link) {
				resolved = filepath.Join(filepath.Dir(path), link)
			}
			inside, err := filepath.Rel(sourceRoot, resolved)
			if err != nil || inside == ".." || strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
				return fmt.Errorf("external symlink cannot be reused: %s", path)
			}
			// Reject chains that look internal lexically but resolve outside.
			if real, err := filepath.EvalSymlinks(path); err == nil {
				actual, err := filepath.Rel(realSourceRoot, real)
				if err != nil || actual == ".." || strings.HasPrefix(actual, ".."+string(filepath.Separator)) {
					return fmt.Errorf("external symlink chain cannot be reused: %s", path)
				}
			} else if !os.IsNotExist(err) {
				return err
			}
			// Convert even absolute internal links into portable relative links.
			finalPath := filepath.Join(destination, rel)
			link, err = filepath.Rel(filepath.Dir(finalPath), filepath.Join(targetRoot, inside))
			if err != nil {
				return err
			}
			if err := os.Symlink(link, dest); err != nil {
				return err
			}
			stats.Links++
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(dest, 0700)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("cannot reuse special file: %s", path)
		}
		cloned, err := copyFile(path, dest, info.Mode().Perm(), mode)
		if err != nil {
			return err
		}
		if cloned {
			stats.Cloned++
		} else {
			stats.Copied++
		}
		return nil
	})
	if err != nil {
		return stats, err
	}
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return stats, fmt.Errorf("destination appeared during copy: %s", destination)
	}
	return stats, os.Rename(staged, destination)
}

func copyFile(source, destination string, mode fs.FileMode, strategy string) (bool, error) {
	if strategy != "copy" {
		if err := cloneFile(source, destination); err == nil {
			return true, os.Chmod(destination, mode)
		} else if strategy == "required" {
			return false, fmt.Errorf("filesystem CoW unavailable for %s: %w", source, err)
		}
		// A failed native clone may have created an empty/partial destination.
		if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
			return false, err
		}
	}
	in, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return false, err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	return false, errors.Join(copyErr, closeErr)
}
