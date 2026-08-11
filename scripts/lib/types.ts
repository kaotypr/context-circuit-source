export interface RepositoryConfig {
  path: string;
  mode: "ignored-clone" | "submodule";
  role: string;
  agent: string;
  default_branch: string;
}

export interface WorkspaceConfig {
  version: 1;
  template_version: "0.1.0";
  workspace: {
    name: string;
    mode: "solo" | "team";
    default_branch: string;
  };
  repositories: Record<string, RepositoryConfig>;
  activity: {
    provider: string;
    access: "auto" | "manual";
    required_capabilities: string[];
    optional_capabilities: string[];
  };
  workflow: {
    human_gates: string[];
    maximum_repair_attempts: number;
    wrapper_change_policy: "pull-request" | "direct-commit";
  };
}

export interface TaskBrief {
  contract_version: 1;
  work_id: string;
  run_id: string;
  source: { kind: "direct-request"; reference?: string };
  requested_outcome: string;
  scope: string[];
  implementation_scope?: string[];
  test_expectation?: TestExpectation;
  acceptance_criteria: string[];
  repositories: Array<{ name: string; dependency_order: number }>;
  plan: { reference: null; approval_state: "not-applicable" };
  activity: {
    reference: null;
    claim_status: "not-applicable";
    duplicate_effort_warning: true;
  };
  assumptions: string[];
  risks: string[];
  verification_commands: string[];
  authorization: { kind: "explicit-user-request"; evidence: string };
  created_at: string;
}

export type TestExpectationPolicy = "required" | "existing-coverage" | "verifier-only" | "not-required";

export interface TestExpectation {
  policy: TestExpectationPolicy;
  paths: string[];
  rationale: string;
}

export interface RuntimeRepository {
  name: string;
  base_path: string;
  base_commit: string;
  branch: string;
  worktree: string;
  worker_input: string;
  verifier_input: string;
}

export type RuntimeStatus = "preparing" | "prepared" | "running" | "verifying" | "passed" | "failed" | "blocked" | "cancelled" | "closed";

export interface ExecutionEvent {
  stage: "worker-started" | "worker-result" | "verifier-result";
  repository: string;
  from_status: "prepared" | "running" | "verifying";
  to_status: "running" | "verifying" | "passed" | "failed" | "blocked";
  inferred: boolean;
  result_path?: string;
  idempotency_key: string;
  occurred_at: string;
}

export interface RuntimeManifest {
  contract_version: 1;
  work_id: string;
  run_id: string;
  status: RuntimeStatus;
  created_at: string;
  updated_at: string;
  task_brief: string;
  repositories: RuntimeRepository[];
  evidence: string[];
  warnings: string[];
  execution_events?: ExecutionEvent[];
  lifecycle_events: Array<{
    event: "task.starting";
    status: "skipped";
    idempotency_key: string;
    occurred_at: string;
  }>;
}
