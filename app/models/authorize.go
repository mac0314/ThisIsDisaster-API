package models

import (
	"github.com/revel/revel"
)

type Authorize struct {
	Id       int64  `db:"auth_id, primarykey, autoincrement" json:"id"`
	UId      int64  `db:"user_id" json:"uid"`
	Email    string `db:"email_mn" json:"email"`
	Platform string `db:"platform_sn" json:"platform"`
	Create   int64  `db:"create_n" json:"create"`
}

func (b *Authorize) Validate(v *revel.Validation) {
	v.Check(b.Email,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Platform,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
}
