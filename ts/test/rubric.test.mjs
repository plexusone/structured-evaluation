import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import { RubricSchema } from '../dist/index.js'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const fixture = JSON.parse(
  readFileSync(path.join(__dirname, 'fixtures/rubric-report.json'), 'utf-8')
)

test('RubricSchema parses real Go-generated output', () => {
  const result = RubricSchema.safeParse(fixture)
  assert.equal(result.success, true, result.success ? '' : JSON.stringify(result.error?.issues))
})

test('the 5-level intScore round-trips correctly', () => {
  const result = RubricSchema.parse(fixture)
  assert.equal(result.intScore, 5)
  assert.equal(result.categories[0].intScore, 5)
})

test('category severity (computed from findings) round-trips', () => {
  const result = RubricSchema.parse(fixture)
  assert.equal(result.categories[0].severity, 'low')
})

test('decision enum values are strictly validated, not free-form strings', () => {
  // A report that skipped Finalize()/Evaluate() has decision.status/
  // overallDecision as empty strings, which are not valid enum members.
  // Go's json.Unmarshal into a plain string field would accept this
  // silently; Zod correctly rejects it.
  const unfinalized = {
    ...fixture,
    decision: { ...fixture.decision, status: '' },
    overallDecision: '',
  }
  const result = RubricSchema.safeParse(unfinalized)
  assert.equal(result.success, false)
})

test('rejects an eval report with the wrong field name (regression guard)', () => {
  // This is the exact class of bug this package exists to prevent: a
  // hand-maintained consumer type using the wrong field name for the score
  // (visionstudio's api.EvalResult used "scoreV2"; the real field is
  // "intScore"). z.object().strict() rejects unrecognized keys outright —
  // a loud parse failure instead of a silently-undefined field.
  const wrongShape = { ...fixture, scoreV2: fixture.intScore }
  delete wrongShape.intScore

  const result = RubricSchema.safeParse(wrongShape)
  assert.equal(result.success, false)
  assert.ok(
    result.error.issues.some((i) => i.code === 'unrecognized_keys' && i.keys.includes('scoreV2')),
    `expected an unrecognized_keys error for "scoreV2", got: ${JSON.stringify(result.error.issues)}`
  )
})
