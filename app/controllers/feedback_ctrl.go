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
	feedback := c.parseFeedback()
	fmt.Println(feedback)
	// Validate the model
	feedback.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your feedback.")
	} else {
		if err := c.Txn.Insert(&feedback); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(feedback)
		}
	}

}

func (c FeedbackCtrl) Get(id int64) revel.Result {
	feedback := new(models.Feedback)
	err := c.Txn.SelectOne(feedback,
		`SELECT * FROM feedback WHERE feedback_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. feedback probably doesn't exist.")
	}
	return c.RenderJSON(feedback)
}

func (c FeedbackCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	feedback, err := c.Txn.Select(models.Feedback{},
		`SELECT * FROM feedback WHERE feedback_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(feedback)
}

func (c FeedbackCtrl) Update(id int64) revel.Result {
	feedback := c.parseFeedback()
	// Ensure the Id is set.
	feedback.Id = id
	success, err := c.Txn.Update(&feedback)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update feedback.")
	}
	return c.RenderText("Updated %v", id)
}

func (c FeedbackCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Feedback{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove feedback")
	}
	return c.RenderText("Deleted %v", id)
}
