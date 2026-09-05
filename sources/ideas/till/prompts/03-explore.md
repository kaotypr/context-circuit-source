# 03 — Explore (work with you live)

Use this for small, reversible polish on **one** connected folder. You are the
judge. There is no independent check and no "finished / shipped".

Do **not** use Explore for sign-in or recording payments (see case 13).

---

## 16. direct-collaboration-explore

**When:** `till-web` is connected. You want a tiny visual check you can inspect
yourself.

**Say:**

1. `Work with me directly on till-web. Add a single line of comment at the top of the invoice list stylesheet that says "till list spacing — paired locally". Leave the change uncommitted so I can inspect it.`
2. `That looks right. Leave the change uncommitted and available to continue later; don't open a pull request.`

**You should see:** Live edits in `till-web` only. Described as something you
supervised, never as independently checked. Uncommitted. No pull request.

If the stylesheet does not exist yet, say instead:

`Work with me directly on till-web. Create a small note file at src/pair-marker.txt containing exactly one line: paired locally. Leave it uncommitted so I can inspect it.`

---

## 17. explore-promote-to-standard

**When:** A try-out is turning into logic you actually care about.

**Say:**

1. `Let's work directly on till-web — add a small helper that formats money for display, like 1500 cents → "15.00".`
2. `Actually this is turning into real logic I care about. Can we make this a proper, checked change without starting over?`
3. `Yes, that goal's right — approve it with the independent check on. Don't build more or ship it yet.`
4. `Thanks.`

**You should see:** The try-out becomes a real agreed change with an independent
check, without throwing the work away. After the yes, it may derive the plan
and stop if you said not to build more. It does not open a pull request.
