package orm

import "fmt"

type FeedType struct {
	ModelCore
	Name    string
	Owner   User
	OwnerID uint
}

type FeedTypeAcl struct {
	ModelCore
	UserID     uint
	FeedTypeID uint
}

func (ft FeedType) String() string {
	s := fmt.Sprintf("FT%d: \"%s\"\n", ft.ID, ft.Name)
	s += fmt.Sprintf("    owner: U%d\n", ft.OwnerID)
	return s
}

func (acl FeedTypeAcl) String() string {
	s := fmt.Sprintf("A%d: FT%d U%d\n",
		acl.ID, acl.FeedTypeID, acl.UserID)
	return s
}
