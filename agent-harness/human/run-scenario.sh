#!/bin/sh
# Human-simulated agent test harness — prepare and (optionally) drive one case.
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
ROLE_TIERING_SOURCE=${CC_ROLE_TIERING_SOURCE:-$HERE/../fixtures/role-tiering.yaml}

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
[ -f "$ROLE_TIERING_SOURCE" ] || { printf 'FAIL: missing harness role-tiering fixture: %s\n' "$ROLE_TIERING_SOURCE" >&2; exit 1; }

# The fixture is grouped by host. The generated workspace receives the complete
# file so every host run has the same input; the coordinator selects its own
# hosts.<host> group. This is intentionally a small structural check rather than
# a second YAML parser.
role_tiering_has_host() {
	awk -v want="$1" '
		/^hosts:[[:space:]]*$/{in_hosts=1; next}
		in_hosts && /^  [^[:space:]][^:]*:[[:space:]]*$/{h=$0; sub(/^  /,"",h); sub(/:.*/,"",h); if(h==want){found=1; exit}}
		END{exit found?0:1}
	' "$ROLE_TIERING_SOURCE"
}
role_tiering_has_host "$HOST" || { printf 'FAIL: role-tiering fixture has no host group: %s\n' "$HOST" >&2; exit 1; }

# Minimal top-level scalar reader (controlled subset; same discipline as engine.sh).
yscalar() { awk -v k="$2" '$0 ~ "^" k ":[[:space:]]" { sub("^" k ":[[:space:]]*",""); sub(/[[:space:]]*#.*$/,""); sub(/[[:space:]]+$/,""); print; exit }' "$1"; }
CASE_ID=$(yscalar "$CASE_FILE" id); [ -n "$CASE_ID" ] || CASE_ID=$(basename -- "$CASE_DIR")
CASE_MODE=$(yscalar "$CASE_FILE" mode); [ -n "$CASE_MODE" ] || CASE_MODE=conversation-only
CASE_FAULT=$(yscalar "$CASE_FILE" fault)

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
sh "$ROOT/scripts/release-artifact.sh" "$STAGE" "$OUT" v1.0.0 >/dev/null
ARTIFACT="$OUT/context-circuit-v1.0.0"
[ -d "$ARTIFACT" ] || { printf 'FAIL: assembly produced no artifact\n' >&2; exit 1; }

# 2. Instantiate an isolated workspace from the artifact only.
cp -R "$ARTIFACT" "$WORKSPACE"
# Keep the shipped workspace clean while giving the host adapter the same
# host-local role configuration for every scenario.
cp "$ROLE_TIERING_SOURCE" "$WORKSPACE/role-tiering.local.yaml"
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

# Emit one TSV row per setup.intents entry:
#   id title repository objective open_question tier state path
# Intents are the v1.0 Gate-1 fixture. A draft intent has no derived plan; an
# approved intent is ready for the coordinator to derive a plan automatically.
intent_fixtures() {
	awk 'BEGIN{S=sprintf("%c",31)}
		/^  intents:/{ini=1; next}
		ini && /^  [A-Za-z]/ && $0 !~ /^    /{ini=0}
		ini && /^    -[[:space:]]*id:[[:space:]]*/{
			if(id!="") print id S title S repo S obj S oq S tier S state S path;
			id=$0; sub(/^    -[[:space:]]*id:[[:space:]]*/,"",id); gsub(/[[:space:]]+$/,"",id);
			title="";repo="";obj="";oq="";tier="";state="";path=""; next
		}
		ini && id!="" && /^      title:/{v=$0;sub(/^      title:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);title=v;next}
		ini && id!="" && /^      repository:/{v=$0;sub(/^      repository:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);repo=v;next}
		ini && id!="" && /^      objective:/{v=$0;sub(/^      objective:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);obj=v;next}
		ini && id!="" && /^      open_question:/{v=$0;sub(/^      open_question:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);oq=v;next}
		ini && id!="" && /^      tier:/{v=$0;sub(/^      tier:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);tier=v;next}
		ini && id!="" && /^      state:/{v=$0;sub(/^      state:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);state=v;next}
		ini && id!="" && /^      path:/{v=$0;sub(/^      path:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);path=v;next}
		END{if(id!="") print id S title S repo S obj S oq S tier S state S path}
	' "$CASE_FILE"
}

