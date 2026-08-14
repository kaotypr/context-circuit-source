# Legacy command reference

This document describes the previous command-oriented implementation. It is
retained temporarily during the workspace refactor and is not the user-facing
workflow.

The intended interface is an agent session that starts or resumes the workspace.
The agent may use deterministic helpers internally, but users should not need
to know their names or invoke them directly.

| Legacy capability | New workspace behavior |
| --- | --- |
| Configure workspace | Workspace-entry session inspects and records configuration |
| Import Product Knowledge | Root session gathers source evidence and proposes context |
| Create plan | Root session drafts a human-reviewable plan |
| Choose next work | Root session inspects active plans and session ownership |
| Run task | Session claims an approved plan and works in its assigned worktree |
| Review plan | Root or review subagent verifies the implementation and reports findings |
| Set status | Human gate changes remain explicit and authoritative |
| Publish or archive | Separate human-authorized actions |

The old command implementation, generated bundle, and command adapters are
legacy migration material. They may be removed after the filesystem-based
workspace workflow passes its acceptance scenarios.
