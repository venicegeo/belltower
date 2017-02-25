package orm

import "fmt"

type FeedType struct {
	Core
	FeedTypeAttributes
}

type FeedTypeAttributes struct {
	Name       string
	ConfigInfo string
	IsEnabled  bool
}

func (ft FeedType) String() string {
	s := fmt.Sprintf("ft.%d: %s", ft.ID, ft.Name)
	return s
}
