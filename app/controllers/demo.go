package controllers

import (
	"github.com/revel/revel"
)

type Demo struct {
	*revel.Controller
}

func (c Demo) Success() revel.Result {
	var code int = 200
	var msg string = "Success"
	var rType string = "Response"

	// JSON response
	data := make(map[string]interface{})
	data["result_code"] = code
	data["result_msg"] = msg
	data["response_type"] = rType

	return c.RenderJSON(data)
}

func (c Demo) Fail() revel.Result {
	var code int = 400
	var msg string = "Fail"
	var rType string = "response"

	// JSON response
	data := make(map[string]interface{})
	data["result_code"] = code
	data["result_msg"] = msg
	data["response_type"] = rType

	return c.RenderJSON(data)
}
