package controllers

import (
	"fmt"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	fmt.Println(c.Session["email"])
	fmt.Println(c.Session["role"])

	delete(c.Session, "email")
	delete(c.Session, "role")

	var signinCheck bool

	if c.Session["email"] != "" {
		signinCheck = true
	} else {
		signinCheck = false
	}

	return c.Render(signinCheck)
}

func (c App) News() revel.Result {
	return c.Render()
}

func (c App) Overview() revel.Result {
	return c.Render()
}

func (c App) Media() revel.Result {
	return c.Render()
}

func (c App) Community() revel.Result {
	return c.Render()
}

func (c App) MultiPlay() revel.Result {
	return c.Render()
}

// /robots.txt - Only allow spiders on prod site
func (c App) RobotsTxt() revel.Result {

	txt := "User-agent: *\n"
	if revel.Config.BoolDefault("site.live", false) == false {
		txt += "Disallow: /\n"
	}
	txt += "\n"

	return c.RenderText(txt)
}
