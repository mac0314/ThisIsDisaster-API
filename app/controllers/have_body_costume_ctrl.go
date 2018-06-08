package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type HaveBodyCostumeCtrl struct {
	GorpController
}

func defineHaveBodyCostumeTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.HaveBodyCostume{}, "have_body_costume").SetKeys(true, "have_bc_id")

	t.ColMap("state_sn").SetMaxSize(20)
}

func (c HaveBodyCostumeCtrl) parseHaveBodyCostume() models.HaveBodyCostume {
	var jsonData models.HaveBodyCostume

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c HaveBodyCostumeCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	haveBodyCostume := c.parseHaveBodyCostume()

	haveBodyCostume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your have_body_costume."
	} else {
		if err := c.Txn.Insert(&haveBodyCostume); err != nil {
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

func (c HaveBodyCostumeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	haveBodyCostume := new(models.HaveBodyCostume)
	err := c.Txn.SelectOne(haveBodyCostume,
		`SELECT * FROM have_body_costume WHERE have_bc_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error have_body_costume probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = haveBodyCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAVE_BODY_COSTUME

	return c.RenderJSON(response)
}

func (c HaveBodyCostumeCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	haveBodyCostume, err := c.Txn.Select(models.HaveBodyCostume{},
		`SELECT * FROM have_body_costume WHERE have_bc_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = haveBodyCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAVE_BODY_COSTUMES

	return c.RenderJSON(response)
}

func (c HaveBodyCostumeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	haveBodyCostume := c.parseHaveBodyCostume()
	// Ensure the Id is set.
	haveBodyCostume.Id = id
	success, err := c.Txn.Update(&haveBodyCostume)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update have_body_costume."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c HaveBodyCostumeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.HaveBodyCostume{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove have_body_costume"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
