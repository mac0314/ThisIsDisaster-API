package models

import (
	"github.com/revel/revel"
)

type Map struct {
	Id        int64  `db:"map_id, primarykey, autoincrement" json:"id"`
	SId       int64  `db:"stage_id" json:"sid"`
	Formation string `db:"formation_ln" json:"formation"`
	Create    int64  `db:"create_n" json:"create"`
}

func (b *Map) Validate(v *revel.Validation) {
	v.Check(b.Formation,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
