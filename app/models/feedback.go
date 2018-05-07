package models

import (
	"github.com/revel/revel"
)

type Feedback struct {
	Id      int64  `db:"feedback_id, primarykey, autoincrement" json:"id"`
	Email   string `db:"email_mn" json:"email"`
	Title   string `db:"title_mn" json:"title"`
	Content string `db:"content_ln" json:"content"`
	Create  int64  `db:"create_n" json:"create"`
}

func (b *Feedback) Validate(v *revel.Validation) {
	v.Check(b.Email,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Content,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