# Emit one TSV row per setup.plans entry:
#   id title repository objective open_question seed_state path deps tier intent_state intent_path
# `path` scopes the task and its verifiable target (a live worker creates
# <path>/mod.txt; verification is `test -f <path>/mod.txt`). `deps` (an inline
# list) adds the schema-3 plan_dependencies block. Both are
# optional: a plan with neither behaves exactly as before (case 05/06).
plan_fixtures() {
	awk 'BEGIN{S=sprintf("%c",31)}
		/^  plans:/{inp=1; next}
		inp && /^  [A-Za-z]/ && $0 !~ /^    /{inp=0}
		inp && /^    -[[:space:]]*id:[[:space:]]*/{
			if(id!="") print id S title S repo S obj S oq S ss S path S deps S tier S istate S ipath;
			id=$0; sub(/^    -[[:space:]]*id:[[:space:]]*/,"",id); gsub(/[[:space:]]+$/,"",id);
			title="";repo="";obj="";oq="";ss="";path="";deps="";tier="";istate="";ipath=""; next
		}
		inp && id!="" && /^      title:/{v=$0;sub(/^      title:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);title=v;next}
		inp && id!="" && /^      repository:/{v=$0;sub(/^      repository:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);repo=v;next}
		inp && id!="" && /^      objective:/{v=$0;sub(/^      objective:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);obj=v;next}
		inp && id!="" && /^      open_question:/{v=$0;sub(/^      open_question:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);oq=v;next}
		inp && id!="" && /^      seed_state:/{v=$0;sub(/^      seed_state:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);ss=v;next}
		inp && id!="" && /^      path:/{v=$0;sub(/^      path:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);path=v;next}
		inp && id!="" && /^      deps:/{v=$0;sub(/^      deps:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/^\[|\]$/,"",v);gsub(/[[:space:]]/,"",v);deps=v;next}
		inp && id!="" && /^      tier:/{v=$0;sub(/^      tier:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);tier=v;next}
		inp && id!="" && /^      intent_state:/{v=$0;sub(/^      intent_state:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);istate=v;next}
		inp && id!="" && /^      intent_path:/{v=$0;sub(/^      intent_path:[[:space:]]*/,"",v);sub(/[[:space:]]+#.*$/,"",v);gsub(/[[:space:]]+$/,"",v);ipath=v;next}
		END{ if(id!="") print id S title S repo S obj S oq S ss S path S deps S tier S istate S ipath }
	' "$CASE_FILE"
}

# Drive the shipped engine to bring a seeded plan to a pre-execution state. A plan is
# authorized by its approved intent — there is no separate plan-approval step and no
# automated scope gate — so a plan is `draft` until it completes.
# seed_state=draft: leave the authorized draft as-is. seed_state=verified-after-repair:
# execute, one FAILED verify + a repair commit, then a PASSED verify (worker_failures=1).
# verified-stale adds a new worker commit after the pass, leaving the prior evidence
# stale. verified-accepted records current-candidate acceptance; completed-standard
# records an explicit mark-done. Delivery is not required to mark done and
# does not start Product Knowledge reconcile.
seed_plan_state() {
	sps_pid="$1"; sps_repo="$2"; sps_state="$3"
	( . "$ENGINE_CLI"
		# a draft plan is already authorized by its approved intent; nothing to seed
		[ "$sps_state" = "draft" ] && exit 0
		# archived: archive the authorized draft, so a restore case starts from a plan
		# sitting in plans/archive/ (status must survive the restore).
		[ "$sps_state" = "archived" ] && { cc_plan_archive "$WORKSPACE" "$sps_pid" >/dev/null; exit $?; }
		case "$sps_state" in verified|verified-after-repair|verified-stale|verified-accepted|completed-standard|worker-committed) : ;; *) exit 0 ;; esac
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
		cc_candidate_digest "$WORKSPACE" "$sps_pid" "$sps_exec" >/dev/null
		if [ "$sps_state" = "verified-stale" ]; then
			# A later worker commit moves the candidate after verification. The old
			# verifier result remains preserved, but no longer authorizes completion.
			cc_attempt_begin "$sps_edir" >/dev/null
			printf 'post-verification change\n' > "$sps_wt/post-verification.txt"
			git -C "$sps_wt" add -A; git -C "$sps_wt" commit -q -m 'post-verification change'
			cc_worker_commit_record "$sps_edir" "$sps_repo" repair >/dev/null
			exit 0
		fi
		if [ "$sps_state" = "verified-accepted" ] || [ "$sps_state" = "completed-standard" ]; then
			cc_human_acceptance_record "$sps_edir" harness-human >/dev/null
		fi
		if [ "$sps_state" = "completed-standard" ]; then
			cc_plan_complete "$WORKSPACE" "$sps_pid" >/dev/null
		fi
	) || return 1
}

FIXTURE_NOTE=role-tiering
if ! setup_empty repositories; then
	FIXTURE_NOTE="${FIXTURE_NOTE},repos"
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
			printf '[setup] seeded + connected repo %s at %s (base %s)\n' "$id" "$dest" "$conn"
		else
			printf '[setup] seeded fixture repo %s at %s (default %s; branches %s)\n' "$id" "$dest" "$defb" "${branches:-$defb}"
		fi
	done
fi

