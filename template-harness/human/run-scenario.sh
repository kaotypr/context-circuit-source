#!/bin/sh
# Human-simulated template test harness — prepare and (optionally) drive one case.
#
# Deterministic responsibilities (done here, no live model needed):
#   1. assemble context-circuit-template with scripts/release-artifact.sh;
#   2. instantiate an isolated workspace from that artifact;
#   3. apply the case `setup` fixtures (seed repos / source files);
#   4. capture a pristine baseline snapshot (context/ + workspace.yaml);
#   5. write a run manifest under human/.out/<run-id>/ (git-ignored).
#
# Agent-driven responsibility (delegated to a per-host driver via --driver/
#   $CC_HUMAN_DRIVER): spawn the product coordinator with cwd = the instantiated
#   workspace, spawn the human-simulator with only the case `human:` block, run
#   the turns, and record transcript.txt (+ optional file-access-trace.tsv and
#   telemetry.tsv). When no driver is wired the run stops at `prepared`, which is
#   a normal outcome — the prepared workspace can be driven and graded later.
#
# This suite is opt-in and NOT part of `sh test/acceptance.sh` (harness §7).
#
# usage: sh run-scenario.sh [--host H] [--driver CMD] [--no-grade] <case-id|case-dir>
set -eu

HOST=claude-code
DRIVER="${CC_HUMAN_DRIVER:-}"
LIVE=0
AUTOGRADE=1
CASE_ARG=

while [ "$#" -gt 0 ]; do
	case "$1" in
		--host) HOST=${2:?--host needs a value}; shift 2 ;;
		--host=*) HOST=${1#--host=}; shift ;;
		--driver) DRIVER=${2:?--driver needs a value}; shift 2 ;;
		--driver=*) DRIVER=${1#--driver=}; shift ;;
		--live) LIVE=1; shift ;;
		--no-grade) AUTOGRADE=0; shift ;;
		-h|--help)
			printf 'usage: sh run-scenario.sh [--host H] [--live|--driver CMD] [--no-grade] <case>\n'
			printf 'hosts: codex | claude-code | cursor-agent (default: claude-code)\n'
			printf '  --live    use the built-in driver for the host (human/drivers/<host>.sh)\n'
			printf '  --driver  use a custom driver command instead\n'
			exit 0 ;;
		--*) printf 'FAIL: unknown flag: %s\n' "$1" >&2; exit 2 ;;
		*) CASE_ARG=$1; shift ;;
	esac
done
[ -n "$CASE_ARG" ] || { printf 'FAIL: no case given\n' >&2; exit 2; }
case "$HOST" in codex|claude-code|cursor-agent) : ;; *) printf 'FAIL: unknown host: %s\n' "$HOST" >&2; exit 2 ;; esac

HERE=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT=$(git -C "$HERE" rev-parse --show-toplevel 2>/dev/null) || { printf 'FAIL: not a git source checkout\n' >&2; exit 1; }
SCENARIOS="$HERE/../scenarios"

# --live selects the built-in per-host driver unless --driver overrode it.
if [ "$LIVE" -eq 1 ] && [ -z "$DRIVER" ]; then
	BUILTIN="$HERE/drivers/$HOST.sh"
	[ -f "$BUILTIN" ] || { printf 'FAIL: no built-in driver for host %s (%s)\n' "$HOST" "$BUILTIN" >&2; exit 1; }
	DRIVER="sh '$BUILTIN'"
fi

# Resolve the case directory (accept an id or a path).
if [ -d "$CASE_ARG" ]; then CASE_DIR=$(CDPATH= cd -- "$CASE_ARG" && pwd)
elif [ -d "$SCENARIOS/$CASE_ARG" ]; then CASE_DIR="$SCENARIOS/$CASE_ARG"
else printf 'FAIL: case not found: %s\n' "$CASE_ARG" >&2; exit 1; fi
CASE_FILE="$CASE_DIR/case.yaml"
[ -f "$CASE_FILE" ] || { printf 'FAIL: missing case.yaml in %s\n' "$CASE_DIR" >&2; exit 1; }

