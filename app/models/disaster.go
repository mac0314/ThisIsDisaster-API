package models

import (
	"github.com/revel/revel"
)

type Disaster struct {
	Id     int64  `db:"disaster_id, primarykey, autoincrement" json:"id" xml:"id"`
	Name   string `db:"name_mn" json:"name" xml:"name"`
	Info   string `db:"info_ln" json:"info" xml:"info"`
	Create int64  `db:"create_n" json:"create" xml:"create"`
}

func (b *Disaster) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Info,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
