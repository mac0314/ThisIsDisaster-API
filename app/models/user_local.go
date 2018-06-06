package models

import (
	"github.com/revel/revel"
)

type UserLocal struct {
	Id    int64  `db:"user_id, primarykey, autoincrement" json:"id"`
	Email string `db:"email_mn" json:"email"`
	Name  string `db:"name_sn" json:"name"`
	Role  string `db:"role_sn" json:"role"`
	Ip    string `db:"ip_sn" json:"ip"`
}

func (b *UserLocal) Validate(v *revel.Validation) {
	v.Check(b.Email,
		revel.ValidRequired(),
		revel.ValidMaxSize(30))
	v.Check(b.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
	v.Check(b.Role,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
	v.Check(b.Ip,
		revel.ValidRequired(),
		revel.ValidMaxSize(20))
}
