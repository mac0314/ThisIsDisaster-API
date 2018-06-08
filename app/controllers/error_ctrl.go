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

	c.Params.BindJSON(&jsonData)

	jsonData.Create = makeTimestamp()

	return jsonData
}

func (c ErrorCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	error := c.parseError()
	// Validate the model
	error.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your error."
	} else {
		if err := c.Txn.Insert(&error); err != nil {
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

func (c ErrorCtrl) selectErrorById(id int64) (bool, *models.Error) {
	var err bool
	error := new(models.Error)
	_err := c.Txn.SelectOne(error,
		`SELECT * FROM error WHERE error_id = ?`, id)
	if _err != nil {
		err = true
	} else {
		err = false
	}
	return err, error
}

func (c ErrorCtrl) selectErrorByUserId(id int64) (bool, *models.Error) {
	var err bool
	var error *models.Error

	_err := c.Txn.SelectOne(&error,
		`SELECT * FROM error WHERE user_id = ? ORDER BY error_id DESC LIMIT 1`, id)
	if _err != nil {
		err = true
	} else {
		err = false
	}
	return err, error
}

func (c ErrorCtrl) selectErrorByEmail(email string) (bool, *models.Error) {
	var err bool
	var error *models.Error

	_err := c.Txn.SelectOne(&error,
		`SELECT * FROM error WHERE error.user_id = (SELECT user_id FROM user WHERE email_mn = ?) ORDER BY error_id DESC LIMIT 1`, email)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, error
}

func (c ErrorCtrl) selectErrors(lastId int64, limit uint64) (bool, []models.Error) {
	var err bool
	var errors []models.Error

	_, _err := c.Txn.Select(&errors,
		`SELECT * FROM error WHERE error_id > ? LIMIT ?`, lastId, limit)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, errors
}

func (c ErrorCtrl) selectErrorsByUserId(userId int64) (bool, []models.Error) {
	var err bool
	var errors []models.Error

	_, _err := c.Txn.Select(&errors,
		`SELECT * FROM error WHERE user_id = ?`, userId)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, errors
}

func (c ErrorCtrl) Get() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	var _err bool
	var error *models.Error

	id := parseIntOrDefault(c.Params.Get("id"), -1)
	if id > -1 {
		_err, error = c.selectErrorById(id)
	} else {
		userId := parseIntOrDefault(c.Params.Get("uid"), -1)
		if userId > -1 {
			_err, error = c.selectErrorByUserId(userId)
		} else {
			fmt.Println("check")
			email := c.Params.Get("email")
			_err, error = c.selectErrorByEmail(email)
		}

	}

	if _err {
		msg = msg + " Error error probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = error
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_ERROR

	return c.RenderJSON(response)
}

func (c ErrorCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	var err bool
	var errors []models.Error

	userId := parseIntOrDefault(c.Params.Get("uid"), -1)
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))

	if userId > -1 {
		err, errors = c.selectErrorsByUserId(userId)
	} else {
		err, errors = c.selectErrors(lastId, limit)
	}

	if err {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = errors
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_ERRORS

	return c.RenderJSON(response)
}

func (c ErrorCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	error := c.parseError()
	// Ensure the Id is set.
	error.Id = id
	success, err := c.Txn.Update(&error)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update error."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c ErrorCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Error{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove error"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
