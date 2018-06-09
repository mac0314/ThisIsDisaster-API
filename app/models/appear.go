package models

import (
	"github.com/revel/revel"
)

type Appear struct {
	Id     int64  `db:"appear_id, primarykey, autoincrement" json:"id"`
	SId    int64  `db:"stage_id" json:"sid"`
	MId    int64  `db:"monster_id" json:"mid"`
	State  string `db:"state_sn" json:"state"`
	Create int64  `db:"create_n" json:"create"`
}

func (b *Appear) Validate(v *revel.Validation) {
	v.Check(b.State,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
}