# Minimal top-level scalar reader (controlled subset; same discipline as engine.sh).
yscalar() { awk -v k="$2" '$0 ~ "^" k ":[[:space:]]" { sub("^" k ":[[:space:]]*",""); sub(/[[:space:]]*#.*$/,""); sub(/[[:space:]]+$/,""); print; exit }' "$1"; }
CASE_ID=$(yscalar "$CASE_FILE" id); [ -n "$CASE_ID" ] || CASE_ID=$(basename -- "$CASE_DIR")
CASE_MODE=$(yscalar "$CASE_FILE" mode); [ -n "$CASE_MODE" ] || CASE_MODE=conversation-only

RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-$$-${CASE_ID}-${HOST}"
OUT_ROOT="$HERE/.out"
RUN_DIR="$OUT_ROOT/$RUN_ID"
WORKSPACE="$RUN_DIR/workspace"
BASELINE="$RUN_DIR/baseline"
TRANSCRIPT="$RUN_DIR/transcript.txt"
TRACE="$RUN_DIR/file-access-trace.tsv"
TELEMETRY="$RUN_DIR/telemetry.tsv"
ROLE_EVIDENCE="$RUN_DIR/role-evidence.tsv"
mkdir -p "$RUN_DIR"

# Temp assembly area (disposable; never the source .runtime).
STAGE=$(mktemp -d "${TMPDIR:-/tmp}/cc-hs-stage.XXXXXX")
OUT=$(mktemp -d "${TMPDIR:-/tmp}/cc-hs-out.XXXXXX")
trap 'rm -rf "$STAGE" "$OUT"' EXIT HUP INT TERM

# 1. Assemble the released template (the same artifact a user receives).
sh "$ROOT/scripts/release-artifact.sh" "$STAGE" "$OUT" v0.5.0 >/dev/null
ARTIFACT="$OUT/context-circuit-v0.5.0"
[ -d "$ARTIFACT" ] || { printf 'FAIL: assembly produced no artifact\n' >&2; exit 1; }

# 2. Instantiate an isolated workspace from the artifact only.
cp -R "$ARTIFACT" "$WORKSPACE"
# Safety: the seed must not carry any source-side state.
for leak in .runtime plans/context-circuit-plans repositories.local.yaml repositories; do
	[ ! -e "$WORKSPACE/$leak" ] || { printf 'FAIL: source state leaked into workspace: %s\n' "$leak" >&2; exit 1; }
done

# 3. Apply case setup fixtures.
#    `setup.repositories` seeds real git repos INSIDE the workspace (the human's
#    "already-existing code"), created at prep time so the coordinator sees them
#    as pre-existing and only has to CONNECT them. `setup.sources` is still a TODO.
setup_empty() { awk -v k="$1" '$0 ~ "^  " k ":[[:space:]]*\\[\\]" {f=1} END{exit f?0:1}' "$CASE_FILE"; }
US=$(printf '\037')   # non-whitespace field separator: preserves EMPTY middle fields (TAB, being IFS-whitespace, collapses them)

ENGINE_CLI="$WORKSPACE/wrapper/runtime/engine.sh"

# Emit one TSV row per setup.repositories entry:
#   id dest default_branch branches seed_files connect agents_md
# `agents_md` (optional) seeds the repo's OWN agent guidance (an AGENTS.md with a
# distinctive convention) so repository-grounding discovery finds real guidance.
repo_fixtures() {
	awk 'BEGIN{S=sprintf("%c",31)}
		/^  repositories:/{inr=1; next}
		inr && /^  [A-Za-z]/ && $0 !~ /^    /{inr=0}
		inr && /^    -[[:space:]]*id:[[:space:]]*/{
			if(id!="") print id S dest S defb S br S sf S conn S am;
			id=$0; sub(/^    -[[:space:]]*id:[[:space:]]*/,"",id); gsub(/[[:space:]]+$/,"",id);
			dest="";defb="";br="";sf="";conn="";am=""; next
		}
		inr && id!="" && /^      dest:/{v=$0;sub(/^      dest:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);dest=v;next}
		inr && id!="" && /^      default_branch:/{v=$0;sub(/^      default_branch:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);defb=v;next}
		inr && id!="" && /^      branches:/{v=$0;sub(/^      branches:[[:space:]]*/,"",v);gsub(/^\[|\]$/,"",v);gsub(/[[:space:]]/,"",v);br=v;next}
		inr && id!="" && /^      seed_files:/{v=$0;sub(/^      seed_files:[[:space:]]*/,"",v);gsub(/^\[|\]$/,"",v);gsub(/[[:space:]]/,"",v);sf=v;next}
		inr && id!="" && /^      connect:/{v=$0;sub(/^      connect:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);conn=v;next}
		inr && id!="" && /^      agents_md:/{v=$0;sub(/^      agents_md:[[:space:]]*/,"",v);gsub(/[[:space:]]+$/,"",v);gsub(/^"|"$/,"",v);am=v;next}
		END{ if(id!="") print id S dest S defb S br S sf S conn S am }
	' "$CASE_FILE"
}

