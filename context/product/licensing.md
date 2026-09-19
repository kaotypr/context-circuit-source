# What each product is licensed under

Two things leave this checkout, and the people who receive them stand in
different relationships to them, so they carry different terms.

| What | License | Text |
| --- | --- | --- |
| This checkout and `context-circuit-cli` | Apache-2.0 | `context-circuit-source@LICENSE` |
| Everything a workspace receives | 0BSD | `context-circuit-source@product/LICENSE` |

The CLI is a binary an organization installs across its machines. Apache's
explicit patent grant is what carries that past a legal review without a
conversation, and its notice conventions already match a product shipping
`THIRD_PARTY_NOTICES.txt` in every archive.

The template is the side that decided the split. Those files are copied into
somebody else's repository and edited there, which makes them theirs in every
sense that matters. A notice-retention clause would oblige every workspace to
carry a notice of ours inside a repository that is not ours, forever, as a
condition of using a scaffold. 0BSD asks for nothing, which is the only terms a
scaffold can honestly impose on the thing it scaffolds.

## Where each text lives

The 0BSD text travels with the files it covers rather than being quoted at them.
The manifest installs it beside the shipped documentation inside a workspace, and
publication puts the same file at the root of the published template repository,
so the terms on display and the terms a workspace holds cannot disagree.

It is deliberately not installed at a workspace root. GitHub reads the root
`LICENSE` and labels the repository from it, and a workspace root belongs to
somebody else's project — a license there would label their work with ours.

The same reasoning governs the repository's other landing-page files. A
`CONTRIBUTING.md` or `SECURITY.md` about Context Circuit belongs on the published
repository and nowhere near a workspace, so they are held outside the manifest,
in a directory the published repository owns, and copied in after the release
artifact is extracted. Publication refuses a `README.md` among them, because the
published README is the product guide the manifest assembles.

## A binary is read where it was downloaded

The CLI's assets are published on the template repository, which is labeled 0BSD
because that repository is the template. The assets are not 0BSD, and a
repository label is the wrong thing to be answering that question — the
misreading it invites is the permissive one, where somebody redistributes an
Apache-2.0 binary without the notices it asks for.

Two things close that. Every archive carries its own `LICENSE` beside the binary,
which Apache-2.0 asks of a redistributor and which makes the archive answer for
itself wherever it was downloaded. The release notes say which license applies and
name the other product's. This is the same reasoning that keeps the asset filename
long while the tag is short: a release list supplies context and a download
directory does not.

The exact archive inventory is asserted by the release check, so a file added or
dropped there fails rather than shipping.

Owner:

- `context-circuit-source@LICENSE` — Apache-2.0, and what GitHub labels this
  checkout with.
- `context-circuit-source@product/LICENSE` — 0BSD, the terms every workspace
  receives alongside its shipped documentation.
- `context-circuit-source@scripts/release-manifest.txt` — where that text lands
  in a workspace, and the reason no license lands at its root.
- `context-circuit-source@release/template-repo/` — the landing-page files the
  published repository owns and no workspace receives.
- `context-circuit-source@scripts/publish-template.sh` — where those files are
  restored over the artifact and held out of the publish comparison.
- `context-circuit-source@scripts/build-cli.sh` — where each binary archive is
  given its own license beside the third-party notices.
