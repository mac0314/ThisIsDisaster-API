package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type NoticeCtrl struct {
	GorpController
}

func defineNoticeTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Notice{}, "notice").SetKeys(true, "notice_id")
	// e.g. VARCHAR(25)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("content_ln").SetMaxSize(255)
}

func (c NoticeCtrl) parseNotice() models.Notice {
	var jsonData models.Notice

	fmt.Println("parseNotice")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c NoticeCtrl) Add() revel.Result {
	notice := c.parseNotice()
	fmt.Println(notice)
	// Validate the model
	notice.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your notice.")
	} else {
		if err := c.Txn.Insert(&notice); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(notice)
		}
	}

}

func (c NoticeCtrl) Get(id int64) revel.Result {
	notice := new(models.Notice)
	err := c.Txn.SelectOne(notice,
		`SELECT * FROM notice WHERE notice_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. notice probably doesn't exist.")
	}
	return c.RenderJSON(notice)
}

func (c NoticeCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	notice, err := c.Txn.Select(models.Notice{},
		`SELECT * FROM notice WHERE notice_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(notice)
}

func (c NoticeCtrl) Update(id int64) revel.Result {
	notice := c.parseNotice()
	// Ensure the Id is set.
	notice.Id = id
	success, err := c.Txn.Update(&notice)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update notice.")
	}
	return c.RenderText("Updated %v", id)
}

func (c NoticeCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Notice{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove notice")
	}
	return c.RenderText("Deleted %v", id)
}
