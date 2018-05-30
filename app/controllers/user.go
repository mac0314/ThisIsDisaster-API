package controllers

import (
	"encoding/json"

	"github.com/revel/revel"
)

type User struct {
	*revel.Controller
}

func (c User) Index() revel.Result {
	// TODO modify demo data
	var code int = 200
	var msg string = "Success"
	var nickname string = "mac"

	// JSON response
	response := make(map[string]interface{})
	response["result_code"] = code
	response["result_msg"] = msg
	data := make(map[string]interface{})

	data["nickname"] = nickname
	data["email"] = "admin@thisisdisaster.com"
	data["signin_time"] = "2018-05-06T15:04:05Z07:00"

	response["result_data"] = data

	return c.RenderJSON(response)
}

func (c User) Lobby() revel.Result {
	// TODO modify demo data
	var code int = 200
	var msg string = "Success"
	var nickname string = "mac"
	var rType string = "lobby"

	// JSON response
	response := make(map[string]interface{})
	response["result_code"] = code
	response["result_msg"] = msg
	response["response_type"] = rType
	data := make(map[string]interface{})

	data["nickname"] = nickname
	data["level"] = "15"
	data["exp"] = "12412"
	data["score"] = "980"
	data["gold"] = "9000"

	response["result_data"] = data

	return c.RenderJSON(response)
}

func (c User) Costume() revel.Result {
	// TODO modify demo data
	//var code int = 200
	//var msg string = "Success"

	/*
		// JSON response
		response := make(map[string]interface{})
		response["result_code"] = code
		response["result_msg"] = msg

		data := make(map[string]interface{})

		costumes := []interface{}{}

		costume := make(map[string]interface{}) // 문자열을 키로하고 모든 자료형을 저장할 수 있는 맵 생성

		costume["name"] = "costume01"
		costume["info"] = "abc"

		costumes = costume

		data["costume_list"] = costumes

		response["result_data"] = data
	*/

	// Demo data
	testData := `
	{
			"result_code": 200,
			"result_msg": "Success",
			"result_data": {
				"gold": 100000,
				"costume_list": [
					{
						"name": "costume01",
						"info": "Basic costume"
					},
					{
						"name": "costume02",
						"info": "Unique costume"
					}
				]
			}
		}
		`

	var response map[string]interface{} // JSON 문서의 데이터를 저장할 공간을 맵으로 선언

	json.Unmarshal([]byte(testData), &response) // doc를 바이트 슬라이스로 변환하여 넣고,
	// data의 포인터를 넣어줌

	return c.RenderJSON(response)
}

func (c User) Score() revel.Result {
	// TODO modify demo data
	// Demo data
	testData := `
	{
			"result_code": 200,
			"result_msg": "Success",
			"result_data": {
				"achievement_list": [
					{
						"name": "Fashion newbie",
						"info": "Get First costume",
						"score": 10
					},
					{
						"name": "Fashionista",
						"info": "Get Unique costume 5",
						"score": 10
					},
					{
						"name": "Beginner",
						"info": "Clear tutorial",
						"score": 10
					},
					{
						"name": "Explorer",
						"info": "Clear all stage",
						"score": 10
					},
					{
						"name": "Slayer",
						"info": "Kill 1000 monsters",
						"score": 10
					},
					{
						"name": "Coward",
						"info": "Hide in shelter",
						"score": 10
					}
				]
			}
		}
		`

	var response map[string]interface{} // JSON 문서의 데이터를 저장할 공간을 맵으로 선언

	json.Unmarshal([]byte(testData), &response) // doc를 바이트 슬라이스로 변환하여 넣고,
	// data의 포인터를 넣어줌

	return c.RenderJSON(response)
}

func (c User) Matching() revel.Result {
	// TODO modify demo data
	// Demo data
	testData := `
	{
		"result_code": 200,
		"result_msg": "Success"
	}
	`

	var response map[string]interface{} // JSON 문서의 데이터를 저장할 공간을 맵으로 선언

	json.Unmarshal([]byte(testData), &response) // doc를 바이트 슬라이스로 변환하여 넣고,
	// data의 포인터를 넣어줌

	return c.RenderJSON(response)
}
