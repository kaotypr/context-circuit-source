import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { parse as parseYaml } from "yaml";
import { validateContract, workspaceDocumentErrors, workspaceSemanticErrors } from "../scripts/lib/validation.js";
import type { WorkspaceConfig } from "../scripts/lib/types.js";
import { projectRoot } from "./helpers.js";

test("workspace configuration validates", async () => {
  const config = parseYaml(await readFile(join(projectRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  assert.deepEqual(await validateContract("workspace", config), []);
  assert.deepEqual(workspaceSemanticErrors(config), []);
});

test("deterministic TypeScript commands avoid sandbox-incompatible tsx IPC", async () => {
  const packageJson = JSON.parse(await readFile(join(projectRoot, "package.json"), "utf8")) as {
    scripts?: Record<string, string>;
  };
  const deterministicCommands = [
    "configure-workspace",
    "create-plan",
    "confirm-merge",
    "fixture:create",
    "finish-work",
    "initialize-workspace",
    "prepare-repair",
    "prepare-review",
    "prepare-lifecycle",
    "prepare-plan-publication",
    "record-lifecycle-action",
    "record-plan-publication",
    "record-result",
    "run-task",
    "set-plan-state",
    "validate",
    "validate-plan",
    "whats-next",
  ];

  for (const command of deterministicCommands) {
    assert.match(packageJson.scripts?.[command] ?? "", /^node --import tsx\b/, command);
  }
});

test("workspace rejects credential fields and duplicate repository paths", async () => {
  const config = {
    version: 1,
    template_version: "0.2.1",
    workspace: { name: "bad", mode: "team", default_branch: "main" },
    repositories: {
      one: { path: "repositories/same", mode: "ignored-clone", role: "app", agent: "frontend", default_branch: "main", token_env: "SECRET" },
      two: { path: "repositories/same", mode: "submodule", role: "app", agent: "frontend", default_branch: "main" },
    },
    activity: { provider: "none", access: "auto", required_capabilities: [], optional_capabilities: [] },
    workflow: { human_gates: ["merge"], maximum_repair_attempts: 2, wrapper_change_policy: "pull-request" },
  } as unknown as WorkspaceConfig;
  assert.notEqual((await validateContract("workspace", config)).length, 0);
  assert.match(workspaceSemanticErrors(config).join("\n"), /duplicates/);
});

test("workspace lifecycle policy rejects undeclared and misclassified capabilities", () => {
  const config = {
    version: 1,
    template_version: "0.2.1",
    workspace: { name: "bad-lifecycle", mode: "team", default_branch: "main" },
    repositories: { frontend: { path: "repositories/frontend", mode: "ignored-clone", role: "app", agent: "frontend", default_branch: "main" } },
    activity: {
      provider: "example", access: "auto", required_capabilities: [], optional_capabilities: ["timers"],
      lifecycle: {
        "task.starting": [
          { id: "claim", capability: "assign-task", policy: "required", description: "Claim the task." },
          { id: "claim", capability: "timers", policy: "optional", description: "Start the timer." },
        ],
      },
    },
    workflow: { human_gates: ["merge"], maximum_repair_attempts: 2, wrapper_change_policy: "pull-request" },
  } as unknown as WorkspaceConfig;
  const errors = workspaceSemanticErrors(config).join("\n");
  assert.match(errors, /undeclared capability/);
  assert.match(errors, /required lifecycle action/);
  assert.match(errors, /duplicate action id/);
});

test("workspace document validation reports missing required files", async () => {
  const config = parseYaml(await readFile(join(projectRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  const errors = await workspaceDocumentErrors(join(projectRoot, "fixtures", "react-app"), config);
  assert.ok(errors.some((error) => error.includes("README.md")));
  assert.ok(errors.some((error) => error.includes("workspace-bootstrap-request.schema.json")));
});

test("machine-readable contract schemas reject incomplete data", async () => {
  for (const schema of ["workspace-bootstrap-request", "workspace-configure-request", "task-brief", "worker-result", "verifier-result", "runtime-manifest", "review-preparation", "merge-confirmation-record", "closeout-record", "plan-index", "plan-work-breakdown", "plan-draft-request", "work-candidate", "fake-activity-source", "whats-next-result", "activity-lifecycle-record", "plan-publication-discovery", "plan-publication-record"] as const) {
    assert.notEqual((await validateContract(schema, { contract_version: 1 })).length, 0, schema);
  }
});

test("complete worker and verifier results validate", async () => {
  const common = {
    contract_version: 1,
    work_id: "ADHOC-20260811-001",
    run_id: "20260811T083000Z-abcd1234",
    repository: "frontend",
    summary: "Implemented and checked the requested change.",
  };
  const worker = {
    ...common,
    status: "completed",
    branch: "agent/adhoc-20260811-001-reset-abcd1234",
    worktree: "/tmp/worktree",
    commits: ["0123456789abcdef0123456789abcdef01234567"],
    changed_files: ["src/App.tsx"],
    checks: [{ command: "npm test", status: "passed", evidence: "1 test passed" }],
    risks: [],
  };
  const verifier = {
    ...common,
    status: "pass",
    acceptance: [{ criterion: "Reset returns count to zero", status: "passed", evidence: "Test passed" }],
    checks: ["npm test"],
    findings: [],
    verified_at: "2026-08-11T08:45:00.000Z",
  };
  assert.deepEqual(await validateContract("worker-result", worker), []);
  assert.deepEqual(await validateContract("verifier-result", verifier), []);
});

test("task brief rejects required test policy without a test path", async () => {
  const invalid = {
    contract_version: 1,
    work_id: "ADHOC-20260811-001",
    run_id: "20260811T083000Z-abcd1234",
    source: { kind: "direct-request" },
    requested_outcome: "Change behavior",
    scope: ["src/App.tsx"],
    implementation_scope: ["src/App.tsx"],
    test_expectation: { policy: "required", paths: [], rationale: "Tests must change." },
    acceptance_criteria: ["Behavior changes"],
    repositories: [{ name: "frontend", dependency_order: 0 }],
    plan: { reference: null, approval_state: "not-applicable" },
    activity: { reference: null, claim_status: "not-applicable", duplicate_effort_warning: true },
    assumptions: [],
    risks: [],
    verification_commands: ["npm test"],
    authorization: { kind: "explicit-user-request", evidence: "Explicit request" },
    created_at: "2026-08-11T08:30:00.000Z",
  };
  assert.notEqual((await validateContract("task-brief", invalid)).length, 0);
});
