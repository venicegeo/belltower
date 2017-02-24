package orm

import "fmt"

type FeedRuleAssociation struct {
	Core
	Rule   Rule
	RuleID uint
	Feed   Feed
	FeedID uint
}

func (fra FeedRuleAssociation) String() string {
	s := fmt.Sprintf("A%d: F%d\n",
		fra.ID, fra.FeedID)
	return s
}
