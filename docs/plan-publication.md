# Publication

Publication is an optional human-authorized action after a plan is approved and
before external execution or collaboration.

The root session may prepare publication information and explain:

- which plan and tasks would be published;
- which external provider or destination is involved;
- what identifiers or URLs would be stored;
- what publication cannot change;
- what risks or irreversible effects exist.

The human explicitly authorizes publication. The agent must not publish, merge,
deploy, create external issues, or synchronize external status implicitly.

Publication does not:

- approve a plan;
- change plan or task status;
- start a session;
- monitor external status;
- synchronize completion back into the workspace.

Publication is outside the filesystem execution protocol. If a host or
provider integration is later added, it must remain an explicitly authorized
adapter and must not redefine session, plan, or completion state.
