package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type AuthorizeCtrl struct {
	GorpController
}

func defineAuthorizeTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Authorize{}, "authorize").SetKeys(true, "auth_id")
	// e.g. VARCHAR(25)
	t.ColMap("email_mn").SetMaxSize(30)
	t.ColMap("platform_sn").SetMaxSize(20)
}

func (c AuthorizeCtrl) parseAuthorize() models.Authorize {
	var jsonData models.Authorize

	fmt.Println("parseAuthorize")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c AuthorizeCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	authorize := c.parseAuthorize()
	fmt.Println(authorize)
	// Validate the model
	authorize.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your authorize."
	} else {
		if err := c.Txn.Insert(&authorize); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = authorize
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c AuthorizeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	authorize := new(models.Authorize)
	err := c.Txn.SelectOne(authorize,
		`SELECT * FROM authorize WHERE auth_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error authorize probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = authorize
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c AuthorizeCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	authorize, err := c.Txn.Select(models.Authorize{},
		`SELECT * FROM authorize WHERE auth_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = authorize
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c AuthorizeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	authorize := c.parseAuthorize()
	// Ensure the Id is set.
	authorize.Id = id
	success, err := c.Txn.Update(&authorize)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update authorize."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = authorize
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c AuthorizeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Authorize{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove authorize"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
