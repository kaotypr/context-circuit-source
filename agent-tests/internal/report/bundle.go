package report

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// BundleOptions selects the run directory and the paths to clean.
type BundleOptions struct {
	RunDir   string
	Dest     string
	Home     string
	User     string
	RunRoots []string
}

type zipEntry struct {
	name string
	body []byte
}

// WriteBundle writes a cleaned zip of the run.
// A key-shaped string in anything that would be included is an error,
// and no archive file is left behind.
func WriteBundle(opts BundleOptions) error {
	if opts.Dest == "" {
		opts.Dest = filepath.Join(opts.RunDir, "bundle.zip")
	}
	entries, err := collectEntries(opts)
	if err != nil {
		if secretErr(err) {
			_ = os.Remove(opts.Dest)
		}
		return err
	}
	if err := writeZip(opts.Dest, entries); err != nil {
		_ = os.Remove(opts.Dest)
		return err
	}
	return nil
}

func secretErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "key-shaped string")
}

func collectEntries(opts BundleOptions) ([]zipEntry, error) {
	var entries []zipEntry
	err := filepath.WalkDir(opts.RunDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == opts.RunDir {
			return nil
		}
		rel, err := filepath.Rel(opts.RunDir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if skipAuthDir(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		zipName, ok := safeZipName(rel)
		if !ok || droppedFile(zipName) || !includeRel(zipName) {
			return nil
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		text := cleanPaths(string(body), opts.Home, opts.User, opts.RunRoots)
		if record.LooksLikeSecret(zipName) || record.LooksLikeSecret(text) {
			return fmt.Errorf("bundle: key-shaped string in %s", zipName)
		}
		entries = append(entries, zipEntry{name: zipName, body: []byte(text)})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	return entries, nil
}

func includeRel(rel string) bool {
	base := path.Base(rel)
	switch base {
	case "runs.jsonl", "reply.md", "reply", "reads.txt", "check.json", "check-before.json":
		return true
	}
	if base == "stream.jsonl" || (strings.HasPrefix(base, "stream") && strings.HasSuffix(base, ".jsonl")) {
		return true
	}
	return strings.HasPrefix(rel, "artifacts/") || strings.Contains(rel, "/artifacts/")
}

func isAuthDir(name string) bool {
	switch strings.ToLower(name) {
	case "claude", "codex", "cursor-home", ".claude", ".codex", ".cursor":
		return true
	default:
		return false
	}
}

// skipAuthDir drops a host config directory. runs/<scenario>/<host> keeps
// the host id claude or codex; that path element is the run layout, not
// the config directory the host writes under the run root.
func skipAuthDir(rel string) bool {
	base := path.Base(rel)
	if !isAuthDir(base) {
		return false
	}
	parts := strings.Split(rel, "/")
	if len(parts) == 3 && parts[0] == "runs" && (strings.EqualFold(base, "claude") || strings.EqualFold(base, "codex")) {
		return false
	}
	return true
}

func droppedFile(name string) bool {
	base := strings.ToLower(path.Base(name))
	switch base {
	case "auth.json", ".netrc", ".env":
		return true
	}
	if strings.HasPrefix(base, ".credentials") || strings.HasPrefix(base, "id_rsa") {
		return true
	}
	return strings.HasSuffix(base, ".token") || strings.HasSuffix(base, ".pem")
}

func safeZipName(rel string) (string, bool) {
	if rel == "" || strings.HasPrefix(rel, "/") || strings.Contains(rel, `\`) {
		return "", false
	}
	clean := path.Clean(rel)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || path.IsAbs(clean) {
		return "", false
	}
	return clean, true
}

func writeZip(dest string, entries []zipEntry) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".bundle-*.zip.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	done := false
	defer func() {
		if !done {
			_ = tmp.Close()
			_ = os.Remove(tmpName)
		}
	}()
	zw := zip.NewWriter(tmp)
	stamp := time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC)
	for _, entry := range entries {
		hdr := &zip.FileHeader{Name: entry.name, Method: zip.Deflate}
		hdr.SetModTime(stamp)
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			_ = zw.Close()
			return err
		}
		if _, err := w.Write(entry.body); err != nil {
			_ = zw.Close()
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(tmpName, dest); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	done = true
	return nil
}

// readZip is used by tests in this package. It is also a small helper
// so a caller can confirm an entry without reimplementing the reader.
func readZip(path string) (map[string]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	out := map[string]string{}
	for _, file := range r.File {
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return nil, err
		}
		out[file.Name] = string(body)
	}
	return out, nil
}
