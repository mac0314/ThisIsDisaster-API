package models

import (
	"github.com/revel/revel"
)

type User struct {
	Id       int64  `db:"user_id, primarykey, autoincrement" json:"id"`
	Email    string `db:"email_mn" json:"email"`
	Nickname string `db:"nickname_mn" json:"nickname"`
	Password string `db:"password_ln" json:"password"`
	Level    int64  `db:"level_n" json:"level"`
	Gold     int64  `db:"gold_n" json:"gold"`
	Score    int64  `db:"score_n" json:"score"`
	Signin   int64  `db:"signin_n" json:"signin"`
}

func (b *User) Validate(v *revel.Validation) {
	v.Check(b.Email,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Nickname,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Nickname,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))
}
