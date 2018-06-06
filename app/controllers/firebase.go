package controllers

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/revel/revel"

	"google.golang.org/api/option"
)

type Firebase struct {
	*revel.Controller
}

func (c Firebase) Index() revel.Result {
	opt := option.WithCredentialsFile(revel.SourcePath + "/ThisisDisaster-API/conf/thisisdisaster-204407-firebase-adminsdk-zvtoj-0c39c33f36.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access Auth service from default app
	defaultClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	user, _ := defaultClient.GetUser(context.Background(), "3OM8rAkdlYV7EPxGkWrMGyPjJs53")

	fmt.Println(user.Email)

	c.Session["email"] = user.Email
	c.Session["role"] = "admin"

	return c.RenderText("ok")
}
