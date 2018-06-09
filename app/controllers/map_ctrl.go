package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type MapCtrl struct {
	GorpController
}

func defineMapTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Map{}, "map").SetKeys(true, "map_id")

	t.ColMap("formation_ln").SetMaxSize(255)
}

func (c MapCtrl) parseMap() models.Map {
	var jsonData models.Map

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c MapCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	_map := c.parseMap()

	_map.Create = makeTimestamp()

	_map.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your map."
	} else {
		if err := c.Txn.Insert(&_map); err != nil {
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

func (c MapCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	_map := new(models.Map)
	err := c.Txn.SelectOne(_map,
		`SELECT * FROM map WHERE map_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + "Error. _map probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = _map
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MAP

	return c.RenderJSON(response)
}

func (c MapCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	_map, err := c.Txn.Select(models.Map{},
		`SELECT * FROM map WHERE map_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		msg = msg + "Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = _map
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MAPS

	return c.RenderJSON(response)
}

func (c MapCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	_map := c.parseMap()
	// Ensure the Id is set.
	_map.Id = id
	success, err := c.Txn.Update(&_map)
	if err != nil || success == 0 {
		msg = msg + " Unable to update map."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c MapCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Map{Id: id})
	if err != nil || success == 0 {
		msg = msg + " Failed to remove map"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
