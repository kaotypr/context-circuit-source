package workspace

import (
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

// A workspace's rules are of two kinds. Fixed rules hold in every workspace
// whatever an extension says: approving an intent, authorizing anything that
// leaves the workspace, starting planning and execution, leaving the plan
// record alone, no attribution, no secrets, borrowed knowledge read-only. The
// overridable defaults are a closed list of four, and an extension changes one
// only where its manifest declares the override and a person allowed it in this
// workspace with `extension override`. The CLI is where that is enforced: an
// override reaches the agent only through `extension hooks`, and a dispatch that
// needs one is accepted only when it names an allowed one.

// The overridable defaults.
const (
	OverrideRepairLoop = "repair-loop"
	OverrideVerify     = "verify-during-execution"
	OverrideRiskReview = "risk-review"
	OverrideDelegation = "delegation"
)

// OverridableDefaults is the closed list, in the order they are documented.
var OverridableDefaults = []string{OverrideRepairLoop, OverrideVerify, OverrideRiskReview, OverrideDelegation}

// maxRounds caps a repair loop, so no manifest can ask for an unbounded one.
const maxRounds = 10

// delegationStages are the stages delegation may name.
var delegationStages = []string{"plan", "implement", "review"}

// ManifestOverride is one overridable default a manifest asks to change, with
// the fields that default takes.
type ManifestOverride struct {
	Default string `yaml:"default" json:"default"`
	// Instructions say what counts as clean and what to fix first for a loop
	// or a verification, and when a change warrants review for risk-review.
	Instructions string `yaml:"instructions,omitempty" json:"instructions,omitempty"`
	// Verifier is the read-only role a repair loop or a verification runs: a
	// role this extension declares, or the shipped reviewer.
	Verifier string `yaml:"verifier,omitempty" json:"verifier,omitempty"`
	// Reviewer is the role risk-review runs, the shipped reviewer by default.
	Reviewer string `yaml:"reviewer,omitempty" json:"reviewer,omitempty"`
	// Rounds is the most verifications a repair loop asks for, 1 to 10.
	Rounds int `yaml:"rounds,omitempty" json:"rounds,omitempty"`
	// Stage and Role are the stage delegation always hands to the role.
	Stage string `yaml:"stage,omitempty" json:"stage,omitempty"`
	Role  string `yaml:"role,omitempty" json:"role,omitempty"`
}

// OverrideAllowance is a person's allowance of one override in extensions.yaml.
// Rounds may lower the manifest's rounds, never raise them.
type OverrideAllowance struct {
	Rounds int `yaml:"rounds,omitempty" json:"rounds,omitempty"`
}

// fixedRulesNote is said wherever something outside the closed list is named.
const fixedRulesNote = "approving an intent, authorizing push, pull requests, merges, deployment and deletion, starting planning and execution, leaving the plan record alone, no attribution, no secrets, and borrowed knowledge read-only are fixed whatever an extension says"

func overridableList() string { return strings.Join(OverridableDefaults, ", ") }

// validateManifestOverrides checks the overrides a manifest declares. Their
// instruction files are checked by loadExtension, which has the folder.
func validateManifestOverrides(m Manifest) error {
	seen := map[string]bool{}
	for _, o := range m.Overrides {
		if !slices.Contains(OverridableDefaults, o.Default) {
			return fmt.Errorf("override %q is not an overridable default; the overridable defaults are %s, and %s", o.Default, overridableList(), fixedRulesNote)
		}
		if seen[o.Default] {
			return fmt.Errorf("override %s appears twice; each default appears once", o.Default)
		}
		seen[o.Default] = true
		label := "override " + o.Default
		unexpected := func(field string, set bool) error {
			if set {
				return fmt.Errorf("%s takes no %s", label, field)
			}
			return nil
		}
		var errs []error
		switch o.Default {
		case OverrideRepairLoop, OverrideVerify:
			if o.Verifier == "" {
				return fmt.Errorf("%s names no verifier", label)
			}
			if err := verifierRole(m, o.Verifier); err != nil {
				return fmt.Errorf("%s: %w", label, err)
			}
			if o.Default == OverrideRepairLoop && (o.Rounds < 1 || o.Rounds > maxRounds) {
				return fmt.Errorf("%s: rounds is between 1 and %d, not %d", label, maxRounds, o.Rounds)
			}
			errs = append(errs, unexpected("reviewer", o.Reviewer != ""), unexpected("stage", o.Stage != ""), unexpected("role", o.Role != ""))
			if o.Default == OverrideVerify {
				errs = append(errs, unexpected("rounds; it verifies once", o.Rounds != 0))
			}
		case OverrideRiskReview:
			if o.Instructions == "" {
				return fmt.Errorf("%s names no instructions saying when a change warrants review", label)
			}
			if o.Reviewer != "" {
				if err := verifierRole(m, o.Reviewer); err != nil {
					return fmt.Errorf("%s: %w", label, err)
				}
			}
			errs = append(errs, unexpected("verifier", o.Verifier != ""), unexpected("rounds", o.Rounds != 0), unexpected("stage", o.Stage != ""), unexpected("role", o.Role != ""))
		case OverrideDelegation:
			if !slices.Contains(delegationStages, o.Stage) {
				return fmt.Errorf("%s: stage is %s, not %q", label, strings.Join(delegationStages, ", "), o.Stage)
			}
			if err := delegationRole(m, o.Stage, o.Role); err != nil {
				return fmt.Errorf("%s: %w", label, err)
			}
			errs = append(errs, unexpected("instructions", o.Instructions != ""), unexpected("verifier", o.Verifier != ""), unexpected("reviewer", o.Reviewer != ""), unexpected("rounds", o.Rounds != 0))
		}
		if err := errors.Join(errs...); err != nil {
			return err
		}
	}
	return nil
}

// verifierRole accepts the shipped reviewer, or a role the manifest declares
// with read-only access: the verifier must not be able to change the code it
// checks.
func verifierRole(m Manifest, name string) error {
	if name == "reviewer" {
		return nil
	}
	for _, role := range m.Roles {
		if role.ID == name {
			if role.Access != AccessReadOnly {
				return fmt.Errorf("verifier %s has %s access; a verifier is read-only: the shipped reviewer, or a role this extension declares with access read-only", name, role.Access)
			}
			return nil
		}
	}
	return fmt.Errorf("verifier %s is neither the shipped reviewer nor a role this extension declares", name)
}

// delegationRole accepts a shipped role or one the manifest declares, with the
// access its stage needs: implementing edits, planning and reviewing read.
func delegationRole(m Manifest, stage, name string) error {
	if name == "" {
		return errors.New("names no role")
	}
	readOnly, known := false, false
	if shipped, ok := Roles[name]; ok {
		readOnly, known = shipped.ReadOnly, true
		if name == "explorer" {
			return errors.New("the explorer does no stage; name the planner, worker, reviewer, or a role this extension declares")
		}
	}
	for _, role := range m.Roles {
		if role.ID == name {
			readOnly, known = role.Access != AccessEdit, true
		}
	}
	if !known {
		return fmt.Errorf("role %s is neither a shipped role nor one this extension declares", name)
	}
	if (stage == "implement") == readOnly {
		need := "edit"
		if stage != "implement" {
			need = "read-only"
		}
		return fmt.Errorf("the %s stage needs a role with %s access, and %s is not one", stage, need, name)
	}
	return nil
}

// manifestOverrideFiles are the instruction files a manifest's overrides name.
func manifestOverrideFiles(m Manifest) []roleFileRef {
	var files []roleFileRef
	for _, o := range m.Overrides {
		if o.Instructions != "" {
			files = append(files, roleFileRef{"override " + o.Default, o.Instructions})
		}
	}
	return files
}

// overrideRole is the name an override's role is dispatched as: a shipped
// role's own, or an extension role's x- name.
func overrideRole(extension, name string) string {
	if name == "" {
		return ""
	}
	if _, ok := Roles[name]; ok {
		return name
	}
	return extensionRoleName(extension, name)
}

// overridePoint is the hook point that returns an override.
func overridePoint(o ManifestOverride) HookPoint {
	if o.Default == OverrideDelegation {
		return HookPoint{"before", o.Stage}
	}
	return HookPoint{"after", "implement"}
}

// HookOverride is an allowed override as the coordinator receives it: which
// default, from which extension, and what it runs.
type HookOverride struct {
	Default   string `json:"default"`
	Extension string `json:"extension"`
	// Rounds is the round limit of a repair loop.
	Rounds int `json:"rounds,omitempty"`
	// Verifier is the role a loop or a verification dispatches, Reviewer the
	// role risk-review dispatches, each by the name `agent dispatch` takes.
	Verifier string `json:"verifier,omitempty"`
	Reviewer string `json:"reviewer,omitempty"`
	// Stage and Role are delegation's: the stage always goes to the role.
	Stage string `json:"stage,omitempty"`
	Role  string `json:"role,omitempty"`
	// Instructions is the workspace-relative instruction file, if any.
	Instructions string `json:"instructions,omitempty"`
}

// allowedOverride resolves one declared override against the entry's
// allowance, or reports false when the person has not allowed it. Rounds above
// the manifest's are capped, because the manifest says the most it asks for.
func allowedOverride(loaded loadedExtension, o ManifestOverride) (HookOverride, bool) {
	allowance, ok := loaded.Entry.Overrides[o.Default]
	if !ok {
		return HookOverride{}, false
	}
	resolved := HookOverride{Default: o.Default, Extension: loaded.Entry.ID, Stage: o.Stage}
	if o.Default == OverrideRepairLoop {
		resolved.Rounds = o.Rounds
		if allowance.Rounds >= 1 && allowance.Rounds < o.Rounds {
			resolved.Rounds = allowance.Rounds
		}
	}
	resolved.Verifier = overrideRole(loaded.Entry.ID, o.Verifier)
	if o.Default == OverrideRiskReview {
		resolved.Reviewer = overrideRole(loaded.Entry.ID, orPlaceholder(o.Reviewer, "reviewer"))
	}
	resolved.Role = overrideRole(loaded.Entry.ID, o.Role)
	if o.Instructions != "" {
		resolved.Instructions = loaded.Folder + "/" + filepath.ToSlash(filepath.Clean(filepath.FromSlash(o.Instructions)))
	}
	return resolved, true
}

// pointOverrides are the allowed overrides of a usable extension that a point
// returns. Whether the extension is on for the member and covers the plan is
// the caller's to check, as it is for the extension's hooks.
func pointOverrides(loaded loadedExtension, point HookPoint) []HookOverride {
	if loaded.Manifest == nil || len(loaded.Entry.Overrides) == 0 {
		return nil
	}
	var out []HookOverride
	for _, o := range loaded.Manifest.Overrides {
		if overridePoint(o) != point {
			continue
		}
		if resolved, ok := allowedOverride(loaded, o); ok {
			out = append(out, resolved)
		}
	}
	return out
}

// dispatchOverride finds the override a dispatch names for the active member
// and the plan's repositories, preferring the extension whose override runs
// role. It says why when none applies.
func (s *Store) dispatchOverride(name, role string, repositories []string) (HookOverride, error) {
	if !slices.Contains(OverridableDefaults, name) {
		return HookOverride{}, fmt.Errorf("--override names an overridable default (%s), not %q; %s", overridableList(), name, fixedRulesNote)
	}
	file, present, err := s.Extensions()
	if err != nil {
		return HookOverride{}, err
	}
	member, _ := s.ActiveMember()
	var candidates []HookOverride
	var why []string
	for index, entry := range file.Extensions {
		if _, allowed := entry.Overrides[name]; !allowed {
			continue
		}
		loaded := s.loadExtension(entry, index)
		if loaded.Fault != extensionUsable {
			why = append(why, entry.ID+" cannot be read: "+loaded.Problem)
			continue
		}
		var declared *ManifestOverride
		for i := range loaded.Manifest.Overrides {
			if loaded.Manifest.Overrides[i].Default == name {
				declared = &loaded.Manifest.Overrides[i]
			}
		}
		if declared == nil {
			why = append(why, entry.ID+"'s manifest does not declare it")
			continue
		}
		if on, reason := loaded.memberState(member); !on {
			why = append(why, entry.ID+" is "+reason)
			continue
		}
		if len(entry.Repositories) > 0 && !overlaps(entry.Repositories, repositories) {
			why = append(why, entry.ID+" is limited to repositories "+strings.Join(entry.Repositories, ", ")+", which do not cover this plan")
			continue
		}
		resolved, _ := allowedOverride(loaded, *declared)
		candidates = append(candidates, resolved)
	}
	for _, candidate := range candidates {
		if candidate.Verifier == role || candidate.Reviewer == role {
			return candidate, nil
		}
	}
	if len(candidates) > 0 {
		return candidates[0], nil
	}
	if !present || len(why) == 0 {
		return HookOverride{}, fmt.Errorf("%s is not allowed in this workspace: only a person allows an override, with `extension override --id ID --allow %s`", name, name)
	}
	return HookOverride{}, fmt.Errorf("%s does not apply to this dispatch: %s", name, strings.Join(why, "; "))
}

// checkOverrideUse accepts the one use of each override that unlocks a
// dispatch, and refuses any other, naming what is missing. round is 0 when the
// dispatch names none.
func checkOverrideUse(role string, r Role, o HookOverride, task, plan string, round int) error {
	if o.Default == OverrideDelegation {
		return errors.New("--override delegation unlocks no dispatch: it says which role a stage goes to; dispatch that role as usual")
	}
	if plan == "" {
		return fmt.Errorf("--override %s runs a step of a plan's execution: pass --plan", o.Default)
	}
	limit := 1
	if o.Default == OverrideRepairLoop {
		limit = o.Rounds
	}
	if round < 0 || round > limit {
		return fmt.Errorf("--round is between 1 and %d for %s from %s, not %d", limit, o.Default, o.Extension, round)
	}
	switch {
	case role == "worker":
		if o.Default != OverrideRepairLoop {
			return fmt.Errorf("--override %s unlocks no worker dispatch; only repair-loop dispatches a worker, with the verifier's findings", o.Default)
		}
		if strings.TrimSpace(task) == "" {
			return errors.New("a repair-loop worker fixes what the verifier found: pass the findings as --task")
		}
		if round != 0 && round >= limit {
			return fmt.Errorf("round %d is the last of %d, so no fix follows it: report the findings to the person", round, limit)
		}
		return nil
	case o.Default == OverrideRiskReview && role == o.Reviewer, o.Default != OverrideRiskReview && role == o.Verifier:
		if !r.ReadOnly {
			return fmt.Errorf("%s is not read-only, so it cannot verify or review the code", role)
		}
		return nil
	}
	unlocked := o.Verifier
	if o.Default == OverrideRiskReview {
		unlocked = o.Reviewer
	}
	if o.Default == OverrideRepairLoop {
		unlocked += ", or a worker with the findings as --task"
	}
	return fmt.Errorf("--override %s from %s unlocks %s, not %s", o.Default, o.Extension, unlocked, role)
}

// overrideBrief is the line a brief carries for a dispatch an override
// unlocked.
func (s *Store) overrideBrief(o HookOverride, round int) string {
	line := "This dispatch runs under the override " + o.Default + " that this workspace allowed for the extension " + o.Extension
	if round > 0 {
		limit := 1
		if o.Default == OverrideRepairLoop {
			limit = o.Rounds
		}
		line += fmt.Sprintf(", round %d of %d", round, limit)
	}
	line += "."
	if o.Instructions != "" {
		line += " Follow " + filepath.Join(s.Root, filepath.FromSlash(o.Instructions)) + " for what counts as clean and what to report."
	}
	return "# Override\n\n" + line
}

// OverrideExtension allows one override for an extension in this workspace,
// or revokes it. Only a person's request runs it. An allowance names a default
// the manifest declares, and its rounds may lower the manifest's, never raise
// them.
func (s *Store) OverrideExtension(id, name string, rounds int, revoke bool) (map[string]OverrideAllowance, error) {
	file, index, err := s.extensionEntry(id)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(OverridableDefaults, name) {
		return nil, fmt.Errorf("%q is not an overridable default; the overridable defaults are %s, and %s", name, overridableList(), fixedRulesNote)
	}
	entry := file.Extensions[index]
	keys := []string{"extensions", ListIndex(index), "overrides"}
	if revoke {
		if rounds != 0 {
			return nil, errors.New("--rounds goes with --allow")
		}
		if _, ok := entry.Overrides[name]; !ok {
			return nil, fmt.Errorf("%s has no allowance for %s to revoke", id, name)
		}
		if len(entry.Overrides) == 1 {
			err = s.RemoveKey(ExtensionsPath, keys, 0644)
		} else {
			err = s.RemoveKey(ExtensionsPath, append(keys, name), 0644)
		}
		remaining := maps.Clone(entry.Overrides)
		delete(remaining, name)
		return remaining, err
	}
	loaded := s.loadExtension(entry, index)
	if loaded.Fault != extensionUsable {
		return nil, fmt.Errorf("%s cannot be used: %s", id, loaded.Problem)
	}
	var declared *ManifestOverride
	for i := range loaded.Manifest.Overrides {
		if loaded.Manifest.Overrides[i].Default == name {
			declared = &loaded.Manifest.Overrides[i]
		}
	}
	if declared == nil {
		return nil, fmt.Errorf("%s's manifest does not declare the override %s, so there is nothing to allow", id, name)
	}
	allowance := OverrideAllowance{Rounds: rounds}
	if rounds != 0 {
		if name != OverrideRepairLoop {
			return nil, fmt.Errorf("--rounds belongs to repair-loop; %s takes none", name)
		}
		if rounds < 1 || rounds > declared.Rounds {
			return nil, fmt.Errorf("--rounds is between 1 and the manifest's %d; an allowance may lower the rounds an extension asks for, never raise them", declared.Rounds)
		}
	}
	if len(entry.Overrides) == 0 {
		err = s.Update(ExtensionsPath, keys, map[string]OverrideAllowance{name: allowance}, 0644)
	} else {
		err = s.Update(ExtensionsPath, append(keys, name), allowance, 0644)
	}
	allowed := maps.Clone(entry.Overrides)
	if allowed == nil {
		allowed = map[string]OverrideAllowance{}
	}
	allowed[name] = allowance
	return allowed, err
}

// OverrideStatus is one allowed override as status reports it.
type OverrideStatus struct {
	Default     string `json:"default"`
	Extension   string `json:"extension"`
	Rounds      int    `json:"rounds,omitempty"`
	OnForMember bool   `json:"on_for_member"`
}

// overrideStatuses lists each override this workspace allowed and its
// extension declares, or nil where none is.
func (s *Store) overrideStatuses(member string) []OverrideStatus {
	file, present, err := s.Extensions()
	if err != nil || !present {
		return nil
	}
	var out []OverrideStatus
	for index, entry := range file.Extensions {
		if len(entry.Overrides) == 0 {
			continue
		}
		loaded := s.loadExtension(entry, index)
		if loaded.Fault != extensionUsable {
			continue
		}
		on, _ := loaded.memberState(member)
		for _, o := range loaded.Manifest.Overrides {
			if resolved, ok := allowedOverride(loaded, o); ok {
				out = append(out, OverrideStatus{resolved.Default, entry.ID, resolved.Rounds, on})
			}
		}
	}
	return out
}

// overrideIssues reports what is wrong with the allowances in extensions.yaml.
// A manifest's own overrides are checked when it is loaded, so an override
// outside the closed list, a verifier that is not read-only, or rounds out of
// range already make the extension invalid.
func (s *Store) overrideIssues(file ExtensionsFile) []Finding {
	var issues []Finding
	delegations := map[string]HookOverride{} // stage -> first allowed delegation
	for index, entry := range file.Extensions {
		if len(entry.Overrides) == 0 {
			continue
		}
		label := fmt.Sprintf("extensions.yaml entry %d (%s)", index+1, entry.ID)
		revoke := func(name string) string { return "`extension override --id " + entry.ID + " --revoke " + name + "`" }
		loaded := s.loadExtension(entry, index)
		for _, name := range slices.Sorted(maps.Keys(entry.Overrides)) {
			allowance := entry.Overrides[name]
			if !slices.Contains(OverridableDefaults, name) {
				issues = append(issues, found(label+" allows "+name+", which is not an overridable default, so it has no effect",
					needsAPerson+fixedRulesNote+"; remove overrides."+name+" from extensions.yaml"))
				continue
			}
			if loaded.Fault != extensionUsable {
				issues = append(issues, found(label+" allows "+name+", but its manifest cannot be read, so the override has no effect",
					needsAPerson+"make the extension readable again, or revoke the allowance with "+revoke(name)))
				continue
			}
			var declared *ManifestOverride
			for i := range loaded.Manifest.Overrides {
				if loaded.Manifest.Overrides[i].Default == name {
					declared = &loaded.Manifest.Overrides[i]
				}
			}
			if declared == nil {
				issues = append(issues, found(label+" allows "+name+", which its manifest does not declare, so it has no effect",
					needsAPerson+"decide whether to revoke it with "+revoke(name)+", or ask the extension's author to declare it"))
				continue
			}
			if allowance.Rounds != 0 {
				switch {
				case name != OverrideRepairLoop:
					issues = append(issues, found(fmt.Sprintf("%s gives %s rounds %d, which only repair-loop takes", label, name, allowance.Rounds),
						needsAPerson+"remove the rounds from overrides."+name+" in extensions.yaml"))
				case allowance.Rounds < 1 || allowance.Rounds > maxRounds:
					issues = append(issues, found(fmt.Sprintf("%s allows repair-loop with rounds %d, outside 1-%d", label, allowance.Rounds, maxRounds),
						needsAPerson+fmt.Sprintf("choose rounds between 1 and %d with `extension override --id %s --allow repair-loop --rounds N`", declared.Rounds, entry.ID)))
				case allowance.Rounds > declared.Rounds:
					issues = append(issues, found(fmt.Sprintf("%s allows repair-loop with rounds %d, above the %d its manifest asks for; the loop stops at %d", label, allowance.Rounds, declared.Rounds, declared.Rounds),
						needsAPerson+fmt.Sprintf("choose rounds up to %d with `extension override --id %s --allow repair-loop --rounds N`", declared.Rounds, entry.ID)))
				}
			}
			if name == OverrideDelegation {
				resolved, _ := allowedOverride(loaded, *declared)
				if other, ok := delegations[resolved.Stage]; ok && other.Role != resolved.Role {
					issues = append(issues, found(fmt.Sprintf("extensions %s and %s both allow delegation of the %s stage, to %s and to %s", other.Extension, entry.ID, resolved.Stage, other.Role, resolved.Role),
						needsAPerson+"choose which role the stage goes to, and revoke the other allowance with `extension override --id ID --revoke delegation`"))
					continue
				}
				delegations[resolved.Stage] = resolved
			}
		}
	}
	return issues
}
