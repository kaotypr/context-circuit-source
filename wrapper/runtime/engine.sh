#!/bin/sh
# Context Circuit v1.0.0 runtime engine.
#
# A small, host-neutral, deterministic runtime library for workspace, Git, and
# execution-state operations. It is sourced by host adapters and tests, and can
# also be invoked as a thin CLI: `sh engine.sh <command> [--flag value ...]`.
#
# This runtime OWNS: safe path/identifier checks, atomic writes and digests,
# workspace and repository-binding validation, base branch/commit validation,
# branch/worktree preparation, plan structure and intent-derived authorization validation,
# active plan-index maintenance, exact archive/restore moves, execution and
# attempt records, commit capture, verifier-result and read-only enforcement,
# the three-failure counter, one-worker locking, completion eligibility,
# implementation completion records, context-impact handoff references, light
# direct-collaboration session pointers, explicit leftover cleanup, and recovery
# inspection.
#
# This runtime does NOT own: provider-specific child launch, model prompts,
# Product Knowledge interpretation, plan-writing intelligence, conversational
# routing policy, confirmation cards/tokens, product test semantics, or any
# automatic pull-request/merge/push/deploy/cleanup/completion action.
#
# POSIX sh only (runs under dash). No bashisms, no `local`, no arrays.
# Every function returns 0 on success and non-zero with a reason code on stderr.

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

CC_RUNTIME_VERSION="1.0.0"

# Each persisted runtime record owns its schema version independently. These
# constants are record-format versions, not runtime or template versions.
CC_REPOSITORIES_LOCAL_SCHEMA_VERSION="2"
CC_PLAN_SCHEMA_VERSION="3"
CC_PAIRING_SESSION_SCHEMA_VERSION="1"
CC_GROUNDING_MANIFEST_SCHEMA_VERSION="1"
CC_LEASE_SCHEMA_VERSION="1"
CC_EXECUTION_SCHEMA_VERSION="2"
CC_CANDIDATE_SCHEMA_VERSION="1"
CC_CHANGE_SET_SCHEMA_VERSION="1"
CC_VERIFIER_RESULT_SCHEMA_VERSION="1"
CC_HUMAN_ACCEPTANCE_SCHEMA_VERSION="1"
CC_COMPLETION_SCHEMA_VERSION="1"
CC_DELIVERY_SCHEMA_VERSION="1"
CC_KNOWLEDGE_DEBT_SCHEMA_VERSION="1"

# ---------------------------------------------------------------------------
# Diagnostics
# ---------------------------------------------------------------------------

# cc_fail CODE [detail...] -> print reason code to stderr, return 1
cc_fail() {
	cc_code="$1"
	shift 2>/dev/null || true
	if [ "$#" -gt 0 ]; then
		printf '%s %s\n' "$cc_code" "$*" >&2
	else
		printf '%s\n' "$cc_code" >&2
	fi
	return 1
}

# cc_emit KEY VALUE -> structured stdout line
cc_emit() { printf '%s: %s\n' "$1" "$2"; }

cc_now() { date -u +%Y-%m-%dT%H:%M:%SZ; }

# ---------------------------------------------------------------------------
# Safe identifiers and paths
# ---------------------------------------------------------------------------

cc_lower() { printf '%s' "$1" | tr '[:upper:]' '[:lower:]'; }

# cc_safe_id VALUE -> ok if [A-Za-z0-9._-]+ and <=96 chars, no leading dash/dot
cc_safe_id() {
	case "$1" in
		'' ) return 1 ;;
		-* | .* ) return 1 ;;
		*[!A-Za-z0-9._-]* ) return 1 ;;
	esac
	[ "${#1}" -le 96 ]
}

# cc_safe_slug VALUE -> ok if lowercase kebab-case, no dates/underscores/spaces
cc_safe_slug() {
	case "$1" in
		'' ) return 1 ;;
		-* | *- ) return 1 ;;
		*[!a-z0-9-]* ) return 1 ;;
	esac
	[ "${#1}" -le 80 ]
}

# cc_safe_relative PATH -> reject absolute paths, traversal, and NUL/space tricks
cc_safe_relative() {
	case "$1" in
		'' ) return 1 ;;
		/* ) return 1 ;;
		*..* ) return 1 ;;
		*' '* ) return 1 ;;
	esac
	return 0
}

# cc_root_abs ROOT -> physical absolute directory path. Runtime-authored paths
# must not depend on whether a host invokes the CLI with `.` or an absolute root.
cc_root_abs() { (CDPATH= cd "$1" 2>/dev/null && pwd -P); }

# cc_plan_id_valid ID -> NNNN-<kebab-slug>
cc_plan_id_valid() {
	cc_pid="$1"
	case "$cc_pid" in
		[0-9][0-9][0-9][0-9]-* ) : ;;
		* ) return 1 ;;
	esac
	cc_pid_slug=${cc_pid#????-}
	cc_safe_slug "$cc_pid_slug"
}

# ---------------------------------------------------------------------------
# Digests and atomic writes
# ---------------------------------------------------------------------------

cc_digest() {
	if command -v sha256sum >/dev/null 2>&1; then
		printf 'sha256:%s' "$(sha256sum "$1" | cut -d' ' -f1)"
	elif command -v shasum >/dev/null 2>&1; then
		printf 'sha256:%s' "$(shasum -a 256 "$1" | cut -d' ' -f1)"
	else
		printf 'cksum:%s' "$(cksum "$1" | cut -d' ' -f1)"
	fi
}

cc_digest_text() {
	cc_dt_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-dg.XXXXXX") || return 1
	printf '%s' "$1" >"$cc_dt_tmp"
	cc_digest "$cc_dt_tmp"
	rm -f "$cc_dt_tmp"
}

# cc_atomic_write DEST  (content on stdin) -> write via temp + mv
cc_atomic_write() {
	cc_aw_dest="$1"
	cc_aw_dir=$(dirname -- "$cc_aw_dest")
	mkdir -p "$cc_aw_dir" || return 1
	cc_aw_tmp=$(mktemp "$cc_aw_dir/.cc-aw.XXXXXX") || return 1
	cat >"$cc_aw_tmp" || { rm -f "$cc_aw_tmp"; return 1; }
	mv -f "$cc_aw_tmp" "$cc_aw_dest"
}

# ---------------------------------------------------------------------------
# Minimal YAML readers (controlled subset; the engine writes its own records)
# ---------------------------------------------------------------------------

# cc_scalar FILE KEY -> value of a top-level `key: value` scalar (empty if bare)
cc_scalar() {
	[ -f "$1" ] || return 1
	awk -v k="$2" '
		$0 ~ "^" k ":[[:space:]]" { sub("^" k ":[[:space:]]*", ""); print; found=1; exit }
		$0 ~ "^" k ":$" { print ""; found=1; exit }
		END { if (!found) exit 3 }
	' "$1"
}

# cc_list_ids FILE SECTION -> `- id: X` values under a top-level SECTION: block
cc_list_ids() {
	[ -f "$1" ] || return 1
	awk -v sec="$2" '
		/^[A-Za-z_][A-Za-z0-9_]*:/ { in_sec = ($0 ~ "^" sec ":") ; next }
		in_sec && /^  -[[:space:]]*id:[[:space:]]*/ {
			sub("^  -[[:space:]]*id:[[:space:]]*", "")
			gsub(/[[:space:]]+$/, "")
			print
		}
	' "$1"
}

# cc_inline_list "[a, b, c]" -> newline-separated items
cc_inline_list() {
	cc_il_raw="$1"
	cc_il_raw=${cc_il_raw#[}
	cc_il_raw=${cc_il_raw%]}
	printf '%s\n' "$cc_il_raw" | tr ',' '\n' | sed 's/^[[:space:]]*//; s/[[:space:]]*$//' | sed '/^$/d'
}

# cc_task_field FILE TASK_ID FIELD -> raw value of a per-task field under tasks:
# Handles both inline (`field: [a, b]` / `field: value`) forms.
cc_task_field() {
	[ -f "$1" ] || return 1
	awk -v want="$2" -v field="$3" '
		/^[A-Za-z_][A-Za-z0-9_]*:/ { in_tasks = ($0 ~ "^tasks:") ; cur=0 ; next }
		in_tasks && /^  -[[:space:]]*id:[[:space:]]*/ {
			v=$0; sub("^  -[[:space:]]*id:[[:space:]]*", "", v); gsub(/[[:space:]]+$/, "", v)
			cur = (v == want) ; next
		}
		in_tasks && cur && $0 ~ "^    " field ":[[:space:]]*" {
			sub("^[[:space:]]+" field ":[[:space:]]*", "")
			gsub(/[[:space:]]+$/, "")
			print; exit
		}
	' "$1"
}

# cc_task_list FILE TASK_ID FIELD -> items of a per-task list field, whether the
# field is authored inline (`field: [a, b]`) or as a block list (`field:` then
# `      - a`). Prints one item per line.
cc_task_list() {
	[ -f "$1" ] || return 1
	awk -v want="$2" -v field="$3" '
		/^[A-Za-z_][A-Za-z0-9_]*:/ { in_tasks = ($0 ~ "^tasks:") ; cur=0 ; infield=0 ; next }
		in_tasks && /^  -[[:space:]]*id:[[:space:]]*/ {
			v=$0; sub("^  -[[:space:]]*id:[[:space:]]*", "", v); gsub(/[[:space:]]+$/, "", v)
			cur = (v == want) ; infield=0 ; next
		}
		in_tasks && cur && $0 ~ "^    " field ":[[:space:]]*" {
			val=$0; sub("^    " field ":[[:space:]]*", "", val); gsub(/[[:space:]]+$/, "", val)
			if (val ~ /^\[/) {
				gsub(/^\[|\]$/, "", val); n=split(val, a, ",")
				for (i=1;i<=n;i++){ gsub(/^[[:space:]]+|[[:space:]]+$/, "", a[i]); if (a[i]!="") print a[i] }
				infield=0
			} else if (val == "") { infield=1 }
			else { print val; infield=0 }
			next
		}
		in_tasks && cur && infield && /^      -[[:space:]]*/ {
			it=$0; sub("^      -[[:space:]]*", "", it); gsub(/[[:space:]]+$/, "", it)
			if (it != "") print it; next
		}
		in_tasks && cur && infield && /^    [A-Za-z]/ { infield=0 }
	' "$1"
}

# cc_task_ids FILE -> all task ids in order
cc_task_ids() { cc_list_ids "$1" "tasks"; }

# cc_plan_repositories FILE -> repository ids declared in the plan header
cc_plan_repositories() { cc_list_ids "$1" "repositories"; }

# cc_plan_repo_paths FILE REPO -> union of bounded paths of every task mapped to REPO
cc_plan_repo_paths() {
	cc_prp_file="$1"; cc_prp_repo="$2"
	for cc_prp_t in $(cc_task_ids "$cc_prp_file"); do
		if cc_task_list "$cc_prp_file" "$cc_prp_t" "repositories" | grep -Fxq "$cc_prp_repo"; then
			cc_task_list "$cc_prp_file" "$cc_prp_t" "paths"
		fi
	done | sort -u | sed '/^$/d'
}

# cc_plan_affected_repositories FILE -> union of every task's repositories
cc_plan_affected_repositories() {
	cc_par_file="$1"
	for cc_par_t in $(cc_task_ids "$cc_par_file"); do
		cc_task_list "$cc_par_file" "$cc_par_t" "repositories"
	done | sort -u | sed '/^$/d'
}

# cc_plan_dependencies FILE -> inter-plan dependency plan ids (one per line)
cc_plan_dependencies() { cc_list_ids "$1" "plan_dependencies"; }

# cc_plan_dep_closure PLANS_DIR PLAN -> every transitive plan_dependencies id of
# PLAN (one per line, sorted). PLANS_DIR is the directory that holds <plan>/plan.yaml
# for every plan. A plan appears in its OWN closure iff it participates in a cycle.
cc_plan_dep_closure() {
	cc_pdc_root="$1"; cc_pdc_start="$2"
	cc_pdc_seen=$(mktemp "${TMPDIR:-/tmp}/cc-cl.XXXXXX") || return 1
	cc_pdc_front=$(mktemp "${TMPDIR:-/tmp}/cc-cl.XXXXXX") || { rm -f "$cc_pdc_seen"; return 1; }
	cc_pdc_next=$(mktemp "${TMPDIR:-/tmp}/cc-cl.XXXXXX") || { rm -f "$cc_pdc_seen" "$cc_pdc_front"; return 1; }
	: >"$cc_pdc_seen"
	if [ -f "$cc_pdc_root/$cc_pdc_start/plan.yaml" ]; then
		cc_plan_dependencies "$cc_pdc_root/$cc_pdc_start/plan.yaml" >"$cc_pdc_front"
	else
		: >"$cc_pdc_front"
	fi
	while [ -s "$cc_pdc_front" ]; do
		: >"$cc_pdc_next"
		while IFS= read -r cc_pdc_d; do
			[ -n "$cc_pdc_d" ] || continue
			grep -Fxq "$cc_pdc_d" "$cc_pdc_seen" 2>/dev/null && continue
			printf '%s\n' "$cc_pdc_d" >>"$cc_pdc_seen"
			[ -f "$cc_pdc_root/$cc_pdc_d/plan.yaml" ] \
				&& cc_plan_dependencies "$cc_pdc_root/$cc_pdc_d/plan.yaml" >>"$cc_pdc_next" || :
		done <"$cc_pdc_front"
		cp "$cc_pdc_next" "$cc_pdc_front"
	done
	sort -u "$cc_pdc_seen"
	rm -f "$cc_pdc_seen" "$cc_pdc_front" "$cc_pdc_next"
}

# cc_plan_is_descendant PLANS_DIR CANDIDATE ANCESTOR -> ok if CANDIDATE transitively
# depends on ANCESTOR (ANCESTOR is in CANDIDATE's dependency closure).
cc_plan_is_descendant() {
	cc_plan_dep_closure "$1" "$2" | grep -Fxq "$3"
}

# ---------------------------------------------------------------------------
# Intent (Context Circuit v1.0, Mechanism 1) — the first-class decision for one
# change. The intent tree mirrors plans/ conventions (stable ids never reused,
# a status-blind archive move, an INDEX catalog) so it reuses the same id, index,
# and archive machinery. contract.yaml is the frozen definition of correct;
# approval freezes contract_digest (INV-INTENT-01) and also confirms the coordinator
# understood the plain ask, which is what lets the tracer read the real code next.
# The runtime owns only the deterministic mechanics — id allocation, structure
# validation, digest freezing, archive/restore, and the plan's derives-from-an-
# approved-intent authorization (INV-INTENT-02). It never authors the criteria,
# reads the code, runs the feasibility check, or decides a tier (INV-RUNTIME-01):
# tracing and the feasibility check are coordinator judgment (intent-feasibility),
# not engine verbs, and scope-safety is settled at delivery (Gate 2).
# ---------------------------------------------------------------------------

# cc_intent_id_valid ID -> i<NNN>-<kebab-slug>
cc_intent_id_valid() {
	cc_iiv_id="$1"
	case "$cc_iiv_id" in
		i[0-9][0-9][0-9]-* ) : ;;
		* ) return 1 ;;
	esac
	cc_iiv_slug=${cc_iiv_id#i???-}
	cc_safe_slug "$cc_iiv_slug"
}

cc_intent_dir() { printf '%s/intent/%s' "$1" "$2"; }

# cc_intent_contract_digest FILE -> digest over the criteria-bearing content only:
# every line except the top-level status and contract_digest lines (which change on
# approval and would otherwise make the digest self-referential). Deterministic:
# identical criteria -> identical digest; any criteria change -> a new digest, which
# is what re-enters Gate 1 (re-approval) and voids candidate evidence.
cc_intent_contract_digest() {
	[ -f "$1" ] || return 1
	cc_icd_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-icd.XXXXXX") || return 1
	awk '/^status:[[:space:]]/{next} /^status:$/{next} /^contract_digest:/{next} {print}' "$1" >"$cc_icd_tmp"
	cc_digest "$cc_icd_tmp"
	rm -f "$cc_icd_tmp"
}

# cc_intent_scope_repos FILE -> repository ids declared under scope.repositories
cc_intent_scope_repos() {
	[ -f "$1" ] || return 1
	awk '
		/^[A-Za-z_]/ { in_scope = ($0 ~ /^scope:/) ; in_repos=0 ; next }
		in_scope && /^  repositories:/ { in_repos=1 ; next }
		in_scope && /^  [A-Za-z]/ { in_repos=0 }
		in_scope && in_repos && /^    -[[:space:]]*id:[[:space:]]*/ {
			v=$0; sub(/^    -[[:space:]]*id:[[:space:]]*/, "", v); gsub(/[[:space:]]+$/, "", v); print v
		}
	' "$1"
}

# cc_intent_scope_paths FILE REPO -> path regions declared for REPO under
# scope.repositories, whether inline ([a, b]) or as a block list. One per line.
cc_intent_scope_paths() {
	[ -f "$1" ] || return 1
	awk -v want="$2" '
		/^[A-Za-z_]/ { in_scope = ($0 ~ /^scope:/) ; in_repos=0 ; cur=0 ; infield=0 ; next }
		in_scope && /^  repositories:/ { in_repos=1 ; next }
		in_scope && /^  [A-Za-z]/ { in_repos=0 }
		in_scope && in_repos && /^    -[[:space:]]*id:[[:space:]]*/ {
			v=$0; sub(/^    -[[:space:]]*id:[[:space:]]*/, "", v); gsub(/[[:space:]]+$/, "", v); cur=(v==want); infield=0; next
		}
		in_scope && in_repos && cur && /^      paths:[[:space:]]*/ {
			val=$0; sub(/^      paths:[[:space:]]*/, "", val); gsub(/[[:space:]]+$/, "", val)
			if (val ~ /^\[/) {
				gsub(/^\[|\]$/, "", val); n=split(val, a, ",")
				for (i=1;i<=n;i++){ gsub(/^[[:space:]]+|[[:space:]]+$/, "", a[i]); if (a[i]!="") print a[i] }
				infield=0
			} else if (val == "") { infield=1 }
			else { print val; infield=0 }
			next
		}
		in_scope && in_repos && cur && infield && /^        -[[:space:]]*/ {
			it=$0; sub(/^        -[[:space:]]*/, "", it); gsub(/[[:space:]]+$/, "", it); if (it != "") print it; next
		}
		in_scope && in_repos && cur && infield && /^      [A-Za-z]/ { infield=0 }
	' "$1"
}

# cc_intent_criteria_count FILE -> number of acceptance criteria items. Criteria are
# OUTCOME level (id + statement, in terms a human can approve); the runnable check
# that proves each one is a tracing output carried in the plan, never authored here
# (tracing-and-grounding). Counts one item per criterion `- id:` entry.
cc_intent_criteria_count() {
	[ -f "$1" ] || return 1
	awk '
		/^[A-Za-z_]/ { in_ac = ($0 ~ /^acceptance_criteria:/) ; next }
		in_ac && /^[[:space:]]+-[[:space:]]*id:[[:space:]]*/ { n++ }
		END { print n+0 }
	' "$1"
}

# cc_intent_validate DIR -> confirm intent/<id>/ structure and contract.yaml fields
cc_intent_validate() {
	cc_iv_dir="$1"
	[ -d "$cc_iv_dir" ] || { cc_fail INTENT_DIR_MISSING; return 1; }
	[ -f "$cc_iv_dir/contract.yaml" ] || { cc_fail INTENT_CONTRACT_MISSING; return 1; }
	[ -f "$cc_iv_dir/INTENT.md" ] || { cc_fail INTENT_MD_MISSING; return 1; }
	cc_iv_yaml="$cc_iv_dir/contract.yaml"
	cc_iv_id=$(cc_scalar "$cc_iv_yaml" "intent") || cc_iv_id=""
	cc_intent_id_valid "$cc_iv_id" || { cc_fail INTENT_ID_INVALID "$cc_iv_id"; return 1; }
	[ "$(basename -- "$cc_iv_dir")" = "$cc_iv_id" ] || { cc_fail INTENT_ID_DIR_MISMATCH "$cc_iv_id"; return 1; }
	cc_iv_status=$(cc_scalar "$cc_iv_yaml" "status") || cc_iv_status=""
	case "$cc_iv_status" in
		draft|approved) : ;;
		*) cc_fail INTENT_STATUS_INVALID "$cc_iv_status"; return 1 ;;
	esac
	cc_iv_tier=$(cc_scalar "$cc_iv_yaml" "tier") || cc_iv_tier=""
	case "$cc_iv_tier" in
		explore|standard|critical) : ;;
		*) cc_fail INTENT_TIER_INVALID "$cc_iv_tier"; return 1 ;;
	esac
	# goal and at least one outcome-level acceptance criterion. Criteria are NOT
	# frozen as executable here — the tracer earns the runnable check against the real
	# code after approval (tracing-and-grounding). scope is COARSE and OPTIONAL: it may
	# name no repository at all (a lay human draws almost no paths), so an empty scope
	# is valid — the tracer reports where the change lands and scope-safety is settled
	# at delivery (Gate 2), not by an automated gate here (intent-feasibility).
	cc_iv_goal=$(cc_scalar "$cc_iv_yaml" "goal") || cc_iv_goal=""
	[ -n "$cc_iv_goal" ] || { cc_fail INTENT_GOAL_MISSING; return 1; }
	cc_iv_ncrit=$(cc_intent_criteria_count "$cc_iv_yaml")
	[ "${cc_iv_ncrit:-0}" -ge 1 ] || { cc_fail INTENT_NO_CRITERIA; return 1; }
	cc_emit intent "$cc_iv_id"
	cc_emit status "$cc_iv_status"
	cc_emit tier "$cc_iv_tier"
	return 0
}

# cc_intent_allocate_id ROOT SLUG -> next i<NNN>-slug after the highest ever
# allocated (active + archived), the intent tree's own never-reused sequence.
cc_intent_allocate_id() {
	cc_ia_root="$1"; cc_ia_slug="$2"
	cc_safe_slug "$cc_ia_slug" || { cc_fail INTENT_SLUG_INVALID "$cc_ia_slug"; return 1; }
	cc_ia_max=0
	for cc_ia_d in "$cc_ia_root/intent"/*/ "$cc_ia_root/intent/archive"/*/; do
		[ -d "$cc_ia_d" ] || continue
		cc_ia_base=$(basename -- "$cc_ia_d")
		case "$cc_ia_base" in
			i[0-9][0-9][0-9]-*)
				cc_ia_seq=${cc_ia_base#i}; cc_ia_seq=${cc_ia_seq%%-*}
				cc_ia_seq=$(printf '%s' "$cc_ia_seq" | sed 's/^0*//'); [ -n "$cc_ia_seq" ] || cc_ia_seq=0
				[ "$cc_ia_seq" -gt "$cc_ia_max" ] && cc_ia_max=$cc_ia_seq ;;
		esac
	done
	[ "$cc_ia_max" -lt 999 ] || { cc_fail INTENT_ID_EXHAUSTED; return 1; }
	cc_ia_next=$((cc_ia_max + 1))
	printf 'i%03d-%s\n' "$cc_ia_next" "$cc_ia_slug"
}

cc_intent_index_init() {
	printf '# Active intents\n\n| Intent ID | Title | Status | Tier | Goal | Path |\n| --- | --- | --- | --- | --- | --- |\n' \
		| cc_atomic_write "$1/intent/INDEX.md"
}

# cc_intent_index_remove ROOT INTENT -> drop the intent's row from the active index
cc_intent_index_remove() {
	cc_iir_root="$1"; cc_iir_id="$2"
	[ -f "$cc_iir_root/intent/INDEX.md" ] || return 0
	grep -Fv "| $cc_iir_id |" "$cc_iir_root/intent/INDEX.md" | cc_atomic_write "$cc_iir_root/intent/INDEX.md"
}

# cc_intent_index_upsert ROOT INTENT -> add/update the active row
cc_intent_index_upsert() {
	cc_iiu_root="$1"; cc_iiu_id="$2"
	[ -f "$cc_iiu_root/intent/INDEX.md" ] || cc_intent_index_init "$cc_iiu_root"
	cc_iiu_dir="$cc_iiu_root/intent/$cc_iiu_id"
	cc_iiu_title=""; cc_iiu_status=""; cc_iiu_tier=""; cc_iiu_goal=""
	if [ -f "$cc_iiu_dir/contract.yaml" ]; then
		cc_iiu_title=$(cc_scalar "$cc_iiu_dir/contract.yaml" "title")
		cc_iiu_status=$(cc_scalar "$cc_iiu_dir/contract.yaml" "status")
		cc_iiu_tier=$(cc_scalar "$cc_iiu_dir/contract.yaml" "tier")
		cc_iiu_goal=$(cc_scalar "$cc_iiu_dir/contract.yaml" "goal")
	fi
	cc_intent_index_remove "$cc_iiu_root" "$cc_iiu_id"
	printf '| %s | %s | %s | %s | %s | %s |\n' \
		"$cc_iiu_id" "${cc_iiu_title:-}" "${cc_iiu_status:-draft}" "${cc_iiu_tier:-}" "${cc_iiu_goal:-}" "intent/$cc_iiu_id/INTENT.md" \
		>>"$cc_iiu_root/intent/INDEX.md"
	cc_emit index_row "$cc_iiu_id"
	return 0
}

# cc_intent_human_status_sync DIR STATUS -> keep the human-facing status mirror
# aligned with the canonical contract status after a lifecycle transition.
cc_intent_human_status_sync() {
	cc_ihs_file="$1/INTENT.md"; cc_ihs_status="$2"
	[ -f "$cc_ihs_file" ] || { cc_fail INTENT_MD_MISSING; return 1; }
	awk -v s="$cc_ihs_status" '
		!done && /^_Status:[[:space:]]/ { print "_Status: " s "."; done=1; next }
		!done && NR == 1 && /^#[[:space:]]/ { print; print ""; print "_Status: " s "."; done=1; next }
		{ print }
		END { if (!done) print "_Status: " s "." }
	' "$cc_ihs_file" | cc_atomic_write "$cc_ihs_file"
}

# cc_intent_human_status ROOT INTENT PHRASE -> rewrite INTENT.md _Status without
# changing contract.yaml. Used after a feasible tracer so the human file cannot
# stay stale while approval identity remains on the contract.
cc_intent_human_status() {
	cc_ihsr_root="$1"; cc_ihsr_id="$2"; cc_ihsr_phrase="$3"
	[ -n "$cc_ihsr_phrase" ] || { cc_fail INTENT_HUMAN_STATUS_MISSING; return 1; }
	cc_ihsr_flat=$(printf '%s' "$cc_ihsr_phrase" | tr -d '\n\r')
	[ "$cc_ihsr_flat" = "$cc_ihsr_phrase" ] || { cc_fail INTENT_HUMAN_STATUS_INVALID; return 1; }
	cc_ihsr_dir=$(cc_intent_dir "$cc_ihsr_root" "$cc_ihsr_id") || return 1
	cc_intent_human_status_sync "$cc_ihsr_dir" "$cc_ihsr_phrase" || return 1
	cc_emit intent "$cc_ihsr_id"
	cc_emit human_status "$cc_ihsr_phrase"
	return 0
}

# cc_intent_approve ROOT INTENT -> draft->approved (Gate 1); freeze contract_digest
cc_intent_approve() {
	cc_iap_root="$1"; cc_iap_id="$2"
	cc_iap_dir=$(cc_intent_dir "$cc_iap_root" "$cc_iap_id")
	cc_intent_validate "$cc_iap_dir" >/dev/null || { cc_fail INTENT_APPROVE_INVALID "$cc_iap_id"; return 1; }
	cc_iap_yaml="$cc_iap_dir/contract.yaml"
	# Tier floor (the one safety-critical automated check, INV-ASSURE-01, fail upward):
	# approving an intent at the Explore tier — which drops the independent verifier —
	# is refused when any risk signal is present in scope. Standard/Critical always
	# keep the verifier.
	cc_iap_tier=$(cc_scalar "$cc_iap_yaml" "tier")
	if [ "$cc_iap_tier" = "explore" ]; then
		cc_tier_lower_check "$cc_iap_root" "$cc_iap_id" explore >/dev/null || { cc_fail INTENT_TIER_UNSAFE "$cc_iap_id"; return 1; }
	fi
	cc_iap_status=$(cc_scalar "$cc_iap_yaml" "status")
	cc_iap_digest=$(cc_intent_contract_digest "$cc_iap_yaml") || { cc_fail INTENT_DIGEST_FAILED; return 1; }
	# A draft approves normally; an already-approved intent may be RE-approved only
	# when its criteria changed (its content digest no longer matches the frozen one)
	# — the "changing the criteria re-enters the gate" rule (INV-INTENT-01). A
	# re-approval with nothing changed is refused, so approval is never a no-op flip.
	if [ "$cc_iap_status" = "approved" ]; then
		cc_iap_frozen=$(cc_scalar "$cc_iap_yaml" "contract_digest") || cc_iap_frozen=""
		[ "$cc_iap_digest" != "$cc_iap_frozen" ] || { cc_fail INTENT_ALREADY_APPROVED "$cc_iap_id"; return 1; }
	elif [ "$cc_iap_status" != "draft" ]; then
		cc_fail INTENT_NOT_DRAFT "$cc_iap_status"; return 1
	fi
	# set status approved and freeze contract_digest atomically
	awk -v d="$cc_iap_digest" '
		!sdone && /^status:[[:space:]]/ { print "status: approved"; sdone=1; next }
		/^contract_digest:/ { next }
		{ print }
		END { print "contract_digest: " d }
	' "$cc_iap_yaml" | cc_atomic_write "$cc_iap_yaml" \
		|| { cc_fail INTENT_APPROVAL_RECORD_WRITE_FAILED "$cc_iap_id"; return 1; }
	cc_intent_human_status_sync "$cc_iap_dir" approved \
		|| { cc_fail INTENT_APPROVAL_HUMAN_STATUS_WRITE_FAILED "$cc_iap_id"; return 1; }
	cc_intent_index_upsert "$cc_iap_root" "$cc_iap_id" >/dev/null
	cc_emit intent "$cc_iap_id"
	cc_emit status approved
	cc_emit contract_digest "$cc_iap_digest"
	return 0
}

# cc_intent_archive ROOT INTENT -> move intent/ID -> intent/archive/ID; status-blind
cc_intent_archive() {
	cc_iar_root="$1"; cc_iar_id="$2"
	cc_intent_id_valid "$cc_iar_id" || { cc_fail INTENT_ID_INVALID "$cc_iar_id"; return 1; }
	cc_iar_src="$cc_iar_root/intent/$cc_iar_id"
	cc_iar_dst="$cc_iar_root/intent/archive/$cc_iar_id"
	[ -d "$cc_iar_src" ] || { cc_fail INTENT_ARCHIVE_SOURCE_MISSING "$cc_iar_id"; return 1; }
	[ -e "$cc_iar_dst" ] && { cc_fail INTENT_ARCHIVE_TARGET_COLLISION "$cc_iar_id"; return 1; }
	cc_plan_org_lock "$cc_iar_root" || return 1
	mkdir -p "$cc_iar_root/intent/archive"
	if mv "$cc_iar_src" "$cc_iar_dst" 2>/dev/null; then
		if cc_intent_index_remove "$cc_iar_root" "$cc_iar_id"; then
			cc_plan_org_unlock "$cc_iar_root"
			cc_emit archived "$cc_iar_id"
			return 0
		fi
		mv "$cc_iar_dst" "$cc_iar_src" 2>/dev/null || true
		cc_plan_org_unlock "$cc_iar_root"
		cc_fail INTENT_ARCHIVE_INDEX_FAILED "$cc_iar_id"; return 1
	fi
	cc_plan_org_unlock "$cc_iar_root"
	cc_fail INTENT_ARCHIVE_MOVE_FAILED "$cc_iar_id"; return 1
}

