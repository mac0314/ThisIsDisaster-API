package models

import (
	"github.com/revel/revel"
)

type Error struct {
	Id      int64  `db:"error_id, primarykey, autoincrement" json:"id"`
	Title   string `db:"title_mn" json:"title"`
	Content string `db:"log_ln" json:"log"`
	Create  int64  `db:"create_n" json:"create"`
}

func (b *Error) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Content,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
