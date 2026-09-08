# 08 — Whole flow (one conversation)

Play this **instead of** chaining orientation → intent → execute → deliver →
reconcile as separate sessions, once Till exists and at least `till-web` (and
usually `till-api`) is connected. It is the golden path as a lay user would
actually talk.

---

## 40. standard-feature-whole-flow

**When:** Till is connected. You want the invoice list, end to end, in one
sitting. You will make real decisions; you will not rubber-stamp.

**Say:**

1. `I want a page that lists my invoices and lets me filter them by status.`
2. `That's right — approve it.`
3. `Looks right to me.`
4. `Go ahead and open it onto main.`
5. `Reconcile it.`

**You should see, in order:**

1. A plain restatement of the goal, what is out, likely projects, any open
   question, and how carefully it will be checked — then a **single** ask to
   approve. Not a plan approval.
2. After your yes: it looks at the real code, derives the plan, builds, and
   independently checks. It does **not** call it done or shipped yet.
3. "Looks right" accepts **that** result. Still not shipped.
4. Opening onto `main` is the ship step. For Standard work, finished follows
   from accept + delivery — no extra "please complete" yes.
5. It asks to fold the shipped behavior into what the project knows. You say
   reconcile (or "no update needed") — it never accepts knowledge silently.

If it raises a question before line 2 (empty state, which project, API vs web
only), answer in ordinary words first, then say `That's right — approve it.`

---

## Suggested first real build (optional path)

If you are using this pack to actually get Till off the ground, a sensible
order of **whole decisions** is:

1. Cases 1–5 (orientation, folders, brief).
2. Case 8 (design → one draft per concern).
3. Case 40 for **invoice list** (or cases 10 → 18 → 27 → 34).
4. Case 16 for a visual tweak when you want to feel the Explore path.
5. A **new** Critical conversation for **record a payment**, then case 28.
6. A **new** Critical conversation for **sign-in**, using case 13 if it offers
   a shortcut.
7. Case 22 when API and web work are both ready.
8. Cases 38–39 when you want the backlog visible to someone else.

Keep Stripe, public invoice links, and reminders on the shelf (cases 35–37)
until you decide they are their own v0.2 decisions.
