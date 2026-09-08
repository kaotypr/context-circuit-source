# Decisions

Status: uninitialized

Record accepted tradeoffs and constraints as durable decisions with a stable
anchor so plans and other context units can reference them. A decision changes
only through an explicit human decision, never silently from code or a plan.

A knowledge-change entry records decision, rationale, and consequence — what is
now true about the product — not edited paths or the ephemeral artifact behind
the change. Live context files hold durable product knowledge only
(INV-KNOWLEDGE-03).
