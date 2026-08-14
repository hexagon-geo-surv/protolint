---
name: release
description: Cut a protolint release — preflight checks, tagging, and verifying that all six distribution channels actually published. Use when asked to release, publish, ship a version, or when a release partially failed and needs diagnosing.
---

# Releasing protolint

One git tag fans out to six channels. `release.sh` only pushes the tag; everything
else is CI. The job is mostly **verifying that each channel really published**,
because they fail independently and silently.

## 1. Preflight

Run from a clean `master` that is in sync with `origin`:

```bash
go build ./... && go test ./...
(cd bdist/js && npm ci --ignore-scripts && npx tsc && npx eslint . && npx prettier --check .)
```

`.github/workflows/go.yml` and the JS workflows only run on `pull_request`, so
nothing has verified `master` as a whole. Do not skip this.

## 2. Pick the version

Every channel takes its version from the git tag — `.goreleaser.yml` ldflags,
`build.gradle` from `GITHUB_REF_NAME`, and `npm version from-git`. Nothing needs
bumping by hand; the tag *is* the version.

Look at what changed since the last tag, and remember the npm package ships to a
different audience than the Go binary:

- **patch** — Go-side fixes, dependency bumps, docs
- **minor** — new lint rules, or anything that changes what npm consumers need:
  a raised `engines.node`, a module-system switch, renamed install scripts

## 3. Tag

```bash
./release.sh vX.Y.Z "Release vX.Y.Z"
```

## 4. Watch the pipeline

`goreleaser` runs on the tag push. `publish_npm` and `publish_wheel` then trigger
off its **completion, not its success** — so a green goreleaser does not mean the
downstream jobs passed.

```bash
gh run watch <goreleaser-run-id> --exit-status
gh run list --workflow publish_npm.yml --limit 1
gh run list --workflow publish_wheel.yml --limit 1
```

## 5. Verify every channel

Check the registries, not the workflow badges.

```bash
gh release view vX.Y.Z --json assets --jq '.assets|length'          # expect 9
npm view protolint version
curl -s https://pypi.org/pypi/protolint-bin/json | jq -r .info.version
curl -s https://repo1.maven.org/maven2/io/github/yoheimuta/protoc-gen-protolint/maven-metadata.xml | grep release
```

Maven Central publishes automatically via the Central Publisher Portal step in
`goreleaser.yml` (`publishing_type=automatic`). No manual staging promotion.

Then prove the npm package actually works end to end — this exercises
`postinstall.js`, which downloads the release binary and is the part most likely
to break:

```bash
cd "$(mktemp -d)" && npm init -y >/dev/null && npm install protolint@X.Y.Z
./node_modules/.bin/protolint version          # must print X.Y.Z
./node_modules/.bin/protoc-gen-protolint version
```

## Troubleshooting

**`npm error 404 ... PUT https://registry.npmjs.org/protolint`**
Authentication, not a missing package. npm publishes through
[trusted publishing](https://docs.npmjs.com/trusted-publishers) — there is no
token. Check the trusted publisher on npmjs.com still points at `yoheimuta` /
`protolint` / `publish_npm.yml`. **Renaming `publish_npm.yml` breaks publishing**,
because the registration matches on the workflow filename.

**`fatal: No tags can describe '<sha>'` in *Set package version***
`npm version from-git` needs history back to the tag. `publish_npm.yml` uses
`fetch-depth: 0` and checks the tag out explicitly for this reason; a shallow
checkout reintroduces the failure.

**`npm warn publish "bin[protolint]" script name ... was invalid and removed`**
Harmless. npm only strips the `./` prefix; the published `bin` entries survive.
Confirm with `npm view protolint@X.Y.Z bin` if in doubt.

**A channel failed after the tag is already public**
Do not re-tag. Fix the workflow, merge, then re-run the failed workflow with
`gh workflow run <name>.yml --ref master`. `publish_npm` and `publish_wheel` both
have `workflow_dispatch`, and `publish_npm` checks out the latest tag, so a re-run
publishes the tagged tree even after later merges have landed on `master`.

## After the release

Post the version on any issue that asked for it, and close it.
