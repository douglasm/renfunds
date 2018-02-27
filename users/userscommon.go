package users

import (
	// "fmt"

	"github.com/gorilla/schema"
	"github.com/kataras/iris"
	// "ssafa/crypto"
	// "ssafa/types"
)

const (
	KLoginForm = 150 + iota
)

type (
	User struct {
		Name  string
		Towns []string
		Num   int
	}

	Session struct {
		Id         string `bson:"_id"`
		UserNumber int    `bson:"usernum"`
		LoggedIn   bool   `bson:"logged,omitempty"`
		Admin      bool   `bson:"admin,omitempty"`
	}

	userChoice struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	userRequest struct {
		SearchTerm string `schema:"search"`
	}
)

var (
	decoder = schema.NewDecoder()
)

func SetRoutes(app *iris.Application) {
	app.Get("/login", login)
	app.Post("/login", login)
	app.Get("/logout", logout)

	app.Post("/usersget", getUserList)

	app.Get("/user/{usernum:int}", showUser)
	app.Get("/useredit/{usernum:int}", editUser)
	app.Post("/useredit/{usernum:int}", editUser)
}
