# Workspace

Status: uninitialized

This universal project workspace has not been configured for a project yet.
During initialization, record the workspace name and purpose in
`workspace.yaml`, register the repositories the project spans, and set each
repository's local `base_branch` in the ignored `repositories.local.yaml`.

Portable identity (name, purpose, logical repository keys, optional
credential-free URLs, optional `default_branch`) lives in `workspace.yaml`.
Machine-specific paths and base branches live only in
`repositories.local.yaml`.
