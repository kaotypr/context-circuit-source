# roles-and-expertise

Continues [design.md](./design.md) §"The central decision — no new role". This file
is the full argument for **why v0.7.0 adds no role**, and where the frontend
competence actually lives instead.

## Roles encode authority, not domain

Every Context Circuit role is defined by its position in the verification chain:

| Role | May write? | May verify its own work? | Independence |
| --- | --- | --- | --- |
| Coordinator | no | n/a | the interface; holds no write or verify authority |
| Worker | yes, bounded | **no** | executes one bounded brief |
| Independent verifier | no (read-only) | verifies *others'* work | independent of the worker |

Not one of these is defined by a domain. A worker editing a Rust service and a
worker editing a React component are the *same role* — the boundary is "bounded
write, no self-verify," not "knows Rust" or "knows React." Domain is a property of
the **work**, delivered as **context**; it is never a property of the **role**.

## Why "make the root session a frontend expert" fails

The request — a role that turns the **root session** into an expert frontend
engineer — fails on the authority axis two different ways:

1. **Root-becomes-writer.** The root session is the coordinator. Granting it
   frontend "expertise" that lets it edit code gives the coordinator write
   authority. That collapses coordinator/worker separation and removes the point of
   having an independent verifier at all — the interface would be writing and
   effectively self-approving. This is the invariant the product exists to protect.
2. **A new writing identity per domain.** Suppose instead a distinct "Frontend
   Expert" writer identity. It adds no new authority distinction (it is just a
   worker with a domain hat), and it opens an unbounded series — Backend Expert,
   Data Expert, Infra Expert — each duplicating the worker role with a different
   sticker. The role set should enumerate *authority positions*, which are few and
   fixed, not *domains*, which are infinite.

Either way, "expert" is being modeled at the wrong layer.

## Where the competence actually attaches

Frontend competence is real and necessary — it just attaches as **loaded context on
existing roles**, split by how universal it is:

- **The worker becomes the expert, via its brief.** Under `cc-refine` the worker is
  grounded (through the existing repository-grounding path, INV-GROUND) in the
  target repository's *own* frontend context — component library, design tokens,
  dev-server command, breakpoints, visual-regression tooling — and the `cc-refine`
  skill supplies the *universal* UI-engineering procedure (responsive thinking,
  a11y defaults, layout discipline, how to read taste language into concrete
  edits). The assembled brief is where "expert frontend engineer" literally lives.
- **The coordinator gains fluency, not authority.** From the same skill (read as
  procedure), the coordinator gains enough frontend vocabulary to host the
  conversation — interpret "make it pop / calmer / tighter," know which breakpoints
  to check, narrate what changed — while writing nothing. It still delegates every
  edit to the worker.

So the human's experience — "I'm pairing with a senior frontend engineer" — is
produced by two existing roles wearing loaded competence, not by a new agent.

## This is INV-SKILL-01, restated

A skill grants no route, role, or authority; it is read-as-procedure resolved by
path. `cc-refine` is precisely that: it makes the coordinator and worker *competent*
at frontend refinement and defines the *mode*, but it moves no authority. The
worker still cannot self-verify; the coordinator still cannot write; the verifier is
still independent and read-only. Competence up, authority unchanged.

## Boundary restated for hosts

Native "expert persona" or "custom agent" features on a host (Claude Code
subagents, Codex/Cursor equivalents) map only to the existing worker or independent
verifier packets carrying this loaded context. A host persona never becomes a
Context Circuit role and never authorizes a route; host identity stays in
`host_evidence`.
