package models

import (
	"github.com/revel/revel"
)

type Apply struct {
	Id     int64  `db:"apply_id, primarykey, autoincrement" json:"id"`
	FId    int64  `db:"feedback_id" json:"fid"`
	Title  string `db:"title_mn" json:"title"`
	Log    string `db:"log_ln" json:"log"`
	Create int64  `db:"create_n" json:"create"`
}

func (b *Apply) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Log,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
