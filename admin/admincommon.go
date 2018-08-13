package admin

import (
	"errors"
	// "html/template"
	// "fmt"
	// "log"
	// "sort"
	// "strings"

	"github.com/gorilla/schema"
	"github.com/kataras/iris"
	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"
	// "ssafa/crypto"
	// "ssafa/db"
	// "ssafa/users"
	// "ssafa/utils"
	// "ssafa/utils"
)

type (
	managePerson struct {
		ID         int
		Name       string
		First      string
		Surname    string
		Role       string
		Based      string
		UserName   string
		Admin      bool
		AdminStr   string
		AdminLink  string
		ActiveStr  string
		ActiveLink string
	}

	personAdd struct {
		ID       int    `schema:"id"`
		First    string `schema:"first"`
		Surname  string `schema:"surname"`
		UserName string `schema:"username"`
		EMail    string `schema:"email"`
		Phone    string `schema:"phone"`
		Mobile   string `schema:"mobile"`
		Role     string `schema:"role"`
		Based    string `schema:"based"`
		Address  string `schema:"address"`
		Postcode string `schema:"postcode"`
		Admin    bool   `schema:"-"`
		AdminStr string `schema:"admin"`
		Commit   string `schema:"commit"`
	}

	manageList []managePerson
	BySurname  struct{ manageList }
	ByFirst    struct{ manageList }
	ByRole     struct{ manageList }
	ByAdmin    struct{ manageList }
	ByBased    struct{ manageList }
	ByActive   struct{ manageList }
)

var (
	decoder = schema.NewDecoder()
)

var (
	errorNoFirst       = errors.New("Please enter a first name")
	errorNoSurname     = errors.New("Please enter a surname")
	errorNoEMail       = errors.New("Please enter an e-mail address")
	errorBadEMail      = errors.New("Please enter a valid e-mail address")
	errorNoRole        = errors.New("Please enter a position")
	errorNoUsername    = errors.New("Please enter a username")
	errorShortUsername = errors.New("The username has to be at least 4 characters")
	errorUsernameUsed  = errors.New("The username is already in use")
)

func SetRoutes(app *iris.Application) {
	app.Get("/admin", adminMain)
	app.Get("/adminperson", adminPerson)
	app.Get("/adminperson/{sortnum:int}", adminPerson)
	app.Get("/adminpersonadd", adminAddPerson)
	app.Post("/adminpersonadd", adminAddPerson)
	app.Get("/adminpersonreset", adminReset)

	app.Get("/adminswap/{usernum:int}/{sortnum:int}", adminSwap)
	app.Get("/activeswap/{usernum:int}/{sortnum:int}", activeSwap)
}
