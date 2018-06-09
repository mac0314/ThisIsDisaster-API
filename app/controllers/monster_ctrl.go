package controllers

import (
	"ThisIsDisaster-API/app/models"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type MonsterCtrl struct {
	GorpController
}

type Monsters struct {
	Monster []models.Monster
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

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c MonsterCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	monster := c.parseMonster()

	monster.Create = makeTimestamp()

	monster.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your monster."
	} else {
		if err := c.Txn.Insert(&monster); err != nil {
			fmt.Println(err)

			msg = msg + " Error inserting record into database!"
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

func (c MonsterCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	monster := new(models.Monster)
	err := c.Txn.SelectOne(monster,
		`SELECT * FROM monster WHERE monster_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error monster probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = monster
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MONSTER

	return c.RenderJSON(response)
}

func (c MonsterCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	monster, err := c.Txn.Select(models.Monster{},
		`SELECT * FROM monster WHERE monster_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = monster
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_MONSTERS

	return c.RenderJSON(response)
}

func (c MonsterCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	monster := c.parseMonster()
	// Ensure the Id is set.
	monster.Id = id
	success, err := c.Txn.Update(&monster)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update monster."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c MonsterCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Monster{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove monster"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c MonsterCtrl) Load() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	fmt.Println(path.Join(revel.BasePath, "/app/models/data/monster.xml"))
	// xml 파일 오픈
	fp, err := os.Open(path.Join(revel.BasePath, "/app/models/data/monster.xml"))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// xml 파일 읽기
	xmlData, err := ioutil.ReadAll(fp)

	// xml 디코딩
	var monsterData Monsters
	xmlerr := xml.Unmarshal(xmlData, &monsterData)
	fmt.Println(monsterData)
	if xmlerr != nil {
		panic(xmlerr)
	} else {
		createTime := makeTimestamp()
		for _, monster := range monsterData.Monster {
			monster.Validate(c.Validation)

			monster.Create = createTime
			if c.Validation.HasErrors() {
				msg = msg + " You have error in your monster."
			} else {
				if err := c.Txn.Insert(&monster); err != nil {
					fmt.Println(err)

					msg = msg + " Error inserting record into database!"
				} else {

					code = RESULT_CODE_SUCCESS
					msg = "Success."
				}
			}
		}

	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}
