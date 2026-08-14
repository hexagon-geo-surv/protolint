---
name: add-lint-rule
description: Add a new lint rule to protolint — every file that has to change, and how the official/unofficial, fixable and auto-disable behaviours are wired. Use when asked to add, implement, or write a new rule.
---

# Adding a lint rule

A rule is not done when `Apply` works — it has to be registered, made
configurable, and documented, or it is invisible to users. Copy the closest
existing rule rather than starting from scratch.

## Files to touch

| File | Always? | What goes in |
|---|---|---|
| `internal/addon/rules/<camelCase>Rule.go` | yes | the rule |
| `internal/addon/rules/<camelCase>Rule_test.go` | yes | table-driven tests |
| `internal/cmd/subcmds/rules.go` | yes | register in `newAllInternalRules` |
| `internal/linter/config/rulesOption.go` | yes | one field on `RulesOption` |
| `README.md` | yes | a row in the `## Rules` table |
| `_example/config/.protolint.yaml` | yes | a commented example entry |
| `internal/linter/config/<name>Option.go` (+ test) | only if configurable | a dedicated option type |
| `_testdata/rules/<name>/*.proto` | when fixtures help | input protos |

Naming is mechanical: rule ID `FIELD_NUMBERS_ORDER_ASCENDING` → type
`FieldNumbersOrderAscendingRule` → file `fieldNumbersOrderAscendingRule.go` →
config key `field_numbers_order_ascending`.

## The rule type

Embed `RuleWithSeverity` and implement:

- `ID() string` — `UPPER_SNAKE_CASE`, unique
- `Purpose() string` — one sentence; it is user-facing and goes in the README verbatim
- `IsOfficial() bool` — **this decides whether the rule is on by default**

`Rules.Default()` in `internal/linter/rule/rules.go` keeps only rules where
`IsOfficial()` is true. Return `true` only for rules from the
[official style guide](https://protobuf.dev/programming-guides/style/); anything
opinionated returns `false` and users opt in via `.protolint.yaml`.

`Apply(proto *parser.Proto) ([]report.Failure, error)` walks the AST with a
visitor from `linter/visitor`.

## Fixable and auto-disable

Both are opt-in and both show up as columns in the README table.

- **Fixable** — build the visitor with `visitor.NewBaseFixableVisitor` and take a
  `fixMode bool` in the constructor. Users get it via `-fix`.
- **AutoDisable** — take `autoDisableType autodisable.PlacementType`. The
  constructor must force `fixMode = false` when `autoDisableType != autodisable.Noop`;
  the two cannot both apply. Users get it via `-auto_disable`.

Rules whose fix cannot break wire compatibility are marked `*1` in the README
instead of supporting auto-disable — for those, fixing is always safe.

## Configuration

If severity is the only knob, use the shared `CustomizableSeverityOption` in
`rulesOption.go`. Add a dedicated option type only when the rule needs its own
settings, and give it a test — these types are parsed from YAML, JSON and TOML,
so all three struct tags must be present and identical.

## Verify

```bash
make test/lint && make test
make run/cmd/protolint ARG="lint _example/proto"
```

Confirm the rule appears where users will look for it:

```bash
go run cmd/protolint/main.go list | grep YOUR_RULE_ID
```

A rule missing from that list was implemented but never registered.
