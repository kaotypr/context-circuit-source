#!/bin/sh
# Context Circuit v0.7.0 runtime engine.
#
# A small, host-neutral, deterministic runtime library for workspace, Git, and
# execution-state operations. It is sourced by host adapters and tests, and can
# also be invoked as a thin CLI: `sh engine.sh <command> [--flag value ...]`.
#
# This runtime OWNS: safe path/identifier checks, atomic writes and digests,
# workspace and repository-binding validation, anchor branch/commit validation,
# branch/worktree preparation, plan structure and approval-state validation,
# active plan-index maintenance, exact archive/restore moves, execution and
# attempt records, commit capture, verifier-result and read-only enforcement,
# the three-failure counter, one-worker locking, completion eligibility,
# implementation completion records, context-impact handoff references, light
# direct-collaboration session pointers, and recovery inspection.
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

CC_RUNTIME_VERSION="0.7.0"
CC_SCHEMA_VERSION="1"

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
		"$cc_wi_root/plans/archive" "$cc_wi_root/.runtime/executions" \
		"$cc_wi_root/.runtime/worktrees" "$cc_wi_root/.runtime/pairing" \
		"$cc_wi_root/.runtime/locks" \
		"$cc_wi_root/repositories" || return 1
	[ -f "$cc_wi_root/plans/INDEX.md" ] || cc_plan_index_init "$cc_wi_root"
	[ -f "$cc_wi_root/context/INDEX.md" ] || printf '# Context index\n\nNo accepted context units yet.\n' | cc_atomic_write "$cc_wi_root/context/INDEX.md"
	[ -f "$cc_wi_root/repositories.local.yaml" ] || printf 'schema_version: %s\nbindings: {}\n' "$CC_SCHEMA_VERSION" | cc_atomic_write "$cc_wi_root/repositories.local.yaml"
	cc_emit workspace_init ok
	return 0
}

cc_plan_index_init() {
	printf '# Active plans\n\n| Plan ID | Title | Status | Objective | Repositories | Path |\n| --- | --- | --- | --- | --- | --- |\n' \
		| cc_atomic_write "$1/plans/INDEX.md"
}

# ---------------------------------------------------------------------------
# Repository bindings, anchor branches, worktrees
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
			print; exit
		}
	' "$cc_bf_file"
}

# cc_repository_register ROOT ID PATH ANCHOR [CANONICAL_URL] [DEFAULT_BRANCH]
# Records portable logical identity in workspace.yaml and a host-local binding
# in repositories.local.yaml. It does not clone, init, or execute anything.
cc_repository_register() {
	cc_reg_root="$1"; cc_reg_id="$2"; cc_reg_path="$3"; cc_reg_anchor="$4"
	cc_reg_url="${5:-}"; cc_reg_default="${6:-}"
	cc_safe_id "$cc_reg_id" || { cc_fail REPOSITORY_ID_UNSAFE "$cc_reg_id"; return 1; }
	[ -n "$cc_reg_path" ] || { cc_fail REPOSITORY_PATH_MISSING; return 1; }
	[ -n "$cc_reg_anchor" ] || { cc_fail REPOSITORY_ANCHOR_UNSET "$cc_reg_id"; return 1; }
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
		printf 'schema_version: %s\nbindings:\n' "$CC_SCHEMA_VERSION" | cc_atomic_write "$cc_reg_local"
	fi
	if cc_binding_field "$cc_reg_root" "$cc_reg_id" "path" >/dev/null 2>&1 && [ -n "$(cc_binding_field "$cc_reg_root" "$cc_reg_id" path)" ]; then
		cc_fail REPOSITORY_BINDING_EXISTS "$cc_reg_id"; return 1
	fi
	printf '  %s:\n    path: %s\n    anchor_branch: %s\n' "$cc_reg_id" "$cc_reg_path" "$cc_reg_anchor" >>"$cc_reg_local"
	cc_emit repository "$cc_reg_id"
	cc_emit registered ok
	return 0
}

# cc_repo_resolve ROOT REPO -> validate binding; emit path/anchor; REPOSITORY_* codes
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
	cc_rr_anchor=$(cc_binding_field "$cc_rr_root" "$cc_rr_id" "anchor_branch") || cc_rr_anchor=""
	[ -n "$cc_rr_anchor" ] || { cc_fail REPOSITORY_ANCHOR_UNSET "$cc_rr_id"; return 1; }
	cc_emit repository "$cc_rr_id"
	cc_emit path "$cc_rr_abs"
	cc_emit anchor_branch "$cc_rr_anchor"
	return 0
}

# cc_repo_anchor_commit ROOT REPO -> print the anchor branch tip commit
cc_repo_anchor_commit() {
	cc_rac_abs=$(cc_repo_resolve "$1" "$2" | sed -n 's/^path: //p')
	cc_rac_anchor=$(cc_binding_field "$1" "$2" "anchor_branch")
	git -C "$cc_rac_abs" rev-parse --verify "refs/heads/$cc_rac_anchor" 2>/dev/null \
		|| { cc_fail ANCHOR_BRANCH_MISSING "$2:$cc_rac_anchor"; return 1; }
}

