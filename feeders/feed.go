package drivers

type Config map[string]string
type StatusF func(mssg string) (bool, error)
type MssgF func(map[string]string) error

type Feed interface {
	Run(map[string]string) error
}
