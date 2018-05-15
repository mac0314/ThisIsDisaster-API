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
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	userSetting := c.parseUserSetting()
	fmt.Println(userSetting)
	// Validate the model
	userSetting.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your user setting."
	} else {
		if err := c.Txn.Insert(&userSetting); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = userSetting
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserSettingCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	userSetting := new(models.UserSetting)
	err := c.Txn.SelectOne(userSetting,
		`SELECT * FROM user_setting WHERE user_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error user setting probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = userSetting
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserSettingCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	userSetting, err := c.Txn.Select(models.UserSetting{},
		`SELECT * FROM user_setting WHERE user_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = userSetting
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserSettingCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	userSetting := c.parseUserSetting()
	// Ensure the Id is set.
	userSetting.Id = id
	success, err := c.Txn.Update(&userSetting)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update user setting."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = userSetting
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserSettingCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.User{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove user setting"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
