package models

import (
	"github.com/revel/revel"
)

type Media struct {
	Id        int64  `db:"media_id, primarykey, autoincrement" json:"id"`
	Title     string `db:"title_mn" json:"title"`
	Content   string `db:"content_ln" json:"content"`
	Thumbnail string `db:"thumbnail_ln" json:"thumbnail"`
	Image     string `db:"image_ln" json:"image"`
	Create    int64  `db:"create_n" json:"create"`
}

func (b *Media) Validate(v *revel.Validation) {
	v.Check(b.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))

	v.Check(b.Content,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))

	v.Check(b.Thumbnail,
		revel.ValidMaxSize(255))

	v.Check(b.Image,
		revel.ValidMaxSize(255))
}
