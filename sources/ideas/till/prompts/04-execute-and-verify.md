# 04 — Execute and verify

Building starts only after an approved change. **Checked is not finished.**
**Finished is not shipped.**

Name the work by its ordinary title ("the invoice list", "the payment recording
change"), not by internal ids, unless you already saw a number and want to
reuse it.

---

## 18. execute-standard-verify-not-complete

**When:** Invoice list (or another Standard change) is approved and ready.

**Say:**

1. `The approved invoice-list change looks good — go ahead and build it.`
2. `Great. Is it finished / all done now?`
3. `Okay, thanks.`

**You should see:** Built and independently checked. When you ask if it is
finished: **not yet** — you still accept the result and say whether to open it
onto `main`. It must not say shipped.

---

## 19. repository-grounding

**When:** A Till folder has its own agent notes (for example `AGENTS.md` in
`till-api` that says "never add a new HTTP framework"). A plan for that folder
is ready.

**Say:**

1. `Please build the plan that's ready to go.`
2. `How did that go — is it built and checked?`

**You should see:** The build honors that folder's own guidance. If it cannot,
it stops and says so rather than ignoring it. Still not finished or shipped.

---

## 20. verifier-host-blocked

**When:** You ask for the independent check and the environment cannot spawn it
(host limit, nesting unavailable). Hard to stage on purpose; use it if you
actually hit the wall.

**Say:**

1. `The code's written — can you double-check it's actually good before I rely on it?`
2. `Ah, okay — I understand. I'll deal with that.`
3. `Thanks.`

**You should see:** Honest block. It does **not** check its own work and pretend
that counts. Work is preserved.

---

## 21. run-stack-single-repo

**When:** Several approved, dependent changes all live in `till-api` (for
example: invoice schema, then invoice API, then payment log).

**Say:**

1. `Please build the invoice schema, the invoice API, and the payment log — run them all in one go, in order.`
2. `How did that go — is everything built and checked?`

**You should see:** They run in dependency order in that one project. Each still
gets its own build and check. Nothing marked done or shipped as a side effect.
If a middle one fails, later ones wait.

---

## 22. run-stack-multi-repo

**When:** An API change must land before the web app can show it, both approved.

**Say:**

1. `Please build all the ready work across till-api and till-web — run them together, API first.`
2. `How did that go — is everything built and checked?`

**You should see:** Cross-project order is respected (web does not build against
an API that is not there). Same-project stacking still applies. Progress in
plain language. Nothing shipped.

---

## 23. execution-tiering-hidden

**When:** You run a batch (case 21 or 22). This is a **watch** case, not extra
words.

**Say:** the same as case 21 or 22.

**You should see:** Ordinary progress ("built and checked"). You should **not**
see model names, effort sliders, role-tiering, or which model wrote vs checked.

---

## 24. three-failure-stop

**When:** A build was attempted and the independent check failed three times
(you will usually hear this from the coordinator, not stage it).

**Say:**

1. `Did that invoice-list change build okay in the end?`
2. `Alright — I'll take a look at it myself. Nothing got lost, right?`
3. `Okay, thanks.`

**You should see:** Honest stop after the third failed attempt. Evidence kept.
No silent fourth try. Not finished, not shipped.

---

## 25. interrupted-recovery

**When:** A build was cut off (you closed the session, the host died).

**Say:**

1. `I think that invoice-list build got cut off earlier — what's the state? Can we pick it up without starting over?`
2. `Good — just confirm nothing's lost and don't redo it from scratch. Leave it there for now.`
3. `Thanks.`

**You should see:** Resume or a read-only hold. Nothing thrown away. No
pretend-success.

---

## 26. lease-ownership-conflict

**When:** A build is **currently running** on a path (invoice API) and you ask
to start another overlapping change.

**Say:**

1. `While that's going, can you also start building another tweak in the invoice API right now?`
2. `Ah, right — okay, wait until the current one's done then. Don't force it.`
3. `Thanks.`

**You should see:** It refuses to steal in-progress work. Read-only / wait. It
does not start a second overlapping build on the same area.
