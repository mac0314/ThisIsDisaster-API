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
	character := c.parseCharacter()
	fmt.Println(character)
	// Validate the model
	character.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your character.")
	} else {
		if err := c.Txn.Insert(&character); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(character)
		}
	}

}

func (c CharacterCtrl) Get(id int64) revel.Result {
	character := new(models.Character)
	err := c.Txn.SelectOne(character,
		`SELECT * FROM u_character WHERE u_character_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. character probably doesn't exist.")
	}
	return c.RenderJSON(character)
}

func (c CharacterCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	character, err := c.Txn.Select(models.Character{},
		`SELECT * FROM u_character WHERE u_character_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(character)
}

func (c CharacterCtrl) Update(id int64) revel.Result {
	character := c.parseCharacter()
	// Ensure the Id is set.
	character.Id = id
	success, err := c.Txn.Update(&character)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update character.")
	}
	return c.RenderText("Updated %v", id)
}

func (c CharacterCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Character{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove character")
	}
	return c.RenderText("Deleted %v", id)
}
