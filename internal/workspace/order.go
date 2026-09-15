package workspace

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Ordering is read-only derivation over recorded plan dependencies. It reports
// waves, start references, and the integration merges a dependent plan needs.
// It reserves nothing, runs nothing, and merges nothing: the coordinator
// decides what to act on and performs the Git work itself.

type OrderStart struct {
	Base  string   `yaml:"base" json:"base"`
	Merge []string `yaml:"merge,omitempty" json:"merge,omitempty"`
}

type OrderPlan struct {
	ID           string                `yaml:"id" json:"id"`
	Repositories []string              `yaml:"repositories" json:"repositories"`
	DependsOn    []string              `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Start        map[string]OrderStart `yaml:"start" json:"start"`
	Shares       []string              `yaml:"shares_repository_with,omitempty" json:"shares_repository_with,omitempty"`
}

type OrderWave struct {
	Wave        int         `yaml:"wave" json:"wave"`
	Concurrency int         `yaml:"concurrency" json:"concurrency"`
	Plans       []OrderPlan `yaml:"plans" json:"plans"`
}

type OrderIntegration struct {
	Plan       string   `yaml:"plan" json:"plan"`
	Wave       int      `yaml:"wave" json:"wave"`
	Repository string   `yaml:"repository" json:"repository"`
	Base       string   `yaml:"base" json:"base"`
	Merge      []string `yaml:"merge" json:"merge"`
}

type OrderOption struct {
	Layers            int `yaml:"layers" json:"layers"`
	ParallelPeak      int `yaml:"parallel_peak" json:"parallel_peak"`
	IntegrationMerges int `yaml:"integration_merges" json:"integration_merges"`
}

type OrderStrategy struct {
	Repositories    int         `yaml:"repositories" json:"repositories"`
	FanIns          int         `yaml:"fan_ins" json:"fan_ins"`
	ConflictSurface string      `yaml:"conflict_surface" json:"conflict_surface"`
	Waves           OrderOption `yaml:"waves" json:"waves"`
	Linear          OrderOption `yaml:"linear" json:"linear"`
	Recommended     string      `yaml:"recommended" json:"recommended"`
	Reason          string      `yaml:"reason" json:"reason"`
}

type OrderBlocked struct {
	Plan   string `yaml:"plan" json:"plan"`
	Reason string `yaml:"reason" json:"reason"`
}

type Order struct {
	Intent      string             `yaml:"intent,omitempty" json:"intent,omitempty"`
	Selected    string             `yaml:"selected" json:"selected"`
	Strategy    OrderStrategy      `yaml:"strategy" json:"strategy"`
	Layers      []OrderWave        `yaml:"layers,omitempty" json:"layers,omitempty"`
	Chain       []OrderPlan        `yaml:"chain,omitempty" json:"chain,omitempty"`
	Integration []OrderIntegration `yaml:"integration_points,omitempty" json:"integration_points,omitempty"`
	Completed   []string           `yaml:"completed,omitempty" json:"completed,omitempty"`
	Blocked     []OrderBlocked     `yaml:"blocked,omitempty" json:"blocked,omitempty"`
	Remaining   int                `yaml:"remaining" json:"remaining"`
}

// PlanBranch is the default branch for one plan's work in one repository.
func PlanBranch(plan, repo string) string { return "cc/" + plan + "/" + repo }

func (s *Store) Order(intent, mode string) (Order, error) {
	switch mode {
	case "", "auto":
		mode = "auto"
	case "waves", "linear":
	default:
		return Order{}, errors.New("order mode must be auto, waves, or linear")
	}
	if intent != "" {
		parent, err := s.FindRecord(intent)
		if err != nil {
			return Order{}, err
		}
		if !strings.HasPrefix(parent.ID, "i") {
			return Order{}, errors.New("--intent needs an intent ID")
		}
	}
	cfg, err := s.Config()
	if err != nil {
		return Order{}, err
	}
	// An order is derived where it will be run, so a repository starts from the
	// branch this machine bases work on rather than the shared default.
	bindings, err := s.Bindings()
	if err != nil {
		return Order{}, err
	}
	records, err := s.ListRecords(false)
	if err != nil {
		return Order{}, err
	}
	plans := map[string]Record{}
	var ids []string
	for _, record := range records {
		if !strings.HasPrefix(record.ID, "p") {
			continue
		}
		if intent != "" && record.Intent != intent {
			continue
		}
		plans[record.ID] = record
		ids = append(ids, record.ID)
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return Order{}, errors.New("no plans to order")
	}

	result := Order{Intent: intent}
	satisfied := map[string]bool{}
	blocked := map[string]string{}
	for _, id := range ids {
		if plans[id].CompletedAt != "" {
			satisfied[id] = true
			result.Completed = append(result.Completed, id)
		}
	}
	// A dependency outside the selection still constrains readiness. A completed
	// one is satisfied; an unfinished one holds its dependents rather than
	// silently widening the run to another intent's work.
	for _, id := range ids {
		for _, dep := range plans[id].DependsOn {
			if _, ok := plans[dep]; ok {
				continue
			}
			external, err := s.FindRecord(dep)
			if err != nil {
				return Order{}, fmt.Errorf("%s: %w", id, err)
			}
			if external.CompletedAt != "" {
				satisfied[dep] = true
				continue
			}
			blocked[id] = "dependency " + dep + " is outside this selection and not completed"
		}
	}

	var remaining []string
	for _, id := range ids {
		if satisfied[id] {
			continue
		}
		if reason, held := blocked[id]; held {
			result.Blocked = append(result.Blocked, OrderBlocked{id, reason})
			continue
		}
		remaining = append(remaining, id)
	}
	result.Remaining = len(remaining)

	var waves [][]string
	for len(remaining) > 0 {
		var ready, rest []string
		for _, id := range remaining {
			if slices.ContainsFunc(plans[id].DependsOn, func(dep string) bool { return !satisfied[dep] }) {
				rest = append(rest, id)
				continue
			}
			ready = append(ready, id)
		}
		if len(ready) == 0 {
			return Order{}, fmt.Errorf("unsatisfiable or cyclic plan dependencies among: %s", strings.Join(remaining, ", "))
		}
		waves = append(waves, ready)
		for _, id := range ready {
			satisfied[id] = true
		}
		remaining = rest
	}

	repositories := map[string]bool{}
	for _, id := range ids {
		for _, repo := range plans[id].Repositories {
			repositories[repo] = true
		}
	}
	result.Strategy.Repositories = len(repositories)

	for number, wave := range waves {
		entry := OrderWave{Wave: number + 1, Concurrency: len(wave)}
		for _, id := range wave {
			plan := plans[id]
			item := OrderPlan{
				ID:           id,
				Repositories: plan.Repositories,
				DependsOn:    plan.DependsOn,
				Start:        map[string]OrderStart{},
			}
			for _, repo := range plan.Repositories {
				start := s.startFor(plans, plan, repo, bindings.BaseBranch(cfg, repo))
				item.Start[repo] = start
				if len(start.Merge) > 0 {
					result.Strategy.FanIns++
					result.Strategy.Waves.IntegrationMerges += len(start.Merge)
					result.Integration = append(result.Integration, OrderIntegration{
						Plan: id, Wave: number + 1, Repository: repo, Base: start.Base, Merge: start.Merge,
					})
				}
			}
			for _, other := range wave {
				if other != id && sharesRepository(plan.Repositories, plans[other].Repositories) {
					item.Shares = append(item.Shares, other)
				}
			}
			entry.Plans = append(entry.Plans, item)
		}
		if entry.Concurrency > result.Strategy.Waves.ParallelPeak {
			result.Strategy.Waves.ParallelPeak = entry.Concurrency
		}
		result.Layers = append(result.Layers, entry)
	}
	result.Strategy.Waves.Layers = len(waves)

	var flat []string
	for _, wave := range waves {
		flat = append(flat, wave...)
	}
	chain := s.linearChain(plans, flat, cfg, bindings)
	linearPeak := 0
	if len(flat) > 0 {
		linearPeak = 1
	}
	result.Strategy.Linear = OrderOption{Layers: len(flat), ParallelPeak: linearPeak}
	result.Strategy.ConflictSurface = conflictSurface(result.Strategy.Repositories, result.Layers)
	result.Strategy.Recommended, result.Strategy.Reason = recommend(result.Strategy)

	result.Selected = mode
	if mode == "auto" {
		result.Selected = result.Strategy.Recommended
	}
	if result.Selected == "linear" {
		result.Layers, result.Integration, result.Chain = nil, nil, chain
	}
	return result, nil
}

// startFor resolves one repository's starting point. A predecessor that also
// includes the repository supplies real Git ancestry; several predecessors need
// an integration merge, which the coordinator performs with ordinary Git. The
// base is the lowest predecessor ID so the same graph always resolves alike.
func (s *Store) startFor(plans map[string]Record, plan Record, repo, baseBranch string) OrderStart {
	var predecessors []string
	for _, dep := range plan.DependsOn {
		record, ok := plans[dep]
		if !ok {
			external, err := s.FindRecord(dep)
			if err != nil {
				continue
			}
			record = external
		}
		if slices.Contains(record.Repositories, repo) {
			predecessors = append(predecessors, dep)
		}
	}
	sort.Strings(predecessors)
	if len(predecessors) == 0 {
		return OrderStart{Base: baseBranch}
	}
	start := OrderStart{Base: PlanBranch(predecessors[0], repo)}
	for _, extra := range predecessors[1:] {
		start.Merge = append(start.Merge, PlanBranch(extra, repo))
	}
	return start
}

// linearChain stacks every plan on the previous chain member that shares its
// repository. A topological order guarantees each link already contains its own
// dependencies, so a chain never needs an integration merge.
func (s *Store) linearChain(plans map[string]Record, order []string, cfg Config, bindings Bindings) []OrderPlan {
	var chain []OrderPlan
	for position, id := range order {
		plan := plans[id]
		item := OrderPlan{
			ID:           id,
			Repositories: plan.Repositories,
			DependsOn:    plan.DependsOn,
			Start:        map[string]OrderStart{},
		}
		for _, repo := range plan.Repositories {
			base := bindings.BaseBranch(cfg, repo)
			for earlier := position - 1; earlier >= 0; earlier-- {
				if slices.Contains(plans[order[earlier]].Repositories, repo) {
					base = PlanBranch(order[earlier], repo)
					break
				}
			}
			item.Start[repo] = OrderStart{Base: base}
		}
		chain = append(chain, item)
	}
	return chain
}

func sharesRepository(left, right []string) bool {
	return slices.ContainsFunc(left, func(repo string) bool { return slices.Contains(right, repo) })
}

// conflictSurface reports how much undelivered work shares a repository. Plan
// records name repositories, not paths, so this is a coarse signal: it cannot
// tell neighbouring modules from the same file.
func conflictSurface(repositories int, layers []OrderWave) string {
	shared := false
	peak := 0
	for _, wave := range layers {
		if wave.Concurrency > peak {
			peak = wave.Concurrency
		}
		for _, plan := range wave.Plans {
			if len(plan.Shares) > 0 {
				shared = true
			}
		}
	}
	switch {
	case peak <= 1 || !shared:
		return "none"
	case repositories == 1:
		return "high"
	default:
		return "medium"
	}
}

func recommend(strategy OrderStrategy) (string, string) {
	if strategy.Waves.ParallelPeak <= 1 {
		return "linear", "no plan runs beside another; the chain is the same work"
	}
	if strategy.Repositories == 1 && strategy.FanIns > 0 {
		return "linear", fmt.Sprintf("single repository with %d fan-in(s); waves buy wall-clock and cost %d integration merge(s)",
			strategy.FanIns, strategy.Waves.IntegrationMerges)
	}
	return "waves", fmt.Sprintf("%d repositories with independent plans; waves overlap them without an integration merge", strategy.Repositories)
}
