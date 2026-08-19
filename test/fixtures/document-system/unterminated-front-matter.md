---
type: Design Note
title: Unterminated front matter
description: This fixture has no closing front-matter delimiter.
status: stable

# Expected extraction failure

The body cannot begin until a later standalone `---` closes the initial
front-matter block.
