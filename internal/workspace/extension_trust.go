package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Whether this workspace uses an extension is decided in extensions.yaml and
// reviewed like any commit. Whether its scripts may run on a machine is decided
// by the member using that machine, here: a script runs only after that member
// ran `extension allow` for the extension's exact version, which is the commit
// of a remote source and a hash of a folder source's files. Any change to
// either asks again. Nothing in this CLI allows anything by itself, and no flag
// or setting skips the allow.

// ExtensionTrustPath is this machine's record of which extension versions its
// members allowed. It sits under .context-circuit/local/, which every workspace
// ignores, because an allow speaks for one member on one machine only.
const ExtensionTrustPath = ".context-circuit/local/extension-trust.yaml"

// ExtensionAllow is one member's allow of one extension version on one machine.
type ExtensionAllow struct {
	Member    string  `yaml:"member" json:"member"`
	Machine   string  `yaml:"machine" json:"machine"`
	Extension string  `yaml:"extension" json:"extension"`
	Version   string  `yaml:"version" json:"version"`
	AllowedAt ISOTime `yaml:"allowed_at" json:"allowed_at"`
}

// ExtensionTrust is the trust file.
type ExtensionTrust struct {
	Allows []ExtensionAllow `yaml:"allows" json:"allows"`
}

// Trust states an extension with scripts can be in for this member and machine.
const (
	TrustAllowed    = "allowed"
	TrustNotAllowed = "not allowed"
	TrustChanged    = "changed since allowed"
)

// machineName names this machine in an allow. The trust file is already local
// to a machine; naming it as well keeps a copied or shared folder from carrying
// one machine's allows to another.
func machineName() string {
	name, err := os.Hostname()
	if err != nil || strings.TrimSpace(name) == "" {
		return "unknown"
	}
	return name
}

func (s *Store) extensionTrust() (ExtensionTrust, error) {
	trust := ExtensionTrust{Allows: []ExtensionAllow{}}
	err := s.YAML(ExtensionTrustPath, &trust)
	if os.IsNotExist(err) {
		return ExtensionTrust{Allows: []ExtensionAllow{}}, nil
	}
	if err != nil {
		return trust, fmt.Errorf("%s: %w", ExtensionTrustPath, err)
	}
	return trust, nil
}

