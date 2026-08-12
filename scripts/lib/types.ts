export interface RepositoryConfig {
  path: string;
  mode: "ignored-clone" | "submodule";
  role: string;
  agent: string;
  default_branch: string;
}

export interface WorkspaceConfig {
  version: 1;
  template_version: "0.2.1";
  workspace: {
    name: string;
    mode: "solo" | "team";
    default_branch: string;
  };
  repositories: Record<string, RepositoryConfig>;
  activity: {
    provider: string;
    access: "auto" | "manual";
    required_capabilities: ActivityCapability[];
    optional_capabilities: ActivityCapability[];
    lifecycle?: Partial<Record<ActivityEvent, ActivityActionPolicy[]>>;
  };
  workflow: {
    human_gates: string[];
    maximum_repair_attempts: number;
    wrapper_change_policy: "pull-request" | "direct-commit";
  };
}

export interface BootstrapGitCommit {
  authorize_initial_commit: boolean;
  commit_message: string;
  author_name?: string;
  author_email?: string;
}

export interface WorkspaceBootstrapContext {
  project_summary: string;
  architecture: string[];
  conventions: string[];
  decisions: string[];
}

export interface WorkspaceBootstrapRequest {
  contract_version: 1;
  configuration: WorkspaceConfig;
  context: WorkspaceBootstrapContext;
  wrapper: BootstrapGitCommit & { initialize_git: boolean };
  repositories: Array<{
    name: string;
    source: "new" | "clone" | "existing" | "submodule";
    url?: string;
    authorize_initial_commit: boolean;
    commit_message?: string;
    author_name?: string;
    author_email?: string;
  }>;
}

export type ActivityCapability = "read-tasks" | "update-status" | "create-tasks" | "assign-task" | "timers";
export type ActivityEvent = "task.starting" | "task.review-ready" | "task.completed" | "task.blocked" | "task.cancelled";
export type ActivityActionPolicyKind = "required" | "optional" | "manual";
export type ActivityActionStatus = "pending" | "completed" | "failed" | "skipped" | "manual";

export interface ActivityActionPolicy {
  id: string;
  capability: ActivityCapability;
  policy: ActivityActionPolicyKind;
  description: string;
}

export interface ActivityActionResult extends ActivityActionPolicy {
  status: ActivityActionStatus;
  idempotency_key: string;
  evidence: string | null;
  external_reference: string | null;
}

export interface ActivityLifecycleRecord {
  contract_version: 1;
  work_id: string;
  run_id: string;
  provider: string;
  event: ActivityEvent;
  status: "pending" | "completed" | "failed" | "skipped" | "manual";
  actions: ActivityActionResult[];
  warnings: string[];
  manual_fallbacks: string[];
  prepared_at: string;
  updated_at: string;
}

export type PlanSourceKind = "idea" | "prd" | "document" | "issue" | "pull-request";

export interface PlanIndex {
  contract_version: 1;
  plan_id: string;
  title: string;
  status: "draft" | "approved";
  plan_version: number;
  approved_at: string | null;
  approved_by: string | null;
  revision_reason: string;
  source: { kind: PlanSourceKind; reference: string };
  work_prefix: string;
  documents: string[];
  work_breakdown: string;
  material_digest: string;
  approved_digest: string | null;
  created_at: string;
  updated_at: string;
}

export interface PlanWorkItem {
  work_id: string;
  title: string;
  parent: string | null;
  depends_on: string[];
  area: string;
  repository: string;
  scope: string[];
  test_scope: string[];
  test_policy: TestExpectationPolicy;
  test_rationale?: string;
  verification_commands: string[];
  acceptance_criteria: string[];
  external_reference: string | null;
}

export interface PlanWorkBreakdown {
  contract_version: 1;
  plan_id: string;
  work_prefix: string;
  items: PlanWorkItem[];
}

export type CandidateKind =
  | "assigned-work"
  | "review"
  | "verification-failure"
  | "ci-failure"
  | "activity-task"
  | "repository-item"
  | "contribution-follow-up"
  | "plan-work-item";

export type CandidateState = "ready" | "in-progress" | "review" | "failed" | "completed" | "cancelled";
export type DependencyState = "completed" | "pending" | "unknown";
export type CandidatePlanState = "approved" | "draft" | "not-applicable" | "unknown";

export interface WorkCandidate {
  contract_version: 1;
  candidate_id: string;
  kind: CandidateKind;
  work_id: string | null;
  title: string;
  state: CandidateState;
  urgent: boolean;
  priority: number;
  owner: string | null;
  plan_reference: string | null;
  plan_approval_state: CandidatePlanState;
  dependencies: Array<{ reference: string; state: DependencyState }>;
  scope_sufficient: boolean;
  acceptance_sufficient: boolean;
  repositories: string[];
  access_available: boolean;
  contract_blocked: boolean;
  source_reference: string;
  risks: string[];
}

