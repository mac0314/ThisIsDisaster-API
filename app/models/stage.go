package models

import (
	"github.com/revel/revel"
)

type Stage struct {
	Id     int64  `db:"stage_id, primarykey, autoincrement" json:"id"`
	Name   string `db:"name_sn" json:"name"`
	Mode   string `db:"mode_sn" json:"mode"`
	Create int64  `db:"create_n" json:"create"`
}

func (b *Stage) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Mode,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
}