# cc_repo_clean ROOT REPO -> ok if the anchor checkout has no uncommitted changes
cc_repo_clean() {
	cc_rc_abs=$(cc_repo_resolve "$1" "$2" | sed -n 's/^path: //p') || return 1
	if [ -n "$(git -C "$cc_rc_abs" status --porcelain 2>/dev/null)" ]; then
		cc_fail REPOSITORY_ANCHOR_DIRTY "$2"; return 1
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
		cc_repo_anchor_commit "$cc_rp_root" "$cc_rp_id" >/dev/null || { cc_fail ANCHOR_PREFLIGHT_FAILED "$cc_rp_id"; return 1; }
		cc_repo_clean "$cc_rp_root" "$cc_rp_id" || { cc_fail REPOSITORY_ANCHOR_DIRTY "$cc_rp_id"; return 1; }
		cc_rp_ok=$((cc_rp_ok + 1))
	done
	[ "$cc_rp_ok" -gt 0 ] || { cc_fail PLAN_NO_AFFECTED_REPOSITORIES; return 1; }
	cc_emit repositories_ready "$cc_rp_ok"
	return 0
}

# cc_worktree_prepare ROOT PLAN_ID REPO -> create branch + worktree from anchor tip
cc_worktree_prepare() {
	cc_wp_root="$1"; cc_wp_plan="$2"; cc_wp_id="$3"
	cc_wp_abs=$(cc_repo_resolve "$cc_wp_root" "$cc_wp_id" | sed -n 's/^path: //p') || return 1
	cc_wp_anchor=$(cc_binding_field "$cc_wp_root" "$cc_wp_id" "anchor_branch")
	cc_wp_base=$(git -C "$cc_wp_abs" rev-parse --verify "refs/heads/$cc_wp_anchor" 2>/dev/null) \
		|| { cc_fail ANCHOR_BRANCH_MISSING "$cc_wp_id"; return 1; }
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
	[ "$cc_pv_schema" = "1" ] || { cc_fail PAIR_SCHEMA_UNSUPPORTED "$cc_pv_schema"; return 1; }
	cc_pv_repo=$(cc_scalar "$cc_pv_file" repo) || cc_pv_repo=""
	cc_pv_wt=$(cc_scalar "$cc_pv_file" worktree) || cc_pv_wt=""
	cc_pv_branch=$(cc_scalar "$cc_pv_file" branch) || cc_pv_branch=""
	cc_pv_base=$(cc_scalar "$cc_pv_file" base) || cc_pv_base=""
	cc_safe_id "$cc_pv_repo" || { cc_fail PAIR_REPOSITORY_INVALID "$cc_pv_repo"; return 1; }
	cc_pv_expected_wt="$cc_pv_root/.runtime/worktrees/cc-pair/$cc_pv_session/$cc_pv_repo"
	cc_pv_expected_branch="cc-pair/$cc_pv_session"
	[ "$cc_pv_wt" = "$cc_pv_expected_wt" ] || { cc_fail PAIR_WORKTREE_MISMATCH "$cc_pv_session"; return 1; }
	[ "$cc_pv_branch" = "$cc_pv_expected_branch" ] || { cc_fail PAIR_BRANCH_MISMATCH "$cc_pv_session"; return 1; }
	[ -n "$cc_pv_base" ] || { cc_fail PAIR_BASE_MISSING "$cc_pv_session"; return 1; }
	return 0
}

