package models

import (
	"github.com/revel/revel"
)

type Handle struct {
	Id     int64  `db:"handle_id, primarykey, autoincrement" json:"id"`
	EId    int64  `db:"error_id" json:"eid"`
	Title  string `db:"title_mn" json:"title"`
	Log    string `db:"log_ln" json:"log"`
	Create int64  `db:"create_n" json:"create"`
}

func (b *Handle) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Log,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
