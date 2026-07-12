package rubric

// ReasonCode is a standardized finding identifier using category prefixes.
// Format: {CATEGORY}-{ISSUE}
// Example: REQ-AMBIGUOUS, SEC-MISSING_AUTH
//
// Reason codes enable automated repair workflows by providing
// machine-readable categorization of issues that map to specific
// repair strategies.
type ReasonCode string

// Reason code categories (prefixes).
const (
	CategoryREQ    = "REQ"    // Requirements
	CategoryMETRIC = "METRIC" // Metrics and measurements
	CategoryUSER   = "USER"   // User personas and journeys
	CategoryARCH   = "ARCH"   // Architecture and technical design
	CategorySEC    = "SEC"    // Security
	CategorySCALE  = "SCALE"  // Scalability and performance
	CategoryINFRA  = "INFRA"  // Infrastructure and operations
	CategoryDOC    = "DOC"    // Documentation
	CategorySCOPE  = "SCOPE"  // Scope and constraints
	CategoryUX     = "UX"     // UX and accessibility
)

// Requirements codes (REQ-*)
const (
	CodeREQAmbiguous     ReasonCode = "REQ-AMBIGUOUS"
	CodeREQNoCriteria    ReasonCode = "REQ-NO_CRITERIA"
	CodeREQConflict      ReasonCode = "REQ-CONFLICT"
	CodeREQIncomplete    ReasonCode = "REQ-INCOMPLETE"
	CodeREQUntestable    ReasonCode = "REQ-UNTESTABLE"
	CodeREQMissingReason ReasonCode = "REQ-MISSING_REASON"
)

// Metrics codes (METRIC-*)
const (
	CodeMETRICUnmeasurable ReasonCode = "METRIC-UNMEASURABLE"
	CodeMETRICNoBaseline   ReasonCode = "METRIC-NO_BASELINE"
	CodeMETRICNoTarget     ReasonCode = "METRIC-NO_TARGET"
	CodeMETRICUnrealistic  ReasonCode = "METRIC-UNREALISTIC"
	CodeMETRICNoTracking   ReasonCode = "METRIC-NO_TRACKING"
	CodeMETRICMissingKPI   ReasonCode = "METRIC-MISSING_KPI"
	CodeMETRICVanity       ReasonCode = "METRIC-VANITY"
)

// User codes (USER-*)
const (
	CodeUSERNoPersona      ReasonCode = "USER-NO_PERSONA"
	CodeUSERIncomplete     ReasonCode = "USER-INCOMPLETE"
	CodeUSERNoJourney      ReasonCode = "USER-NO_JOURNEY"
	CodeUSERUnclearProblem ReasonCode = "USER-UNCLEAR_PROBLEM"
	CodeUSERNoGoals        ReasonCode = "USER-NO_GOALS"
	CodeUSERNoPainPoints   ReasonCode = "USER-NO_PAIN_POINTS"
)

// Architecture codes (ARCH-*)
const (
	CodeARCHNoErrorHandling ReasonCode = "ARCH-NO_ERROR_HANDLING"
	CodeARCHNoAPI           ReasonCode = "ARCH-NO_API"
	CodeARCHNoDataModel     ReasonCode = "ARCH-NO_DATA_MODEL"
	CodeARCHMissingDep      ReasonCode = "ARCH-MISSING_DEP"
	CodeARCHGap             ReasonCode = "ARCH-GAP"
	CodeARCHNoInterface     ReasonCode = "ARCH-NO_INTERFACE"
	CodeARCHCircularDep     ReasonCode = "ARCH-CIRCULAR_DEP"
	CodeARCHTightCoupling   ReasonCode = "ARCH-TIGHT_COUPLING"
)

// Security codes (SEC-*)
const (
	CodeSECGap             ReasonCode = "SEC-GAP"
	CodeSECNoAuth          ReasonCode = "SEC-NO_AUTH"
	CodeSECNoAuthz         ReasonCode = "SEC-NO_AUTHZ"
	CodeSECPrivacy         ReasonCode = "SEC-PRIVACY"
	CodeSECNoValidation    ReasonCode = "SEC-NO_VALIDATION"
	CodeSECNoEncryption    ReasonCode = "SEC-NO_ENCRYPTION"
	CodeSECHardcodedSecret ReasonCode = "SEC-HARDCODED_SECRET" //nolint:gosec // G101 false positive: this is a reason code identifier, not a credential
	CodeSECInjectionRisk   ReasonCode = "SEC-INJECTION_RISK"
)

// Scalability codes (SCALE-*)
const (
	CodeSCALEConcern     ReasonCode = "SCALE-CONCERN"
	CodeSCALEPerformance ReasonCode = "SCALE-PERFORMANCE"
	CodeSCALENoCapacity  ReasonCode = "SCALE-NO_CAPACITY"
	CodeSCALESPOF        ReasonCode = "SCALE-SPOF"
	CodeSCALENoRateLimit ReasonCode = "SCALE-NO_RATE_LIMIT"
	CodeSCALENoCache     ReasonCode = "SCALE-NO_CACHE"
	CodeSCALEBlockingOp  ReasonCode = "SCALE-BLOCKING_OP"
)

