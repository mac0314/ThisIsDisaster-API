package models

import (
	"github.com/revel/revel"
)

type Character struct {
	Id   int64  `db:"u_character_id, primarykey, autoincrement" json:"id"`
	Uid  int64  `db:"user_id" json:"uid"`
	Name string `db:"name_sn" json:"name"`
}

func (b *Character) Validate(v *revel.Validation) {
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))
}
