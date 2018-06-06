package controllers

import (
	"ThisIsDisaster-API/app/models"
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
)

type Game struct {
	*revel.Controller
}

func (c Game) SinglePlayResult() revel.Result {
	// TODO modify demo data
	testData := `
  {
    "result_code": 200,
    "result_msg": "Success",
		"response_type": "SingleResult",
    "result_data": {
      "mode": "normal",
      "exp": "3000",
      "rank": "A",
      "disaster":[
        "earthquake", "tsunami"
      ],
      "clear_time": "2018-05-06T15:04:05Z07:00",
      "gold": "1000"
    }
  }
  `

	var response map[string]interface{} // JSON 문서의 데이터를 저장할 공간을 맵으로 선언

	json.Unmarshal([]byte(testData), &response) // doc를 바이트 슬라이스로 변환하여 넣고,
	// data의 포인터를 넣어줌

	return c.RenderJSON(response)
}

func (c Game) parseUser() models.UserLocal {
	var jsonData models.UserLocal

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c Game) MultiPlay() revel.Result {
	testData := c.parseUser()

	fmt.Println(testData)
	/*
		// JSON 문서의 데이터를 저장할 공간을 맵으로 선언
		jsonBytes, err := json.Marshal(testData) // doc를 바이트 슬라이스로 변환하여 넣고,
		if err != nil {
			panic(err)
		}
		user := string(jsonBytes)

		fmt.Println(user)
	*/
	Matching(testData)

	d := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
	}

	return c.RenderJSON(d)
}

func (c Game) MultiPlayLobby() revel.Result {
	// TODO modify demo data

	room := c.Params.Get("room")

	data := GetMatchingRoom(room)

	var userData = []interface{}{}

	if data != nil {
		userData = data
	}

	response := map[string]interface{}{
		"result_code":   200,
		"result_msg":    "Success",
		"response_type": "MultiPlayLobby",
		"result_data": map[string]interface{}{
			"name":      "mac",
			"level":     "15",
			"role":      "client",
			"ip":        "192.168.1.9",
			"user_list": userData,
		},
	}

	/*
			testData := `
		  {
		    "result_code": 200,
		    "result_msg": "Success",
				"response_type": "MultiLobby",
		    "result_data": {
		      "name": "mac",
		      "level": "15",
					"role": "client",
					"ip": "192.168.1.9",
		      "user_list":[
		        {
							"name": "fuck",
				      "level": "20",
							"role": "host",
							"ip": "192.168.1.10"
						},
						{
							"name": "helloworld",
				      "level": "15",
							"role": "client",
							"ip": "192.168.1.11"
						},
						{
							"name": "happyface",
				      "level": "10",
							"role": "client",
							"ip": "192.168.1.12"
						}
		      ]
		    }
		  }
		  `

			var response map[string]interface{} // JSON 문서의 데이터를 저장할 공간을 맵으로 선언

			json.Unmarshal([]byte(testData), &response) // doc를 바이트 슬라이스로 변환하여 넣고,
			// data의 포인터를 넣어줌
	*/
	return c.RenderJSON(response)
}

func (c Game) LeaveMultiPlayLobby() revel.Result {
	email := c.Params.Get("email")

	LeaveMatchingRoom(email)

	response := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
	}

	return c.RenderJSON(response)
}

func (c Game) StartGame() revel.Result {
	room := c.Params.Get("room")

	ClearMatchingRoom(room)

	response := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
	}

	return c.RenderJSON(response)
}

func (c Game) MultiPlayResult() revel.Result {
	// TODO modify demo data
	testData := `
  {
    "result_code": 200,
    "result_msg": "Success",
		"response_type": "MultiResult",
    "result_data": {
      "mode": "normal",
      "exp": "3000",
      "rank": "A",
      "disaster":[
        "earthquake", "tsunami"
      ],
      "clear_time": "2018-05-06T15:04:05Z07:00",
      "gold": "1000"
    }
  }
  `

	var response map[string]interface{} // JSON 문서의 데이터를 저장할 공간을 맵으로 선언

	json.Unmarshal([]byte(testData), &response) // doc를 바이트 슬라이스로 변환하여 넣고,
	// data의 포인터를 넣어줌

	return c.RenderJSON(response)
}
