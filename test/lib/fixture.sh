#!/bin/sh
# Shared fixture helpers for Context Circuit v0.5 semantic tests.
# Sourced after test/lib/assert.sh (which provides ROOT and the engine).

# cc_fx_ws -> create an initialized tmp workspace; prints its root path
cc_fx_ws() {
	cc_fx_root=$(mktemp -d "${TMPDIR:-/tmp}/cc-fx.XXXXXX") || return 1
	cc_fx_root=$(cd "$cc_fx_root" && pwd -P) || return 1
	printf 'schema_version: 1\nworkspace: fx-ws\ntitle: Fixture\nrepositories: []\n' >"$cc_fx_root/workspace.yaml"
	cc_workspace_init "$cc_fx_root" >/dev/null
	printf '# Project\n\nFixture project knowledge.\n' >"$cc_fx_root/context/PROJECT.md"
	printf '%s' "$cc_fx_root"
}

# cc_fx_bindings_header WS -> ensure repositories.local.yaml has a real bindings block
cc_fx_bindings_header() {
	if [ ! -f "$1/repositories.local.yaml" ] || grep -q '^bindings: {}' "$1/repositories.local.yaml" 2>/dev/null || ! grep -q '^bindings:' "$1/repositories.local.yaml" 2>/dev/null; then
		printf 'schema_version: 1\nbindings:\n' >"$1/repositories.local.yaml"
	fi
}

# cc_fx_repo WS ID BRANCH [PATH] -> create a git repo, seed a commit, bind it
cc_fx_repo() {
	cc_fxr_ws=$1; cc_fxr_id=$2; cc_fxr_br=$3; cc_fxr_path=${4:-repositories/$2}
	mkdir -p "$cc_fxr_ws/$cc_fxr_path/src"
	git -C "$cc_fxr_ws/$cc_fxr_path" init -q -b "$cc_fxr_br"
	git -C "$cc_fxr_ws/$cc_fxr_path" config user.email t@t.t
	git -C "$cc_fxr_ws/$cc_fxr_path" config user.name t
	printf 'seed\n' >"$cc_fxr_ws/$cc_fxr_path/src/seed.txt"
	git -C "$cc_fxr_ws/$cc_fxr_path" add -A
	git -C "$cc_fxr_ws/$cc_fxr_path" commit -q -m seed
	cc_fx_bindings_header "$cc_fxr_ws"
	printf '  %s:\n    path: %s\n    anchor_branch: %s\n' "$cc_fxr_id" "$cc_fxr_path" "$cc_fxr_br" \
		>>"$cc_fxr_ws/repositories.local.yaml"
}

