#!/usr/bin/env node
// Generates Zod schemas (src/generated/*.ts) from the Go-first canonical
// JSON Schema in ../schema/*.schema.json. Never hand-edit src/generated/ —
// re-run `npm run generate` after the Go structs (and their regenerated
// JSON Schema) change.
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import { jsonSchemaToZod } from 'json-schema-to-zod'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const schemaDir = path.resolve(__dirname, '../../schema')
const outDir = path.resolve(__dirname, '../src/generated')

mkdirSync(outDir, { recursive: true })

const targets = [
  { file: 'rubric.schema.json', name: 'Rubric', out: 'rubric.ts' },
  { file: 'rubricset.schema.json', name: 'RubricSet', out: 'rubricset.ts' },
  { file: 'claims.schema.json', name: 'ClaimsReport', out: 'claims.ts' },
  { file: 'summary.schema.json', name: 'SummaryReport', out: 'summary.ts' },
]

const banner = `// Code generated from ../../schema/*.schema.json by scripts/generate.mjs. DO NOT EDIT.
// The Go structs in rubric/, claims/, and summary/ are the source of truth;
// regenerate via \`go run ./cmd/genschema schema/\` in the repo root, then
// \`npm run generate\` here.
`

for (const t of targets) {
  const schemaPath = path.join(schemaDir, t.file)
  const schema = JSON.parse(readFileSync(schemaPath, 'utf-8'))

  const zodSource = jsonSchemaToZod(schema, {
    name: `${t.name}Schema`,
    module: 'esm',
    type: t.name,
  })

  const outPath = path.join(outDir, t.out)
  writeFileSync(outPath, banner + '\n' + zodSource + '\n')
  console.log(`Generated ${path.relative(process.cwd(), outPath)}`)
}
