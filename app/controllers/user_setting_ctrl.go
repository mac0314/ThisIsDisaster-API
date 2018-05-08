package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type UserSettingCtrl struct {
	GorpController
}

func defineUserSettingTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	dbm.AddTableWithName(models.UserSetting{}, "user_setting").SetKeys(true, "user_id")
}

func (c UserSettingCtrl) parseUserSetting() models.UserSetting {
	var jsonData models.UserSetting

	fmt.Println("parseUser")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c UserSettingCtrl) Add() revel.Result {
	userSetting := c.parseUserSetting()
	fmt.Println(userSetting)
	// Validate the model
	userSetting.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your user setting.")
	} else {
		if err := c.Txn.Insert(&userSetting); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(userSetting)
		}
	}

}

func (c UserSettingCtrl) Get(id int64) revel.Result {
	userSetting := new(models.UserSetting)
	err := c.Txn.SelectOne(userSetting,
		`SELECT * FROM user_setting WHERE user_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. user setting probably doesn't exist.")
	}
	return c.RenderJSON(userSetting)
}

func (c UserSettingCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	userSetting, err := c.Txn.Select(models.UserSetting{},
		`SELECT * FROM user_setting WHERE user_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(userSetting)
}

func (c UserSettingCtrl) Update(id int64) revel.Result {
	userSetting := c.parseUserSetting()
	// Ensure the Id is set.
	userSetting.Id = id
	success, err := c.Txn.Update(&userSetting)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update user setting.")
	}
	return c.RenderText("Updated %v", id)
}

func (c UserSettingCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.User{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove user setting")
	}
	return c.RenderText("Deleted %v", id)
}
