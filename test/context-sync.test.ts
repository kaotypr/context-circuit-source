import assert from "node:assert/strict";
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { prepareContextReview, prepareContextSync } from "../scripts/lib/context-sync.js";
import { git } from "../scripts/lib/git.js";
import type { ContextSyncRequest } from "../scripts/lib/types.js";

const contribution = "contributions/general/20260812T030000Z-kao-reset-contract.md";

async function workspace(): Promise<{ root: string; cleanup: () => Promise<void> }> {
  const root = await mkdtemp(join(tmpdir(), "kao-context-sync-"));
  await mkdir(join(root, "context"), { recursive: true });
  await mkdir(join(root, "contributions", "general"), { recursive: true });
  await writeFile(join(root, ".gitignore"), ".runtime/\n", "utf8");
  await writeFile(join(root, "README.md"), "# Wrapper\n", "utf8");
  for (const name of ["PROJECT", "ARCHITECTURE", "CONVENTIONS", "DECISIONS"]) await writeFile(join(root, "context", `${name}.md`), `# ${name}\n`, "utf8");
  await writeFile(join(root, contribution), `# ADHOC-1: Reset contract

- Run: \`20260812T030000Z-1234abcd\`

## Outcome

Merged.

## Affected repositories

- \`backend\` on branch \`agent/reset\`.

## Pull requests and commits

- Commit recorded.

## Verification

- Passed.

## Decisions and deviations

- Reset responses use state.count.

## Remaining risks and follow-up

- None.

## Candidate durable learnings

- The shared reset response is state.count.
`, "utf8");
  await writeFile(join(root, "workspace.yaml"), `version: 1
template_version: 0.2.1
workspace: { name: sync-test, mode: team, default_branch: main }
repositories:
  backend: { path: repositories/backend, mode: ignored-clone, role: application-api, agent: backend, default_branch: main }
activity: { provider: none, access: auto, required_capabilities: [], optional_capabilities: [] }
workflow: { human_gates: [plan-approval, task-selection, merge], maximum_repair_attempts: 2, wrapper_change_policy: pull-request }
`, "utf8");
  await git(root, ["init", "--initial-branch=main"]);
  await git(root, ["add", "."]);
  await git(root, ["-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "fixture"]);
  return { root, cleanup: async () => rm(root, { recursive: true, force: true }) };
}

function request(): ContextSyncRequest {
  return {
    contract_version: 1,
    author: "kao",
    contributions: [contribution],
    proposals: [
      { source_contribution: contribution, classification: "durable-wrapper", summary: "Record the shared reset shape.", proposed_change: "Document state.count.", target: "context/ARCHITECTURE.md" },
      { source_contribution: contribution, classification: "repository-local", summary: "Keep handler tests beside the handler.", target_repository: "backend" },
      { source_contribution: contribution, classification: "future-task", summary: "Consider versioning the response." },
    ],
  };
}

test("prepares a scoped wrapper worktree and review handoff", async (t) => {
  const fixture = await workspace();
  t.after(fixture.cleanup);
  const record = await prepareContextSync({ workspaceRoot: fixture.root, request: request(), now: new Date("2026-08-12T04:00:00Z") });
  assert.equal(record.wrapper_mode, "pull-request");
  assert.deepEqual(record.allowed_wrapper_paths, ["context/ARCHITECTURE.md"]);
  assert.equal(record.repository_follow_ups[0]?.repository, "backend");
  const architecture = join(record.worktree, "context", "ARCHITECTURE.md");
  await writeFile(architecture, `${await readFile(architecture, "utf8")}\nReset uses state.count. Source: ${contribution}\n`, "utf8");
  await git(record.worktree, ["add", "context/ARCHITECTURE.md"]);
  await git(record.worktree, ["-c", "user.name=Curator", "-c", "user.email=curator@example.invalid", "commit", "-m", "docs: sync reset contract"]);
  const review = await prepareContextReview({ workspaceRoot: fixture.root, syncId: record.sync_id, now: new Date("2026-08-12T04:05:00Z") });
  assert.equal(review.status, "blocked");
  assert.deepEqual(review.changed_files, ["context/ARCHITECTURE.md"]);
  assert.match(review.blockers[0]!, /no origin remote/);
});

test("rejects out-of-scope wrapper edits", async (t) => {
  const fixture = await workspace();
  t.after(fixture.cleanup);
  const record = await prepareContextSync({ workspaceRoot: fixture.root, request: request(), now: new Date("2026-08-12T05:00:00Z") });
  await writeFile(join(record.worktree, "README.md"), "# Changed outside scope\n", "utf8");
  await git(record.worktree, ["add", "README.md"]);
  await git(record.worktree, ["-c", "user.name=Curator", "-c", "user.email=curator@example.invalid", "commit", "-m", "docs: unsafe change"]);
  await assert.rejects(prepareContextReview({ workspaceRoot: fixture.root, syncId: record.sync_id }), /outside approved wrapper scope/);
});

test("dirty wrapper protection preserves existing work", async (t) => {
  const fixture = await workspace();
  t.after(fixture.cleanup);
  const path = join(fixture.root, "unrecorded.txt");
  await writeFile(path, "preserve\n", "utf8");
  await assert.rejects(prepareContextSync({ workspaceRoot: fixture.root, request: request() }), /unresolved local changes/);
  assert.equal(await readFile(path, "utf8"), "preserve\n");
});

test("review recomputes allowed scope instead of trusting a tampered runtime record", async (t) => {
  const fixture = await workspace();
  t.after(fixture.cleanup);
  const record = await prepareContextSync({ workspaceRoot: fixture.root, request: request(), now: new Date("2026-08-12T06:00:00Z") });
  const recordPath = join(fixture.root, ".runtime", "context-sync", record.sync_id, "record.json");
  const tampered = JSON.parse(await readFile(recordPath, "utf8"));
  tampered.allowed_wrapper_paths.push("README.md");
  await writeFile(recordPath, `${JSON.stringify(tampered, null, 2)}\n`, "utf8");
  await assert.rejects(prepareContextReview({ workspaceRoot: fixture.root, syncId: record.sync_id }), /allowed paths do not match/);
});

test("review requires the canonical update to cite its contribution", async (t) => {
  const fixture = await workspace();
  t.after(fixture.cleanup);
  const record = await prepareContextSync({ workspaceRoot: fixture.root, request: request(), now: new Date("2026-08-12T07:00:00Z") });
  await writeFile(join(record.worktree, "context", "ARCHITECTURE.md"), "# ARCHITECTURE\n\nReset uses state.count.\n", "utf8");
  await git(record.worktree, ["add", "context/ARCHITECTURE.md"]);
  await git(record.worktree, ["-c", "user.name=Curator", "-c", "user.email=curator@example.invalid", "commit", "-m", "docs: omit evidence"]);
  await assert.rejects(prepareContextReview({ workspaceRoot: fixture.root, syncId: record.sync_id }), /must cite source contribution/);
});
