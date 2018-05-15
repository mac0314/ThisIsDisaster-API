package tests

import (
	"github.com/revel/revel/testing"
)

type WebTest struct {
	testing.TestSuite
}

type APITest struct {
	testing.TestSuite
}

func (t *WebTest) Before() {
	println("Web Set up")
}

func (t *WebTest) After() {
	println("Web Tear down")
}

// Test Web routes
func (t *WebTest) TestIndexPage() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *WebTest) TestNewsPage() {
	t.Get("/news")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *WebTest) TestOverviewPage() {
	t.Get("/overview")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *WebTest) TestMediaPage() {
	t.Get("/media")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *WebTest) TestCommunityPage() {
	t.Get("/community")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

// robots.txt
func (t *WebTest) TestRobotsPage() {
	t.Get("/robots.txt")
	t.AssertOk()
	t.AssertContentType("text/plain; charset=utf-8")
}

// Will not appear in panel as it does not start with `Test`.
func (t *WebTest) TestFavicon() {
	t.Get("/favicon.ico")
	t.AssertOk()
	t.AssertContentType("image/x-icon")
}

func (t *APITest) Before() {
	println("API Set up")
}

func (t *APITest) After() {
	println("API Tear down")
}

// API routes
func (t *APITest) TestUser() {
	t.Get("/api/user")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *APITest) TestGetUserId() {
	t.Get("/api/user/:id")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *APITest) TestGetUsers() {
	t.Get("/api/users")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *APITest) TestPostUser() {
	/*
		// Success case
		user := models.User{1, "admin@thisisdisaster.com", "admin", 150, 99999, 99999, 1525685748965}
		ubytes, _ := json.Marshal(user)
		buff := bytes.NewBuffer(ubytes)

		t.Post("/api/user", "application/json", buff)
	*/
	// Failure case
	t.Post("/api/user", "application/json", nil)
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}
