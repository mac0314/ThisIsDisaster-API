package models

import (
	"github.com/revel/revel"
)

type HaveHeadCostume struct {
	Id     int64  `db:"have_hc_id, primarykey, autoincrement" json:"id"`
	UId    int64  `db:"user_id" json:"uid"`
	HCId   int64  `db:"head_costume_id" json:"hcid"`
	State  string `db:"state_sn" json:"state"`
	Update int64  `db:"update_n" json:"update"`
}

func (b *HaveHeadCostume) Validate(v *revel.Validation) {
	v.Check(b.State,
		revel.ValidMaxSize(20))
}
