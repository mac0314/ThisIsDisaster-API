package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type ErrorCtrl struct {
	GorpController
}

func defineErrorTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Error{}, "error").SetKeys(true, "error_id")
	// e.g. VARCHAR(25)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("log_ln").SetMaxSize(255)
}

func (c ErrorCtrl) parseError() models.Error {
	var jsonData models.Error

	fmt.Println("parseError")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c ErrorCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	error := c.parseError()
	fmt.Println(error)
	// Validate the model
	error.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your error."
	} else {
		if err := c.Txn.Insert(&error); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = error
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c ErrorCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	error := new(models.Error)
	err := c.Txn.SelectOne(error,
		`SELECT * FROM error WHERE error_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error error probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = error
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c ErrorCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	error, err := c.Txn.Select(models.Error{},
		`SELECT * FROM error WHERE error_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = error
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c ErrorCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	error := c.parseError()
	// Ensure the Id is set.
	error.Id = id
	success, err := c.Txn.Update(&error)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update error."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = error
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c ErrorCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Error{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove error"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
