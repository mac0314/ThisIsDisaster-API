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
	stage := c.parseStage()
	fmt.Println(stage)
	// Validate the model
	stage.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your stage.")
	} else {
		if err := c.Txn.Insert(&stage); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(stage)
		}
	}

}

func (c StageCtrl) Get(id int64) revel.Result {
	stage := new(models.Stage)
	err := c.Txn.SelectOne(stage,
		`SELECT * FROM stage WHERE stage_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. stage probably doesn't exist.")
	}
	return c.RenderJSON(stage)
}

func (c StageCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	stage, err := c.Txn.Select(models.Stage{},
		`SELECT * FROM stage WHERE stage_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(stage)
}

func (c StageCtrl) Update(id int64) revel.Result {
	stage := c.parseStage()
	// Ensure the Id is set.
	stage.Id = id
	success, err := c.Txn.Update(&stage)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update stage.")
	}
	return c.RenderText("Updated %v", id)
}

func (c StageCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Stage{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove stage")
	}
	return c.RenderText("Deleted %v", id)
}
