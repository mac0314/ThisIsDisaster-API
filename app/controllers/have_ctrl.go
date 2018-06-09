package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type HaveCtrl struct {
	GorpController
}

func defineHaveTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Have{}, "have").SetKeys(true, "have_id")

	t.ColMap("state_sn").SetMaxSize(20)
}

func (c HaveCtrl) parseHave() models.Have {
	var jsonData models.Have

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c HaveCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	have := c.parseHave()

	have.Create = makeTimestamp()

	have.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your have."
	} else {
		if err := c.Txn.Insert(&have); err != nil {
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

func (c HaveCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	have := new(models.Have)
	err := c.Txn.SelectOne(have,
		`SELECT * FROM have WHERE have_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error have probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = have
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAVE

	return c.RenderJSON(response)
}

func (c HaveCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	have, err := c.Txn.Select(models.Have{},
		`SELECT * FROM have WHERE have_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = have
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAVES

	return c.RenderJSON(response)
}

func (c HaveCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	have := c.parseHave()
	// Ensure the Id is set.
	have.Id = id
	success, err := c.Txn.Update(&have)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update have."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c HaveCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Have{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove have"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
