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

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c UserSettingCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	userSetting := c.parseUserSetting()

	userSetting.Update = makeTimestamp()

	userSetting.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your user setting."
	} else {
		if err := c.Txn.Insert(&userSetting); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = RESULT_CODE_SUCCESS
			msg = "Success."
		}
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c UserSettingCtrl) Get() revel.Result {
	id := parseIntOrDefault(c.Params.Get("uid"), 0)
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	userSetting := new(models.UserSetting)
	err := c.Txn.SelectOne(userSetting,
		`SELECT * FROM user_setting WHERE user_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error user setting probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = userSetting
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_USER_SETTING

	return c.RenderJSON(response)
}

func (c UserSettingCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	userSetting, err := c.Txn.Select(models.UserSetting{},
		`SELECT * FROM user_setting WHERE user_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = userSetting
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_USER_SETTINGS

	return c.RenderJSON(response)
}

func (c UserSettingCtrl) Update() revel.Result {
	id := parseIntOrDefault(c.Params.Get("uid"), 0)
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	userSetting := c.parseUserSetting()
	// Ensure the Id is set.
	userSetting.Id = id
	success, err := c.Txn.Update(&userSetting)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update user setting."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c UserSettingCtrl) Delete() revel.Result {
	id := parseIntOrDefault(c.Params.Get("uid"), 0)
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.User{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove user setting"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
