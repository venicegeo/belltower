package orm

import "fmt"

type ActionRuleAssociation struct {
	Core
	Rule     Rule
	RuleID   uint
	Action   Action
	ActionID uint
}

func (raa ActionRuleAssociation) String() string {
	s := fmt.Sprintf("raa.%d", raa.ID)
	return s
}
