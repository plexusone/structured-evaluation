// Code generated from ../../schema/*.schema.json by scripts/generate.mjs. DO NOT EDIT.
// The Go structs in rubric/, claims/, and summary/ are the source of truth;
// regenerate via `go run ./cmd/genschema schema/` in the repo root, then
// `npm run generate` here.

import { z } from "zod"

export const SummaryReportSchema = z.object({ "$schema": z.string().optional(), "project": z.string().optional(), "version": z.string().optional(), "target": z.string().optional(), "phase": z.string().optional(), "teams": z.array(z.object({ "id": z.string().optional(), "name": z.string().optional(), "agent_id": z.string().optional(), "model": z.string().optional(), "depends_on": z.array(z.string()).optional(), "tasks": z.array(z.object({ "id": z.string().optional(), "status": z.string().optional(), "detail": z.string().optional(), "duration_ms": z.number().int().optional(), "metadata": z.record(z.string(), z.any()).optional() }).strict()).optional(), "status": z.string().optional() }).strict()).optional(), "status": z.string().optional(), "generated_at": z.string().datetime({ offset: true }).optional(), "generated_by": z.string().optional(), "embeddedReports": z.object({ "evaluations": z.record(z.string(), z.any()).optional(), "claims": z.record(z.string(), z.any()).optional(), "custom": z.record(z.string(), z.any()).optional() }).strict().optional() }).strict().describe("Schema for GO/NO-GO summary reports from deterministic checks")
export type SummaryReport = z.infer<typeof SummaryReportSchema>

