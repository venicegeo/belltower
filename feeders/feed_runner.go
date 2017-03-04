package feeders

type StatusF func(mssg string) (bool, error)
type MssgF func(map[string]string) error

type FeedRunner interface {
	ID() uint
	Name() string
	Run(statusF StatusF, mssgF MssgF) error
}
