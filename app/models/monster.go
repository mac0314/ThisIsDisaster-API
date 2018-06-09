package models

import (
	"github.com/revel/revel"
)

type Monster struct {
	Id       int64  `db:"monster_id, primarykey, autoincrement" json:"id" xml:"id"`
	Name     string `db:"name_mn" json:"name" xml:"name"`
	Info     string `db:"info_ln" json:"info" xml:"info"`
	Health   int64  `db:"health_n" json:"health" xml:"health"`
	Resource string `db:"resource_mn" json:"resource" xml:"resource"`
	Create   int64  `db:"create_n" json:"create" xml:"create"`
}

func (b *Monster) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Info,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))

	v.Check(b.Resource,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))
}
