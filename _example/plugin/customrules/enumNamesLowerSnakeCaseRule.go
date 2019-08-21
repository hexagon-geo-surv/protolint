package customrules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

// EnumNamesLowerSnakeCaseRule verifies that all enum names are LowerSnakeCase.
type EnumNamesLowerSnakeCaseRule struct{}

// NewEnumNamesLowerSnakeCaseRule creates a new EnumNamesLowerSnakeCaseRule.
func NewEnumNamesLowerSnakeCaseRule() EnumNamesLowerSnakeCaseRule {
	return EnumNamesLowerSnakeCaseRule{}
}

// ID returns the ID of this rule.
func (r EnumNamesLowerSnakeCaseRule) ID() string {
	return "ENUM_NAMES_LOWER_SNAKE_CASE"
}

// Purpose returns the purpose of this rule.
func (r EnumNamesLowerSnakeCaseRule) Purpose() string {
	return "Verifies that all enum names are LowerSnakeCase."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumNamesLowerSnakeCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumNamesLowerSnakeCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	// TODO:
	return []report.Failure{
		report.Failuref(meta.Position{}, r.ID(), "Custom Rule"),
	}, nil
}
