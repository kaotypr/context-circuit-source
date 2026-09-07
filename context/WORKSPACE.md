# Workspace

This is the Context Circuit product source workspace. Its source identity is
owned by `workspace.yaml`; `template/` is the blank mutable seed copied into a
new user workspace. The source checkout's `context/` describes the product and
its maintainer conventions, not an instantiated customer project.

Released root adapters live under `.context-circuit/wrapper/adapters/` and are staged at the
workspace root by the maintainer-only release assembler.
