package models

import (
	"github.com/revel/revel"
)

type Achievement struct {
	Id      int64  `db:"achievement_id, primarykey, autoincrement" json:"id"`
	Title   string `db:"title_mn" json:"title"`
	Content string `db:"content_ln" json:"content"`
	Score   int64  `db:"score_n" json:"score"`
}

func (b *Achievement) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Content,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
