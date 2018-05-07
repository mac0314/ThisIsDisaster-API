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
	user := c.parseUser()
	fmt.Println(user)
	// Validate the model
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your user.")
	} else {
		if err := c.Txn.Insert(&user); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(user)
		}
	}

}

func (c UserCtrl) Get(id int64) revel.Result {
	user := new(models.User)
	err := c.Txn.SelectOne(user,
		`SELECT * FROM user WHERE user_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. user probably doesn't exist.")
	}
	return c.RenderJSON(user)
}

func (c UserCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	user, err := c.Txn.Select(models.User{},
		`SELECT * FROM user WHERE user_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(user)
}

func (c UserCtrl) Update(id int64) revel.Result {
	user := c.parseUser()
	// Ensure the Id is set.
	user.Id = id
	success, err := c.Txn.Update(&user)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update user.")
	}
	return c.RenderText("Updated %v", id)
}

func (c UserCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.User{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove user")
	}
	return c.RenderText("Deleted %v", id)
}
