package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type ItemCtrl struct {
	GorpController
}

func defineItemTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Item{}, "item").SetKeys(true, "item_id")
	// e.g. VARCHAR(25)
	t.ColMap("name_sn").SetMaxSize(20)
	t.ColMap("type_sn").SetMaxSize(20)
	t.ColMap("rank_sn").SetMaxSize(20)
	t.ColMap("effect_ln").SetMaxSize(255)
	t.ColMap("resource_mn").SetMaxSize(50)
}

func (c ItemCtrl) parseItem() models.Item {
	var jsonData models.Item

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c ItemCtrl) Add() revel.Result {
	item := c.parseItem()
	fmt.Println(item)
	// Validate the model
	item.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your item.")
	} else {
		if err := c.Txn.Insert(&item); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(item)
		}
	}

}

func (c ItemCtrl) Get(id int64) revel.Result {
	item := new(models.Item)
	err := c.Txn.SelectOne(item,
		`SELECT * FROM item WHERE item_id = ?`, id)
	if err != nil {
		return c.RenderText("Error.  item probably doesn't exist.")
	}
	return c.RenderJSON(item)
}

func (c ItemCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	item, err := c.Txn.Select(models.Item{},
		`SELECT * FROM item WHERE item_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(item)
}

func (c ItemCtrl) Update(id int64) revel.Result {
	item := c.parseItem()
	// Ensure the Id is set.
	item.Id = id
	success, err := c.Txn.Update(&item)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update item.")
	}
	return c.RenderText("Updated %v", id)
}

func (c ItemCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Item{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove item")
	}
	return c.RenderText("Deleted %v", id)
}
