package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type BodyCostumeCtrl struct {
	GorpController
}

func defineBodyCostumeTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.BodyCostume{}, "body_costume").SetKeys(true, "body_costume_id")
	// e.g. VARCHAR(25)
	t.ColMap("name_mn").SetMaxSize(30)
	t.ColMap("resource_mn").SetMaxSize(50)
}

func (c BodyCostumeCtrl) parseBodyCostume() models.BodyCostume {
	var jsonData models.BodyCostume

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c BodyCostumeCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	bodyCostume := c.parseBodyCostume()

	bodyCostume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your body costume."
	} else {
		if err := c.Txn.Insert(&bodyCostume); err != nil {
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

func (c BodyCostumeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	bodyCostume := new(models.BodyCostume)
	err := c.Txn.SelectOne(bodyCostume,
		`SELECT * FROM body_costume WHERE body_costume_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error body costume probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = bodyCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_BODY_COSTUME

	return c.RenderJSON(response)
}

func (c BodyCostumeCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	bodyCostume, err := c.Txn.Select(models.BodyCostume{},
		`SELECT * FROM body_costume WHERE body_costume_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = bodyCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_BODY_COSTUMES

	return c.RenderJSON(response)
}

func (c BodyCostumeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	bodyCostume := c.parseBodyCostume()
	// Ensure the Id is set.
	bodyCostume.Id = id
	success, err := c.Txn.Update(&bodyCostume)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update body costume."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c BodyCostumeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.BodyCostume{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove body costume"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
