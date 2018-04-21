package controllers

import (
	"github.com/revel/revel"
)

type User struct {
	*revel.Controller
}

func (c User) Index() revel.Result {
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
