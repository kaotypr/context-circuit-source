# Context Circuit source, template, and universal project workspace

This document defines the relationship between the current source repository and
the product template that users initialize for their projects.

## 1. Three distinct things

Context Circuit v0.5 has three distinct identities:

~~~text
Cc source: context-circuit-source
  → maintainer source repository
  → self-hosts Cc while developing and testing the product

Cc template: context-circuit-template
  → distributable product repository/template
  → contains the universal project workspace

instantiated project workspace
  → a configured copy or installation of the template
  → holds one project's knowledge, plans, repositories, and runtime evidence
~~~

The source repository is not itself the universal project workspace, even when
it self-hosts Context Circuit files for development. The template is the
product artifact. An instantiated workspace is the user's project environment.

“Cc” is the short prompt form for Context Circuit. “Cc source” and “Cc
template” are the corresponding short forms for the two repositories defined
here.

## 2. context-circuit-source

context-circuit-source is the current repository maintained by the Context
Circuit team. It contains:

- source design and terminology;
- product implementation;
- the product template seed;
- contracts and runtime code;
- host adapters and role guidance;
- semantic, security, release, and migration tests;
- release assembly and maintainer instructions;
- source-only evidence about design and implementation.

The source repository may contain its own Product Knowledge because it is a
project maintained with Context Circuit. That self-hosted project knowledge is
not automatically copied into context-circuit-template or any customer
workspace.

Source-only directories and maintainer evidence must not leak into the product
template or an instantiated workspace.

### 2.1 The source repository self-hosts Cc

The source repository has nearly the same workspace-facing structure as the
product template because the Cc team uses Cc to understand and improve Cc. It
has its own Product Knowledge, plans, repository mappings, and `.runtime/` so
the team can create, approve, execute, verify, repair, and preserve plans for
the source project itself.

The source-side workspace is a dogfooding environment, not customer Product
Knowledge. It is allowed to describe the Context Circuit product in detail;
that knowledge remains source-owned unless a generic part is deliberately
turned into a template file or product contract.

The intended source shape is:

~~~text
context-circuit-source/
├── context/                         # Cc source Product Knowledge
├── plans/                           # Cc source plans
├── .runtime/                        # persistent source-side execution state
├── template/                        # blank/product workspace seed
├── wrapper/                         # shipped product layer (runtime, contracts, adapters)
├── test/                            # source-only acceptance and test suites
├── template-harness/                # built-template behavior laboratory
├── sources/system-design/context-circuit/v0.5/
│                                     # canonical maintainer design
└── maintainer files and release assembly
~~~

The `context/`, `plans/`, and `.runtime/` surfaces follow the same semantics as
the corresponding surfaces in an initialized Cc template. The source-only
`template/`, implementation paths, tests, design files, and release machinery
are additional maintainer surfaces. They are not copied merely because the
source repository has them.

The source `.runtime/` is real workspace runtime state. It may contain active
source plans, branches, worktrees, commits, verifier results, repair attempts,
and recovery evidence. It must be private to the source checkout, preserved
across normal work, and excluded from template assembly.

### 2.2 Testing the built Cc template with realistic work

`template-harness/` is a source-only behavior laboratory for the built
Cc template. It answers a different question from unit tests of source
functions: does the assembled product provide the intended experience when an
agent receives realistic project requests?

Each scenario may define:

- a natural-language project prompt or conversation;
- a small fixture project with one or more Git repositories;
- starting Product Knowledge and source evidence;
- expected plan, repository mapping, approval, execution, verification, and
  repair outcomes;
- expected runtime evidence and bounded failure behavior.

The harness must assemble or select the same `context-circuit-template` artifact
that would be distributed, initialize an isolated temporary project workspace,
connect fixture repositories, and exercise the product through its normal
conversation and runtime boundaries. It must not accidentally import the
source repository's Product Knowledge, plans, `.runtime/`, or implementation
state as if they were installed product files.

The template-harness laboratory is source-owned and is not part of
`context-circuit-template`. Generated workspaces, fixture worktrees, commits,
and runtime records are disposable test state and must live outside the
source-side `.runtime/`. The scenarios should include, at minimum, context
gathering, detailed plan creation, `review plan`, approval, multi-repository
execution, independent verification, worker repair, the three-failure limit,
explicit human completion, archive or restore without plan-status validation,
binding the workspace root as the optional `workspace` repository, cloning into
the gitignored `repositories/` directory, initializing an empty repository
with `git init` and an initial commit, setting personal or team anchor branches,
and the separate pull-request action that uses execution branches as sources
and anchor branches as default targets.

## 3. context-circuit-template

context-circuit-template is the product repository/template. It is the universal
project workspace that can be initialized for any project.

It provides:

