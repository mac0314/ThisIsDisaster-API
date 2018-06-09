package models

import (
	"github.com/revel/revel"
)

type Have struct {
	Id     int64  `db:"have_id, primarykey, autoincrement" json:"id"`
	SId    int64  `db:"stage_id" json:"sid"`
	IId    int64  `db:"item_id" json:"iid"`
	UId    int64  `db:"user_id" json:"uid"`
	State  string `db:"state_sn" json:"state"`
	Create int64  `db:"create_n" json:"create"`
}

func (b *Have) Validate(v *revel.Validation) {
	v.Check(b.State,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
}
