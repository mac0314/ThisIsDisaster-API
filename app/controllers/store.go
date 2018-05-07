package controllers

import (
	"encoding/json"

	"github.com/revel/revel"
)

type Store struct {
	*revel.Controller
}

func (c Store) Index() revel.Result {
	// TODO modify demo data
	// Demo data
	testData := `
	{
			"result_code": "200",
			"result_msg": "Success",
			"result_data": {
				"costume_list": [
					{
						"name": "costume01",
						"info": "Basic costume",
            "cost": 1
					},
					{
						"name": "costume02",
						"info": "Unique costume",
            "cost": 1000
					},
					{
						"name": "costume03",
						"info": "Unique costume",
            "cost": 1500
					},
					{
						"name": "costume04",
						"info": "Unique costume",
            "cost": 3000
					},
					{
						"name": "costume05",
						"info": "Unique costume",
            "cost": 2000
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

func (c Store) BuyCostume() revel.Result {
	// TODO modify demo data
	var code string = "200"
	var msg string = "Success"

	// JSON response
	response := make(map[string]interface{})
	response["result_code"] = code
	response["result_msg"] = msg

	return c.RenderJSON(response)
}