# cc_intent_restore ROOT INTENT -> move intent/archive/ID -> intent/ID; re-add index
cc_intent_restore() {
	cc_ire_root="$1"; cc_ire_id="$2"
	cc_intent_id_valid "$cc_ire_id" || { cc_fail INTENT_ID_INVALID "$cc_ire_id"; return 1; }
	cc_ire_src="$cc_ire_root/intent/archive/$cc_ire_id"
	cc_ire_dst="$cc_ire_root/intent/$cc_ire_id"
	[ -d "$cc_ire_src" ] || { cc_fail INTENT_RESTORE_SOURCE_MISSING "$cc_ire_id"; return 1; }
	[ -e "$cc_ire_dst" ] && { cc_fail INTENT_RESTORE_TARGET_COLLISION "$cc_ire_id"; return 1; }
	cc_plan_org_lock "$cc_ire_root" || return 1
	if mv "$cc_ire_src" "$cc_ire_dst" 2>/dev/null; then
		if cc_intent_index_upsert "$cc_ire_root" "$cc_ire_id" >/dev/null; then
			cc_plan_org_unlock "$cc_ire_root"
			cc_emit restored "$cc_ire_id"
			return 0
		fi
		mv "$cc_ire_dst" "$cc_ire_src" 2>/dev/null || true
		cc_plan_org_unlock "$cc_ire_root"
		cc_fail INTENT_RESTORE_INDEX_FAILED "$cc_ire_id"; return 1
	fi
	cc_plan_org_unlock "$cc_ire_root"
	cc_fail INTENT_RESTORE_MOVE_FAILED "$cc_ire_id"; return 1
}

# cc_intent_authorized ROOT PLAN -> the plan's derives-from-an-approved-intent
# authorization (INV-INTENT-02 reworked, INV-EXEC-01). It confirms the plan validates,
# names a parent intent, and that intent is APPROVED with its criteria unchanged since
# approval (current contract digest == frozen contract_digest). It does NOT compare
# scope or path regions against the intent: v1.0 has no automated scope gate — the
# tracer reports where the change lands, the coordinator's feasibility check surfaces a
# required change beyond a bound scope, and scope-safety is settled at delivery
# (Gate 2, INV-DELIVER-01). A criteria change after approval re-enters Gate 1 (the
# frozen digest no longer matches) and also voids candidate evidence (INV-CANDIDATE-01).
# Fails upward: an invalid plan, a missing/unapproved intent, or a criteria drift is
# unauthorized. Emits `authorized: yes` (return 0) or `authorized: no` + a reason
# (return 1) so a caller re-gates rather than proceeds.
cc_intent_authorized() {
	cc_au_root="$1"; cc_au_plan="$2"
	cc_au_pdir="$cc_au_root/plans/$cc_au_plan"
	[ -f "$cc_au_pdir/plan.yaml" ] || { cc_fail AUTH_PLAN_MISSING "$cc_au_plan"; return 1; }
	cc_plan_validate "$cc_au_pdir" >/dev/null || { cc_emit authorized no; cc_emit reason INVALID_PLAN; cc_fail AUTH_PLAN_INVALID "$cc_au_plan"; return 1; }
	cc_au_intent=$(cc_scalar "$cc_au_pdir/plan.yaml" "intent") || cc_au_intent=""
	[ -n "$cc_au_intent" ] || { cc_emit authorized no; cc_emit reason NO_INTENT; cc_fail AUTH_PLAN_NO_INTENT "$cc_au_plan"; return 1; }
	cc_au_contract="$(cc_intent_dir "$cc_au_root" "$cc_au_intent")/contract.yaml"
	[ -f "$cc_au_contract" ] || { cc_emit authorized no; cc_emit reason INTENT_MISSING; cc_fail AUTH_INTENT_MISSING "$cc_au_intent"; return 1; }
	cc_au_istatus=$(cc_scalar "$cc_au_contract" "status")
	[ "$cc_au_istatus" = "approved" ] || { cc_emit authorized no; cc_emit reason INTENT_NOT_APPROVED; cc_fail AUTH_INTENT_NOT_APPROVED "$cc_au_intent"; return 1; }
	# criteria change: current digest must equal the frozen one, else re-gate to Gate 1
	cc_au_frozen=$(cc_scalar "$cc_au_contract" "contract_digest") || cc_au_frozen=""
	cc_au_now=$(cc_intent_contract_digest "$cc_au_contract") || cc_au_now=""
	if [ -z "$cc_au_frozen" ] || [ "$cc_au_now" != "$cc_au_frozen" ]; then
		cc_emit authorized no; cc_emit reason CRITERIA_CHANGED
		cc_fail AUTH_CRITERIA_CHANGED "$cc_au_intent"; return 1
	fi
	cc_emit authorized yes
	cc_emit intent "$cc_au_intent"
	return 0
}

# ---------------------------------------------------------------------------
# Workspace
# ---------------------------------------------------------------------------

# cc_workspace_validate ROOT -> confirm a Cc workspace root
cc_workspace_validate() {
	cc_ws_root="$1"
	[ -n "$cc_ws_root" ] || { cc_fail WORKSPACE_ROOT_MISSING; return 1; }
	[ -d "$cc_ws_root" ] || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$cc_ws_root"; return 1; }
	[ -f "$cc_ws_root/workspace.yaml" ] || { cc_fail WORKSPACE_IDENTITY_MISSING; return 1; }
	cc_ws_name=$(cc_scalar "$cc_ws_root/workspace.yaml" "workspace") || cc_ws_name=""
	if [ -z "$cc_ws_name" ]; then
		# tolerate nested `workspace:\n  name:` identity form
		cc_ws_name=$(awk '/^workspace:$/{f=1;next} f&&/^[[:space:]]+name:/{sub("^[[:space:]]+name:[[:space:]]*","");print;exit} /^[A-Za-z]/{f=0}' "$cc_ws_root/workspace.yaml")
	fi
	[ -d "$cc_ws_root/plans" ] || { cc_fail WORKSPACE_PLANS_MISSING; return 1; }
	cc_emit workspace "${cc_ws_name:-unknown}"
	cc_emit runtime_version "$CC_RUNTIME_VERSION"
	return 0
}

# cc_workspace_init ROOT -> create the minimal deterministic structure (idempotent)
cc_workspace_init() {
	cc_wi_root="$1"
	[ -n "$cc_wi_root" ] || { cc_fail WORKSPACE_ROOT_MISSING; return 1; }
	mkdir -p "$cc_wi_root/context/domains" "$cc_wi_root/context/roles" \
		"$cc_wi_root/context/proposals" "$cc_wi_root/sources" \
		"$cc_wi_root/intent/archive" \
		"$cc_wi_root/plans/archive" "$cc_wi_root/.runtime/executions" \
		"$cc_wi_root/.runtime/worktrees" "$cc_wi_root/.runtime/pairing" \
		"$cc_wi_root/.runtime/explore" \
		"$cc_wi_root/.runtime/locks" "$cc_wi_root/.runtime/knowledge-debt" \
		"$cc_wi_root/repositories" || return 1
	# Leftover tool cache is never project material; drop it on init so an
	# existing workspace that still has it no longer keeps it.
	if [ -e "$cc_wi_root/.code-review-graph" ]; then
		rm -rf "$cc_wi_root/.code-review-graph" \
			|| { cc_fail LEFTOVER_DROP_FAILED .code-review-graph; return 1; }
	fi
	[ -f "$cc_wi_root/plans/INDEX.md" ] || cc_plan_index_init "$cc_wi_root"
	[ -f "$cc_wi_root/intent/INDEX.md" ] || cc_intent_index_init "$cc_wi_root"
	[ -f "$cc_wi_root/context/INDEX.md" ] || printf '# Context index\n\nNo accepted context units yet.\n' | cc_atomic_write "$cc_wi_root/context/INDEX.md"
	[ -f "$cc_wi_root/repositories.local.yaml" ] || printf 'schema_version: %s\nbindings: {}\n' "$CC_REPOSITORIES_LOCAL_SCHEMA_VERSION" | cc_atomic_write "$cc_wi_root/repositories.local.yaml"
	cc_emit workspace_init ok
	return 0
}

cc_plan_index_init() {
	printf '# Active plans\n\n| Plan ID | Title | Status | Objective | Repositories | Path |\n| --- | --- | --- | --- | --- | --- |\n' \
		| cc_atomic_write "$1/plans/INDEX.md"
}

# ---------------------------------------------------------------------------
# Repository bindings, base branches, worktrees
# ---------------------------------------------------------------------------

# cc_binding_field ROOT REPO FIELD -> read a field from repositories.local.yaml
cc_binding_field() {
	cc_bf_file="$1/repositories.local.yaml"
	[ -f "$cc_bf_file" ] || return 1
	awk -v repo="$2" -v field="$3" '
		/^bindings:/ { in_b=1; next }
		in_b && /^[[:space:]][[:space:]][A-Za-z0-9._-]+:/ {
			k=$0; sub("^[[:space:]]+", "", k); sub(":.*$", "", k)
			cur=(k==repo); next
		}
		in_b && cur && $0 ~ "^[[:space:]]+" field ":[[:space:]]*" {
			sub("^[[:space:]]+" field ":[[:space:]]*", "")
			gsub(/[[:space:]]+$/, ""); gsub(/^["'"'"']|["'"'"']$/, "")
			print; found=1; exit
		}
		in_b && cur && field == "base_branch" && $0 ~ "^[[:space:]]+anchor_branch:[[:space:]]*" {
			legacy=$0
			sub("^[[:space:]]+anchor_branch:[[:space:]]*", "", legacy)
			gsub(/[[:space:]]+$/, "", legacy); gsub(/^['"'"']|['"'"']$/, "", legacy)
			legacy_found=1; next
		}
		END { if (!found && legacy_found) print legacy }
	' "$cc_bf_file"
}

# cc_binding_branch ROOT REPO -> print the canonical base branch. Existing
# alpha.4 bindings may still use anchor_branch; read that legacy spelling until
# the user explicitly runs repository-binding-migrate. A file containing both
# spellings is rejected by returning non-zero rather than guessing.
cc_binding_branch() {
	cc_bb_base=$(cc_binding_field "$1" "$2" "base_branch") || cc_bb_base=""
	cc_bb_legacy=$(cc_binding_field "$1" "$2" "anchor_branch") || cc_bb_legacy=""
	if [ -n "$cc_bb_base" ] && [ -n "$cc_bb_legacy" ]; then
		[ "$cc_bb_base" = "$cc_bb_legacy" ] || return 1
	fi
	[ -n "$cc_bb_base" ] && { printf '%s\n' "$cc_bb_base"; return 0; }
	[ -n "$cc_bb_legacy" ] && { printf '%s\n' "$cc_bb_legacy"; return 0; }
	return 1
}

# cc_record_base_branch FILE -> print the canonical base branch from a runtime
# repository record. Records written by older releases may still contain the
# anchor_branch spelling; reading it is safe because the branch value has not
# changed, only its name.
cc_record_base_branch() {
	cc_rb_base=$(cc_scalar "$1" "base_branch") || cc_rb_base=""
	cc_rb_legacy=$(cc_scalar "$1" "anchor_branch") || cc_rb_legacy=""
	if [ -n "$cc_rb_base" ] && [ -n "$cc_rb_legacy" ]; then
		[ "$cc_rb_base" = "$cc_rb_legacy" ] || return 1
	fi
	[ -n "$cc_rb_base" ] && { printf '%s\n' "$cc_rb_base"; return 0; }
	[ -n "$cc_rb_legacy" ] && { printf '%s\n' "$cc_rb_legacy"; return 0; }
	return 1
}

# cc_repository_binding_migrate ROOT -> rewrite legacy anchor_branch bindings
	# to base_branch and mark repositories.local.yaml as the current schema. This is explicit
# because repositories.local.yaml is workspace-owned host state; normal runtime
# reads remain non-mutating. Conflicting dual spellings are refused.
cc_repository_binding_migrate() {
	cc_rbm_root="$1"
	cc_rbm_file="$cc_rbm_root/repositories.local.yaml"
	[ -f "$cc_rbm_file" ] || { cc_fail REPOSITORY_BINDING_MISSING; return 1; }
	cc_rbm_schema=$(cc_scalar "$cc_rbm_file" schema_version) || cc_rbm_schema=""
	case "$cc_rbm_schema" in
		1|"$CC_REPOSITORIES_LOCAL_SCHEMA_VERSION") : ;;
		*) cc_fail REPOSITORY_BINDING_SCHEMA_UNSUPPORTED "$cc_rbm_schema"; return 1 ;;
	esac
	cc_rbm_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-binding-migrate.XXXXXX") || return 1
	if awk -v target_schema="$CC_REPOSITORIES_LOCAL_SCHEMA_VERSION" '
		BEGIN { in_bindings=0; in_binding=0; has_base=0; has_legacy=0; conflict=0 }
		/^schema_version:[[:space:]]*/ { print "schema_version: " target_schema; next }
		/^bindings:/ {
			if (in_binding && has_base && has_legacy) conflict=1
			in_bindings=1; in_binding=0; has_base=0; has_legacy=0; print; next
		}
		in_bindings && /^  [A-Za-z0-9._-]+:/ {
			if (in_binding && has_base && has_legacy) conflict=1
			in_binding=1; has_base=0; has_legacy=0; print; next
		}
		in_bindings && in_binding && /^    base_branch:[[:space:]]*/ {
			has_base=1; print; next
		}
		in_bindings && in_binding && /^    anchor_branch:[[:space:]]*/ {
			if (has_base) conflict=1
			has_legacy=1; sub(/^    anchor_branch:/, "    base_branch:"); print; next
		}
		{ print }
		END {
			if (in_binding && has_base && has_legacy) conflict=1
			if (conflict) exit 2
		}
	' "$cc_rbm_file" >"$cc_rbm_tmp"; then
		:
	else
		cc_rbm_status=$?
		rm -f "$cc_rbm_tmp"
		[ "$cc_rbm_status" -eq 2 ] && { cc_fail REPOSITORY_BASE_BRANCH_CONFLICT; return 1; }
		cc_fail REPOSITORY_BINDING_MIGRATION_FAILED
		return 1
	fi
	cc_atomic_write "$cc_rbm_file" <"$cc_rbm_tmp" || { rm -f "$cc_rbm_tmp"; cc_fail REPOSITORY_BINDING_MIGRATION_FAILED; return 1; }
	rm -f "$cc_rbm_tmp"
	cc_emit migrated repositories.local.yaml
	cc_emit schema_version "$CC_REPOSITORIES_LOCAL_SCHEMA_VERSION"
	return 0
}

# cc_repository_register ROOT ID PATH BASE_BRANCH [CANONICAL_URL] [DEFAULT_BRANCH]
# Records portable logical identity in workspace.yaml and a host-local binding
# in repositories.local.yaml. It does not clone, init, or execute anything.
cc_repository_register() {
	cc_reg_root="$1"; cc_reg_id="$2"; cc_reg_path="$3"; cc_reg_base="$4"
	cc_reg_url="${5:-}"; cc_reg_default="${6:-}"
	cc_safe_id "$cc_reg_id" || { cc_fail REPOSITORY_ID_UNSAFE "$cc_reg_id"; return 1; }
	[ -n "$cc_reg_path" ] || { cc_fail REPOSITORY_PATH_MISSING; return 1; }
	[ -n "$cc_reg_base" ] || { cc_fail REPOSITORY_BASE_BRANCH_UNSET "$cc_reg_id"; return 1; }
	case "$cc_reg_url" in *[Pp]assword@*|*://*:*@*) cc_fail REPOSITORY_URL_HAS_CREDENTIALS; return 1 ;; esac
	# portable identity in workspace.yaml (never a machine path)
	cc_reg_ws="$cc_reg_root/workspace.yaml"
	[ -f "$cc_reg_ws" ] || { cc_fail WORKSPACE_IDENTITY_MISSING; return 1; }
	if cc_plan_repositories "$cc_reg_ws" | grep -Fxq "$cc_reg_id"; then
		: # identity already present; leave it
	else
		# Insert the new identity entry immediately AFTER the `repositories:` line,
		# not at end-of-file: workspace.yaml keeps other top-level blocks (e.g.
		# `identity:`) after `repositories:`, and an EOF append would place the
		# entry in the wrong section where cc_plan_repositories cannot see it.
		awk -v id="$cc_reg_id" -v url="$cc_reg_url" -v defb="$cc_reg_default" '
			/^repositories:[[:space:]]*(\[\])?[[:space:]]*$/ && !ins {
				print "repositories:"
				print "  - id: " id
				if (url  != "") print "    canonical_url: "  url
				if (defb != "") print "    default_branch: " defb
				ins=1; next
			}
			{ print }
			END { if (!ins) { print "repositories:"; print "  - id: " id;
				if (url != "") print "    canonical_url: " url;
				if (defb != "") print "    default_branch: " defb } }
		' "$cc_reg_ws" | cc_atomic_write "$cc_reg_ws"
	fi
	# host-local binding in repositories.local.yaml
	cc_reg_local="$cc_reg_root/repositories.local.yaml"
	if [ ! -f "$cc_reg_local" ] || grep -q '^bindings:[[:space:]]*{}' "$cc_reg_local" 2>/dev/null || ! grep -q '^bindings:' "$cc_reg_local" 2>/dev/null; then
		printf 'schema_version: %s\nbindings:\n' "$CC_REPOSITORIES_LOCAL_SCHEMA_VERSION" | cc_atomic_write "$cc_reg_local"
	else
		cc_reg_schema=$(cc_scalar "$cc_reg_local" schema_version 2>/dev/null) || cc_reg_schema=""
		case "$cc_reg_schema" in
			1) cc_repository_binding_migrate "$cc_reg_root" >/dev/null || return 1 ;;
			"$CC_REPOSITORIES_LOCAL_SCHEMA_VERSION") : ;;
			*) cc_fail REPOSITORY_BINDING_SCHEMA_UNSUPPORTED "$cc_reg_schema"; return 1 ;;
		esac
	fi
	if cc_binding_field "$cc_reg_root" "$cc_reg_id" "path" >/dev/null 2>&1 && [ -n "$(cc_binding_field "$cc_reg_root" "$cc_reg_id" path)" ]; then
		cc_fail REPOSITORY_BINDING_EXISTS "$cc_reg_id"; return 1
	fi
	printf '  %s:\n    path: %s\n    base_branch: %s\n' "$cc_reg_id" "$cc_reg_path" "$cc_reg_base" >>"$cc_reg_local"
	cc_emit repository "$cc_reg_id"
	cc_emit registered ok
	return 0
}