# cc_fx_plan WS PID TITLE "repo1 repo2..." -> write a minimal valid draft plan
# One task per repository, named <REPO>-001, dependencies chained in order.
cc_fx_plan() {
	cc_fxp_ws=$1; cc_fxp_pid=$2; cc_fxp_title=$3; cc_fxp_repos=$4
	cc_fxp_dir="$cc_fxp_ws/plans/$cc_fxp_pid"
	# v1.0: every plan derives from an approved parent intent within a scope envelope.
	# Author + approve an intent i<pid> scoped to each repo's src, then bind the plan.
	cc_fxp_iid="i$cc_fxp_pid"
	mkdir -p "$cc_fxp_ws/intent/$cc_fxp_iid"
	{
		printf 'schema_version: 1\nintent: %s\ntitle: %s\ngoal: %s goal.\n' "$cc_fxp_iid" "$cc_fxp_title" "$cc_fxp_title"
		printf 'non_goals:\n  - none\nconstraints:\n  - none\n'
		printf 'acceptance_criteria:\n  - id: ac-1\n    statement: %s works.\n    method: test\n' "$cc_fxp_title"
		printf 'done_when: ac-1 passes and a human accepts the candidate.\n'
		printf 'scope:\n  repositories:\n'
		for cc_fxp_r in $cc_fxp_repos; do
			printf '    - id: %s\n      paths: [src]\n' "$cc_fxp_r"
		done
		printf 'tier: standard\nstatus: draft\ncontract_digest:\n'
	} >"$cc_fxp_ws/intent/$cc_fxp_iid/contract.yaml"
	printf '# %s\n' "$cc_fxp_title" >"$cc_fxp_ws/intent/$cc_fxp_iid/INTENT.md"
	cc_intent_index_upsert "$cc_fxp_ws" "$cc_fxp_iid" >/dev/null
	cc_intent_approve "$cc_fxp_ws" "$cc_fxp_iid" >/dev/null
	mkdir -p "$cc_fxp_dir/tasks"
	{
		printf 'schema_version: 3\nplan: %s\ntitle: %s\nstatus: draft\nobjective: %s\nintent: %s\n' \
			"$cc_fxp_pid" "$cc_fxp_title" "$cc_fxp_title objective" "$cc_fxp_iid"
		printf 'repositories:\n'
		for cc_fxp_r in $cc_fxp_repos; do
			printf '  - id: %s\n    purpose: %s scope\n' "$cc_fxp_r" "$cc_fxp_r"
		done
		printf 'product_knowledge:\n  - id: project.core\n    path: context/PROJECT.md\n    reason: Grounds the objective.\n'
		printf 'context_grounding:\n  summary: Fixture grounding.\n  constraints: []\n  decisions: []\n'
		printf 'knowledge_impact:\n  expected_context_units: []\n  review_on_completion: true\n'
		printf 'tasks:\n'
		cc_fxp_prev=""
		for cc_fxp_r in $cc_fxp_repos; do
			cc_fxp_tid=$(printf '%s' "$cc_fxp_r" | tr '[:lower:]' '[:upper:]')-001
			printf '  - id: %s\n    title: Work in %s\n    repositories: [%s]\n    paths: [src]\n' \
				"$cc_fxp_tid" "$cc_fxp_r" "$cc_fxp_r"
			if [ -n "$cc_fxp_prev" ]; then
				printf '    depends_on: [%s]\n' "$cc_fxp_prev"
			else
				printf '    depends_on: []\n'
			fi
			printf '    changes:\n      - Change %s.\n    acceptance:\n      - id: %s-AC\n        statement: %s changed.\n    verification:\n      - id: %s-VT\n        command: test/%s.sh\n' \
				"$cc_fxp_r" "$cc_fxp_tid" "$cc_fxp_r" "$cc_fxp_tid" "$cc_fxp_r"
			cc_fxp_prev=$cc_fxp_tid
		done
		printf 'execution:\n  worker: one\n  independent_verifier: required\n  max_worker_failures: 3\n'
	} >"$cc_fxp_dir/plan.yaml"
	printf '# %s\n\nObjective: %s objective.\n\nStatus: draft.\n' "$cc_fxp_title" "$cc_fxp_title" >"$cc_fxp_dir/PLAN.md"
	cc_plan_index_upsert "$cc_fxp_ws" "$cc_fxp_pid" >/dev/null
}

# cc_fx_plan_ex WS PID TITLE REPO PATH "dep1 dep2..." -> write a valid draft v1.0
# plan (schema 3) with ONE task scoped to PATH in REPO, whose acceptance is trivially
# achievable (a file PATH/mod.txt), deriving from an approved parent intent i<pid>
# scoped to REPO:PATH. When deps are given a plan_dependencies block is added. Used by
# the run-stack and lease gates to build inter-plan dependency graphs.
cc_fx_plan_ex() {
	cc_fxe_ws=$1; cc_fxe_pid=$2; cc_fxe_title=$3; cc_fxe_repo=$4; cc_fxe_path=$5; cc_fxe_deps=${6:-}
	cc_fxe_dir="$cc_fxe_ws/plans/$cc_fxe_pid"
	cc_fxe_iid="i$cc_fxe_pid"
	mkdir -p "$cc_fxe_ws/intent/$cc_fxe_iid"
	{
		printf 'schema_version: 1\nintent: %s\ntitle: %s\ngoal: %s goal.\n' "$cc_fxe_iid" "$cc_fxe_title" "$cc_fxe_title"
		printf 'non_goals:\n  - none\nconstraints:\n  - none\n'
		printf 'acceptance_criteria:\n  - id: ac-1\n    statement: %s works.\n    method: test\n' "$cc_fxe_title"
		printf 'done_when: ac-1 passes and a human accepts the candidate.\n'
		printf 'scope:\n  repositories:\n    - id: %s\n      paths: [%s]\n' "$cc_fxe_repo" "$cc_fxe_path"
		printf 'tier: standard\nstatus: draft\ncontract_digest:\n'
	} >"$cc_fxe_ws/intent/$cc_fxe_iid/contract.yaml"
	printf '# %s\n' "$cc_fxe_title" >"$cc_fxe_ws/intent/$cc_fxe_iid/INTENT.md"
	cc_intent_index_upsert "$cc_fxe_ws" "$cc_fxe_iid" >/dev/null
	cc_intent_approve "$cc_fxe_ws" "$cc_fxe_iid" >/dev/null
	mkdir -p "$cc_fxe_dir/tasks"
	cc_fxe_tid=$(printf '%s' "$cc_fxe_repo" | tr '[:lower:]' '[:upper:]')-001
	{
		printf 'schema_version: 3\nplan: %s\ntitle: %s\nstatus: draft\nobjective: %s objective\nintent: %s\n' \
			"$cc_fxe_pid" "$cc_fxe_title" "$cc_fxe_title" "$cc_fxe_iid"
		printf 'repositories:\n  - id: %s\n    purpose: %s scope\n' "$cc_fxe_repo" "$cc_fxe_repo"
		if [ -n "$cc_fxe_deps" ]; then
			printf 'plan_dependencies:\n'
			for cc_fxe_d in $cc_fxe_deps; do
				printf '  - id: %s\n    reason: builds on %s\n' "$cc_fxe_d" "$cc_fxe_d"
			done
		fi
		printf 'product_knowledge:\n  - id: project.core\n    path: context/PROJECT.md\n    reason: Grounds the objective.\n'
		printf 'context_grounding:\n  summary: Fixture grounding.\n  constraints: []\n  decisions: []\n'
		printf 'knowledge_impact:\n  expected_context_units: []\n  review_on_completion: true\n'
		printf 'tasks:\n  - id: %s\n    title: Work in %s\n    repositories: [%s]\n    paths: [%s]\n    depends_on: []\n' \
			"$cc_fxe_tid" "$cc_fxe_path" "$cc_fxe_repo" "$cc_fxe_path"
		printf '    changes:\n      - Create %s/mod.txt.\n' "$cc_fxe_path"
		printf '    acceptance:\n      - id: %s-AC\n        statement: %s/mod.txt exists.\n' "$cc_fxe_tid" "$cc_fxe_path"
		printf '    verification:\n      - id: %s-VT\n        command: test -f %s/mod.txt\n' "$cc_fxe_tid" "$cc_fxe_path"
		printf 'execution:\n  worker: one\n  independent_verifier: required\n  max_worker_failures: 3\n'
	} >"$cc_fxe_dir/plan.yaml"
	printf '# %s\n\nObjective: %s objective.\n' "$cc_fxe_title" "$cc_fxe_title" >"$cc_fxe_dir/PLAN.md"
	cc_plan_index_upsert "$cc_fxe_ws" "$cc_fxe_pid" >/dev/null
}

