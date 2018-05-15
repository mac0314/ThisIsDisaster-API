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

	fmt.Println("parseStage")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c StageCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	stage := c.parseStage()
	fmt.Println(stage)
	// Validate the model
	stage.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your stage."
	} else {
		if err := c.Txn.Insert(&stage); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = stage
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c StageCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	stage := new(models.Stage)
	err := c.Txn.SelectOne(stage,
		`SELECT * FROM stage WHERE stage_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + "Error. stage probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = stage
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c StageCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	stage, err := c.Txn.Select(models.Stage{},
		`SELECT * FROM stage WHERE stage_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		msg = msg + "Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = stage
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c StageCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	stage := c.parseStage()
	// Ensure the Id is set.
	stage.Id = id
	success, err := c.Txn.Update(&stage)
	if err != nil || success == 0 {
		msg = msg + " Unable to update stage."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c StageCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Stage{Id: id})
	if err != nil || success == 0 {
		msg = msg + " Failed to remove stage"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
