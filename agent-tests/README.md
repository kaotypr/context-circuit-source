# Agent tests

A contributor tool that runs a coding host against a fresh fixture workspace.
It is a separate Go module. It is never shipped, and CI does not run it.

## How to run

From this directory:

```sh
cd agent-tests
go test ./...
go run ./cmd/agent-tests probe
go run ./cmd/agent-tests run --dry-run
```

`probe` prints each host's version, whether its pinned model is set, and which
auth variable names are present. `run --dry-run` prints the plan and does not
start a host. A live run needs `--yes`.

## What a run touches

The run root is a temporary directory outside any repository. The checkout is
read through git and is not written. The environment given to a host is an
allowlist. Auth comes only from the environment variables named in
`hosts.yaml`. Nothing is read from or written to the real home, the real
workspace, or user config. The one shared place is the Go module cache: building the
CLI uses the cache that `go env GOMODCACHE` reports, and may download a missing module
into it. Scenario setup commands run as `sh -c` in the fixture workspace, so only run
scenarios from a checkout you trust.

## What it costs

Every live run uses model usage on the account whose key it is given. A fixed
reply of OK measured 15k to 45k input tokens per host with the maintainer's
config loaded. A scenario is many times that. Use a spend-limited key. Run
someone else's branch in a container or a VM, because host permissions are
skipped inside the run.

## What it does not prove

Permission prompts. A headless run cannot measure them.
