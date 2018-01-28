package cases

import (
	"errors"
	"html/template"
	// "fmt"
	// "log"
	// "sort"
	// "strings"

	"github.com/gorilla/schema"
	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/users"
	"ssafa/utils"
)

type (
	CaseList struct {
		Id         int
		Name       string
		ClientId   int
		CaseNumber string
		CMSNumber  string
		CaseWorker string
		Opened     string
		Updated    string
		State      string
	}

	CaseRec struct {
		Id          int         `bson:"_id"`
		ClientNum   int         `bson:"clientnum"`
		CaseNumber  string      `bson:"case,omitempty"`
		CMSId       string      `bson:"cms,omitempty"`
		CaseFirst   string      `bson:"casefirst,omitempty"`
		CaseSurname string      `bson:"casesurn,omitempty"`
		Closed      bool        `bson:"closed"`
		CWNum       int         `bson:"cwnum"`
		Created     int         `bson:"create,omitempty"`
		Updated     int         `bson:"updated,omitempty"`
		Client      []db.Client `bson:"client"`
		User        []db.User   `bson:"user"`
	}

	CommentDisplay struct {
		Date    string
		Comment template.HTML
		Name    template.HTML
	}

	CommentRec struct {
		Id      int    `schema:"id"`
		Comment string `schema:"comment"`
		Commit  string `schema:"commit"`
	}
)

var (
	key     []byte
	decoder = schema.NewDecoder()

	ErrorDateCaseUsed = errors.New("That case number is used elsewhere")
	ErrorDateCMSUsed  = errors.New("That CMS number is used elsewhere")
)

func SetKey(theKey []byte) {
	key = theKey
}

func SetRoutes(app *iris.Application) {
	app.Get("/case/{casenum:int}", showCase)
	app.Get("/cases", listCases)
	app.Get("/cases/{page:int}", listCases)
	app.Get("/casesopen", openCases)
	app.Get("/casesopen/{page:int}", openCases)
	app.Get("/casesunassign", unassignedCases)
	app.Get("/casesunassign/{page:int}", unassignedCases)
	app.Get("/casesinactive", inactiveCases)
	app.Get("/casesinactive/{page:int}", inactiveCases)

	app.Get("/caseclosed/{casenum:int}", closeCase)
	app.Get("/caseopened/{casenum:int}", openCase)

	app.Get("/caseedit/{casenum:int}", editCase)
	app.Post("/caseedit/{casenum:int}", editCase)

	app.Post("/commentcase/{casenum:int}", addComment)
}

func (cd *CommentDisplay) GetCommentDisplay(comment db.Comment) {
	cd.Date = utils.DateToString(comment.Date)
	cd.Comment = template.HTML(crypto.Decrypt(comment.Comment))

	theStr := ""
	if comment.User != 0 {
		theStr = users.GetUserName(comment.User)
	}
	if len(theStr) == 0 {
		theStr = comment.Name
	}
	cd.Name = template.HTML(theStr)
}

func getClientName(clientNum int) string {
	var (
		theClient db.Client
	)
	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	clientColl.FindId(clientNum).One(&theClient)
	return crypto.Decrypt(theClient.First) + " " + crypto.Decrypt(theClient.Surname)
}
