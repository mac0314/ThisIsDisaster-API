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

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c EventCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	response := make(map[string]interface{})

	event := c.parseEvent()

	event.Create = makeTimestamp()

	event.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your event."
	} else {
		if err := c.Txn.Insert(&event); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
		}
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c EventCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	response := make(map[string]interface{})

	event := new(models.Event)
	err := c.Txn.SelectOne(event,
		`SELECT * FROM event WHERE event_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error event probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		response["result_data"] = event
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_EVENT

	return c.RenderJSON(response)
}

func (c EventCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	event, err := c.Txn.Select(models.Event{},
		`SELECT * FROM event WHERE event_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		response["result_data"] = event
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_EVENTS

	return c.RenderJSON(response)
}

func (c EventCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	response := make(map[string]interface{})

	event := c.parseEvent()
	// Ensure the Id is set.
	event.Id = id
	success, err := c.Txn.Update(&event)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update event."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c EventCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Event{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove event"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
