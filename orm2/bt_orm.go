package orm2

import "github.com/venicegeo/belltower/common"

//---------------------------------------------------------------------

type BtOrm struct {
	Orm *Orm
}

func (orm *BtOrm) CreateAction(requestorID common.Ident, fields *ActionFieldsForCreate) (common.Ident, error) {
	action := &Action{}
	return orm.Orm.CreateThing(requestorID, action, fields)
}
