package common

type Resourcer interface {
	Create(JSON)
	Read(id int) (JSON, error)
	Update(id int) (JSON, error)
	Delete(id int) (JSON, error)
}
