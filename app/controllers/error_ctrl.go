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
	error := c.parseError()
	fmt.Println(error)
	// Validate the model
	error.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your error.")
	} else {
		if err := c.Txn.Insert(&error); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(error)
		}
	}

}

func (c ErrorCtrl) Get(id int64) revel.Result {
	error := new(models.Error)
	err := c.Txn.SelectOne(error,
		`SELECT * FROM error WHERE error_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. error probably doesn't exist.")
	}
	return c.RenderJSON(error)
}

func (c ErrorCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	error, err := c.Txn.Select(models.Error{},
		`SELECT * FROM error WHERE error_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(error)
}

func (c ErrorCtrl) Update(id int64) revel.Result {
	error := c.parseError()
	// Ensure the Id is set.
	error.Id = id
	success, err := c.Txn.Update(&error)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update error.")
	}
	return c.RenderText("Updated %v", id)
}

func (c ErrorCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Error{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove error")
	}
	return c.RenderText("Deleted %v", id)
}
