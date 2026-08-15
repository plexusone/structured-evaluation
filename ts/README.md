# @plexusone/structured-evaluation

Zod schemas and TypeScript types for `structured-evaluation`'s LLM-as-Judge
report types — `Rubric`, `RubricSet`, `ClaimsReport`, `SummaryReport` —
generated from the same Go structs' canonical JSON Schema the Go library
itself embeds. The Go structs are the source of truth; this package is
downstream of them, never hand-maintained.

## Why this exists

A hand-maintained TypeScript mirror of a Go report type drifts silently: a
renamed field (`intScore` vs. the guessed `scoreV2`), a type mismatch
(`score` is a `"pass"|"partial"|"fail"` string, not an int), or a restructured
field (`decision` became an object, not a bare string) compiles fine on both
sides and fails at runtime — or worse, doesn't fail at all, it just silently
produces `undefined`. This package makes that class of bug a `zod.parse()`
error instead.

## Usage

```bash
npm install @plexusone/structured-evaluation
```

```ts
import { RubricSchema, type Rubric } from '@plexusone/structured-evaluation'

const report: Rubric = RubricSchema.parse(JSON.parse(rawJson))
console.log(report.intScore) // the 1-5 score, correctly typed
console.log(report.categories[0].severity) // "critical" | "high" | "medium" | "low" | "info" | undefined
```

Every object schema is `.strict()` — an unrecognized key (e.g. a stale
consumer still expecting a field that was renamed upstream) fails parsing
loudly instead of being silently dropped.

### Claims

```ts
import { ClaimsReportSchema, type ClaimsReport } from '@plexusone/structured-evaluation'

const report: ClaimsReport = ClaimsReportSchema.parse(JSON.parse(rawJson))

for (const claim of report.claims ?? []) {
  if (claim.validation?.external?.sourceType === 'aggregator') {
    // Sourced from a stats-roundup site with no original reporting —
    // rejected by default even if the excerpt matches verbatim.
    continue
  }
  if (claim.statistical) {
    // Structured value/unit/precision/as-of-date, independent of the
    // rendered claim.text (e.g. "4.7M paid subscribers").
    console.log(claim.statistical.value, claim.statistical.unit, claim.statistical.precision)
  }
}
```

## Regenerating

```bash
# 1. From the repo root: regenerate the Go-side JSON Schema after any
#    change to rubric.Rubric, rubric.RubricSet, claims.ClaimsReport, or
#    summary.SummaryReport.
go run ./cmd/genschema schema/

# 2. From ts/: regenerate the Zod schemas from that JSON Schema.
npm run generate

# 3. Rebuild and verify against real Go-generated output.
npm test
```

Never hand-edit `src/generated/*.ts` — it's regenerated wholesale from
`../schema/*.schema.json` by `scripts/generate.mjs`.

## Publishing

**A `structured-evaluation` release that changes an exported Go type is not
finished until this package ships the matching Zod.** npm consumers pin a
version; if the Go module gains a field and npm doesn't, their
`ClaimsReportSchema.parse()` silently drops it (or `.strict()` rejects it) —
the exact drift this package exists to prevent, now reintroduced at the
package boundary. Publish npm in lockstep with the Git tag.

> npm versions need not be contiguous — if a release was skipped on npm
> (e.g. v0.13.0 was never published), just publish the current version. The
> latest package always reflects the current schema; there's no need to
> backfill intermediate versions.

There is **no `prepublishOnly`/`prepack` hook**, so `pnpm publish` ships
whatever is already in `dist/` — you must regenerate and build first.

```bash
# 0. Authenticate — pnpm reads npm's ~/.npmrc token; needs @plexusone
#    publish rights. Verify with `npm whoami`.
npm login

# 1. Regenerate from the current JSON Schema and rebuild dist/.
pnpm run generate
pnpm run build
pnpm test              # build + parse real Go-generated fixtures

# 2. Bump the version to match the Go module (see Versioning below), then:
pnpm publish --access public
```

Notes:

- **Build before publish.** With no prepublish hook, a stale `dist/` would
  ship silently. `pnpm publish --dry-run` prints the exact tarball contents
  (should be `dist/` + `package.json` only) — eyeball it first.
- **Git checks.** `pnpm publish` verifies a clean tree on the publish branch
  by default; publish from the tagged commit. Append `--no-git-checks` only
  if you must override.
- **Lockfile.** This package is npm-managed (`package-lock.json`, no
  `pnpm-lock.yaml`), so `pnpm install --frozen-lockfile` fails — use plain
  `pnpm install` (or rely on the existing `node_modules`).

## Known limitation

The JSON Schema currently declares no `required` fields (no Go struct has a
`jsonschema:"required"` tag), so every field here is `.optional()` — a
missing field validates successfully instead of failing. This does **not**
affect the bug class this package exists to prevent (wrong names, wrong
types, restructured shapes — all caught), but it means an accidentally
omitted field won't be caught either. See `test/rubric.test.mjs`'s
regression-guard test for what's covered today.

## Versioning

This package's version tracks the Go module's version — use the same
version number here as the `structured-evaluation` Go module you're
consuming reports from (e.g. Go v0.12.0 → `@plexusone/structured-evaluation`
v0.12.0).