# Emit one TSV row per setup.plans entry:
#   id title repository objective open_question seed_state path deps
# `path` scopes the task and its verifiable target (a live worker creates
# <path>/mod.txt; verification is `test -f <path>/mod.txt`). `deps` (an inline
# list) makes the plan schema_version 2 with a plan_dependencies block. Both are
# optional: a plan with neither behaves exactly as before (case 05/06).
plan_fixtures() {
	awk 'BEGIN{S=sprintf("%c",31)}
		/^  plans:/{inp=1; next}
		inp && /^  [A-Za-z]/ && $0 !~ /^    /{inp=0}
		inp && /^    -[[:space:]]*id:[[:space:]]*/{
			if(id!="") print id S title S repo S obj S oq S ss S path S deps;
			id=$0; sub(/^    -[[:space:]]*id:[[:space:]]*/,"",id); gsub(/[[:space:]]+$/,"",id);
			title="";repo="";obj="";oq="";ss="";path="";deps=""; next
		}
		inp && id!="" && /^      title:/{v=$0;sub(/^      title:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);title=v;next}
		inp && id!="" && /^      repository:/{v=$0;sub(/^      repository:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);repo=v;next}
		inp && id!="" && /^      objective:/{v=$0;sub(/^      objective:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);obj=v;next}
		inp && id!="" && /^      open_question:/{v=$0;sub(/^      open_question:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);oq=v;next}
		inp && id!="" && /^      seed_state:/{v=$0;sub(/^      seed_state:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);ss=v;next}
		inp && id!="" && /^      path:/{v=$0;sub(/^      path:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);path=v;next}
		inp && id!="" && /^      deps:/{v=$0;sub(/^      deps:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/^\[|\]$/,"",v);gsub(/[[:space:]]/,"",v);deps=v;next}
		END{ if(id!="") print id S title S repo S obj S oq S ss S path S deps }
	' "$CASE_FILE"
}

# Drive the shipped engine to bring a seeded plan to a pre-execution state.
# seed_state=approved: approve only. seed_state=verified-after-repair: approve,
# execute, one FAILED verify + a repair commit, then a PASSED verify (worker_failures=1).
seed_plan_state() {
	sps_pid="$1"; sps_repo="$2"; sps_state="$3"
	( . "$ENGINE_CLI"
		cc_plan_approve "$WORKSPACE" "$sps_pid" >/dev/null || { printf 'FAIL: seed approve %s\n' "$sps_pid" >&2; exit 1; }
		[ "$sps_state" = "approved" ] && exit 0
		# archived: approve then archive, so a restore case starts from an approved
		# plan sitting in plans/archive/ (status must survive the restore).
		[ "$sps_state" = "archived" ] && { cc_plan_archive "$WORKSPACE" "$sps_pid" >/dev/null; exit $?; }
		case "$sps_state" in verified|verified-after-repair|worker-committed) : ;; *) exit 0 ;; esac
		sps_exec=$(cc_execution_begin "$WORKSPACE" "$sps_pid" seed-worker | sed -n 's/^execution_id: //p')
		sps_edir="$WORKSPACE/.runtime/executions/$sps_pid/$sps_exec"
		sps_wt=$(cc_scalar "$sps_edir/repositories/$sps_repo.yaml" worktree)
		[ -d "$sps_wt" ] || { printf 'FAIL: seed worktree missing for %s\n' "$sps_pid" >&2; exit 1; }
		sps_att=1
		if [ "$sps_state" = "verified-after-repair" ]; then
			# attempt 1: an implementation that does NOT satisfy verification -> failed
			cc_attempt_begin "$sps_edir" >/dev/null
			printf 'incomplete first attempt\n' > "$sps_wt/notes-wip.txt"
			git -C "$sps_wt" add -A; git -C "$sps_wt" commit -q -m 'attempt 1 (incomplete)'
			cc_worker_commit_record "$sps_edir" "$sps_repo" implementation >/dev/null
			cc_verifier_prepare "$sps_edir" >/dev/null
			cc_verifier_result_record "$sps_edir" 1 failed >/dev/null   # worker_failures -> 1, repairing
			sps_att=2
		fi
		# passing attempt: satisfy verification (create export.py) -> passed -> verified
		cc_attempt_begin "$sps_edir" >/dev/null
		printf 'print("export")\n' > "$sps_wt/export.py"
		git -C "$sps_wt" add -A; git -C "$sps_wt" commit -q -m "attempt $sps_att (implementation)"
		cc_worker_commit_record "$sps_edir" "$sps_repo" implementation >/dev/null
		# worker-committed stops here: a real worker commit exists (status verifying),
		# but verification has NOT run (case 09 forces the verifier child unavailable).
		[ "$sps_state" = "worker-committed" ] && exit 0
		cc_verifier_prepare "$sps_edir" >/dev/null
		cc_verifier_result_record "$sps_edir" "$sps_att" passed >/dev/null
	) || return 1
}