# cc_fx_run_ok WS PID REPO PATH -> approve/execute/worker-commit (PATH/mod.txt)/verify
# passed. Leaves the plan verified. Uses OWNER=<pid>-w.
cc_fx_run_ok() {
	cc_fxo_ws=$1; cc_fxo_pid=$2; cc_fxo_repo=$3; cc_fxo_path=$4
	[ "$(cc_plan_status "$cc_fxo_ws" "$cc_fxo_pid" 2>/dev/null)" = draft ] \
		&& { cc_plan_approve "$cc_fxo_ws" "$cc_fxo_pid" >/dev/null || return 1; } || :
	cc_fxo_exec=$(cc_execution_begin "$cc_fxo_ws" "$cc_fxo_pid" "$cc_fxo_pid-w" | sed -n 's/^execution_id: //p') || return 1
	cc_fxo_edir="$cc_fxo_ws/.runtime/executions/$cc_fxo_pid/$cc_fxo_exec"
	cc_fxo_wt="$cc_fxo_ws/.runtime/worktrees/$cc_fxo_pid/$cc_fxo_repo"
	cc_attempt_begin "$cc_fxo_edir" >/dev/null || return 1
	mkdir -p "$cc_fxo_wt/$cc_fxo_path"
	printf 'mod\n' >"$cc_fxo_wt/$cc_fxo_path/mod.txt"
	git -C "$cc_fxo_wt" add -A
	git -C "$cc_fxo_wt" commit -q -m "feat($cc_fxo_repo): add $cc_fxo_path/mod.txt"
	cc_worker_commit_record "$cc_fxo_edir" "$cc_fxo_repo" implementation >/dev/null || return 1
	cc_verifier_prepare "$cc_fxo_edir" >/dev/null || return 1
	cc_verifier_result_record "$cc_fxo_edir" 1 passed >/dev/null || return 1
}

# cc_fx_commit WS PID REPO [MSG] -> make a worker commit in the execution worktree
cc_fx_commit() {
	cc_fxc_ws=$1; cc_fxc_pid=$2; cc_fxc_repo=$3; cc_fxc_msg=${4:-work}
	cc_fxc_wt="$cc_fxc_ws/.runtime/worktrees/$cc_fxc_pid/$cc_fxc_repo"
	printf '%s\n' "$cc_fxc_msg $(date -u +%s%N 2>/dev/null || echo x)" >>"$cc_fxc_wt/src/change.txt"
	git -C "$cc_fxc_wt" add -A
	git -C "$cc_fxc_wt" commit -q -m "$cc_fxc_msg"
}

# cc_fx_exec_dir WS PID EXEC -> path to an execution directory
cc_fx_exec_dir() { printf '%s/.runtime/executions/%s/%s' "$1" "$2" "$3"; }

