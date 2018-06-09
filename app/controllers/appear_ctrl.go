package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type AppearCtrl struct {
	GorpController
}

func defineAppearTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Appear{}, "appear").SetKeys(true, "appear_id")

	t.ColMap("state_sn").SetMaxSize(20)
}

func (c AppearCtrl) parseAppear() models.Appear {
	var jsonData models.Appear

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c AppearCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	appear := c.parseAppear()

	appear.Create = makeTimestamp()

	appear.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your appear."
	} else {
		if err := c.Txn.Insert(&appear); err != nil {
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

func (c AppearCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	appear := new(models.Appear)
	err := c.Txn.SelectOne(appear,
		`SELECT * FROM appear WHERE appear_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error appear probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = appear
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_APPEAR

	return c.RenderJSON(response)
}

func (c AppearCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	appear, err := c.Txn.Select(models.Appear{},
		`SELECT * FROM appear WHERE appear_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = appear
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_APPEARS

	return c.RenderJSON(response)
}

func (c AppearCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	appear := c.parseAppear()
	// Ensure the Id is set.
	appear.Id = id
	success, err := c.Txn.Update(&appear)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update appear."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c AppearCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Appear{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove appear"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
