---
product: both
kind: added
---

Extensions can now add agent roles and add to the shipped ones. A manifest's `roles` declare a role with its description, `read-only` or `edit` access, the stages it serves, the shipped role it is `like` for model and effort, its instructions, and the headings it returns. `agent setup` installs each role of an extension on for you beside the shipped four, as `cc-x-<extension>-<role>`, and `agent dispatch --role x-<extension>-<role>` briefs it like a shipped role, with limits no instruction lifts. `role_additions` add instructions and return headings to the explorer, planner, worker, or reviewer through their brief, without changing their definition or access. After an extension change, `extension add`, `update`, `remove`, `enable`, `disable`, `sync`, and `status` return `setup_required` until `agent setup` runs again; setup removes a definition whose extension is gone or off unless you edited it, and its result now lists the files `written`, `removed`, and `kept`. `check` reports a role or addition that is not valid, and two extensions adding the same heading to one role.
