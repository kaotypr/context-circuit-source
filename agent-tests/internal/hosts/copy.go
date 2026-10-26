package hosts

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// FileCopy is one regular file the runner places in the isolated run root
// before the host starts. Mode is the permission of the destination.
type FileCopy struct {
	Src  string
	Dst  string
	Mode os.FileMode
}

// ApplyCopies copies each file and returns secret strings found in the
// destination bytes. The source path is not included in errors.
func ApplyCopies(copies []FileCopy) ([]string, error) {
	var secrets []string
	for _, copy := range copies {
		body, err := copyOne(copy)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secretsFromLogin(body)...)
	}
	return uniqueStrings(secrets), nil
}

func copyOne(copy FileCopy) ([]byte, error) {
	mode := copy.Mode
	if mode == 0 {
		mode = 0o600
	}
	body, err := readRegular(copy.Src)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(copy.Dst), 0o755); err != nil {
		return nil, errCopy
	}
	copied, err := writeExclusive(copy.Dst, body, mode)
	if err != nil {
		return nil, err
	}
	return copied, nil
}

var (
	errLoginMissing = errors.New("subscription login file: not found")
	errLoginSymlink = errors.New("subscription login file is a symlink")
	errLoginKind    = errors.New("subscription login file is not a regular file")
	errLoginRead    = errors.New("subscription login file: unavailable")
	errCopy         = errors.New("subscription login file: could not be copied")
)

func readRegular(src string) ([]byte, error) {
	info, err := os.Lstat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errLoginMissing
		}
		return nil, errLoginRead
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, errLoginSymlink
	}
	if !info.Mode().IsRegular() {
		return nil, errLoginKind
	}
	in, err := os.Open(src)
	if err != nil {
		return nil, errLoginRead
	}
	body, readErr := io.ReadAll(in)
	closeErr := in.Close()
	if readErr != nil || closeErr != nil {
		return nil, errLoginRead
	}
	return body, nil
}

func writeExclusive(dst string, body []byte, mode os.FileMode) ([]byte, error) {
	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return nil, errCopy
	}
	_, werr := out.Write(body)
	if werr != nil {
		out.Close()
		os.Remove(dst)
		return nil, errCopy
	}
	if _, err := out.Seek(0, io.SeekStart); err != nil {
		out.Close()
		os.Remove(dst)
		return nil, errCopy
	}
	copied, rerr := io.ReadAll(out)
	cerr := out.Close()
	if rerr != nil || cerr != nil {
		os.Remove(dst)
		return nil, errCopy
	}
	st, err := os.Lstat(dst)
	if err != nil || st.Mode()&os.ModeSymlink != 0 || !st.Mode().IsRegular() || st.Mode().Perm() != mode.Perm() {
		os.Remove(dst)
		return nil, errCopy
	}
	return copied, nil
}

// secretsFromLogin reads secret strings from the copied bytes.
// A JSON document contributes every string value of length 16 or more.
// A document that does not parse contributes JWT and refresh-token shapes.
func secretsFromLogin(b []byte) []string {
	var doc any
	if err := json.Unmarshal(b, &doc); err != nil {
		return record.ShapeSecrets(string(b))
	}
	var found []string
	walkJSONStrings(doc, func(s string) {
		if len(s) >= 16 {
			found = append(found, s)
		}
	})
	return found
}

func walkJSONStrings(v any, fn func(string)) {
	switch t := v.(type) {
	case map[string]any:
		for _, child := range t {
			walkJSONStrings(child, fn)
		}
	case []any:
		for _, child := range t {
			walkJSONStrings(child, fn)
		}
	case string:
		fn(t)
	}
}

func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
