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

func defineAchievementTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Achievement{}, "achievement").SetKeys(true, "achievement_id")
	// e.g. VARCHAR(25)
	t.ColMap("title_mn").SetMaxSize(50)
	t.ColMap("content_ln").SetMaxSize(255)
}

func (c AchievementCtrl) parseAchievement() models.Achievement {
	var jsonData models.Achievement

	fmt.Println("parseAchievement")

	fmt.Println(makeTimestamp())

	c.Params.BindJSON(&jsonData)

	fmt.Println(jsonData)

	return jsonData
}

func (c AchievementCtrl) Add() revel.Result {
	achievement := c.parseAchievement()
	fmt.Println(achievement)
	// Validate the model
	achievement.Validate(c.Validation)
	if c.Validation.HasErrors() {
		// Do something better here!
		return c.RenderText("You have error in your achievement.")
	} else {
		if err := c.Txn.Insert(&achievement); err != nil {
			fmt.Println(err)
			return c.RenderText(
				"Error inserting record into database!")
		} else {
			return c.RenderJSON(achievement)
		}
	}

}

func (c AchievementCtrl) Get(id int64) revel.Result {
	achievement := new(models.Achievement)
	err := c.Txn.SelectOne(achievement,
		`SELECT * FROM achievement WHERE achievement_id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return c.RenderText("Error. achievement probably doesn't exist.")
	}
	return c.RenderJSON(achievement)
}

func (c AchievementCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	achievement, err := c.Txn.Select(models.Achievement{},
		`SELECT * FROM achievement WHERE achievement_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}
	return c.RenderJSON(achievement)
}

func (c AchievementCtrl) Update(id int64) revel.Result {
	achievement := c.parseAchievement()
	// Ensure the Id is set.
	achievement.Id = id
	success, err := c.Txn.Update(&achievement)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update achievement.")
	}
	return c.RenderText("Updated %v", id)
}

func (c AchievementCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Achievement{Id: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove achievement")
	}
	return c.RenderText("Deleted %v", id)
}
