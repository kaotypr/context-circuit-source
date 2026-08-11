import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { parse as parseYaml } from "yaml";
import { validateContract, workspaceSemanticErrors } from "../scripts/lib/validation.js";
import type { WorkspaceConfig } from "../scripts/lib/types.js";
import { projectRoot } from "./helpers.js";

test("workspace configuration validates", async () => {
  const config = parseYaml(await readFile(join(projectRoot, "workspace.yaml"), "utf8")) as WorkspaceConfig;
  assert.deepEqual(await validateContract("workspace", config), []);
  assert.deepEqual(workspaceSemanticErrors(config), []);
});

test("workspace rejects credential fields and duplicate repository paths", async () => {
  const config = {
    version: 1,
    template_version: "0.1.0",
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

test("worker, verifier, task brief, and manifest schemas reject incomplete data", async () => {
  for (const schema of ["task-brief", "worker-result", "verifier-result", "runtime-manifest"] as const) {
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