# cc_repo_resolve ROOT REPO -> validate binding; emit path/base_branch; REPOSITORY_* codes
cc_repo_resolve() {
	cc_rr_root="$1"; cc_rr_id="$2"
	cc_safe_id "$cc_rr_id" || { cc_fail REPOSITORY_ID_UNSAFE "$cc_rr_id"; return 1; }
	cc_rr_path=$(cc_binding_field "$cc_rr_root" "$cc_rr_id" "path") || cc_rr_path=""
	[ -n "$cc_rr_path" ] || { cc_fail REPOSITORY_BINDING_MISSING "$cc_rr_id"; return 1; }
	# resolve relative to workspace root
	case "$cc_rr_path" in
		/*) cc_rr_abs="$cc_rr_path" ;;
		*) cc_rr_abs="$cc_rr_root/$cc_rr_path" ;;
	esac
	[ -e "$cc_rr_abs" ] || { cc_fail REPOSITORY_PATH_UNAVAILABLE "$cc_rr_id"; return 1; }
	# reject unsafe symlink escape when path is inside the workspace
	if [ -L "$cc_rr_abs" ]; then
		cc_fail REPOSITORY_PATH_UNSAFE_SYMLINK "$cc_rr_id"; return 1
	fi
	git -C "$cc_rr_abs" rev-parse --git-dir >/dev/null 2>&1 || { cc_fail REPOSITORY_NOT_GIT "$cc_rr_id"; return 1; }
	cc_rr_base_branch=$(cc_binding_branch "$cc_rr_root" "$cc_rr_id" "base_branch") || cc_rr_base_branch=""
	[ -n "$cc_rr_base_branch" ] || { cc_fail REPOSITORY_BASE_BRANCH_UNSET "$cc_rr_id"; return 1; }
	cc_emit repository "$cc_rr_id"
	cc_emit path "$cc_rr_abs"
	cc_emit base_branch "$cc_rr_base_branch"
	return 0
}

# cc_repo_base_commit ROOT REPO -> print the base branch tip commit
cc_repo_base_commit() {
	cc_rac_abs=$(cc_repo_resolve "$1" "$2" | sed -n 's/^path: //p')
	cc_rac_base_branch=$(cc_binding_branch "$1" "$2" "base_branch")
	git -C "$cc_rac_abs" rev-parse --verify "refs/heads/$cc_rac_base_branch" 2>/dev/null \
		|| { cc_fail BASE_BRANCH_MISSING "$2:$cc_rac_base_branch"; return 1; }
}

# cc_repo_clean ROOT REPO -> ok if the base checkout has no uncommitted changes
cc_repo_clean() {
	cc_rc_abs=$(cc_repo_resolve "$1" "$2" | sed -n 's/^path: //p') || return 1
	if [ -n "$(git -C "$cc_rc_abs" status --porcelain 2>/dev/null)" ]; then
		cc_fail REPOSITORY_BASE_DIRTY "$2"; return 1
	fi
	return 0
}

# cc_repository_preflight ROOT PLAN_DIR -> validate all affected repos for a plan
cc_repository_preflight() {
	cc_rp_root="$1"; cc_rp_plandir="$2"
	[ -f "$cc_rp_plandir/plan.yaml" ] || { cc_fail PLAN_YAML_MISSING; return 1; }
	cc_rp_ok=0
	for cc_rp_id in $(cc_plan_affected_repositories "$cc_rp_plandir/plan.yaml"); do
		cc_repo_resolve "$cc_rp_root" "$cc_rp_id" >/dev/null || { cc_fail REPOSITORY_PREFLIGHT_FAILED "$cc_rp_id"; return 1; }
		cc_repo_base_commit "$cc_rp_root" "$cc_rp_id" >/dev/null || { cc_fail BASE_PREFLIGHT_FAILED "$cc_rp_id"; return 1; }
		cc_repo_clean "$cc_rp_root" "$cc_rp_id" || { cc_fail REPOSITORY_BASE_DIRTY "$cc_rp_id"; return 1; }
		cc_rp_ok=$((cc_rp_ok + 1))
	done
	[ "$cc_rp_ok" -gt 0 ] || { cc_fail PLAN_NO_AFFECTED_REPOSITORIES; return 1; }
	cc_emit repositories_ready "$cc_rp_ok"
	return 0
}

# cc_worktree_prepare ROOT PLAN_ID REPO -> create branch + worktree from base tip
cc_worktree_prepare() {
	cc_wp_root="$1"; cc_wp_plan="$2"; cc_wp_id="$3"
	cc_wp_abs=$(cc_repo_resolve "$cc_wp_root" "$cc_wp_id" | sed -n 's/^path: //p') || return 1
	cc_wp_base_branch=$(cc_binding_branch "$cc_wp_root" "$cc_wp_id" "base_branch")
	cc_wp_base=$(git -C "$cc_wp_abs" rev-parse --verify "refs/heads/$cc_wp_base_branch" 2>/dev/null) \
		|| { cc_fail BASE_BRANCH_MISSING "$cc_wp_id"; return 1; }
	cc_wp_branch="cc/$cc_wp_plan/$cc_wp_id"
	cc_wp_tree="$cc_wp_root/.runtime/worktrees/$cc_wp_plan/$cc_wp_id"
	if [ -d "$cc_wp_tree" ]; then
		cc_emit worktree "$cc_wp_tree"
		cc_emit branch "$cc_wp_branch"
		cc_emit base_commit "$cc_wp_base"
		cc_emit worktree_reused true
		return 0
	fi
	mkdir -p "$(dirname -- "$cc_wp_tree")"
	if git -C "$cc_wp_abs" show-ref --verify --quiet "refs/heads/$cc_wp_branch"; then
		git -C "$cc_wp_abs" worktree add "$cc_wp_tree" "$cc_wp_branch" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_wp_id"; return 1; }
	else
		git -C "$cc_wp_abs" worktree add -b "$cc_wp_branch" "$cc_wp_tree" "$cc_wp_base" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_wp_id"; return 1; }
	fi
	cc_emit worktree "$cc_wp_tree"
	cc_emit branch "$cc_wp_branch"
	cc_emit base_commit "$cc_wp_base"
	cc_emit worktree_reused false
	return 0
}

# ---------------------------------------------------------------------------
# Direct collaboration (INV-PAIR-01) — one repository, one fresh pairing
# branch/worktree, and one light resumable pointer. This is not an execution.
# ---------------------------------------------------------------------------

# cc_pair_pointer_file ROOT SESSION -> print the active or closed pointer path.
# Active and closed pointers are mutually exclusive; contradictory state fails.
cc_pair_pointer_file() {
	cc_pp_root="$1"; cc_pp_session="$2"
	cc_safe_slug "$cc_pp_session" || { cc_fail PAIR_SESSION_INVALID "$cc_pp_session"; return 1; }
	cc_pp_dir="$cc_pp_root/.runtime/pairing/$cc_pp_session"
	cc_pp_active="$cc_pp_dir/pointer.yaml"
	cc_pp_closed="$cc_pp_dir/closed.yaml"
	if [ -f "$cc_pp_active" ] && [ -f "$cc_pp_closed" ]; then
		cc_fail PAIR_STATE_CONTRADICTORY "$cc_pp_session"; return 1
	fi
	if [ -f "$cc_pp_active" ]; then printf '%s\n' "$cc_pp_active"; return 0; fi
	if [ -f "$cc_pp_closed" ]; then printf '%s\n' "$cc_pp_closed"; return 0; fi
	cc_fail PAIR_SESSION_MISSING "$cc_pp_session"; return 1
}

# cc_pair_pointer_validate ROOT SESSION POINTER -> validate runtime-authored
# identity and deterministic paths before any later operation consumes them.
cc_pair_pointer_validate() {
	cc_pv_root="$1"; cc_pv_session="$2"; cc_pv_file="$3"
	[ -f "$cc_pv_file" ] || { cc_fail PAIR_POINTER_MISSING "$cc_pv_session"; return 1; }
	cc_pv_schema=$(cc_scalar "$cc_pv_file" schema_version) || cc_pv_schema=""
	[ "$cc_pv_schema" = "$CC_PAIRING_SESSION_SCHEMA_VERSION" ] || { cc_fail PAIR_SCHEMA_UNSUPPORTED "$cc_pv_schema"; return 1; }
	cc_pv_repo=$(cc_scalar "$cc_pv_file" repo) || cc_pv_repo=""
	cc_pv_wt=$(cc_scalar "$cc_pv_file" worktree) || cc_pv_wt=""
	cc_pv_branch=$(cc_scalar "$cc_pv_file" branch) || cc_pv_branch=""
	cc_pv_base=$(cc_scalar "$cc_pv_file" base) || cc_pv_base=""
	cc_safe_id "$cc_pv_repo" || { cc_fail PAIR_REPOSITORY_INVALID "$cc_pv_repo"; return 1; }
	cc_pv_expected_wt="$cc_pv_root/.runtime/explore/$cc_pv_session/$cc_pv_repo"
	cc_pv_legacy_wt="$cc_pv_root/.runtime/worktrees/cc-pair/$cc_pv_session/$cc_pv_repo"
	cc_pv_expected_branch="cc-pair/$cc_pv_session"
	[ "$cc_pv_wt" = "$cc_pv_expected_wt" ] || [ "$cc_pv_wt" = "$cc_pv_legacy_wt" ] \
		|| { cc_fail PAIR_WORKTREE_MISMATCH "$cc_pv_session"; return 1; }
	[ "$cc_pv_branch" = "$cc_pv_expected_branch" ] || { cc_fail PAIR_BRANCH_MISMATCH "$cc_pv_session"; return 1; }
	[ -n "$cc_pv_base" ] || { cc_fail PAIR_BASE_MISSING "$cc_pv_session"; return 1; }
	return 0
}

# cc_pair_begin ROOT REPO SESSION [BASE] -> create a fresh pairing branch and
# worktree from BASE (the connected repository's base tip by default), then
# atomically write the only resumable session state.
cc_pair_begin() {
	cc_pb_root="$1"; cc_pb_repo="$2"; cc_pb_session="$3"; cc_pb_requested_base="${4:-}"
	cc_workspace_validate "$cc_pb_root" >/dev/null || return 1
	cc_pb_root=$(cc_root_abs "$cc_pb_root") || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$cc_pb_root"; return 1; }
	cc_safe_id "$cc_pb_repo" || { cc_fail PAIR_REPOSITORY_INVALID "$cc_pb_repo"; return 1; }
	cc_safe_slug "$cc_pb_session" || { cc_fail PAIR_SESSION_INVALID "$cc_pb_session"; return 1; }
	cc_pb_dir="$cc_pb_root/.runtime/pairing/$cc_pb_session"
	cc_pb_pointer="$cc_pb_dir/pointer.yaml"
	cc_pb_closed="$cc_pb_dir/closed.yaml"
	[ ! -e "$cc_pb_pointer" ] && [ ! -e "$cc_pb_closed" ] \
		|| { cc_fail PAIR_SESSION_EXISTS "$cc_pb_session"; return 1; }
	cc_pb_abs=$(cc_repo_resolve "$cc_pb_root" "$cc_pb_repo" | sed -n 's/^path: //p') || return 1
	cc_pb_base_branch=$(cc_binding_branch "$cc_pb_root" "$cc_pb_repo" base_branch) || cc_pb_base_branch=""
	if [ -n "$cc_pb_requested_base" ]; then
		case "$cc_pb_requested_base" in
			-*|*' '*|*..*) cc_fail PAIR_BASE_INVALID "$cc_pb_requested_base"; return 1 ;;
		esac
		cc_pb_base=$(git -C "$cc_pb_abs" rev-parse --verify "$cc_pb_requested_base^{commit}" 2>/dev/null) \
			|| { cc_fail PAIR_BASE_INVALID "$cc_pb_requested_base"; return 1; }
	else
		cc_pb_base=$(git -C "$cc_pb_abs" rev-parse --verify "refs/heads/$cc_pb_base_branch" 2>/dev/null) \
			|| { cc_fail BASE_BRANCH_MISSING "$cc_pb_repo:$cc_pb_base_branch"; return 1; }
	fi
	cc_pb_branch="cc-pair/$cc_pb_session"
	cc_pb_wt="$cc_pb_root/.runtime/explore/$cc_pb_session/$cc_pb_repo"
	git -C "$cc_pb_abs" show-ref --verify --quiet "refs/heads/$cc_pb_branch" \
		&& { cc_fail PAIR_BRANCH_EXISTS "$cc_pb_branch"; return 1; }
	[ ! -e "$cc_pb_wt" ] || { cc_fail PAIR_WORKTREE_EXISTS "$cc_pb_wt"; return 1; }
	mkdir -p "$(dirname -- "$cc_pb_wt")" "$cc_pb_dir" || { cc_fail PAIR_STATE_CREATE_FAILED "$cc_pb_session"; return 1; }
	git -C "$cc_pb_abs" worktree add -b "$cc_pb_branch" "$cc_pb_wt" "$cc_pb_base" >/dev/null 2>&1 \
		|| { cc_fail PAIR_WORKTREE_CREATE_FAILED "$cc_pb_repo"; return 1; }
	{
		printf 'schema_version: %s\n' "$CC_PAIRING_SESSION_SCHEMA_VERSION"
		printf 'repo: %s\n' "$cc_pb_repo"
		printf 'worktree: %s\n' "$cc_pb_wt"
		printf 'branch: %s\n' "$cc_pb_branch"
		printf 'base: %s\n' "$cc_pb_base"
	} | cc_atomic_write "$cc_pb_pointer" \
		|| { cc_fail PAIR_POINTER_WRITE_FAILED "$cc_pb_session"; return 1; }
	cc_emit session "$cc_pb_session"
	cc_emit repository "$cc_pb_repo"
	cc_emit worktree "$cc_pb_wt"
	cc_emit branch "$cc_pb_branch"
	cc_emit base_commit "$cc_pb_base"
	cc_emit supervision human-supervised
	return 0
}

# cc_pair_inspect ROOT SESSION -> report resumable state without changing it.
cc_pair_inspect() {
	cc_pi_root="$1"; cc_pi_session="$2"
	cc_workspace_validate "$cc_pi_root" >/dev/null || return 1
	cc_pi_root=$(cc_root_abs "$cc_pi_root") || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$cc_pi_root"; return 1; }
	cc_pi_file=$(cc_pair_pointer_file "$cc_pi_root" "$cc_pi_session") || return 1
	cc_pair_pointer_validate "$cc_pi_root" "$cc_pi_session" "$cc_pi_file" || return 1
	cc_pi_repo=$(cc_scalar "$cc_pi_file" repo)
	cc_pi_wt=$(cc_scalar "$cc_pi_file" worktree)
	cc_pi_branch=$(cc_scalar "$cc_pi_file" branch)
	cc_pi_base=$(cc_scalar "$cc_pi_file" base)
	cc_pi_abs=$(cc_repo_resolve "$cc_pi_root" "$cc_pi_repo" | sed -n 's/^path: //p') || return 1
	cc_pi_status=active; [ "$(basename -- "$cc_pi_file")" = closed.yaml ] && cc_pi_status=closed
	cc_pi_resumable=false
	cc_pi_branch_present=false
	cc_pi_worktree_present=false
	cc_pi_dirty=unknown
	if git -C "$cc_pi_abs" show-ref --verify --quiet "refs/heads/$cc_pi_branch"; then cc_pi_branch_present=true; fi
	if [ -d "$cc_pi_wt" ]; then
		cc_pi_worktree_present=true
		if [ -n "$(git -C "$cc_pi_wt" status --porcelain 2>/dev/null)" ]; then cc_pi_dirty=true; else cc_pi_dirty=false; fi
	fi
	if [ "$cc_pi_status" = active ] && [ "$cc_pi_branch_present" = true ] && [ "$cc_pi_worktree_present" = true ]; then
		cc_pi_head=$(git -C "$cc_pi_wt" symbolic-ref --quiet --short HEAD 2>/dev/null) || cc_pi_head=""
		[ "$cc_pi_head" = "$cc_pi_branch" ] && cc_pi_resumable=true
	fi
	cc_emit session "$cc_pi_session"
	cc_emit status "$cc_pi_status"
	cc_emit repository "$cc_pi_repo"
	cc_emit worktree "$cc_pi_wt"
	cc_emit branch "$cc_pi_branch"
	cc_emit base_commit "$cc_pi_base"
	cc_emit branch_present "$cc_pi_branch_present"
	cc_emit worktree_present "$cc_pi_worktree_present"
	cc_emit dirty "$cc_pi_dirty"
	cc_emit resumable "$cc_pi_resumable"
	cc_emit supervision human-supervised
	return 0
}

# cc_pair_close ROOT SESSION -> close only a clean session pointer. The branch
# and worktree remain. A dirty worktree stays active until the human explicitly
# decides whether the worker should commit its changes.
cc_pair_close() {
	cc_pc_root="$1"; cc_pc_session="$2"
	cc_workspace_validate "$cc_pc_root" >/dev/null || return 1
	cc_pc_root=$(cc_root_abs "$cc_pc_root") || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$cc_pc_root"; return 1; }
	cc_safe_slug "$cc_pc_session" || { cc_fail PAIR_SESSION_INVALID "$cc_pc_session"; return 1; }
	cc_pc_dir="$cc_pc_root/.runtime/pairing/$cc_pc_session"
	cc_pc_pointer="$cc_pc_dir/pointer.yaml"
	cc_pc_closed="$cc_pc_dir/closed.yaml"
	if [ -f "$cc_pc_closed" ] && [ ! -f "$cc_pc_pointer" ]; then
		cc_emit session "$cc_pc_session"; cc_emit status closed; cc_emit reused true; return 0
	fi
	[ -f "$cc_pc_pointer" ] || { cc_fail PAIR_SESSION_MISSING "$cc_pc_session"; return 1; }
	[ ! -e "$cc_pc_closed" ] || { cc_fail PAIR_STATE_CONTRADICTORY "$cc_pc_session"; return 1; }
	cc_pair_pointer_validate "$cc_pc_root" "$cc_pc_session" "$cc_pc_pointer" || return 1
	cc_pc_wt=$(cc_scalar "$cc_pc_pointer" worktree)
	[ -d "$cc_pc_wt" ] || { cc_fail PAIR_WORKTREE_MISSING "$cc_pc_session"; return 1; }
	[ -z "$(git -C "$cc_pc_wt" status --porcelain 2>/dev/null)" ] \
		|| { cc_fail PAIR_WORKTREE_DIRTY "$cc_pc_session"; return 1; }
	mv "$cc_pc_pointer" "$cc_pc_closed" || { cc_fail PAIR_CLOSE_FAILED "$cc_pc_session"; return 1; }
	cc_emit session "$cc_pc_session"
	cc_emit status closed
	cc_emit branch_preserved true
	cc_emit worktree_preserved true
	return 0
}

# cc_pair_delivery_targets ROOT SESSION -> report the closed pairing branch and
# current connected base target. Drift blocks; this function never rebases,
# pushes, merges, or opens a pull request.
cc_pair_delivery_targets() {
	cc_pd_root="$1"; cc_pd_session="$2"
	cc_workspace_validate "$cc_pd_root" >/dev/null || return 1
	cc_pd_root=$(cc_root_abs "$cc_pd_root") || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$cc_pd_root"; return 1; }
	cc_safe_slug "$cc_pd_session" || { cc_fail PAIR_SESSION_INVALID "$cc_pd_session"; return 1; }
	cc_pd_dir="$cc_pd_root/.runtime/pairing/$cc_pd_session"
	[ ! -f "$cc_pd_dir/pointer.yaml" ] || { cc_fail PAIR_SESSION_ACTIVE "$cc_pd_session"; return 1; }
	cc_pd_file="$cc_pd_dir/closed.yaml"
	cc_pair_pointer_validate "$cc_pd_root" "$cc_pd_session" "$cc_pd_file" || return 1
	cc_pd_repo=$(cc_scalar "$cc_pd_file" repo)
	cc_pd_wt=$(cc_scalar "$cc_pd_file" worktree)
	cc_pd_branch=$(cc_scalar "$cc_pd_file" branch)
	cc_pd_base=$(cc_scalar "$cc_pd_file" base)
	if [ -d "$cc_pd_wt" ] && [ -n "$(git -C "$cc_pd_wt" status --porcelain 2>/dev/null)" ]; then
		cc_fail PAIR_WORKTREE_DIRTY "$cc_pd_session"; return 1
	fi
	cc_pd_abs=$(cc_repo_resolve "$cc_pd_root" "$cc_pd_repo" | sed -n 's/^path: //p') || return 1
	cc_pd_base_branch=$(cc_binding_branch "$cc_pd_root" "$cc_pd_repo" base_branch) || cc_pd_base_branch=""
	cc_pd_tip=$(git -C "$cc_pd_abs" rev-parse --verify "refs/heads/$cc_pd_branch" 2>/dev/null) \
		|| { cc_fail PAIR_BRANCH_MISSING "$cc_pd_branch"; return 1; }
	cc_pd_base_tip=$(git -C "$cc_pd_abs" rev-parse --verify "refs/heads/$cc_pd_base_branch" 2>/dev/null) \
		|| { cc_fail BASE_BRANCH_MISSING "$cc_pd_repo:$cc_pd_base_branch"; return 1; }
	cc_pd_drift=true
	if git -C "$cc_pd_abs" merge-base --is-ancestor "$cc_pd_base_tip" "$cc_pd_tip" 2>/dev/null; then cc_pd_drift=false; fi
	cc_emit session "$cc_pd_session"
	cc_emit repository "$cc_pd_repo"
	cc_emit source_branch "$cc_pd_branch"
	cc_emit target_branch "$cc_pd_base_branch"
	cc_emit base_commit "$cc_pd_base"
	cc_emit branch_tip "$cc_pd_tip"
	cc_emit base_tip "$cc_pd_base_tip"
	cc_emit drift_detected "$cc_pd_drift"
	cc_emit result_label human-supervised
	if [ "$cc_pd_drift" = true ]; then cc_fail PAIR_BASE_DRIFT "$cc_pd_session"; return 1; fi
	return 0
}

# cc_worktree_drop PATH -> remove a git worktree directory, or rm -rf if it is
# already unregistered. Never called on a live Explore session or live execution.
cc_worktree_drop() {
	cc_wd_wt="$1"
	[ -e "$cc_wd_wt" ] || return 0
	if [ -d "$cc_wd_wt" ] && git -C "$cc_wd_wt" rev-parse --git-dir >/dev/null 2>&1; then
		git -C "$cc_wd_wt" worktree remove --force "$cc_wd_wt" >/dev/null 2>&1 \
			|| rm -rf "$cc_wd_wt"
	else
		rm -rf "$cc_wd_wt"
	fi
}

# cc_runtime_plan_live ROOT PLAN_ID -> 0 when an execution is still in flight
cc_runtime_plan_live() {
	cc_rpl_root="$1"; cc_rpl_plan="$2"
	cc_rpl_dir="$cc_rpl_root/.runtime/executions/$cc_rpl_plan"
	[ -d "$cc_rpl_dir" ] || return 1
	for cc_rpl_ex in "$cc_rpl_dir"/*; do
		[ -f "$cc_rpl_ex/execution.yaml" ] || continue
		cc_rpl_st=$(cc_scalar "$cc_rpl_ex/execution.yaml" status) || continue
		case "$cc_rpl_st" in
			running|verifying|repairing) return 0 ;;
		esac
	done
	return 1
}

# cc_runtime_cleanup ROOT -> explicit leftover cleanup. Removes closed Explore
# worktrees under .runtime/explore/ (and leftover .runtime/worktrees/cc-pair/),
# idle plan-execution worktrees, and the leftover .code-review-graph cache.
# Never deletes a still-live Explore session or an in-flight execution.
# Never runs from pair-close.
cc_runtime_cleanup() {
	cc_rc_root="$1"
	cc_workspace_validate "$cc_rc_root" >/dev/null || return 1
	cc_rc_root=$(cc_root_abs "$cc_rc_root") || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$cc_rc_root"; return 1; }
	if [ -e "$cc_rc_root/.code-review-graph" ]; then
		rm -rf "$cc_rc_root/.code-review-graph" \
			|| { cc_fail LEFTOVER_DROP_FAILED .code-review-graph; return 1; }
		cc_emit leftover_removed .code-review-graph
	fi
	for cc_rc_base in "$cc_rc_root/.runtime/explore" "$cc_rc_root/.runtime/worktrees/cc-pair"; do
		[ -d "$cc_rc_base" ] || continue
		for cc_rc_sess in "$cc_rc_base"/*; do
			[ -d "$cc_rc_sess" ] || continue
			cc_rc_name=$(basename -- "$cc_rc_sess")
			if [ -f "$cc_rc_root/.runtime/pairing/$cc_rc_name/pointer.yaml" ]; then
				cc_emit explore_skipped_live "$cc_rc_name"
				continue
			fi
			for cc_rc_wt in "$cc_rc_sess"/*; do
				[ -e "$cc_rc_wt" ] || continue
				cc_worktree_drop "$cc_rc_wt"
			done
			rmdir "$cc_rc_sess" 2>/dev/null || rm -rf "$cc_rc_sess"
			cc_emit explore_removed "$cc_rc_name"
		done
	done
	if [ -d "$cc_rc_root/.runtime/worktrees" ]; then
		for cc_rc_plan in "$cc_rc_root/.runtime/worktrees"/*; do
			[ -d "$cc_rc_plan" ] || continue
			cc_rc_pid=$(basename -- "$cc_rc_plan")
			[ "$cc_rc_pid" != "cc-pair" ] || continue
			if cc_runtime_plan_live "$cc_rc_root" "$cc_rc_pid"; then
				cc_emit execution_skipped_live "$cc_rc_pid"
				continue
			fi
			for cc_rc_wt in "$cc_rc_plan"/*; do
				[ -e "$cc_rc_wt" ] || continue
				cc_worktree_drop "$cc_rc_wt"
			done
			rmdir "$cc_rc_plan" 2>/dev/null || rm -rf "$cc_rc_plan"
			cc_emit execution_removed "$cc_rc_pid"
		done
	fi
	cc_emit runtime_cleanup ok
	return 0
}

# ---------------------------------------------------------------------------
# Execution bases (INV-CONCURRENCY-02) — base selection for a dependent plan.
# ---------------------------------------------------------------------------

# cc_plan_same_repo_preds ROOT PLAN REPO -> same-repo predecessor plan ids: the
# plan_dependencies of PLAN whose own plan also touches REPO (one per line).
cc_plan_same_repo_preds() {
	cc_srp_root="$1"; cc_srp_plan="$2"; cc_srp_repo="$3"
	cc_srp_pf="$cc_srp_root/plans/$cc_srp_plan/plan.yaml"
	[ -f "$cc_srp_pf" ] || return 0
	for cc_srp_dep in $(cc_plan_dependencies "$cc_srp_pf"); do
		cc_srp_df="$cc_srp_root/plans/$cc_srp_dep/plan.yaml"
		[ -f "$cc_srp_df" ] || continue
		if cc_plan_affected_repositories "$cc_srp_df" | grep -Fxq "$cc_srp_repo"; then
			printf '%s\n' "$cc_srp_dep"
		fi
	done
}

# cc_plan_has_same_repo_pred ROOT PLAN REPO -> ok(0) if PLAN has >=1 same-repo pred.
cc_plan_has_same_repo_pred() {
	[ -n "$(cc_plan_same_repo_preds "$1" "$2" "$3")" ]
}

# cc_base_prepare ROOT PLAN REPO -> base-aware branch + worktree for a dependent
# plan. Selects the base (base tip / single predecessor branch / runtime-authored
# integration merge), keeps the base ref at refs/cc-base/<plan>/<repo>, and detects
# a stale base (predecessor repaired) to rebuild it. Emits the same keys as
# cc_worktree_prepare plus based_on and base_kind. BASE_UNBUILDABLE -> blocked.
cc_base_prepare() {
	cc_bp_root="$1"; cc_bp_plan="$2"; cc_bp_repo="$3"
	cc_bp_abs=$(cc_repo_resolve "$cc_bp_root" "$cc_bp_repo" | sed -n 's/^path: //p') || return 1
	cc_bp_base_branch=$(cc_binding_branch "$cc_bp_root" "$cc_bp_repo" "base_branch")
	cc_bp_base_tip=$(git -C "$cc_bp_abs" rev-parse --verify "refs/heads/$cc_bp_base_branch" 2>/dev/null) \
		|| { cc_fail BASE_BRANCH_MISSING "$cc_bp_repo"; return 1; }
	cc_bp_branch="cc/$cc_bp_plan/$cc_bp_repo"
	cc_bp_tree="$cc_bp_root/.runtime/worktrees/$cc_bp_plan/$cc_bp_repo"
	cc_bp_baseref="refs/cc-base/$cc_bp_plan/$cc_bp_repo"
	# same-repo predecessors and their current branch tips
	cc_bp_preds=""; cc_bp_tips=""; cc_bp_n=0
	for cc_bp_dep in $(cc_plan_same_repo_preds "$cc_bp_root" "$cc_bp_plan" "$cc_bp_repo"); do
		cc_bp_ptip=$(git -C "$cc_bp_abs" rev-parse --verify "refs/heads/cc/$cc_bp_dep/$cc_bp_repo" 2>/dev/null) \
			|| { cc_fail BASE_UNBUILDABLE "$cc_bp_repo:$cc_bp_dep:branch-missing"; return 1; }
		cc_bp_preds="${cc_bp_preds:+$cc_bp_preds }$cc_bp_dep"
		cc_bp_tips="${cc_bp_tips:+$cc_bp_tips }$cc_bp_ptip"
		cc_bp_n=$((cc_bp_n + 1))
	done
	# stale/reuse: if a worktree already exists, keep it only when every current
	# predecessor tip is still reachable from the recorded base; otherwise rebuild.
	if [ -d "$cc_bp_tree" ]; then
		cc_bp_rec=$(git -C "$cc_bp_abs" rev-parse --verify "$cc_bp_baseref" 2>/dev/null) || cc_bp_rec=""
		cc_bp_stale=0
		[ -n "$cc_bp_rec" ] || cc_bp_stale=1
		for cc_bp_t in $cc_bp_tips; do
			git -C "$cc_bp_abs" merge-base --is-ancestor "$cc_bp_t" "$cc_bp_rec" 2>/dev/null || cc_bp_stale=1
		done
		if [ "$cc_bp_stale" -eq 0 ]; then
			cc_bp_based=$(printf '%s' "$cc_bp_preds" | tr ' ' ',' | sed 's/,/, /g')
			cc_emit worktree "$cc_bp_tree"
			cc_emit branch "$cc_bp_branch"
			cc_emit base_commit "$cc_bp_rec"
			cc_emit based_on "[$cc_bp_based]"
			cc_emit base_kind reused
			cc_emit worktree_reused true
			return 0
		fi
		# invalidate a stale base entirely before rebuilding
		git -C "$cc_bp_abs" worktree remove --force "$cc_bp_tree" >/dev/null 2>&1 || rm -rf "$cc_bp_tree"
		git -C "$cc_bp_abs" branch -D "$cc_bp_branch" >/dev/null 2>&1 || :
		git -C "$cc_bp_abs" update-ref -d "$cc_bp_baseref" >/dev/null 2>&1 || :
		git -C "$cc_bp_abs" worktree prune >/dev/null 2>&1 || :
	fi
	mkdir -p "$(dirname -- "$cc_bp_tree")"
	# select and build the base
	if [ "$cc_bp_n" -eq 0 ]; then
		cc_bp_kind=base; cc_bp_start="$cc_bp_base_tip"
		git -C "$cc_bp_abs" worktree add -b "$cc_bp_branch" "$cc_bp_tree" "$cc_bp_start" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_bp_repo"; return 1; }
		cc_bp_base="$cc_bp_base_tip"
	elif [ "$cc_bp_n" -eq 1 ]; then
		cc_bp_kind=stack; cc_bp_start="$cc_bp_tips"
		git -C "$cc_bp_abs" worktree add -b "$cc_bp_branch" "$cc_bp_tree" "$cc_bp_start" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_bp_repo"; return 1; }
		cc_bp_base="$cc_bp_start"
	else
		cc_bp_kind=integration
		git -C "$cc_bp_abs" worktree add -b "$cc_bp_branch" "$cc_bp_tree" "$cc_bp_base_tip" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_bp_repo"; return 1; }
		cc_bp_mbr=""
		for cc_bp_dep in $cc_bp_preds; do
			cc_bp_mbr="${cc_bp_mbr:+$cc_bp_mbr }cc/$cc_bp_dep/$cc_bp_repo"
		done
		cc_bp_deplist=$(printf '%s' "$cc_bp_preds" | tr ' ' ',' | sed 's/,/, /g')
		# runtime-authored integration merge — NOT a worker attempt, NOT delivery.
		if ! git -C "$cc_bp_tree" merge --no-ff \
				-m "cc: integration base $cc_bp_plan (merge $cc_bp_deplist)" \
				$cc_bp_mbr >/dev/null 2>&1; then
			# Preserve the conflicted integration worktree and branch. The caller records
			# BASE_UNBUILDABLE as blocked; retry/recovery must be able to inspect the
			# conflict instead of losing the failed setup state (INV-PRESERVE-01).
			cc_fail BASE_UNBUILDABLE "$cc_bp_repo:integration-conflict"; return 1
		fi
		cc_bp_base=$(git -C "$cc_bp_tree" rev-parse HEAD)
	fi
	git -C "$cc_bp_abs" update-ref "$cc_bp_baseref" "$cc_bp_base" >/dev/null 2>&1 || :
	cc_bp_based=$(printf '%s' "$cc_bp_preds" | tr ' ' ',' | sed 's/,/, /g')
	cc_emit worktree "$cc_bp_tree"
	cc_emit branch "$cc_bp_branch"
	cc_emit base_commit "$cc_bp_base"
	cc_emit based_on "[$cc_bp_based]"
	cc_emit base_kind "$cc_bp_kind"
	cc_emit worktree_reused false
	return 0
}

# ---------------------------------------------------------------------------
# Repository grounding (INV-GROUND-01/02/03) — discover the target repository's
# own agent guidance from the worktree, harden the worktree, and assemble the
# worker brief by deterministic slot substitution of a shipped template. The
# runtime emits DATA (a manifest) and a report-style directive, never a model
# prompt (INV-RUNTIME-01); the brief prose lives in wrapper/adapters/, not here.
# ---------------------------------------------------------------------------

# cc_harden_worktree WORKTREE -> detect the toolchain and report the prepared
# environment. Detection is a table lookup; no network, no per-worktree install
# here (full provisioning is a later phase). Emits environment ready|no-toolchain.
cc_harden_worktree() {
	cc_hw_wt="$1"
	[ -d "$cc_hw_wt" ] || { cc_fail GROUNDING_WORKTREE_MISSING "$cc_hw_wt"; return 1; }
	cc_hw_tc=""
	for cc_hw_lf in package-lock.json pnpm-lock.yaml yarn.lock bun.lockb bun.lock \
		go.mod Cargo.lock Gemfile.lock requirements.txt poetry.lock composer.lock \
		pom.xml build.gradle build.gradle.kts; do
		[ -f "$cc_hw_wt/$cc_hw_lf" ] && { cc_hw_tc="$cc_hw_lf"; break; }
	done
	if [ -n "$cc_hw_tc" ]; then
		cc_emit environment ready
		cc_emit toolchain "$cc_hw_tc"
	else
		cc_emit environment no-toolchain
	fi
	return 0
}

# cc_skill_desc SKILL_FILE -> the first line of a skill's frontmatter description
cc_skill_desc() {
	awk '
		/^---[[:space:]]*$/{ fm++; if(fm>=2) exit; next }
		fm==1 && /^description:[[:space:]]*/{
			d=$0; sub(/^description:[[:space:]]*/,"",d)
			if (d=="" || d==">" || d=="|" || d==">-" || d=="|-"){ folded=1; next }
			print d; exit
		}
		fm==1 && folded && /^[[:space:]]+/{ line=$0; sub(/^[[:space:]]+/,"",line); if(line!=""){ print line; exit } }
		fm==1 && folded && /^[^[:space:]-]/{ exit }
	' "$1"
}

# cc_discover_repo_grounding WORKTREE [REPO] -> print the grounding manifest for
# the worktree: the agent-guidance files present, the skills (name + description),
# and the prepared-environment status. Deterministic scan; data, not prompt text.
cc_discover_repo_grounding() {
	cc_dg_wt="$1"; cc_dg_repo="${2:-}"
	[ -d "$cc_dg_wt" ] || { cc_fail GROUNDING_WORKTREE_MISSING "$cc_dg_wt"; return 1; }
	cc_dg_env=$(cc_harden_worktree "$cc_dg_wt" | sed -n 's/^environment: //p'); [ -n "$cc_dg_env" ] || cc_dg_env=no-toolchain
	printf 'schema_version: %s\n' "$CC_GROUNDING_MANIFEST_SCHEMA_VERSION"
	[ -n "$cc_dg_repo" ] && printf 'repository: %s\n' "$cc_dg_repo" || :
	printf 'worktree: %s\n' "$cc_dg_wt"
	# files — accumulate to a temp so the result is independent of shell word-splitting
	cc_dg_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-dg.XXXXXX") || return 1
	: >"$cc_dg_tmp"
	for cc_dg_f in AGENTS.md CLAUDE.md .github/copilot-instructions.md; do
		[ -f "$cc_dg_wt/$cc_dg_f" ] && printf '  - %s\n' "$cc_dg_f" >>"$cc_dg_tmp" || :
	done
	if [ -d "$cc_dg_wt/.cursor/rules" ] && [ -n "$(ls -A "$cc_dg_wt/.cursor/rules" 2>/dev/null)" ]; then
		printf '  - .cursor/rules/\n' >>"$cc_dg_tmp"
	fi
	if [ -s "$cc_dg_tmp" ]; then printf 'files:\n'; cat "$cc_dg_tmp"; else printf 'files: []\n'; fi
	: >"$cc_dg_tmp"
	if [ -d "$cc_dg_wt/.agents/skills" ]; then
		for cc_dg_s in "$cc_dg_wt"/.agents/skills/*/SKILL.md; do
			[ -f "$cc_dg_s" ] || continue
			cc_dg_nm=$(cc_scalar "$cc_dg_s" name 2>/dev/null); [ -n "$cc_dg_nm" ] || cc_dg_nm=$(basename -- "$(dirname -- "$cc_dg_s")")
			cc_dg_de=$(cc_skill_desc "$cc_dg_s"); [ -n "$cc_dg_de" ] || cc_dg_de="(no description)"
			printf '  - name: %s\n    description: %s\n' "$cc_dg_nm" "$cc_dg_de" >>"$cc_dg_tmp"
		done
	fi
	if [ -s "$cc_dg_tmp" ]; then printf 'skills:\n'; cat "$cc_dg_tmp"; else printf 'skills: []\n'; fi
	rm -f "$cc_dg_tmp"
	printf 'environment: %s\n' "$cc_dg_env"
	return 0
}

# cc_grounding_directive MANIFEST_FILE -> the worker-facing grounding directive,
# rendered deterministically from the manifest (empty vs non-empty variants).
cc_grounding_directive() {
	cc_gd_mf="$1"; [ -f "$cc_gd_mf" ] || { cc_fail GROUNDING_MANIFEST_MISSING "$cc_gd_mf"; return 1; }
	cc_gd_nf=$(awk '/^files:/{f=1;next} f&&/^  - /{c++} /^[a-z]/&&!/^files:/{f=0} END{print c+0}' "$cc_gd_mf")
	cc_gd_ns=$(awk '/^skills:/{s=1;next} s&&/^  - name:/{c++} /^[a-z]/&&!/^skills:/{s=0} END{print c+0}' "$cc_gd_mf")
	if [ "${cc_gd_nf:-0}" -eq 0 ] && [ "${cc_gd_ns:-0}" -eq 0 ]; then
		printf 'No repository agent guidance was discovered (no AGENTS.md, CLAUDE.md, skills, or rules). Ground your work in the plan and Product Knowledge, and follow any patterns already present in the code.\n'
		return 0
	fi
	printf "Before you write any code, read and apply this repository's own agent guidance.\n"
	printf 'It is authoritative on HOW to write code here, within the scope this brief sets.\n\n'
	awk '/^files:/{f=1;next} f&&/^  - /{v=$0;sub(/^  - /,"",v);print "- " v " — read and honor it."} /^[a-z]/&&!/^files:/{f=0}' "$cc_gd_mf"
	if [ "${cc_gd_ns:-0}" -gt 0 ]; then
		printf -- "- Skills — all of this repo's skills are listed below; read and follow only the ones RELEVANT TO YOUR TASK, chosen by description; skip the rest:\n"
		awk '/^skills:/{s=1;next} s&&/^  - name:/{n=$0;sub(/^  - name:[[:space:]]*/,"",n);getline;d=$0;sub(/^    description:[[:space:]]*/,"",d);print "    - " n " — " d} /^[a-z]/&&!/^skills:/{s=0}' "$cc_gd_mf"
	fi
	printf "\nIf anything here conflicts with this brief's scope or safety rules, STOP and report.\n"
	return 0
}

# cc_brief_preflight BRIEF_FILE -> refuse a brief missing or with an unfilled
# repository-grounding slot (INV-GROUND-03).
cc_brief_preflight() {
	[ -f "$1" ] || { cc_fail BRIEF_MISSING "$1"; return 1; }
	grep -q '^## Repository grounding' "$1" || { cc_fail BRIEF_GROUNDING_SLOT_MISSING "$1"; return 1; }
	grep -q '@@GROUNDING@@' "$1" && { cc_fail BRIEF_GROUNDING_UNFILLED "$1"; return 1; }
	cc_emit brief_preflight ok
	return 0
}

# cc_worker_brief_assemble ROOT EXEC_DIR REPO TASK_FOCUS -> assemble the worker
# brief by deterministic slot substitution of the shipped template, filling the
# grounding directive + environment from the recorded manifest and the rest from
# the execution record and plan snapshot. Writes brief-<repo>.md and preflights it.
cc_worker_brief_assemble() {
	cc_wb_root="$1"; cc_wb_edir="$2"; cc_wb_repo="$3"; cc_wb_tf="${4:-Implement the plan.}"
	cc_wb_tpl=""
	[ -f "$cc_wb_root/wrapper/runtime/worker-brief.md" ] && cc_wb_tpl="$cc_wb_root/wrapper/runtime/worker-brief.md"
	[ -z "$cc_wb_tpl" ] && [ -f "$cc_wb_root/wrapper/adapters/worker-brief.md" ] && cc_wb_tpl="$cc_wb_root/wrapper/adapters/worker-brief.md"
	[ -n "$cc_wb_tpl" ] || { cc_fail BRIEF_TEMPLATE_MISSING; return 1; }
	cc_wb_rf="$cc_wb_edir/repositories/$cc_wb_repo.yaml"
	[ -f "$cc_wb_rf" ] || { cc_fail EXECUTION_REPOSITORY_UNKNOWN "$cc_wb_repo"; return 1; }
	cc_wb_pid=$(cc_scalar "$cc_wb_edir/execution.yaml" plan)
	cc_wb_w=$(cc_scalar "$cc_wb_rf" worktree); cc_wb_br=$(cc_scalar "$cc_wb_rf" branch)
	cc_wb_bc=$(cc_scalar "$cc_wb_rf" base_commit); cc_wb_ap=$(cc_scalar "$cc_wb_rf" allowed_paths)
	cc_wb_mf="$cc_wb_edir/grounding/$cc_wb_repo.yaml"
	cc_wb_env=no-toolchain; [ -f "$cc_wb_mf" ] && cc_wb_env=$(cc_scalar "$cc_wb_mf" environment)
	cc_wb_grf="$cc_wb_edir/grounding/$cc_wb_repo.directive.txt"
	if [ -f "$cc_wb_mf" ]; then cc_grounding_directive "$cc_wb_mf" >"$cc_wb_grf"; else printf 'No repository agent guidance was discovered.\n' >"$cc_wb_grf"; fi
	cc_wb_envf="$cc_wb_edir/grounding/$cc_wb_repo.env.txt"
	if [ "$cc_wb_env" = "ready" ]; then
		printf "Ready: dependencies provisioned, commit hooks handled. Do NOT install or modify dependencies. Use the repository's own build/test/lint commands.\n" >"$cc_wb_envf"
	else
		printf 'No dependency toolchain detected. If this plan scaffolds one, create it within the allowed paths; add no dependencies beyond what the plan specifies.\n' >"$cc_wb_envf"
	fi
	cc_wb_snap="$cc_wb_edir/snapshot/plan.yaml"
	cc_wb_out="$cc_wb_edir/brief-$cc_wb_repo.md"
	awk -v pid="$cc_wb_pid" -v wt="$cc_wb_w" -v br="$cc_wb_br" -v bc="$cc_wb_bc" \
	    -v ap="$cc_wb_ap" -v tf="$cc_wb_tf" -v snap="$cc_wb_snap" \
	    -v envf="$cc_wb_envf" -v grf="$cc_wb_grf" '
		/@@ENVIRONMENT@@/{ while((getline l<envf)>0) print l; close(envf); next }
		/@@GROUNDING@@/{ while((getline l<grf)>0) print l; close(grf); next }
		{
			gsub(/\{plan_id\}/,pid); gsub(/\{worktree_path\}/,wt); gsub(/\{branch\}/,br)
			gsub(/\{base_commit\}/,bc); gsub(/\{allowed_paths\}/,ap)
			gsub(/\{plan_snapshot\}/,snap); gsub(/\{task_focus\}/,tf)
			print
		}
	' "$cc_wb_tpl" | cc_atomic_write "$cc_wb_out"
	cc_brief_preflight "$cc_wb_out" >/dev/null || { cc_fail BRIEF_PREFLIGHT_FAILED "$cc_wb_repo"; return 1; }
	cc_emit brief "$cc_wb_out"
	cc_emit grounding_manifest "$cc_wb_mf"
	return 0
}

