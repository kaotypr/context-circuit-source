# Coding agent

Reads the instruction, runs the executable, and does the work a member asked for.

## [How workspace behavior is taught](../architecture/instruction-layers.md)

1. **As the agent, I want the rules that must always hold to be always
   loaded,** so that a prohibition does not depend on my having recognized the
   moment it applies.
2. **As the agent, I want each stage's procedure reachable by the words a
   person actually uses,** so that a skill arrives when its moment comes rather
   than being read once at session start.
3. **As the agent, I want everything shipped to be named by something I load,**
   so that no instruction reaches a workspace that nothing will ever open.

## [Command surface and output contract](../architecture/command-surface.md)

1. **As the agent, I want structured output on request,** so that I read values
   instead of parsing prose I wrote the format of.
2. **As the agent, I want a refusal to say what it needs,** so that a failure is
   an instruction rather than a dead end I work around.
3. **As the agent, I want the help text to carry what a command refuses,** so
   that I can tell what is unsupported from what I have asked for wrongly.

## [Subagent roles and host settings](../domains/subagent-roles.md)

1. **As the agent, I want a brief that carries its own record in full,** so that
   a delegate starts from what it was handed rather than searching a repository
   to reconstruct it.
2. **As the agent, I want to integrate from what a worker reported,** so that
   delegation stays cheaper than doing the work myself.
3. **As the agent, I want a written role definition reported as written rather
   than loaded,** so that I do not claim a capability the host has not
   registered.

## [Authorization boundaries](../domains/authorization.md)

1. **As the agent, I want to know exactly which act a person's word
   authorized,** so that I neither stop for permission already given nor take a
   step nobody asked for.
2. **As the agent, I want a lane for a change whose outcome is already
   settled,** so that a person choosing to skip the ceremony gets the change
   rather than the process.

Owner:

- `context-circuit-source@product/AGENTS.md.in` — what is always in force.
- `context-circuit-source@product/skills/` — the procedure behind each stage.
- `context-circuit-source@internal/workspace/agents.go` — the briefs and role
  definitions this actor composes and launches.
