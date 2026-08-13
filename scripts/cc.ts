const command = process.argv[2];
if (!command) throw new Error("Usage: cc <command> [arguments]");
process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT = process.cwd();
process.argv.splice(2, 1);

switch (command) {
  case "validate": await import("./validate.js"); break;
  case "initialize-workspace": await import("./initialize-workspace.js"); break;
  case "configure-workspace": await import("./configure-workspace.js"); break;
  case "run-task": await import("./run-task.js"); break;
  case "record-result": await import("./record-result.js"); break;
  case "prepare-repair": await import("./prepare-repair.js"); break;
  case "prepare-review": await import("./prepare-review.js"); break;
  case "record-review-publication": await import("./record-review-publication.js"); break;
  case "confirm-merge": await import("./confirm-merge.js"); break;
  case "finish-work": await import("./finish-work.js"); break;
  case "create-plan": await import("./create-plan.js"); break;
  case "validate-plan": await import("./validate-plan.js"); break;
  case "set-plan-state": await import("./set-plan-state.js"); break;
  case "whats-next": await import("./whats-next.js"); break;
  case "prepare-lifecycle": await import("./prepare-lifecycle.js"); break;
  case "record-lifecycle-action": await import("./record-lifecycle-action.js"); break;
  case "prepare-plan-publication": await import("./prepare-plan-publication.js"); break;
  case "record-plan-publication": await import("./record-plan-publication.js"); break;
  case "sync-context": await import("./sync-context.js"); break;
  case "prepare-context-review": await import("./prepare-context-review.js"); break;
  case "import-context": await import("./import-context.js"); break;
  case "onboarding-pack": await import("./onboarding-pack.js"); break;
  default: throw new Error(`Unknown Context Circuit command: ${command}`);
}