# setup.intents: seed a first-class draft or approved intent. This is separate
# from setup.plans because v1.0 plans are derivations of approved intents; a
# draft intent must not be smuggled into the normal plan lifecycle.
if grep -q '^  intents:' "$CASE_FILE" 2>/dev/null; then
	FIXTURE_NOTE="${FIXTURE_NOTE},intents"
	intent_fixtures | while IFS="$US" read -r iid title irepo iobj iquestion itier istate ipath; do
		[ -n "$iid" ] || continue
		itier=${itier:-standard}; istate=${istate:-draft}; ipath=${ipath:-.}
		idir="$WORKSPACE/intent/$iid"; mkdir -p "$idir"
		{
			printf 'schema_version: 2\nintent: %s\ntitle: %s\ngoal: %s\n' "$iid" "$title" "$iobj"
			printf 'non_goals:\n  - none\nconstraints:\n  - none\n'
			printf 'acceptance_criteria:\n  - id: ac-1\n    statement: %s\n' "$iobj"
			printf 'scope:\n  repositories:\n    - id: %s\n      paths: [%s]\n' "$irepo" "$ipath"
			printf 'tier: %s\nstatus: draft\ncontract_digest:\n' "$itier"
		} > "$idir/contract.yaml"
		{
			printf '# %s\n\nGoal: %s\n' "$title" "$iobj"
			[ -n "$iquestion" ] && printf '\n## Open question\n\n- %s\n' "$iquestion"
		} > "$idir/INTENT.md"
		sh "$ENGINE_CLI" intent-index-upsert "$WORKSPACE" "$iid" >/dev/null || :
		if [ "$istate" = approved ]; then
			sh "$ENGINE_CLI" intent-approve "$WORKSPACE" "$iid" >/dev/null \
				|| { printf 'FAIL: could not approve seeded intent %s\n' "$iid" >&2; exit 1; }
		fi
		printf '[setup] seeded intent %s (%s, tier %s)\n' "$iid" "$istate" "$itier"
	done
fi

# setup.plans: seed a schema-3 DRAFT plan (plan.yaml + PLAN.md + INDEX row) via the
# engine. These are post-Gate-1 fixtures for execution, organization, and delivery
# cases; pre-Gate-1 cases use setup.intents instead.
if grep -q '^  plans:' "$CASE_FILE" 2>/dev/null; then
	FIXTURE_NOTE="${FIXTURE_NOTE},plans"
	plan_fixtures | while IFS="$US" read -r pid title prepo obj oq seedstate path deps ptier istate intentpath; do
		[ -n "$pid" ] || continue
		ptier=${ptier:-standard}; istate=${istate:-}
		pdir="$WORKSPACE/plans/$pid"; mkdir -p "$pdir"
		# path scopes the task; without it the task is repo-wide.
		# The verifiable target is <path>/mod.txt when a path is given, else export.py.
		if [ -n "$path" ]; then taskpath="$path"; target="$path/mod.txt"; else taskpath="."; target="export.py"; fi
		iscopepath=${intentpath:-$taskpath}
		# v1.0: every plan derives from an approved parent intent (no automated scope
		# gate). Author an intent i<pid>, then bind the schema-3 plan to it. A plan with
		# an unresolved open question keeps its intent DRAFT (so review / approve and
		# refuse-before-approval scenarios have a real Gate 1 to settle); a clean plan
		# approves the intent, so it is authorized and ready to run.
		intent_seq=${pid%%-*}; intent_seq=${intent_seq#0}; intent_seq=${intent_seq#0}; intent_seq=${intent_seq#0}; intent_seq=${intent_seq#0}; [ -n "$intent_seq" ] || intent_seq=0
		intent_slug=${pid#*-}; iid=$(printf 'i%03d-%s' "$intent_seq" "$intent_slug"); idir="$WORKSPACE/intent/$iid"; mkdir -p "$idir"
		{
			printf 'schema_version: 2\nintent: %s\ntitle: %s\ngoal: %s\n' "$iid" "$title" "$obj"
			printf 'non_goals:\n  - none\nconstraints:\n  - none\n'
			printf 'acceptance_criteria:\n  - id: ac-1\n    statement: %s\n' "$obj"
			printf 'scope:\n  repositories:\n    - id: %s\n      paths: [%s]\n' "$prepo" "$iscopepath"
			printf 'tier: %s\nstatus: draft\ncontract_digest:\n' "$ptier"
		} > "$idir/contract.yaml"
		printf '# %s\n\nGoal: %s\n' "$title" "$obj" > "$idir/INTENT.md"
		sh "$ENGINE_CLI" intent-index-upsert "$WORKSPACE" "$iid" >/dev/null || :
		# A clean fixture intent is approved by default. Explicit intent_state: draft
		# keeps the pre-Gate-1 negative fixture genuinely unauthorized.
		if [ "$istate" != draft ] && [ -z "$oq" ]; then
			sh "$ENGINE_CLI" intent-approve "$WORKSPACE" "$iid" >/dev/null \
				|| { printf 'FAIL: could not approve seeded intent %s\n' "$iid" >&2; exit 1; }
		fi
		{
			printf 'schema_version: 3\nplan: %s\ntitle: %s\nstatus: draft\nobjective: %s\nintent: %s\n' "$pid" "$title" "$obj" "$iid"
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

printf '[setup] copied shared host role-tiering fixture for hosts.codex, hosts.claude-code, and hosts.cursor-agent\n'

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
	[ -n "$CASE_FAULT" ] && printf 'fault: %s\n' "$CASE_FAULT" || :
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
	CC_FAULT="$CASE_FAULT" \
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
