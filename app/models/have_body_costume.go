package models

import (
	"github.com/revel/revel"
)

type HaveBodyCostume struct {
	Id     int64  `db:"have_bc_id, primarykey, autoincrement" json:"id"`
	UId    int64  `db:"user_id" json:"uid"`
	BCId   int64  `db:"body_costume_id" json:"bcid"`
	State  string `db:"state_sn" json:"state"`
	Update int64  `db:"update_n" json:"update"`
}

func (b *HaveBodyCostume) Validate(v *revel.Validation) {
	v.Check(b.State,
		revel.ValidMaxSize(20))
}
