package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type MediaCtrl struct {
	GorpController
}

func defineMediaTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Media{}, "media").SetKeys(true, "media_id")
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("content_ln").SetMaxSize(255)
	t.ColMap("thumbnail_ln").SetMaxSize(255)
	t.ColMap("image_ln").SetMaxSize(255)
}

func (c MediaCtrl) parseMedia() models.Media {
	var jsonData models.Media

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c MediaCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	media := c.parseMedia()

	media.Create = makeTimestamp()

	media.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your media."
	} else {
		if err := c.Txn.Insert(&media); err != nil {
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

func (c MediaCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	media := new(models.Media)
	err := c.Txn.SelectOne(media,
		`SELECT * FROM media WHERE media_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error media probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = media
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MEDIUM

	return c.RenderJSON(response)
}

func (c MediaCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	media, err := c.Txn.Select(models.Media{},
		`SELECT * FROM media WHERE media_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = media
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MEDIA

	return c.RenderJSON(response)
}

func (c MediaCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	media := c.parseMedia()
	// Ensure the Id is set.
	media.Id = id
	success, err := c.Txn.Update(&media)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update media."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c MediaCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Media{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove media"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
