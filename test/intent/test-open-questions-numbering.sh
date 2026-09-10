#!/bin/sh
# Numbered Open questions on intent authoring surfaces (template, example,
# cc-intent skill, intent domain page) plus a fixture INTENT.md that authors
# numbered questions and the unnumbered empty state.
set -eu
. "$(dirname -- "$0")/../lib/assert.sh"
. "$ROOT/test/lib/fixture.sh"

tmpl="$ROOT/.context-circuit/docs/templates/intent.md"
ex="$ROOT/.context-circuit/docs/templates/intent.example.md"
ci="$ROOT/.agents/skills/cc-intent/SKILL.md"
dom="$ROOT/context/domains/intent/README.md"

# --- NQ-AC-NUMBERED / NQ-AC-AUTHORING: shipped surfaces teach numbered questions ---
contains "$tmpl" "Open questions"
contains "$ex" "Open questions"
contains "$ci" "Open questions"
contains "$dom" "Open questions"
contains "$tmpl" "1, 2, 3"
contains "$ci" "1, 2, 3"
contains "$dom" "1, 2, 3"
contains "$tmpl" "1. **<the question>**"
contains "$ci" "answer by number"
grep -E '^1\. \*\*' "$ex" >/dev/null || fail "expected numbered question '1. **' in $ex"

# --- NQ-AC-STYLE: number, bold question, italic answer (The plans shape) ---
contains "$tmpl" "_Answer: <the decision>._"
contains "$ex" "_Answer:"
contains "$ci" "_Answer:_"
contains "$dom" "italic answer"

# --- NQ-AC-STABLE: answering keeps the number; later questions take the next unused ---
contains "$tmpl" "next unused"
contains "$ci" "next unused"
contains "$dom" "next unused"
contains "$tmpl" "Keep the question"
contains "$ci" "keeps the question"
contains "$dom" "keeps the question"

# --- NQ-AC-EMPTY: empty-state stays unnumbered; no dummy numbered item ---
contains "$tmpl" "no known unresolved"
contains "$tmpl" "Do not invent a dummy numbered item"
contains "$ci" "Do not invent a dummy numbered item"
contains "$dom" "do not invent a dummy numbered item"
contains "$tmpl" "The empty-state line stays unnumbered"
contains "$ci" "The empty-state line stays unnumbered"
# visible empty-state line is not a numbered list item
grep -E '^_At draft time: no known unresolved human decisions\._$' "$tmpl" >/dev/null \
	|| fail "expected unnumbered empty-state line in $tmpl"

# --- NQ-AC-EXISTING-UNTOUCHED: already-written intents are not rewritten ---
contains "$tmpl" "Already-written intents"
contains "$ci" "Already-written intents"
contains "$dom" "already-written intents"
contains "$ci" "Do not rewrite"
i026="$ROOT/intent/i026-compiled-runtime/INTENT.md"
if [ -f "$i026" ]; then
	awk '
		/^## Open questions$/ { in_oq=1; next }
		/^## / { in_oq=0 }
		in_oq && /^[0-9]+\. \*\*/ { found=1 }
		END { exit found ? 0 : 1 }
	' "$i026" && fail "historical $i026 Open questions must stay unnumbered" || :
fi

# --- fixture intent: numbered questions and empty state author correctly ---
ws=$(cc_fx_ws)
trap 'rm -rf "$ws"' EXIT HUP INT TERM

iid_q=i001-numbered-questions
cc_fx_intent "$ws" "$iid_q" "Numbered questions" api "src"
{
	printf '# Intention — %s\n\n' "$iid_q"
	printf '_Status: draft, waiting for your approval._\n\n'
	printf '## Intention\n\nNumbered open questions fixture.\n\n'
	printf '## Expectations\n\n- Numbered questions author correctly.\n\n'
	printf '## The plans\n\n1. **Author numbered questions.**\n   _After this:_ questions are numbered.\n\n'
	printf '## How carefully this is checked\n\n**`Explore`**\n\n'
	printf '## Open questions\n\n'
	printf '1. **Should everyone go through the queue?**\n'
	printf '   _Answer: Yes — one simple rule._\n'
	printf '2. **Should trusted members skip review?**\n'
} >"$ws/intent/$iid_q/INTENT.md"
grep -E '^1\. \*\*' "$ws/intent/$iid_q/INTENT.md" >/dev/null \
	|| fail 'fixture numbered intent missing 1. **'
grep -E '^2\. \*\*' "$ws/intent/$iid_q/INTENT.md" >/dev/null \
	|| fail 'fixture numbered intent missing next unused number 2. **'
contains "$ws/intent/$iid_q/INTENT.md" "_Answer:"
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-validate "$ws/intent/$iid_q" >/dev/null

iid_empty=i002-empty-questions
cc_fx_intent "$ws" "$iid_empty" "Empty questions" api "src"
{
	printf '# Intention — %s\n\n' "$iid_empty"
	printf '_Status: draft, waiting for your approval._\n\n'
	printf '## Intention\n\nEmpty open-questions fixture.\n\n'
	printf '## Expectations\n\n- Empty state stays unnumbered.\n\n'
	printf '## The plans\n\n1. **Keep the empty state.**\n   _After this:_ no dummy item.\n\n'
	printf '## How carefully this is checked\n\n**`Explore`**\n\n'
	printf '## Open questions\n\n'
	printf '_At draft time: no known unresolved human decisions._\n'
} >"$ws/intent/$iid_empty/INTENT.md"
contains "$ws/intent/$iid_empty/INTENT.md" "no known unresolved"
# The plans may still be numbered; only Open questions must stay unnumbered.
awk '
	/^## Open questions$/ { in_oq=1; next }
	/^## / { in_oq=0 }
	in_oq && /^[0-9]+\. \*\*/ { found=1 }
	END { exit found ? 0 : 1 }
' "$ws/intent/$iid_empty/INTENT.md" && fail 'empty-state Open questions must stay unnumbered' || :
sh "$ROOT/.context-circuit/wrapper/runtime/engine.sh" intent-validate "$ws/intent/$iid_empty" >/dev/null

pass 'intent open-questions numbering'
