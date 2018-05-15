package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type MonsterCtrl struct {
	GorpController
}

func defineMonsterTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Monster{}, "monster").SetKeys(true, "monster_id")
	// e.g. VARCHAR(25)
	t.ColMap("name_mn").SetMaxSize(30)
	t.ColMap("info_ln").SetMaxSize(255)
	t.ColMap("resource_mn").SetMaxSize(50)
}

func (c MonsterCtrl) parseMonster() models.Monster {
	var jsonData models.Monster

	fmt.Println("parseMonster")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c MonsterCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	monster := c.parseMonster()
	fmt.Println(monster)
	// Validate the model
	monster.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your monster."
	} else {
		if err := c.Txn.Insert(&monster); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = monster
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c MonsterCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	monster := new(models.Monster)
	err := c.Txn.SelectOne(monster,
		`SELECT * FROM monster WHERE monster_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error monster probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = monster
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c MonsterCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	monster, err := c.Txn.Select(models.Monster{},
		`SELECT * FROM monster WHERE monster_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = monster
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c MonsterCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	monster := c.parseMonster()
	// Ensure the Id is set.
	monster.Id = id
	success, err := c.Txn.Update(&monster)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update monster."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = monster
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c MonsterCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Monster{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove monster"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
