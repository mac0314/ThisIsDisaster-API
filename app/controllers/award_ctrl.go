package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type AwardCtrl struct {
	GorpController
}

func defineAwardTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	dbm.AddTableWithName(models.Award{}, "award").SetKeys(true, "award_id")
}

func (c AwardCtrl) parseAward() models.Award {
	var jsonData models.Award

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c AwardCtrl) Add(data models.Award) (bool, string) {
	var err bool
	var msg string

	data.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your award."
	} else {
		if _err := c.Txn.Insert(&data); _err != nil {
			fmt.Println(_err)

			msg = msg + " Error inserting record into database!"
		} else {
			msg = "Success."
		}
	}

	return err, msg
}

func (c AwardCtrl) Post() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	award := c.parseAward()

	err, msg := c.Add(award)

	fmt.Println(err)

	if !err {
		code = RESULT_CODE_SUCCESS
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)

}

func (c AwardCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	award := new(models.Award)
	err := c.Txn.SelectOne(award,
		`SELECT * FROM award WHERE award_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error award probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = award
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_AWARD

	return c.RenderJSON(response)
}

func (c AwardCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	award, err := c.Txn.Select(models.Award{},
		`SELECT * FROM award WHERE award_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = award
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_AWARDS

	return c.RenderJSON(response)
}

func (c AwardCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	award := c.parseAward()
	// Ensure the Id is set.
	award.Id = id
	success, err := c.Txn.Update(&award)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update award."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		response["result_data"] = award
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c AwardCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Award{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove award"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