# cc_pair_begin ROOT REPO SESSION [BASE] -> create a fresh pairing branch and
# worktree from BASE (the connected repository's anchor tip by default), then
# atomically write the only resumable session state.
cc_pair_begin() {
	cc_pb_root="$1"; cc_pb_repo="$2"; cc_pb_session="$3"; cc_pb_requested_base="${4:-}"
	cc_workspace_validate "$cc_pb_root" >/dev/null || return 1
	cc_safe_id "$cc_pb_repo" || { cc_fail PAIR_REPOSITORY_INVALID "$cc_pb_repo"; return 1; }
	cc_safe_slug "$cc_pb_session" || { cc_fail PAIR_SESSION_INVALID "$cc_pb_session"; return 1; }
	cc_pb_dir="$cc_pb_root/.runtime/pairing/$cc_pb_session"
	cc_pb_pointer="$cc_pb_dir/pointer.yaml"
	cc_pb_closed="$cc_pb_dir/closed.yaml"
	[ ! -e "$cc_pb_pointer" ] && [ ! -e "$cc_pb_closed" ] \
		|| { cc_fail PAIR_SESSION_EXISTS "$cc_pb_session"; return 1; }
	cc_pb_abs=$(cc_repo_resolve "$cc_pb_root" "$cc_pb_repo" | sed -n 's/^path: //p') || return 1
	cc_pb_anchor=$(cc_binding_field "$cc_pb_root" "$cc_pb_repo" anchor_branch) || cc_pb_anchor=""
	if [ -n "$cc_pb_requested_base" ]; then
		case "$cc_pb_requested_base" in
			-*|*' '*|*..*) cc_fail PAIR_BASE_INVALID "$cc_pb_requested_base"; return 1 ;;
		esac
		cc_pb_base=$(git -C "$cc_pb_abs" rev-parse --verify "$cc_pb_requested_base^{commit}" 2>/dev/null) \
			|| { cc_fail PAIR_BASE_INVALID "$cc_pb_requested_base"; return 1; }
	else
		cc_pb_base=$(git -C "$cc_pb_abs" rev-parse --verify "refs/heads/$cc_pb_anchor" 2>/dev/null) \
			|| { cc_fail ANCHOR_BRANCH_MISSING "$cc_pb_repo:$cc_pb_anchor"; return 1; }
	fi
	cc_pb_branch="cc-pair/$cc_pb_session"
	cc_pb_wt="$cc_pb_root/.runtime/worktrees/cc-pair/$cc_pb_session/$cc_pb_repo"
	git -C "$cc_pb_abs" show-ref --verify --quiet "refs/heads/$cc_pb_branch" \
		&& { cc_fail PAIR_BRANCH_EXISTS "$cc_pb_branch"; return 1; }
	[ ! -e "$cc_pb_wt" ] || { cc_fail PAIR_WORKTREE_EXISTS "$cc_pb_wt"; return 1; }
	mkdir -p "$(dirname -- "$cc_pb_wt")" "$cc_pb_dir" || { cc_fail PAIR_STATE_CREATE_FAILED "$cc_pb_session"; return 1; }
	git -C "$cc_pb_abs" worktree add -b "$cc_pb_branch" "$cc_pb_wt" "$cc_pb_base" >/dev/null 2>&1 \
		|| { cc_fail PAIR_WORKTREE_CREATE_FAILED "$cc_pb_repo"; return 1; }
	{
		printf 'schema_version: 1\n'
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
# current connected anchor target. Drift blocks; this function never rebases,
# pushes, merges, or opens a pull request.
cc_pair_delivery_targets() {
	cc_pd_root="$1"; cc_pd_session="$2"
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
	cc_pd_anchor=$(cc_binding_field "$cc_pd_root" "$cc_pd_repo" anchor_branch) || cc_pd_anchor=""
	cc_pd_tip=$(git -C "$cc_pd_abs" rev-parse --verify "refs/heads/$cc_pd_branch" 2>/dev/null) \
		|| { cc_fail PAIR_BRANCH_MISSING "$cc_pd_branch"; return 1; }
	cc_pd_anchor_tip=$(git -C "$cc_pd_abs" rev-parse --verify "refs/heads/$cc_pd_anchor" 2>/dev/null) \
		|| { cc_fail ANCHOR_BRANCH_MISSING "$cc_pd_repo:$cc_pd_anchor"; return 1; }
	cc_pd_drift=true
	if git -C "$cc_pd_abs" merge-base --is-ancestor "$cc_pd_anchor_tip" "$cc_pd_tip" 2>/dev/null; then cc_pd_drift=false; fi
	cc_emit session "$cc_pd_session"
	cc_emit repository "$cc_pd_repo"
	cc_emit source_branch "$cc_pd_branch"
	cc_emit target_branch "$cc_pd_anchor"
	cc_emit base_commit "$cc_pd_base"
	cc_emit branch_tip "$cc_pd_tip"
	cc_emit anchor_tip "$cc_pd_anchor_tip"
	cc_emit drift_detected "$cc_pd_drift"
	cc_emit result_label human-supervised
	if [ "$cc_pd_drift" = true ]; then cc_fail PAIR_ANCHOR_DRIFT "$cc_pd_session"; return 1; fi
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
# plan. Selects the base (anchor tip / single predecessor branch / runtime-authored
# integration merge), keeps the base ref at refs/cc-base/<plan>/<repo>, and detects
# a stale base (predecessor repaired) to rebuild it. Emits the same keys as
# cc_worktree_prepare plus based_on and base_kind. BASE_UNBUILDABLE -> blocked.
cc_base_prepare() {
	cc_bp_root="$1"; cc_bp_plan="$2"; cc_bp_repo="$3"
	cc_bp_abs=$(cc_repo_resolve "$cc_bp_root" "$cc_bp_repo" | sed -n 's/^path: //p') || return 1
	cc_bp_anchor=$(cc_binding_field "$cc_bp_root" "$cc_bp_repo" "anchor_branch")
	cc_bp_anchortip=$(git -C "$cc_bp_abs" rev-parse --verify "refs/heads/$cc_bp_anchor" 2>/dev/null) \
		|| { cc_fail ANCHOR_BRANCH_MISSING "$cc_bp_repo"; return 1; }
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
		cc_bp_kind=anchor; cc_bp_start="$cc_bp_anchortip"
		git -C "$cc_bp_abs" worktree add -b "$cc_bp_branch" "$cc_bp_tree" "$cc_bp_start" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_bp_repo"; return 1; }
		cc_bp_base="$cc_bp_anchortip"
	elif [ "$cc_bp_n" -eq 1 ]; then
		cc_bp_kind=stack; cc_bp_start="$cc_bp_tips"
		git -C "$cc_bp_abs" worktree add -b "$cc_bp_branch" "$cc_bp_tree" "$cc_bp_start" >/dev/null 2>&1 \
			|| { cc_fail WORKTREE_CREATE_FAILED "$cc_bp_repo"; return 1; }
		cc_bp_base="$cc_bp_start"
	else
		cc_bp_kind=integration
		git -C "$cc_bp_abs" worktree add -b "$cc_bp_branch" "$cc_bp_tree" "$cc_bp_anchortip" >/dev/null 2>&1 \
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
			git -C "$cc_bp_tree" merge --abort >/dev/null 2>&1 || :
			git -C "$cc_bp_abs" worktree remove --force "$cc_bp_tree" >/dev/null 2>&1 || rm -rf "$cc_bp_tree"
			git -C "$cc_bp_abs" branch -D "$cc_bp_branch" >/dev/null 2>&1 || :
			git -C "$cc_bp_abs" worktree prune >/dev/null 2>&1 || :
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
	printf 'schema_version: 1\n'
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
		draft|approved|done) : ;;
		*) cc_fail PLAN_STATUS_INVALID "$cc_pv_status"; return 1 ;;
	esac
	# schema version and inter-plan dependencies (INV-PLAN-05)
	cc_pv_schema=$(cc_scalar "$cc_pv_dir/plan.yaml" "schema_version") || cc_pv_schema=""
	case "$cc_pv_schema" in
		1|2) : ;;
		*) cc_fail PLAN_SCHEMA_UNSUPPORTED "$cc_pv_schema"; return 1 ;;
	esac
	# optional complexity hint (additive; absent by default). A coordinator
	# tiering hint only (INV-HOST-01) — never a gate; the runtime stays model-blind.
	cc_pv_cx=$(cc_scalar "$cc_pv_dir/plan.yaml" "complexity") || cc_pv_cx=""
	case "$cc_pv_cx" in
		''|standard|high) : ;;
		*) cc_fail PLAN_COMPLEXITY_INVALID "$cc_pv_cx"; return 1 ;;
	esac
	cc_pv_deps=$(cc_plan_dependencies "$cc_pv_dir/plan.yaml")
	if [ -n "$cc_pv_deps" ]; then
		[ "$cc_pv_schema" = "2" ] || { cc_fail PLAN_DEPS_REQUIRE_SCHEMA_2 "$cc_pv_id"; return 1; }
		cc_pv_plansdir=$(dirname -- "$cc_pv_dir")
		for cc_pv_dep in $cc_pv_deps; do
			[ "$cc_pv_dep" != "$cc_pv_id" ] || { cc_fail PLAN_DEP_SELF "$cc_pv_dep"; return 1; }
			[ -f "$cc_pv_plansdir/$cc_pv_dep/plan.yaml" ] || { cc_fail PLAN_DEP_UNKNOWN "$cc_pv_dep"; return 1; }
		done
		if cc_plan_dep_closure "$cc_pv_plansdir" "$cc_pv_id" | grep -Fxq "$cc_pv_id"; then
			cc_fail PLAN_DEP_CYCLE "$cc_pv_id"; return 1
		fi
	fi
	# every task maps to at least one declared repository; deps reference tasks
	cc_pv_repos=$(cc_plan_repositories "$cc_pv_dir/plan.yaml")
	cc_pv_tasks=$(cc_task_ids "$cc_pv_dir/plan.yaml")
	[ -n "$cc_pv_tasks" ] || { cc_fail PLAN_NO_TASKS; return 1; }
	for cc_pv_t in $cc_pv_tasks; do
		cc_pv_treps=$(cc_task_list "$cc_pv_dir/plan.yaml" "$cc_pv_t" "repositories")
		[ -n "$cc_pv_treps" ] || { cc_fail TASK_NO_REPOSITORY "$cc_pv_t"; return 1; }
		for cc_pv_tr in $cc_pv_treps; do
			printf '%s\n' "$cc_pv_repos" | grep -Fxq "$cc_pv_tr" \
				|| { cc_fail TASK_UNDECLARED_REPOSITORY "$cc_pv_t:$cc_pv_tr"; return 1; }
		done
		for cc_pv_dep in $(cc_task_list "$cc_pv_dir/plan.yaml" "$cc_pv_t" "depends_on"); do
			printf '%s\n' "$cc_pv_tasks" | grep -Fxq "$cc_pv_dep" \
				|| { cc_fail TASK_UNKNOWN_DEPENDENCY "$cc_pv_t:$cc_pv_dep"; return 1; }
		done
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

