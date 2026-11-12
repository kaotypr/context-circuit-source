I'll run the QA pass for plan p1000. First I'm reading the QA skill and the workspace setup it depends on.The QA pass on p1000's worktree (`.worktrees/p1000/api`) finished with both checks passing. Nothing in the worktree was changed or committed.

## Commands run

- `go vet ./...` — passed
- `go test ./...` — passed (`ok example.com/api/internal/invoice`)

The before and after QA hooks ran cleanly. The after hook only asked for this command list.

Context Circuit plugin 0.5.0 is available (this session has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.