FIXTURE_NOTE=none
if ! setup_empty repositories; then
	FIXTURE_NOTE=repos
	repo_fixtures | while IFS="$US" read -r id dest defb branches seeds conn agentsmd; do
		[ -n "$id" ] || continue
		dest=${dest:-$id}; defb=${defb:-main}
		case "$dest" in /*|*..*) printf 'FAIL: unsafe fixture dest: %s\n' "$dest" >&2; exit 1 ;; esac
		repo="$WORKSPACE/$dest"
		mkdir -p "$repo"
		git -C "$repo" init -q -b "$defb"
		git -C "$repo" config user.email 'cc-fixture@example.invalid'
		git -C "$repo" config user.name 'cc-fixture'
		for f in $(printf '%s' "${seeds:-README.md}" | tr ',' ' '); do
			mkdir -p "$repo/$(dirname -- "$f")"
			printf '# %s\n\nseed content for the %s fixture.\n' "$f" "$id" > "$repo/$f"
		done
		# optional: the repository's OWN agent guidance (discovered by grounding)
		if [ -n "$agentsmd" ]; then
			printf '# %s — agent guide\n\n%s\n' "$id" "$agentsmd" > "$repo/AGENTS.md"
			printf '[setup] seeded repo agent guidance (AGENTS.md) for %s\n' "$id"
		fi
		git -C "$repo" add -A
		git -C "$repo" commit -q -m 'seed'
		for b in $(printf '%s' "$branches" | tr ',' ' '); do
			[ -n "$b" ] && [ "$b" != "$defb" ] && git -C "$repo" branch "$b" 2>/dev/null || :
		done
		# optional pre-connect: register + bind via the shipped engine so a seeded
		# plan can reference it (case 03). Without `connect:` the repo is left for the
		# coordinator to connect in-conversation (case 02).
		if [ -n "$conn" ]; then
			sh "$ENGINE_CLI" repository-register "$WORKSPACE" "$id" "$dest" "$conn" >/dev/null \
				|| { printf 'FAIL: could not pre-connect fixture %s\n' "$id" >&2; exit 1; }
			printf '[setup] seeded + connected repo %s at %s (anchor %s)\n' "$id" "$dest" "$conn"
		else
			printf '[setup] seeded fixture repo %s at %s (default %s; branches %s)\n' "$id" "$dest" "$defb" "${branches:-$defb}"
		fi
	done
fi

# setup.plans: seed a DRAFT plan (plan.yaml + PLAN.md + INDEX row) via the engine,
# so review/approve cases start from an existing plan with an open question.
if grep -q '^  plans:' "$CASE_FILE" 2>/dev/null; then
	FIXTURE_NOTE="${FIXTURE_NOTE},plans"
	plan_fixtures | while IFS="$US" read -r pid title prepo obj oq seedstate path deps; do
		[ -n "$pid" ] || continue
		pdir="$WORKSPACE/plans/$pid"; mkdir -p "$pdir"
		# path scopes the task; without it the task is repo-wide (backward compatible).
		# The verifiable target is <path>/mod.txt when a path is given, else export.py.
		if [ -n "$path" ]; then taskpath="$path"; target="$path/mod.txt"; else taskpath="."; target="export.py"; fi
		# deps present => schema_version 2 with a plan_dependencies block (v0.6).
		schema=1; [ -n "$deps" ] && schema=2
		{
			printf 'schema_version: %s\nplan: %s\ntitle: %s\nstatus: draft\nobjective: %s\n' "$schema" "$pid" "$title" "$obj"
			printf 'repositories:\n  - id: %s\n' "$prepo"
			if [ -n "$deps" ]; then
				printf 'plan_dependencies:\n'
				for d in $(printf '%s' "$deps" | tr ',' ' '); do
					[ -n "$d" ] && printf '  - id: %s\n    reason: builds on %s\n' "$d" "$d" || :
				done
			fi
			[ -n "$oq" ] && printf 'open_questions:\n  - %s\n' "$oq" || :
			printf 'tasks:\n  - id: EXPORT-001\n    title: %s\n    repositories: [%s]\n    paths: [%s]\n    depends_on: []\n' "$title" "$prepo" "$taskpath"
			# a concrete, trivially-verifiable target so full-execution cases have real
			# worker work (create the target file) and a real verifier check (test -f).
			printf '    changes: [Create %s.]\n' "$target"
			printf '    acceptance:\n      - id: EXPORT-AC-001\n        statement: %s exists in the repository.\n' "$target"
			printf '    verification:\n      - id: EXPORT-VT-001\n        command: test -f %s\n' "$target"
		} > "$pdir/plan.yaml"
		if [ -n "$oq" ]; then
			printf '# %s\n\n%s\n\n## Open question\n\n- %s\n' "$title" "$obj" "$oq" > "$pdir/PLAN.md"
		else
			printf '# %s\n\n%s\n' "$title" "$obj" > "$pdir/PLAN.md"
		fi
		sh "$ENGINE_CLI" plan-validate "$pdir" >/dev/null || { printf 'FAIL: seeded plan %s is invalid\n' "$pid" >&2; exit 1; }
		sh "$ENGINE_CLI" plan-index-upsert "$WORKSPACE" "$pid" >/dev/null || :
		if [ -n "$seedstate" ]; then
			seed_plan_state "$pid" "$prepo" "$seedstate" || { printf 'FAIL: could not seed state %s for %s\n' "$seedstate" "$pid" >&2; exit 1; }
			printf '[setup] seeded plan %s at state "%s" (repo %s)\n' "$pid" "$seedstate" "$prepo"
		else
			printf '[setup] seeded draft plan %s (repo %s%s)\n' "$pid" "$prepo" "$([ -n "$oq" ] && printf ', one open question')"
		fi
	done
fi

# setup.role_tiering: seed a host-local role-tiering.local.yaml at the workspace
# root — the per-user, gitignored, opt-in file that names a concrete (model,
# effort) per role, grouped by host. The case block gives a flat worker/verifier;
# we wrap it under `hosts.<this-run's-host>` so the coordinator on that host reads
# its own group. Seeded here so a case can prove the coordinator reads and honors
# it; a real user writes it by hand and nothing auto-creates it.
rt_read() { # role field -> value from setup.role_tiering
	awk -v role="$1" -v field="$2" '
		/^  role_tiering:/{inrt=1; next}
		inrt && /^  [A-Za-z]/ && $0 !~ /^    / {inrt=0}
		inrt && /^    [A-Za-z]/ {cur=$0; sub(/^    /,"",cur); sub(/:.*/,"",cur)}
		inrt && cur==role && $0 ~ ("^      " field ":") {v=$0; sub(/^[^:]*:[[:space:]]*/,"",v); sub(/[[:space:]]+$/,"",v); print v; exit}
	' "$CASE_FILE"
}
if grep -q '^  role_tiering:' "$CASE_FILE" 2>/dev/null; then
	rt_wm=$(rt_read worker model);   rt_we=$(rt_read worker effort)
	rt_vm=$(rt_read verifier model); rt_ve=$(rt_read verifier effort)
	{
		printf 'hosts:\n  %s:\n' "$HOST"
		if [ -n "$rt_wm" ]; then printf '    worker:\n      model: %s\n' "$rt_wm"; [ -n "$rt_we" ] && printf '      effort: %s\n' "$rt_we" || :; printf '      escalate_on_repair: false\n'; fi
		if [ -n "$rt_vm" ]; then printf '    verifier:\n      model: %s\n' "$rt_vm"; [ -n "$rt_ve" ] && printf '      effort: %s\n' "$rt_ve" || :; printf '      escalate_on_repair: false\n'; fi
	} > "$WORKSPACE/role-tiering.local.yaml"
	FIXTURE_NOTE="${FIXTURE_NOTE},role-tiering"
	printf '[setup] seeded host-local role-tiering.local.yaml under hosts.%s (worker=%s/%s, verifier=%s/%s)\n' "$HOST" "${rt_wm:-default}" "${rt_we:-default}" "${rt_vm:-default}" "${rt_ve:-default}"
