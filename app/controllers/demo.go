package controllers

import (
	"github.com/revel/revel"
)

type Demo struct {
	*revel.Controller
}

func (c Demo) Success() revel.Result {
	var code string = "200"
	var msg string = "Success"

	// JSON response
	data := make(map[string]interface{})
	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}

func (c Demo) Fail() revel.Result {
	var code string = "400"
	var msg string = "Fail"

	// JSON response
	data := make(map[string]interface{})
	data["result_code"] = code
	data["result_msg"] = msg

	return c.RenderJSON(data)
}
