package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// runDigestLimit is the most entries a run digest records. Past this the
// digest is a single "~truncated" entry whose value is the count.
var runDigestLimit = 20000

// digestSkipDir lists top-level run-root directories that are build
// output, tool cache, or host state. Host state names are the directories
// the adapters in internal/hosts create directly under the run root.
var digestSkipDir = map[string]bool{
	"out":         true,
	"gocache":     true,
	"gopath":      true,
	"tmp":         true,
	"bin":         true,
	"src":         true,
	"setup-home":  true,
	"home":        true,
	"claude":      true,
	"claude-home": true,
	"codex":       true,
	"codex-home":  true,
	"cursor-home": true,
}

// runDigest hashes the run root. A regular file maps to its sha256 hex.
// A directory maps to "dir" and a symlink to "link:<target>". Symlinks are
// not followed, and a .git directory is recorded as "dir" without reading
// its contents. Skipped top-level directories, host-stream, and decoy.log
// are absent.
func runDigest(root string) (map[string]string, error) {
	out := map[string]string{}
	n := 0
	record := func(rel, value string) {
		n++
		if n <= runDigestLimit {
			out[rel] = value
		}
	}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if digestSkipped(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == ".git" {
			record(rel, "dir")
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			record(rel, "link:"+target)
			return nil
		}
		if d.IsDir() {
			record(rel, "dir")
			return nil
		}
		if !d.Type().IsRegular() {
			record(rel, "dir")
			return nil
		}
		sum, err := hashFile(path)
		if err != nil {
			return err
		}
		record(rel, sum)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if n > runDigestLimit {
		return map[string]string{"~truncated": strconv.Itoa(n)}, nil
	}
	return out, nil
}

func digestSkipped(rel string) bool {
	top, _, _ := strings.Cut(rel, "/")
	if digestSkipDir[top] {
		return true
	}
	return rel == "host-stream" || rel == "decoy.log"
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sum := sha256.New()
	if _, err := io.Copy(sum, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(sum.Sum(nil)), nil
}
