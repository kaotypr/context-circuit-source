package world

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

// ExtractArchive reads git archive of ref from checkout and unpacks it under
// dest. The checkout is not written. The child environment is PATH plus the
// two Git settings that keep the user's config out of the command.
func ExtractArchive(ctx context.Context, checkout, ref, dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	tarFile, err := os.CreateTemp(filepath.Dir(dest), "archive-*.tar")
	if err != nil {
		return err
	}
	tarName := tarFile.Name()
	tarFile.Close()
	defer os.Remove(tarName)

	res, err := proc.Run(ctx, proc.Spec{
		Argv:       []string{"git", "-C", checkout, "archive", "--format=tar", ref},
		Env:        archiveEnv(),
		StreamPath: tarName,
	})
	if err != nil {
		return fmt.Errorf("git archive: %w\n%s", err, res.Stderr)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("git archive: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	f, err := os.Open(tarName)
	if err != nil {
		return err
	}
	defer f.Close()
	return extractTar(f, dest)
}

func archiveEnv() []string {
	return []string{
		"PATH=" + os.Getenv("PATH"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
	}
}

func extractTar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag == tar.TypeXGlobalHeader || hdr.Typeflag == tar.TypeXHeader {
			continue
		}
		target, err := safeJoin(dest, hdr.Name)
		if err != nil {
			return err
		}
		if err := ensureWithin(dest, target); err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := writeTarFile(target, tr, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeSymlink, tar.TypeLink:
			if linkEscapes(dest, target, hdr.Linkname) {
				return fmt.Errorf("tar link escapes dest: %s -> %s", hdr.Name, hdr.Linkname)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if hdr.Typeflag == tar.TypeSymlink {
				if err := os.Symlink(hdr.Linkname, target); err != nil {
					return err
				}
				continue
			}
			linkTarget, err := safeJoin(dest, hdr.Linkname)
			if err != nil {
				return err
			}
			if err := os.Link(linkTarget, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("tar entry %s has unsupported type %q", hdr.Name, hdr.Typeflag)
		}
	}
}

func writeTarFile(target string, r io.Reader, mode os.FileMode) error {
	f, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode.Perm())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, r)
	closeErr := f.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func safeJoin(dest, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty tar path")
	}
	slash := strings.ReplaceAll(name, "\\", "/")
	if strings.HasPrefix(slash, "/") || (len(slash) >= 2 && slash[1] == ':') {
		return "", fmt.Errorf("absolute tar path %q", name)
	}
	clean := path.Clean(slash)
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("tar path escapes dest: %q", name)
	}
	if clean == "." {
		return dest, nil
	}
	target := filepath.Join(dest, filepath.FromSlash(clean))
	if !insideOrEqual(dest, target) {
		return "", fmt.Errorf("tar path escapes dest: %q", name)
	}
	return target, nil
}

func linkEscapes(dest, linkPath, linkname string) bool {
	if linkname == "" {
		return true
	}
	slash := strings.ReplaceAll(linkname, "\\", "/")
	if strings.HasPrefix(slash, "/") || (len(slash) >= 2 && slash[1] == ':') {
		return true
	}
	clean := path.Clean(slash)
	if clean == ".." || strings.HasPrefix(clean, "../") {
		resolved := filepath.Join(filepath.Dir(linkPath), filepath.FromSlash(clean))
		return !insideOrEqual(dest, resolved)
	}
	resolved := filepath.Join(filepath.Dir(linkPath), filepath.FromSlash(clean))
	return !insideOrEqual(dest, resolved)
}

func ensureWithin(dest, target string) error {
	rel, err := filepath.Rel(dest, target)
	if err != nil {
		return err
	}
	cur := dest
	for _, part := range strings.Split(rel, string(os.PathSeparator)) {
		if part == "" || part == "." {
			continue
		}
		cur = filepath.Join(cur, part)
		info, err := os.Lstat(cur)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		raw, err := os.Readlink(cur)
		if err != nil {
			return err
		}
		resolved := raw
		if !filepath.IsAbs(raw) {
			resolved = filepath.Join(filepath.Dir(cur), raw)
		}
		if !insideOrEqual(dest, resolved) {
			return fmt.Errorf("symlink %s escapes dest", cur)
		}
	}
	return nil
}

func insideOrEqual(parent, child string) bool {
	if parent == child {
		return true
	}
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}
