package controllers

import (
	"ThisIsDisaster-API/app/models"
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
)

type Game struct {
	*revel.Controller
	Matching
}

func (c Game) SinglePlayResult() revel.Result {
	// TODO modify demo data
	testData := `
  {
    "result_code": 200,
    "result_msg": "Success",
		"result_type": "SinglePlayResult",
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

func (c Game) parseUser() models.User {
	var jsonData models.User

	c.Params.BindJSON(&jsonData)

	return jsonData
}

func (c Game) MultiPlayJoin() revel.Result {
	user := c.parseUser()

	c.UpdateIP(user.Email, user.IP)

	GoMatching(user)

	response := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
		"result_type": RESULT_TYPE_RESPONSE,
	}

	return c.RenderJSON(response)
}

func (c Game) MultiPlayLobby() revel.Result {
	email := c.Params.Get("email")

	var userList = []interface{}{}

	room, users := c.GetMyMatchingRoom(email)

	fmt.Println(room)

	fmt.Println(len(users))
	fmt.Println(users)

	if len(users) > 0 {
		host := LoadHost(room)

		for _, user := range users {
			var role string
			if host == user.Email {
				role = "host"
			} else {
				role = "client"
			}

			userData := map[string]interface{}{
				"id":       user.Id,
				"email":    user.Email,
				"nickname": user.Nickname,
				"level":    user.Level,
				"ip":       user.IP,
				"role":     role,
			}

			userList = append(userList, userData)
		}
	}

	response := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
		"result_type": "MultiPlayLobby",
		"result_data": map[string]interface{}{
			"email":     email,
			"stage":     room,
			"user_list": userList,
		},
	}

	return c.RenderJSON(response)
}

func (c Game) LeaveMultiPlayLobby() revel.Result {
	email := c.Params.Get("email")
	fmt.Println(email)

	LeaveMatchingRoom(email)

	response := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
		"result_type": RESULT_TYPE_RESPONSE,
	}

	return c.RenderJSON(response)
}

func (c Game) StartGame() revel.Result {
	email := c.Params.Get("email")
	mode := c.Params.Get("mode")

	room, _ := c.GetMyMatchingRoom(email)

	ClearMatchingRoom(room)

	createTime := makeTimestamp()

	_stage := &models.Stage{0, room, mode, createTime}

	c.Txn.Insert(_stage)

	response := map[string]interface{}{
		"result_code": 200,
		"result_msg":  "Success",
		"result_type": RESULT_TYPE_RESPONSE,
	}

	return c.RenderJSON(response)
}

func (c Game) MultiPlayResult() revel.Result {
	// TODO modify demo data
	testData := `
  {
    "result_code": 200,
    "result_msg": "Success",
		"result_type": "MultiPlayResult",
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
