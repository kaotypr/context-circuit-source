# Project

Context Circuit 0.2.1 is a static, provider-neutral wrapper for safe,
human-reviewed AI-assisted delivery across one or more product repositories.
The wrapper retains compact context, workflow policy, approved plans, decisions,
and contribution summaries. Product repositories remain authoritative for code;
configured activity systems remain authoritative for live task state.

The product works without a plan, activity provider, database, dashboard,
background service, or provider SDK. A direct request with sufficient scope and
acceptance evidence can start a planless run. Plans add reviewed delivery intent
and activity integrations add optional lifecycle state; neither is mandatory.

Version 0.2.0 introduces the Context Circuit identity and a neutral release
template. Interactive initialization can create the configured wrapper and new
product Git roots, clone repositories, register local roots, or add submodules
without shipping a placeholder registration. The built wrapper receives neutral
project context populated from the approved initialization request, four focused
human guides, and a metadata-clean release archive. It also provides numbered
planning and approval, duplicate-safe task-publication preparation, read-only
next-work recommendation,
isolated single- and multi-repository execution, contract-first sequencing,
independent verification, bounded repair, review preparation, human-invoked
closeout, append-only contributions, and curated context synchronization.

Version 0.2.x continues this into one linked journey: approved plans execute
under their stable work IDs and stay tied to their runtime, review, and closeout
evidence; `whats-next` reconciles that evidence read-only and excludes completed
work; review advances through explicit local-review, publication, merge-
confirmation, and closeout states behind the human merge gate; and
`configure-workspace` sets up fresh or existing wrappers with recoverable,
reviewable configuration and cited authoritative context sources.

This release does not implement automatic merge or deployment, scheduling, a
workflow UI, a central task database, permanent agent sessions, Cursor support,
or a universal activity-provider abstraction.
