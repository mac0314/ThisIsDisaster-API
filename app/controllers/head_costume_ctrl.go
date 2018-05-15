package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type HeadCostumeCtrl struct {
	GorpController
}

func defineHeadCostumeTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.HeadCostume{}, "head_costume").SetKeys(true, "head_costume_id")
	// e.g. VARCHAR(25)
	t.ColMap("name_mn").SetMaxSize(30)
	t.ColMap("resource_mn").SetMaxSize(50)
}

func (c HeadCostumeCtrl) parseHeadCostume() models.HeadCostume {
	var jsonData models.HeadCostume

	fmt.Println("parseHeadCostume")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c HeadCostumeCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	headCostume := c.parseHeadCostume()
	fmt.Println(headCostume)
	// Validate the model
	headCostume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your head costume."
	} else {
		if err := c.Txn.Insert(&headCostume); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = headCostume
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c HeadCostumeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	headCostume := new(models.HeadCostume)
	err := c.Txn.SelectOne(headCostume,
		`SELECT * FROM head_costume WHERE head_costume_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error head costume probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = headCostume
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c HeadCostumeCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	headCostume, err := c.Txn.Select(models.HeadCostume{},
		`SELECT * FROM head_costume WHERE head_costume_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = headCostume
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c HeadCostumeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	headCostume := c.parseHeadCostume()
	// Ensure the Id is set.
	headCostume.Id = id
	success, err := c.Txn.Update(&headCostume)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update head costume."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = headCostume
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c HeadCostumeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.HeadCostume{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove head costume"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
