package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type UserCtrl struct {
	GorpController
}

func defineUserTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.User{}, "user").SetKeys(true, "user_id")
	// e.g. VARCHAR(25)
	t.ColMap("email_mn").SetMaxSize(30)
	t.ColMap("nickname_mn").SetMaxSize(30)
}

func (c UserCtrl) parseUser() models.User {
	var jsonData models.User

	fmt.Println("parseUser")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c UserCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	user := c.parseUser()
	fmt.Println(user)
	// Validate the model
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your user."
	} else {
		if err := c.Txn.Insert(&user); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = user
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	fmt.Println(data)

	return c.RenderJSON(data)
}

func (c UserCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	user := new(models.User)
	err := c.Txn.SelectOne(user,
		`SELECT * FROM user WHERE user_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error user probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = user
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	user, err := c.Txn.Select(models.User{},
		`SELECT * FROM user WHERE user_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = user
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	user := c.parseUser()
	// Ensure the Id is set.
	user.Id = id
	success, err := c.Txn.Update(&user)
	if err != nil || success == 0 {
		msg = msg + " Unable to update user."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = user
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c UserCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.User{Id: id})
	if err != nil || success == 0 {
		msg = msg + " Failed to remove user"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
