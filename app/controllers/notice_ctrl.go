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
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	notice := c.parseNotice()
	fmt.Println(notice)
	// Validate the model
	notice.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your notice."
	} else {
		if err := c.Txn.Insert(&notice); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = notice
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c NoticeCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	notice := new(models.Notice)
	err := c.Txn.SelectOne(notice,
		`SELECT * FROM notice WHERE notice_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error notice probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = notice
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c NoticeCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	notice, err := c.Txn.Select(models.Notice{},
		`SELECT * FROM notice WHERE notice_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = notice
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c NoticeCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	notice := c.parseNotice()
	// Ensure the Id is set.
	notice.Id = id
	success, err := c.Txn.Update(&notice)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update notice."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = notice
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c NoticeCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Notice{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove notice"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