// Infrastructure codes (INFRA-*)
const (
	CodeINFRANoDeploy    ReasonCode = "INFRA-NO_DEPLOY"
	CodeINFRANoMonitor   ReasonCode = "INFRA-NO_MONITOR"
	CodeINFRANoAlert     ReasonCode = "INFRA-NO_ALERT"
	CodeINFRANoRecovery  ReasonCode = "INFRA-NO_RECOVERY"
	CodeINFRANoBackup    ReasonCode = "INFRA-NO_BACKUP"
	CodeINFRANoRunbook   ReasonCode = "INFRA-NO_RUNBOOK"
	CodeINFRANoRollback  ReasonCode = "INFRA-NO_ROLLBACK"
	CodeINFRAEnvMismatch ReasonCode = "INFRA-ENV_MISMATCH"
)

// Documentation codes (DOC-*)
const (
	CodeDOCInsufficient ReasonCode = "DOC-INSUFFICIENT"
	CodeDOCOutdated     ReasonCode = "DOC-OUTDATED"
	CodeDOCNoDiagram    ReasonCode = "DOC-NO_DIAGRAM"
	CodeDOCNoExamples   ReasonCode = "DOC-NO_EXAMPLES"
	CodeDOCInconsistent ReasonCode = "DOC-INCONSISTENT"
)

// Scope codes (SCOPE-*)
const (
	CodeSCOPECreep      ReasonCode = "SCOPE-CREEP"
	CodeSCOPEUnbounded  ReasonCode = "SCOPE-UNBOUNDED"
	CodeSCOPENoConstr   ReasonCode = "SCOPE-NO_CONSTRAINTS"
	CodeSCOPENoNonGoals ReasonCode = "SCOPE-NO_NON_GOALS"
	CodeSCOPEMVPUnclear ReasonCode = "SCOPE-MVP_UNCLEAR"
	CodeSCOPENoTimeline ReasonCode = "SCOPE-NO_TIMELINE"
)

// UX codes (UX-*)
const (
	CodeUXNoARIA        ReasonCode = "UX-NO_ARIA"
	CodeUXNoErrorState  ReasonCode = "UX-NO_ERROR_STATE"
	CodeUXNoLoading     ReasonCode = "UX-NO_LOADING"
	CodeUXNoEmpty       ReasonCode = "UX-NO_EMPTY"
	CodeUXNoResponsive  ReasonCode = "UX-NO_RESPONSIVE"
	CodeUXNoKeyboard    ReasonCode = "UX-NO_KEYBOARD"
	CodeUXIncompleteNav ReasonCode = "UX-INCOMPLETE_NAV"
	CodeUXNoFeedback    ReasonCode = "UX-NO_FEEDBACK"
)

// Generic codes
const (
	CodeOther ReasonCode = "OTHER"
)

// ReasonCodeInfo provides metadata about a reason code including
// information needed for AI-assisted automated repair.
type ReasonCodeInfo struct {
	// Code is the reason code identifier.
	Code ReasonCode `json:"code"`

	// Category is the prefix category (REQ, SEC, ARCH, etc.).
	Category string `json:"category"`

	// Description explains what this code means.
	Description string `json:"description"`

	// DefaultSeverity is the typical severity for this issue.
	DefaultSeverity Severity `json:"defaultSeverity"`

	// RepairPrompt is the AI prompt for automated repair.
	// This should be a clear instruction that an LLM can follow
	// to fix the issue in the spec document.
	RepairPrompt string `json:"repairPrompt"`

	// RequiresHuman indicates if human review is needed after AI repair.
	// True for security-critical, business-critical, or subjective issues.
	RequiresHuman bool `json:"requiresHuman"`

	// SpecTypes lists which spec types this code applies to.
	SpecTypes []string `json:"specTypes,omitempty"`
}