# cc_plan_approve ROOT PLAN -> draft->approved after readiness checks
cc_plan_approve() {
	cc_ap_root="$1"; cc_ap_plan="$2"
	cc_ap_dir="$cc_ap_root/plans/$cc_ap_plan"
	cc_plan_validate "$cc_ap_dir" >/dev/null || { cc_fail APPROVAL_PLAN_INVALID "$cc_ap_plan"; return 1; }
	cc_ap_status=$(cc_scalar "$cc_ap_dir/plan.yaml" "status")
	[ "$cc_ap_status" = "draft" ] || { cc_fail APPROVAL_NOT_DRAFT "$cc_ap_status"; return 1; }
	cc_plan_set_status "$cc_ap_dir/plan.yaml" "approved" || { cc_fail APPROVAL_WRITE_FAILED; return 1; }
	cc_plan_index_upsert "$cc_ap_root" "$cc_ap_plan" >/dev/null
	cc_emit plan "$cc_ap_plan"
	cc_emit status approved
	return 0
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
			| cc_atomic_write "$cc_la_dir/owner.yaml"
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
	printf 'schema_version: 1\nplan: %s\nrepository: %s\nregions: [%s]\nacquired_at: %s\nreleased_at:\n' \
		"$cc_la2_plan" "$cc_la2_repo" "$cc_la2_inline" "$(cc_now)" \
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
	cc_eb_status=$(cc_scalar "$cc_eb_dir/plan.yaml" "status")
	[ "$cc_eb_status" = "approved" ] || { cc_fail EXECUTION_NOT_APPROVED "$cc_eb_status"; return 1; }
	cc_repository_preflight "$cc_eb_root" "$cc_eb_dir" >/dev/null || { cc_fail EXECUTION_PREFLIGHT_FAILED; return 1; }
	cc_lock_acquire "$cc_eb_root" "$cc_eb_plan" "$cc_eb_owner" >/dev/null || { cc_fail EXECUTION_LOCK_FAILED; return 1; }
	cc_eb_exec=$(cc_execution_next_id "$cc_eb_root" "$cc_eb_plan")
	cc_eb_edir=$(cc_execution_dir "$cc_eb_root" "$cc_eb_plan" "$cc_eb_exec")
	mkdir -p "$cc_eb_edir/attempts" "$cc_eb_edir/repositories" "$cc_eb_edir/snapshot"
	# immutable plan snapshot
	cp "$cc_eb_dir/plan.yaml" "$cc_eb_edir/snapshot/plan.yaml"
	cp "$cc_eb_dir/PLAN.md" "$cc_eb_edir/snapshot/PLAN.md"
	[ -d "$cc_eb_dir/tasks" ] && cp -R "$cc_eb_dir/tasks" "$cc_eb_edir/snapshot/tasks"
	cc_eb_rev=$(cc_digest "$cc_eb_dir/plan.yaml")
	# per-repository branch + worktree. A plan with same-repo predecessors is
	# base-aware (INV-CONCURRENCY-02): its base is the predecessor branch (stack)
	# or a runtime-authored integration merge. A plan with no dependency keeps the
	# v0.5 anchor-tip worktree unchanged. Leases are NOT acquired here; the
	# run-stack loop manages them (v0.5 fixtures run overlapping-path plans).
	for cc_eb_id in $(cc_plan_affected_repositories "$cc_eb_dir/plan.yaml"); do
		cc_eb_based=""
		if cc_plan_has_same_repo_pred "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id"; then
			if ! cc_eb_out=$(cc_base_prepare "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id"); then
				# a base that cannot be built cleanly is a blocked execution, not a
				# worker failure (INV-CONCURRENCY-02); preserve evidence, no worker runs.
				printf 'schema_version: %s\nexecution_id: %s\nplan: %s\nplan_revision: %s\nowner: %s\nstatus: blocked\nblocked_reason: BASE_UNBUILDABLE\nworker_failures: 0\ncurrent_attempt: 0\ncreated_at: %s\nupdated_at: %s\n' \
					"$CC_SCHEMA_VERSION" "$cc_eb_exec" "$cc_eb_plan" "$cc_eb_rev" "$cc_eb_owner" "$(cc_now)" "$(cc_now)" \
					| cc_atomic_write "$cc_eb_edir/execution.yaml"
				cc_emit execution_id "$cc_eb_exec"
				cc_emit status blocked
				cc_emit blocked_reason BASE_UNBUILDABLE
				cc_fail EXECUTION_BASE_UNBUILDABLE "$cc_eb_id"; return 1
			fi
			cc_eb_based=$(printf '%s' "$cc_eb_out" | sed -n 's/^based_on: //p')
		else
			cc_eb_out=$(cc_worktree_prepare "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id") \
				|| { cc_fail EXECUTION_WORKTREE_FAILED "$cc_eb_id"; return 1; }
		fi
		cc_eb_wt=$(printf '%s' "$cc_eb_out" | sed -n 's/^worktree: //p')
		cc_eb_br=$(printf '%s' "$cc_eb_out" | sed -n 's/^branch: //p')
		cc_eb_bc=$(printf '%s' "$cc_eb_out" | sed -n 's/^base_commit: //p')
		cc_eb_an=$(cc_binding_field "$cc_eb_root" "$cc_eb_id" "anchor_branch")
		cc_eb_paths=$(cc_plan_repo_paths "$cc_eb_dir/plan.yaml" "$cc_eb_id" | tr '\n' ',' | sed 's/,$//; s/,/, /g')
		{
			printf 'repository: %s\nworktree: %s\nbranch: %s\nanchor_branch: %s\nbase_commit: %s\nlatest_commit: %s\nallowed_paths: [%s]\n' \
				"$cc_eb_id" "$cc_eb_wt" "$cc_eb_br" "$cc_eb_an" "$cc_eb_bc" "$cc_eb_bc" "$cc_eb_paths"
			[ -n "$cc_eb_based" ] && printf 'based_on: %s\n' "$cc_eb_based" || :
		} | cc_atomic_write "$cc_eb_edir/repositories/$cc_eb_id.yaml"
		# repository grounding (INV-GROUND-01): discover the target repo's own agent
		# guidance from the prepared worktree and record it as execution evidence.
		mkdir -p "$cc_eb_edir/grounding"
		cc_discover_repo_grounding "$cc_eb_wt" "$cc_eb_id" | cc_atomic_write "$cc_eb_edir/grounding/$cc_eb_id.yaml"
	done
	printf 'schema_version: %s\nexecution_id: %s\nplan: %s\nplan_revision: %s\nowner: %s\nstatus: running\nworker_failures: 0\ncurrent_attempt: 0\ncreated_at: %s\nupdated_at: %s\n' \
		"$CC_SCHEMA_VERSION" "$cc_eb_exec" "$cc_eb_plan" "$cc_eb_rev" "$cc_eb_owner" "$(cc_now)" "$(cc_now)" \
		| cc_atomic_write "$cc_eb_edir/execution.yaml"
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
	' "$cc_es_dir/execution.yaml" | cc_atomic_write "$cc_es_dir/execution.yaml"
}

