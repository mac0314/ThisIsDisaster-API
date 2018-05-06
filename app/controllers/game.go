package controllers

import (
	"encoding/json"

	"github.com/revel/revel"
)

type Game struct {
	*revel.Controller
}

func (c Game) SinglePlay(nickname string) revel.Result {
	// TODO modify demo data
	testData := `
  {
    "result_code": "200",
    "result_msg": "Success",
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

func (c Game) MultiPlay(nickname string) revel.Result {
	// TODO modify demo data
	testData := `
  {
    "result_code": "200",
    "result_msg": "Success",
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
