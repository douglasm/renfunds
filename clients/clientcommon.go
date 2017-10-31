package clients

import (
	// "log"
	// "fmt"
	// "sort"
	// "strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"
	// "ssafa/crypto"
	// "ssafa/db"
)

func SetRoutes(app *iris.Application) {
	app.Get("/client/{clientnum:int}", showClient)
	// app.Post("/login", login)
}