export interface FakeActivitySource {
  contract_version: 1;
  current_user: string;
  candidates: WorkCandidate[];
}

export interface NextAction {
  action: "execute" | "enable";
  candidate_id: string | null;
  title: string;
  why: string;
  readiness_evidence: string[];
  source_references: string[];
  repositories: string[];
  agent_sequence: string[];
  blockers: string[];
  risks: string[];
}

export interface WhatsNextResult {
  contract_version: 1;
  generated_at: string;
  recommendation: NextAction;
  alternatives: NextAction[];
  considered: { total: number; executable: number; blocked: number; excluded: number };
  warnings: string[];
  no_state_changed: true;
}

export interface PlanDraftWorkItem {
  key: string;
  title: string;
  area: string;
  repository: string;
  scope: string[];
  test_scope: string[];
  test_policy: TestExpectationPolicy;
  test_rationale?: string;
  verification_commands: string[];
  acceptance_criteria: string[];
  parent?: string;
  depends_on?: string[];
}

export interface PlanDraftRequest {
  contract_version: 1;
  plan_id: string;
  title: string;
  source: { kind: PlanSourceKind; reference: string };
  work_prefix: string;
  summary: string;
  affected_repositories: string[];
  assumptions: string[];
  open_questions: string[];
  requirements: string[];
  solution: string[];
  delivery: string[];
  verification: string[];
  risks: string[];
  work_items: PlanDraftWorkItem[];
}

export interface PlanPublicationDiscovery {
  contract_version: 1;
  provider: string;
  destination: string;
  mappings: Array<{ work_id: string; external_reference: string; evidence: string }>;
}

export interface PlanPublicationItem {
  work_id: string;
  title: string;
  parent: string | null;
  depends_on: string[];
  area: string;
  action: "create" | "skip-existing";
  status: "proposed" | "existing" | "created" | "failed";
  external_reference: string | null;
  evidence: string | null;
  idempotency_key: string;
}

export interface PlanPublicationRecord {
  contract_version: 1;
  plan_id: string;
  plan_version: number;
  approved_digest: string;
  provider: string;
  destination: string;
  status: "proposed" | "in-progress" | "partial" | "completed" | "failed";
  items: PlanPublicationItem[];
  warnings: string[];
  prepared_at: string;
  updated_at: string;
}

export interface TaskBrief {
  contract_version: 1;
  work_id: string;
  run_id: string;
  source: { kind: "direct-request"; reference?: string } | { kind: "plan"; reference: string };
  requested_outcome: string;
  scope: string[];
  implementation_scope?: string[];
  test_expectation?: TestExpectation;
  acceptance_criteria: string[];
  repositories: TaskRepositoryTarget[];
  shared_contract?: SharedContract;
  plan:
    | { reference: null; approval_state: "not-applicable" }
    | { reference: string; approval_state: "approved"; plan_version: number; approved_digest: string; work_ids: string[] };
  activity: {
    reference: null;
    claim_status: "not-applicable";
    duplicate_effort_warning: true;
  };
  assumptions: string[];
  risks: string[];
  verification_commands: string[];
  authorization: { kind: "explicit-user-request" | "confirmed-selection"; evidence: string };
  created_at: string;
}

export interface TaskRepositoryTarget {
  name: string;
  dependency_order: number;
  depends_on?: string[];
  scope?: string[];
  implementation_scope?: string[];
  test_expectation?: TestExpectation;
  verification_commands?: string[];
  acceptance_criteria?: string[];
}

export interface SharedContract {
  repository: string;
  paths: string[];
}

export interface RunTaskRepositoryRequest {
  name: string;
  depends_on: string[];
  scope: string[];
  test_scope: string[];
  test_policy: TestExpectationPolicy;
  test_rationale?: string;
  verification_commands: string[];
  acceptance_criteria: string[];
}

export interface RunTaskRequest {
  contract_version: 1;
  request: string;
  acceptance_criteria: string[];
  shared_contract: SharedContract;
  repositories: RunTaskRepositoryRequest[];
}

export interface PlanRunTaskRequest {
  contract_version: 1;
  source: {
    kind: "plan";
    reference: string;
    plan_version: number;
    approved_digest: string;
  };
  work_ids: string[];
}