# ---------------------------------------------------------------------------
# Plan structure, status, and active index
# ---------------------------------------------------------------------------

# cc_plan_validate PLAN_DIR -> validate structure and repository mapping
cc_plan_validate() {
	cc_pv_dir="$1"
	[ -d "$cc_pv_dir" ] || { cc_fail PLAN_DIR_MISSING; return 1; }
	[ -f "$cc_pv_dir/plan.yaml" ] || { cc_fail PLAN_YAML_MISSING; return 1; }
	[ -f "$cc_pv_dir/PLAN.md" ] || { cc_fail PLAN_MD_MISSING; return 1; }
	cc_pv_id=$(cc_scalar "$cc_pv_dir/plan.yaml" "plan") || cc_pv_id=""
	cc_plan_id_valid "$cc_pv_id" || { cc_fail PLAN_ID_INVALID "$cc_pv_id"; return 1; }
	# directory basename must equal the plan id
	[ "$(basename -- "$cc_pv_dir")" = "$cc_pv_id" ] || { cc_fail PLAN_ID_DIR_MISMATCH "$cc_pv_id"; return 1; }
	cc_pv_status=$(cc_scalar "$cc_pv_dir/plan.yaml" "status") || cc_pv_status=""
	case "$cc_pv_status" in
		draft|done) : ;;
		*) cc_fail PLAN_STATUS_INVALID "$cc_pv_status"; return 1 ;;
	esac
	# schema version and inter-plan dependencies (INV-PLAN-05)
	cc_pv_schema=$(cc_scalar "$cc_pv_dir/plan.yaml" "schema_version") || cc_pv_schema=""
	case "$cc_pv_schema" in
		"$CC_PLAN_SCHEMA_VERSION") : ;;
		*) cc_fail PLAN_SCHEMA_UNSUPPORTED "$cc_pv_schema"; return 1 ;;
	esac
	# parent intent (INV-INTENT-02 / INV-PLAN-01). Every plan names its parent intent
	# — the decision it derives from. There is no pre-intent plan. The id form is
	# validated here; cc_intent_authorized confirms the intent is approved (the
	# derives-from-an-approved-intent authorization) as an execution preflight.
	cc_pv_intent=$(cc_scalar "$cc_pv_dir/plan.yaml" "intent") || cc_pv_intent=""
	[ -n "$cc_pv_intent" ] || { cc_fail PLAN_INTENT_REQUIRED "$cc_pv_id"; return 1; }
	cc_intent_id_valid "$cc_pv_intent" || { cc_fail PLAN_INTENT_INVALID "$cc_pv_intent"; return 1; }
	# Explore is the planless, human-supervised route in v1.0. A plan of record
	# begins only after promotion to Standard/Critical; accepting an Explore plan
	# here would create a second execution lifecycle with no independent verifier.
	cc_pv_root=$(cd -- "$cc_pv_dir/../.." 2>/dev/null && pwd) || cc_pv_root=""
	cc_pv_intent_dir=$(cc_intent_dir "$cc_pv_root" "$cc_pv_intent")
	cc_pv_tier=$(cc_scalar "$cc_pv_intent_dir/contract.yaml" "tier" 2>/dev/null) || cc_pv_tier=""
	[ "$cc_pv_tier" != explore ] || { cc_fail PLAN_EXPLORE_PLANLESS "$cc_pv_id"; return 1; }
	cc_pv_deps=$(cc_plan_dependencies "$cc_pv_dir/plan.yaml")
	if [ -n "$cc_pv_deps" ]; then
		cc_pv_plansdir=$(dirname -- "$cc_pv_dir")
		for cc_pv_dep in $cc_pv_deps; do
			[ "$cc_pv_dep" != "$cc_pv_id" ] || { cc_fail PLAN_DEP_SELF "$cc_pv_dep"; return 1; }
			[ -f "$cc_pv_plansdir/$cc_pv_dep/plan.yaml" ] || { cc_fail PLAN_DEP_UNKNOWN "$cc_pv_dep"; return 1; }
		done
		if cc_plan_dep_closure "$cc_pv_plansdir" "$cc_pv_id" | grep -Fxq "$cc_pv_id"; then
			cc_fail PLAN_DEP_CYCLE "$cc_pv_id"; return 1
		fi
	fi
	# every task maps to the plan's single declared repository; deps reference tasks
	cc_pv_repos=$(cc_plan_repositories "$cc_pv_dir/plan.yaml")
	cc_pv_repo_n=0
	for cc_pv_rcount in $cc_pv_repos; do
		cc_pv_repo_n=$((cc_pv_repo_n + 1))
	done
	[ "$cc_pv_repo_n" -eq 1 ] || { cc_fail PLAN_MULTIPLE_REPOSITORIES "$cc_pv_id"; return 1; }
	cc_pv_tasks=$(cc_task_ids "$cc_pv_dir/plan.yaml")
	[ -n "$cc_pv_tasks" ] || { cc_fail PLAN_NO_TASKS; return 1; }
	for cc_pv_t in $cc_pv_tasks; do
		cc_pv_treps=$(cc_task_list "$cc_pv_dir/plan.yaml" "$cc_pv_t" "repositories")
		[ -n "$cc_pv_treps" ] || { cc_fail TASK_NO_REPOSITORY "$cc_pv_t"; return 1; }
		cc_pv_tcount=0
		for cc_pv_tr in $cc_pv_treps; do
			cc_pv_tcount=$((cc_pv_tcount + 1))
			printf '%s\n' "$cc_pv_repos" | grep -Fxq "$cc_pv_tr" \
				|| { cc_fail TASK_UNDECLARED_REPOSITORY "$cc_pv_t:$cc_pv_tr"; return 1; }
		done
		[ "$cc_pv_tcount" -eq 1 ] || { cc_fail TASK_MULTIPLE_REPOSITORIES "$cc_pv_t"; return 1; }
		for cc_pv_dep in $(cc_task_list "$cc_pv_dir/plan.yaml" "$cc_pv_t" "depends_on"); do
			printf '%s\n' "$cc_pv_tasks" | grep -Fxq "$cc_pv_dep" \
				|| { cc_fail TASK_UNKNOWN_DEPENDENCY "$cc_pv_t:$cc_pv_dep"; return 1; }
		done
	done
	# Every declared repository must be represented by a task. Otherwise a plan can
	# declare an extra repository whose scope is never checked as part of execution.
	cc_pv_affected=$(cc_plan_affected_repositories "$cc_pv_dir/plan.yaml")
	for cc_pv_r in $cc_pv_repos; do
		printf '%s\n' "$cc_pv_affected" | grep -Fxq "$cc_pv_r" \
			|| { cc_fail PLAN_REPOSITORY_UNUSED "$cc_pv_r"; return 1; }
	done
	cc_emit plan "$cc_pv_id"
	cc_emit status "$cc_pv_status"
	cc_emit tasks "$(printf '%s' "$cc_pv_tasks" | wc -w | tr -d ' ')"
	return 0
}

# cc_plan_allocate_id ROOT SLUG -> next NNNN-slug after the highest ever allocated
cc_plan_allocate_id() {
	cc_ai_root="$1"; cc_ai_slug="$2"
	cc_safe_slug "$cc_ai_slug" || { cc_fail PLAN_SLUG_INVALID "$cc_ai_slug"; return 1; }
	cc_ai_max=0
	for cc_ai_d in "$cc_ai_root/plans"/*/ "$cc_ai_root/plans/archive"/*/; do
		[ -d "$cc_ai_d" ] || continue
		cc_ai_base=$(basename -- "$cc_ai_d")
		case "$cc_ai_base" in
			[0-9][0-9][0-9][0-9]-*)
				cc_ai_seq=${cc_ai_base%%-*}
				cc_ai_seq=$(printf '%s' "$cc_ai_seq" | sed 's/^0*//'); [ -n "$cc_ai_seq" ] || cc_ai_seq=0
				[ "$cc_ai_seq" -gt "$cc_ai_max" ] && cc_ai_max=$cc_ai_seq ;;
		esac
	done
	cc_ai_next=$((cc_ai_max + 1))
	printf '%04d-%s\n' "$cc_ai_next" "$cc_ai_slug"
}

# cc_plan_index_row_present ROOT PLAN
cc_plan_index_row_present() {
	grep -Fq "| $2 |" "$1/plans/INDEX.md" 2>/dev/null
}

# cc_plan_index_remove ROOT PLAN -> drop the plan's row from the active index
cc_plan_index_remove() {
	cc_ir_root="$1"; cc_ir_plan="$2"
	[ -f "$cc_ir_root/plans/INDEX.md" ] || return 0
	grep -Fv "| $cc_ir_plan |" "$cc_ir_root/plans/INDEX.md" | cc_atomic_write "$cc_ir_root/plans/INDEX.md"
}

# cc_plan_index_upsert ROOT PLAN [TITLE] [STATUS] [OBJECTIVE] -> add/update active row
cc_plan_index_upsert() {
	cc_iu_root="$1"; cc_iu_plan="$2"
	[ -f "$cc_iu_root/plans/INDEX.md" ] || cc_plan_index_init "$cc_iu_root"
	cc_iu_pdir="$cc_iu_root/plans/$cc_iu_plan"
	cc_iu_title="${3:-}"; cc_iu_status="${4:-}"; cc_iu_obj="${5:-}"
	if [ -f "$cc_iu_pdir/plan.yaml" ]; then
		[ -n "$cc_iu_title" ] || cc_iu_title=$(cc_scalar "$cc_iu_pdir/plan.yaml" "title")
		[ -n "$cc_iu_status" ] || cc_iu_status=$(cc_scalar "$cc_iu_pdir/plan.yaml" "status")
		[ -n "$cc_iu_obj" ] || cc_iu_obj=$(cc_scalar "$cc_iu_pdir/plan.yaml" "objective")
		cc_iu_repos=$(cc_plan_repositories "$cc_iu_pdir/plan.yaml" | tr '\n' ' ' | sed 's/ *$//; s/ /, /g')
	fi
	cc_plan_index_remove "$cc_iu_root" "$cc_iu_plan"
	printf '| %s | %s | %s | %s | %s | %s |\n' \
		"$cc_iu_plan" "${cc_iu_title:-}" "${cc_iu_status:-draft}" "${cc_iu_obj:-}" "${cc_iu_repos:-}" "plans/$cc_iu_plan/PLAN.md" \
		>>"$cc_iu_root/plans/INDEX.md"
	cc_emit index_row "$cc_iu_plan"
	return 0
}

# cc_plan_status ROOT PLAN -> current status from plan.yaml
cc_plan_status() {
	cc_scalar "$1/plans/$2/plan.yaml" "status"
}

# cc_plan_set_status PLAN_YAML STATUS -> atomically replace the top-level status
cc_plan_set_status() {
	cc_ss_file="$1"; cc_ss_new="$2"
	awk -v s="$cc_ss_new" '
		!done && /^status:[[:space:]]/ { print "status: " s; done=1; next }
		{ print }
	' "$cc_ss_file" | cc_atomic_write "$cc_ss_file"
}

# ---------------------------------------------------------------------------
# Archive / restore (organization only; never inspects plan or execution status)
# ---------------------------------------------------------------------------

# cc_plan_org_lock ROOT -> acquire the short-lived plan-organization lock (mkdir)
cc_plan_org_lock() {
	cc_ol_dir="$1/.runtime/locks/plan-organization.lock"
	mkdir -p "$1/.runtime/locks"
	if mkdir "$cc_ol_dir" 2>/dev/null; then return 0; fi
	cc_fail PLAN_ORG_LOCK_HELD; return 1
}
cc_plan_org_unlock() { rmdir "$1/.runtime/locks/plan-organization.lock" 2>/dev/null || true; }

# cc_plan_archive ROOT PLAN -> move plans/PLAN -> plans/archive/PLAN; drop index row
cc_plan_archive() {
	cc_ar_root="$1"; cc_ar_plan="$2"
	cc_safe_id "$cc_ar_plan" || { cc_fail PLAN_ID_UNSAFE "$cc_ar_plan"; return 1; }
	cc_ar_src="$cc_ar_root/plans/$cc_ar_plan"
	cc_ar_dst="$cc_ar_root/plans/archive/$cc_ar_plan"
	[ -d "$cc_ar_src" ] || { cc_fail ARCHIVE_SOURCE_MISSING "$cc_ar_plan"; return 1; }
	[ -e "$cc_ar_dst" ] && { cc_fail ARCHIVE_TARGET_COLLISION "$cc_ar_plan"; return 1; }
	cc_plan_org_lock "$cc_ar_root" || return 1
	mkdir -p "$cc_ar_root/plans/archive"
	if mv "$cc_ar_src" "$cc_ar_dst" 2>/dev/null; then
		if cc_plan_index_remove "$cc_ar_root" "$cc_ar_plan"; then
			cc_plan_org_unlock "$cc_ar_root"
			cc_emit archived "$cc_ar_plan"
			return 0
		fi
		# index update failed: roll the move back so files and index stay intact
		mv "$cc_ar_dst" "$cc_ar_src" 2>/dev/null || true
		cc_plan_org_unlock "$cc_ar_root"
		cc_fail ARCHIVE_INDEX_FAILED "$cc_ar_plan"; return 1
	fi
	cc_plan_org_unlock "$cc_ar_root"
	cc_fail ARCHIVE_MOVE_FAILED "$cc_ar_plan"; return 1
}

# cc_plan_restore ROOT PLAN -> move plans/archive/PLAN -> plans/PLAN; re-add index
cc_plan_restore() {
	cc_re_root="$1"; cc_re_plan="$2"
	cc_safe_id "$cc_re_plan" || { cc_fail PLAN_ID_UNSAFE "$cc_re_plan"; return 1; }
	cc_re_src="$cc_re_root/plans/archive/$cc_re_plan"
	cc_re_dst="$cc_re_root/plans/$cc_re_plan"
	[ -d "$cc_re_src" ] || { cc_fail RESTORE_SOURCE_MISSING "$cc_re_plan"; return 1; }
	[ -e "$cc_re_dst" ] && { cc_fail RESTORE_TARGET_COLLISION "$cc_re_plan"; return 1; }
	cc_plan_org_lock "$cc_re_root" || return 1
	if mv "$cc_re_src" "$cc_re_dst" 2>/dev/null; then
		if cc_plan_index_upsert "$cc_re_root" "$cc_re_plan" >/dev/null; then
			cc_plan_org_unlock "$cc_re_root"
			cc_emit restored "$cc_re_plan"
			return 0
		fi
		# index update failed: roll the move back so files and index stay intact
		mv "$cc_re_dst" "$cc_re_src" 2>/dev/null || true
		cc_plan_org_unlock "$cc_re_root"
		cc_fail RESTORE_INDEX_FAILED "$cc_re_plan"; return 1
	fi
	cc_plan_org_unlock "$cc_re_root"
	cc_fail RESTORE_MOVE_FAILED "$cc_re_plan"; return 1
}

# ---------------------------------------------------------------------------
# One-worker execution locks
# ---------------------------------------------------------------------------

# cc_lock_acquire ROOT PLAN OWNER -> exclusive-create worker lock for a plan
cc_lock_acquire() {
	cc_la_root="$1"; cc_la_plan="$2"; cc_la_owner="$3"
	[ -n "$cc_la_owner" ] || { cc_fail LOCK_OWNER_MISSING; return 1; }
	cc_la_dir="$cc_la_root/.runtime/locks/$cc_la_plan.lock"
	mkdir -p "$cc_la_root/.runtime/locks"
	if mkdir "$cc_la_dir" 2>/dev/null; then
		printf 'owner: %s\nplan: %s\nacquired_at: %s\n' "$cc_la_owner" "$cc_la_plan" "$(cc_now)" \
			| cc_atomic_write "$cc_la_dir/owner.yaml" \
			|| { rmdir "$cc_la_dir" 2>/dev/null || :; cc_fail LOCK_RECORD_WRITE_FAILED "$cc_la_plan"; return 1; }
		cc_emit lock acquired
		cc_emit owner "$cc_la_owner"
		return 0
	fi
	cc_la_cur=$(cc_scalar "$cc_la_dir/owner.yaml" "owner" 2>/dev/null)
	if [ "$cc_la_cur" = "$cc_la_owner" ]; then
		cc_emit lock reentrant
		return 0
	fi
	cc_fail LOCK_HELD "$cc_la_cur"; return 1
}

cc_lock_owner() { cc_scalar "$1/.runtime/locks/$2.lock/owner.yaml" "owner" 2>/dev/null; }

# cc_lock_release ROOT PLAN OWNER -> only the matching owner may release
cc_lock_release() {
	cc_lr_dir="$1/.runtime/locks/$2.lock"
	cc_lr_cur=$(cc_scalar "$cc_lr_dir/owner.yaml" "owner" 2>/dev/null)
	[ -d "$cc_lr_dir" ] || { cc_emit lock absent; return 0; }
	[ "$cc_lr_cur" = "$3" ] || { cc_fail LOCK_OWNER_MISMATCH "$cc_lr_cur"; return 1; }
	rm -rf "$cc_lr_dir"
	cc_emit lock released
	return 0
}

# ---------------------------------------------------------------------------
# Path leases (INV-CONCURRENCY-01) — composes with the one-worker lock above.
# A lease reserves (repository, path-region) scope. Records live at
# .runtime/locks/paths/<repo>/<holder>.yaml and are preserved on release.
# ---------------------------------------------------------------------------

# cc_region_overlap A B -> ok(0) when the two path regions overlap:
# equal, one a path-prefix ancestor of the other, or either repository-wide ".".
cc_region_overlap() {
	[ "$1" = "." ] && return 0
	[ "$2" = "." ] && return 0
	[ "$1" = "$2" ] && return 0
	case "$2/" in "$1/"*) return 0 ;; esac
	case "$1/" in "$2/"*) return 0 ;; esac
	return 1
}

cc_lease_dir() { printf '%s/.runtime/locks/paths/%s' "$1" "$2"; }
cc_lease_file() { printf '%s/.runtime/locks/paths/%s/%s.yaml' "$1" "$2" "$3"; }

# cc_lease_regions FILE -> the held regions of a lease record, one per line.
cc_lease_regions() {
	cc_lg_raw=$(cc_scalar "$1" "regions" 2>/dev/null) || cc_lg_raw=""
	cc_inline_list "$cc_lg_raw"
}

# cc_lease_check ROOT REPO PLAN "region ..." -> ok when no non-descendant holder
# reserves an overlapping region in REPO. Own lease and descendants are exempt.
cc_lease_check() {
	cc_lc_root="$1"; cc_lc_repo="$2"; cc_lc_plan="$3"; cc_lc_regions="$4"
	cc_lc_dir=$(cc_lease_dir "$cc_lc_root" "$cc_lc_repo")
	[ -d "$cc_lc_dir" ] || { cc_emit lease free; return 0; }
	for cc_lc_f in "$cc_lc_dir"/*.yaml; do
		[ -f "$cc_lc_f" ] || continue
		cc_lc_holder=$(cc_scalar "$cc_lc_f" "plan" 2>/dev/null) || cc_lc_holder=""
		[ -n "$cc_lc_holder" ] || continue
		[ "$cc_lc_holder" = "$cc_lc_plan" ] && continue           # own lease (reentrant)
		cc_lc_rel=$(cc_scalar "$cc_lc_f" "released_at" 2>/dev/null) || cc_lc_rel=""
		[ -n "$cc_lc_rel" ] && continue                            # released; not held
		# a declared descendant of the holder builds on it, never competes
		cc_plan_is_descendant "$cc_lc_root/plans" "$cc_lc_plan" "$cc_lc_holder" && continue
		for cc_lc_held in $(cc_lease_regions "$cc_lc_f"); do
			for cc_lc_want in $cc_lc_regions; do
				if cc_region_overlap "$cc_lc_held" "$cc_lc_want"; then
					cc_fail LEASE_CONFLICT "$cc_lc_repo:$cc_lc_holder:$cc_lc_held"
					return 1
				fi
			done
		done
	done
	cc_emit lease free
	return 0
}

# cc_lease_acquire ROOT REPO PLAN "region ..." -> reserve regions in REPO for PLAN
# after a conflict check. Re-acquiring PLAN's own lease is idempotent (reentrant).
cc_lease_acquire() {
	cc_la2_root="$1"; cc_la2_repo="$2"; cc_la2_plan="$3"; cc_la2_regions="$4"
	[ -n "$cc_la2_regions" ] || { cc_fail LEASE_NO_REGIONS "$cc_la2_plan"; return 1; }
	cc_lease_check "$cc_la2_root" "$cc_la2_repo" "$cc_la2_plan" "$cc_la2_regions" >/dev/null || return 1
	cc_la2_file=$(cc_lease_file "$cc_la2_root" "$cc_la2_repo" "$cc_la2_plan")
	cc_la2_inline=$(printf '%s' "$cc_la2_regions" | tr ' ' '\n' | sed '/^$/d' | paste -sd',' - 2>/dev/null | sed 's/,/, /g')
	[ -n "$cc_la2_inline" ] || cc_la2_inline=$(printf '%s' "$cc_la2_regions" | tr ' ' ',' | sed 's/,/, /g')
	printf 'schema_version: %s\nplan: %s\nrepository: %s\nregions: [%s]\nacquired_at: %s\nreleased_at:\n' \
		"$CC_LEASE_SCHEMA_VERSION" "$cc_la2_plan" "$cc_la2_repo" "$cc_la2_inline" "$(cc_now)" \
		| cc_atomic_write "$cc_la2_file"
	cc_emit lease acquired
	cc_emit repository "$cc_la2_repo"
	cc_emit holder "$cc_la2_plan"
	return 0
}

# cc_lease_release ROOT REPO PLAN -> mark PLAN's lease in REPO released, preserving
# the record (released_at set; regions kept for inspection).
cc_lease_release() {
	cc_lr2_file=$(cc_lease_file "$1" "$2" "$3")
	[ -f "$cc_lr2_file" ] || { cc_emit lease absent; return 0; }
	awk -v ts="$(cc_now)" '/^released_at:/{print "released_at: " ts; seen=1; next}{print} END{if(!seen)print "released_at: " ts}' \
		"$cc_lr2_file" | cc_atomic_write "$cc_lr2_file"
	cc_emit lease released
	return 0
}

# ---------------------------------------------------------------------------
# Execution records
# ---------------------------------------------------------------------------

cc_execution_dir() { printf '%s/.runtime/executions/%s/%s' "$1" "$2" "$3"; }

# cc_execution_next_id ROOT PLAN -> exec-<plan>-NNN (next sequence)
cc_execution_next_id() {
	cc_ni_base="$1/.runtime/executions/$2"
	cc_ni_max=0
	if [ -d "$cc_ni_base" ]; then
		for cc_ni_d in "$cc_ni_base"/exec-*/; do
			[ -d "$cc_ni_d" ] || continue
			cc_ni_n=$(basename -- "$cc_ni_d"); cc_ni_n=${cc_ni_n##*-}
			cc_ni_n=$(printf '%s' "$cc_ni_n" | sed 's/^0*//'); [ -n "$cc_ni_n" ] || cc_ni_n=0
			[ "$cc_ni_n" -gt "$cc_ni_max" ] && cc_ni_max=$cc_ni_n
		done
	fi
	printf 'exec-%s-%03d' "$2" "$((cc_ni_max + 1))"
}

# cc_execution_begin ROOT PLAN OWNER -> preflight, snapshot, worktrees, record
cc_execution_begin() {
	cc_eb_root="$1"; cc_eb_plan="$2"; cc_eb_owner="$3"
	# Canonicalize the workspace root to an absolute path so every persisted
	# worktree/branch path is cwd-independent — a worker or verifier child that
	# runs from a different directory must still resolve them (a "." root would
	# otherwise store relative worktree paths that break across cwds).
	cc_eb_root=$(CDPATH= cd -- "$cc_eb_root" 2>/dev/null && pwd) || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$1"; return 1; }
	cc_eb_dir="$cc_eb_root/plans/$cc_eb_plan"
	cc_plan_validate "$cc_eb_dir" >/dev/null || { cc_fail EXECUTION_PLAN_INVALID; return 1; }
	# An intent-bearing (v1.0) plan re-confirms its authorization at execution start
	# (INV-INTENT-02 reworked, INV-EXEC-01): the plan derives from an approved intent
	# whose criteria are unchanged since approval. Authorization is derived from the
	# approved intent — there is no separate plan-approval gate and no intermediate
	# plan status, and no automated scope gate (scope-safety is settled at delivery,
	# Gate 2). A plan stays `draft` until it completes (INV-EXEC-01 reworked); an
	# already-completed plan is not re-executed.
	cc_eb_intent=$(cc_scalar "$cc_eb_dir/plan.yaml" "intent") || cc_eb_intent=""
	cc_eb_status=$(cc_scalar "$cc_eb_dir/plan.yaml" "status")
	[ "$cc_eb_status" != "done" ] || { cc_fail EXECUTION_PLAN_DONE "$cc_eb_plan"; return 1; }
	cc_intent_authorized "$cc_eb_root" "$cc_eb_plan" >/dev/null || { cc_fail EXECUTION_UNAUTHORIZED "$cc_eb_plan"; return 1; }
	# Capture the intent's frozen contract digest for this execution's candidate
	# identity (INV-CANDIDATE-01). Recorded once, immutably, so candidate computation
	# never has to chase the live intent file.
	cc_eb_icontract="$(cc_intent_dir "$cc_eb_root" "$cc_eb_intent")/contract.yaml"
	cc_eb_cdigest=$(cc_scalar "$cc_eb_icontract" "contract_digest" 2>/dev/null) || cc_eb_cdigest=""
	[ -n "$cc_eb_cdigest" ] || { cc_fail EXECUTION_INTENT_NOT_FROZEN "$cc_eb_intent"; return 1; }
	cc_eb_tier=$(cc_scalar "$cc_eb_icontract" "tier" 2>/dev/null) || cc_eb_tier=standard
	[ -n "$cc_eb_tier" ] || cc_eb_tier=standard
	cc_repository_preflight "$cc_eb_root" "$cc_eb_dir" >/dev/null || { cc_fail EXECUTION_PREFLIGHT_FAILED; return 1; }
	cc_lock_acquire "$cc_eb_root" "$cc_eb_plan" "$cc_eb_owner" >/dev/null || { cc_fail EXECUTION_LOCK_FAILED; return 1; }
	cc_eb_exec=$(cc_execution_next_id "$cc_eb_root" "$cc_eb_plan")
	cc_eb_edir=$(cc_execution_dir "$cc_eb_root" "$cc_eb_plan" "$cc_eb_exec")
	cc_eb_rev=$(cc_digest "$cc_eb_dir/plan.yaml")
	mkdir -p "$cc_eb_edir/attempts" "$cc_eb_edir/repositories" "$cc_eb_edir/snapshot"
	# Persist the execution record before creating any branch/worktree. If setup is
	# interrupted, recovery can see the owner and the execution remains an honest,
	# blocked/running record instead of an orphaned lock with no evidence.
	printf 'schema_version: %s\nexecution_id: %s\nplan: %s\nintent: %s\ncontract_digest: %s\ntier: %s\nplan_revision: %s\nowner: %s\nstatus: running\nblocked_reason:\nworker_failures: 0\ncurrent_attempt: 0\ncreated_at: %s\nupdated_at: %s\n' \
		"$CC_EXECUTION_SCHEMA_VERSION" "$cc_eb_exec" "$cc_eb_plan" "${cc_eb_intent:-}" "$cc_eb_cdigest" "$cc_eb_tier" "$cc_eb_rev" "$cc_eb_owner" "$(cc_now)" "$(cc_now)" \
		| cc_atomic_write "$cc_eb_edir/execution.yaml" \
		|| { cc_lock_release "$cc_eb_root" "$cc_eb_plan" "$cc_eb_owner" >/dev/null 2>&1 || :; cc_fail EXECUTION_RECORD_WRITE_FAILED "$cc_eb_exec"; return 1; }
	# immutable plan snapshot
	cp "$cc_eb_dir/plan.yaml" "$cc_eb_edir/snapshot/plan.yaml" \
		|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason SNAPSHOT_PLAN_FAILED; cc_fail EXECUTION_SNAPSHOT_FAILED; return 1; }
	cp "$cc_eb_dir/PLAN.md" "$cc_eb_edir/snapshot/PLAN.md" \
		|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason SNAPSHOT_PLAN_MD_FAILED; cc_fail EXECUTION_SNAPSHOT_FAILED; return 1; }
	if [ -d "$cc_eb_dir/tasks" ]; then
		cp -R "$cc_eb_dir/tasks" "$cc_eb_edir/snapshot/tasks" \
			|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason SNAPSHOT_TASKS_FAILED; cc_fail EXECUTION_SNAPSHOT_FAILED; return 1; }
	fi
	# per-repository branch + worktree. A plan with same-repo predecessors is
	# base-aware (INV-CONCURRENCY-02): its base is the predecessor branch (stack)
	# or a runtime-authored integration merge. A plan with no dependency keeps the
	# v0.5 base-tip worktree unchanged. Leases are NOT acquired here; the
	# run-stack loop manages them for stacked executions and overlapping paths.
	for cc_eb_id in $(cc_plan_affected_repositories "$cc_eb_dir/plan.yaml"); do
		cc_eb_based=""
		if cc_plan_has_same_repo_pred "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id"; then
			if ! cc_eb_out=$(cc_base_prepare "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id"); then
				# a base that cannot be built cleanly is a blocked execution, not a
				# worker failure (INV-CONCURRENCY-02); preserve evidence, no worker runs.
				cc_exec_set "$cc_eb_edir" status blocked
				cc_exec_set "$cc_eb_edir" blocked_reason BASE_UNBUILDABLE
				cc_emit execution_id "$cc_eb_exec"
				cc_emit status blocked
				cc_emit blocked_reason BASE_UNBUILDABLE
				cc_fail EXECUTION_BASE_UNBUILDABLE "$cc_eb_id"; return 1
			fi
			cc_eb_based=$(printf '%s' "$cc_eb_out" | sed -n 's/^based_on: //p')
		else
			cc_eb_out=$(cc_worktree_prepare "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id") \
				|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason WORKTREE_SETUP_FAILED; cc_fail EXECUTION_WORKTREE_FAILED "$cc_eb_id"; return 1; }
		fi
		cc_eb_wt=$(printf '%s' "$cc_eb_out" | sed -n 's/^worktree: //p')
		cc_eb_br=$(printf '%s' "$cc_eb_out" | sed -n 's/^branch: //p')
		cc_eb_bc=$(printf '%s' "$cc_eb_out" | sed -n 's/^base_commit: //p')
		cc_eb_base_branch=$(cc_binding_branch "$cc_eb_root" "$cc_eb_id" "base_branch")
		cc_eb_paths=$(cc_plan_repo_paths "$cc_eb_dir/plan.yaml" "$cc_eb_id" | tr '\n' ',' | sed 's/,$//; s/,/, /g')
			{
				printf 'repository: %s\nworktree: %s\nbranch: %s\nbase_branch: %s\nbase_commit: %s\nlatest_commit: %s\nallowed_paths: [%s]\n' \
					"$cc_eb_id" "$cc_eb_wt" "$cc_eb_br" "$cc_eb_base_branch" "$cc_eb_bc" "$cc_eb_bc" "$cc_eb_paths"
				[ -n "$cc_eb_based" ] && printf 'based_on: %s\n' "$cc_eb_based" || :
			} | cc_atomic_write "$cc_eb_edir/repositories/$cc_eb_id.yaml" \
				|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason REPOSITORY_RECORD_WRITE_FAILED; cc_fail EXECUTION_RECORD_WRITE_FAILED "$cc_eb_id"; return 1; }
		# repository grounding (INV-GROUND-01): discover the target repo's own agent
		# guidance from the prepared worktree and record it as execution evidence.
		mkdir -p "$cc_eb_edir/grounding"
		cc_discover_repo_grounding "$cc_eb_wt" "$cc_eb_id" | cc_atomic_write "$cc_eb_edir/grounding/$cc_eb_id.yaml" \
			|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason GROUNDING_DISCOVERY_FAILED; cc_fail EXECUTION_GROUNDING_FAILED "$cc_eb_id"; return 1; }
	done
	printf 'schema_version: %s\nexecution_id: %s\nplan: %s\nintent: %s\ncontract_digest: %s\ntier: %s\nplan_revision: %s\nowner: %s\nstatus: running\nworker_failures: 0\ncurrent_attempt: 0\ncreated_at: %s\nupdated_at: %s\n' \
		"$CC_EXECUTION_SCHEMA_VERSION" "$cc_eb_exec" "$cc_eb_plan" "${cc_eb_intent:-}" "$cc_eb_cdigest" "$cc_eb_tier" "$cc_eb_rev" "$cc_eb_owner" "$(cc_now)" "$(cc_now)" \
		| cc_atomic_write "$cc_eb_edir/execution.yaml" \
		|| { cc_exec_set "$cc_eb_edir" status blocked; cc_exec_set "$cc_eb_edir" blocked_reason EXECUTION_RECORD_WRITE_FAILED; cc_fail EXECUTION_RECORD_WRITE_FAILED "$cc_eb_exec"; return 1; }
	cc_emit execution_id "$cc_eb_exec"
	cc_emit status running
	return 0
}

