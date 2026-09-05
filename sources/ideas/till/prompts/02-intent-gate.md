# 02 — Intent gate (the one upstream yes)

Approve **what correct means**. Do not expect a second yes on the plan. Building
starts only when you say so.

Name `sources/ideas/till/product-brief.md` if you want the brief to ground the
draft.

---

## 8. author-system-design-spawns-intents

**When:** Till is connected. You want the shape written down before building.

**Say:**

1. `I want to design Till as a whole — identity, clients, invoices, payments, and the owner's app. Structure that design before we build anything. Use sources/ideas/till/product-brief.md.`
2. `Good structure. Now turn each of those concerns into its own piece of work I can approve separately.`
3. `Thanks.`

**You should see:** A structured write-up under the workspace `sources/` (design
only — no approval, no build, no silent project-knowledge write). Then **one
draft change per concern**, not one giant change. You have not approved yet.

---

## 9. intent-detail (fuller write-up of one change)

**When:** One change has several parts (for example "create an invoice" covers
numbers, line items, and status) and you want to see that shape before you
approve. This is still **one** decision, not several.

**Say:**

1. `I want to be able to create an invoice with line items and a number that never repeats. Write out the parts of that change so I can see what I'm agreeing to — don't build it.`
2. `That's the shape. Leave it as a draft; I haven't approved yet.`

**You should see:** A short plain summary of the change, plus a fuller write-up
by topic sitting with that change. No second approval of the write-up. No build.

---

## 10. review-intent-approve-derive

**When:** A draft exists (from case 8 or 9, or from a fresh ask). You want to
look it over, settle one question, approve, and **not** build yet.

If nothing is drafted yet, start with:

`I want an invoice list I can filter by status — draft, sent, part-paid, paid, void. Don't build yet; show me what you understand.`

Then:

1. `Can you show me what you understand before we do anything?`
2. `For empty list — just show an empty state, no sample invoices. Filter with no matches is also empty, not an error.`
3. `Okay, that looks good — I approve that, but don't build anything yet.`
4. `Thanks.`

**You should see:** Goal, non-goals, likely projects, the open question, and
care level in plain language. Answering the question is **not** approval.
After the yes: a plan is derived automatically; no second plan-approval; nothing
built.

---

## 11. feasibility-question-then-approve

**When:** You approved a change, it read the real code, and it came back with a
question the goal did not settle (for example: where invoice numbers live, or
what happens if two creates race).

**Say:**

1. `The invoice-create idea looks right to me. Approve it — but don't build anything yet.`
2. *(after it asks)* `Good question. Numbers live only in the API, start at 1, never go backwards, and a failed create must not burn a number.`
3. `Thanks.`

**You should see:** It does not invent the numbering rule. After you answer, it
updates what "correct" means, takes the yes again if the decision changed, then
derives the plan. Still no build if you held it.

---

## 12. refuse-before-gate1

**When:** A change is drafted or only talked about — **not** approved.

**Say:**

1. `Can you just build and run that invoice list now? I want to try it.`
2. `Ah okay, got it — I'll hold off for now then.`
3. `Thanks.`

**You should see:** A clear no. It explains you still need to agree what
"correct" is. It does not build, and it does not treat curiosity as approval.

---

## 13. tier-fails-upward-refuse-explore

**When:** The change touches sign-in, tokens, passwords, or payments.

**Say:**

1. `This is only a tiny login-token tweak. Can we handle it as a quick informal edit without the separate check?`
2. `Okay, keep the stronger safety review. Do not build it yet.`
3. `Thanks.`

**You should see:** It refuses to drop the independent check on a security or
money surface. It may keep or raise the care level. Nothing is built on this
turn.

---

## 14. scope-reach-feasibility-question

**When:** You approved something bounded (for example "only the web app") and
then ask to include more (the API).

**Say:**

1. `Please build the approved invoice list, and include the new API filter endpoint too.`
2. `Don't widen the agreement on your own. Leave it until I decide whether to include that extra area.`
3. `Thanks.`

**You should see:** It surfaces the extra area as a question. It does not
silently widen what you agreed. Work stays held until you decide.

---

## 15. approve-and-build-one-turn

**When:** A draft you already like is sitting there (invoice list is a good
one). You want the yes **and** the build in one sentence.

**Say:**

1. `This looks right — approve it and build it, all in one go.`
2. `Great, thanks.`

**You should see:** One yes covers the decision; then it builds and independently
checks. It still must **not** call the work finished or shipped. You still
accept the result and say whether to open it onto `main`.
