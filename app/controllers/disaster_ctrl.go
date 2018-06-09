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

type DisasterCtrl struct {
	GorpController
}

type Disasters struct {
	Disaster []models.Disaster
}

func defineDisasterTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Disaster{}, "disaster").SetKeys(true, "disaster_id")
	// e.g. VARCHAR(25)
	t.ColMap("name_mn").SetMaxSize(30)
	t.ColMap("info_ln").SetMaxSize(255)
}

func (c DisasterCtrl) parseDisaster() models.Disaster {
	var jsonData models.Disaster

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c DisasterCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	disaster := c.parseDisaster()

	disaster.Create = makeTimestamp()

	disaster.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your disaster."
	} else {
		if err := c.Txn.Insert(&disaster); err != nil {
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

func (c DisasterCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	disaster := new(models.Disaster)
	err := c.Txn.SelectOne(disaster,
		`SELECT * FROM disaster WHERE disaster_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error disaster probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = disaster
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_DISASTER

	return c.RenderJSON(response)
}

func (c DisasterCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	disaster, err := c.Txn.Select(models.Disaster{},
		`SELECT * FROM disaster WHERE disaster_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = disaster
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_DISASTERS

	return c.RenderJSON(response)
}

func (c DisasterCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	disaster := c.parseDisaster()
	// Ensure the Id is set.
	disaster.Id = id
	success, err := c.Txn.Update(&disaster)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update disaster."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c DisasterCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Disaster{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove disaster"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c DisasterCtrl) Load() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	fmt.Println(path.Join(revel.BasePath, "/app/models/data/disaster.xml"))
	// xml 파일 오픈
	fp, err := os.Open(path.Join(revel.BasePath, "/app/models/data/disaster.xml"))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// xml 파일 읽기
	xmlData, err := ioutil.ReadAll(fp)

	// xml 디코딩
	var disasterData Disasters
	xmlerr := xml.Unmarshal(xmlData, &disasterData)
	fmt.Println(disasterData)
	if xmlerr != nil {
		panic(xmlerr)
	} else {
		createTime := makeTimestamp()
		for _, disaster := range disasterData.Disaster {
			disaster.Validate(c.Validation)

			disaster.Create = createTime
			if c.Validation.HasErrors() {
				msg = msg + " You have error in your disaster."
			} else {
				if err := c.Txn.Insert(&disaster); err != nil {
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
