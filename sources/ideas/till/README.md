# Till — idea and prompt pack

Passive source material for a **real product** you can build with Context Circuit:
a one-person invoicing studio called **Till**.

Use this pack in an **instantiated product workspace** (a checkout of
`context-circuit-template` with Till as the project). Do not run the prompts in
this maintainer source repository — it is not a product workspace.

```
sources/ideas/till/
├── README.md              # this file
├── product-brief.md       # what Till is (name this file when asking for a design or intent)
├── ideas.md               # feature backlog that maps onto later changes
└── prompts/
    ├── README.md          # how to play the prompts; coverage of every case
    ├── 01-orientation.md
    ├── 02-intent-gate.md
    ├── 03-explore.md
    ├── 04-execute-and-verify.md
    ├── 05-candidate-delivery.md
    ├── 06-knowledge.md
    ├── 07-organization-and-publish.md
    └── 08-whole-flow.md
```

## How to start

1. Instantiate a blank Context Circuit workspace.
2. Copy this folder to that workspace at `sources/ideas/till/` (or name the files
   from here).
3. Open a conversation with the coordinator and play
   [prompts/README.md](prompts/README.md) in order — spine first, then branches.

When you ask the coordinator to read a source, **name the exact file**. This tree
stays passive until you do.

## Suggested repositories

Once you are ready for code, Till is two projects so multi-project cases have a
home:

| Folder | What it is |
| --- | --- |
| `till-api` | HTTP API: clients, invoices, payments |
| `till-web` | The studio owner's app |

You do not create them yourself in this pack. The orientation prompts ask the
coordinator to start them.