// ReasonCodeRegistry maps codes to their metadata.
var ReasonCodeRegistry = map[ReasonCode]ReasonCodeInfo{
	// =========================================================================
	// REQ-* : Requirements
	// =========================================================================
	CodeREQAmbiguous: {
		Code:            CodeREQAmbiguous,
		Category:        CategoryREQ,
		Description:     "Requirement lacks specificity or has multiple interpretations",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Rewrite this requirement to be specific and unambiguous. Include: (1) a clear action verb, (2) measurable outcome, (3) specific constraints or boundaries. Remove any vague terms like 'fast', 'user-friendly', 'easy'.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd", "trd"},
	},
	CodeREQNoCriteria: {
		Code:            CodeREQNoCriteria,
		Category:        CategoryREQ,
		Description:     "Requirement lacks acceptance criteria for verification",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add acceptance criteria using Given/When/Then format. Include: (1) preconditions (Given), (2) action trigger (When), (3) expected outcome (Then). Cover both success and failure scenarios.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd"},
	},
	CodeREQConflict: {
		Code:            CodeREQConflict,
		Category:        CategoryREQ,
		Description:     "Two or more requirements contradict each other",
		DefaultSeverity: SeverityCritical,
		RepairPrompt:    "Identify the conflicting requirements and resolve by: (1) clarifying which takes priority, (2) merging into a consistent requirement, or (3) adding conditional logic to handle both cases. Document the resolution rationale.",
		RequiresHuman:   true, // Conflicts need human judgment
		SpecTypes:       []string{"prd", "mrd", "trd"},
	},
	CodeREQIncomplete: {
		Code:            CodeREQIncomplete,
		Category:        CategoryREQ,
		Description:     "Requirement is missing essential details",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Complete this requirement by adding missing elements: (1) WHO is affected, (2) WHAT action/capability, (3) WHY it's needed (business value), (4) any constraints or dependencies.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeREQUntestable: {
		Code:            CodeREQUntestable,
		Category:        CategoryREQ,
		Description:     "Requirement cannot be objectively verified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Rewrite this requirement to be testable. Add: (1) specific numeric thresholds where applicable, (2) clear pass/fail conditions, (3) measurement method. Replace subjective terms with objective criteria.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd"},
	},
	CodeREQMissingReason: {
		Code:            CodeREQMissingReason,
		Category:        CategoryREQ,
		Description:     "Requirement lacks justification or business rationale",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add a rationale explaining WHY this requirement exists. Include: (1) business value or user benefit, (2) problem it solves, (3) consequences of not implementing it.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},

	// =========================================================================
	// METRIC-* : Metrics
	// =========================================================================
	CodeMETRICUnmeasurable: {
		Code:            CodeMETRICUnmeasurable,
		Category:        CategoryMETRIC,
		Description:     "Success metric cannot be objectively measured",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Make this metric measurable by specifying: (1) exact data source, (2) calculation formula, (3) measurement frequency, (4) responsible team/system for tracking.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeMETRICNoBaseline: {
		Code:            CodeMETRICNoBaseline,
		Category:        CategoryMETRIC,
		Description:     "Success metric lacks a baseline value",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add baseline value for this metric. Include: (1) current value or 'TBD - to be measured in sprint 1', (2) date of measurement, (3) methodology used to obtain baseline.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeMETRICNoTarget: {
		Code:            CodeMETRICNoTarget,
		Category:        CategoryMETRIC,
		Description:     "Success metric lacks a target value",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add a specific target value with: (1) numeric goal, (2) timeframe to achieve it, (3) rationale for why this target was chosen.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeMETRICUnrealistic: {
		Code:            CodeMETRICUnrealistic,
		Category:        CategoryMETRIC,
		Description:     "Target appears unrealistic given constraints",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Revise this target to be achievable. Either: (1) provide evidence/benchmarks supporting the target, (2) adjust to a realistic value with justification, or (3) break into incremental milestones.",
		RequiresHuman:   true, // Business judgment needed
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeMETRICNoTracking: {
		Code:            CodeMETRICNoTracking,
		Category:        CategoryMETRIC,
		Description:     "No plan for how metrics will be tracked",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add tracking plan specifying: (1) tool/system for measurement, (2) reporting frequency, (3) dashboard or alert thresholds, (4) who reviews the metrics.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd"},
	},
	CodeMETRICMissingKPI: {
		Code:            CodeMETRICMissingKPI,
		Category:        CategoryMETRIC,
		Description:     "Key performance indicator not defined for critical feature",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add KPI for this feature covering: (1) primary success metric, (2) leading indicators, (3) guardrail metrics to prevent negative side effects.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeMETRICVanity: {
		Code:            CodeMETRICVanity,
		Category:        CategoryMETRIC,
		Description:     "Metric does not correlate with actual business value",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Replace this vanity metric with an actionable metric that: (1) correlates with business outcomes, (2) can influence decisions, (3) measures user value not just activity.",
		RequiresHuman:   true, // Business judgment needed
		SpecTypes:       []string{"prd", "mrd"},
	},

	// =========================================================================
	// USER-* : User Personas and Journeys
	// =========================================================================
	CodeUSERNoPersona: {
		Code:            CodeUSERNoPersona,
		Category:        CategoryUSER,
		Description:     "Target user persona not defined",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add user persona including: (1) name and role, (2) demographics/context, (3) goals and motivations, (4) pain points and frustrations, (5) technical proficiency level.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd", "uxd"},
	},
	CodeUSERIncomplete: {
		Code:            CodeUSERIncomplete,
		Category:        CategoryUSER,
		Description:     "User persona lacks essential details",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Complete this persona by adding missing elements: goals, pain points, behavioral characteristics, technical context, or typical use scenarios.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "uxd"},
	},
	CodeUSERNoJourney: {
		Code:            CodeUSERNoJourney,
		Category:        CategoryUSER,
		Description:     "User journey or flow not documented",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add user journey showing: (1) entry point/trigger, (2) step-by-step actions, (3) decision points, (4) success outcome, (5) potential failure points and recovery.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "uxd"},
	},
	CodeUSERUnclearProblem: {
		Code:            CodeUSERUnclearProblem,
		Category:        CategoryUSER,
		Description:     "Problem statement is vague or unclear",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Clarify the problem statement using this format: '[User persona] needs [capability] because [reason], but currently [obstacle/pain point].' Be specific about the impact.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeUSERNoGoals: {
		Code:            CodeUSERNoGoals,
		Category:        CategoryUSER,
		Description:     "User goals not articulated",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add user goals distinguishing between: (1) functional goals (tasks to complete), (2) emotional goals (how they want to feel), (3) social goals (how they want to be perceived).",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "uxd"},
	},
	CodeUSERNoPainPoints: {
		Code:            CodeUSERNoPainPoints,
		Category:        CategoryUSER,
		Description:     "User pain points not documented",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Document pain points including: (1) current workarounds users employ, (2) time/money cost of the problem, (3) emotional frustration, (4) frequency of occurrence.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd", "uxd"},
	},

	// =========================================================================
	// ARCH-* : Architecture and Technical Design
	// =========================================================================
	CodeARCHNoErrorHandling: {
		Code:            CodeARCHNoErrorHandling,
		Category:        CategoryARCH,
		Description:     "Error handling strategy is incomplete",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add error handling covering: (1) error categories (network, validation, auth, server), (2) retry strategy per category, (3) user-facing error messages, (4) logging/alerting approach, (5) graceful degradation.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeARCHNoAPI: {
		Code:            CodeARCHNoAPI,
		Category:        CategoryARCH,
		Description:     "API contract or interface not specified",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add API specification including: (1) endpoint paths and methods, (2) request/response schemas with examples, (3) authentication requirements, (4) rate limits, (5) error response format, (6) versioning strategy.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},
	CodeARCHNoDataModel: {
		Code:            CodeARCHNoDataModel,
		Category:        CategoryARCH,
		Description:     "Data model or schema not defined",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add data model with: (1) entity definitions and fields, (2) relationships between entities, (3) data types and constraints, (4) indexes for query patterns, (5) migration strategy for schema changes.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},
	CodeARCHMissingDep: {
		Code:            CodeARCHMissingDep,
		Category:        CategoryARCH,
		Description:     "Required dependency not documented",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Document this dependency including: (1) name and version, (2) why it's needed, (3) license compatibility, (4) maintenance status, (5) fallback if deprecated.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeARCHGap: {
		Code:            CodeARCHGap,
		Category:        CategoryARCH,
		Description:     "Architecture has unexplained gaps or missing components",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Address the architectural gap by: (1) adding the missing component with rationale, (2) explaining why it's intentionally omitted, or (3) marking as future work with timeline.",
		RequiresHuman:   true, // Architecture decisions need review
		SpecTypes:       []string{"trd"},
	},
	CodeARCHNoInterface: {
		Code:            CodeARCHNoInterface,
		Category:        CategoryARCH,
		Description:     "Interface between components not defined",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Define the interface including: (1) method signatures, (2) input/output types, (3) error conditions, (4) idempotency guarantees, (5) versioning approach.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},
	CodeARCHCircularDep: {
		Code:            CodeARCHCircularDep,
		Category:        CategoryARCH,
		Description:     "Circular dependency detected between components",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Resolve circular dependency by: (1) extracting shared logic to a new module, (2) using dependency injection, (3) introducing an interface/abstraction layer, or (4) restructuring the component boundaries.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd"},
	},
	CodeARCHTightCoupling: {
		Code:            CodeARCHTightCoupling,
		Category:        CategoryARCH,
		Description:     "Components are too tightly coupled",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Reduce coupling by: (1) defining clear interfaces, (2) using events/messages instead of direct calls, (3) applying dependency inversion, (4) documenting the bounded context.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd"},
	},

	// =========================================================================
	// SEC-* : Security
	// =========================================================================
	CodeSECGap: {
		Code:            CodeSECGap,
		Category:        CategorySEC,
		Description:     "Security consideration not addressed",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add security controls addressing: (1) threat model for this feature, (2) mitigation strategies, (3) security testing requirements, (4) incident response considerations.",
		RequiresHuman:   true, // Security always needs review
		SpecTypes:       []string{"trd", "ird", "prd"},
	},
	CodeSECNoAuth: {
		Code:            CodeSECNoAuth,
		Category:        CategorySEC,
		Description:     "Authentication mechanism not specified",
		DefaultSeverity: SeverityCritical,
		RepairPrompt:    "Define authentication including: (1) auth method (OAuth, JWT, API key, etc.), (2) token/session management, (3) password/credential requirements, (4) MFA considerations, (5) session timeout policy.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd"},
	},
	CodeSECNoAuthz: {
		Code:            CodeSECNoAuthz,
		Category:        CategorySEC,
		Description:     "Authorization model not defined",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add authorization covering: (1) role definitions, (2) permission matrix, (3) resource-level access rules, (4) principle of least privilege application, (5) authorization check points.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd"},
	},
	CodeSECPrivacy: {
		Code:            CodeSECPrivacy,
		Category:        CategorySEC,
		Description:     "Data privacy requirements not addressed",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Document data privacy including: (1) PII identification, (2) data classification, (3) retention policy, (4) deletion/anonymization procedures, (5) compliance requirements (GDPR, CCPA, etc.).",
		RequiresHuman:   true, // Privacy needs legal/compliance review
		SpecTypes:       []string{"trd", "prd"},
	},
	CodeSECNoValidation: {
		Code:            CodeSECNoValidation,
		Category:        CategorySEC,
		Description:     "Input validation not specified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add input validation covering: (1) allowed characters/formats, (2) length limits, (3) sanitization rules, (4) validation error messages, (5) server-side validation (never trust client).",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},
	CodeSECNoEncryption: {
		Code:            CodeSECNoEncryption,
		Category:        CategorySEC,
		Description:     "Encryption requirements not specified",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Define encryption requirements: (1) data at rest encryption, (2) data in transit (TLS version), (3) key management strategy, (4) sensitive field encryption.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeSECHardcodedSecret: {
		Code:            CodeSECHardcodedSecret,
		Category:        CategorySEC,
		Description:     "Hardcoded secrets or credentials detected",
		DefaultSeverity: SeverityCritical,
		RepairPrompt:    "Remove hardcoded secrets and specify: (1) secret management system to use, (2) environment variable naming convention, (3) rotation policy, (4) access audit logging.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeSECInjectionRisk: {
		Code:            CodeSECInjectionRisk,
		Category:        CategorySEC,
		Description:     "Potential injection vulnerability (SQL, XSS, command)",
		DefaultSeverity: SeverityCritical,
		RepairPrompt:    "Address injection risk by specifying: (1) parameterized queries (SQL), (2) output encoding (XSS), (3) input sanitization, (4) CSP headers, (5) security testing requirements.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd"},
	},

	// =========================================================================
	// SCALE-* : Scalability and Performance
	// =========================================================================
	CodeSCALEConcern: {
		Code:            CodeSCALEConcern,
		Category:        CategorySCALE,
		Description:     "Scalability concern not addressed",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add scaling strategy covering: (1) expected load (requests/sec, concurrent users), (2) horizontal vs vertical scaling approach, (3) bottleneck identification, (4) auto-scaling triggers.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeSCALEPerformance: {
		Code:            CodeSCALEPerformance,
		Category:        CategorySCALE,
		Description:     "Performance risk identified but not mitigated",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add performance requirements: (1) latency targets (p50, p95, p99), (2) throughput requirements, (3) resource limits (CPU, memory), (4) performance testing approach.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeSCALENoCapacity: {
		Code:            CodeSCALENoCapacity,
		Category:        CategorySCALE,
		Description:     "Capacity planning not documented",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add capacity plan including: (1) initial resource sizing, (2) growth projections (6mo, 1yr), (3) cost estimates, (4) scaling thresholds.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeSCALESPOF: {
		Code:            CodeSCALESPOF,
		Category:        CategorySCALE,
		Description:     "Single point of failure identified",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Eliminate SPOF by adding: (1) redundancy strategy (active-active, active-passive), (2) failover mechanism, (3) health check configuration, (4) recovery time objective.",
		RequiresHuman:   true,
		SpecTypes:       []string{"trd", "ird"},
	},
	CodeSCALENoRateLimit: {
		Code:            CodeSCALENoRateLimit,
		Category:        CategorySCALE,
		Description:     "Rate limiting not specified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add rate limiting covering: (1) limits per user/IP/API key, (2) time windows, (3) response when exceeded (429 status), (4) bypass rules for internal services.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},
	CodeSCALENoCache: {
		Code:            CodeSCALENoCache,
		Category:        CategorySCALE,
		Description:     "Caching strategy not defined for frequently accessed data",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add caching strategy: (1) what to cache, (2) cache location (CDN, Redis, in-memory), (3) TTL/expiration policy, (4) cache invalidation approach.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},
	CodeSCALEBlockingOp: {
		Code:            CodeSCALEBlockingOp,
		Category:        CategorySCALE,
		Description:     "Blocking operation in critical path",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Address blocking operation by: (1) making it async with queue/callback, (2) adding timeout, (3) implementing circuit breaker, (4) providing fallback response.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd"},
	},

	// =========================================================================
	// INFRA-* : Infrastructure and Operations
	// =========================================================================
	CodeINFRANoDeploy: {
		Code:            CodeINFRANoDeploy,
		Category:        CategoryINFRA,
		Description:     "Deployment strategy not defined",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add deployment plan covering: (1) deployment method (blue-green, canary, rolling), (2) environment progression (dev→staging→prod), (3) deployment automation, (4) smoke test requirements.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeINFRANoMonitor: {
		Code:            CodeINFRANoMonitor,
		Category:        CategoryINFRA,
		Description:     "Monitoring strategy not defined",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add monitoring covering: (1) key metrics to track, (2) logging strategy and retention, (3) tracing/correlation IDs, (4) dashboard requirements, (5) SLO definitions.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeINFRANoAlert: {
		Code:            CodeINFRANoAlert,
		Category:        CategoryINFRA,
		Description:     "Alerting and on-call strategy not defined",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add alerting strategy: (1) alert thresholds and conditions, (2) severity levels, (3) notification channels, (4) escalation policy, (5) on-call rotation.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeINFRANoRecovery: {
		Code:            CodeINFRANoRecovery,
		Category:        CategoryINFRA,
		Description:     "Disaster recovery plan not documented",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add DR plan including: (1) RTO target, (2) RPO target, (3) recovery procedures, (4) data backup strategy, (5) DR testing schedule.",
		RequiresHuman:   true, // DR plans need stakeholder review
		SpecTypes:       []string{"ird"},
	},
	CodeINFRANoBackup: {
		Code:            CodeINFRANoBackup,
		Category:        CategoryINFRA,
		Description:     "Backup strategy not specified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add backup strategy: (1) what is backed up, (2) backup frequency, (3) retention period, (4) backup location (offsite), (5) restore testing procedure.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeINFRANoRunbook: {
		Code:            CodeINFRANoRunbook,
		Category:        CategoryINFRA,
		Description:     "Operational runbook not provided",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Create runbook covering: (1) common operational tasks, (2) troubleshooting steps for known issues, (3) escalation contacts, (4) maintenance procedures.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeINFRANoRollback: {
		Code:            CodeINFRANoRollback,
		Category:        CategoryINFRA,
		Description:     "Rollback procedure not defined",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add rollback procedure: (1) trigger criteria for rollback, (2) step-by-step rollback process, (3) data migration rollback (if applicable), (4) verification steps post-rollback.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},
	CodeINFRAEnvMismatch: {
		Code:            CodeINFRAEnvMismatch,
		Category:        CategoryINFRA,
		Description:     "Environment configuration mismatch risk",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Address environment parity: (1) document environment differences, (2) use infrastructure-as-code, (3) implement config validation, (4) add environment-specific testing.",
		RequiresHuman:   false,
		SpecTypes:       []string{"ird"},
	},

	// =========================================================================
	// DOC-* : Documentation
	// =========================================================================
	CodeDOCInsufficient: {
		Code:            CodeDOCInsufficient,
		Category:        CategoryDOC,
		Description:     "Documentation is insufficient for implementation",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Expand documentation to include: (1) implementation details, (2) code examples, (3) edge cases and their handling, (4) integration points.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "trd", "ird", "uxd"},
	},
	CodeDOCOutdated: {
		Code:            CodeDOCOutdated,
		Category:        CategoryDOC,
		Description:     "Reference to outdated information or deprecated feature",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Update the outdated reference to: (1) current version/approach, (2) migration path if applicable, (3) deprecation timeline if relevant.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "trd", "ird"},
	},
	CodeDOCNoDiagram: {
		Code:            CodeDOCNoDiagram,
		Category:        CategoryDOC,
		Description:     "Visual diagram would improve clarity",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add a diagram showing: (1) component relationships, (2) data flow, (3) sequence of operations, or (4) architecture overview. Use Mermaid or similar text-based diagram format.",
		RequiresHuman:   false,
		SpecTypes:       []string{"trd", "ird", "uxd"},
	},
	CodeDOCNoExamples: {
		Code:            CodeDOCNoExamples,
		Category:        CategoryDOC,
		Description:     "Missing examples to clarify usage",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add concrete examples showing: (1) typical usage, (2) edge cases, (3) error scenarios. Include request/response examples for APIs.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "trd"},
	},
	CodeDOCInconsistent: {
		Code:            CodeDOCInconsistent,
		Category:        CategoryDOC,
		Description:     "Documentation inconsistent with other sections",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Resolve the inconsistency by: (1) identifying the source of truth, (2) updating conflicting sections, (3) adding cross-references to avoid future drift.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "trd", "ird", "uxd"},
	},

	// =========================================================================
	// SCOPE-* : Scope and Constraints
	// =========================================================================
	CodeSCOPECreep: {
		Code:            CodeSCOPECreep,
		Category:        CategorySCOPE,
		Description:     "Scope includes unnecessary features or complexity",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Address scope creep by: (1) moving non-essential items to 'Future Work', (2) justifying inclusion with business value, or (3) removing entirely with rationale.",
		RequiresHuman:   true, // Scope decisions need stakeholder input
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeSCOPEUnbounded: {
		Code:            CodeSCOPEUnbounded,
		Category:        CategorySCOPE,
		Description:     "Scope is unbounded or unclear",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Define scope boundaries: (1) explicit 'In Scope' list, (2) explicit 'Out of Scope' / 'Non-Goals' list, (3) decision criteria for borderline items.",
		RequiresHuman:   true,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeSCOPENoConstr: {
		Code:            CodeSCOPENoConstr,
		Category:        CategorySCOPE,
		Description:     "Constraints or limitations not documented",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Document constraints including: (1) technical constraints (platform, language, etc.), (2) business constraints (budget, timeline), (3) resource constraints (team size, skills).",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "trd"},
	},
	CodeSCOPENoNonGoals: {
		Code:            CodeSCOPENoNonGoals,
		Category:        CategorySCOPE,
		Description:     "Non-goals not explicitly stated",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add non-goals section listing: (1) features explicitly NOT being built, (2) use cases NOT being supported, (3) rationale for each exclusion.",
		RequiresHuman:   false,
		SpecTypes:       []string{"prd", "mrd"},
	},
	CodeSCOPEMVPUnclear: {
		Code:            CodeSCOPEMVPUnclear,
		Category:        CategorySCOPE,
		Description:     "MVP scope not clearly defined",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Define MVP clearly: (1) minimum feature set for launch, (2) success criteria for MVP, (3) what's deferred to post-MVP, (4) timeline for MVP.",
		RequiresHuman:   true,
		SpecTypes:       []string{"prd"},
	},
	CodeSCOPENoTimeline: {
		Code:            CodeSCOPENoTimeline,
		Category:        CategorySCOPE,
		Description:     "Timeline or milestones not specified",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add timeline with: (1) key milestones, (2) dependencies between milestones, (3) target dates (even if tentative), (4) critical path identification.",
		RequiresHuman:   true, // Timelines need PM input
		SpecTypes:       []string{"prd", "mrd"},
	},

	// =========================================================================
	// UX-* : UX and Accessibility
	// =========================================================================
	CodeUXNoARIA: {
		Code:            CodeUXNoARIA,
		Category:        CategoryUX,
		Description:     "ARIA labels or accessibility attributes not specified",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Add accessibility requirements: (1) ARIA labels for interactive elements, (2) alt text requirements for images, (3) heading hierarchy, (4) focus management, (5) screen reader considerations.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXNoErrorState: {
		Code:            CodeUXNoErrorState,
		Category:        CategoryUX,
		Description:     "Error state UI not designed",
		DefaultSeverity: SeverityHigh,
		RepairPrompt:    "Design error states including: (1) error message content and tone, (2) visual treatment, (3) recovery actions available to user, (4) retry mechanisms, (5) when to show inline vs page-level errors.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXNoLoading: {
		Code:            CodeUXNoLoading,
		Category:        CategoryUX,
		Description:     "Loading state not specified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add loading state design: (1) loading indicator type (spinner, skeleton, progress), (2) placement, (3) timeout handling, (4) partial content display strategy.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXNoEmpty: {
		Code:            CodeUXNoEmpty,
		Category:        CategoryUX,
		Description:     "Empty state not designed",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Design empty state with: (1) friendly message explaining why empty, (2) illustration or icon, (3) call-to-action to populate, (4) help/onboarding content.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXNoResponsive: {
		Code:            CodeUXNoResponsive,
		Category:        CategoryUX,
		Description:     "Responsive/mobile behavior not specified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add responsive design specs: (1) breakpoints, (2) layout changes per breakpoint, (3) touch target sizes, (4) mobile-specific interactions, (5) content prioritization for small screens.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXNoKeyboard: {
		Code:            CodeUXNoKeyboard,
		Category:        CategoryUX,
		Description:     "Keyboard navigation not specified",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Add keyboard accessibility: (1) tab order, (2) keyboard shortcuts, (3) focus indicators, (4) skip links, (5) escape key behavior for modals.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXIncompleteNav: {
		Code:            CodeUXIncompleteNav,
		Category:        CategoryUX,
		Description:     "Navigation flow incomplete",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Complete navigation by adding: (1) entry points to this screen, (2) exit points/next steps, (3) back navigation behavior, (4) breadcrumb requirements, (5) deep linking support.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},
	CodeUXNoFeedback: {
		Code:            CodeUXNoFeedback,
		Category:        CategoryUX,
		Description:     "User feedback mechanism not designed",
		DefaultSeverity: SeverityLow,
		RepairPrompt:    "Add user feedback design: (1) success confirmations, (2) progress indicators for long operations, (3) toast/notification placement, (4) undo capabilities.",
		RequiresHuman:   false,
		SpecTypes:       []string{"uxd"},
	},

	// =========================================================================
	// OTHER
	// =========================================================================
	CodeOther: {
		Code:            CodeOther,
		Category:        "other",
		Description:     "Issue does not fit standard categories",
		DefaultSeverity: SeverityMedium,
		RepairPrompt:    "Review the specific issue and address based on context. Consider whether a new reason code category should be created for this type of issue.",
		RequiresHuman:   true,
		SpecTypes:       []string{},
	},
}

// Legacy code mappings for backwards compatibility.
// Maps old code names to new prefixed codes.
var LegacyCodeMapping = map[ReasonCode]ReasonCode{
	"AMBIGUOUS_REQUIREMENT":       CodeREQAmbiguous,
	"MISSING_ACCEPTANCE_CRITERIA": CodeREQNoCriteria,
	"CONFLICTING_REQUIREMENTS":    CodeREQConflict,
	"INCOMPLETE_REQUIREMENT":      CodeREQIncomplete,
	"UNTESTABLE_REQUIREMENT":      CodeREQUntestable,
	"UNMEASURABLE_SUCCESS_METRIC": CodeMETRICUnmeasurable,
	"MISSING_METRIC_BASELINE":     CodeMETRICNoBaseline,
	"MISSING_METRIC_TARGET":       CodeMETRICNoTarget,
	"UNREALISTIC_TARGET":          CodeMETRICUnrealistic,
	"MISSING_TRACKING_PLAN":       CodeMETRICNoTracking,
	"MISSING_USER_PERSONA":        CodeUSERNoPersona,
	"INCOMPLETE_PERSONA":          CodeUSERIncomplete,
	"MISSING_USER_JOURNEY":        CodeUSERNoJourney,
	"UNCLEAR_PROBLEM_STATEMENT":   CodeUSERUnclearProblem,
	"INCOMPLETE_ERROR_HANDLING":   CodeARCHNoErrorHandling,
	"MISSING_API_CONTRACT":        CodeARCHNoAPI,
	"MISSING_DATA_MODEL":          CodeARCHNoDataModel,
	"MISSING_DEPENDENCY":          CodeARCHMissingDep,
	"ARCHITECTURE_GAP":            CodeARCHGap,
	"SECURITY_GAP":                CodeSECGap,
	"MISSING_AUTHENTICATION":      CodeSECNoAuth,
	"MISSING_AUTHORIZATION":       CodeSECNoAuthz,
	"DATA_PRIVACY_CONCERN":        CodeSECPrivacy,
	"MISSING_INPUT_VALIDATION":    CodeSECNoValidation,
	"SCALABILITY_CONCERN":         CodeSCALEConcern,
	"PERFORMANCE_RISK":            CodeSCALEPerformance,
	"MISSING_CAPACITY_PLAN":       CodeSCALENoCapacity,
	"SINGLE_POINT_OF_FAILURE":     CodeSCALESPOF,
	"MISSING_DEPLOYMENT_PLAN":     CodeINFRANoDeploy,
	"MISSING_MONITORING":          CodeINFRANoMonitor,
	"MISSING_ALERTING_STRATEGY":   CodeINFRANoAlert,
	"MISSING_RECOVERY_PLAN":       CodeINFRANoRecovery,
	"INSUFFICIENT_DOCUMENTATION":  CodeDOCInsufficient,
	"OUTDATED_REFERENCE":          CodeDOCOutdated,
	"MISSING_DIAGRAM":             CodeDOCNoDiagram,
	"SCOPE_CREEP":                 CodeSCOPECreep,
	"UNBOUNDED_SCOPE":             CodeSCOPEUnbounded,
	"MISSING_CONSTRAINTS":         CodeSCOPENoConstr,
	"ACCESSIBILITY_GAP":           CodeUXNoARIA,
	"MISSING_TIMELINE":            CodeSCOPENoTimeline,
}

// NormalizeCode converts legacy codes to new format.
func NormalizeCode(code ReasonCode) ReasonCode {
	if newCode, ok := LegacyCodeMapping[code]; ok {
		return newCode
	}
	return code
}

// GetReasonCodeInfo returns the info for a reason code, or nil if not found.
// Handles both legacy and new code formats.
func GetReasonCodeInfo(code ReasonCode) *ReasonCodeInfo {
	// Try direct lookup first
	if info, ok := ReasonCodeRegistry[code]; ok {
		return &info
	}
	// Try normalized (legacy) lookup
	normalizedCode := NormalizeCode(code)
	if info, ok := ReasonCodeRegistry[normalizedCode]; ok {
		return &info
	}
	return nil
}

// GetReasonCodesByCategory returns all reason codes in a category.
func GetReasonCodesByCategory(category string) []ReasonCode {
	var codes []ReasonCode
	for code, info := range ReasonCodeRegistry {
		if info.Category == category {
			codes = append(codes, code)
		}
	}
	return codes
}

// GetReasonCodesBySpecType returns reason codes applicable to a spec type.
func GetReasonCodesBySpecType(specType string) []ReasonCode {
	var codes []ReasonCode
	for code, info := range ReasonCodeRegistry {
		for _, st := range info.SpecTypes {
			if st == specType {
				codes = append(codes, code)
				break
			}
		}
	}
	return codes
}

// AllReasonCodeCategories returns all unique reason code categories.
func AllReasonCodeCategories() []string {
	return []string{
		CategoryREQ,
		CategoryMETRIC,
		CategoryUSER,
		CategoryARCH,
		CategorySEC,
		CategorySCALE,
		CategoryINFRA,
		CategoryDOC,
		CategorySCOPE,
		CategoryUX,
	}
}

// GetRepairPrompt returns the AI repair prompt for a code.
func GetRepairPrompt(code ReasonCode) string {
	if info := GetReasonCodeInfo(code); info != nil {
		return info.RepairPrompt
	}
	return ""
}

// RequiresHumanReview returns true if the code needs human review after AI repair.
func RequiresHumanReview(code ReasonCode) bool {
	if info := GetReasonCodeInfo(code); info != nil {
		return info.RequiresHuman
	}
	return true // Default to requiring human review for unknown codes
}
