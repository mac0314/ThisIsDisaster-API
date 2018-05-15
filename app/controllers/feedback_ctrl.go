package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type FeedbackCtrl struct {
	GorpController
}

func defineFeedbackTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Feedback{}, "feedback").SetKeys(true, "feedback_id")
	// e.g. VARCHAR(25)
	t.ColMap("email_mn").SetMaxSize(30)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("content_ln").SetMaxSize(255)
}

func (c FeedbackCtrl) parseFeedback() models.Feedback {
	var jsonData models.Feedback

	fmt.Println("parseFeedback")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c FeedbackCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	feedback := c.parseFeedback()
	fmt.Println(feedback)
	// Validate the model
	feedback.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your feedback."
	} else {
		if err := c.Txn.Insert(&feedback); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = feedback
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c FeedbackCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	feedback := new(models.Feedback)
	err := c.Txn.SelectOne(feedback,
		`SELECT * FROM feedback WHERE feedback_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error feedback probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = feedback
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c FeedbackCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	feedback, err := c.Txn.Select(models.Feedback{},
		`SELECT * FROM feedback WHERE feedback_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = feedback
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c FeedbackCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	feedback := c.parseFeedback()
	// Ensure the Id is set.
	feedback.Id = id
	success, err := c.Txn.Update(&feedback)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update feedback."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = feedback
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c FeedbackCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Feedback{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove feedback"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
