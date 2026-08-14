export type PlanStatus = "draft" | "approved" | "done"
export type TaskStatus = PlanStatus

export type PlanSourceKind = "idea" | "prd" | "document" | "issue" | "pull-request"
export type PlanTrack = "epic" | "bau"

export interface RepositoryConfig {
  path: string
  mode: "ignored-clone" | "submodule"
  role: string
  agent: string
  default_branch: string
  remote?: string
}

export interface WorkspaceContextSource {
  id?: string
  kind: string
  location?: string
  reference?: string
  purpose?: string
  repository?: string
  revision?: string
  product_knowledge?: string[]
}

export interface WorkspaceConfig {
  version: number
  template_version?: string
  workspace: {
    name: string
    mode: "solo" | "team"
    default_branch: string
    purpose?: string
    remote?: string
  }
  repositories: Record<string, RepositoryConfig>
  context?: {
    sources_file?: string
    authoritative_sources?: WorkspaceContextSource[]
  }
  workflow?: {
    human_gates?: string[]
    wrapper_change_policy?: "pull-request" | "direct-commit"
  }
}

export interface BootstrapGitCommit {
  authorize_initial_commit?: boolean
  commit_message?: string
  author_name?: string
  author_email?: string
}

export interface ProductKnowledgeBaselineDomain {
  name: string
  workflows?: string[]
}

export interface ProductKnowledgeBaselineSpec {
  title: string
  purpose: string
  sources: string[]
  review_date: string
  roles: string[]
  domains: ProductKnowledgeBaselineDomain[]
  unknowns: string[]
}

export interface WorkspaceBootstrapContext {
  project_summary: string
  architecture: string[]
  conventions: string[]
  decisions: string[]
  sources?: WorkspaceContextSource[]
  product_knowledge?: ProductKnowledgeBaselineSpec
}

export interface WorkspaceBootstrapRequest {
  authorize_reviewable_changes?: boolean
  configuration: WorkspaceConfig
  context: WorkspaceBootstrapContext
  wrapper: BootstrapGitCommit & { initialize_git: boolean }
  repositories: Array<{
    name: string
    source: "new" | "clone" | "existing" | "submodule"
    url?: string
    authorize_initial_commit?: boolean
    commit_message?: string
    author_name?: string
    author_email?: string
  }>
}

export interface PlanConnection {
  type: "depends-on" | "integrates-with" | "blocks" | "related" | "supersedes"
  target: string
  description?: string
}

export interface ProductKnowledgeReference {
  references: string[]
  impact?: "none" | "implementation-only" | "behavior-change" | "new-workflow" | "retired-workflow"
}

export interface PlanYaml {
  id: string
  number: number
  title: string
  track?: PlanTrack
  status: PlanStatus
  source?: { kind: PlanSourceKind | string; reference: string }
  repositories: string[]
  dependencies?: string[]
  connections?: PlanConnection[]
  product_knowledge?: ProductKnowledgeReference
  external_reference?: string
  created_at?: string
  updated_at?: string
}

export interface PlanTask {
  id: string
  plan_id: string
  title: string
  status: TaskStatus
  description?: string
  repository: string
  area?: string
  parent_task?: string
  dependencies?: string[]
  subtasks?: string[]
  connections?: PlanConnection[]
  implementation_scope?: string[]
  test_scope?: string[]
  test_expectations?: string[]
  verification_commands?: string[]
  acceptance_criteria?: string[]
  product_knowledge?: string[]
  external_reference?: string
}

export interface PlanTaskInput extends Partial<Omit<PlanTask, "id" | "plan_id" | "status" | "repository" | "title">> {
  key?: string
  work_id?: string
  task_id?: string
  id?: string
  title: string
  repository: string
  status?: TaskStatus
  plan_id?: string
  depends_on?: string[]
}

export interface PlanCreateRequest {
  id?: string
  plan_id?: string
  number?: number
  title: string
  slug?: string
  track?: PlanTrack
  source?: { kind: PlanSourceKind | string; reference: string }
  repository?: string
  repositories?: string[]
  dependencies?: string[]
  depends_on_plans?: string[]
  connections?: PlanConnection[]
  product_knowledge?: ProductKnowledgeReference
  summary?: string
  assumptions?: string[]
  open_questions?: string[]
  requirements?: string[]
  acceptance_criteria?: string[]
  solution?: string[]
  delivery?: string[]
  verification?: string[]
  risks?: string[]
  work_prefix?: string
  tasks?: PlanTaskInput[]
  work_items?: PlanTaskInput[]
}

export interface PlanCreationSummary {
  plan_id: string
  number: number
  status: "draft"
  directory: string
  plan: string
  tasks: string[]
}

export interface PlanValidationResult {
  plan: PlanYaml | null
  tasks: PlanTask[]
  errors: string[]
  directory: string
}

export interface PlanListItem extends PlanValidationResult {
  archived: boolean
}

export interface NextPlanRecommendation {
  action: "execute" | "review"
  plan_id?: string
  title: string
  why: string
  plan_reference?: string
  product_knowledge_references: string[]
  repository?: string
  remaining_tasks?: string[]
}

export interface WhatsNextResult {
  generated_at: string
  recommendation: NextPlanRecommendation
  alternatives: NextPlanRecommendation[]
  considered: { plans: number; tasks: number; executable: number; blocked: number }
  warnings: string[]
  no_state_changed: true
}

export interface PreparedPlanExecution {
  plan: PlanYaml
  plan_domain: string
  tasks: PlanTask[]
  worktree: string
  branch: string
  base_commit: string
  prompt: string
  prompt_file: string
  review_commands: string[]
  status_changed: false
}

export interface PublicationReference {
  id: string
  url: string
}

export interface PublicationResult {
  plan_id: string
  provider: string
  published: PublicationReference[]
  unchanged_status: PlanStatus
  activity_records: []
}

export interface SourceRegistry {
  sources: SourceRecord[]
}

export interface SourceRecord {
  id: string
  kind: string
  location: string
  revision: string
  product_knowledge: string[]
  imported_at?: string
}

export interface ProductKnowledgeValidationResult {
  present: boolean
  pages: number
  errors: string[]
}

export interface ProductKnowledgeImportRequest {
  source: string
  source_id?: string
  kind?: string
  product_knowledge?: string[]
  repository?: string
  title?: string
  purpose?: string
}

export interface ProductKnowledgeImportResult {
  source: SourceRecord
  updated_pages: string[]
  proposal?: string
}

export interface ProductKnowledgeRefreshResult {
  stale: SourceRecord[]
  proposal: string
  affected_pages: string[]
}

export interface InitializationRepositorySummary {
  name: string
  path: string
  mode: "ignored-clone" | "submodule"
  role: string
  agent: string
  default_branch: string
  current_branch: string
  clean: boolean
  remote: string | null
  instructions: string | null
}

export interface InitializationSummary {
  workspace: string
  mode: "solo" | "team"
  default_branch: string
  repositories: InitializationRepositorySummary[]
  required_documents: string[]
  wrapper_changes: string[]
  gitignore_changed: boolean
  applied: boolean
  warnings: string[]
}

export interface BootstrapWorkspaceSummary extends InitializationSummary {
  status: "initialized"
  bootstrap_actions: string[]
  wrapper_initial_commit: string | null
}

export type WorkspaceConfigurationState = "fresh" | "existing"

export interface ConfigureWorkspaceSummary {
  route: "bootstrap" | "reconfigure" | "inspect-fresh" | "inspect-existing"
  state: WorkspaceConfigurationState
  message: string
  result: BootstrapWorkspaceSummary | InitializationSummary | null
}
