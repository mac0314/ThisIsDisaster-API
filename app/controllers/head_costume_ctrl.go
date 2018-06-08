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

	t.ColMap("name_mn").SetMaxSize(30)
	t.ColMap("resource_mn").SetMaxSize(50)
}

func (c HeadCostumeCtrl) parseHeadCostume() models.HeadCostume {
	var jsonData models.HeadCostume

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c HeadCostumeCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	headCostume := c.parseHeadCostume()

	headCostume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your head costume."
	} else {
		if err := c.Txn.Insert(&headCostume); err != nil {
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

func (c HeadCostumeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	headCostume := new(models.HeadCostume)
	err := c.Txn.SelectOne(headCostume,
		`SELECT * FROM head_costume WHERE head_costume_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error head costume probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = headCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HEAD_COSTUME

	return c.RenderJSON(response)
}

func (c HeadCostumeCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	headCostume, err := c.Txn.Select(models.HeadCostume{},
		`SELECT * FROM head_costume WHERE head_costume_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = headCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HEAD_COSTUMES

	return c.RenderJSON(response)
}

func (c HeadCostumeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	headCostume := c.parseHeadCostume()
	// Ensure the Id is set.
	headCostume.Id = id
	success, err := c.Txn.Update(&headCostume)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update head costume."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c HeadCostumeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.HeadCostume{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove head costume"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
