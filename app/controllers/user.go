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
	var code string = "200"
	var msg string = "Success"
	var nickname string = "mac"

	// JSON response
	response := make(map[string]interface{})
	response["result_code"] = code
	response["result_msg"] = msg
	data := make(map[string]interface{})

	data["nickname"] = nickname
	data["score"] = "1000"
	data["level"] = "15"
	data["gold"] = "90000"

	response["result_data"] = data

	return c.RenderJSON(response)
}

func (c User) Lobby(nickname string) revel.Result {
	// TODO modify demo data
	var code string = "200"
	var msg string = "Success"

	// JSON response
	response := make(map[string]interface{})
	response["result_code"] = code
	response["result_msg"] = msg
	data := make(map[string]interface{})

	data["nickname"] = nickname
	data["level"] = "15"
	data["exp"] = "12412"

	response["result_data"] = data

	return c.RenderJSON(response)
}

func (c User) Costume(nickname string) revel.Result {
	// TODO modify demo data
	//var code string = "200"
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
			"result_code": "200",
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
