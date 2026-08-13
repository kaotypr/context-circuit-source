Read and follow the canonical `.agents/skills/w-run-task/SKILL.md`.

After the deterministic preparation command succeeds, launch a fresh Claude Code
worker with only the emitted worker-input JSON and its referenced instruction
files. Launch a separate fresh, read-only verifier with only verifier-input JSON.
Do not reproduce workflow logic in this command or inherit the coordinator's
full conversation.

For non-interactive `claude -p` proofs, place the positional prompt before the
variadic `--add-dir` option, use `--no-session-persistence`, and grant the wrapper
directory only so the session can read contracts and write its designated result.
