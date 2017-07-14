package orm

type Controller interface {
	Start() error
	Stop() error
	Status() error
}