- small workspace entry instructions;
- workspace identity files;
- Product Knowledge structure and retrieval catalog;
- source intake and provenance structure;
- readable plan and task templates;
- repository registration and local binding support;
- the Context Circuit runtime;
- worker and independent verifier role guidance;
- upgrade, recovery, and delivery boundaries;
- tests or fixtures needed by the installed product where appropriate.

The template contains generic project structure, not the Product Knowledge of one
customer or project. Project-specific knowledge, plans, repository bindings,
runtime records, and work are added after initialization.

The template's `.runtime/` is an installed workspace's private execution
surface. The source repository's `.runtime/` is the source project's own
execution surface. They have the same runtime semantics but never share state.

## 4. Instantiated universal project workspace

A project workspace is created from context-circuit-template and then initialized
for a specific project.

Initialization creates or validates:

~~~text
README.md
.gitignore
AGENTS.md
WORKFLOW.md
workspace.yaml
context/
sources/
plans/
  INDEX.md
  .archived/
repositories/                         # ignored local project-repository checkouts
.runtime/
repositories.local.yaml
~~~

The initialized workspace then gains:

- project identity and purpose;
- connected repository identities and local bindings;
- Product Knowledge units and the context index;
- source artifacts and provenance;
- draft, approved, and completed plans;
- an active plan index and separately stored archived plans;
- execution worktrees, commits, handoffs, and verifier evidence.

The workspace is universal because its structure and behavior apply to many
project types. Its knowledge and plans are specific to the project.

The template-provided root `.gitignore` must exclude at least:

~~~gitignore
/repositories/
/repositories.local.yaml
/.runtime/
~~~

This keeps connected repository source, local anchor settings, and private
execution state out of the workspace repository.

## 5. What belongs where

| Material | context-circuit-source | context-circuit-template | Instantiated project workspace |
| --- | --- | --- | --- |
| Source design | Canonical | Excluded | Excluded |
| Product implementation | Canonical | Released implementation | Installed product files |
| Blank workspace seed | Maintained | Included | Initial files |
| Product Knowledge for Context Circuit itself | May exist for self-hosting | Generic only | Project-specific |
| Project Product Knowledge | Source project only | Not included | Canonical |
| Plan templates | Maintained | Included | Used to create plans |
| Customer/project plans | Maintainer plans only | Not included | Canonical |
| Local repository paths | Source-local only | Never portable | Ignored local binding |
| Connected repository checkouts | Source-local only | Never shipped | Ignored under `repositories/` by default, or an explicitly bound local path |
| Runtime evidence | Source session only | Never shipped | Project execution evidence |
| Template-runtime scenarios and fixtures | Source-owned | Excluded | Never copied |
| Credentials | Never stored | Never shipped | Never stored |

## 6. Template assembly boundary

Product assembly copies only template-owned files and generic product behavior.
It excludes:

- context-circuit-source design files;
- maintainer plans and implementation logs;
- source-only tests and evidence;
- source repository runtime state;
- `template-harness/` scenarios, fixtures, generated workspaces, and
  evidence;
- credentials and local bindings;
- connected repositories;
- project Product Knowledge;
- customer plans and execution records;
- archived plan contents and any archive-only organization state.

The assembled template must identify itself as the Context Circuit product
template and as an uninitialized universal project workspace.

## 7. Initialization and upgrade

Initialization adds project identity and repository registrations without
inventing product knowledge. The user or agent may then gather context, connect
repositories, and create plans.

The initialized workspace ignores `repositories/` and
`repositories.local.yaml`. `repositories/` is the default local destination
for cloned or newly initialized project repositories; the workspace root, when
it is a Git repository, is bound separately as the logical `workspace`
repository.

An upgrade may replace template-owned files, runtime code, contracts, and role
guidance. It must preserve workspace-owned files:

- project identity;
- accepted Product Knowledge;
- sources and provenance;
- plans and plan status;
- local repository bindings;
- connected repositories;
- runtime evidence and active worktrees.

If a template change changes the meaning of a plan, context unit, or runtime
record, the upgrade must report migration-needed and preserve the old state.

## 8. Self-hosting

context-circuit-source may use its own template and runtime while developing
Context Circuit. Self-hosting does not collapse the identities:

- source design remains maintainer design;
- the source repository remains the source repository;
- source Product Knowledge remains source-project knowledge;
- source `.runtime/` records remain source-side execution evidence;
- a generated template remains the distributable product;
- template-runtime scenarios test the generated template in isolated fixture
  workspaces;
- test and release artifacts remain source-owned.

A source change is not a product-template change until release assembly includes
it and the template acceptance checks pass.

## 9. Naming requirements

New user-facing documentation should use:

- Context Circuit source repository or context-circuit-source;
- Context Circuit product template or context-circuit-template;
- universal project workspace;
- instantiated project workspace;
- Context Circuit runtime.

Wrapper is an accepted synonym for the product — the universal project
workspace. The directory wrapper/ holds the product's shipped layer; write the
path explicitly when referring to that directory rather than the product.
