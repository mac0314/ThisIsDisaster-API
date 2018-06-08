package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type HaveHeadCostumeCtrl struct {
	GorpController
}

func defineHaveHeadCostumeTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.HaveHeadCostume{}, "have_head_costume").SetKeys(true, "have_hc_id")

	t.ColMap("state_sn").SetMaxSize(20)
}

func (c HaveHeadCostumeCtrl) parseHaveHeadCostume() models.HaveHeadCostume {
	var jsonData models.HaveHeadCostume

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c HaveHeadCostumeCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	haveHeadCostume := c.parseHaveHeadCostume()

	haveHeadCostume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your have_head_costume."
	} else {
		if err := c.Txn.Insert(&haveHeadCostume); err != nil {
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

func (c HaveHeadCostumeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	haveHeadCostume := new(models.HaveHeadCostume)
	err := c.Txn.SelectOne(haveHeadCostume,
		`SELECT * FROM have_head_costume WHERE have_hc_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error have_head_costume probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = haveHeadCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAVE_HEAD_COSTUME

	return c.RenderJSON(response)
}

func (c HaveHeadCostumeCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	haveHeadCostume, err := c.Txn.Select(models.HaveHeadCostume{},
		`SELECT * FROM have_head_costume WHERE have_hc_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = haveHeadCostume
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_HAVE_HEAD_COSTUMES

	return c.RenderJSON(response)
}

func (c HaveHeadCostumeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	haveHeadCostume := c.parseHaveHeadCostume()
	// Ensure the Id is set.
	haveHeadCostume.Id = id
	success, err := c.Txn.Update(&haveHeadCostume)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update have_head_costume."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c HaveHeadCostumeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.HaveHeadCostume{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove have_head_costume"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
