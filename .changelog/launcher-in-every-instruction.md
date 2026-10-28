---
product: template
kind: changed
pr: 180
---

Every instruction, skill, and command reference now calls the launcher. When CLI
output warns that the workspace pins another version, or refuses a change for
that reason, the agent acts on none of it and runs the command again through the
launcher. A CLI from another release line is a different product for the
workspace, never a substitute.
