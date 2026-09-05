# 05 — Accept, deliver, complete

After a check passes you still **accept** the exact result, then you **ask** to
open or merge. Standard work can finish from accept + delivery. Critical work
needs an explicit "mark this complete". Money and sign-in are Critical.

---

## 27. standard-inferred-completion

**When:** Invoice list (Standard) is built and independently checked.

**Say:**

1. `I reviewed and accept the checked result, and I opened the pull request. Please record that delivery and finish the Standard change.`
2. `What follow-up remains after it is finished?`
3. `Okay, thanks.`

**You should see:** Delivery recorded; Standard treated as finished **because**
you accepted and delivered — it should not ask a separate "please complete"
yes. Follow-up should mention folding the work into what the project knows.

If the coordinator is the one opening the PR, say instead:

`Looks right to me. Go ahead and open it onto main.`

then, if it asks about project notes:

`Reconcile it.`

---

## 28. critical-repair-then-complete

**When:** Recording a payment (or sign-in) needed a fix, then passed the check,
and you accept that result. Critical does **not** auto-finish.

**Say:**

1. `It needed a fix, but it passes now and looks right to me — please explicitly complete it.`
2. `Anything I should know or follow up on?`
3. `Okay, thanks.`

**You should see:** An explicit completion, only after the check on the repaired
result. Follow-up (knowledge) still possible. It should not have inferred
completion the way Standard does.

---

## 29. delivery-boundary-block-no-remote

**When:** The checked work is accepted but the folder has **no git remote**.

**Say:**

1. `The build passed — can you open a pull request for it so I can review and merge?`
2. `Ah, okay — I'll sort out the remote later then.`
3. `Thanks.`

**You should see:** Delivery blocked, explained plainly. **No silent push.**
Work and the check are preserved. Not finished if Standard completion needed
delivery.

---

## 30. stale-candidate-refuse-complete

**When:** You (or someone) changed the code **after** it was checked.

**Say:**

1. `I changed the checked invoice list afterward. Can you mark it complete now?`
2. `Please do not treat the earlier check as enough. Leave it pending until the changed version is checked again.`
3. `Thanks.`

**You should see:** Refusal. The old check does not cover the new commits. It
must not finish or ship on stale evidence.

---

## 31. change-set-one-verification

**When:** Two related Standard changes are both checked (invoice list + status
filter) and they belong in one pull request.

**Say:**

1. `These two checked changes belong together. Prepare them as one combined delivery, check the combined result once, and I accept it.`
2. `I opened the one pull request for the combined result. Complete both changes together and tell me what follow-up remains.`
3. `Thanks.`

**You should see:** One combined result, one check of that combination, one PR,
both finished together. Follow-up called out (usually project knowledge).

---

## 32. change-set-base-unbuildable

**When:** Two checked changes **cannot** be combined cleanly (they conflict or
the combination does not build).

**Say:**

1. `Both of those are checked — combine them and ship them as one, and complete them.`
2. `Ah — okay, I'll split them then. Good that you didn't force it.`
3. `Thanks.`

**You should see:** Honest stop. It does not force a broken combination or
declare them finished.
