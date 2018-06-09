package models

import (
	"github.com/revel/revel"
)

type Notification struct {
	Id        int64  `db:"notification_id, primarykey, autoincrement" json:"id"`
	UId       int64  `db:"user_id" json:"uid"`
	Title     string `db:"title_mn" json:"title"`
	Content   string `db:"content_ln" json:"content"`
	ReadCheck bool   `db:"read_check" json:"readCheck"`
	Create    int64  `db:"create_n" json:"create"`
}

func (b *Notification) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Content,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
