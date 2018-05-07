package models

import (
	"github.com/revel/revel"
)

type Event struct {
	Id      int64  `db:"event_id, primarykey, autoincrement" json:"id"`
	Title   string `db:"title_mn" json:"title"`
	Content string `db:"content_ln" json:"content"`
	Create  int64  `db:"create_n" json:"create"`
}

func (b *Event) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Content,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