export type StructuredRunTaskRequest = RunTaskRequest | PlanRunTaskRequest;

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
  status?: "waiting" | "prepared" | "running" | "verifying" | "passed" | "failed" | "blocked" | "cancelled" | "closing" | "closed";
  depends_on?: string[];
  repair_attempts?: number;
  review_preparation?: string;
  review_publication?: string;
  closeout_record?: string;
  contribution?: string;
}

export type RuntimeStatus = "preparing" | "prepared" | "running" | "verifying" | "passed" | "failed" | "blocked" | "cancelled" | "closing" | "closed";

export interface ExecutionEvent {
  stage: "worker-started" | "worker-result" | "verifier-result" | "repair-prepared" | "repair-exhausted" | "review-prepared" | "closeout-prepared" | "closeout-cleaned";
  repository: string;
  from_status: "prepared" | "running" | "verifying" | "passed" | "failed" | "blocked" | "cancelled" | "closing";
  to_status: "running" | "verifying" | "passed" | "failed" | "blocked" | "closing" | "closed";
  inferred: boolean;
  attempt?: number;
  result_path?: string;
  idempotency_key: string;
  occurred_at: string;
}

export interface ReviewPreparation {
  contract_version: 1;
  work_id: string;
  run_id: string;
  repository: string;
  status: "ready" | "blocked";
  remote: string | null;
  base_branch: string;
  head_branch: string;
  base_commit: string;
  head_commit: string;
  commits: string[];
  changed_files: string[];
  title: string;
  body: string;
  worker_result: string;
  verifier_result: string;
  blockers: string[];
  prepared_at: string;
}

export interface ReviewPublicationRecord {
  contract_version: 1;
  work_id: string;
  run_id: string;
  repository: string;
  status: "published" | "failed";
  tool: "gh" | "glab" | "manual";
  pull_request: string | null;
  evidence: string;
  head_commit: string;
  idempotency_key: string;
  recorded_at: string;
}

export interface CloseoutRecord {
  contract_version: 1;
  work_id: string;
  run_id: string;
  repository: string;
  outcome: "merged" | "abandoned";
  status: "prepared" | "blocked" | "closed";
  author: string;
  reason: string | null;
  contribution: string;
  branch: string;
  base_commit: string;
  head_commit: string;
  merge_commit: string | null;
  pull_requests: string[];
  commits: string[];
  changed_files: string[];
  verification: string[];
  cleanup: {
    requested: boolean;
    worktree_removed: boolean;
    branch_preserved: true;
    runtime_evidence_preserved: true;
  };
  blockers: string[];
  prepared_at: string;
  updated_at: string;
}

export interface RuntimeManifest {
  contract_version: 1;
  work_id: string;
  run_id: string;
  source_kind: "direct-request" | "plan" | "issue" | "pull-request" | "activity-task";
  status: RuntimeStatus;
  created_at: string;
  updated_at: string;
  task_brief: string;
  repositories: RuntimeRepository[];
  plan_work_items?: Array<{
    work_id: string;
    repository: string;
    depends_on: string[];
    outcome: "pending" | "passed" | "failed" | "blocked" | "cancelled";
  }>;
  evidence: string[];
  warnings: string[];
  execution_events?: ExecutionEvent[];
  lifecycle_events: Array<{
    event: ActivityEvent;
    status: "pending" | "completed" | "failed" | "skipped" | "manual";
    idempotency_key: string;
    occurred_at: string;
    record?: string;
    actions?: ActivityActionResult[];
  }>;
}

export type ContextSyncClassification = "durable-wrapper" | "repository-local" | "one-off" | "future-task";

export interface ContextSyncProposal {
  source_contribution: string;
  classification: ContextSyncClassification;
  summary: string;
  proposed_change?: string;
  target?: "context/PROJECT.md" | "context/ARCHITECTURE.md" | "context/CONVENTIONS.md" | "context/DECISIONS.md";
  target_repository?: string;
}

export interface ContextSyncRequest {
  contract_version: 1;
  author: string;
  contributions: string[];
  proposals: ContextSyncProposal[];
}

export interface ContextSyncRecord {
  contract_version: 1;
  sync_id: string;
  status: "prepared" | "review-ready" | "blocked";
  wrapper_mode: "pull-request" | "direct-commit";
  base_branch: string;
  base_commit: string;
  branch: string;
  worktree: string;
  request: string;
  allowed_wrapper_paths: string[];
  repository_follow_ups: Array<{ repository: string; summary: string; source_contribution: string }>;
  future_tasks: Array<{ summary: string; source_contribution: string }>;
  retained_one_offs: string[];
  changed_files: string[];
  commits: string[];
  remote: string | null;
  blockers: string[];
  prepared_at: string;
  updated_at: string;
}
