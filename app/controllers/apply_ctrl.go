package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type ApplyCtrl struct {
	GorpController
}

func defineApplyTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Apply{}, "apply").SetKeys(true, "apply_id")
	// e.g. VARCHAR(25)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("log_ln").SetMaxSize(255)
}

func (c ApplyCtrl) parseApply() models.Apply {
	var jsonData models.Apply

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c ApplyCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	apply := c.parseApply()

	apply.Create = makeTimestamp()

	apply.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your apply."
	} else {
		if err := c.Txn.Insert(&apply); err != nil {
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

func (c ApplyCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	apply := new(models.Apply)
	err := c.Txn.SelectOne(apply,
		`SELECT * FROM apply WHERE apply_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error apply probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = apply
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HANDLE

	return c.RenderJSON(response)
}

func (c ApplyCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	apply, err := c.Txn.Select(models.Apply{},
		`SELECT * FROM apply WHERE apply_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = apply
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HANDLES

	return c.RenderJSON(response)
}

func (c ApplyCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	apply := c.parseApply()
	// Ensure the Id is set.
	apply.Id = id
	success, err := c.Txn.Update(&apply)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update apply."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c ApplyCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Apply{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove apply"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
