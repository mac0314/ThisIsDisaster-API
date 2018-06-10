package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"
	"strconv"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type UserCtrl struct {
	GorpController
}

func defineUserTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.User{}, "user").SetKeys(true, "user_id")
	// e.g. VARCHAR(25)
	t.ColMap("email_mn").SetMaxSize(30)
	t.ColMap("nickname_mn").SetMaxSize(30)
	t.ColMap("password_ln").SetMaxSize(255)
	t.ColMap("ip_sn").SetMaxSize(20)
}

func (c UserCtrl) parseUser() models.User {
	var jsonData models.User

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c UserCtrl) insertUser(data models.User) (bool, string) {
	var err bool
	var msg string

	data.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your user."

	} else {
		if _err := c.Txn.Insert(&data); _err != nil {
			fmt.Println(_err.Error())

			msg = msg + " Error inserting record into database!"
		} else {
			msg = "Success."
		}
	}
	fmt.Println(msg)

	return err, msg
}

func (c UserCtrl) selectUserById(id int64) (bool, *models.User) {
	var err bool
	user := new(models.User)
	_err := c.Txn.SelectOne(user,
		`SELECT * FROM user WHERE user_id = ?`, id)
	if _err != nil {
		fmt.Println(err)

		err = true
	} else {
		err = false
	}

	return err, user
}

func (c UserCtrl) selectUserByEmail(email string) (bool, *models.User) {
	var err bool
	user := new(models.User)
	_err := c.Txn.SelectOne(user,
		`SELECT * FROM user WHERE email_mn = ?`, email)
	if _err != nil {
		fmt.Println(err)

		err = true
	} else {
		err = false
	}

	return err, user
}

func (c UserCtrl) selectUsers(lastId int64, limit uint64) (bool, []models.User) {
	var err bool
	var users []models.User

	_, _err := c.Txn.Select(&users,
		`SELECT * FROM user WHERE user_id > ? LIMIT ?`, lastId, limit)
	if _err != nil {
		err = true
		fmt.Println(_err.Error())
	} else {
		err = false
	}

	return err, users
}

func (c UserCtrl) Post() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	user := c.parseUser()
	user.Level = 1

	fmt.Println(user)

	var dUser models.User
	_dErr := c.Txn.SelectOne(&dUser, "SELECT * FROM user WHERE email_mn='"+user.Email+"'")
	if _dErr != nil {
		platform := c.Params.Get("platform")

		if platform != "local" {
			user.Password = RandStringBytesMaskImprSrc(7)
		}

		if platform == "" {
			platform = DEFAULT_FLATFORM
		}

		hash, _ := HashPassword(user.Password)
		user.Password = hash

		_err, msg := c.insertUser(user)
		if _err {
			fmt.Println(msg)
		}

		create := makeTimestamp()
		_err, _user := c.selectUserByEmail(user.Email)

		_authorize := &models.Authorize{0, _user.Id, user.Email, platform, create}
		_setting := &models.UserSetting{_user.Id, DEFAULT_SETTING_PUSH, DEFAULT_SOUND, create}
		_hhcostume := &models.HaveHeadCostume{0, _user.Id, DEFAULT_HEAD_COSTUME, DEFAULT_COSTUME_STATE, create}
		_hbcostume := &models.HaveBodyCostume{0, _user.Id, DEFAULT_BODY_COSTUME, DEFAULT_COSTUME_STATE, create}

		c.Txn.Insert(_authorize)
		c.Txn.Insert(_setting)
		c.Txn.Insert(_hhcostume)
		c.Txn.Insert(_hbcostume)

		if !_err {
			code = RESULT_CODE_SUCCESS
		}

	}
	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c UserCtrl) Get() revel.Result {
	var err bool
	var user *models.User

	id, _err := strconv.ParseInt(c.Params.Get("id"), 10, 32)
	if _err == nil {
		err, user = c.selectUserById(id)
	} else {
		email := c.Params.Get("email")
		err, user = c.selectUserByEmail(email)
	}

	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	if err {
		msg = msg + " Error user probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."

		userData := map[string]interface{}{
			"id":       user.Id,
			"email":    user.Email,
			"nickname": user.Nickname,
			"level":    user.Level,
			"exp":      user.Exp,
			"gold":     user.Gold,
			"score":    user.Score,
		}

		response["result_data"] = userData
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_USER

	return c.RenderJSON(response)
}

func (c UserCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))

	err, users := c.selectUsers(lastId, limit)
	if err {
		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."

		var userList = []interface{}{}

		for _, user := range users {
			userData := map[string]interface{}{
				"id":       user.Id,
				"email":    user.Email,
				"nickname": user.Nickname,
				"level":    user.Level,
				"exp":      user.Exp,
				"gold":     user.Gold,
				"score":    user.Score,
			}

			userList = append(userList, userData)
		}
		response["result_data"] = userList
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_USERS

	return c.RenderJSON(response)
}

func (c UserCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	user := c.parseUser()
	// Ensure the Id is set.
	user.Id = id
	success, err := c.Txn.Update(&user)
	if err != nil || success == 0 {
		msg = msg + " Unable to update user."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c UserCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.User{Id: id})
	if err != nil || success == 0 {
		msg = msg + " Failed to remove user"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
