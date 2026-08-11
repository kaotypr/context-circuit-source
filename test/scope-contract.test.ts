import assert from "node:assert/strict";
import test from "node:test";
import { normalizeDirectRequest } from "../scripts/lib/run-task.js";
import { validateContract } from "../scripts/lib/validation.js";

const base = {
  request: "Change counter behavior",
  repository: "frontend",
  acceptanceCriteria: ["The behavior is correct"],
  scope: ["src/App.tsx"],
  verificationCommands: ["npm test"],
  workId: "ADHOC-20260811-001",
  runId: "20260811T083000Z-abcd1234",
  createdAt: "2026-08-11T08:30:00.000Z",
};

test("defaults omitted test scope to explicit verifier-only evidence", async () => {
  const brief = normalizeDirectRequest(base);
  assert.deepEqual(brief.scope, ["src/App.tsx"]);
  assert.deepEqual(brief.test_expectation, {
    policy: "verifier-only",
    paths: [],
    rationale: "No test edit scope is authorized; the verifier must supply independent acceptance evidence.",
  });
  assert.deepEqual(await validateContract("task-brief", brief), []);
});

test("existing coverage records test paths without authorizing test edits", async () => {
  const brief = normalizeDirectRequest({
    ...base,
    testScope: ["src/App.test.tsx"],
    testPolicy: "existing-coverage",
  });
  assert.deepEqual(brief.scope, ["src/App.tsx"]);
  assert.deepEqual(brief.test_expectation?.paths, ["src/App.test.tsx"]);
  assert.equal(brief.test_expectation?.policy, "existing-coverage");
});

test("rejects scope paths that escape the repository", () => {
  assert.throws(
    () => normalizeDirectRequest({ ...base, scope: ["../outside.ts"] }),
    /repository-relative paths/,
  );
});