# cc_attempt_begin EXEC_DIR -> increment current_attempt, create attempt dir
cc_attempt_begin() {
	cc_ab_dir="$1"
	[ -f "$cc_ab_dir/execution.yaml" ] || { cc_fail EXECUTION_RECORD_MISSING; return 1; }
	cc_ab_cur=$(cc_scalar "$cc_ab_dir/execution.yaml" "current_attempt")
	cc_ab_next=$((cc_ab_cur + 1))
	cc_ab_pad=$(printf '%03d' "$cc_ab_next")
	mkdir -p "$cc_ab_dir/attempts/$cc_ab_pad"
	cc_exec_set "$cc_ab_dir" current_attempt "$cc_ab_next"
	cc_exec_set "$cc_ab_dir" status running
	printf 'attempt: %s\nstatus: started\nstarted_at: %s\n' "$cc_ab_next" "$(cc_now)" \
		| cc_atomic_write "$cc_ab_dir/attempts/$cc_ab_pad/worker.yaml"
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
		| cc_atomic_write "$cc_wc_rf"
	cc_wc_att=$(cc_scalar "$cc_wc_dir/execution.yaml" "current_attempt")
	cc_wc_pad=$(printf '%03d' "$cc_wc_att")
	cc_wc_wf="$cc_wc_dir/attempts/$cc_wc_pad/worker.yaml"
	{
		[ -f "$cc_wc_wf" ] && cat "$cc_wc_wf" || :
		printf 'commit repository=%s revision=%s kind=%s at=%s\n' \
			"$cc_wc_repo" "$cc_wc_cur" "$cc_wc_kind" "$(cc_now)"
	} | cc_atomic_write "$cc_wc_wf"
	cc_exec_set "$cc_wc_dir" status verifying
	cc_emit repository "$cc_wc_repo"
	cc_emit commit "$cc_wc_cur"
	cc_emit kind "$cc_wc_kind"
	return 0
}

