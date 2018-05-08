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
	head_costume := c.parseHeadCostume()
	fmt.Println(head_costume)
	// Validate the model
	head_costume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your head_costume.")
	} else {
		if err := c.Txn.Insert(&head_costume); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {

			return c.RenderJSON(head_costume)
		}
	}

}

func (c HeadCostumeCtrl) Get(id int64) revel.Result {
	head_costume := new(models.HeadCostume)
	err := c.Txn.SelectOne(head_costume,
		`SELECT * FROM head_costume WHERE head_costume_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. head_costume probably doesn't exist.")
	}
	return c.RenderJSON(head_costume)
}

func (c HeadCostumeCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	head_costume, err := c.Txn.Select(models.HeadCostume{},
		`SELECT * FROM head_costume WHERE head_costume_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(head_costume)
}

func (c HeadCostumeCtrl) Update(id int64) revel.Result {
	head_costume := c.parseHeadCostume()
	// Ensure the Id is set.
	head_costume.Id = id
	success, err := c.Txn.Update(&head_costume)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update head_costume.")
	}
	return c.RenderText("Updated %v", id)
}

func (c HeadCostumeCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.HeadCostume{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove head_costume")
	}
	return c.RenderText("Deleted %v", id)
}
