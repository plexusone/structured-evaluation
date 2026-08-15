# Releasing

This is the maintainer runbook for cutting a `structured-evaluation`
release. It covers the Go module, the JSON Schema, the generated
TypeScript/Zod package, and the docs — all of which ship together.

## The core rule: one source of truth, three artifacts

The Go structs in `rubric/`, `claims/`, and `summary/` are the **single
source of truth**. Two artifacts are generated downstream of them and must
never lag:

```
Go structs ──► JSON Schema (schema/*.schema.json, via cmd/genschema)
           └─► Zod schemas + TS types (ts/src/generated/*, via npm run generate)
```

A release that changes an exported Go type but ships stale schema or stale
Zod is a **drift bug at the package boundary** — the exact class of failure
the generated `@plexusone/structured-evaluation` npm package exists to
prevent. If the Go module gains a field and npm doesn't publish the matching
Zod, a TS consumer's `ClaimsReportSchema.parse()` silently drops that field
(or, because every object schema is `.strict()`, rejects the payload). Keep
all three in lockstep.

!!! warning "Publishing npm is part of the release, not an afterthought"
    v0.13.0 changed several `claims` types (`sourceRole`,
    `minCorroboratingSources`, `corroborationCategories`, `maxClaimAge`) but
    the npm package was never republished — it sat at v0.12.0, silently
    missing all of them, until v0.14.0 caught it up. Don't repeat that:
    publish npm with every release that touches an exported type.

## Release checklist

1. **Regenerate the JSON Schema** after any exported-type change in
   `rubric`, `claims`, or `summary`:

    ```bash
    go run ./cmd/genschema schema/
    ```

    Review the diff — it should be additive/expected. Unrelated schema drift
    signals an accidental type change elsewhere.

2. **Regenerate the Zod bindings** and bump the package version to match the
   Go module (e.g. Go `v0.14.0` → `ts/package.json` `0.14.0`):

    ```bash
    cd ts
    pnpm run generate
    pnpm test              # rebuilds dist/ and parses real Go-generated fixtures
    ```

3. **Update the changelog.** Edit `CHANGELOG.json` (compact single-line leaf
   objects — the existing style, not `json.dump(indent=2)`), then regenerate
   the Markdown:

    ```bash
    schangelog validate CHANGELOG.json
    schangelog generate CHANGELOG.json -o CHANGELOG.md
    ```

4. **Add release notes.** Create `docs/releases/vX.Y.Z.md` (follow the most
   recent release's structure) and add it to the `Releases` nav in
   `mkdocs.yml`.

5. **Verify.** `go test ./...`, `golangci-lint run`, `mkdocs build --strict`,
   and `pnpm test` in `ts/` all clean.

6. **Push, then tag.** Push commits to `origin` first and wait for CI to go
   green — some checks (cross-platform builds, integration) only run in CI.
   Only then:

    ```bash
    git push origin main
    # ...CI green...
    git tag vX.Y.Z
    git push origin vX.Y.Z
    ```

7. **Publish npm** (see below), in lockstep with the tag.

## Publishing the npm package

The package is `@plexusone/structured-evaluation`, published under the
`@plexusone` scope with public access. There is **no `prepublishOnly` /
`prepack` hook**, so `pnpm publish` ships whatever is already in `dist/` —
regenerate and build first.

```bash
# 0. Authenticate. pnpm reads npm's ~/.npmrc token; the account needs
#    publish rights on the @plexusone scope. Verify with `npm whoami`.
npm login

cd ts

# 1. Regenerate from the current JSON Schema and rebuild dist/.
pnpm run generate
pnpm run build

# 2. Sanity-check the tarball — expect dist/ + package.json only.
pnpm publish --dry-run

# 3. Publish.
pnpm publish --access public
```

### Gotchas

- **Build before publish.** With no prepublish hook, a stale `dist/` ships
  silently. The `--dry-run` in step 2 prints the exact file list — check it.
- **Git checks.** `pnpm publish` verifies a clean working tree on the publish
  branch by default; publish from the tagged commit. Use `--no-git-checks`
  only to deliberately override.
- **Lockfile.** This package is npm-managed (`package-lock.json`, no
  `pnpm-lock.yaml`), so `pnpm install --frozen-lockfile` fails — use plain
  `pnpm install`, or rely on the existing `node_modules`.
- **Non-contiguous versions are fine.** npm does not require every version to
  exist. If a release was skipped on npm, just publish the current one — the
  latest package always reflects the current schema. Don't backfill.
- **Publishing is effectively permanent.** npm unpublish is restricted after
  72 hours. Dry-run first; publish deliberately.

## See also

- [`ts/README.md`](https://github.com/plexusone/structured-evaluation/blob/main/ts/README.md#publishing)
  — the package's own regeneration and publishing notes.
- The repo `CLAUDE.md` `Conventions` section — the condensed release
  checklist for automated agents.
