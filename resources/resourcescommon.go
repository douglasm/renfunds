package resources

import (
	"errors"
	"html/template"

	"github.com/gorilla/schema"
	"github.com/kataras/iris"
)

type (
	createResource struct {
		ID            int    `schema:"id"`
		Name          string `schema:"-"`
		Amount        string `schema:"amount"`
		Value         int    `schema:"-"`
		Establishment string `schema:"establish"`
		RENNumber     string `schema:"-"`
		Client        int    `schema:"-"`
		Commit        string `schema:"commit"`
	}

	editResource struct {
		ID         int    `schema:"id"`
		Name       string `schema:"name"`
		Contact    string `schema:"contact"`
		EMail      string `schema:"email"`
		Phone      string `schema:"phone"`
		URL        string `schema:"url"`
		Comment    string `schema:"comment"`
		Address    string `schema:"address"`
		Checkfield string `schema:"checkfield"`
		Commit     string `schema:"commit"`
	}

	listItem struct {
		ID      int
		Name    string
		Contact string
		EMail   string
		Phone   string
	}

	resourceShow struct {
		ID       int
		Name     string
		Contact  string
		EMail    string
		Phone    string
		URL      string
		Link     string
		Comments template.HTML
		Address  template.HTML
		Updated  string
	}
)

var (
	decoder = schema.NewDecoder()

	errNoName = errors.New("You must enter a resource name")
)

// SetRoutes will set where the calls go to.
func SetRoutes(app *iris.Application) {
	app.Get("/resources", listResources)
	app.Get("/resources/{pagenum:int}", listResources)
	app.Get("/resource/{resourcenum:int}", showResource)
	app.Get("/resourceadd", addResource)
	app.Post("/resourceadd", addResource)
	app.Get("/resourceedit/{resourcenum:int}", resourceEdit)
	app.Post("/resourceedit/{resourcenum:int}", resourceEdit)
	app.Get("/resourcedelete/{resourcenum:int}", resourceDelete)
	app.Get("/resourceremove/{resourcenum:int}", resourceRemove)
}