fi

setup_empty sources || grep -q '^  sources:[[:space:]]*\[\]' "$CASE_FILE" || { FIXTURE_NOTE="${FIXTURE_NOTE},sources-todo"; printf 'WARN: setup.sources fixtures are not applied by this scaffold (case %s)\n' "$CASE_ID" >&2; }

# 4. Capture the pristine baseline snapshot (harness §1) for dimension A diffs.
mkdir -p "$BASELINE"
cp -R "$WORKSPACE/context" "$BASELINE/context"
cp "$WORKSPACE/workspace.yaml" "$BASELINE/workspace.yaml"

# 5. Write the run manifest.
{
	printf 'run_id: %s\n' "$RUN_ID"
	printf 'case_id: %s\n' "$CASE_ID"
	printf 'case_file: %s\n' "$CASE_FILE"
	printf 'host: %s\n' "$HOST"
	printf 'mode: %s\n' "$CASE_MODE"
	printf 'workspace: %s\n' "$WORKSPACE"
	printf 'baseline: %s\n' "$BASELINE"
	printf 'transcript: %s\n' "$TRANSCRIPT"
	printf 'trace: %s\n' "$TRACE"
	printf 'telemetry: %s\n' "$TELEMETRY"
	printf 'role_evidence: %s\n' "$ROLE_EVIDENCE"
	printf 'human_simulator: %s\n' "$ROOT/.claude/agents/cc-human-simulator.md"
	printf 'fixtures: %s\n' "$FIXTURE_NOTE"
	printf 'status: prepared\n'
} > "$RUN_DIR/run.yaml"

