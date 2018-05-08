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

	fmt.Println("parseBodyCostume")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c BodyCostumeCtrl) Add() revel.Result {
	body_costume := c.parseBodyCostume()
	fmt.Println(body_costume)
	// Validate the model
	body_costume.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your body_costume.")
	} else {
		if err := c.Txn.Insert(&body_costume); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {

			return c.RenderJSON(body_costume)
		}
	}

}

func (c BodyCostumeCtrl) Get(id int64) revel.Result {
	body_costume := new(models.BodyCostume)
	err := c.Txn.SelectOne(body_costume,
		`SELECT * FROM body_costume WHERE body_costume_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. body_costume probably doesn't exist.")
	}
	return c.RenderJSON(body_costume)
}

func (c BodyCostumeCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	body_costume, err := c.Txn.Select(models.BodyCostume{},
		`SELECT * FROM body_costume WHERE body_costume_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(body_costume)
}

func (c BodyCostumeCtrl) Update(id int64) revel.Result {
	body_costume := c.parseBodyCostume()
	// Ensure the Id is set.
	body_costume.Id = id
	success, err := c.Txn.Update(&body_costume)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update body_costume.")
	}
	return c.RenderText("Updated %v", id)
}

func (c BodyCostumeCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.BodyCostume{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove body_costume")
	}
	return c.RenderText("Deleted %v", id)
}
