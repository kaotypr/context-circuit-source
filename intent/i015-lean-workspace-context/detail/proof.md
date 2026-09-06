# Proof it is an improvement, not a thinner wrong loop

Continues from [design.md](design.md). This is how “90% already correct”
stays true after the change.

## What would count as a downgrade

Any of these, even if files got shorter or greps went down:

- a writing request that skips the intent gate
- a plan written before a feasible tracer (when the loop requires one)
- Standard or Critical work labeled verified without an independent child
- a coordinator that writes when the worker child is missing
- a rule that now exists only inside a skill
- `sources/` scanned because the catalog was silent
- always-on safety spine gone (archive exclusion, credential rule,
  invoke-not-read, host-blocked)

Length is not the test. **Behavior** is the test. Waste gone, loop intact.

## What would count as the improvement

On a typical orientation or intent-authoring turn in a product workspace:

- the knowledge units opened are a subset named by catalog entries (or none,
  if the catalog is empty)
- there is no grep or recursive list of `context/` used to *discover* units
- there is no read of the runtime implementation
- there is no search of tests to learn an invoke line
- the skill did not require opening a contract schema to draft a file that
  validate then accepts

On always-on plus skill load: restated lifecycle essays are gone; pointers
remain; the route still has a complete procedure.

## How that is checked (outcome, not a test plan)

The tracer will earn runnable checks against the real files. The shape those
checks must prove:

1. **Loop preserved** — existing semantic acceptance that encodes gates,
   tracer, verifier, host-blocked, and invoke-not-read still passes. This
   intent does not weaken those fixtures.
2. **Catalog selects** — a fixture or harness access policy where the request
   maps to named units: required reads include the catalog and the named
   units; listing `context/domains/**` or grepping `context/` is forbidden
   (or fails the case). An empty-catalog case forbids a domains-tree scan.
3. **Invoke matches** — each documented skill invoke the change touches is
   asserted against the runtime’s actual action names and argument shape, so
   a wrong line cannot ship.
4. **No schema-as-procedure** — intent (and any other authoring skill this
   change touches) names templates + validate, not “read the schema and
   transcribe fields.”
5. **Product workspace** — template seed carries the catalog shape and the
   shipped adapters/skills; a blank assembled workspace is the proving
   ground, not only this source checkout.

Harness efficiency budgets (tokens, turns) may **observe** the win. They
stay soft. They are not the definition of done. Access discipline (right
files, not a knowledge-tree hunt) is closer to the definition than a token
number.

## What we do not use as proof

- A smaller line count on coordinator or shared instructions by itself.
- A rewritten domain page.
- A new always-on file that “indexes” everything (that would add a read).
- Passing validate on this intent’s own files (necessary, not sufficient).

## Relation to the four plans

Plans 1–3 make the three wastes go away. Plan 4 is this file: the checks that
fail if the loop moved or if discovery is still required.
