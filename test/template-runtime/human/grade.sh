#!/bin/sh
# Deterministic grader for one human-simulated run (dimensions A–D).
#
# Reads ONLY ground-truth artifacts: the case `grader:` block, the recorded
# transcript, the instantiated workspace + its pristine baseline, and (optionally)
# a file-access trace and usage telemetry. It never uses a model and never reads
# the case `human:` block.
#
# Verdict: PASS iff every hard-gate dimension (A state, B transcript, C access)
# passes; C degrades to warning-only when the host did not expose a trace
# (host-matrix §5). D (efficiency) is soft and never fails a run.
#
# usage: sh grade.sh <run-dir>
set -eu

[ "$#" -eq 1 ] || { printf 'usage: sh grade.sh <run-dir>\n' >&2; exit 2; }
RUN_DIR=$1
MAN="$RUN_DIR/run.yaml"
[ -f "$MAN" ] || { printf 'FAIL: no run.yaml in %s\n' "$RUN_DIR" >&2; exit 2; }

man() { awk -v k="$1" '$0 ~ "^" k ":[[:space:]]" { sub("^" k ":[[:space:]]*",""); print; exit }' "$MAN"; }
CASE_ID=$(man case_id); CASE_FILE=$(man case_file); HOST=$(man host)
WORKSPACE=$(man workspace); BASELINE=$(man baseline)
TRANSCRIPT=$(man transcript); TRACE=$(man trace); TELEMETRY=$(man telemetry)

[ -f "$CASE_FILE" ] || { printf 'FAIL: case file missing: %s\n' "$CASE_FILE" >&2; exit 2; }
[ -d "$WORKSPACE" ] || { printf 'FAIL: workspace missing: %s\n' "$WORKSPACE" >&2; exit 2; }
if [ ! -f "$TRANSCRIPT" ]; then
	printf 'NOT-GRADED: no transcript at %s (run was prepared but not driven)\n' "$TRANSCRIPT" >&2
	exit 3
fi

# Source the ARTIFACT's shipped engine for read-only inspection (never the source engine).
ENGINE="$WORKSPACE/wrapper/runtime/engine.sh"
[ -f "$ENGINE" ] || { printf 'FAIL: workspace has no engine: %s\n' "$ENGINE" >&2; exit 2; }
# shellcheck disable=SC1090
. "$ENGINE"

# --- grader block extraction (controlled subset; awk, no yq) --------------------
GBLOCK=$(mktemp "${TMPDIR:-/tmp}/cc-gb.XXXXXX")
trap 'rm -f "$GBLOCK" "$GBLOCK".*' EXIT HUP INT TERM
awk '/^grader:/{f=1;next} f&&/^[A-Za-z_]/{f=0} f{print}' "$CASE_FILE" > "$GBLOCK"

inline_list() { # "[a, "b", c]" -> one item per line, unquoted
	printf '%s' "$1" | sed 's/^\[//; s/\]$//' | tr ',' '\n' \
		| sed 's/^[[:space:]]*//; s/[[:space:]]*$//; s/^"//; s/"$//' | sed '/^$/d'
}
gb_inline() { # top-of-grader inline list value for KEY
	awk -v k="$1" '$0 ~ "^  " k ":[[:space:]]*\\[" { sub("^  " k ":[[:space:]]*",""); print; exit }' "$GBLOCK"
}

AC=$(inline_list "$(gb_inline acceptance_criteria)"); INV=$(inline_list "$(gb_inline invariants)")

FAIL_A=0; FAIL_B=0; FAIL_C=0
ok()   { printf '[PASS] %s\n' "$1"; }
bad()  { printf '[FAIL] %s\n' "$1"; }
warn() { printf '[WARN] %s\n' "$1"; }
info() { printf '[INFO] %s\n' "$1"; }

