package models

import (
	"github.com/revel/revel"
)

type HeadCostume struct {
	Id       int64  `db:"head_costume_id, primarykey, autoincrement" json:"id"`
	Name     string `db:"name_mn" json:"name"`
	Resource string `db:"resource_mn" json:"resource"`
}

func (b *HeadCostume) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Resource,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))
}
