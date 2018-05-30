package models

import (
	"github.com/revel/revel"
)

type User struct {
	Id         int64  `db:"user_id, primarykey, autoincrement" json:"id"`
	Email      string `db:"email_mn" json:"email"`
	Password   string `db:"password_ln" json:"password"`
	Nickname   string `db:"nickname_mn" json:"nickname"`
	FacebookId int64  `db:"facebook_id" json:"facebookId"`
	GoogleId   int64  `db:"google_id" json:"googleId"`
	GithubId   int64  `db:"github_id" json:"githubId"`
	Avatar     string `db:"avatar_ln" json:"avatar"`
	Slug       string `db:"Slug_ln" json:"slug"`
	Level      int64  `db:"level_n" json:"level"`
	Gold       int64  `db:"gold_n" json:"gold"`
	Score      int64  `db:"score_n" json:"score"`
	Signin     int64  `db:"signin_n" json:"signin"`
}

func (b *User) Validate(v *revel.Validation) {
	v.Check(b.Email,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Password,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))

	v.Check(b.Nickname,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))

	v.Check(b.Avatar,
		revel.ValidRequired(),
		revel.ValidMaxSize(255))

	v.Check(b.Slug,
		revel.ValidRequired(),
		revel.ValidMaxSize(128))
}
