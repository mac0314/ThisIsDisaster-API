package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type MovementCtrl struct {
	GorpController
}

func defineMovementTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	dbm.AddTableWithName(models.Movement{}, "movement").SetKeys(true, "movement_id")
}

func (c MovementCtrl) parseMovement() models.Movement {
	var jsonData models.Movement

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c MovementCtrl) Add(data models.Movement) (bool, string) {
	var err bool
	var msg string

	data.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your movement."
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

func (c MovementCtrl) Post() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	movement := c.parseMovement()

	movement.Create = makeTimestamp()

	err, msg := c.Add(movement)

	fmt.Println(err)

	if !err {
		code = RESULT_CODE_SUCCESS
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)

}

func (c MovementCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	movement := new(models.Movement)
	err := c.Txn.SelectOne(movement,
		`SELECT * FROM movement WHERE movement_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error movement probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = movement
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MOVEMENT

	return c.RenderJSON(response)
}

func (c MovementCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	movement, err := c.Txn.Select(models.Movement{},
		`SELECT * FROM movement WHERE movement_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = movement
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MOVEMENTS

	return c.RenderJSON(response)
}

func (c MovementCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	movement := c.parseMovement()
	// Ensure the Id is set.
	movement.Id = id
	success, err := c.Txn.Update(&movement)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update movement."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		response["result_data"] = movement
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c MovementCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Movement{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to removement movement"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
