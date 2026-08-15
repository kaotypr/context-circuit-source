# Optional integrations

An integration is an opt-in adapter for an external activity record or host
convenience. It is never a dependency of session entry, context, planning,
execution, verification, or filesystem handoff.

## Adapter boundary

The portable adapter contract is:

1. **discover** — identify whether the requested provider capability is
   available, without changing workspace state;
2. **explain** — name the provider, data it would read, data it would write,
   authorization needed, and the failure fallback;
3. **authorize** — wait for explicit user authorization for the described
   external mutation; credentials remain in host-managed secure storage;
4. **publish** — send only the selected plan summary, task progress, or
   handoff metadata that the user enabled; and
5. **report** — return `published`, `disabled`, `denied`, or `unavailable`
   without changing canonical plan/task status, runtime ownership, or
   verification evidence.

The adapter may keep a provider-owned `external_status` annotation separate
from canonical task status when an explicitly configured host supports it. The
annotation is opaque to core lifecycle decisions and may be omitted entirely.

## Default and failure behavior

Integrations are `enabled: false` unless the user explicitly opts in. A
missing provider, no network access, denied authorization, or provider failure
falls back to the core filesystem workflow. The result is visible in the
current session handoff as a provider-neutral outcome; no provider payload,
credential, token, or external activity record is copied into workspace plans,
context, or ordinary runtime state.

The core workflow remains complete offline. A failed adapter must not block
implementation, verification, plan evidence, recovery, or the user's choice
of manual delivery.

## Configuration example

```yaml
configuration:
  integrations:
    - id: activity-record
      enabled: true
      provider: user-selected-provider
      reads:
        - plan summaries
        - task progress
      writes:
        - approved handoff metadata
      authorization: approved
      fallback: core-filesystem
```

This example records capability intent only. It does not store provider
credentials or authorize an unspecified mutation.