# cc_exec_set EXEC_DIR KEY VALUE -> update a scalar in execution.yaml (+ touch updated_at)
cc_exec_set() {
	cc_es_dir="$1"; cc_es_k="$2"; cc_es_v="$3"
	awk -v k="$cc_es_k" -v v="$cc_es_v" -v ts="$(cc_now)" '
		$0 ~ "^" k ":[[:space:]]" { print k ": " v; seen=1; next }
		/^updated_at:[[:space:]]/ { print "updated_at: " ts; next }
		{ print }
		END { if (!seen) print k ": " v }
	' "$cc_es_dir/execution.yaml" | cc_atomic_write "$cc_es_dir/execution.yaml" \
		|| { cc_fail EXECUTION_RECORD_UPDATE_FAILED "$cc_es_k"; return 1; }
}

# cc_attempt_norm VALUE -> decimal attempt with leading zeros stripped.
# Folder names are padded (attempts/003/); current_attempt is not (3). Callers
# that take an attempt argument must normalize so 003 and 3 compare equal.
# Leading zeros are stripped before any arithmetic or %d formatting (octal).
cc_attempt_norm() {
	cc_an_raw="$1"
	case "$cc_an_raw" in
		''|*[!0-9]*) cc_fail ATTEMPT_INVALID "$cc_an_raw"; return 1 ;;
	esac
	cc_an_n=$(printf '%s' "$cc_an_raw" | sed 's/^0*//')
	[ -n "$cc_an_n" ] || cc_an_n=0
	printf '%s\n' "$cc_an_n"
}

# cc_attempt_begin EXEC_DIR -> increment current_attempt, create attempt dir
cc_attempt_begin() {
	cc_ab_dir="$1"
	[ -f "$cc_ab_dir/execution.yaml" ] || { cc_fail EXECUTION_RECORD_MISSING; return 1; }
	cc_ab_cur=$(cc_scalar "$cc_ab_dir/execution.yaml" "current_attempt")
	cc_ab_next=$((cc_ab_cur + 1))
	cc_ab_pad=$(printf '%03d' "$cc_ab_next")
	mkdir -p "$cc_ab_dir/attempts/$cc_ab_pad"
	cc_exec_set "$cc_ab_dir" current_attempt "$cc_ab_next" || return 1
	cc_exec_set "$cc_ab_dir" status running || return 1
	printf 'attempt: %s\nstatus: started\nstarted_at: %s\n' "$cc_ab_next" "$(cc_now)" \
		| cc_atomic_write "$cc_ab_dir/attempts/$cc_ab_pad/worker.yaml" \
		|| { cc_fail ATTEMPT_RECORD_WRITE_FAILED "$cc_ab_next"; return 1; }
	cc_emit attempt "$cc_ab_next"
	return 0
}

# cc_worker_commit_record EXEC_DIR REPO [KIND] -> capture worktree HEAD as the commit
cc_worker_commit_record() {
	cc_wc_dir="$1"; cc_wc_repo="$2"; cc_wc_kind="${3:-implementation}"
	cc_wc_rf="$cc_wc_dir/repositories/$cc_wc_repo.yaml"
	[ -f "$cc_wc_rf" ] || { cc_fail EXECUTION_REPOSITORY_UNKNOWN "$cc_wc_repo"; return 1; }
	cc_wc_wt=$(cc_scalar "$cc_wc_rf" "worktree")
	cc_wc_br=$(cc_scalar "$cc_wc_rf" "branch")
	cc_wc_prev=$(cc_scalar "$cc_wc_rf" "latest_commit")
	[ -d "$cc_wc_wt" ] || { cc_fail WORKTREE_MISSING "$cc_wc_repo"; return 1; }
	cc_wc_cur=$(git -C "$cc_wc_wt" rev-parse HEAD 2>/dev/null) || { cc_fail COMMIT_HEAD_UNREADABLE "$cc_wc_repo"; return 1; }
	cc_wc_curbr=$(git -C "$cc_wc_wt" rev-parse --abbrev-ref HEAD 2>/dev/null)
	[ "$cc_wc_curbr" = "$cc_wc_br" ] || { cc_fail COMMIT_WRONG_BRANCH "$cc_wc_repo:$cc_wc_curbr"; return 1; }
	# require a new commit since the last recorded revision (base on the first
	# attempt, latest commit on any repair) so every repair produces a new commit
	[ "$cc_wc_cur" != "$cc_wc_prev" ] || { cc_fail COMMIT_NO_CHANGE "$cc_wc_repo"; return 1; }
	# update latest_commit in the repository record
	awk -v c="$cc_wc_cur" '/^latest_commit:[[:space:]]/{print "latest_commit: " c; next}{print}' "$cc_wc_rf" \
		| cc_atomic_write "$cc_wc_rf" \
		|| { cc_fail EXECUTION_REPOSITORY_RECORD_UPDATE_FAILED "$cc_wc_repo"; return 1; }
	cc_wc_att=$(cc_scalar "$cc_wc_dir/execution.yaml" "current_attempt")
	cc_wc_pad=$(printf '%03d' "$cc_wc_att")
	cc_wc_wf="$cc_wc_dir/attempts/$cc_wc_pad/worker.yaml"
	{
		[ -f "$cc_wc_wf" ] && cat "$cc_wc_wf" || :
		printf 'commit repository=%s revision=%s kind=%s at=%s\n' \
			"$cc_wc_repo" "$cc_wc_cur" "$cc_wc_kind" "$(cc_now)"
	} | cc_atomic_write "$cc_wc_wf" \
		|| { cc_fail ATTEMPT_RECORD_WRITE_FAILED "$cc_wc_repo"; return 1; }
	cc_exec_set "$cc_wc_dir" status verifying || return 1
	cc_emit repository "$cc_wc_repo"
	cc_emit commit "$cc_wc_cur"
	cc_emit kind "$cc_wc_kind"
	return 0
}

# cc_worker_handoff_record EXEC_DIR FILE -> store the worker handoff
cc_worker_handoff_record() {
	[ -f "$2" ] || { cc_fail HANDOFF_FILE_MISSING; return 1; }
	cat "$2" | cc_atomic_write "$1/handoff.md" || { cc_fail HANDOFF_WRITE_FAILED; return 1; }
	cc_emit handoff recorded
	return 0
}

