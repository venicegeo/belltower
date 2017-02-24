package orm

import "fmt"

type FeedType struct {
	Core
	Name       string
	ConfigInfo string
}

func (ft FeedType) String() string {
	s := fmt.Sprintf("ft.%d: %s", ft.ID, ft.Name)
	return s
}
