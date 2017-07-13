package engine

type Component interface {
	Init(config interface{}) error

	Run(in interface{}) (out interface{}, err error)
}
