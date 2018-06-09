package models

import (
	"github.com/revel/revel"
)

type Award struct {
	Id     int64 `db:"award_id, primarykey, autoincrement" json:"id"`
	SId    int64 `db:"stage_id" json:"sid"`
	UId    int64 `db:"user_id" json:"uid"`
	Exp    int64 `db:"exp_n" json:"exp"`
	Gold   int64 `db:"gold_n" json:"gold"`
	Score  int64 `db:"score_n" json:"score"`
	Star   int64 `db:"star_n" json:"star"`
	Create int64 `db:"create_n" json:"create"`
}

func (b *Award) Validate(v *revel.Validation) {

}
