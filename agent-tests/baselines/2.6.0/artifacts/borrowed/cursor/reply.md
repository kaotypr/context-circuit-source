I’ll update the project knowledge so invoice numbers are recorded as restarting each year, and I’ll follow the workspace rules for how that knowledge is written.The handbook is mounted as borrowed knowledge, so I’ll confirm what it says about invoice numbers before changing anything here.The handbook entry on invoice numbers is borrowed knowledge, so it stays unchanged in this workspace.

`context find` for invoice numbers returns one match, marked borrowed, from `handbook/INDEX.md`:

- [Invoice numbers](numbering.md) — invoice numbers are global and never restart

A correction belongs in the handbook repository, through that repository's own review. Copying the new yearly restart rule into `context/` would leave a second copy that still reads as current after the handbook moves on. This workspace's invoices note covers voiding, credit notes, and the plain-text summary, and it does not state a numbering rule.

The Context Circuit plugin here is 0.5.0 (installed 0.2.2). Update it with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.