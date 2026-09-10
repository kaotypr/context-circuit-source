# Session evidence and authentication

## Session identity

A real-host result records the host and version, root session identifier, start
and terminal timestamps, terminal reason, declared model settings when known,
and usage when exposed. Native child records include their identifier, parent,
role evidence, timing, and terminal state.

Tool, command, context, and usage events are attributed to the originating root
or child when the host supports it. Unattributed events remain explicitly
unattributed.

## Authentication

Drivers use an already configured host account or a dedicated test account
through supported host interfaces. They do not copy tokens, cookies, credential
stores, or configuration trees into the run root.

Disposable projects contain synthetic data only. Collectors redact configured
secret patterns before persistence, while preserving enough event shape to
support grading. Redaction never becomes an excuse to store raw provider
payloads first and clean them later.

## Limits and termination

Each scenario declares a maximum duration and may declare a usage ceiling when
the host supports enforcement. A harness limit ending a session is recorded as
such and preserves the run. It is not rewritten into a successful product result
or confused with a native agent completion.

