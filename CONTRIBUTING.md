# Contributing

- Fork it
- Create your feature branch: git checkout -b your-new-feature
- Commit changes: git commit -m 'Add your feature'
- Pass all tests
- Push to the branch: git push origin your-new-feature
- Submit a pull request


## Releasing

`./release.sh <version>` pushes a tag; CI does the rest. A single tag publishes to
GitHub Releases, npm, PyPI, Maven Central, Docker Hub and Homebrew, and every
version number is derived from that tag — nothing is bumped by hand.

The Maven artifacts no longer need to be promoted manually: `goreleaser.yml`
uploads them to the Central Publisher Portal with `publishing_type=automatic`, so
they reach Maven Central without a staging step.

npm publishes through [trusted publishing](https://docs.npmjs.com/trusted-publishers)
rather than a stored token. The registration on npmjs.com matches on the workflow
**filename**, so renaming `.github/workflows/publish_npm.yml` breaks publishing
until the trusted publisher is updated to match.

See [`.claude/skills/release/SKILL.md`](.claude/skills/release/SKILL.md) for the
full runbook, including how to verify each channel and how to recover when one of
them fails after the tag is already public.

## Adding a lint rule

A rule has to be registered and documented, not just implemented, or users will
never see it. [`.claude/skills/add-lint-rule/SKILL.md`](.claude/skills/add-lint-rule/SKILL.md)
lists every file that has to change and explains how the official/unofficial,
fixable and auto-disable behaviours are wired.