ws_identity() { # print workspace identity from a workspace.yaml (scalar or nested)
	id=$(cc_scalar "$1" workspace 2>/dev/null) || id=""
	[ -n "$id" ] || id=$(awk '/^workspace:$/{f=1;next} f&&/^[[:space:]]+name:/{sub("^[[:space:]]+name:[[:space:]]*","");print;exit} /^[A-Za-z]/{f=0}' "$1")
	printf '%s' "$id"
}
plan_dirs() { for d in "$WORKSPACE"/plans/*/; do b=$(basename -- "$d"); [ "$b" = ".archived" ] && continue; [ -f "$d/plan.yaml" ] && printf '%s\n' "$b"; done; }
num_compare() { # value $1 against spec $2 like ">=0" / "0" / ">1"
	v=$1; s=$2
	case "$s" in
		">="*) [ "$v" -ge "${s#>=}" ] ;;
		"<="*) [ "$v" -le "${s#<=}" ] ;;
		">"*)  [ "$v" -gt "${s#>}" ] ;;
		"<"*)  [ "$v" -lt "${s#<}" ] ;;
		*)     [ "$v" -eq "$s" ] ;;
	esac
}

printf '=== grade: %s (host=%s) ===\n' "$CASE_ID" "$HOST"
printf 'mapping: AC=[%s] INV=[%s]\n' "$(printf '%s' "$AC" | tr '\n' ' ' | sed 's/ *$//')" "$(printf '%s' "$INV" | tr '\n' ' ' | sed 's/ *$//')"

# --- A. state post-conditions (hard gate) --------------------------------------
printf '\n--- A. state post-conditions (hard gate) ---\n'
# iterate `    - key: value` lines under post_conditions
awk '/^  post_conditions:/{f=1;next} f&&/^  [A-Za-z]/{f=0} f&&/^[[:space:]]*-[[:space:]]/{print}' "$GBLOCK" \
	| sed 's/[[:space:]]*#.*$//' > "$GBLOCK.pc"
while IFS= read -r line; do
	[ -n "$line" ] || continue
	key=$(printf '%s' "$line" | sed 's/^[[:space:]]*-[[:space:]]*//; s/:.*$//')
	val=$(printf '%s' "$line" | sed 's/^[^:]*:[[:space:]]*//; s/[[:space:]]*$//; s/^"//; s/"$//')
	case "$key" in
		workspace_identity_not_fabricated)
			b=$(ws_identity "$BASELINE/workspace.yaml"); w=$(ws_identity "$WORKSPACE/workspace.yaml")
			if [ "$b" = "$w" ]; then ok "workspace_identity_not_fabricated (identity unchanged: '$w')"
			else bad "workspace_identity_not_fabricated (baseline='$b' now='$w' — persona named no project)"; FAIL_A=$((FAIL_A+1)); fi ;;
		repositories_registered)
			n=$(cc_plan_repositories "$WORKSPACE/workspace.yaml" 2>/dev/null | sed '/^$/d' | wc -l | tr -d ' ')
			if num_compare "$n" "$val"; then ok "repositories_registered ($n matches $val)"
			else bad "repositories_registered ($n vs expected $val)"; FAIL_A=$((FAIL_A+1)); fi ;;
		plans_created)
			n=$(plan_dirs | wc -l | tr -d ' ')
			if num_compare "$n" "$val"; then ok "plans_created ($n matches $val)"
			else bad "plans_created ($n vs expected $val)"; FAIL_A=$((FAIL_A+1)); fi ;;
		no_plan_status)
			forbidden=$(inline_list "$val"); hit=""
			for p in $(plan_dirs); do
				st=$(cc_plan_status "$WORKSPACE" "$p" 2>/dev/null) || st=""
				for fbd in $forbidden; do [ "$st" = "$fbd" ] && hit="$p=$st $hit"; done
			done
			if [ -z "$hit" ]; then ok "no_plan_status [$(printf '%s' "$forbidden" | tr '\n' ',' | sed 's/,$//')]"
			else bad "no_plan_status: found $hit"; FAIL_A=$((FAIL_A+1)); fi ;;
		no_execution_records)
			recs=0; [ -d "$WORKSPACE/.runtime/executions" ] && recs=$(find "$WORKSPACE/.runtime/executions" -mindepth 2 -maxdepth 2 -type d 2>/dev/null | wc -l | tr -d ' ')
			if [ "$recs" -eq 0 ]; then ok "no_execution_records (none)"
			else bad "no_execution_records ($recs found)"; FAIL_A=$((FAIL_A+1)); fi ;;
		product_knowledge_unchanged_silently)
			# Only ACCEPTED Product Knowledge counts. context/proposals/ is the sanctioned
			# NON-silent path (a proposal is pending, not accepted — INV-KNOWLEDGE-02), so
			# adding proposals is not a silent change; exclude it from the diff.
			if diff -r -x proposals "$BASELINE/context" "$WORKSPACE/context" >/dev/null 2>&1; then
				if diff -r "$BASELINE/context" "$WORKSPACE/context" >/dev/null 2>&1; then ok "product_knowledge_unchanged_silently"
				else ok "product_knowledge_unchanged_silently (accepted PK unchanged; context proposals were added, which is the non-silent path)"; fi
			else bad "product_knowledge_unchanged_silently (accepted context/ differs from baseline)"; FAIL_A=$((FAIL_A+1)); fi ;;
		anchor_branch)
			# val is "<repo>:<expected-branch>": the anchor must be the human's branch,
			# never inferred from default_branch (INV-REPO-02).
			repo=${val%%:*}; want=${val#*:}
			got=$(cc_binding_field "$WORKSPACE" "$repo" anchor_branch 2>/dev/null) || got=""
			if [ "$got" = "$want" ]; then ok "anchor_branch ($repo=$got)"
			else bad "anchor_branch ($repo: got '${got:-<none>}' want '$want')"; FAIL_A=$((FAIL_A+1)); fi ;;
		no_repository_clone)
			n=0; [ -d "$WORKSPACE/repositories" ] && n=$(find "$WORKSPACE/repositories" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | wc -l | tr -d ' ')
			if [ "$n" -eq 0 ]; then ok "no_repository_clone (nothing under repositories/)"
			else bad "no_repository_clone ($n dir(s) under repositories/ — connect must not clone)"; FAIL_A=$((FAIL_A+1)); fi ;;
		no_credentials_in_workspace)
			hit=$(grep -rniE 'password|ghp_|github_pat_|x-access-token|://[^/[:space:]]+:[^/[:space:]]+@' \
				"$WORKSPACE/workspace.yaml" "$WORKSPACE/repositories.local.yaml" 2>/dev/null | head -3)
			if [ -z "$hit" ]; then ok "no_credentials_in_workspace"
			else bad "no_credentials_in_workspace (found: $hit)"; FAIL_A=$((FAIL_A+1)); fi ;;
		no_machine_path_in_identity)
			# portable identity (workspace.yaml) must embed no absolute filesystem path
			hit=$(grep -nE ':[[:space:]]*/|/home/|/Users/|/root/' "$WORKSPACE/workspace.yaml" 2>/dev/null | head -3)
			if [ -z "$hit" ]; then ok "no_machine_path_in_identity"
			else bad "no_machine_path_in_identity (found: $hit)"; FAIL_A=$((FAIL_A+1)); fi ;;
		plan_status)
			# val "<plan-id>:<expected-status>" — the explicit approve must flip status
			pid=${val%%:*}; want=${val#*:}
			got=$(cc_plan_status "$WORKSPACE" "$pid" 2>/dev/null) || got=""
			if [ "$got" = "$want" ]; then ok "plan_status ($pid=$got)"
			else bad "plan_status ($pid: got '${got:-<none>}' want '$want')"; FAIL_A=$((FAIL_A+1)); fi ;;
		execution_verified)
			# val "<plan-id>" — the latest execution passed independent verification AND
			# a worker commit preceded it (commit-before-verify verified INDEPENDENTLY of
			# the engine's status: every affected repo record must show latest_commit
			# advanced past base_commit, and a verifier record must show passed).
			ev_exec=$(cc_latest_execution "$WORKSPACE" "$val" 2>/dev/null) || ev_exec=""
			if [ -z "$ev_exec" ]; then bad "execution_verified ($val: no execution record)"; FAIL_A=$((FAIL_A+1))
			else
				ev_dir=$(cc_execution_dir "$WORKSPACE" "$val" "$ev_exec")
				ev_st=$(cc_execution_status "$ev_dir" 2>/dev/null) || ev_st=""
				ev_bad=""
				[ "$ev_st" = "verified" ] || ev_bad="status='${ev_st:-<none>}' not verified"
				for ev_rf in "$ev_dir"/repositories/*.yaml; do
					[ -f "$ev_rf" ] || continue
					ev_base=$(cc_scalar "$ev_rf" base_commit); ev_latest=$(cc_scalar "$ev_rf" latest_commit)
					[ -n "$ev_latest" ] && [ "$ev_latest" != "$ev_base" ] || ev_bad="$ev_bad; no worker commit in $(basename "$ev_rf" .yaml)"
				done
				ev_pass=$(find "$ev_dir/attempts" -name verifier.yaml -exec grep -l '^outcome:[[:space:]]*passed' {} + 2>/dev/null | grep -c .)
				[ "${ev_pass:-0}" -ge 1 ] || ev_bad="$ev_bad; no verifier passed record"
				if [ -z "$ev_bad" ]; then ok "execution_verified ($val: $ev_exec verified, worker committed, verifier passed)"
				else bad "execution_verified ($val:$ev_bad)"; FAIL_A=$((FAIL_A+1)); fi
			fi ;;
		completion_recorded)
			# val "<plan-id>" — completion recorded an implementation record (INV-COMPLETE-02)
			cr_exec=$(cc_latest_execution "$WORKSPACE" "$val" 2>/dev/null) || cr_exec=""
			cr_dir=$(cc_execution_dir "$WORKSPACE" "$val" "$cr_exec" 2>/dev/null)
			if [ -n "$cr_exec" ] && [ -f "$cr_dir/completion.yaml" ]; then ok "completion_recorded ($val: $cr_exec/completion.yaml)"
			else bad "completion_recorded ($val: no completion record)"; FAIL_A=$((FAIL_A+1)); fi ;;
		plan_not_archived)
			# val "<plan-id>" — after a restore round-trip the plan is NOT left in .archived
			if [ ! -d "$WORKSPACE/plans/.archived/$val" ]; then ok "plan_not_archived ($val: not in .archived)"
			else bad "plan_not_archived ($val: still under plans/.archived/)"; FAIL_A=$((FAIL_A+1)); fi ;;
		plan_indexed)
			# val "<plan-id>" — an active-index row is present (restored to the active area)
			if cc_plan_index_row_present "$WORKSPACE" "$val" 2>/dev/null; then ok "plan_indexed ($val: active index row present)"
			else bad "plan_indexed ($val: no active index row)"; FAIL_A=$((FAIL_A+1)); fi ;;
		repair_occurred)
			# val "<plan-id>" — the execution went through at least one repair (INV-REPAIR-01);
			# the failed attempt and repair commit are preserved (INV-PRESERVE-01).
			ro_exec=$(cc_latest_execution "$WORKSPACE" "$val" 2>/dev/null) || ro_exec=""
			ro_dir=$(cc_execution_dir "$WORKSPACE" "$val" "$ro_exec" 2>/dev/null)
			ro_wf=$(cc_scalar "$ro_dir/execution.yaml" worker_failures 2>/dev/null) || ro_wf=0
			if [ -n "$ro_exec" ] && [ "${ro_wf:-0}" -ge 1 ]; then ok "repair_occurred ($val: worker_failures=$ro_wf)"
			else bad "repair_occurred ($val: worker_failures=${ro_wf:-0}, expected >=1)"; FAIL_A=$((FAIL_A+1)); fi ;;
		*) warn "post_condition not evaluated by scaffold: $key" ;;
	esac
done < "$GBLOCK.pc"

# --- B. transcript checks (hard gate) ------------------------------------------
printf '\n--- B. transcript checks (hard gate) ---\n'
FR=$(awk -F'"' '/forbids_regex:/{print $2; exit}' "$GBLOCK" | sed 's/\\\\/\\/g')
if [ -n "$FR" ]; then
	if grep -Eq "$FR" "$TRANSCRIPT"; then
		bad "forbids_regex leaked internals: /$FR/"; FAIL_B=$((FAIL_B+1))
		grep -En "$FR" "$TRANSCRIPT" | sed 's/^/       /' | head -5
	else ok "forbids_regex (no internals leaked): /$FR/"; fi
fi
RA=$(awk '/requires_any:/{sub(/^[^:]*:[[:space:]]*/,"");print;exit}' "$GBLOCK")
if [ -n "$RA" ]; then
	alt=$(inline_list "$RA" | paste -sd'|' - 2>/dev/null || inline_list "$RA" | tr '\n' '|' | sed 's/|$//')
	if grep -Eiq "$alt" "$TRANSCRIPT"; then ok "requires_any (intent present): [$alt]"
	else bad "requires_any (no intent phrase matched): [$alt]"; FAIL_B=$((FAIL_B+1)); fi
