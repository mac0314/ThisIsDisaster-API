package models

import (
	"github.com/revel/revel"
)

type Movement struct {
	Id     int64   `db:"movement_id, primarykey, autoincrement" json:"id"`
	SId    int64   `db:"stage_id" json:"sid"`
	UId    int64   `db:"user_id" json:"uid"`
	X      float64 `db:"x_f" json:"x"`
	Y      float64 `db:"y_f" json:"y"`
	Z      float64 `db:"z_f" json:"z"`
	Create int64   `db:"create_n" json:"create"`
}

func (b *Movement) Validate(v *revel.Validation) {

}