// FolderHash is a folder source's version: a SHA-256 over every file's
// slash-separated relative path and contents, in path order, excluding
// nothing. Line endings are normalized to LF first, so a checkout that
// converts them on Windows hashes as the same version. A symlink is hashed by
// what it points to, not followed.
func FolderHash(folder string) (string, error) {
	type entry struct{ path, digest string }
	var entries []entry
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}
		var content []byte
		switch {
		case d.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			content = []byte("symlink:" + filepath.ToSlash(target))
		case d.Type().IsRegular():
			if content, err = os.ReadFile(path); err != nil {
				return err
			}
			content = []byte(strings.ReplaceAll(string(content), "\r\n", "\n"))
		default:
			return fmt.Errorf("%s is neither a file, a folder, nor a symlink", filepath.ToSlash(relative))
		}
		sum := sha256.Sum256(content)
		entries = append(entries, entry{filepath.ToSlash(relative), hex.EncodeToString(sum[:])})
		return nil
	})
	if err != nil {
		return "", err
	}
	slices.SortFunc(entries, func(a, b entry) int { return strings.Compare(a.path, b.path) })
	hash := sha256.New()
	for _, e := range entries {
		fmt.Fprintf(hash, "%s\x00%s\n", e.path, e.digest)
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

// extensionVersion is what an allow pins: the commit of a remote source, or
// the hash of a folder source's files as they are now.
func (s *Store) extensionVersion(loaded loadedExtension) (string, error) {
	if loaded.Entry.Source.Remote() {
		return loaded.Entry.Source.Commit, nil
	}
	folder, err := s.Path(loaded.Folder)
	if err != nil {
		return "", err
	}
	return FolderHash(folder)
}

// trustState says whether member allowed this version of an extension on this
// machine, allowed another version, or never allowed it.
func (s *Store) trustState(member, id, version string) (string, error) {
	trust, err := s.extensionTrust()
	if err != nil {
		return "", err
	}
	state := TrustNotAllowed
	machine := machineName()
	for _, allow := range trust.Allows {
		if allow.Member != member || allow.Machine != machine || allow.Extension != id {
			continue
		}
		if allow.Version == version {
			return TrustAllowed, nil
		}
		state = TrustChanged
	}
	return state, nil
}

// notAllowedReason is the skip a script not allowed here reports, naming the
// command that clears it and who decides to run it.
func notAllowedReason(state, id string) string {
	prefix := "not allowed on this machine"
	if state == TrustChanged {
		prefix = "changed since allowed on this machine"
	}
	return prefix + ": the person decides whether to run `extension allow --id " + id + "`, which shows what the scripts do before allowing them"
}

// ScriptHookSummary is one script hook as `extension allow` shows it.
type ScriptHookSummary struct {
	Point      string   `json:"point"`
	Run        []string `json:"run"`
	RunWindows []string `json:"run_windows,omitempty"`
	OnFailure  string   `json:"on_failure"`
	Consent    string   `json:"consent"`
	Action     string   `json:"action,omitempty"`
	Timeout    string   `json:"timeout"`
}

// AllowedExtension is what `extension allow` shows and records.
type AllowedExtension struct {
	ID       string              `json:"id"`
	Source   ExtensionSource     `json:"source"`
	Version  string              `json:"version"`
	Member   string              `json:"member"`
	Machine  string              `json:"machine"`
	Scripts  []ScriptHookSummary `json:"scripts"`
	Previous string              `json:"previous_version,omitempty"`
	Note     string              `json:"note"`
}

// AllowExtension records that the active member allows the current version of
// an extension's scripts on this machine, and returns what was allowed: the
// source, the version, and each script hook with its command, its on_failure,
// and the action a standing or each-time consent covers. It runs nothing.
func (s *Store) AllowExtension(id string) (AllowedExtension, error) {
	member, err := s.ActiveMember()
	if err != nil {
		return AllowedExtension{}, err
	}
	if member == "" {
		return AllowedExtension{}, errors.New("select the member using this machine with `member use` first: an allow belongs to one member")
	}
	file, present, err := s.Extensions()
	if err != nil {
		return AllowedExtension{}, err
	}
	index := file.extensionIndex(id)
	if !present || index < 0 {
		return AllowedExtension{}, fmt.Errorf("extensions.yaml lists no extension %s", id)
	}
	loaded := s.loadExtension(file.Extensions[index], index)
	if loaded.Fault != extensionUsable {
		return AllowedExtension{}, fmt.Errorf("%s cannot be read on this machine, so nothing was allowed: %s; %s", id, loaded.Problem, extensionRemedy(loaded))
	}
	allowed := AllowedExtension{ID: id, Source: loaded.Entry.Source, Member: member, Machine: machineName(), Scripts: []ScriptHookSummary{}}
	for _, hook := range loaded.Manifest.Hooks {
		if hook.Script == nil {
			continue
		}
		point, _ := ParseHookPoint(hook.At)
		consent := hook.Script.Consent
		if consent == "" {
			consent = "none"
		}
		allowed.Scripts = append(allowed.Scripts, ScriptHookSummary{
			Point: point.String(), Run: hook.Script.Run, RunWindows: hook.Script.RunWindows,
			OnFailure: hook.Script.OnFailure, Consent: consent, Action: hook.Script.Action,
			Timeout: scriptTimeout(*hook.Script).String(),
		})
	}
	if len(allowed.Scripts) == 0 {
		return AllowedExtension{}, fmt.Errorf("%s declares no scripts, so there is nothing to allow", id)
	}
	if allowed.Version, err = s.extensionVersion(loaded); err != nil {
		return AllowedExtension{}, err
	}
	trust, err := s.extensionTrust()
	if err != nil {
		return AllowedExtension{}, err
	}
	kept := []ExtensionAllow{}
	for _, allow := range trust.Allows {
		if allow.Member == member && allow.Machine == allowed.Machine && allow.Extension == id {
			if allow.Version != allowed.Version {
				allowed.Previous = allow.Version
			}
			continue
		}
		kept = append(kept, allow)
	}
	trust.Allows = append(kept, ExtensionAllow{Member: member, Machine: allowed.Machine, Extension: id, Version: allowed.Version, AllowedAt: ISOTime(isoTime(time.Now()))})
	if err := s.WriteYAML(ExtensionTrustPath, trust, 0600); err != nil {
		return AllowedExtension{}, err
	}
	allowed.Note = "allowed for " + member + " on this machine at this version; nothing was run. These scripts now run when their points fire, each-time ones only after the person's yes; any change to the extension asks again. A standing consent covers only the action named for its hook."
	return allowed, nil
}
