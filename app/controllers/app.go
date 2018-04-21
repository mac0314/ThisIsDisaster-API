package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
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
