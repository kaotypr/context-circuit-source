# 01 — Orientation and getting a home

## 1. onboarding-what-is-this

**When:** Brand-new workspace. No Till folder, no brief ingested.

**Say:**

1. `I just set this up from a template but I'm honestly not sure what it is — what is this?`
2. `Okay — so what can you actually help me do?`
3. `How do I get started?`
4. `Got it, thanks.`

**You should see:** It admits the workspace is blank. It explains, in ordinary
words, that it helps you turn an idea into a change, build it, and ship it when
you say so — with an independent check, and nothing called done or shipped
without you. It does not invent Till or a project folder.

---

## 2. orient-new-project

**When:** You have a goal and no code yet.

**Say:**

1. `I want to build Till — a small invoicing app for my freelance work. I send invoices, record what got paid, and see what's still outstanding.`
2. `It's just for me, one person, nothing fancy. Can you help me get started?`
3. `I don't have an existing code folder yet — what would you need from me before we decide what to build?`
4. `There's a brief at sources/ideas/till/product-brief.md. Don't create a project folder until I say so.`

**You should see:** Orientation only. It may offer to start a new folder or use
one you already have. It does not create, clone, or register anything until you
explicitly ask.

---

## 3. clone-or-init-new-repo

**When:** You are ready for empty code homes. Till needs two folders.

**Say:**

1. `I don't have any code yet — start a fresh project folder called till-api.`
2. `Yeah, create it new, don't clone anything from anywhere.`
3. `I'll work on main, use that.`
4. `Same again for a second folder called till-web — new, on main.`
5. `Great, that's all for now.`

**You should see:** Two new project folders, both on `main`. Nothing planned or
built. If you only want one folder today, stop after `till-api` and come back.

---

## 4. connect-existing-repo

**When:** The code already lives in a folder next to the workspace (you created
it yourself, or you copied it in). Skip this if you used case 3.

**Say:**

1. `I've already got the API — the code's in a folder called till-api right here.`
2. `Can you hook it up so we can work on it together?`
3. `I do all my work on the main branch, so use that one.`
4. `The web app is in till-web. Hook that up on main too.`
5. `Great, thanks.`

**You should see:** It binds the folders you named, on the branch you named. It
does not invent a third project.

---

## 5. build-product-knowledge

**When:** At least one Till folder is connected. You want it to learn the brief
before any change.

**Say:**

1. `It's connected now — can you get up to speed on this project before we plan anything? Read sources/ideas/till/product-brief.md.`
2. `One owner, one currency, payments recorded by hand — that's right. Invoice numbers start at 1 and never repeat. Save that.`
3. `Wait — actually I'm not sure about the wrong-payment rule, don't lock that in.`

**You should see:** It proposes what the project knows. You accept some of it
and hold the open question. It does not treat "save that" as permission to
build. It does not silently decide the wrong-payment rule.

---

## 6. query-product-knowledge

**When:** Some knowledge is already accepted.

**Say:**

1. `Before I ask for any changes — what does this project already know about how invoices are numbered?`
2. `And has anything been decided about charging cards or Stripe?`
3. `Okay, good to know. Thanks.`

**You should see:** Retrieval only. Numbering should match the brief (sequential,
never reused). Stripe should be out of v0.1 / not decided as a feature. No new
folder, plan, or build.

---

## 7. accept-or-defer-context-proposal

**When:** It told you something is waiting for you to accept about the project.

**Say:**

1. `You said there was something waiting for me to look at about the project — what is it?`
2. `The "one owner, payments by hand" summary is right — accept that one. Hold the wrong-payment rule, I'm not sure about it yet.`
3. `That's fine for now.`

**You should see:** Accept and hold are both explicit. It does not accept the
held item for you. It does not start a change.
