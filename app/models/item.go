package models

import (
    "github.com/revel/revel"
    "regexp"
)

type Item struct {
    Id              int64   `db:"item_id" json:"id"`
    Name            string  `db:"name_sn" json:"name"`
    Type            string  `db:"type_sn" json:"type"`
    Rank            string  `db:"rank_sn" json:"rank"`
    Effect          string  `db:"effect_ln" json:"effect"`
    Resource        string  `db:"resource_mn" json:"resource"`
}

func (b *Item) Validate(v *revel.Validation) {

    v.Check(b.Name,
        revel.ValidRequired(),
        revel.ValidMaxSize(25))

    v.Check(b.Type,
        revel.ValidRequired(),
        revel.ValidMatch(
            regexp.MustCompile(
                "^(normal|equipment|consumable)$")))

    v.Check(b.Rank,
        revel.ValidRequired())

    v.Check(b.Effect,
        revel.ValidRequired())

    v.Check(b.Resource,
        revel.ValidRequired())
}
