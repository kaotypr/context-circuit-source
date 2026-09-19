# What each product is licensed under

Two things leave this checkout, and the people who receive them stand in
different relationships to them, so they carry different terms.

| What | License | Text |
| --- | --- | --- |
| This checkout and `context-circuit-cli` | Apache-2.0 | `context-circuit-source@LICENSE` |
| Everything a workspace receives | 0BSD | `context-circuit-source@product/LICENSE` |

A workspace holds no copy of that text, which the next section explains.

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

The 0BSD text appears once: at the root of the published template repository,
where GitHub reads it and labels that repository. Nothing installs it into a
workspace.

A copy shipped inside every workspace until 2.0.0-rc.14, on the reasoning that
terms should travel with the files they cover. They should, when the terms ask
something. 0BSD asks nothing, so that copy answered a question nobody had while
adding a file every workspace had to explain. A license at a workspace root is
worse than useless — GitHub labels the repository from it, and a workspace root
belongs to somebody else's project.

The same reasoning governs the repository's other landing-page files. A
`CONTRIBUTING.md` or `SECURITY.md` about Context Circuit belongs on the published
repository and nowhere near a workspace, so they are held outside the manifest,
in a directory the published repository owns, and copied in after the release
artifact is extracted. The guide and its artwork are restored the same way from
`product/`, which is why publication refuses a `README.md` or `LICENSE` defined
among the landing-page files: each of those has exactly one source.

What makes the arrangement hold is `.gitattributes`, restored with them and
marking every repository-owned path `export-ignore`. A project is created from a
tag with `git archive`, which honors that mark, so the files stay in the tree
for GitHub and leave it for the workspace. A later template update diffs two
tags instead, sees the tree rather than an archive, and cannot read the mark; the
shipped `AGENTS.md` names the same paths so the agent doing that update knows
not to apply them.

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
- `context-circuit-source@product/LICENSE` — 0BSD, the terms everything a
  workspace receives is offered under, shown at the template repository's root.
- `context-circuit-source@scripts/release-manifest.txt` — everything a workspace
  receives, and the reason no license is among it.
- `context-circuit-source@release/template-repo/` — the landing-page files the
  published repository owns, and the `.gitattributes` that keeps every
  repository-owned path out of the archive a project is created from.
- `context-circuit-source@scripts/publish-template.sh` — where those files are
  restored over the artifact and held out of the publish comparison.
- `context-circuit-source@scripts/build-cli.sh` — where each binary archive is
  given its own license beside the third-party notices.
