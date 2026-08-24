---
name: cc-test-case
description: >
  Maintainer-only, source-only. Drives one human-simulated template test case
  IN-SESSION: prepare an isolated workspace, spawn the product coordinator as a
  sub-agent (which may spawn its own worker/verifier for full-execution cases),
  play the case's natural turns, capture the transcript, and grade
  deterministically. Use for full-execution cases (05/06/08/09) where nested
  sub-agents are needed; the shell `run-scenario.sh --live` driver remains for
  conversation-only cases that need a machine-readable file-access trace.
  Invoke as `/cc-test-case <case-id>` (e.g. `/cc-test-case 05-approve-and-execute`).
---

# Driving a template test case in-session

This skill is the **in-session Task driver**. It complements the shell driver:
`run-scenario.sh` still does the deterministic prep and `grade.sh` still grades;
this skill only replaces the *conversation-driving* step, using real nested
sub-agents instead of a headless `claude -p`.

**Boundaries (do not violate):**
- Never let the coordinator sub-agent see the case's `grader:` block or learn it
  is under test. It receives only the isolated workspace and your natural prompts.
- The coordinator operates ONLY inside the instantiated workspace directory; pin
  it there and never let it touch the source checkout.
- Grading is `grade.sh` (mechanical). Never grade with an LLM.
- Dimension C degrades to warning-only in this driver (no machine-readable
  file-access trace); that is expected and allowed (host-matrix §5).

## Procedure

1. **Prepare** the run (deterministic; no live model):
   `sh test/template-runtime/human/run-scenario.sh --host claude-code <case-id>`
   Capture the printed run directory (`prepared run: <RUN>`), and read
   `<RUN>/run.yaml` for `workspace`, `case_file`, `transcript`, `mode`.

2. **Read only the `human:` block** of `<case_file>` — persona, `turns`,
   `reactions`. Do NOT read or use the `grader:` block for anything you send the
   coordinator.

3. **Spawn the coordinator** as a sub-agent (Agent tool, general-purpose) with a
   prompt shaped like:
   > You are the assistant for the software project in `<ABS_WORKSPACE>`. First
   > read `<ABS_WORKSPACE>/AGENTS.md`, `CLAUDE.md`, and `WORKFLOW.md` and follow
   > them exactly. Operate ONLY inside `<ABS_WORKSPACE>` — use absolute paths and
   > `cd "<ABS_WORKSPACE>"` for any shell. A user is talking to you; reply to them
   > normally. First message: "<turn 1 say text>"

   Do not tell it it is a test, a simulation, or being graded.

4. **Play the turns** in order, continuing the same coordinator via SendMessage
   (its context stays intact):
   - `say:` — send the text verbatim.
   - `on_offer_plan:` — send only if the last reply offered/《asked about》a plan.
   - `on_open_questions: answer_plainly` — if the last reply asked a question,
     answer plainly in lay terms without inventing repo/tech facts the persona
     would not know; otherwise skip.
   - In `full-execution` mode, when the human approves+executes, let the
     coordinator run the real flow: it spawns its own worker and independent
     verifier sub-agents (nesting is supported) and drives the workspace engine.
     Wait for each reply.

5. **Capture the transcript** to `<RUN>/transcript.txt` as you go, one block per
   turn:
   ```
   human: <text>
   coordinator: <full reply text>
   ```
   (Append with the Write/Bash tools; start the file with `# transcript: <id>`.)

6. **Conversational verdict (optional but preferred):** spawn the
   `cc-human-simulator` agent with ONLY the `human:` block and the finished
   transcript; save its verdict to `<RUN>/conversational-verdict.txt`.

7. **Grade:** `sh test/template-runtime/human/grade.sh <RUN>` and report the
   result. C will show its forbidden checks (hard) and warn that required-reads /
   the trace are unavailable in this driver.

8. **Report** to the maintainer: transcript highlights, the grader verdict, any
   isolation warning in `run.yaml`, and whether `~/…` outside the workspace was
   touched. Leave the run under `.out/` (git-ignored, disposable).

## Notes

- If a sub-agent cannot be spawned (nesting/host limit), stop and report
  `host-blocked` for that case — never self-drive the worker/verifier.
- Keep the coordinator's cwd pinned to the workspace; if you see it reading the
  source checkout, that is an isolation finding to report, not to ignore.
