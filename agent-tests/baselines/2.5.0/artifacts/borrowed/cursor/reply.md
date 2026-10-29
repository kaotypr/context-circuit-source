The handbook still treats invoice numbers as global, while they restart each year. I’ll read how this workspace records knowledge and update that fact so later work uses the yearly rule.Invoice numbering is recorded in the shared handbook, which this workspace only reads. I’ll locate that note and confirm what it says before deciding where a correction belongs.The yearly restart is not recorded here, and it should not be copied into this workspace.

Invoice numbering lives in the borrowed handbook, which this workspace only reads. `context find` marks that entry as borrowed:

- `handbook/numbering.md` — "Invoice numbers are global and never restart."
- Catalog line in `handbook/INDEX.md` — "invoice numbers are global and never restart."

This workspace's own notes under `context/` do not describe numbering. Copying the yearly restart into `context/` would leave a second copy that goes stale while still looking current, and nothing here can move the handbook note.

The correction belongs upstream, through the handbook repository's own review: change `numbering.md` and its index line so they say numbers restart each year. This checkout cannot push that change. I have not edited `context/` or the handbook.