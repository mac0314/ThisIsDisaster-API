package models

import (
	"regexp"

	"github.com/revel/revel"
)

type Item struct {
	Id       int64  `db:"item_id, primarykey, autoincrement" json:"id" xml:"id"`
	Name     string `db:"name_sn" json:"name" xml:"name"`
	Type     string `db:"type_sn" json:"type" xml:"type"`
	Rank     string `db:"rank_sn" json:"rank" xml:"rank"`
	Effect   string `db:"effect_ln" json:"effect" xml:"effect"`
	Resource string `db:"resource_mn" json:"resource" xml:"resource"`
}

func (b *Item) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))

	v.Check(b.Type,
		revel.ValidRequired(),
		revel.ValidMaxSize(20),
		revel.ValidMatch(
			regexp.MustCompile(
				"^(normal|equipment|consumable)$")))

	v.Check(b.Rank,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))

	v.Check(b.Effect,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))

	v.Check(b.Resource,
		revel.ValidRequired(),
		revel.ValidMaxSize(50))
}
