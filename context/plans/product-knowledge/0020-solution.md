# Solution

## Proposed solution

- Define the Product Knowledge content conventions and minimal metadata in canonical contracts and host-neutral skills, using stable file paths and Git revisions instead of a large hand-maintained knowledge graph.
- Use context/PROJECT.md as the entry point, context/roles/README.md as the role index, context/roles/<role>.md for one cross-domain role story per business role, and context/domains/<domain>/README.md plus workflows/<workflow>.md for detailed behavior.
- Add neutral templates and initialization inputs that create only the known baseline; retain unknowns explicitly and allow later work to grow coverage just in time.
- Teach gather-context to traverse PROJECT.md to a relevant role or domain, load only referenced workflow pages, compare them with requested sources and repository evidence, and emit a bounded context package plus candidate additions or contradictions under ignored runtime state.
- Extend plan contracts and generated plan documents with Product Knowledge references, impact classification, and a concise proposed delta while keeping canonical workflow pages unchanged during planning.
- Carry the resolved context package and declared delta into scoped worker and independent verifier inputs, and validate their Product Knowledge impact reports alongside existing Git and acceptance evidence.
- Extend closeout contributions with Product Knowledge impact and add a reviewed synchronization workflow that updates only referenced current-behavior pages after explicit effectiveness confirmation.
- Generate onboarding packs by selecting the product map, one or more role files, relevant domain summaries and workflows, architecture context, and known gaps; stamp every pack with its source revision and included paths.
- Update the deterministic bundle, Codex and Claude adapters, human documentation, neutral template build, fixtures, and tests without adding a database, provider SDK, background service, or automatic external write.
