const command = process.argv[2];
if (!command) throw new Error("Usage: cc <command> [arguments]");
process.env.CONTEXT_CIRCUIT_WORKSPACE_ROOT = process.cwd();
process.argv.splice(2, 1);

switch (command) {
  case "validate": await import("./validate.js"); break;
  case "initialize-workspace": await import("./initialize-workspace.js"); break;
  case "configure-workspace": await import("./configure-workspace.js"); break;
  case "run-task": await import("./run-task.js"); break;
  case "review-plan": await import("./review-plan.js"); break;
  case "create-plan": await import("./create-plan.js"); break;
  case "validate-plan": await import("./validate-plan.js"); break;
  case "set-plan-state": await import("./set-plan-state.js"); break;
  case "set-task-state": await import("./set-task-state.js"); break;
  case "whats-next": await import("./whats-next.js"); break;
  case "publish-plan": await import("./publish-plan.js"); break;
  case "archive-plan": await import("./archive-plan.js"); break;
  case "unarchive-plan": await import("./unarchive-plan.js"); break;
  case "import-product-knowledge": await import("./import-product-knowledge.js"); break;
  case "refresh-product-knowledge": await import("./refresh-product-knowledge.js"); break;
  default: throw new Error(`Unknown Context Circuit command: ${command}`);
}
