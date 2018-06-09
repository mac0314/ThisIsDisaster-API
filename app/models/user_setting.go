package models

import (
	"github.com/revel/revel"
)

type UserSetting struct {
	Id        int64 `db:"user_id, primarykey" json:"id"`
	PushCheck bool  `db:"push_check" json:"pushCheck"`
	Sound     int64 `db:"sound_n" json:"sound"`
	Update    int64 `db:"update_n" json:"update"`
}

func (b *UserSetting) Validate(v *revel.Validation) {

}
