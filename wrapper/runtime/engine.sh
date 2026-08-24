#!/bin/sh
# Context Circuit v0.5 runtime engine.
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
# the three-failure counter, one-writer locking, completion eligibility,
# implementation completion records, context-impact handoff references, and
# recovery inspection.
#
# This runtime does NOT own: provider-specific child launch, model prompts,
# Product Knowledge interpretation, plan-writing intelligence, conversational
# routing policy, confirmation cards/tokens, product test semantics, or any
# automatic pull-request/merge/push/publish/deploy/cleanup/completion action.
#
# POSIX sh only (runs under dash). No bashisms, no `local`, no arrays.
# Every function returns 0 on success and non-zero with a reason code on stderr.

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

CC_RUNTIME_VERSION="0.5.0"
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
		"$cc_wi_root/plans/.archived" "$cc_wi_root/.runtime/executions" \
		"$cc_wi_root/.runtime/worktrees" "$cc_wi_root/.runtime/locks" \
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
		{
			awk '/^repositories:[[:space:]]*(\[\])?[[:space:]]*$/{print "repositories:"; next}{print}' "$cc_reg_ws"
			printf '  - id: %s\n' "$cc_reg_id"
			[ -n "$cc_reg_url" ] && printf '    canonical_url: %s\n' "$cc_reg_url" || :
			[ -n "$cc_reg_default" ] && printf '    default_branch: %s\n' "$cc_reg_default" || :
		} | cc_atomic_write "$cc_reg_ws"
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
	for cc_ai_d in "$cc_ai_root/plans"/*/ "$cc_ai_root/plans/.archived"/*/; do
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

# cc_plan_archive ROOT PLAN -> move plans/PLAN -> plans/.archived/PLAN; drop index row
cc_plan_archive() {
	cc_ar_root="$1"; cc_ar_plan="$2"
	cc_safe_id "$cc_ar_plan" || { cc_fail PLAN_ID_UNSAFE "$cc_ar_plan"; return 1; }
	cc_ar_src="$cc_ar_root/plans/$cc_ar_plan"
	cc_ar_dst="$cc_ar_root/plans/.archived/$cc_ar_plan"
	[ -d "$cc_ar_src" ] || { cc_fail ARCHIVE_SOURCE_MISSING "$cc_ar_plan"; return 1; }
	[ -e "$cc_ar_dst" ] && { cc_fail ARCHIVE_TARGET_COLLISION "$cc_ar_plan"; return 1; }
	cc_plan_org_lock "$cc_ar_root" || return 1
	mkdir -p "$cc_ar_root/plans/.archived"
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

# cc_plan_restore ROOT PLAN -> move plans/.archived/PLAN -> plans/PLAN; re-add index
cc_plan_restore() {
	cc_re_root="$1"; cc_re_plan="$2"
	cc_safe_id "$cc_re_plan" || { cc_fail PLAN_ID_UNSAFE "$cc_re_plan"; return 1; }
	cc_re_src="$cc_re_root/plans/.archived/$cc_re_plan"
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
# One-writer execution locks
# ---------------------------------------------------------------------------

# cc_lock_acquire ROOT PLAN OWNER -> exclusive-create writer lock for a plan
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
	# per-repository branch + worktree
	for cc_eb_id in $(cc_plan_affected_repositories "$cc_eb_dir/plan.yaml"); do
		cc_eb_out=$(cc_worktree_prepare "$cc_eb_root" "$cc_eb_plan" "$cc_eb_id") \
			|| { cc_fail EXECUTION_WORKTREE_FAILED "$cc_eb_id"; return 1; }
		cc_eb_wt=$(printf '%s' "$cc_eb_out" | sed -n 's/^worktree: //p')
		cc_eb_br=$(printf '%s' "$cc_eb_out" | sed -n 's/^branch: //p')
		cc_eb_bc=$(printf '%s' "$cc_eb_out" | sed -n 's/^base_commit: //p')
		cc_eb_an=$(cc_binding_field "$cc_eb_root" "$cc_eb_id" "anchor_branch")
		cc_eb_paths=$(cc_plan_repo_paths "$cc_eb_dir/plan.yaml" "$cc_eb_id" | tr '\n' ',' | sed 's/,$//; s/,/, /g')
		{
			printf 'repository: %s\nworktree: %s\nbranch: %s\nanchor_branch: %s\nbase_commit: %s\nlatest_commit: %s\nallowed_paths: [%s]\n' \
				"$cc_eb_id" "$cc_eb_wt" "$cc_eb_br" "$cc_eb_an" "$cc_eb_bc" "$cc_eb_bc" "$cc_eb_paths"
		} | cc_atomic_write "$cc_eb_edir/repositories/$cc_eb_id.yaml"
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
