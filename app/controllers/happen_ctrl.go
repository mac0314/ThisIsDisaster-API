package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type HappenCtrl struct {
	GorpController
}

func defineHappenTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	dbm.AddTableWithName(models.Happen{}, "happen").SetKeys(true, "happen_id")
}

func (c HappenCtrl) parseHappen() models.Happen {
	var jsonData models.Happen

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c HappenCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	happen := c.parseHappen()

	happen.Create = makeTimestamp()

	happen.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your happen."
	} else {
		if err := c.Txn.Insert(&happen); err != nil {
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

func (c HappenCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	happen := new(models.Happen)
	err := c.Txn.SelectOne(happen,
		`SELECT * FROM happen WHERE happen_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error happen probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = happen
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAPPEN

	return c.RenderJSON(response)
}

func (c HappenCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	happen, err := c.Txn.Select(models.Happen{},
		`SELECT * FROM happen WHERE happen_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = happen
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAPPENS

	return c.RenderJSON(response)
}

func (c HappenCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	happen := c.parseHappen()
	// Ensure the Id is set.
	happen.Id = id
	success, err := c.Txn.Update(&happen)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update happen."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c HappenCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Happen{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove happen"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
