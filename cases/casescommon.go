package cases

import (
	// "log"
	// "fmt"
	// "sort"
	// "strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "ssafa/crypto"
	// "ssafa/db"
)

type (
	CaseList struct {
		Id         int
		CaseNumber string
		CMSNumber  string
		CaseWorker string
		State      string
	}
)

var (
	key []byte
)

func SetKey(theKey []byte) {
	key = theKey
}

func SetRoutes(app *iris.Application) {
	app.Get("/case/{casenum:int}", showCase)
}
