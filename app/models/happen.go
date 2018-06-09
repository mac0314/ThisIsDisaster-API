package models

import (
	"github.com/revel/revel"
)

type Happen struct {
	Id     int64 `db:"happen_id, primarykey, autoincrement" json:"id"`
	DId    int64 `db:"disaster_id" json:"did"`
	SId    int64 `db:"stage_id" json:"sid"`
	Create int64 `db:"create_n" json:"create"`
}

func (b *Happen) Validate(v *revel.Validation) {

}
