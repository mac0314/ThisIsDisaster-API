package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type GameLogCtrl struct {
	GorpController
}

func defineGameLogTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.GameLog{}, "game_log").SetKeys(true, "game_log_id")

	t.ColMap("room_sn").SetMaxSize(20)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("log_ln").SetMaxSize(255)
}

func (c GameLogCtrl) parseGameLog() models.GameLog {
	var jsonData models.GameLog

	c.Params.BindJSON(&jsonData)

	jsonData.Create = makeTimestamp()

	return jsonData
}

func (c GameLogCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	gameLog := c.parseGameLog()

	gameLog.Create = makeTimestamp()

	gameLog.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your game_log."
	} else {
		if err := c.Txn.Insert(&gameLog); err != nil {
			msg = msg + " GameLog inserting record into database!"
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

func (c GameLogCtrl) selectGameLogById(id int64) (bool, *models.GameLog) {
	var err bool
	gameLog := new(models.GameLog)
	_err := c.Txn.SelectOne(gameLog,
		`SELECT * FROM game_log WHERE game_log_id = ?`, id)
	if _err != nil {
		err = true
	} else {
		err = false
	}
	return err, gameLog
}

func (c GameLogCtrl) selectGameLogByUserId(id int64) (bool, *models.GameLog) {
	var err bool
	var gameLog *models.GameLog

	_err := c.Txn.SelectOne(&gameLog,
		`SELECT * FROM game_log WHERE user_id = ? ORDER BY game_log_id DESC LIMIT 1`, id)
	if _err != nil {
		err = true
	} else {
		err = false
	}
	return err, gameLog
}

func (c GameLogCtrl) selectGameLogByEmail(email string) (bool, *models.GameLog) {
	var err bool
	var gameLog *models.GameLog

	_err := c.Txn.SelectOne(&gameLog,
		`SELECT * FROM game_log WHERE game_log.user_id = (SELECT user_id FROM user WHERE email_mn = ?) ORDER BY game_log_id DESC LIMIT 1`, email)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, gameLog
}

func (c GameLogCtrl) selectGameRoomLogs(lastId int64, limit uint64, room string) (bool, []models.GameLog) {
	var err bool
	var gameLogs []models.GameLog

	_, _err := c.Txn.Select(&gameLogs,
		`SELECT * FROM game_log WHERE room_sn = ? AND game_log_id > ? LIMIT ?`, room, lastId, limit)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, gameLogs
}

func (c GameLogCtrl) selectGameRoomUserLogs(lastId int64, limit uint64, userId int64, room string) (bool, []models.GameLog) {
	var err bool
	var gameLogs []models.GameLog

	_, _err := c.Txn.Select(&gameLogs,
		`SELECT * FROM game_log WHERE user_id = ? AND room_sn = ? AND game_log_id > ? LIMIT ?`, userId, room, lastId, limit)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, gameLogs
}

func (c GameLogCtrl) selectGameLogs(lastId int64, limit uint64) (bool, []models.GameLog) {
	var err bool
	var gameLogs []models.GameLog

	_, _err := c.Txn.Select(&gameLogs,
		`SELECT * FROM game_log WHERE game_log_id > ? LIMIT ?`, lastId, limit)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, gameLogs
}

func (c GameLogCtrl) selectGameLogsByUserId(userId int64) (bool, []models.GameLog) {
	var err bool
	var gameLogs []models.GameLog

	_, _err := c.Txn.Select(&gameLogs,
		`SELECT * FROM game_log WHERE user_id = ?`, userId)
	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, gameLogs
}

func (c GameLogCtrl) Get() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	var _err bool
	var gameLog *models.GameLog

	id := parseIntOrDefault(c.Params.Get("id"), -1)
	if id > -1 {
		_err, gameLog = c.selectGameLogById(id)
	} else {
		userId := parseIntOrDefault(c.Params.Get("uid"), -1)
		if userId > -1 {
			_err, gameLog = c.selectGameLogByUserId(userId)
		} else {
			fmt.Println("check")
			email := c.Params.Get("email")
			_err, gameLog = c.selectGameLogByEmail(email)
		}

	}

	if _err {
		msg = msg + " GameLog game_log probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = gameLog
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_GAME_LOG

	return c.RenderJSON(response)
}

func (c GameLogCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	var err bool
	var gameLogs []models.GameLog

	email := c.Params.Get("email")
	room := c.Params.Get("room")

	userId := parseIntOrDefault(c.Params.Get("uid"), -1)
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))

	if email != "" {
		var user models.User

		query := "SELECT * FROM user WHERE email_mn='" + email + "'"

		_err := c.Txn.SelectOne(&user, query)
		if _err != nil {
			fmt.Println(_err.Error())

			if room != "" {
				err, gameLogs = c.selectGameRoomUserLogs(lastId, limit, userId, room)
			} else {
				err, gameLogs = c.selectGameLogsByUserId(userId)
			}
		} else {
			userId = user.Id

			if room != "" {
				err, gameLogs = c.selectGameRoomLogs(lastId, limit, room)
			} else {
				err, gameLogs = c.selectGameLogs(lastId, limit)
			}
		}

	} else {
		if room != "" {
			err, gameLogs = c.selectGameRoomLogs(lastId, limit, room)
		} else {
			err, gameLogs = c.selectGameLogs(lastId, limit)
		}
	}

	if err {
		fmt.Println(err)

		msg = msg + " GameLog trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = gameLogs
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_GAME_LOGS

	return c.RenderJSON(response)
}

func (c GameLogCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	gameLog := c.parseGameLog()
	// Ensure the Id is set.
	gameLog.Id = id
	success, err := c.Txn.Update(&gameLog)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update game_log."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c GameLogCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.GameLog{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove game_log"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
