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

type ItemCtrl struct {
	GorpController
}

type Items struct {
	Item []models.Item
}

func defineItemTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTableWithName(models.Item{}, "item").SetKeys(true, "item_id")
	// e.g. VARCHAR(25)
	t.ColMap("name_sn").SetMaxSize(20)
	t.ColMap("type_sn").SetMaxSize(20)
	t.ColMap("description_ln").SetMaxSize(255)
	t.ColMap("resource_mn").SetMaxSize(50)
}

func (c ItemCtrl) parseItem() models.Item {
	var jsonData models.Item

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c ItemCtrl) Add() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	item := c.parseItem()

	item.Create = makeTimestamp()

	item.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg = msg + " You have error in your item."
	} else {
		if err := c.Txn.Insert(&item); err != nil {
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

func (c ItemCtrl) Get(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	item := new(models.Item)
	err := c.Txn.SelectOne(item,
		`SELECT * FROM item WHERE item_id = ?`, id)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error item probably doesn't exist."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = item
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_ITEM

	return c.RenderJSON(response)
}

func (c ItemCtrl) List() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	item, err := c.Txn.Select(models.Item{},
		`SELECT * FROM item WHERE item_id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		fmt.Println(err)

		msg = msg + " Error trying to get records from DB."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success."
		response["result_data"] = item
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_ITEMS

	return c.RenderJSON(response)
}

func (c ItemCtrl) Update(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	item := c.parseItem()
	// Ensure the Id is set.
	item.Id = id
	success, err := c.Txn.Update(&item)
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Unable to update item."
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Updated %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c ItemCtrl) Delete(id int64) revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	success, err := c.Txn.Delete(&models.Item{Id: id})
	if err != nil || success == 0 {
		fmt.Println(err)

		msg = msg + " Failed to remove item"
	} else {
		code = RESULT_CODE_SUCCESS
		msg = "Success. " + fmt.Sprintf("Deleted %d", id)
	}

	response["result_code"] = code
	response["result_msg"] = msg
	response["result_type"] = RESULT_TYPE_RESPONSE

	return c.RenderJSON(response)
}

func (c ItemCtrl) Load() revel.Result {
	// JSON response
	code := RESULT_CODE_FAILURE
	msg := "Fail."
	response := make(map[string]interface{})

	fmt.Println(path.Join(revel.BasePath, "/app/models/data/item.xml"))
	// xml 파일 오픈
	fp, err := os.Open(path.Join(revel.BasePath, "/app/models/data/item.xml"))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// xml 파일 읽기
	xmlData, err := ioutil.ReadAll(fp)

	// xml 디코딩
	var itemData Items
	xmlerr := xml.Unmarshal(xmlData, &itemData)
	fmt.Println(itemData)
	if xmlerr != nil {
		panic(xmlerr)
	} else {
		createTime := makeTimestamp()

		for _, item := range itemData.Item {
			item.Validate(c.Validation)
			item.Create = createTime

			if c.Validation.HasErrors() {
				msg = msg + " You have error in your item."
			} else {
				if err := c.Txn.Insert(&item); err != nil {
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
