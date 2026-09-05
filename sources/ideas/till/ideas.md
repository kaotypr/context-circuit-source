# Till — change ideas

Backlog of real changes. Each row is one decision you would approve separately.
Use these as the *ask* in the prompt pack; do not dump the whole list on the
coordinator in one turn.

## v0.1 — build these

| Idea | Why it exists in the pack | Typical care |
| --- | --- | --- |
| Start the two project folders (`till-api`, `till-web`) | Orientation / new repo | — |
| Write down how Till fits together (identity, clients, invoices, payments, the app) | Larger design, then one change per concern | — |
| Client directory: add, list, open, edit | Ordinary feature | Standard |
| Invoice list with status filter | Ordinary feature; the whole-flow walkthrough | Standard |
| Create an invoice with line items and a never-repeating number | Core behavior; open question on numbering if the brief is not named | Standard |
| Record a payment against an invoice; remaining balance and paid/part-paid | Money | Critical |
| Owner sign-in | Security | Critical; refuse a "quick informal" path |
| Invoice detail page in the web app | UI on top of the API | Standard; may need both projects |

## Try-out / polish

| Idea | Why it exists in the pack | Typical care |
| --- | --- | --- |
| Tweak invoice-card spacing and type on the list | Direct, human-supervised edit | Explore |
| A small display-name helper that grows into real logic | Promote a try-out into a checked change | Explore → Standard |

## Later, each as its own decision

| Idea | Notes |
| --- | --- |
| PDF of an invoice | Own change; not implied by "create invoice" |
| CSV export of invoices | Good publish-to-tracker example; own change |
| Public invoice link a client can open without signing in | Scope-sensitive: web *and* API, plus identity |
| Email when you mark an invoice sent | Depends on the "sent means what" question |
| Automatic overdue reminder | Out of v0.1; easy to archive and restore |
| Stripe / card charge | Out of v0.1; shelve it, do not sneak it into payments |
| Edit or reverse a wrong payment | Blocked on open question 1 in the brief |
| Multi-currency | Out of v0.1 |

## Good combinations (stacks)

- **One project, several steps:** invoice schema → invoice API → invoice list.
- **Two projects, ordered:** API can create an invoice, *then* the web app can
  show it.
- **Two checked changes, one delivery:** invoice list + status filter shipped
  together, or they refuse to combine and you split them.