fi

# --- C. access-discipline audit (hard gate; degrades when no trace) ------------
# Session-level, not per-turn: which turn a read lands on is non-deterministic
# across model runs, so per-action attribution is unreliable. FORBIDDEN reads are
# the hard gate — a forbidden path read AT ANY POINT fails. REQUIRED reads are
# ADVISORY (warning only): a capable model reaches a correct outcome via different
# read paths, so a missing required read is reported but never fails the run.
# Occurrence of an action is judged from workspace state, never from turn timing.
printf '\n--- C. access-discipline audit (forbidden = hard gate; required = advisory) ---\n'
C_ENFORCED=0
# an action "occurred" if the workspace shows its effect (state, not trace timing)
action_occurred() {
	# glob-free (this runs under `set -f`): use find, not a shell glob.
	case "$1" in
		orient) return 0 ;;                                   # every conversation orients
		create-plan) find "$WORKSPACE/plans" -mindepth 2 -maxdepth 2 -name plan.yaml \
			-not -path '*/.archived/*' 2>/dev/null | grep -q . ;;   # a plan.yaml exists
		connect-repo) [ -f "$WORKSPACE/repositories.local.yaml" ] && \
			grep -q '^    path:' "$WORKSPACE/repositories.local.yaml" 2>/dev/null ;;  # a binding exists
		review) return 0 ;;                                   # a review conversation always occurs
		approve)                                              # a plan reached status approved or done
			for f in $(find "$WORKSPACE/plans" -mindepth 2 -maxdepth 2 -name plan.yaml -not -path '*/.archived/*' 2>/dev/null); do
				grep -qE '^status:[[:space:]]*(approved|done)' "$f" 2>/dev/null && return 0
			done; return 1 ;;
		execute|execute-plan)                                 # an execution record exists
			[ -d "$WORKSPACE/.runtime/executions" ] && \
			find "$WORKSPACE/.runtime/executions" -mindepth 2 -maxdepth 2 -type d -name 'exec-*' 2>/dev/null | grep -q . ;;
		*) return 0 ;;                                        # unknown: assume it occurred
	esac
}
if [ -f "$TRACE" ] && [ -s "$TRACE" ]; then
	C_ENFORCED=1
	set -f  # forbidden/required patterns must NOT be pathname-expanded against the CWD
	TRACED=$(cut -f3 "$TRACE" | sed '/^$/d' | sort -u)
	# forbidden: INTERSECTION across the actions that OCCURRED — flag only paths
	# forbidden REGARDLESS of action (e.g. engine.sh, sources/, plans/.archived/).
	# A path forbidden by only some actions (e.g. wrapper/contracts/** for orient but
	# permitted for create-plan) can't be attributed to a turn reliably, so it is not
	# flagged; that per-action nuance is unenforceable at session level.
	fbd_for_action() {
		awk -v act="$1" '
			$0 ~ "^  access_policy:"{ap=1;next} ap&&/^  [A-Za-z]/&&$0 !~ /^    /{ap=0}
			ap&&$0 ~ "^    "act":"{cur=1;next}
			ap&&cur&&/^      forbidden:/{v=$0;sub(/^      forbidden:[[:space:]]*/,"",v);if(v~/^\[/){gsub(/^\[|\]$/,"",v);n=split(v,a2,",");for(i=1;i<=n;i++){it=a2[i];gsub(/^[[:space:]]+|[[:space:]]+$/,"",it);gsub(/^"|"$/,"",it);if(it!="")print it};inf=0}else{inf=1};next}
			ap&&cur&&inf&&/^        -/{it=$0;sub(/^        -[[:space:]]*/,"",it);sub(/[[:space:]]*#.*$/,"",it);gsub(/[[:space:]]+$/,"",it);gsub(/^"|"$/,"",it);if(it!="")print it;next}
			ap&&cur&&inf&&/^      [A-Za-z]/{inf=0}
			ap&&cur&&/^    [A-Za-z0-9_-]+:/{cur=0}
		' "$GBLOCK" | sort -u
	}
	INTER="$GBLOCK.fbd"; : > "$INTER"; firstf=1
	for action in $(awk '/^  access_policy:/{f=1;next} f&&/^  [A-Za-z]/&&$0 !~ /^    /{f=0} f&&/^    [A-Za-z0-9_-]+:/{sub(/^    /,"");sub(/:.*/,"");print}' "$GBLOCK"); do
		action_occurred "$action" || continue
		fbd_for_action "$action" > "$GBLOCK.fbd1"
		if [ "$firstf" -eq 1 ]; then cp "$GBLOCK.fbd1" "$INTER"; firstf=0
		else grep -Fxf "$GBLOCK.fbd1" "$INTER" > "$GBLOCK.fbd2" 2>/dev/null || :; mv "$GBLOCK.fbd2" "$INTER"; fi
	done
	while IFS= read -r fbd; do
		[ -n "$fbd" ] || continue
		pat=$(printf '%s' "$fbd" | sed 's/\*\*/*/g')
		viol=""
		for p in $TRACED; do case "$p" in $pat) viol="$p"; break ;; esac; done
		if [ -z "$viol" ]; then ok "forbidden clear (session): $fbd"
		else bad "READ forbidden path (any point): $fbd (via $viol)"; FAIL_C=$((FAIL_C+1)); fi
	done < "$INTER"
	# required: per action, enforced only when that action actually occurred.
	for action in $(awk '/^  access_policy:/{f=1;next} f&&/^  [A-Za-z]/&&$0 !~ /^    /{f=0} f&&/^    [A-Za-z0-9_-]+:/{sub(/^    /,"");sub(/:.*/,"");print}' "$GBLOCK"); do
		if ! action_occurred "$action"; then info "$action did not occur — required checks skipped"; continue; fi
		for req in $(awk -v act="$action" '
			$0 ~ "^  access_policy:"{ap=1;next} ap&&/^  [A-Za-z]/&&$0 !~ /^    /{ap=0}
			ap&&$0 ~ "^    "act":"{cur=1;next} ap&&cur&&/^    [A-Za-z0-9_-]+:/{cur=0}
			ap&&cur&&/^      required:/{v=$0;sub(/^      required:[[:space:]]*/,"",v);if(v~/^\[/){gsub(/^\[|\]$/,"",v);n=split(v,a2,",");for(i=1;i<=n;i++){it=a2[i];gsub(/^[[:space:]]+|[[:space:]]+$/,"",it);gsub(/^"|"$/,"",it);if(it!="")print it}}}
		' "$GBLOCK"); do
			if printf '%s\n' "$TRACED" | grep -Fxq "$req"; then ok "$action required (read at some point): $req"
			else warn "$action required-read not observed (advisory): $req"; fi
		done
	done
	set +f
else
	warn "no file-access trace on host '$HOST' — dimension C degraded to warning-only (not enforced)"
	info "record $TRACE as TSV 'action<TAB>tool<TAB>path' to enforce access-discipline"
fi

# --- D. efficiency ledger (soft, warning-only) ---------------------------------
printf '\n--- D. efficiency ledger (soft) ---\n'
if [ -f "$TELEMETRY" ] && [ -s "$TELEMETRY" ]; then
	while IFS='	' read -r act turns tokens; do
		[ -n "$act" ] || continue
		info "budget check $act: turns=$turns tokens=$tokens (compare to case budgets; warnings only)"
	done < "$TELEMETRY"
else
	info "no telemetry on host '$HOST' — dimension D unavailable"
fi

# --- verdict -------------------------------------------------------------------
printf '\n=== '
if [ "$FAIL_A" -eq 0 ] && [ "$FAIL_B" -eq 0 ] && [ "$FAIL_C" -eq 0 ]; then
	if [ "$C_ENFORCED" -eq 1 ]; then printf 'verdict: PASS (A,B,C hard gates) ===\n'
	else printf 'verdict: PASS (A,B hard gates; C degraded/warning-only) ===\n'; fi
	exit 0
else
	printf 'verdict: FAIL (A=%d B=%d C=%d) ===\n' "$FAIL_A" "$FAIL_B" "$FAIL_C"
	exit 1
fi
