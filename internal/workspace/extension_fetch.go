package workspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// A remote extension is pinned by commit. The ref is the name a person chose,
// and the commit is what every machine fetches, so a tag moved upstream changes
// nothing until somebody records the move with `extension update`.

// extensionPushDisabled replaces the push URL of every fetched source. Nothing
// here pushes, and Git echoes the value back if anything tries.
const extensionPushDisabled = "READ-ONLY-EXTENSION-SOURCE"

// ErrRefMoved reports a ref that no longer resolves to the recorded commit.
var ErrRefMoved = errors.New("the ref no longer resolves to the recorded commit")

// resolveExtensionRef asks the remote which commit ref names, peeling an
// annotated tag to the commit it points at. A full commit name is its own
// answer.
func (s *Store) resolveExtensionRef(ctx context.Context, url, ref string) (string, error) {
	if err := checkSource(ExtensionSource{URL: url, Ref: ref}); err != nil {
		return "", err
	}
	// ls-remote matches a pattern against whole ref tails, so the peeled
	// entry of an annotated tag is asked for by its own name.
	out, err := Git(ctx, s.Root, "ls-remote", url, ref, ref+"^{}")
	if err != nil {
		return "", fmt.Errorf("could not read %s: %w", url, err)
	}
	refs := map[string]string{}
	for line := range strings.SplitSeq(out, "\n") {
		hash, name, ok := strings.Cut(strings.TrimSpace(line), "\t")
		if ok {
			refs[name] = hash
		}
	}
	candidates := []string{ref}
	if !strings.HasPrefix(ref, "refs/") {
		candidates = []string{"refs/tags/" + ref, "refs/heads/" + ref}
	}
	for _, name := range candidates {
		if hash, ok := refs[name+"^{}"]; ok {
			return hash, nil
		}
		if hash, ok := refs[name]; ok {
			return hash, nil
		}
	}
	if commitPattern.MatchString(ref) {
		return ref, nil
	}
	return "", fmt.Errorf("%s has no branch or tag %q", url, ref)
}

// fetchExtension makes sure the entry's commit is in the cache and returns the
// checkout's workspace-relative path. It tries the commit itself first; a
// server refusing a fetch by commit name is asked for the ref instead, which
// must still name the recorded commit.
func (s *Store) fetchExtension(ctx context.Context, entry ExtensionEntry) (string, error) {
	src := entry.Source
	if err := Name(entry.ID); err != nil {
		return "", err
	}
	if err := checkSource(src); err != nil || !src.Remote() {
		return "", fmt.Errorf("%s is not a remote source that can be fetched", entry.ID)
	}
	if !commitPattern.MatchString(src.Commit) {
		return "", fmt.Errorf("%s records no resolved commit; record one with `extension update --id %s --ref %s`", entry.ID, entry.ID, src.Ref)
	}
	rel := cacheCheckout(entry.ID, src.Commit)
	final, err := s.Path(rel)
	if err != nil {
		return "", err
	}
	if info, err := os.Stat(final); err == nil && info.IsDir() {
		return rel, nil
	}
	parent := filepath.Dir(final)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp(parent, ".fetch-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	steps := [][]string{
		{"init", "--quiet"},
		{"remote", "add", "origin", src.URL},
		{"remote", "set-url", "--push", "origin", extensionPushDisabled},
	}
	for _, step := range steps {
		if _, err := Git(ctx, tmp, step...); err != nil {
			return "", err
		}
	}
	if _, err := Git(ctx, tmp, "fetch", "--quiet", "--depth", "1", "origin", src.Commit); err != nil {
		if _, err := Git(ctx, tmp, "fetch", "--quiet", "--depth", "1", "origin", src.Ref); err != nil {
			return "", fmt.Errorf("could not fetch %s: %w", entry.ID, err)
		}
		fetched, err := Git(ctx, tmp, "rev-parse", "--verify", "FETCH_HEAD^{commit}")
		if err != nil {
			return "", err
		}
		if fetched != src.Commit {
			return "", fmt.Errorf("%s: %w: %s is now %s, and %s is recorded; the recorded commit was kept, and `extension update --id %s --ref %s` records the move", entry.ID, ErrRefMoved, src.Ref, fetched, src.Commit, entry.ID, src.Ref)
		}
	}
	if _, err := Git(ctx, tmp, "checkout", "--quiet", "--detach", src.Commit); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, final); err != nil {
		return "", err
	}
	return rel, nil
}

