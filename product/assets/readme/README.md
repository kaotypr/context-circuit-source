# Context Circuit README assets

These assets apply the Kaotypr Design Language System to Context Circuit. They
use product concepts that exist today and do not depict a fictional graphical
interface.

| Asset | Size | Intended use |
| --- | --- | --- |
| `context-circuit-mark.svg` / `.webp` | 512 × 512 | Repository avatar, favicon source, square brand mark |
| `context-circuit-logo.webp` | 520 × 160 | README brand lockup or documentation header |
| `knowledge-circuit.webp` | 1400 × 820 | Primary product explanation |
| `social-preview.png` | 1280 × 640 | Uploadable GitHub social preview; header of the source README |

The workspace release ships only the assets a workspace README renders —
`context-circuit-logo.webp` and `knowledge-circuit.webp` — mapping them to
`.context-circuit/assets/readme/`. A generated workspace README uses that path.

The rest are source-side: an avatar, a favicon source, and a social preview are
uploaded to a repository or a site rather than rendered in a workspace, so no
workspace carries them. This checkout links to `product/assets/readme/`
directly. The social preview stays PNG because GitHub's social-preview upload
accepts PNG, JPG and GIF but not WebP, and nothing here can regenerate that
composed image; the square mark has no such copy, since its `.svg` rasterizes
on demand for an avatar upload.

## DLS application

- Schoolbell is the display face for the wordmark and loud moments.
- Fredoka carries headings, labels, and body copy.
- Fira Code is reserved for commands, filenames, and record IDs.
- Near-black background, card, muted, foreground, and border values follow the
  canonical dark theme hierarchy.
- Lavender is Context Circuit's secondary brand accent. It marks coordination
  and selected emphasis; it is not a status color.
- Status hues remain reserved for destructive, success, warning, and info.
- Labels use sentence case. All caps are reserved for proper wordmarks and
  conventional unit labels.
- Surfaces use semantic card and chip relationships rather than arbitrary radii.
- Icons use regular line weight, with solid treatment only where emphasis is
  necessary.
- There are no glow effects, arbitrary gradients, or opacity-based text colors.

## Public-facing language

- Lead with the outcome a developer or team receives, not the internal record or
  storage model that produces it.
- Prefer ordinary terms such as project context, goal, plan, working copy, and
  verified result.
- Use **project context** for the complete coordinated picture: product
  knowledge, repositories, goals, plans, dependencies, and current work. Use
  **product knowledge** only for durable architecture, domain rules, decisions,
  conventions, and vocabulary stored under `context/`.
- Do not use Context Circuit-specific terms such as intent, durable knowledge,
  worktree association, host role, or reconcile without explaining them in the
  same sentence.
- Name people, AI agents, and sub-agents directly instead of hiding them behind
  implementation terms such as member or host.
- A reader should understand each asset before learning Context Circuit's file
  layout or command vocabulary.

The mark SVG retains semantic `title`, `desc`, and `aria-labelledby` attributes
for scalable brand use. Composed README assets are WebP wherever they are only
rendered, so their layout and Kaotypr font roles remain stable on every host at
a fraction of the weight a workspace would otherwise carry.