# cc_verifier_prepare EXEC_DIR -> validate latest revisions and read-only scope
cc_verifier_prepare() {
	cc_vp_dir="$1"
	[ -f "$cc_vp_dir/execution.yaml" ] || { cc_fail VERIFIER_EXECUTION_MISSING; return 1; }
	cc_vp_att=$(cc_scalar "$cc_vp_dir/execution.yaml" current_attempt 2>/dev/null) || cc_vp_att=0
	[ "$cc_vp_att" -gt 0 ] || { cc_fail VERIFIER_NO_ATTEMPT; return 1; }
	cc_vp_pad=$(printf '%03d' "$cc_vp_att")
	[ -f "$cc_vp_dir/attempts/$cc_vp_pad/worker.yaml" ] || { cc_fail VERIFIER_WORKER_RECORD_MISSING; return 1; }
	for cc_vp_rf in "$cc_vp_dir"/repositories/*.yaml; do
		[ -f "$cc_vp_rf" ] || { cc_fail VERIFIER_REPOSITORY_RECORD_MISSING; return 1; }
		cc_vp_latest=$(cc_scalar "$cc_vp_rf" "latest_commit")
		cc_vp_base=$(cc_scalar "$cc_vp_rf" "base_commit")
		[ "$cc_vp_latest" != "$cc_vp_base" ] || { cc_fail VERIFIER_NO_WORKER_COMMIT "$(basename "$cc_vp_rf" .yaml)"; return 1; }
	done
	cc_vp_candidate=$(cc_candidate_id "$cc_vp_dir") || return 1
	printf 'attempt: %s\ncandidate_id: %s\nread_only: true\nprepared_at: %s\n' \
		"$cc_vp_att" "$cc_vp_candidate" "$(cc_now)" \
		| cc_atomic_write "$cc_vp_dir/attempts/$cc_vp_pad/verifier-scope.yaml" \
		|| { cc_fail VERIFIER_SCOPE_WRITE_FAILED; return 1; }
	cc_emit verifier_scope read-only
	cc_emit candidate_id "$cc_vp_candidate"
	return 0
}

# cc_verifier_result_record EXEC_DIR ATTEMPT OUTCOME [--wrote-products] -> record + count
cc_verifier_result_record() {
	cc_vr_dir="$1"; cc_vr_att="$2"; cc_vr_out="$3"; cc_vr_flag="${4:-}"
	[ -f "$cc_vr_dir/execution.yaml" ] || { cc_fail VERIFIER_EXECUTION_MISSING; return 1; }
	cc_vr_att=$(cc_attempt_norm "$cc_vr_att") || return 1
	cc_vr_current=$(cc_scalar "$cc_vr_dir/execution.yaml" current_attempt 2>/dev/null) || cc_vr_current=0
	cc_vr_current=$(cc_attempt_norm "$cc_vr_current") || return 1
	[ "$cc_vr_att" = "$cc_vr_current" ] || { cc_fail VERIFIER_ATTEMPT_NOT_CURRENT "$cc_vr_att"; return 1; }
	cc_vr_pad=$(printf '%03d' "$cc_vr_att")
	[ -f "$cc_vr_dir/attempts/$cc_vr_pad/verifier-scope.yaml" ] \
		|| { cc_fail VERIFIER_SCOPE_NOT_PREPARED "$cc_vr_att"; return 1; }
	case "$cc_vr_out" in
		passed|failed|blocked|waived) : ;;
		*) cc_fail VERIFIER_OUTCOME_INVALID "$cc_vr_out"; return 1 ;;
	esac
	# read-only enforcement: an explicit product write is rejected outright
	if [ "$cc_vr_flag" = "--wrote-products" ]; then
		cc_fail VERIFIER_WRITE_REJECTED; return 1
	fi
	# structural read-only enforcement: branch tips must be unchanged since worker commit
	for cc_vr_rf in "$cc_vr_dir"/repositories/*.yaml; do
		[ -f "$cc_vr_rf" ] || { cc_fail VERIFIER_REPOSITORY_RECORD_MISSING; return 1; }
		cc_vr_wt=$(cc_scalar "$cc_vr_rf" "worktree")
		cc_vr_rec=$(cc_scalar "$cc_vr_rf" "latest_commit")
		[ -d "$cc_vr_wt" ] || { cc_fail VERIFIER_WORKTREE_MISSING "$(basename "$cc_vr_rf" .yaml)"; return 1; }
		cc_vr_head=$(git -C "$cc_vr_wt" rev-parse HEAD 2>/dev/null)
		if [ -z "$cc_vr_head" ] || [ "$cc_vr_head" != "$cc_vr_rec" ]; then
			cc_fail VERIFIER_MODIFIED_PRODUCT "$(basename "$cc_vr_rf" .yaml)"; return 1
		fi
	done
	# candidate binding (INV-CANDIDATE-01): the result names the candidate it
	# observed, so any later commit or criteria change voids it by construction.
	cc_vr_cand=$(cc_candidate_id "$cc_vr_dir") || return 1
	cc_vr_scope_cand=$(cc_scalar "$cc_vr_dir/attempts/$cc_vr_pad/verifier-scope.yaml" candidate_id 2>/dev/null) || cc_vr_scope_cand=""
	[ -n "$cc_vr_scope_cand" ] && [ "$cc_vr_scope_cand" = "$cc_vr_cand" ] \
		|| { cc_fail VERIFIER_SCOPE_STALE; return 1; }
	mkdir -p "$cc_vr_dir/attempts/$cc_vr_pad"
	printf 'schema_version: %s\nattempt: %s\noutcome: %s\nread_only: true\ncandidate_id: %s\nchecked_at: %s\n' \
		"$CC_VERIFIER_RESULT_SCHEMA_VERSION" "$cc_vr_att" "$cc_vr_out" "$cc_vr_cand" "$(cc_now)" \
		| cc_atomic_write "$cc_vr_dir/attempts/$cc_vr_pad/verifier.yaml" \
		|| { cc_fail VERIFIER_RESULT_WRITE_FAILED; return 1; }
	if [ "$cc_vr_out" = "passed" ]; then
		cc_exec_set "$cc_vr_dir" status verified || return 1
		cc_exec_set "$cc_vr_dir" verified_candidate "$cc_vr_cand" || return 1
		cc_emit outcome passed
		cc_emit status verified
		cc_emit candidate_id "$cc_vr_cand"
		return 0
	fi
	if [ "$cc_vr_out" = "blocked" ]; then
		cc_exec_set "$cc_vr_dir" status blocked || return 1
		cc_emit outcome blocked
		cc_emit status blocked
		return 0
	fi
	if [ "$cc_vr_out" = "waived" ]; then
		# A human limitation decision: non-passing, but not a worker fault.
		# It never satisfies verification and never increments the failure counter.
		cc_exec_set "$cc_vr_dir" status blocked || return 1
		cc_emit outcome waived
		cc_emit status blocked
		return 0
	fi
	# failed: the independent verifier rejected the worker result
	cc_vr_wf=$(cc_scalar "$cc_vr_dir/execution.yaml" "worker_failures")
	cc_vr_wf=$((cc_vr_wf + 1))
	cc_exec_set "$cc_vr_dir" worker_failures "$cc_vr_wf" || return 1
	if [ "$cc_vr_wf" -ge 3 ]; then
		cc_exec_set "$cc_vr_dir" status failed || return 1
		cc_emit outcome "$cc_vr_out"
		cc_emit worker_failures "$cc_vr_wf"
		cc_emit status failed
		cc_emit stop FAILURE_LIMIT_REACHED
		return 0
	fi
	cc_exec_set "$cc_vr_dir" status repairing || return 1
	cc_emit outcome "$cc_vr_out"
	cc_emit worker_failures "$cc_vr_wf"
	cc_emit status repairing
	return 0
}

# cc_attempt_evidence_record EXEC_DIR ATTEMPT KEY=VALUE ... -> record bounded,
# provider-neutral host evidence for one attempt (tiering): the
# (model, effort) each role ran at — the one thing only the coordinator knows,
# since the engine is model-blind (INV-RUNTIME-01). Inference timing is NOT
# recorded here: the attempt already carries started_at (worker.yaml) and
# checked_at (verifier.yaml), whose span is the attempt duration, so no fragile
# coordinator stopwatch is needed. This is EVIDENCE ONLY (INV-HOST-01): the engine
# stores it,
# never selects a model, never interprets it, and no gate, route, lease, verdict,
# or the worker-failure counter (INV-REPAIR-01) ever reads it. The key set is a
# fixed allowlist and values are credential-free (INV-SEC-01); an unknown key or
# unsafe value is refused, keeping the surface bounded. Written atomically
# (INV-RUNTIME-02).
cc_attempt_evidence_record() {
	cc_ae_dir="$1"; cc_ae_att="$2"
	[ -n "$cc_ae_dir" ] && [ -n "$cc_ae_att" ] || { cc_fail ATTEMPT_EVIDENCE_ARGS; return 1; }
	shift 2
	[ -f "$cc_ae_dir/execution.yaml" ] || { cc_fail EXECUTION_RECORD_MISSING; return 1; }
	cc_ae_att=$(cc_attempt_norm "$cc_ae_att") || return 1
	cc_ae_pad=$(printf '%03d' "$cc_ae_att")
	[ -d "$cc_ae_dir/attempts/$cc_ae_pad" ] || { cc_fail ATTEMPT_UNKNOWN "$cc_ae_att"; return 1; }
	[ "$#" -gt 0 ] || { cc_fail ATTEMPT_EVIDENCE_EMPTY; return 1; }
	cc_ae_f="$cc_ae_dir/attempts/$cc_ae_pad/host-evidence.yaml"
	if [ ! -f "$cc_ae_f" ]; then
		printf 'attempt: %s\n' "$cc_ae_att" | cc_atomic_write "$cc_ae_f" \
			|| { cc_fail ATTEMPT_EVIDENCE_WRITE_FAILED "$cc_ae_att"; return 1; }
	fi
	for cc_ae_kv in "$@"; do
		case "$cc_ae_kv" in *=*) : ;; *) cc_fail ATTEMPT_EVIDENCE_MALFORMED "$cc_ae_kv"; return 1 ;; esac
		cc_ae_k=${cc_ae_kv%%=*}; cc_ae_v=${cc_ae_kv#*=}
		case "$cc_ae_k" in
			worker_model|worker_effort|verifier_model|verifier_effort) : ;;
			*) cc_fail ATTEMPT_EVIDENCE_KEY_UNKNOWN "$cc_ae_k"; return 1 ;;
		esac
		case "$cc_ae_v" in
			''|*[!A-Za-z0-9._-]*) cc_fail ATTEMPT_EVIDENCE_VALUE_INVALID "$cc_ae_k"; return 1 ;;
		esac
		awk -v k="$cc_ae_k" -v v="$cc_ae_v" '
			$0 ~ ("^" k ":[[:space:]]") { next }
			{ print }
			END { print k ": " v }
		' "$cc_ae_f" | cc_atomic_write "$cc_ae_f" \
			|| { cc_fail ATTEMPT_EVIDENCE_WRITE_FAILED "$cc_ae_k"; return 1; }
		cc_emit "$cc_ae_k" "$cc_ae_v"
	done
	return 0
}

# cc_repair_allowed EXEC_DIR -> yes/no with remaining attempts
cc_repair_allowed() {
	cc_ra_wf=$(cc_scalar "$1/execution.yaml" "worker_failures")
	if [ "$cc_ra_wf" -lt 3 ]; then
		cc_emit repair allowed
		cc_emit remaining "$((3 - cc_ra_wf))"
		return 0
	fi
	cc_emit repair denied
	cc_fail FAILURE_LIMIT_REACHED; return 1
}

# cc_execution_status EXEC_DIR -> current execution state
cc_execution_status() { cc_scalar "$1/execution.yaml" "status"; }

# ---------------------------------------------------------------------------
# Candidate identity (Context Circuit v1.0, Mechanism 2, INV-CANDIDATE-01).
# A candidate is a deterministic digest over the exact proposed result: the
# per-repository commit map, the selected base commits, and the intent's frozen
# contract_digest. It is an identity computed over records the engine already
# keeps plus the contract digest — not a new state machine. Any new commit OR a
# criteria change yields a new candidate id, which is what voids prior evidence
# and human acceptance. The engine already detects a commit change (verifier
# read-only tip check); the candidate generalizes that to "candidate unchanged
# since the evidence was recorded", and extends it to human acceptance.
# ---------------------------------------------------------------------------

# cc_candidate_id EXEC_DIR -> deterministic cand-<hash> over the current commit
# map + bases + the execution's recorded contract_digest. Recomputed from the live
# repository records, so a new worker commit (which updates latest_commit) changes
# the id automatically.
cc_candidate_id() {
	cc_cid_dir="$1"
	[ -f "$cc_cid_dir/execution.yaml" ] || { cc_fail EXECUTION_RECORD_MISSING; return 1; }
	# Resolve the contract digest the candidate folds in. An active intent must still
	# have the exact frozen digest it was approved with; a changed-but-unapproved
	# contract is stale and cannot retain the old candidate. Archived intents are no
	# longer editable through the active path, so the execution snapshot is the
	# portable fallback for those records.
	cc_cid_cdig=""
	cc_cid_intent=$(cc_scalar "$cc_cid_dir/execution.yaml" "intent" 2>/dev/null) || cc_cid_intent=""
	if [ -n "$cc_cid_intent" ]; then
		cc_cid_root=$(cd -- "$cc_cid_dir/../../../.." 2>/dev/null && pwd) || cc_cid_root=""
		if [ -n "$cc_cid_root" ] && [ -f "$cc_cid_root/intent/$cc_cid_intent/contract.yaml" ]; then
			cc_cid_contract="$cc_cid_root/intent/$cc_cid_intent/contract.yaml"
			cc_cid_cdig=$(cc_scalar "$cc_cid_contract" "contract_digest" 2>/dev/null) || cc_cid_cdig=""
			cc_cid_now=$(cc_intent_contract_digest "$cc_cid_contract" 2>/dev/null) || cc_cid_now=""
			[ -n "$cc_cid_cdig" ] && [ "$cc_cid_cdig" = "$cc_cid_now" ] \
				|| { cc_fail CANDIDATE_CONTRACT_STALE; return 1; }
		fi
	fi
	[ -n "$cc_cid_cdig" ] || cc_cid_cdig=$(cc_scalar "$cc_cid_dir/execution.yaml" "contract_digest" 2>/dev/null) || cc_cid_cdig=""
	[ -n "$cc_cid_cdig" ] || { cc_fail CANDIDATE_NO_CONTRACT_DIGEST; return 1; }
	cc_cid_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-candidate.XXXXXX") || return 1
	: >"$cc_cid_tmp"
	cc_cid_found=no
	for cc_cid_rf in "$cc_cid_dir"/repositories/*.yaml; do
		[ -f "$cc_cid_rf" ] || continue
		cc_cid_found=yes
		cc_cid_wt=$(cc_scalar "$cc_cid_rf" worktree 2>/dev/null) || cc_cid_wt=""
		cc_cid_latest=$(cc_scalar "$cc_cid_rf" latest_commit 2>/dev/null) || cc_cid_latest=""
		[ -n "$cc_cid_wt" ] && [ -d "$cc_cid_wt" ] \
			|| { rm -f "$cc_cid_tmp"; cc_fail CANDIDATE_WORKTREE_MISSING; return 1; }
		cc_cid_head=$(git -C "$cc_cid_wt" rev-parse HEAD 2>/dev/null) || { rm -f "$cc_cid_tmp"; cc_fail CANDIDATE_WORKTREE_UNREADABLE; return 1; }
		[ "$cc_cid_head" = "$cc_cid_latest" ] || { rm -f "$cc_cid_tmp"; cc_fail CANDIDATE_COMMIT_DRIFT; return 1; }
		printf '%s\t%s\t%s\n' \
			"$(cc_scalar "$cc_cid_rf" repository)" \
			"$(cc_scalar "$cc_cid_rf" latest_commit)" \
			"$(cc_scalar "$cc_cid_rf" base_commit)" >>"$cc_cid_tmp"
	done
	[ "$cc_cid_found" = yes ] || { rm -f "$cc_cid_tmp"; cc_fail CANDIDATE_REPOSITORIES_MISSING; return 1; }
	printf 'contract\t%s\n' "$cc_cid_cdig" >>"$cc_cid_tmp"
	cc_cid_canon=$(LC_ALL=C sort "$cc_cid_tmp")
	rm -f "$cc_cid_tmp"
	cc_cid_hash=$(cc_digest_text "$cc_cid_canon")
	cc_cid_hash=${cc_cid_hash#sha256:}; cc_cid_hash=${cc_cid_hash#cksum:}
	printf 'cand-%s' "$cc_cid_hash"
}

# cc_candidate_digest ROOT PLAN EXEC -> compute the current candidate and record
# candidate.yaml under the execution dir; prints candidate_id. Idempotent for an
# unchanged result; writes a new identity after any commit or criteria change.
cc_candidate_digest() {
	cc_cd_dir=$(cc_execution_dir "$1" "$2" "$3")
	[ -f "$cc_cd_dir/execution.yaml" ] || { cc_fail CANDIDATE_NO_EXECUTION "$2/$3"; return 1; }
	cc_cd_id=$(cc_candidate_id "$cc_cd_dir") || return 1
	cc_cd_cdig=$(cc_scalar "$cc_cd_dir/execution.yaml" "contract_digest" 2>/dev/null) || cc_cd_cdig=""
	cc_cd_intent=$(cc_scalar "$cc_cd_dir/execution.yaml" "intent" 2>/dev/null) || cc_cd_intent=""
	cc_cd_root=$(cd -- "$cc_cd_dir/../../../.." 2>/dev/null && pwd) || cc_cd_root=""
	if [ -n "$cc_cd_root" ] && [ -n "$cc_cd_intent" ] && [ -f "$cc_cd_root/intent/$cc_cd_intent/contract.yaml" ]; then
		cc_cd_cdig=$(cc_scalar "$cc_cd_root/intent/$cc_cd_intent/contract.yaml" "contract_digest" 2>/dev/null) || cc_cd_cdig=""
	fi
	[ -n "$cc_cd_cdig" ] || { cc_fail CANDIDATE_NO_CONTRACT_DIGEST; return 1; }
	{
		printf 'schema_version: %s\ncandidate_id: %s\nplan: %s\nexecution_id: %s\ncontract_digest: %s\nrepositories:\n' \
			"$CC_CANDIDATE_SCHEMA_VERSION" "$cc_cd_id" "$2" "$3" "$cc_cd_cdig"
		for cc_cd_rf in "$cc_cd_dir"/repositories/*.yaml; do
			[ -f "$cc_cd_rf" ] || continue
			printf '  %s: %s\n' "$(cc_scalar "$cc_cd_rf" repository)" "$(cc_scalar "$cc_cd_rf" latest_commit)"
		done
		printf 'bases:\n'
		for cc_cd_rf in "$cc_cd_dir"/repositories/*.yaml; do
			[ -f "$cc_cd_rf" ] || continue
			printf '  %s: %s\n' "$(cc_scalar "$cc_cd_rf" repository)" "$(cc_scalar "$cc_cd_rf" base_commit)"
		done
		printf 'computed_at: %s\n' "$(cc_now)"
	} | cc_atomic_write "$cc_cd_dir/candidate.yaml"
	cc_emit candidate_id "$cc_cd_id"
	cc_emit plan "$2"
	cc_emit execution_id "$3"
	return 0
}

# cc_candidate_current ROOT PLAN -> report the current candidate id for the plan's
# latest execution (recomputed live; not necessarily the last recorded one).
cc_candidate_current() {
	cc_cc_exec=$(cc_latest_execution "$1" "$2") || { cc_fail CANDIDATE_NO_EXECUTION "$2"; return 1; }
	[ -n "$cc_cc_exec" ] || { cc_fail CANDIDATE_NO_EXECUTION "$2"; return 1; }
	cc_cc_dir=$(cc_execution_dir "$1" "$2" "$cc_cc_exec")
	cc_cc_id=$(cc_candidate_id "$cc_cc_dir") || return 1
	cc_emit candidate_id "$cc_cc_id"
	cc_emit execution_id "$cc_cc_exec"
	return 0
}

# cc_plan_single_repository ROOT PLAN -> the plan's one repository id. A plan
# that names zero or several repositories cannot be a covering-tip member
# (INV-DELIVER-01 / one-repo-per-plan).
cc_plan_single_repository() {
	cc_psr_root="$1"; cc_psr_plan="$2"
	cc_psr_repos=$(cc_plan_affected_repositories "$cc_psr_root/plans/$cc_psr_plan/plan.yaml")
	cc_psr_n=0
	cc_psr_one=""
	for cc_psr_r in $cc_psr_repos; do
		cc_psr_n=$((cc_psr_n + 1))
		cc_psr_one="$cc_psr_r"
	done
	[ "$cc_psr_n" -eq 1 ] || { cc_fail CHANGE_SET_CROSS_REPO "$cc_psr_plan"; return 1; }
	printf '%s' "$cc_psr_one"
}

# cc_change_set_require_same_repo ROOT PLAN... -> ok when the members are a
# change set of one, or when every member names exactly the same single
# repository. Two or more members that do not share one repository cannot
# converge to one pull request (INV-DELIVER-01); partition them first.
cc_change_set_require_same_repo() {
	cc_cssr_root="$1"; shift
	[ "$#" -le 1 ] && return 0
	cc_cssr_expected=""
	for cc_cssr_m in "$@"; do
		cc_cssr_one=$(cc_plan_single_repository "$cc_cssr_root" "$cc_cssr_m") || return 1
		if [ -z "$cc_cssr_expected" ]; then
			cc_cssr_expected="$cc_cssr_one"
		elif [ "$cc_cssr_one" != "$cc_cssr_expected" ]; then
			cc_fail CHANGE_SET_CROSS_REPO "$cc_cssr_m:$cc_cssr_one"
			return 1
		fi
	done
	return 0
}

# cc_change_set_covering_plan ROOT REPO PLAN... -> the one member whose execution
# branch contains every other member (the shippable tip for one pull request).
# Parallel tips that do not contain each other fail CHANGE_SET_NO_SINGLE_TIP;
# delivery never authors a merge to invent one (INV-DELIVER-01).
cc_change_set_covering_plan() {
	cc_cscv_root="$1"; cc_cscv_repo="$2"; shift 2
	[ "$#" -ge 1 ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_cscv_abs=$(cc_repo_resolve "$cc_cscv_root" "$cc_cscv_repo" | sed -n 's/^path: //p') \
		|| { cc_fail CHANGE_SET_REPO_UNRESOLVED "$cc_cscv_repo"; return 1; }
	cc_cscv_cover=""
	cc_cscv_cover_tip=""
	for cc_cscv_m in "$@"; do
		cc_cscv_br="refs/heads/cc/$cc_cscv_m/$cc_cscv_repo"
		git -C "$cc_cscv_abs" show-ref --verify --quiet "$cc_cscv_br" \
			|| { cc_fail CHANGE_SET_BRANCH_MISSING "$cc_cscv_m:$cc_cscv_repo"; return 1; }
		cc_cscv_ok=yes
		for cc_cscv_o in "$@"; do
			[ "$cc_cscv_o" = "$cc_cscv_m" ] && continue
			git -C "$cc_cscv_abs" merge-base --is-ancestor \
				"refs/heads/cc/$cc_cscv_o/$cc_cscv_repo" "$cc_cscv_br" \
				|| cc_cscv_ok=no
		done
		[ "$cc_cscv_ok" = yes ] || continue
		cc_cscv_tip=$(git -C "$cc_cscv_abs" rev-parse "$cc_cscv_br")
		if [ -z "$cc_cscv_cover" ]; then
			cc_cscv_cover="$cc_cscv_m"
			cc_cscv_cover_tip="$cc_cscv_tip"
		elif [ "$cc_cscv_tip" != "$cc_cscv_cover_tip" ]; then
			cc_fail CHANGE_SET_NO_SINGLE_TIP "$cc_cscv_cover,$cc_cscv_m"
			return 1
		fi
	done
	[ -n "$cc_cscv_cover" ] || { cc_fail CHANGE_SET_NO_SINGLE_TIP; return 1; }
	printf '%s' "$cc_cscv_cover"
}

# cc_change_set_covering_tips ROOT REPO PLAN... -> maximal shippable heads among
# the named members (space-separated plan ids, sorted). A member is a covering
# tip when no other named member strictly contains it. Sibling stacks in one
# repository are several tips and several pull requests; CHANGE_SET_NO_SINGLE_TIP
# applies only when those tips are forced into one change-set-prepare.
cc_change_set_covering_tips() {
	cc_cstp_root="$1"; cc_cstp_repo="$2"; shift 2
	[ "$#" -ge 1 ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_cstp_abs=$(cc_repo_resolve "$cc_cstp_root" "$cc_cstp_repo" | sed -n 's/^path: //p') \
		|| { cc_fail CHANGE_SET_REPO_UNRESOLVED "$cc_cstp_repo"; return 1; }
	cc_cstp_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-covers.XXXXXX") || return 1
	: >"$cc_cstp_tmp"
	for cc_cstp_m in "$@"; do
		cc_cstp_br="refs/heads/cc/$cc_cstp_m/$cc_cstp_repo"
		git -C "$cc_cstp_abs" show-ref --verify --quiet "$cc_cstp_br" \
			|| { rm -f "$cc_cstp_tmp"; cc_fail CHANGE_SET_BRANCH_MISSING "$cc_cstp_m:$cc_cstp_repo"; return 1; }
		cc_cstp_mtip=$(git -C "$cc_cstp_abs" rev-parse "$cc_cstp_br")
		cc_cstp_dominated=no
		for cc_cstp_o in "$@"; do
			[ "$cc_cstp_o" = "$cc_cstp_m" ] && continue
			git -C "$cc_cstp_abs" merge-base --is-ancestor \
				"$cc_cstp_br" "refs/heads/cc/$cc_cstp_o/$cc_cstp_repo" 2>/dev/null || continue
			cc_cstp_otip=$(git -C "$cc_cstp_abs" rev-parse "refs/heads/cc/$cc_cstp_o/$cc_cstp_repo")
			[ "$cc_cstp_mtip" != "$cc_cstp_otip" ] && cc_cstp_dominated=yes
		done
		[ "$cc_cstp_dominated" = yes ] || printf '%s\t%s\n' "$cc_cstp_mtip" "$cc_cstp_m" >>"$cc_cstp_tmp"
	done
	cc_cstp_out=""
	for cc_cstp_sha in $(cut -f1 "$cc_cstp_tmp" | LC_ALL=C sort -u); do
		cc_cstp_cover=$(awk -F '\t' -v s="$cc_cstp_sha" '$1 == s { print $2 }' "$cc_cstp_tmp" | LC_ALL=C sort | sed -n '1p')
		cc_cstp_out="${cc_cstp_out:+$cc_cstp_out }$cc_cstp_cover"
	done
	rm -f "$cc_cstp_tmp"
	[ -n "$cc_cstp_out" ] || { cc_fail CHANGE_SET_NO_SINGLE_TIP; return 1; }
	printf '%s' "$cc_cstp_out"
}

# cc_change_set_partition ROOT PLAN... -> one group per covering tip. Mixed
# repositories and sibling stacks in one repository are both several pull
# requests. A linear stack is one group. Forcing sibling tips into one
# change-set-prepare still fails CHANGE_SET_NO_SINGLE_TIP.
cc_change_set_partition() {
	cc_cspp_root=$(cc_root_abs "$1" 2>/dev/null) || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$1"; return 1; }
	shift
	[ "$#" -ge 1 ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_cspp_members=$(printf '%s\n' "$@" | LC_ALL=C sort -u | sed '/^$/d')
	[ -n "$cc_cspp_members" ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_cspp_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-partition.XXXXXX") || return 1
	: >"$cc_cspp_tmp"
	for cc_cspp_m in $cc_cspp_members; do
		cc_plan_validate "$cc_cspp_root/plans/$cc_cspp_m" >/dev/null \
			|| { rm -f "$cc_cspp_tmp"; cc_fail CHANGE_SET_PLAN_INVALID "$cc_cspp_m"; return 1; }
		cc_cspp_repo=$(cc_plan_single_repository "$cc_cspp_root" "$cc_cspp_m") \
			|| { rm -f "$cc_cspp_tmp"; return 1; }
		printf '%s\t%s\n' "$cc_cspp_repo" "$cc_cspp_m" >>"$cc_cspp_tmp"
	done
	cc_cspp_n=0
	cc_cspp_ready=0
	for cc_cspp_repo in $(cut -f1 "$cc_cspp_tmp" | LC_ALL=C sort -u); do
		cc_cspp_group=$(awk -F '\t' -v r="$cc_cspp_repo" '$1 == r { print $2 }' "$cc_cspp_tmp" | LC_ALL=C sort)
		cc_cspp_abs=$(cc_repo_resolve "$cc_cspp_root" "$cc_cspp_repo" | sed -n 's/^path: //p') \
			|| { rm -f "$cc_cspp_tmp"; cc_fail CHANGE_SET_REPO_UNRESOLVED "$cc_cspp_repo"; return 1; }
		cc_cspp_covers=$(cc_change_set_covering_tips "$cc_cspp_root" "$cc_cspp_repo" $cc_cspp_group) \
			|| { rm -f "$cc_cspp_tmp"; return 1; }
		for cc_cspp_cover in $(printf '%s' "$cc_cspp_covers" | tr ' ' '\n' | LC_ALL=C sort); do
			cc_cspp_subset=""
			for cc_cspp_m in $cc_cspp_group; do
				git -C "$cc_cspp_abs" merge-base --is-ancestor \
					"refs/heads/cc/$cc_cspp_m/$cc_cspp_repo" \
					"refs/heads/cc/$cc_cspp_cover/$cc_cspp_repo" 2>/dev/null \
					|| continue
				cc_cspp_subset="${cc_cspp_subset:+$cc_cspp_subset }$cc_cspp_m"
			done
			cc_cspp_subset=$(printf '%s' "$cc_cspp_subset" | tr ' ' '\n' | LC_ALL=C sort | tr '\n' ' ' | sed 's/ *$//')
			cc_cspp_memblock=$(printf '%s' "$cc_cspp_subset" | sed 's/ /, /g')
			cc_cspp_gn=0
			for cc_cspp_gm in $cc_cspp_subset; do cc_cspp_gn=$((cc_cspp_gn + 1)); done
			if [ "$cc_cspp_gn" -eq 1 ]; then
				cc_cspp_one=$(printf '%s\n' "$cc_cspp_subset" | tr ' ' '\n' | sed -n '1p')
				cc_cspp_exec=$(cc_latest_execution "$cc_cspp_root" "$cc_cspp_one") \
					|| { rm -f "$cc_cspp_tmp"; return 1; }
				cc_cspp_cand=$(cc_candidate_id "$(cc_execution_dir "$cc_cspp_root" "$cc_cspp_one" "$cc_cspp_exec")") \
					|| { rm -f "$cc_cspp_tmp"; return 1; }
			else
				cc_cspp_cand=$(cc_change_set_tips_digest "$cc_cspp_root" $cc_cspp_subset) \
					|| { rm -f "$cc_cspp_tmp"; return 1; }
			fi
			cc_cspp_n=$((cc_cspp_n + 1))
			cc_cspp_ready=$((cc_cspp_ready + 1))
			cc_emit group "$cc_cspp_n"
			cc_emit repository "$cc_cspp_repo"
			cc_emit members "$cc_cspp_memblock"
			cc_emit covering_plan "$cc_cspp_cover"
			cc_emit source_branch "cc/$cc_cspp_cover/$cc_cspp_repo"
			cc_emit candidate_id "$cc_cspp_cand"
			cc_emit status ready
		done
	done
	rm -f "$cc_cspp_tmp"
	cc_emit groups "$cc_cspp_n"
	cc_emit pull_requests "$cc_cspp_ready"
	return 0
}

# cc_change_set_tips_digest ROOT PLAN... -> cand-<hash> over the live member tip
# map (commits, bases, per-plan candidates). This IS the change-set candidate;
# delivery does not replace it with a newly authored merge.
cc_change_set_tips_digest() {
	cc_cstd_root="$1"; shift
	[ "$#" -ge 1 ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_cstd_tmp=$(mktemp "${TMPDIR:-/tmp}/cc-candidate.XXXXXX") || return 1
	: >"$cc_cstd_tmp"
	for cc_cstd_plan in "$@"; do
		cc_cstd_pdir="$cc_cstd_root/plans/$cc_cstd_plan"
		cc_plan_validate "$cc_cstd_pdir" >/dev/null \
			|| { rm -f "$cc_cstd_tmp"; cc_fail CHANGE_SET_PLAN_INVALID "$cc_cstd_plan"; return 1; }
		cc_cstd_exec=$(cc_latest_execution "$cc_cstd_root" "$cc_cstd_plan") \
			|| { rm -f "$cc_cstd_tmp"; cc_fail CHANGE_SET_NO_EXECUTION "$cc_cstd_plan"; return 1; }
		[ -n "$cc_cstd_exec" ] || { rm -f "$cc_cstd_tmp"; cc_fail CHANGE_SET_NO_EXECUTION "$cc_cstd_plan"; return 1; }
		cc_cstd_edir=$(cc_execution_dir "$cc_cstd_root" "$cc_cstd_plan" "$cc_cstd_exec")
		[ -f "$cc_cstd_edir/execution.yaml" ] || { rm -f "$cc_cstd_tmp"; cc_fail CHANGE_SET_NO_EXECUTION "$cc_cstd_plan"; return 1; }
		cc_cstd_plan_cand=$(cc_candidate_id "$cc_cstd_edir") \
			|| { rm -f "$cc_cstd_tmp"; cc_fail CHANGE_SET_MEMBER_CANDIDATE_INVALID "$cc_cstd_plan"; return 1; }
		printf 'plan\t%s\t%s\n' "$cc_cstd_plan" "$cc_cstd_plan_cand" >>"$cc_cstd_tmp"
		cc_cstd_found=no
		for cc_cstd_rf in "$cc_cstd_edir"/repositories/*.yaml; do
			[ -f "$cc_cstd_rf" ] || continue
			cc_cstd_found=yes
			printf '%s\t%s\t%s\t%s\n' \
				"$(cc_scalar "$cc_cstd_rf" repository)" \
				"$(cc_scalar "$cc_cstd_rf" latest_commit)" \
				"$(cc_scalar "$cc_cstd_rf" base_commit)" \
				"$cc_cstd_plan_cand" >>"$cc_cstd_tmp"
		done
		[ "$cc_cstd_found" = yes ] || { rm -f "$cc_cstd_tmp"; cc_fail CHANGE_SET_REPOSITORIES_MISSING "$cc_cstd_plan"; return 1; }
	done
	cc_cstd_canon=$(LC_ALL=C sort -u "$cc_cstd_tmp")
	rm -f "$cc_cstd_tmp"
	[ -n "$cc_cstd_canon" ] || { cc_fail CHANGE_SET_NO_EXECUTIONS; return 1; }
	cc_cstd_hash=$(cc_digest_text "$cc_cstd_canon")
	cc_cstd_hash=${cc_cstd_hash#sha256:}; cc_cstd_hash=${cc_cstd_hash#cksum:}
	printf 'cand-%s' "$cc_cstd_hash"
}

# cc_change_set_candidate ROOT PLAN... -> the identity of a change set: the set of
# same-repository plans delivered together as one pull request. A change set of
# one is exactly that plan's candidate (delegates to cc_candidate_id, so it
# equals candidate-current). For several same-repo plans it digests the live
# member tip map. Members that do not share one repository fail
# CHANGE_SET_CROSS_REPO. Delivery never authors an integration merge.
cc_change_set_candidate() {
	cc_cs_root="$1"; shift
	[ "$#" -ge 1 ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_cs_id=$(cc_change_set_id "$cc_cs_root" "$@")
	cc_cs_existing=$(cc_change_set_dir "$cc_cs_root" "$cc_cs_id")
	if [ "$#" -gt 1 ] && [ -f "$cc_cs_existing/change-set.yaml" ]; then
		cc_change_set_candidate_from_record "$cc_cs_root" "$cc_cs_id"
		return $?
	fi
	# a change set of one is exactly that plan's candidate; delegate so the two
	# never diverge (candidate-current == change-set-candidate for one plan).
	if [ "$#" -eq 1 ]; then
		cc_cs_exec=$(cc_latest_execution "$cc_cs_root" "$1") || { cc_fail CHANGE_SET_NO_EXECUTIONS; return 1; }
		[ -n "$cc_cs_exec" ] || { cc_fail CHANGE_SET_NO_EXECUTIONS; return 1; }
		cc_cs_one=$(cc_candidate_id "$(cc_execution_dir "$cc_cs_root" "$1" "$cc_cs_exec")") || return 1
		cc_emit change_set_candidate "$cc_cs_one"
		return 0
	fi
	for cc_cs_plan in "$@"; do
		cc_plan_validate "$cc_cs_root/plans/$cc_cs_plan" >/dev/null \
			|| { cc_fail CHANGE_SET_PLAN_INVALID "$cc_cs_plan"; return 1; }
	done
	cc_change_set_require_same_repo "$cc_cs_root" "$@" || return 1
	cc_cs_hash=$(cc_change_set_tips_digest "$cc_cs_root" "$@") || return 1
	cc_emit change_set_candidate "$cc_cs_hash"
	return 0
}

# ---------------------------------------------------------------------------
# Change-set delivery (INV-DELIVER-01): same-repository plans shipped as one
# pull request. The candidate is the live member tip map. Delivery records that
# set, names the one covering execution branch, and reuses each member's
# already-bound independent pass. It never authors a merge and never spawns a
# verifier. Plans in different repositories never form a change set.
# ---------------------------------------------------------------------------

cc_change_set_dir() { printf '%s/.runtime/change-sets/%s' "$1" "$2"; }

# cc_change_set_record_status DIR STATUS REASON -> update a preserved preparation
# record after setup fails. Failed/interrupted integration is state, not disposable
# scratch, so callers never remove the change-set directory to hide the failure.
cc_change_set_record_status() {
	cc_csrst_dir="$1"; cc_csrst_status="$2"; cc_csrst_reason="${3:-}"
	[ -f "$cc_csrst_dir/change-set.yaml" ] || return 1
	awk -v s="$cc_csrst_status" -v r="$cc_csrst_reason" '
		$0 ~ /^status:[[:space:]]/ { print "status: " s; seen=1; next }
		$0 ~ /^blocked_reason:[[:space:]]/ { print "blocked_reason: " r; brec=1; next }
		{ print }
		END { if (!seen) print "status: " s; if (r != "" && !brec) print "blocked_reason: " r }
	' "$cc_csrst_dir/change-set.yaml" | cc_atomic_write "$cc_csrst_dir/change-set.yaml"
}

# cc_change_set_candidate_from_record ROOT CS_ID -> recompute the tip-map
# candidate from the recorded members. A member that moved yields a new id.
cc_change_set_candidate_from_record() {
	cc_cscr_root="$1"; cc_cscr_id="$2"; cc_cscr_dir=$(cc_change_set_dir "$cc_cscr_root" "$cc_cscr_id")
	[ -f "$cc_cscr_dir/change-set.yaml" ] || { cc_fail CHANGE_SET_UNKNOWN "$cc_cscr_id"; return 1; }
	cc_cscr_members=$(cc_inline_list "$(cc_scalar "$cc_cscr_dir/change-set.yaml" members)")
	[ -n "$cc_cscr_members" ] || { cc_fail CHANGE_SET_MEMBERS_MISSING "$cc_cscr_id"; return 1; }
	set -- $cc_cscr_members
	if [ "$#" -eq 1 ]; then
		cc_cscr_exec=$(cc_latest_execution "$cc_cscr_root" "$1") \
			|| { cc_fail CHANGE_SET_NO_EXECUTION "$1"; return 1; }
		[ -n "$cc_cscr_exec" ] || { cc_fail CHANGE_SET_NO_EXECUTION "$1"; return 1; }
		cc_cscr_one=$(cc_candidate_id "$(cc_execution_dir "$cc_cscr_root" "$1" "$cc_cscr_exec")") || return 1
		cc_emit change_set_candidate "$cc_cscr_one"
		return 0
	fi
	cc_cscr_hash=$(cc_change_set_tips_digest "$cc_cscr_root" "$@") || return 1
	cc_emit change_set_candidate "$cc_cscr_hash"
}

# cc_change_set_id ROOT PLAN... -> deterministic id over the sorted member set.
cc_change_set_id() {
	shift
	cc_csi_members=$(printf '%s\n' "$@" | LC_ALL=C sort -u | sed '/^$/d')
	cc_csi_h=$(cc_digest_text "$cc_csi_members"); cc_csi_h=${cc_csi_h#sha256:}; cc_csi_h=${cc_csi_h#cksum:}
	printf 'cs-%s' "$cc_csi_h"
}

# cc_change_set_tier ROOT PLAN... -> the max consequence tier across members
# (fail upward: critical > standard > explore).
cc_change_set_tier() {
	cc_cst_root="$1"; shift
	cc_cst_tier=explore
	for cc_cst_m in "$@"; do
		cc_cst_e=$(cc_latest_execution "$cc_cst_root" "$cc_cst_m" 2>/dev/null) || continue
		[ -n "$cc_cst_e" ] || continue
		cc_cst_t=$(cc_scalar "$(cc_execution_dir "$cc_cst_root" "$cc_cst_m" "$cc_cst_e")/execution.yaml" tier 2>/dev/null) || cc_cst_t=standard
		case "$cc_cst_t" in
			critical) cc_cst_tier=critical ;;
			standard) [ "$cc_cst_tier" = critical ] || cc_cst_tier=standard ;;
		esac
	done
	printf '%s' "$cc_cst_tier"
}

# cc_change_set_prepare ROOT PLAN... -> record the member tip-map candidate and
# the one covering execution branch. Delivery never authors a merge and never
# spawns a verifier. Members that do not share one repository fail
# CHANGE_SET_CROSS_REPO — partition first (`change-set-partition`) and prepare
# one group per covering tip. Tips that do not converge to one shippable branch
# fail CHANGE_SET_NO_SINGLE_TIP.
cc_change_set_prepare() {
	cc_csp_root=$(cc_root_abs "$1" 2>/dev/null) || { cc_fail WORKSPACE_ROOT_NOT_FOUND "$1"; return 1; }
	shift
	[ "$#" -ge 1 ] || { cc_fail CHANGE_SET_EMPTY; return 1; }
	cc_csp_members=$(printf '%s\n' "$@" | LC_ALL=C sort -u | sed '/^$/d')
	cc_csp_id=$(cc_change_set_id "$cc_csp_root" $cc_csp_members)
	cc_csp_dir=$(cc_change_set_dir "$cc_csp_root" "$cc_csp_id")
	[ ! -e "$cc_csp_dir/change-set.yaml" ] || { cc_fail CHANGE_SET_EXISTS "$cc_csp_id"; return 1; }
	for cc_csp_m in $cc_csp_members; do
		cc_plan_validate "$cc_csp_root/plans/$cc_csp_m" >/dev/null \
			|| { cc_fail CHANGE_SET_MEMBER_INVALID "$cc_csp_m"; return 1; }
		cc_intent_authorized "$cc_csp_root" "$cc_csp_m" >/dev/null \
			|| { cc_fail CHANGE_SET_MEMBER_UNAUTHORIZED "$cc_csp_m"; return 1; }
		cc_csp_exec=$(cc_latest_execution "$cc_csp_root" "$cc_csp_m") || { cc_fail CHANGE_SET_NO_EXECUTION "$cc_csp_m"; return 1; }
		[ -n "$cc_csp_exec" ] || { cc_fail CHANGE_SET_NO_EXECUTION "$cc_csp_m"; return 1; }
		cc_csp_edir=$(cc_execution_dir "$cc_csp_root" "$cc_csp_m" "$cc_csp_exec")
		[ -f "$cc_csp_edir/execution.yaml" ] || { cc_fail CHANGE_SET_NO_EXECUTION "$cc_csp_m"; return 1; }
	done
	cc_change_set_require_same_repo "$cc_csp_root" $cc_csp_members || return 1
	cc_csp_repos=$(for cc_csp_m in $cc_csp_members; do cc_plan_affected_repositories "$cc_csp_root/plans/$cc_csp_m/plan.yaml"; done | LC_ALL=C sort -u | sed '/^$/d')
	[ -n "$cc_csp_repos" ] || { cc_fail CHANGE_SET_NO_REPOSITORIES; return 1; }
	cc_csp_repo=$(printf '%s\n' "$cc_csp_repos" | sed -n '1p')
	cc_csp_cover=$(cc_change_set_covering_plan "$cc_csp_root" "$cc_csp_repo" $cc_csp_members) || return 1
	cc_csp_src="cc/$cc_csp_cover/$cc_csp_repo"
	cc_csp_memblock=$(printf '%s' "$cc_csp_members" | tr '\n' ' ' | sed 's/ *$//; s/ /, /g')
	cc_csp_tierv=$(cc_change_set_tier "$cc_csp_root" $cc_csp_members)
	cc_csp_n=0
	for cc_csp_m in $cc_csp_members; do cc_csp_n=$((cc_csp_n + 1)); done
	if [ "$cc_csp_n" -eq 1 ]; then
		cc_csp_one=$(printf '%s\n' "$cc_csp_members" | sed -n '1p')
		cc_csp_exec=$(cc_latest_execution "$cc_csp_root" "$cc_csp_one") || return 1
		cc_csp_cand=$(cc_candidate_id "$(cc_execution_dir "$cc_csp_root" "$cc_csp_one" "$cc_csp_exec")") || return 1
	else
		cc_csp_cand=$(cc_change_set_tips_digest "$cc_csp_root" $cc_csp_members) || return 1
	fi
	mkdir -p "$cc_csp_dir"
	{
		printf 'schema_version: %s\nchange_set: %s\ncandidate_id: %s\ntier: %s\nstatus: prepared\nmembers: [%s]\nsource_plan: %s\nsource_branch: %s\ncreated_at: %s\nrepositories:\n  - repository: %s\n    branch: %s\n' \
			"$CC_CHANGE_SET_SCHEMA_VERSION" "$cc_csp_id" "$cc_csp_cand" "$cc_csp_tierv" "$cc_csp_memblock" "$cc_csp_cover" "$cc_csp_src" "$(cc_now)" "$cc_csp_repo" "$cc_csp_src"
	} | cc_atomic_write "$cc_csp_dir/change-set.yaml" \
		|| { cc_fail CHANGE_SET_RECORD_WRITE_FAILED "$cc_csp_id"; return 1; }
	cc_emit change_set "$cc_csp_id"
	cc_emit candidate_id "$cc_csp_cand"
	cc_emit tier "$cc_csp_tierv"
	cc_emit status prepared
	cc_emit source_plan "$cc_csp_cover"
	cc_emit source_branch "$cc_csp_src"
	return 0
}

# Delivery never records a change-set verifier. Each member's existing
# candidate-bound pass is the floor (INV-DELIVER-01).
cc_change_set_verifier_prepare() {
	cc_fail CHANGE_SET_DELIVERY_HAS_NO_VERIFIER "${2:-}"
	return 1
}

cc_change_set_verifier_record() {
	cc_fail CHANGE_SET_DELIVERY_HAS_NO_VERIFIER "${2:-}"
	return 1
}

# cc_change_set_accept ROOT CS_ID ACCEPTED_BY -> one human acceptance bound to the
# change-set candidate.
cc_change_set_accept() {
	cc_csa_dir=$(cc_change_set_dir "$1" "$2")
	[ -f "$cc_csa_dir/change-set.yaml" ] || { cc_fail CHANGE_SET_UNKNOWN "$2"; return 1; }
	[ -n "${3:-}" ] || { cc_fail ACCEPTANCE_NO_HUMAN; return 1; }
	case "$3" in *[!A-Za-z0-9._@-]*) cc_fail ACCEPTANCE_HUMAN_INVALID "$3"; return 1 ;; esac
	cc_csa_cand=$(cc_scalar "$cc_csa_dir/change-set.yaml" candidate_id)
	printf 'schema_version: %s\ncandidate_id: %s\naccepted_by: %s\naccepted_at: %s\n' \
		"$CC_HUMAN_ACCEPTANCE_SCHEMA_VERSION" "$cc_csa_cand" "$3" "$(cc_now)" | cc_atomic_write "$cc_csa_dir/human-acceptance.yaml" \
		|| { cc_fail ACCEPTANCE_RECORD_WRITE_FAILED; return 1; }
	cc_emit change_set "$2"
	cc_emit accepted_by "$3"
	cc_emit candidate_id "$cc_csa_cand"
	return 0
}

# cc_change_set_ready ROOT CS_ID -> the floor for delivering a change set as one
# PR: the recorded tip-map candidate must still describe the live members, a
# human acceptance must bind to it, and each Standard/Critical member must still
# have its own candidate-bound independent pass. Delivery does not add a
# change-set verifier. Explore is planless and cannot form a change set.
cc_change_set_ready() {
	cc_csr_root="$1"; cc_csr_dir=$(cc_change_set_dir "$1" "$2")
	[ -f "$cc_csr_dir/change-set.yaml" ] || { cc_fail CHANGE_SET_UNKNOWN "$2"; return 1; }
	cc_csr_status=$(cc_scalar "$cc_csr_dir/change-set.yaml" status 2>/dev/null) || cc_csr_status=""
	[ "$cc_csr_status" = prepared ] || { cc_fail CHANGE_SET_NOT_PREPARED "$2"; return 1; }
	cc_csr_recorded=$(cc_scalar "$cc_csr_dir/change-set.yaml" candidate_id)
	cc_csr_mlist=$(cc_inline_list "$(cc_scalar "$cc_csr_dir/change-set.yaml" members)")
	cc_csr_live_out=$(cc_change_set_candidate_from_record "$cc_csr_root" "$2") || return 1
	cc_csr_live=$(printf '%s\n' "$cc_csr_live_out" | sed -n 's/^change_set_candidate: //p')
	if [ "$cc_csr_live" != "$cc_csr_recorded" ]; then
		cc_emit change_set_ready stale
		cc_fail CHANGE_SET_STALE "$2"; return 1
	fi
	cc_csr_tier=$(cc_scalar "$cc_csr_dir/change-set.yaml" tier)
	cc_csr_acc=$(cc_scalar "$cc_csr_dir/human-acceptance.yaml" candidate_id 2>/dev/null) || cc_csr_acc=""
	if [ "$cc_csr_acc" != "$cc_csr_recorded" ]; then
		cc_emit change_set_ready blocked
		cc_fail CHANGE_SET_NO_ACCEPTANCE "$2"; return 1
	fi
	if [ "$cc_csr_tier" = standard ] || [ "$cc_csr_tier" = critical ]; then
		for cc_csr_m in $cc_csr_mlist; do
			cc_csr_exec=$(cc_latest_execution "$cc_csr_root" "$cc_csr_m") \
				|| { cc_emit change_set_ready blocked; cc_fail CHANGE_SET_FLOOR_UNMET "$cc_csr_m"; return 1; }
			cc_csr_edir=$(cc_execution_dir "$cc_csr_root" "$cc_csr_m" "$cc_csr_exec")
			[ "$(cc_execution_status "$cc_csr_edir")" = verified ] \
				|| { cc_emit change_set_ready blocked; cc_fail CHANGE_SET_FLOOR_UNMET "$cc_csr_m"; return 1; }
			cc_csr_vc=$(cc_scalar "$cc_csr_edir/execution.yaml" verified_candidate 2>/dev/null) || cc_csr_vc=""
			cc_csr_now=$(cc_candidate_id "$cc_csr_edir") || { cc_emit change_set_ready blocked; cc_fail CHANGE_SET_FLOOR_UNMET "$cc_csr_m"; return 1; }
			[ -n "$cc_csr_vc" ] && [ "$cc_csr_vc" = "$cc_csr_now" ] \
				|| { cc_emit change_set_ready blocked; cc_fail CHANGE_SET_FLOOR_UNMET "$cc_csr_m"; return 1; }
		done
		cc_emit assurance independent
	else
		cc_fail CHANGE_SET_EXPLORE_PLANLESS "$2"; return 1
	fi
	cc_emit change_set_ready eligible
	cc_emit candidate_id "$cc_csr_recorded"
	cc_emit tier "$cc_csr_tier"
	return 0
}

# cc_change_set_complete ROOT CS_ID -> complete every member of a delivered change set
# from the ONE change-set acceptance (pain 6: accept once, complete the set).
# Assurance is each member's already-bound independent pass, not a delivery-time
# verifier.
cc_change_set_complete() {
	cc_csc_root="$1"; cc_csc_dir=$(cc_change_set_dir "$1" "$2")
	[ -f "$cc_csc_dir/change-set.yaml" ] || { cc_fail CHANGE_SET_UNKNOWN "$2"; return 1; }
	cc_change_set_ready "$cc_csc_root" "$2" >/dev/null || { cc_fail CHANGE_SET_NOT_READY "$2"; return 1; }
	cc_csc_members=$(cc_inline_list "$(cc_scalar "$cc_csc_dir/change-set.yaml" members)")
	cc_csc_candidate=$(cc_scalar "$cc_csc_dir/change-set.yaml" candidate_id)
	cc_csc_acceptor=$(cc_scalar "$cc_csc_dir/human-acceptance.yaml" accepted_by)
	cc_csc_tier=$(cc_scalar "$cc_csc_dir/change-set.yaml" tier)
	cc_csc_verifier_outcome=passed
	cc_csc_kind=inferred
	[ "$cc_csc_tier" = critical ] && cc_csc_kind=accepted
	# Preflight every member before changing any member status. A set is complete only
	# when all members still have a current execution, current candidate, and open plan.
	for cc_csc_m in $cc_csc_members; do
		cc_csc_yaml="$cc_csc_root/plans/$cc_csc_m/plan.yaml"
		[ "$(cc_scalar "$cc_csc_yaml" status 2>/dev/null)" = draft ] \
			|| { cc_fail CHANGE_SET_MEMBER_NOT_OPEN "$cc_csc_m"; return 1; }
		cc_csc_exec=$(cc_latest_execution "$cc_csc_root" "$cc_csc_m") || { cc_fail CHANGE_SET_NO_EXECUTION "$cc_csc_m"; return 1; }
		cc_csc_edir=$(cc_execution_dir "$cc_csc_root" "$cc_csc_m" "$cc_csc_exec")
		cc_candidate_id "$cc_csc_edir" >/dev/null || { cc_fail CHANGE_SET_MEMBER_CANDIDATE_INVALID "$cc_csc_m"; return 1; }
		cc_csc_revision=$(cc_scalar "$cc_csc_edir/execution.yaml" plan_revision 2>/dev/null) || cc_csc_revision=""
		[ -n "$cc_csc_revision" ] || { cc_fail CHANGE_SET_MEMBER_REVISION_MISSING "$cc_csc_m"; return 1; }
		cc_csc_commit_count=$(for cc_csc_rf in "$cc_csc_edir"/repositories/*.yaml; do [ -f "$cc_csc_rf" ] || continue; cc_scalar "$cc_csc_rf" latest_commit; done | sed '/^$/d' | wc -l | tr -d ' ')
		[ "$cc_csc_commit_count" -gt 0 ] || { cc_fail CHANGE_SET_MEMBER_NO_COMMIT "$cc_csc_m"; return 1; }
	done
	# Gate 2 happened for the whole set: record the change-set delivery once.
	printf 'schema_version: %s\ncandidate_id: %s\ndelivered_at: %s\n' \
		"$CC_DELIVERY_SCHEMA_VERSION" "$cc_csc_candidate" "$(cc_now)" \
		| cc_atomic_write "$cc_csc_dir/delivered.yaml" \
			|| { cc_fail CHANGE_SET_DELIVERY_RECORD_FAILED "$2"; return 1; }
	cc_knowledge_debt_emit_change_set "$cc_csc_root" "$2" >/dev/null \
		|| { cc_fail CHANGE_SET_DEBT_FAILED "$2"; return 1; }
	cc_csc_done=0
	for cc_csc_m in $cc_csc_members; do
		cc_csc_exec=$(cc_latest_execution "$cc_csc_root" "$cc_csc_m")
		cc_csc_edir=$(cc_execution_dir "$cc_csc_root" "$cc_csc_m" "$cc_csc_exec")
		cc_csc_revision=$(cc_scalar "$cc_csc_edir/execution.yaml" plan_revision)
		{
			printf 'schema_version: %s\nexecution_id: %s\nplan: %s\ncandidate_id: %s\nplan_revision: %s\naccepted_by: %s\nverifier_outcome: %s\nchange_set: %s\nhuman_completion: %s\ncompleted_at: %s\ncommits:\n' \
				"$CC_COMPLETION_SCHEMA_VERSION" "$cc_csc_exec" "$cc_csc_m" "$cc_csc_candidate" "$cc_csc_revision" "$cc_csc_acceptor" "$cc_csc_verifier_outcome" "$2" "$cc_csc_kind" "$(cc_now)"
			for cc_csc_rf in "$cc_csc_edir"/repositories/*.yaml; do
				[ -f "$cc_csc_rf" ] || continue
				printf '  %s: %s\n' "$(cc_scalar "$cc_csc_rf" repository)" "$(cc_scalar "$cc_csc_rf" latest_commit)"
			done
		} | cc_atomic_write "$cc_csc_edir/completion.yaml" \
			|| { cc_fail CHANGE_SET_MEMBER_RECORD_FAILED "$cc_csc_m"; return 1; }
		cc_plan_set_status "$cc_csc_root/plans/$cc_csc_m/plan.yaml" done \
			|| { cc_fail CHANGE_SET_MEMBER_STATUS_FAILED "$cc_csc_m"; return 1; }
		cc_plan_index_upsert "$cc_csc_root" "$cc_csc_m" >/dev/null
		cc_csc_done=$((cc_csc_done + 1))
	done
	cc_change_set_record_status "$cc_csc_dir" delivered-and-completed >/dev/null || return 1
	cc_emit change_set "$2"
	cc_emit completed "$cc_csc_done"
	cc_emit status delivered-and-completed
	return 0
}

# cc_human_acceptance_record EXEC_DIR ACCEPTED_BY [CHECKLIST_FILE] -> write a
# first-class human acceptance bound to the CURRENT candidate (INV-CANDIDATE-01).
# Acceptance names the candidate it observed, so any later commit or criteria
# change (which changes the candidate) leaves this acceptance no longer matching —
# void by construction, not by discipline. A checklist file, when given, records
# which criteria the human confirmed (including any manual ones).
cc_human_acceptance_record() {
	cc_ha_dir="$1"; cc_ha_by="$2"; cc_ha_file="${3:-}"
	[ -f "$cc_ha_dir/execution.yaml" ] || { cc_fail ACCEPTANCE_NO_EXECUTION; return 1; }
	[ -n "$cc_ha_by" ] || { cc_fail ACCEPTANCE_NO_ACCEPTER; return 1; }
	case "$cc_ha_by" in *[!A-Za-z0-9._@-]*) cc_fail ACCEPTANCE_ACCEPTER_INVALID "$cc_ha_by"; return 1 ;; esac
	cc_ha_cand=$(cc_candidate_id "$cc_ha_dir") || return 1
	{
		printf 'schema_version: %s\ncandidate_id: %s\naccepted_by: %s\naccepted_at: %s\nchecklist:\n' \
			"$CC_HUMAN_ACCEPTANCE_SCHEMA_VERSION" "$cc_ha_cand" "$cc_ha_by" "$(cc_now)"
		if [ -n "$cc_ha_file" ] && [ -f "$cc_ha_file" ]; then
			while IFS= read -r cc_ha_line; do
				[ -n "$cc_ha_line" ] || continue
				printf '  - %s\n' "$cc_ha_line"
			done <"$cc_ha_file"
		fi
	} | cc_atomic_write "$cc_ha_dir/human-acceptance.yaml" \
		|| { cc_fail ACCEPTANCE_RECORD_WRITE_FAILED; return 1; }
	cc_emit acceptance recorded
	cc_emit candidate_id "$cc_ha_cand"
	cc_emit accepted_by "$cc_ha_by"
	return 0
}

# cc_human_acceptance_current EXEC_DIR -> ok when a human-acceptance record exists
# AND still binds to the current candidate. A candidate change voids it: the record
# names an old candidate_id that no longer matches (INV-CANDIDATE-01).
cc_human_acceptance_current() {
	cc_hac_dir="$1"
	[ -f "$cc_hac_dir/human-acceptance.yaml" ] || { cc_fail ACCEPTANCE_MISSING; return 1; }
	cc_hac_bound=$(cc_scalar "$cc_hac_dir/human-acceptance.yaml" "candidate_id")
	cc_hac_now=$(cc_candidate_id "$cc_hac_dir") || return 1
	[ "$cc_hac_bound" = "$cc_hac_now" ] || { cc_fail ACCEPTANCE_VOID "$cc_hac_bound"; return 1; }
	cc_emit acceptance current
	cc_emit candidate_id "$cc_hac_now"
	return 0
}

# ---------------------------------------------------------------------------
# Consequence tiering (Context Circuit v1.0, Mechanism 3, the one safety-critical
# automated check, INV-ASSURE-01). Tier is DECLARED on the intent by the coordinator; the runtime
# NEVER selects a tier (INV-RUNTIME-01). What the runtime provides is a
# deterministic, fail-upward SIGNAL classifier over declared facts (scope paths and
# repository count) — the same kind of deterministic data extraction as the
# grounding scan, not model intelligence — plus the enforcement of the tier floor
# (completion-ready) and a guard that refuses dropping the independent verifier
# (tier explore) when any risk signal is present. Skipping the verifier at Explore
# is safe only if this fails upward: when unsure, tier higher, never lower.
# ---------------------------------------------------------------------------

# cc_tier_signals CONTRACT_FILE -> emit each detected risk-signal category, one per
# line, by matching the intent's scope path regions (and repository count) against
# transparent patterns. Deterministic; emits data, decides nothing.
cc_tier_signals() {
	cc_ts_file="$1"
	cc_ts_repos=$(cc_intent_scope_repos "$cc_ts_file")
	cc_ts_rc=$(printf '%s\n' "$cc_ts_repos" | sed '/^$/d' | wc -l | tr -d ' ')
	[ "$cc_ts_rc" -gt 1 ] && printf 'multi-repository\n'
	for cc_ts_r in $cc_ts_repos; do
		for cc_ts_p in $(cc_intent_scope_paths "$cc_ts_file" "$cc_ts_r"); do
			[ "$cc_ts_p" = "." ] && printf 'repository-wide-scope\n'
			cc_ts_lp=$(cc_lower "$cc_ts_p")
			case "$cc_ts_lp" in
				*auth*|*secret*|*credential*|*passw*|*oauth*|*token*|*security*|*/keys*|*/key/*|*login*|*session*) printf 'security\n' ;;
			esac
			case "$cc_ts_lp" in
				*payment*|*billing*|*charge*|*invoice*|*money*|*wallet*|*checkout*|*refund*) printf 'money\n' ;;
			esac
			case "$cc_ts_lp" in
				*migration*|*migrate*|*schema*|*/db/*|*database*) printf 'data-migration\n' ;;
			esac
			case "$cc_ts_lp" in
				*deploy*|*production*|*/prod/*|*infra*|*release*|*/ci/*|*/k8s/*) printf 'production\n' ;;
			esac
		done
	done | LC_ALL=C sort -u
}

# cc_tier_classify ROOT INTENT -> emit the detected signals, a classified tier, and
# explore_ok. Fail upward: a hard signal (security/money/migration/production) is
# critical; a soft signal (multi-repo/repo-wide) is standard; and a bounded single-
# repo intent with no signal ALSO defaults to standard — the engine cannot see
# reversibility, coverage, or novelty, so absence of a signal is not proof of low
# risk. explore_ok reports whether a human MAY lower to Explore: yes only for a
# bounded single repo with no risk signal; no whenever any signal is present.
cc_tier_classify() {
	cc_tc_dir=$(cc_intent_dir "$1" "$2")
	cc_tc_file="$cc_tc_dir/contract.yaml"
	[ -f "$cc_tc_file" ] || { cc_fail TIER_INTENT_MISSING "$2"; return 1; }
	cc_tc_sigs=$(cc_tier_signals "$cc_tc_file")
	cc_tc_hard=no
	printf '%s\n' "$cc_tc_sigs" | grep -Eq '^(security|money|data-migration|production)$' && cc_tc_hard=yes
	cc_tc_soft=no
	printf '%s\n' "$cc_tc_sigs" | grep -Eq '^(multi-repository|repository-wide-scope)$' && cc_tc_soft=yes
	if [ "$cc_tc_hard" = "yes" ]; then
		cc_tc_tier=critical; cc_tc_explore=no
	elif [ "$cc_tc_soft" = "yes" ]; then
		cc_tc_tier=standard; cc_tc_explore=no
	else
		# No detected risk signal on a single bounded repository. The engine
		# cannot see reversibility, coverage, or novelty, so it does NOT auto-
		# classify Explore: the default fails upward to Standard. Explore stays
		# available as an explicit human lowering (explore_ok=yes), recorded by
		# the human declaring tier: explore on the intent.
		cc_tc_tier=standard; cc_tc_explore=yes
	fi
	for cc_tc_s in $cc_tc_sigs; do cc_emit signal "$cc_tc_s"; done
	cc_emit classified_tier "$cc_tc_tier"
	cc_emit explore_ok "$cc_tc_explore"
	return 0
}

# cc_tier_lower_check ROOT INTENT REQUESTED -> refuse dropping the independent
# verifier (requesting Explore) when any risk signal is present. Standard<->Critical
# both keep the verifier and are a human decision. Fails upward on any risk.
cc_tier_lower_check() {
	cc_tl_req="$3"
	case "$cc_tl_req" in
		explore|standard|critical) : ;;
		*) cc_fail TIER_REQUEST_INVALID "$cc_tl_req"; return 1 ;;
	esac
	cc_tl_ok=$(cc_tier_classify "$1" "$2" | sed -n 's/^explore_ok: //p') || return 1
	if [ "$cc_tl_req" = "explore" ] && [ "$cc_tl_ok" != "yes" ]; then
		cc_emit tier_lower refused
		cc_fail TIER_EXPLORE_UNSAFE "$2"; return 1
	fi
	cc_emit tier_lower allowed
	cc_emit tier "$cc_tl_req"
	return 0
}

# ---------------------------------------------------------------------------
# Readiness and the run-stack partition (deterministic; the coordinator decides
# how many ready plans to launch — the runtime never schedules, INV-RUNTIME-01).
# ---------------------------------------------------------------------------

# cc_plan_latest_status ROOT PLAN -> the plan's latest execution status, or "none".
cc_plan_latest_status() {
	cc_pls_e=$(cc_latest_execution "$1" "$2" 2>/dev/null) || cc_pls_e=""
	[ -n "$cc_pls_e" ] || { printf 'none'; return 0; }
	cc_execution_status "$(cc_execution_dir "$1" "$2" "$cc_pls_e")"
}

# cc_plan_ready ROOT PLAN -> readiness AND-join: every plan dependency verified
# AND every needed path region free in each repository. Emits readiness (ready |
# waiting | blocked) and a reason; returns 0 only when ready.
cc_plan_ready() {
	# tolerate a missing workspace root when invoked from within the workspace:
	# `plan-ready <plan>` means root ".".
	if [ -n "$2" ] && { [ -f "$1/workspace.yaml" ] || [ -d "$1/plans" ]; }; then
		cc_pr_root="$1"; cc_pr_plan="$2"
	else
		cc_pr_root="."; cc_pr_plan="$1"
	fi
	cc_pr_pf="$cc_pr_root/plans/$cc_pr_plan/plan.yaml"
	[ -f "$cc_pr_pf" ] || { cc_fail PLAN_YAML_MISSING "$cc_pr_plan"; return 1; }
	cc_plan_validate "$cc_pr_root/plans/$cc_pr_plan" >/dev/null || {
		cc_emit readiness blocked; cc_emit reason PLAN_INVALID; return 1;
	}
	cc_intent_authorized "$cc_pr_root" "$cc_pr_plan" >/dev/null || {
		cc_emit readiness blocked; cc_emit reason INTENT_UNAUTHORIZED; return 1;
	}
	# 1. dependency AND-join
	for cc_pr_dep in $(cc_plan_dependencies "$cc_pr_pf"); do
		cc_pr_ds=$(cc_plan_latest_status "$cc_pr_root" "$cc_pr_dep")
		case "$cc_pr_ds" in
			verified) : ;;
			failed|blocked)
				cc_emit readiness blocked
				cc_emit reason "DEP_FAILED:$cc_pr_dep"
				return 1 ;;
			*)
				cc_emit readiness waiting
				cc_emit reason "DEP_NOT_VERIFIED:$cc_pr_dep"
				return 1 ;;
		esac
	done
	# 2. lease gate: each affected repository's needed regions must be free
	for cc_pr_repo in $(cc_plan_affected_repositories "$cc_pr_pf"); do
		cc_pr_reg=$(cc_plan_repo_paths "$cc_pr_pf" "$cc_pr_repo" | tr '\n' ' ' | sed 's/  */ /g; s/^ //; s/ $//')
		[ -n "$cc_pr_reg" ] || cc_pr_reg="."
		if ! cc_lease_check "$cc_pr_root" "$cc_pr_repo" "$cc_pr_plan" "$cc_pr_reg" >/dev/null 2>&1; then
			cc_emit readiness waiting
			cc_emit reason "LEASE_BUSY:$cc_pr_repo"
			return 1
		fi
	done
	cc_emit readiness ready
	return 0
}