// ExtensionSync is what sync did for one extension.
type ExtensionSync struct {
	ID     string `json:"id"`
	Commit string `json:"commit,omitempty"`
	// Status is cached, fetched, local, unavailable, or ref-moved.
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// SyncResult is every extension's sync and the cache folders pruned.
type SyncResult struct {
	Extensions []ExtensionSync `json:"extensions"`
	Pruned     []string        `json:"pruned,omitempty"`
}

// SyncExtensions fetches every recorded commit this machine lacks, and reports
// a ref that has moved away from its recorded commit without following it.
// prune removes cached commits no entry records.
func (s *Store) SyncExtensions(ctx context.Context, prune bool) (SyncResult, error) {
	file, present, err := s.Extensions()
	if err != nil {
		return SyncResult{}, err
	}
	result := SyncResult{Extensions: []ExtensionSync{}}
	if !present {
		return result, nil
	}
	keep := map[string]bool{}
	for _, entry := range file.Extensions {
		state := ExtensionSync{ID: entry.ID, Commit: entry.Source.Commit}
		if !entry.Source.Remote() {
			state.Status = "local"
			result.Extensions = append(result.Extensions, state)
			continue
		}
		if commitPattern.MatchString(entry.Source.Commit) {
			keep[cacheCheckout(entry.ID, entry.Source.Commit)] = true
		}
		loaded := s.loadExtension(entry, 0)
		state.Status = "cached"
		if loaded.Fault == extensionUncached {
			if _, err := s.fetchExtension(ctx, entry); err != nil {
				state.Status, state.Detail = "unavailable", err.Error()
				if errors.Is(err, ErrRefMoved) {
					state.Status = "ref-moved"
				}
				result.Extensions = append(result.Extensions, state)
				continue
			}
			state.Status = "fetched"
		}
		if commit, err := s.resolveExtensionRef(ctx, entry.Source.URL, entry.Source.Ref); err == nil && commit != entry.Source.Commit {
			state.Status = "ref-moved"
			state.Detail = fmt.Sprintf("%s now resolves to %s; %s is recorded and still used. Record the move with `extension update --id %s --ref %s` once a person agrees", entry.Source.Ref, commit, entry.Source.Commit, entry.ID, entry.Source.Ref)
		}
		result.Extensions = append(result.Extensions, state)
	}
	if prune {
		pruned, err := s.pruneExtensionCache(keep)
		result.Pruned = pruned
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// pruneExtensionCache removes every cached commit not in keep.
func (s *Store) pruneExtensionCache(keep map[string]bool) ([]string, error) {
	root, err := s.Path(extensionCache)
	if err != nil {
		return nil, err
	}
	ids, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var pruned []string
	for _, id := range ids {
		if !id.IsDir() {
			continue
		}
		commits, err := os.ReadDir(filepath.Join(root, id.Name()))
		if err != nil {
			return pruned, err
		}
		left := 0
		for _, commit := range commits {
			rel := extensionCache + "/" + id.Name() + "/" + commit.Name()
			if keep[rel] {
				left++
				continue
			}
			if err := os.RemoveAll(filepath.Join(root, id.Name(), commit.Name())); err != nil {
				return pruned, err
			}
			pruned = append(pruned, rel)
		}
		if left == 0 {
			_ = os.Remove(filepath.Join(root, id.Name()))
		}
	}
	slices.Sort(pruned)
	return pruned, nil
}
