package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type StageCtrl struct {
	GorpController
}

func defineStageTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Stage{}, "stage").SetKeys(true, "stage_id")
	// e.g. VARCHAR(25)
	t.ColMap("mode_sn").SetMaxSize(20)
}

func (c StageCtrl) parseStage() models.Stage {
	var jsonData models.Stage

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c StageCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	stage := c.parseStage()

	stage.Create = makeTimestamp()

	stage.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your stage."
	} else {
		if err := c.Txn.Insert(&stage); err != nil {
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

func (c StageCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	stage := new(models.Stage)
	err := c.Txn.SelectOne(stage,
		`SELECT * FROM stage WHERE stage_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + "Error. stage probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = stage
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_STAGE

	return c.RenderJSON(response)
}

func (c StageCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	stage, err := c.Txn.Select(models.Stage{},
		`SELECT * FROM stage WHERE stage_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		msg = msg + "Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = stage
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_STAGES

	return c.RenderJSON(response)
}

func (c StageCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	stage := c.parseStage()
	// Ensure the Id is set.
	stage.Id = id
	success, err := c.Txn.Update(&stage)
	if err != nil || success == 0 {
		msg = msg + " Unable to update stage."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c StageCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Stage{Id: id})
	if err != nil || success == 0 {
		msg = msg + " Failed to remove stage"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