# cc_run_stack_ready ROOT PLAN... -> partition a set of plans into
# verified / ready / waiting / failed / blocked / refused. A plan whose latest
# execution is terminal (verified/failed/blocked) or in-progress
# (running/verifying/repairing) is NOT runnable; only an authorized, never-run plan
# that passes readiness is `ready`. Emits one "<plan>: <bucket>" line per plan.
cc_run_stack_ready() {
	# tolerate invocation from within the workspace: the root defaults to "." when
	# the first argument is not a workspace root (a real root has plans/), so
	# `run-stack-ready <plan> <plan> ...` works with the current directory.
	if [ -f "$1/workspace.yaml" ] || [ -d "$1/plans" ]; then
		cc_rsr_root="$1"; shift 2>/dev/null || true
	else
		cc_rsr_root="."
	fi
	for cc_rsr_plan in "$@"; do
		[ -n "$cc_rsr_plan" ] || continue
		# resolve a bare four-digit prefix (e.g. 0001) to its unique full plan id
		if [ ! -f "$cc_rsr_root/plans/$cc_rsr_plan/plan.yaml" ]; then
			case "$cc_rsr_plan" in
				[0-9][0-9][0-9][0-9])
					for cc_rsr_m in "$cc_rsr_root/plans/$cc_rsr_plan-"*/; do
						[ -d "$cc_rsr_m" ] && { cc_rsr_plan=$(basename -- "$cc_rsr_m"); break; }
					done ;;
			esac
		fi
		# v1.0 authorization: a plan is authorized to run when it derives from an
		# approved intent whose criteria are unchanged (INV-INTENT-02 / INV-EXEC-01).
		# There is no separate plan-approval status and no automated scope gate; a plan
		# that cannot be authorized — because its intent is not approved or its criteria
		# drifted since approval — is refused.
		if ! cc_intent_authorized "$cc_rsr_root" "$cc_rsr_plan" >/dev/null 2>&1; then
			cc_emit "$cc_rsr_plan" refused
			continue
		fi
		cc_rsr_le=$(cc_plan_latest_status "$cc_rsr_root" "$cc_rsr_plan")
		case "$cc_rsr_le" in
			verified) cc_emit "$cc_rsr_plan" verified; continue ;;
			failed)   cc_emit "$cc_rsr_plan" failed;   continue ;;
			blocked)  cc_emit "$cc_rsr_plan" blocked;  continue ;;
			running|verifying|repairing) cc_emit "$cc_rsr_plan" waiting; continue ;;
		esac
		# never executed: evaluate readiness
		cc_rsr_out=$(cc_plan_ready "$cc_rsr_root" "$cc_rsr_plan" 2>/dev/null) || :
		cc_rsr_r=$(printf '%s' "$cc_rsr_out" | sed -n 's/^readiness: //p')
		case "$cc_rsr_r" in
			ready)   cc_emit "$cc_rsr_plan" ready ;;
			blocked) cc_emit "$cc_rsr_plan" blocked ;;
			*)       cc_emit "$cc_rsr_plan" waiting ;;
		esac
	done
	return 0
}

# ---------------------------------------------------------------------------
# Completion (human-controlled) and knowledge-impact handoff
# ---------------------------------------------------------------------------

# cc_latest_execution ROOT PLAN -> newest execution id for a plan
cc_latest_execution() {
	cc_le_base="$1/.runtime/executions/$2"
	[ -d "$cc_le_base" ] || return 1
	ls -1 "$cc_le_base" 2>/dev/null | grep '^exec-' | sort | tail -n1
}

# cc_completion_ready ROOT PLAN -> eligible only if the latest execution's evidence
# still describes the current candidate (INV-CANDIDATE-01): a new commit or criteria
# change after the pass voids the evidence, so completion is refused until re-verified.
# The tier floor (INV-ASSURE-01) decides what evidence is required.
cc_completion_ready() {
	cc_cr_root="$1"; cc_cr_plan="$2"
	cc_intent_authorized "$cc_cr_root" "$cc_cr_plan" >/dev/null \
		|| { cc_fail COMPLETION_INTENT_UNAUTHORIZED "$cc_cr_plan"; return 1; }
	cc_cr_exec=$(cc_latest_execution "$cc_cr_root" "$cc_cr_plan") || { cc_fail COMPLETION_NO_EXECUTION; return 1; }
	[ -n "$cc_cr_exec" ] || { cc_fail COMPLETION_NO_EXECUTION; return 1; }
	cc_cr_dir=$(cc_execution_dir "$cc_cr_root" "$cc_cr_plan" "$cc_cr_exec")
	cc_cr_tier=$(cc_scalar "$cc_cr_dir/execution.yaml" "tier" 2>/dev/null) || cc_cr_tier=""
	[ -n "$cc_cr_tier" ] || cc_cr_tier=standard
	if [ "$cc_cr_tier" = "explore" ]; then
		# Explore is human-supervised (INV-ASSURE-01): no independent verifier, and
		# the result is NEVER labeled "verified". Eligibility rests on a current
		# candidate-bound human acceptance, not a verifier pass.
		cc_human_acceptance_current "$cc_cr_dir" >/dev/null || { cc_fail COMPLETION_NO_ACCEPTANCE "$cc_cr_plan"; return 1; }
		cc_emit completion eligible
		cc_emit execution_id "$cc_cr_exec"
		cc_emit assurance human-supervised
		return 0
	fi
	# Standard/Critical require the independent verifier floor: a candidate-bound
	# `passed` that still describes the current candidate (INV-VERIFY-01, INV-CANDIDATE-01).
	cc_cr_st=$(cc_execution_status "$cc_cr_dir")
	[ "$cc_cr_st" = "verified" ] || { cc_fail COMPLETION_NOT_VERIFIED "$cc_cr_st"; return 1; }
	cc_cr_vc=$(cc_scalar "$cc_cr_dir/execution.yaml" "verified_candidate" 2>/dev/null) || cc_cr_vc=""
	[ -n "$cc_cr_vc" ] || { cc_fail COMPLETION_VERIFIER_CANDIDATE_MISSING "$cc_cr_plan"; return 1; }
	cc_cr_now=$(cc_candidate_id "$cc_cr_dir" 2>/dev/null) || cc_cr_now=""
	[ "$cc_cr_vc" = "$cc_cr_now" ] || { cc_fail COMPLETION_CANDIDATE_STALE "$cc_cr_vc"; return 1; }
	cc_human_acceptance_current "$cc_cr_dir" >/dev/null \
		|| { cc_fail COMPLETION_NO_ACCEPTANCE "$cc_cr_plan"; return 1; }
	cc_emit completion eligible
	cc_emit execution_id "$cc_cr_exec"
	cc_emit assurance independent
	[ -n "$cc_cr_vc" ] && cc_emit candidate_id "$cc_cr_vc" || :
	return 0
}

