# Plan review

`review plan X` starts a non-executing, status-preserving discussion. It never
approves, executes, or changes plan status.

The agent walks through, as needed:

- the original request and how each important detail is represented;
- objective, desired behavior, constraints, and non-goals;
- Product Knowledge and repository evidence used to ground the plan;
- repository and task mapping;
- task order, dependencies, and expected changes;
- acceptance and verification evidence;
- assumptions, open questions, and risks;
- delivery effects intentionally left outside execution.

When you resolve an open question, correct a requirement, change scope, or ask
for a task to be more specific, the agent updates the draft plan and continues
the discussion until the details are clear. There is no separate plan-approval
step: a plan derived from an approved intent executes on that approval, and there is
no automated scope gate — scope-safety is settled at delivery (Gate 2). Changing the
intent's criteria is a new decision that re-enters Gate 1.
