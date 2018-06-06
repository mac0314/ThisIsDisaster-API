package controllers

import (
	"github.com/revel/revel"
)

type Admin struct {
	*revel.Controller
}

func (c Admin) Index() revel.Result {
	return c.Render()
}

func (c Admin) Members() revel.Result {
	return c.Render()
}

func (c Admin) Score() revel.Result {
	return c.Render()
}

func (c Admin) Games() revel.Result {
	return c.Render()
}

func (c Admin) Feedbacks() revel.Result {
	return c.Render()
}
