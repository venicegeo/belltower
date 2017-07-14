package common

type Resourcer interface {
	Create(Json)
	Read(id int) (Json, error)
	Update(id int) (Json, error)
	Delete(id int) (Json, error)
}
