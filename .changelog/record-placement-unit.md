---
product: cli
kind: changed
pr: 209
---

Every path the CLI computes for an intent or a plan now goes through one
placement unit, `internal/workspace/placement.go`, with no change in where
records go. A downstream build that wants another layout, such as new records
under a folder per member, changes that one file rather than path handling
across the package.
