package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type AchievementCtrl struct {
	GorpController
}

//RESULT_TYPE := "achievement"

func defineAchievementTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Achievement{}, "achievement").SetKeys(true, "achievement_id")
	// e.g. VARCHAR(25)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("content_ln").SetMaxSize(255)
}

func (c AchievementCtrl) parseAchievement() models.Achievement {
	var jsonData models.Achievement

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c AchievementCtrl) Add(data models.Achievement) (bool, string) {
	var err bool
	var msg string

	data.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your achievement."
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

func (c AchievementCtrl) Post() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	achievement := c.parseAchievement()

	err, msg := c.Add(achievement)

	fmt.Println(err)

	if !err {
		code = RESULT_CODE_SUCCESS
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)

}

func (c AchievementCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	achievement := new(models.Achievement)
	err := c.Txn.SelectOne(achievement,
		`SELECT * FROM achievement WHERE achievement_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error achievement probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = achievement
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_ACHIEVEMENT

	return c.RenderJSON(response)
}

func (c AchievementCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	achievement, err := c.Txn.Select(models.Achievement{},
		`SELECT * FROM achievement WHERE achievement_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = achievement
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_ACHIEVEMENTS

	return c.RenderJSON(response)
}

func (c AchievementCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	achievement := c.parseAchievement()
	// Ensure the Id is set.
	achievement.Id = id
	success, err := c.Txn.Update(&achievement)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update achievement."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		response["result_data"] = achievement
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c AchievementCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Achievement{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove achievement"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