printf 'prepared run: %s\n' "$RUN_DIR"
printf '  case=%s host=%s mode=%s\n' "$CASE_ID" "$HOST" "$CASE_MODE"
printf '  workspace=%s\n' "$WORKSPACE"

# 6. Drive the conversation via the per-host driver, if one is wired.
if [ -n "$DRIVER" ]; then
	printf 'driving via: %s\n' "$DRIVER"
	CC_RUN_DIR="$RUN_DIR" CC_WORKSPACE="$WORKSPACE" CC_BASELINE="$BASELINE" \
	CC_CASE_FILE="$CASE_FILE" CC_CASE_ID="$CASE_ID" CC_HOST="$HOST" CC_MODE="$CASE_MODE" \
	CC_HUMAN_SIM="$ROOT/.claude/agents/cc-human-simulator.md" \
	CC_TRANSCRIPT="$TRANSCRIPT" CC_TRACE="$TRACE" CC_TELEMETRY="$TELEMETRY" \
	CC_ROLE_EVIDENCE="$ROLE_EVIDENCE" \
		sh -c "$DRIVER" || { printf 'FAIL: driver returned non-zero\n' >&2; exit 1; }
	awk '/^status:[[:space:]]/{print "status: driven";next}{print}' "$RUN_DIR/run.yaml" > "$RUN_DIR/run.yaml.tmp" && mv "$RUN_DIR/run.yaml.tmp" "$RUN_DIR/run.yaml"
	if [ "$AUTOGRADE" -eq 1 ]; then
		printf '\n'
		sh "$HERE/grade.sh" "$RUN_DIR"
		exit $?
	fi
else
	printf '\nNo driver wired (status: prepared). To drive on this host:\n'
	printf '  - spawn the product coordinator with cwd=%s\n' "$WORKSPACE"
	printf '  - spawn the human-simulator (.claude/agents/cc-human-simulator.md) with the\n'
	printf "    case's human: block only, run the turns, and write:\n"
	printf '      transcript:  %s\n' "$TRANSCRIPT"
	printf '      trace (opt): %s   (TSV: action<TAB>tool<TAB>path)\n' "$TRACE"
	printf '      telem (opt): %s   (TSV: action<TAB>turns<TAB>tokens)\n' "$TELEMETRY"
	printf '  then grade with:  sh %s %s\n' "$HERE/grade.sh" "$RUN_DIR"
fi
