package tests

import (
	"github.com/revel/revel/testing"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) After() {
	println("Tear down")
}

// Check if robots.txt exists
func (t *AppTest) TestRobotsPage() {
	t.Get("/robots.txt")
	t.AssertOk()
	t.AssertContentType("text/html")
}

// Will not appear in panel as it does not start with `Test`.
func (t *AppTest) TEstFavIcon() {
	t.Get("/favicon.ico")
	t.AssertOk()
	t.AssertContentType("text/html")
}
