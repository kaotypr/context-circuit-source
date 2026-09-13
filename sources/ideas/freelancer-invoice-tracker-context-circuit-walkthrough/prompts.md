# Exact Prompts

The quoted text is the exact text to give Context Circuit. It is written for
someone who knows what they want to build but does not know how Context Circuit
works. Replace values in angle brackets with real values.

## 1. Workspace initialization

> I want to build a Freelancer Invoice Tracker. It should help freelancers keep
> client details, create invoices, download them as PDFs, record payments, and
> see what is overdue. This is a completely new project and there are no code
> repositories yet. Help me set up the project so we can start working on it.

## 2. Repository bindings

> Create two new Git repositories for this project: `invoice-tracker-api` for
> the backend and `invoice-tracker-web` for the web app. Use `main` as the main
> branch for both and connect them to this workspace.

## 3. Agent role-tiering

This setup can use sensible defaults without asking the user anything. If the
person wants to understand how the work will be handled, they can ask:

> How are you going to plan, build, and check this project?

## 4. Team member records

> I'm `<your-name>`. Set this machine up for me.

## 5. Gathering context

There is no separate prompt at this point. This is a new project, so there are no
existing product documents or implementation details to gather. The idea already
given by the user is the starting evidence.

## Design the whole system

> Lets design the whole system for this idea

Read through the proposed system and answer any questions needed to settle the
product boundaries and important design choices. When the design matches the
idea, ask:

## Ask how to implement the design

> How can we implement all of this designed system?

This request points to the system design as the source for the implementation
path. Review the answer before starting the first piece of work.

## 6. Intent creation and detail

> Lets start with the first piece of work from that implementation path. Write
> down what we are trying to achieve so I can check it before we build it.

## 7. Intent approval and plan creation

> This matches what I want. Go ahead with it and show me how the API and web work
> should be built, including which parts need to happen first.

## 8. Plan execution and verification

> Start building the first part that's ready. Check the result carefully when
> it's finished and show me what happened.

## 9. Stack execution

> Continue with the rest of the invoice tracker work that's ready. Work on
> separate parts at the same time when that's safe, follow the API and web
> dependencies, and tell me which parts passed, failed, or need my help.

## 10. Delivery

First, create the remote repositories:

> Create private GitHub repositories for `invoice-tracker-api` and
> `invoice-tracker-web` under `<github-account-or-organization>` and connect our
> local projects to them.

Then publish the reviewed branches:

> Push the verified invoice tracker branches to GitHub.

Then open the pull requests:

> Open pull requests for the invoice tracker changes against `main` and explain
> what is included in each one.

After the pull requests have been reviewed:

> The invoice tracker pull requests `<api-pr-number>` and `<web-pr-number>` are
> ready. Merge them for me.

When the merged release is ready for an environment:

> Release the latest invoice tracker changes to `<environment>` and tell me
> whether everything completed successfully.

## 11. Context reconciliation

> The first invoice tracker release is finished. Mark that work complete and
> update our project knowledge with what we actually built, including client
> records, invoice rules, PDFs, payment status, the dashboard, and which codebase
> owns each part.

## 12. Publications

Preview first:

> Show me how the completed invoice tracker work would look in our `<provider>`
> project called `<target-name>`. Use plain English and include the main tasks
> and checks so I can review it first.

After reviewing the preview:

> That preview looks right. Send it to the `<provider>` project called
> `<target-name>` and share the link with me.
