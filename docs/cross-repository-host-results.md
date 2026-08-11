# Cross-repository host compatibility results

Date: 2026-08-12

Codex desktop and Claude Code 2.1.220 each executed the same packaged-template
task in separate ignored wrapper copies. The API repository established and
independently verified `contractVersion: 1` before the frontend input unlocked;
the frontend then displayed `Contract v1` and received a separate verifier.

| Host | Run ID | Backend commit | Frontend commit | Final status |
| --- | --- | --- | --- | --- |
| Codex desktop | `20260811T174954Z-f93568dc` | `7495b184103603ad64dc93b874a18b5bb39b8918` | `f257c08fc3d43d37e4ec0811d0c30f8999622bef` | passed |
| Claude Code 2.1.220 | `20260811T174954Z-bfdd68dc` | `706c38c37531bd5978465d2ff8f796b6120b35ad` | `cfc77b4ec62ba905caa06886e8e26cf3f86ef5e2` | passed |

Both runs created separate backend and frontend branches/worktrees. An attempted
frontend start before contract verification was covered by automated rejection;
in the host runs, frontend remained `waiting` until the backend verifier result
was recorded. Each host ran repository tests and type/build checks, committed
only its repository scope, and used a different fresh verifier. Both manifests
contain six validated execution events and ended `passed`.

Claude Code initially wrote decorated `git log --oneline` text into the backend
worker `commits` array. The deterministic recorder rejected it. A fresh recovery
session reconstructed the result using the full object ID without changing
code, proving that malformed host output cannot advance the run.

The Claude frontend verifier also reported that its verifier input still showed
the contract approval as pending even though the authoritative backend verifier
had passed. The unlock transition was corrected to refresh both dependent
worker and verifier inputs, and an automated regression assertion now covers
the field. The verifier also recorded the legitimate drift risk of mirroring a
cross-repository contract constant without a shared generated package; this is
task-specific evidence, not hidden workflow state.

No host added remotes, pushed, opened pull requests, merged, deployed, mutated
an activity tool, or removed worktrees/runtime evidence. Preserved proof data is
ignored under `.runtime/cross-host-proof/20260812-release/`.