# cc_worker_handoff_record EXEC_DIR FILE -> store the worker handoff
cc_worker_handoff_record() {
	[ -f "$2" ] || { cc_fail HANDOFF_FILE_MISSING; return 1; }
	cp "$2" "$1/handoff.md"
	cc_emit handoff recorded
	return 0
}

# cc_verifier_prepare EXEC_DIR -> validate latest revisions and read-only scope
cc_verifier_prepare() {
	cc_vp_dir="$1"
	for cc_vp_rf in "$cc_vp_dir"/repositories/*.yaml; do
		[ -f "$cc_vp_rf" ] || continue
		cc_vp_latest=$(cc_scalar "$cc_vp_rf" "latest_commit")
		cc_vp_base=$(cc_scalar "$cc_vp_rf" "base_commit")
		[ "$cc_vp_latest" != "$cc_vp_base" ] || { cc_fail VERIFIER_NO_WORKER_COMMIT "$(basename "$cc_vp_rf" .yaml)"; return 1; }
	done
	cc_emit verifier_scope read-only
	return 0
}

# cc_verifier_result_record EXEC_DIR ATTEMPT OUTCOME [--wrote-products] -> record + count
cc_verifier_result_record() {
	cc_vr_dir="$1"; cc_vr_att="$2"; cc_vr_out="$3"; cc_vr_flag="${4:-}"
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
		[ -f "$cc_vr_rf" ] || continue
		cc_vr_wt=$(cc_scalar "$cc_vr_rf" "worktree")
		cc_vr_rec=$(cc_scalar "$cc_vr_rf" "latest_commit")
		[ -d "$cc_vr_wt" ] || continue
		cc_vr_head=$(git -C "$cc_vr_wt" rev-parse HEAD 2>/dev/null)
		if [ -n "$cc_vr_head" ] && [ "$cc_vr_head" != "$cc_vr_rec" ]; then
			cc_fail VERIFIER_MODIFIED_PRODUCT "$(basename "$cc_vr_rf" .yaml)"; return 1
		fi
	done
	cc_vr_pad=$(printf '%03d' "$cc_vr_att")
	mkdir -p "$cc_vr_dir/attempts/$cc_vr_pad"
	printf 'attempt: %s\noutcome: %s\nread_only: true\nchecked_at: %s\n' \
		"$cc_vr_att" "$cc_vr_out" "$(cc_now)" \
		| cc_atomic_write "$cc_vr_dir/attempts/$cc_vr_pad/verifier.yaml"
	if [ "$cc_vr_out" = "passed" ]; then
		cc_exec_set "$cc_vr_dir" status verified
		cc_emit outcome passed
		cc_emit status verified
		return 0
	fi
	if [ "$cc_vr_out" = "blocked" ]; then
		cc_exec_set "$cc_vr_dir" status blocked
		cc_emit outcome blocked
		cc_emit status blocked
		return 0
	fi
	if [ "$cc_vr_out" = "waived" ]; then
		# A human limitation decision: non-passing, but not a worker fault.
		# It never satisfies verification and never increments the failure counter.
		cc_exec_set "$cc_vr_dir" status blocked
		cc_emit outcome waived
		cc_emit status blocked
		return 0
	fi
	# failed: the independent verifier rejected the worker result
	cc_vr_wf=$(cc_scalar "$cc_vr_dir/execution.yaml" "worker_failures")
	cc_vr_wf=$((cc_vr_wf + 1))
	cc_exec_set "$cc_vr_dir" worker_failures "$cc_vr_wf"
	if [ "$cc_vr_wf" -ge 3 ]; then
		cc_exec_set "$cc_vr_dir" status failed
		cc_emit outcome "$cc_vr_out"
		cc_emit worker_failures "$cc_vr_wf"
		cc_emit status failed
		cc_emit stop FAILURE_LIMIT_REACHED
		return 0
	fi
	cc_exec_set "$cc_vr_dir" status repairing
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
	case "$cc_ae_att" in ''|*[!0-9]*) cc_fail ATTEMPT_INVALID "$cc_ae_att"; return 1 ;; esac
	cc_ae_pad=$(printf '%03d' "$cc_ae_att")
	[ -d "$cc_ae_dir/attempts/$cc_ae_pad" ] || { cc_fail ATTEMPT_UNKNOWN "$cc_ae_att"; return 1; }
	[ "$#" -gt 0 ] || { cc_fail ATTEMPT_EVIDENCE_EMPTY; return 1; }
	cc_ae_f="$cc_ae_dir/attempts/$cc_ae_pad/host-evidence.yaml"
	[ -f "$cc_ae_f" ] || printf 'attempt: %s\n' "$cc_ae_att" | cc_atomic_write "$cc_ae_f"
	for cc_ae_kv in "$@"; do
		case "$cc_ae_kv" in *=*) : ;; *) cc_fail ATTEMPT_EVIDENCE_MALFORMED "$cc_ae_kv"; return 1 ;; esac
		cc_ae_k=${cc_ae_kv%%=*}; cc_ae_v=${cc_ae_kv#*=}
		case "$cc_ae_k" in
			worker_model|worker_effort|verifier_model|verifier_effort|complexity|escalated) : ;;
			*) cc_fail ATTEMPT_EVIDENCE_KEY_UNKNOWN "$cc_ae_k"; return 1 ;;
		esac
		case "$cc_ae_v" in
			''|*[!A-Za-z0-9._-]*) cc_fail ATTEMPT_EVIDENCE_VALUE_INVALID "$cc_ae_k"; return 1 ;;
		esac
		awk -v k="$cc_ae_k" -v v="$cc_ae_v" '
			$0 ~ ("^" k ":[[:space:]]") { next }
			{ print }
			END { print k ": " v }
		' "$cc_ae_f" | cc_atomic_write "$cc_ae_f"
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
# (running/verifying/repairing) is NOT runnable; only an approved, never-run plan
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
		cc_rsr_status=$(cc_plan_status "$cc_rsr_root" "$cc_rsr_plan" 2>/dev/null) || cc_rsr_status=""
		if [ "$cc_rsr_status" != "approved" ]; then
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

# cc_completion_ready ROOT PLAN -> eligible only if latest execution is verified
cc_completion_ready() {
	cc_cr_exec=$(cc_latest_execution "$1" "$2") || { cc_fail COMPLETION_NO_EXECUTION; return 1; }
	[ -n "$cc_cr_exec" ] || { cc_fail COMPLETION_NO_EXECUTION; return 1; }
	cc_cr_dir=$(cc_execution_dir "$1" "$2" "$cc_cr_exec")
	cc_cr_st=$(cc_execution_status "$cc_cr_dir")
	[ "$cc_cr_st" = "verified" ] || { cc_fail COMPLETION_NOT_VERIFIED "$cc_cr_st"; return 1; }
	cc_emit completion eligible
	cc_emit execution_id "$cc_cr_exec"
	return 0
}

# cc_plan_complete ROOT PLAN -> approved->done ONLY when verified; write completion record
cc_plan_complete() {
	cc_pc_root="$1"; cc_pc_plan="$2"
	cc_completion_ready "$cc_pc_root" "$cc_pc_plan" >/dev/null || { cc_fail COMPLETION_BLOCKED; return 1; }
	cc_pc_exec=$(cc_latest_execution "$cc_pc_root" "$cc_pc_plan")
	cc_pc_edir=$(cc_execution_dir "$cc_pc_root" "$cc_pc_plan" "$cc_pc_exec")
	cc_pc_yaml="$cc_pc_root/plans/$cc_pc_plan/plan.yaml"
	cc_pc_status=$(cc_scalar "$cc_pc_yaml" "status")
	[ "$cc_pc_status" = "approved" ] || { cc_fail COMPLETION_PLAN_NOT_APPROVED "$cc_pc_status"; return 1; }
	# implementation completion record with the accepted commits
	{
		printf 'schema_version: %s\nexecution_id: %s\nplan: %s\nhuman_completion: accepted\ncompleted_at: %s\ncommits:\n' \
			"$CC_SCHEMA_VERSION" "$cc_pc_exec" "$cc_pc_plan" "$(cc_now)"
		for cc_pc_rf in "$cc_pc_edir"/repositories/*.yaml; do
			[ -f "$cc_pc_rf" ] || continue
			printf '  %s: %s\n' "$(cc_scalar "$cc_pc_rf" repository)" "$(cc_scalar "$cc_pc_rf" latest_commit)"
		done
	} | cc_atomic_write "$cc_pc_edir/completion.yaml"
	cc_plan_set_status "$cc_pc_yaml" "done"
	cc_plan_index_upsert "$cc_pc_root" "$cc_pc_plan" >/dev/null
	cc_emit plan "$cc_pc_plan"
	cc_emit status done
	cc_emit implementation_completion recorded
	cc_emit knowledge_impact review-pending
	return 0
}

# cc_context_impact_record EXEC_DIR FILE -> store reconciliation refs (no PK interpretation)
cc_context_impact_record() {
	[ -f "$2" ] || { cc_fail CONTEXT_IMPACT_FILE_MISSING; return 1; }
	cp "$2" "$1/context-impact.yaml"
	cc_emit context_impact recorded
	return 0
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
		cc_dt_tgt=$(cc_scalar "$cc_dt_rf" anchor_branch)
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
# latest execution, whether the recorded base has drifted from the current anchor
# tip (a sibling merged and advanced the anchor). Drift = the current anchor tip
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
		cc_dd_anchor=$(cc_scalar "$cc_dd_rf" anchor_branch)
		cc_dd_base=$(cc_scalar "$cc_dd_rf" base_commit)
		cc_dd_latest=$(cc_scalar "$cc_dd_rf" latest_commit)
		cc_dd_abs=$(cc_repo_resolve "$cc_dd_root" "$cc_dd_id" 2>/dev/null | sed -n 's/^path: //p')
		cc_dd_tip=$(git -C "$cc_dd_abs" rev-parse --verify "refs/heads/$cc_dd_anchor" 2>/dev/null) || cc_dd_tip=""
		cc_dd_drift=unknown
		if [ -n "$cc_dd_abs" ] && [ -n "$cc_dd_tip" ]; then
			if git -C "$cc_dd_abs" merge-base --is-ancestor "$cc_dd_tip" "$cc_dd_latest" 2>/dev/null; then
				cc_dd_drift=false
			else
				cc_dd_drift=true; cc_dd_any=true
			fi
		fi
		cc_emit "repository" "$cc_dd_id"
		cc_emit "  anchor_tip" "${cc_dd_tip:-unknown}"
		cc_emit "  recorded_base" "$cc_dd_base"
		cc_emit "  drifted" "$cc_dd_drift"
	done
	cc_emit drift_detected "$cc_dd_any"
	return 0
}

# cc_delivery_rebase ROOT PLAN -> rebase every drifted repository's execution
# branch onto the current anchor tip, update base_commit/latest_commit, and flag
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
		cc_dr_anchor=$(cc_scalar "$cc_dr_rf" anchor_branch)
		cc_dr_wt=$(cc_scalar "$cc_dr_rf" worktree)
		cc_dr_latest=$(cc_scalar "$cc_dr_rf" latest_commit)
		cc_dr_abs=$(cc_repo_resolve "$cc_dr_root" "$cc_dr_id" 2>/dev/null | sed -n 's/^path: //p')
		[ -n "$cc_dr_abs" ] || continue
		cc_dr_tip=$(git -C "$cc_dr_abs" rev-parse --verify "refs/heads/$cc_dr_anchor" 2>/dev/null) || continue
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
		cc_exec_set "$cc_dr_edir" status verifying
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
		repository-resolve)      cc_repo_resolve "$@" ;;
		repository-preflight)    cc_repository_preflight "$@" ;;
		delivery-targets)        cc_delivery_targets "$@" ;;
		worktree-prepare)        cc_worktree_prepare "$@" ;;
		pair-begin)              cc_pair_begin "$@" ;;
		pair-inspect)            cc_pair_inspect "$@" ;;
		pair-close)              cc_pair_close "$@" ;;
		pair-delivery-targets)   cc_pair_delivery_targets "$@" ;;
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
		plan-approve)            cc_plan_approve "$@" ;;
		plan-archive)            cc_plan_archive "$@" ;;
		plan-restore)            cc_plan_restore "$@" ;;
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
		completion-ready)        cc_completion_ready "$@" ;;
		plan-complete)           cc_plan_complete "$@" ;;
		context-impact-record)   cc_context_impact_record "$@" ;;
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
