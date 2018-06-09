package models

import (
	"regexp"

	"github.com/revel/revel"
)

type Item struct {
	Id          int64  `db:"item_id, primarykey, autoincrement" json:"id" xml:"id"`
	Type        string `db:"type_sn" json:"type" xml:"type"`
	Name        string `db:"name_sn" json:"name" xml:"name"`
	Description string `db:"description_ln" json:"description" xml:"description"`
	Damage      int64  `db:"damage_n" json:"damage" xml:"damage"`
	RangeX      int64  `db:"range_x" json:"rangeX" xml:"range_x"`
	RangeY      int64  `db:"range_y" json:"rangeY" xml:"range_y"`
	Resource    string `db:"resource_mn" json:"resource" xml:"sprite"`
	Create      int64  `db:"create_n" json:"create" xml:"create"`
}

func (b *Item) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))

	v.Check(b.Type,
		revel.ValidMaxSize(20),
		revel.ValidRequired(),
		revel.ValidMatch(
			regexp.MustCompile(
				"^(normal|weapon|consumable)$")))

	v.Check(b.Description,
		revel.ValidMaxSize(255))

	v.Check(b.Resource,
		revel.ValidMaxSize(50))
}
