package btorm

type Controller interface {
	Start() error
	Stop() error
	Status() error
}
