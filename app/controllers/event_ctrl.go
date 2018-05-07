package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type EventCtrl struct {
	GorpController
}

func defineEventTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Event{}, "event").SetKeys(true, "event_id")
	// e.g. VARCHAR(25)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("content_ln").SetMaxSize(255)
}

func (c EventCtrl) parseEvent() models.Event {
	var jsonData models.Event

	fmt.Println("parseEvent")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c EventCtrl) Add() revel.Result {
	event := c.parseEvent()
	fmt.Println(event)
	// Validate the model
	event.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your event.")
	} else {
		if err := c.Txn.Insert(&event); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(event)
		}
	}

}

func (c EventCtrl) Get(id int64) revel.Result {
	event := new(models.Event)
	err := c.Txn.SelectOne(event,
		`SELECT * FROM event WHERE event_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. event probably doesn't exist.")
	}
	return c.RenderJSON(event)
}

func (c EventCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	event, err := c.Txn.Select(models.Event{},
		`SELECT * FROM event WHERE event_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(event)
}

func (c EventCtrl) Update(id int64) revel.Result {
	event := c.parseEvent()
	// Ensure the Id is set.
	event.Id = id
	success, err := c.Txn.Update(&event)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update event.")
	}
	return c.RenderText("Updated %v", id)
}

func (c EventCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Event{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove event")
	}
	return c.RenderText("Deleted %v", id)
}