# cc_completion_finalize ROOT PLAN KIND -> shared completion writer (KIND is the
# human_completion value: accepted for an explicit human completion, inferred for a
# projection of candidate acceptance + delivery). Writes the completion record, sets
# status done, and emits the reconciliation-debt marker (INV-COMPLETE-02).
cc_completion_finalize() {
	cc_cf_root="$1"; cc_cf_plan="$2"; cc_cf_kind="$3"
	cc_cf_exec=$(cc_latest_execution "$cc_cf_root" "$cc_cf_plan")
	cc_cf_edir=$(cc_execution_dir "$cc_cf_root" "$cc_cf_plan" "$cc_cf_exec")
	cc_cf_yaml="$cc_cf_root/plans/$cc_cf_plan/plan.yaml"
	cc_cf_status=$(cc_scalar "$cc_cf_yaml" "status")
	[ "$cc_cf_status" = "draft" ] || { cc_fail COMPLETION_PLAN_NOT_OPEN "$cc_cf_status"; return 1; }
	cc_cf_candidate=$(cc_candidate_id "$cc_cf_edir") || return 1
	cc_cf_revision=$(cc_scalar "$cc_cf_edir/execution.yaml" "plan_revision" 2>/dev/null) || cc_cf_revision=""
	[ -n "$cc_cf_revision" ] || { cc_fail COMPLETION_PLAN_REVISION_MISSING "$cc_cf_plan"; return 1; }
	cc_cf_acceptor=$(cc_scalar "$cc_cf_edir/human-acceptance.yaml" "accepted_by" 2>/dev/null) || cc_cf_acceptor=""
	[ -n "$cc_cf_acceptor" ] || { cc_fail COMPLETION_ACCEPTOR_MISSING "$cc_cf_plan"; return 1; }
	cc_cf_verified=$(cc_scalar "$cc_cf_edir/execution.yaml" "verified_candidate" 2>/dev/null) || cc_cf_verified=""
	cc_cf_verdict=not-run
	[ -n "$cc_cf_verified" ] && cc_cf_verdict=passed
	# Emit reconciliation debt before changing the plan status. A debt write failure
	# therefore leaves the plan draft and retryable, rather than silently completing
	# while losing the knowledge-loop marker (INV-COMPLETE-02).
	cc_knowledge_debt_emit "$cc_cf_root" "$cc_cf_plan" "$cc_cf_exec" >/dev/null \
		|| { cc_fail COMPLETION_DEBT_FAILED "$cc_cf_plan"; return 1; }
	{
		printf 'schema_version: %s\nexecution_id: %s\nplan: %s\ncandidate_id: %s\nplan_revision: %s\naccepted_by: %s\nverifier_outcome: %s\nhuman_completion: %s\ncompleted_at: %s\ncommits:\n' \
			"$CC_COMPLETION_SCHEMA_VERSION" "$cc_cf_exec" "$cc_cf_plan" "$cc_cf_candidate" "$cc_cf_revision" "$cc_cf_acceptor" "$cc_cf_verdict" "$cc_cf_kind" "$(cc_now)"
		for cc_cf_rf in "$cc_cf_edir"/repositories/*.yaml; do
			[ -f "$cc_cf_rf" ] || continue
			printf '  %s: %s\n' "$(cc_scalar "$cc_cf_rf" repository)" "$(cc_scalar "$cc_cf_rf" latest_commit)"
		done
		} | cc_atomic_write "$cc_cf_edir/completion.yaml" \
			|| { cc_fail COMPLETION_RECORD_WRITE_FAILED; return 1; }
	cc_plan_set_status "$cc_cf_yaml" "done" || { cc_fail COMPLETION_STATUS_UPDATE_FAILED "$cc_cf_plan"; return 1; }
	cc_plan_index_upsert "$cc_cf_root" "$cc_cf_plan" >/dev/null
	cc_emit plan "$cc_cf_plan"
	cc_emit status done
	cc_emit human_completion "$cc_cf_kind"
	cc_emit implementation_completion recorded
	cc_emit knowledge_impact review-pending
	cc_emit reconciliation_debt pending
	return 0
}

# cc_plan_complete ROOT PLAN -> the EXPLICIT human completion path. Critical is the
# only tier that uses it; Explore/Standard complete by the delivery + acceptance
# inference path (INV-COMPLETE-01).
cc_plan_complete() {
	cc_pc_exec=$(cc_latest_execution "$1" "$2") || { cc_fail COMPLETION_NO_EXECUTION; return 1; }
	cc_pc_edir=$(cc_execution_dir "$1" "$2" "$cc_pc_exec")
	cc_pc_tier=$(cc_scalar "$cc_pc_edir/execution.yaml" tier 2>/dev/null) || cc_pc_tier=standard
	[ "$cc_pc_tier" = critical ] || { cc_fail COMPLETION_EXPLICIT_CRITICAL_ONLY "$2"; return 1; }
	cc_completion_ready "$1" "$2" >/dev/null || { cc_fail COMPLETION_BLOCKED; return 1; }
	cc_completion_finalize "$1" "$2" accepted
}

# cc_delivery_record ROOT PLAN [EXEC] -> record that the plan's current candidate was
# delivered (Gate 2 happened). It is the delivery signal inferred completion reads;
# it performs no git action itself (delivery is a separate coordinator/host act,
# INV-DELIVER-01). Bound to the candidate so a post-delivery change is visible.
cc_delivery_record() {
	cc_del_root="$1"; cc_del_plan="$2"; cc_del_exec="${3:-}"
	[ -n "$cc_del_exec" ] || cc_del_exec=$(cc_latest_execution "$cc_del_root" "$cc_del_plan") || { cc_fail DELIVERY_NO_EXECUTION "$cc_del_plan"; return 1; }
	[ -n "$cc_del_exec" ] || { cc_fail DELIVERY_NO_EXECUTION "$cc_del_plan"; return 1; }
	cc_del_edir=$(cc_execution_dir "$cc_del_root" "$cc_del_plan" "$cc_del_exec")
	[ -f "$cc_del_edir/execution.yaml" ] || { cc_fail DELIVERY_NO_EXECUTION "$cc_del_plan"; return 1; }
	cc_del_cand=$(cc_candidate_id "$cc_del_edir") || return 1
	printf 'schema_version: %s\ncandidate_id: %s\nplan: %s\nexecution_id: %s\ndelivered_at: %s\n' \
		"$CC_DELIVERY_SCHEMA_VERSION" "$cc_del_cand" "$cc_del_plan" "$cc_del_exec" "$(cc_now)" \
		| cc_atomic_write "$cc_del_edir/delivered.yaml" \
		|| { cc_fail DELIVERY_RECORD_WRITE_FAILED; return 1; }
	cc_knowledge_debt_emit "$cc_del_root" "$cc_del_plan" "$cc_del_exec" >/dev/null \
		|| { cc_fail DELIVERY_DEBT_FAILED "$cc_del_plan"; return 1; }
	cc_emit delivery recorded
	cc_emit candidate_id "$cc_del_cand"
	return 0
}

# cc_completion_infer ROOT PLAN -> INFERRED completion (Context Circuit v1.0). At
# Standard, completion is a projection of "candidate accepted + delivered": it
# requires completion-ready AND a delivery record that still binds to
# the current candidate. At Critical it is refused — an explicit human completion is
# required (INV-COMPLETE-01). Idempotent: an already-done plan reports done.
cc_completion_infer() {
	cc_ci_root="$1"; cc_ci_plan="$2"
	cc_ci_yaml="$cc_ci_root/plans/$cc_ci_plan/plan.yaml"
	[ -f "$cc_ci_yaml" ] || { cc_fail COMPLETION_INFER_PLAN_MISSING "$cc_ci_plan"; return 1; }
	if [ "$(cc_scalar "$cc_ci_yaml" status)" = "done" ]; then
		cc_emit plan "$cc_ci_plan"; cc_emit status done; return 0
	fi
	cc_ci_exec=$(cc_latest_execution "$cc_ci_root" "$cc_ci_plan") || { cc_fail COMPLETION_NO_EXECUTION; return 1; }
	cc_ci_edir=$(cc_execution_dir "$cc_ci_root" "$cc_ci_plan" "$cc_ci_exec")
	cc_ci_tier=$(cc_scalar "$cc_ci_edir/execution.yaml" "tier" 2>/dev/null) || cc_ci_tier=standard
	[ -n "$cc_ci_tier" ] || cc_ci_tier=standard
	if [ "$cc_ci_tier" = "critical" ]; then
		cc_fail COMPLETION_INFER_REQUIRES_EXPLICIT "$cc_ci_plan"; return 1
	fi
	cc_completion_ready "$cc_ci_root" "$cc_ci_plan" >/dev/null || { cc_fail COMPLETION_BLOCKED; return 1; }
	# the delivery signal must exist and still describe the current candidate
	[ -f "$cc_ci_edir/delivered.yaml" ] || { cc_fail COMPLETION_NOT_DELIVERED "$cc_ci_plan"; return 1; }
	cc_ci_delc=$(cc_scalar "$cc_ci_edir/delivered.yaml" "candidate_id")
	cc_ci_now=$(cc_candidate_id "$cc_ci_edir") || return 1
	[ "$cc_ci_delc" = "$cc_ci_now" ] || { cc_fail COMPLETION_DELIVERY_STALE "$cc_ci_delc"; return 1; }
	cc_completion_finalize "$cc_ci_root" "$cc_ci_plan" inferred
}

# cc_context_impact_record EXEC_DIR FILE -> store reconciliation refs (no PK interpretation)
cc_context_impact_record() {
	[ -f "$2" ] || { cc_fail CONTEXT_IMPACT_FILE_MISSING; return 1; }
	cat "$2" | cc_atomic_write "$1/context-impact.yaml" || { cc_fail CONTEXT_IMPACT_WRITE_FAILED; return 1; }
	cc_emit context_impact recorded
	return 0
}

# ---------------------------------------------------------------------------
# Closed knowledge loop (Context Circuit v1.0, Mechanism 4, INV-COMPLETE-02 /
# INV-KNOWLEDGE-02). Completion or delivery of a candidate emits a
# reconciliation-debt marker; the next Standard/Critical plan's grounding preflight
# blocks while delivered work in its knowledge scope remains unreconciled. Explore
# is planless and has no plan grounding preflight. The human still ACCEPTS knowledge (the gate
# never auto-accepts); the loop only refuses to let the debt be FORGOTTEN. The debt
# marker is keyed to the candidate, so it inherits candidate honesty.
# ---------------------------------------------------------------------------

cc_knowledge_debt_dir() { printf '%s/.runtime/knowledge-debt' "$1"; }

# cc_knowledge_debt_emit ROOT PLAN [EXEC] -> record a pending reconciliation-debt
# marker for the plan's delivered/completed candidate, keyed to the candidate id.
cc_knowledge_debt_emit() {
	cc_kde_root="$1"; cc_kde_plan="$2"; cc_kde_exec="${3:-}"
	[ -n "$cc_kde_exec" ] || cc_kde_exec=$(cc_latest_execution "$cc_kde_root" "$cc_kde_plan") || { cc_fail KNOWLEDGE_DEBT_NO_EXECUTION "$cc_kde_plan"; return 1; }
	[ -n "$cc_kde_exec" ] || { cc_fail KNOWLEDGE_DEBT_NO_EXECUTION "$cc_kde_plan"; return 1; }
	cc_kde_edir=$(cc_execution_dir "$cc_kde_root" "$cc_kde_plan" "$cc_kde_exec")
	[ -f "$cc_kde_edir/execution.yaml" ] || { cc_fail KNOWLEDGE_DEBT_NO_EXECUTION "$cc_kde_plan"; return 1; }
	cc_kde_cand=$(cc_candidate_id "$cc_kde_edir") || return 1
	# knowledge scope: the plan's affected repositories and its relied-on PK unit ids
	cc_kde_snap="$cc_kde_edir/snapshot/plan.yaml"
	[ -f "$cc_kde_snap" ] || cc_kde_snap="$cc_kde_root/plans/$cc_kde_plan/plan.yaml"
	cc_kde_repos=$(cc_plan_affected_repositories "$cc_kde_snap" | tr '\n' ',' | sed 's/,$//; s/,/, /g')
	cc_kde_units=$(cc_list_ids "$cc_kde_snap" "product_knowledge" | tr '\n' ',' | sed 's/,$//; s/,/, /g')
	mkdir -p "$(cc_knowledge_debt_dir "$cc_kde_root")"
	cc_kde_file="$(cc_knowledge_debt_dir "$cc_kde_root")/$cc_kde_cand.yaml"
	# preserve an already-resolved marker for the same candidate (idempotent emit)
	if [ -f "$cc_kde_file" ]; then
		cc_kde_prev=$(cc_scalar "$cc_kde_file" "resolved" 2>/dev/null) || cc_kde_prev=""
		[ "$cc_kde_prev" = "reconciled" ] || [ "$cc_kde_prev" = "deferred" ] && { cc_emit debt already-resolved; cc_emit candidate_id "$cc_kde_cand"; return 0; }
	fi
	printf 'schema_version: %s\ncandidate_id: %s\nplan: %s\nexecution_id: %s\nrepositories: [%s]\nknowledge_units: [%s]\nresolved: pending\ncreated_at: %s\n' \
		"$CC_KNOWLEDGE_DEBT_SCHEMA_VERSION" "$cc_kde_cand" "$cc_kde_plan" "$cc_kde_exec" "$cc_kde_repos" "$cc_kde_units" "$(cc_now)" \
		| cc_atomic_write "$cc_kde_file" \
		|| { cc_fail KNOWLEDGE_DEBT_WRITE_FAILED "$cc_kde_cand"; return 1; }
	cc_emit debt recorded
	cc_emit candidate_id "$cc_kde_cand"
	return 0
}

# cc_knowledge_debt_emit_change_set ROOT CS_ID -> emit one reconciliation marker for
# the integrated candidate. The marker is deliberately set-bound: one delivered
# change set creates one candidate and therefore one knowledge-loop obligation.
cc_knowledge_debt_emit_change_set() {
	cc_kdcs_root="$1"; cc_kdcs_id="$2"; cc_kdcs_dir=$(cc_change_set_dir "$cc_kdcs_root" "$cc_kdcs_id")
	[ -f "$cc_kdcs_dir/change-set.yaml" ] || { cc_fail KNOWLEDGE_DEBT_CHANGE_SET_MISSING "$cc_kdcs_id"; return 1; }
	cc_kdcs_cand=$(cc_scalar "$cc_kdcs_dir/change-set.yaml" candidate_id)
	[ -n "$cc_kdcs_cand" ] && [ "$cc_kdcs_cand" != pending ] || { cc_fail KNOWLEDGE_DEBT_CANDIDATE_MISSING "$cc_kdcs_id"; return 1; }
	cc_kdcs_members=$(cc_inline_list "$(cc_scalar "$cc_kdcs_dir/change-set.yaml" members)")
	cc_kdcs_repos=$(for cc_kdcs_m in $cc_kdcs_members; do cc_plan_affected_repositories "$cc_kdcs_root/plans/$cc_kdcs_m/plan.yaml"; done | LC_ALL=C sort -u | tr '\n' ',' | sed 's/,$//; s/,/, /g')
	cc_kdcs_units=$(for cc_kdcs_m in $cc_kdcs_members; do
		cc_kdcs_e=$(cc_latest_execution "$cc_kdcs_root" "$cc_kdcs_m") || exit 1
		cc_kdcs_s="$cc_kdcs_root/.runtime/executions/$cc_kdcs_m/$cc_kdcs_e/snapshot/plan.yaml"
		[ -f "$cc_kdcs_s" ] || cc_kdcs_s="$cc_kdcs_root/plans/$cc_kdcs_m/plan.yaml"
		cc_list_ids "$cc_kdcs_s" product_knowledge
	done | LC_ALL=C sort -u | tr '\n' ',' | sed 's/,$//; s/,/, /g') || { cc_fail KNOWLEDGE_DEBT_CHANGE_SET_SCOPE_FAILED "$cc_kdcs_id"; return 1; }
	mkdir -p "$(cc_knowledge_debt_dir "$cc_kdcs_root")"
	cc_kdcs_file="$(cc_knowledge_debt_dir "$cc_kdcs_root")/$cc_kdcs_cand.yaml"
	if [ -f "$cc_kdcs_file" ]; then
		cc_kdcs_prev=$(cc_scalar "$cc_kdcs_file" resolved 2>/dev/null) || cc_kdcs_prev=""
		[ "$cc_kdcs_prev" = reconciled ] || [ "$cc_kdcs_prev" = deferred ] && { cc_emit debt already-resolved; cc_emit candidate_id "$cc_kdcs_cand"; return 0; }
	fi
	printf 'schema_version: %s\ncandidate_id: %s\nplan: change-set:%s\nexecution_id: %s\nrepositories: [%s]\nknowledge_units: [%s]\nresolved: pending\ncreated_at: %s\n' \
		"$CC_KNOWLEDGE_DEBT_SCHEMA_VERSION" "$cc_kdcs_cand" "$cc_kdcs_id" "$cc_kdcs_id" "$cc_kdcs_repos" "$cc_kdcs_units" "$(cc_now)" \
		| cc_atomic_write "$cc_kdcs_file" \
		|| { cc_fail KNOWLEDGE_DEBT_WRITE_FAILED "$cc_kdcs_cand"; return 1; }
	cc_emit debt recorded
	cc_emit candidate_id "$cc_kdcs_cand"
	return 0
}

# cc_knowledge_debt ROOT -> list every delivered candidate whose reconciliation is
# still pending. Emits one debt line per marker plus a pending_count.
cc_knowledge_debt() {
	cc_kd_dir=$(cc_knowledge_debt_dir "$1")
	cc_kd_n=0
	if [ -d "$cc_kd_dir" ]; then
		for cc_kd_f in "$cc_kd_dir"/*.yaml; do
			[ -f "$cc_kd_f" ] || continue
			cc_kd_res=$(cc_scalar "$cc_kd_f" "resolved" 2>/dev/null) || cc_kd_res=""
			[ "$cc_kd_res" = "pending" ] || continue
			cc_kd_n=$((cc_kd_n + 1))
			cc_emit debt "$(cc_scalar "$cc_kd_f" candidate_id)"
			cc_emit "  plan" "$(cc_scalar "$cc_kd_f" plan)"
			cc_emit "  repositories" "$(cc_scalar "$cc_kd_f" repositories)"
		done
	fi
	cc_emit pending_count "$cc_kd_n"
	return 0
}

# cc_knowledge_reconciled ROOT CANDIDATE [RESOLUTION] -> clear a debt marker once the
# impact proposals have been generated and either accepted or explicitly deferred by
# a human. RESOLUTION is reconciled (default) or deferred; "deferred" is a
# first-class "no durable update needed" outcome (someone decided). It never
# accepts knowledge itself (INV-KNOWLEDGE-02).
cc_knowledge_reconciled() {
	cc_kr_root="$1"; cc_kr_cand="$2"; cc_kr_res="${3:-reconciled}"
	case "$cc_kr_res" in reconciled|deferred) : ;; *) cc_fail KNOWLEDGE_RESOLUTION_INVALID "$cc_kr_res"; return 1 ;; esac
	cc_kr_file="$(cc_knowledge_debt_dir "$cc_kr_root")/$cc_kr_cand.yaml"
	[ -f "$cc_kr_file" ] || { cc_fail KNOWLEDGE_DEBT_UNKNOWN "$cc_kr_cand"; return 1; }
	awk -v r="$cc_kr_res" -v ts="$(cc_now)" '
		/^resolved:[[:space:]]/ { print "resolved: " r; next }
		{ print }
		END { print "resolved_at: " ts }
	' "$cc_kr_file" | cc_atomic_write "$cc_kr_file" \
		|| { cc_fail KNOWLEDGE_RECONCILE_WRITE_FAILED "$cc_kr_cand"; return 1; }
	cc_emit reconciled "$cc_kr_cand"
	cc_emit resolution "$cc_kr_res"
	return 0
}

# cc_knowledge_debt_check ROOT PLAN -> the cc-plan grounding preflight. If any
# pending debt marker's knowledge scope (repositories) overlaps the new plan's
# affected repositories, block at Standard/Critical (non-zero) or loudly warn at
# Explore (zero). No overlap is clear. The new plan's tier comes from its intent.
cc_knowledge_debt_check() {
	cc_kc_root="$1"; cc_kc_plan="$2"
	cc_kc_pfile="$cc_kc_root/plans/$cc_kc_plan/plan.yaml"
	[ -f "$cc_kc_pfile" ] || { cc_fail KNOWLEDGE_CHECK_PLAN_MISSING "$cc_kc_plan"; return 1; }
	cc_kc_newrepos=$(cc_plan_affected_repositories "$cc_kc_pfile")
	cc_kc_newunits=$(cc_list_ids "$cc_kc_pfile" product_knowledge)
	# resolve the new plan's tier from its intent
	cc_kc_intent=$(cc_scalar "$cc_kc_pfile" "intent" 2>/dev/null) || cc_kc_intent=""
	cc_kc_tier=standard
	if [ -n "$cc_kc_intent" ]; then
		cc_kc_tier=$(cc_scalar "$(cc_intent_dir "$cc_kc_root" "$cc_kc_intent")/contract.yaml" "tier" 2>/dev/null) || cc_kc_tier=standard
		[ -n "$cc_kc_tier" ] || cc_kc_tier=standard
	fi
	[ "$cc_kc_tier" != explore ] || { cc_fail KNOWLEDGE_EXPLORE_PLANLESS "$cc_kc_plan"; return 1; }
	cc_kc_dir=$(cc_knowledge_debt_dir "$cc_kc_root")
	cc_kc_hit=no
	if [ -d "$cc_kc_dir" ]; then
		for cc_kc_f in "$cc_kc_dir"/*.yaml; do
			[ -f "$cc_kc_f" ] || continue
			[ "$(cc_scalar "$cc_kc_f" resolved 2>/dev/null)" = "pending" ] || continue
			# skip a marker for the plan's own prior candidate (self-debt does not block)
			[ "$(cc_scalar "$cc_kc_f" plan 2>/dev/null)" = "$cc_kc_plan" ] && continue
			cc_kc_mrepos=$(cc_inline_list "$(cc_scalar "$cc_kc_f" repositories 2>/dev/null)")
			cc_kc_munits=$(cc_inline_list "$(cc_scalar "$cc_kc_f" knowledge_units 2>/dev/null)")
			cc_kc_marker_hit=no
			for cc_kc_mr in $cc_kc_mrepos; do
				if printf '%s\n' "$cc_kc_newrepos" | grep -Fxq "$cc_kc_mr"; then
					cc_kc_hit=yes; cc_kc_marker_hit=yes
					cc_emit debtor "$(cc_scalar "$cc_kc_f" candidate_id)"
					cc_emit "  plan" "$(cc_scalar "$cc_kc_f" plan)"
					cc_emit "  repository" "$cc_kc_mr"
					break
				fi
			done
			if [ "$cc_kc_marker_hit" = no ]; then
				for cc_kc_mu in $cc_kc_munits; do
					if printf '%s\n' "$cc_kc_newunits" | grep -Fxq "$cc_kc_mu"; then
						cc_kc_hit=yes
						cc_emit debtor "$(cc_scalar "$cc_kc_f" candidate_id)"
						cc_emit "  plan" "$(cc_scalar "$cc_kc_f" plan)"
						cc_emit "  knowledge_unit" "$cc_kc_mu"
						break
					fi
				done
			fi
		done
	fi
	if [ "$cc_kc_hit" = "no" ]; then
		cc_emit debt clear
		return 0
	fi
	if [ "$cc_kc_tier" = "explore" ]; then
		cc_emit debt warn
		return 0
	fi
	cc_emit debt blocking
	cc_fail KNOWLEDGE_DEBT_BLOCKING "$cc_kc_plan"; return 1
}

# cc_delivery_targets ROOT PLAN -> read-only report of pull-request source and
# default target for each affected repository of the latest execution. It never
# pushes, merges, or opens a pull request; it only reports and flags gaps.
cc_delivery_targets() {
	cc_dt_root="$1"; cc_dt_plan="$2"
	cc_dt_exec=$(cc_latest_execution "$cc_dt_root" "$cc_dt_plan") || { cc_fail DELIVERY_NO_EXECUTION; return 1; }
	[ -n "$cc_dt_exec" ] || { cc_fail DELIVERY_NO_EXECUTION; return 1; }
	cc_dt_edir=$(cc_execution_dir "$cc_dt_root" "$cc_dt_plan" "$cc_dt_exec")
	cc_emit execution_id "$cc_dt_exec"
	for cc_dt_rf in "$cc_dt_edir"/repositories/*.yaml; do
		[ -f "$cc_dt_rf" ] || continue
		cc_dt_id=$(cc_scalar "$cc_dt_rf" repository)
		cc_dt_src=$(cc_scalar "$cc_dt_rf" branch)
		cc_dt_tgt=$(cc_record_base_branch "$cc_dt_rf") || cc_dt_tgt=""
		cc_dt_abs=$(cc_repo_resolve "$cc_dt_root" "$cc_dt_id" 2>/dev/null | sed -n 's/^path: //p')
		cc_dt_present=unknown
		if [ -n "$cc_dt_abs" ]; then
			if git -C "$cc_dt_abs" show-ref --verify --quiet "refs/heads/$cc_dt_src"; then
				cc_dt_present=true
			else
				cc_dt_present=false
			fi
		fi
		cc_emit "repository" "$cc_dt_id"
		cc_emit "  source_branch" "$cc_dt_src"
		cc_emit "  target_branch" "$cc_dt_tgt"
		cc_emit "  source_present" "$cc_dt_present"
	done
	return 0
}

# cc_delivery_drift ROOT PLAN -> read-only report: for each repository of the
# latest execution, whether the recorded base has drifted from the current base
# tip (a sibling merged and advanced the base). Drift = the current base tip
# is not already reachable from the execution branch (INV-DELIVER-01 drift guard).
cc_delivery_drift() {
	cc_dd_root="$1"; cc_dd_plan="$2"
	cc_dd_exec=$(cc_latest_execution "$cc_dd_root" "$cc_dd_plan") || { cc_fail DELIVERY_NO_EXECUTION; return 1; }
	[ -n "$cc_dd_exec" ] || { cc_fail DELIVERY_NO_EXECUTION; return 1; }
	cc_dd_edir=$(cc_execution_dir "$cc_dd_root" "$cc_dd_plan" "$cc_dd_exec")
	cc_emit execution_id "$cc_dd_exec"
	cc_dd_any=false
	for cc_dd_rf in "$cc_dd_edir"/repositories/*.yaml; do
		[ -f "$cc_dd_rf" ] || continue
		cc_dd_id=$(cc_scalar "$cc_dd_rf" repository)
		cc_dd_base_branch=$(cc_record_base_branch "$cc_dd_rf") || cc_dd_base_branch=""
		cc_dd_base=$(cc_scalar "$cc_dd_rf" base_commit)
		cc_dd_latest=$(cc_scalar "$cc_dd_rf" latest_commit)
		cc_dd_abs=$(cc_repo_resolve "$cc_dd_root" "$cc_dd_id" 2>/dev/null | sed -n 's/^path: //p')
		cc_dd_tip=$(git -C "$cc_dd_abs" rev-parse --verify "refs/heads/$cc_dd_base_branch" 2>/dev/null) || cc_dd_tip=""
		cc_dd_drift=unknown
		if [ -n "$cc_dd_abs" ] && [ -n "$cc_dd_tip" ]; then
			if git -C "$cc_dd_abs" merge-base --is-ancestor "$cc_dd_tip" "$cc_dd_latest" 2>/dev/null; then
				cc_dd_drift=false
			else
				cc_dd_drift=true; cc_dd_any=true
			fi
		fi
		cc_emit "repository" "$cc_dd_id"
		cc_emit "  base_tip" "${cc_dd_tip:-unknown}"
		cc_emit "  recorded_base" "$cc_dd_base"
		cc_emit "  drifted" "$cc_dd_drift"
	done
	cc_emit drift_detected "$cc_dd_any"
	return 0
}

# cc_delivery_rebase ROOT PLAN -> rebase every drifted repository's execution
# branch onto the current base tip, update base_commit/latest_commit, and flag
# a required re-verification. A rebase conflict is DELIVERY_REBASE_CONFLICT (the
# runtime performs no delivery merge). Preserves work on abort.
cc_delivery_rebase() {
	cc_dr_root="$1"; cc_dr_plan="$2"
	cc_dr_exec=$(cc_latest_execution "$cc_dr_root" "$cc_dr_plan") || { cc_fail DELIVERY_NO_EXECUTION; return 1; }
	[ -n "$cc_dr_exec" ] || { cc_fail DELIVERY_NO_EXECUTION; return 1; }
	cc_dr_edir=$(cc_execution_dir "$cc_dr_root" "$cc_dr_plan" "$cc_dr_exec")
	cc_emit execution_id "$cc_dr_exec"
	cc_dr_rebased=0
	for cc_dr_rf in "$cc_dr_edir"/repositories/*.yaml; do
		[ -f "$cc_dr_rf" ] || continue
		cc_dr_id=$(cc_scalar "$cc_dr_rf" repository)
		cc_dr_base_branch=$(cc_record_base_branch "$cc_dr_rf") || cc_dr_base_branch=""
		cc_dr_wt=$(cc_scalar "$cc_dr_rf" worktree)
		cc_dr_latest=$(cc_scalar "$cc_dr_rf" latest_commit)
		cc_dr_abs=$(cc_repo_resolve "$cc_dr_root" "$cc_dr_id" 2>/dev/null | sed -n 's/^path: //p')
		[ -n "$cc_dr_abs" ] || continue
		cc_dr_tip=$(git -C "$cc_dr_abs" rev-parse --verify "refs/heads/$cc_dr_base_branch" 2>/dev/null) || continue
		# already contained: nothing to rebase for this repository
		if git -C "$cc_dr_abs" merge-base --is-ancestor "$cc_dr_tip" "$cc_dr_latest" 2>/dev/null; then
			cc_emit "repository" "$cc_dr_id"
			cc_emit "  rebased" false
			continue
		fi
		[ -d "$cc_dr_wt" ] || { cc_fail DELIVERY_WORKTREE_MISSING "$cc_dr_id"; return 1; }
		if ! git -C "$cc_dr_wt" rebase "$cc_dr_tip" >/dev/null 2>&1; then
			git -C "$cc_dr_wt" rebase --abort >/dev/null 2>&1 || :
			cc_fail DELIVERY_REBASE_CONFLICT "$cc_dr_id"; return 1
		fi
		cc_dr_new=$(git -C "$cc_dr_wt" rev-parse HEAD)
		awk -v b="$cc_dr_tip" -v l="$cc_dr_new" '
			/^base_commit:[[:space:]]/{print "base_commit: " b; next}
			/^latest_commit:[[:space:]]/{print "latest_commit: " l; next}
			{print}
		' "$cc_dr_rf" | cc_atomic_write "$cc_dr_rf"
		cc_dr_rebased=$((cc_dr_rebased + 1))
		cc_emit "repository" "$cc_dr_id"
		cc_emit "  rebased" true
		cc_emit "  new_base" "$cc_dr_tip"
		cc_emit "  new_latest" "$cc_dr_new"
	done
	if [ "$cc_dr_rebased" -gt 0 ]; then
		# a rebased base must be re-verified before its pull request opens
		cc_exec_set "$cc_dr_edir" status verifying || return 1
		cc_emit reverify_required true
	else
		cc_emit reverify_required false
	fi
	return 0
}

# ---------------------------------------------------------------------------
# Recovery
# ---------------------------------------------------------------------------

# cc_recovery_inspect EXEC_DIR [OWNER] -> report status, ownership, and resume
# eligibility. Resume is eligible only when the execution is mid-flight, its
# snapshot is intact, every worktree still exists, and (when OWNER is given) the
# requesting owner matches the recorded owner. Interrupted state is preserved.
cc_recovery_inspect() {
	cc_ri_dir="$1"; cc_ri_req="${2:-}"
	[ -f "$cc_ri_dir/execution.yaml" ] || { cc_fail RECOVERY_NO_EXECUTION; return 1; }
	cc_ri_st=$(cc_execution_status "$cc_ri_dir")
	cc_ri_owner=$(cc_scalar "$cc_ri_dir/execution.yaml" "owner")
	cc_emit status "$cc_ri_st"
	cc_emit owner "$cc_ri_owner"
	cc_emit worker_failures "$(cc_scalar "$cc_ri_dir/execution.yaml" worker_failures)"
	# preservation evidence
	cc_ri_snap=false; [ -f "$cc_ri_dir/snapshot/plan.yaml" ] && cc_ri_snap=true
	cc_emit snapshot_present "$cc_ri_snap"
	cc_ri_wt=intact
	for cc_ri_rf in "$cc_ri_dir"/repositories/*.yaml; do
		[ -f "$cc_ri_rf" ] || continue
		cc_ri_w=$(cc_scalar "$cc_ri_rf" worktree)
		[ -d "$cc_ri_w" ] || cc_ri_wt=missing
	done
	cc_emit worktrees "$cc_ri_wt"
	# resume eligibility
	cc_ri_elig=false; cc_ri_reason=""
	case "$cc_ri_st" in
		running|verifying|repairing)
			if [ "$cc_ri_snap" = true ] && [ "$cc_ri_wt" = intact ]; then
				if [ -z "$cc_ri_req" ] || [ "$cc_ri_req" = "$cc_ri_owner" ]; then
					cc_ri_elig=true
				else
					cc_ri_reason=OWNER_MISMATCH
				fi
			else
				cc_ri_reason=STATE_INCOMPLETE
			fi ;;
		verified|failed|blocked) cc_ri_reason=TERMINAL_STATE ;;
		*) cc_ri_reason=UNKNOWN_STATE ;;
	esac
	cc_emit resume_eligible "$cc_ri_elig"
	[ -n "$cc_ri_reason" ] && cc_emit resume_reason "$cc_ri_reason" || :
	# resumable retained for backward-compatible callers
	case "$cc_ri_st" in
		running|verifying|repairing) cc_emit resumable true ;;
		*) cc_emit resumable false ;;
	esac
	return 0
}

# ---------------------------------------------------------------------------
# Thin CLI dispatch (optional; host adapters may source functions instead)
# ---------------------------------------------------------------------------

cc_main() {
	cc_cmd="${1:-}"
	shift 2>/dev/null || true
	case "$cc_cmd" in
		workspace-validate)      cc_workspace_validate "$@" ;;
		workspace-init)          cc_workspace_init "$@" ;;
		repository-register)     cc_repository_register "$@" ;;
		repository-binding-migrate) cc_repository_binding_migrate "$@" ;;
		repository-resolve)      cc_repo_resolve "$@" ;;
		repository-preflight)    cc_repository_preflight "$@" ;;
		delivery-targets)        cc_delivery_targets "$@" ;;
		worktree-prepare)        cc_worktree_prepare "$@" ;;
		pair-begin)              cc_pair_begin "$@" ;;
		pair-inspect)            cc_pair_inspect "$@" ;;
		pair-close)              cc_pair_close "$@" ;;
		pair-delivery-targets)   cc_pair_delivery_targets "$@" ;;
		runtime-cleanup)         cc_runtime_cleanup "$@" ;;
		base-prepare)            cc_base_prepare "$@" ;;
		discover-repo-grounding) cc_discover_repo_grounding "$@" ;;
		harden-worktree)         cc_harden_worktree "$@" ;;
		grounding-directive)     cc_grounding_directive "$@" ;;
		worker-brief-assemble)   cc_worker_brief_assemble "$@" ;;
		brief-preflight)         cc_brief_preflight "$@" ;;
		lease-check)             cc_lease_check "$@" ;;
		lease-acquire)           cc_lease_acquire "$@" ;;
		lease-release)           cc_lease_release "$@" ;;
		plan-ready)              cc_plan_ready "$@" ;;
		run-stack-ready)         cc_run_stack_ready "$@" ;;
		delivery-drift)          cc_delivery_drift "$@" ;;
		delivery-rebase)         cc_delivery_rebase "$@" ;;
		plan-validate)           cc_plan_validate "$@" ;;
		plan-allocate-id)        cc_plan_allocate_id "$@" ;;
		plan-archive)            cc_plan_archive "$@" ;;
		plan-restore)            cc_plan_restore "$@" ;;
		intent-validate)         cc_intent_validate "$@" ;;
		intent-allocate-id)      cc_intent_allocate_id "$@" ;;
		intent-approve)          cc_intent_approve "$@" ;;
		intent-human-status)     cc_intent_human_status "$@" ;;
		intent-authorized)       cc_intent_authorized "$@" ;;
		intent-archive)          cc_intent_archive "$@" ;;
		intent-restore)          cc_intent_restore "$@" ;;
		intent-index-upsert)     cc_intent_index_upsert "$@" ;;
		intent-index-remove)     cc_intent_index_remove "$@" ;;
		plan-index-upsert)       cc_plan_index_upsert "$@" ;;
		plan-index-remove)       cc_plan_index_remove "$@" ;;
		execution-begin)         cc_execution_begin "$@" ;;
		attempt-begin)           cc_attempt_begin "$@" ;;
		worker-commit-record)    cc_worker_commit_record "$@" ;;
		worker-handoff-record)   cc_worker_handoff_record "$@" ;;
		verifier-prepare)        cc_verifier_prepare "$@" ;;
		verifier-result-record)  cc_verifier_result_record "$@" ;;
		attempt-evidence-record) cc_attempt_evidence_record "$@" ;;
		repair-allowed)          cc_repair_allowed "$@" ;;
		execution-status)        cc_execution_status "$@" ;;
		candidate-digest)        cc_candidate_digest "$@" ;;
		candidate-current)       cc_candidate_current "$@" ;;
		change-set-candidate)    cc_change_set_candidate "$@" ;;
		change-set-prepare)      cc_change_set_prepare "$@" ;;
		change-set-partition)    cc_change_set_partition "$@" ;;
		change-set-verifier-prepare) cc_change_set_verifier_prepare "$@" ;;
		change-set-verifier-record) cc_change_set_verifier_record "$@" ;;
		change-set-accept)       cc_change_set_accept "$@" ;;
		change-set-ready)        cc_change_set_ready "$@" ;;
		change-set-complete)     cc_change_set_complete "$@" ;;
		human-acceptance-record) cc_human_acceptance_record "$@" ;;
		human-acceptance-current) cc_human_acceptance_current "$@" ;;
		tier-classify)           cc_tier_classify "$@" ;;
		tier-lower-check)        cc_tier_lower_check "$@" ;;
		completion-ready)        cc_completion_ready "$@" ;;
		plan-complete)           cc_plan_complete "$@" ;;
		completion-infer)        cc_completion_infer "$@" ;;
		delivery-record)         cc_delivery_record "$@" ;;
		context-impact-record)   cc_context_impact_record "$@" ;;
		knowledge-debt-emit)     cc_knowledge_debt_emit "$@" ;;
		knowledge-debt)          cc_knowledge_debt "$@" ;;
		knowledge-debt-check)    cc_knowledge_debt_check "$@" ;;
		knowledge-reconciled)    cc_knowledge_reconciled "$@" ;;
		recovery-inspect)        cc_recovery_inspect "$@" ;;
		lock-acquire)            cc_lock_acquire "$@" ;;
		lock-release)            cc_lock_release "$@" ;;
		version)                 printf 'context-circuit-runtime %s\n' "$CC_RUNTIME_VERSION" ;;
		''|help|-h|--help)
			printf 'Context Circuit runtime %s\nUsage: engine.sh <command> [args]\n' "$CC_RUNTIME_VERSION" ;;
		*) cc_fail UNKNOWN_COMMAND "$cc_cmd" ;;
	esac
}

# Run as CLI only when invoked directly (not when sourced by tests/adapters).
case "$0" in
	*engine.sh) cc_main "$@" ;;
esac