# cc_fx_intent WS IID TITLE REPO "path1 path2..." [TIER] -> author a draft intent
# contract (goal + one test criterion + scope envelope) and its INTENT.md. Does
# NOT approve. Prints nothing. TIER defaults to standard.
cc_fx_intent() {
	cc_fxi_ws=$1; cc_fxi_id=$2; cc_fxi_title=$3; cc_fxi_repo=$4; cc_fxi_paths=$5; cc_fxi_tier=${6:-standard}
	cc_fxi_dir="$cc_fxi_ws/intent/$cc_fxi_id"
	mkdir -p "$cc_fxi_dir"
	cc_fxi_inline=$(printf '%s' "$cc_fxi_paths" | tr ' ' ',' | sed 's/,/, /g')
	{
		printf 'schema_version: 1\nintent: %s\ntitle: %s\n' "$cc_fxi_id" "$cc_fxi_title"
		printf 'goal: %s goal.\n' "$cc_fxi_title"
		printf 'non_goals:\n  - none\nconstraints:\n  - none\n'
		printf 'acceptance_criteria:\n  - id: ac-1\n    statement: %s works.\n    method: test\n    surface: %s\n' "$cc_fxi_title" "$cc_fxi_repo"
		printf 'done_when: ac-1 passes and a human accepts the candidate.\n'
		printf 'scope:\n  repositories:\n    - id: %s\n      paths: [%s]\n' "$cc_fxi_repo" "$cc_fxi_inline"
		printf 'tier: %s\nstatus: draft\ncontract_digest:\n' "$cc_fxi_tier"
	} >"$cc_fxi_dir/contract.yaml"
	printf '# %s\n\nGoal: %s.\n' "$cc_fxi_title" "$cc_fxi_title" >"$cc_fxi_dir/INTENT.md"
	cc_intent_index_upsert "$cc_fxi_ws" "$cc_fxi_id" >/dev/null
}

# cc_fx_plan_intent WS PID TITLE REPO PATH IID -> a valid schema-3 draft plan bound
# to intent IID, one task scoped to PATH in REPO whose acceptance is a file
# PATH/mod.txt (achievable by cc_fx_run_ok). Adds the active index row.
cc_fx_plan_intent() {
	cc_fxpi_ws=$1; cc_fxpi_pid=$2; cc_fxpi_title=$3; cc_fxpi_repo=$4; cc_fxpi_path=$5; cc_fxpi_iid=$6
	cc_fxpi_dir="$cc_fxpi_ws/plans/$cc_fxpi_pid"
	mkdir -p "$cc_fxpi_dir/tasks"
	cc_fxpi_tid=$(printf '%s' "$cc_fxpi_repo" | tr '[:lower:]' '[:upper:]')-001
	{
		printf 'schema_version: 3\nplan: %s\ntitle: %s\nstatus: draft\nobjective: %s objective\nintent: %s\n' \
			"$cc_fxpi_pid" "$cc_fxpi_title" "$cc_fxpi_title" "$cc_fxpi_iid"
		printf 'repositories:\n  - id: %s\n    purpose: %s scope\n' "$cc_fxpi_repo" "$cc_fxpi_repo"
		printf 'product_knowledge:\n  - id: project.core\n    path: context/PROJECT.md\n    reason: Grounds the objective.\n'
		printf 'context_grounding:\n  summary: Fixture grounding.\n  constraints: []\n  decisions: []\n'
		printf 'knowledge_impact:\n  expected_context_units: []\n  review_on_completion: true\n'
		printf 'tasks:\n  - id: %s\n    title: Work in %s\n    repositories: [%s]\n    paths: [%s]\n    depends_on: []\n' \
			"$cc_fxpi_tid" "$cc_fxpi_path" "$cc_fxpi_repo" "$cc_fxpi_path"
		printf '    changes:\n      - Create %s/mod.txt.\n' "$cc_fxpi_path"
		printf '    acceptance:\n      - id: %s-AC\n        statement: %s/mod.txt exists.\n' "$cc_fxpi_tid" "$cc_fxpi_path"
		printf '    verification:\n      - id: %s-VT\n        command: test -f %s/mod.txt\n' "$cc_fxpi_tid" "$cc_fxpi_path"
		printf 'execution:\n  worker: one\n  independent_verifier: required\n  max_worker_failures: 3\n'
	} >"$cc_fxpi_dir/plan.yaml"
	printf '# %s\n\nObjective: %s objective.\n' "$cc_fxpi_title" "$cc_fxpi_title" >"$cc_fxpi_dir/PLAN.md"
	cc_plan_index_upsert "$cc_fxpi_ws" "$cc_fxpi_pid" >/dev/null
}
