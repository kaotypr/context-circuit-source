---
name: cc-idea-brief
description: Facilitate an early product or project idea when a workspace has no useful plans, source documents, or repositories, then create a concise source-ready Idea Brief for later Product Knowledge import.
---

# Idea Brief

Use this skill when the user has an early idea and the wrapper has little or no usable context yet: no plans, no source documents, no registered product repository, or an explicitly unshaped request.

Discuss the idea before producing a plan. Ask only the focused questions needed to understand:

- the problem or opportunity
- who is affected
- the desired outcome or experience
- the first useful scope and non-goals
- constraints, assumptions, and success signals
- important unknowns and next sources to inspect

Separate what the user said from assumptions and open questions. Do not invent current product behavior, repositories, implementation details, or evidence. Treat the conversation itself as an input, but cite any external facts with their source and location.

After the discussion has enough shape, create or update `context/IDEA-BRIEF.md` with concise Markdown sections for:

- Summary
- Problem and users
- Desired outcome
- Initial scope
- Non-goals
- Constraints and assumptions
- Success signals
- Open questions
- Sources and evidence

Keep the brief about intent, not delivered functionality. Do not create a numbered plan, register a repository, approve work, or import Product Knowledge silently. Present the Idea Brief for human confirmation.

After the human accepts the brief, continue with `$cc-import-product-knowledge` using `context/IDEA-BRIEF.md` as the source. Preserve source citations in the resulting Product Knowledge and record the brief in `context/sources.yaml`. Only then use `$cc-create-plan` when a repository/domain and implementation scope are sufficiently clear.
