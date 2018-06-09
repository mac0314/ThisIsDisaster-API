package models

import (
	"github.com/revel/revel"
)

type GameLog struct {
	Id     int64  `db:"game_log_id, primarykey, autoincrement" json:"id"`
	UId    int64  `db:"user_id" json:"uid"`
	Title  string `db:"title_mn" json:"title"`
	Log    string `db:"log_ln" json:"log"`
	Create int64  `db:"create_n" json:"create"`
}

func (b *GameLog) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Log,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
