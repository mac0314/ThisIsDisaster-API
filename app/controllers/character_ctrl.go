package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type CharacterCtrl struct {
	GorpController
}

func defineCharacterTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Character{}, "u_character").SetKeys(true, "u_character_id")
	t.ColMap("name_sn").SetMaxSize(30)
}

func (c CharacterCtrl) parseCharacter() models.Character {
	var jsonData models.Character

	fmt.Println("parseCharacter")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c CharacterCtrl) Add() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	character := c.parseCharacter()
	fmt.Println(character)
	// Validate the model
	character.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your character."
	} else {
		if err := c.Txn.Insert(&character); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
		} else {
			code = 200
			msg = "Success."
			data["result_data"] = character
		}
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c CharacterCtrl) Get(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	character := new(models.Character)
	err := c.Txn.SelectOne(character,
		`SELECT * FROM u_character WHERE u_character_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error character probably doesn't exist."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = character
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c CharacterCtrl) List() revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	character, err := c.Txn.Select(models.Character{},
		`SELECT * FROM u_character WHERE u_character_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = 200
		msg = "Success."
		data["result_data"] = character
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c CharacterCtrl) Update(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	character := c.parseCharacter()
	// Ensure the Id is set.
	character.Id = id
	success, err := c.Txn.Update(&character)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update character."
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
		data["result_data"] = character
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c CharacterCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := 400
	msg := "Fail."
	data := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Character{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove character"
	} else {
		code = 200
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
