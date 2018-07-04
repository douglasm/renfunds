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
	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/users"
	"ssafa/utils"
)

// Check edit button for case Comments
// Check preserving of case worker details when editing case

type (
	CaseList struct {
		ID         int
		Name       string
		ClientID   int
		RENNumber  string
		CMSNumber  string
		CaseWorker string
		Opened     string
		Updated    string
		State      string
	}

	CaseRec struct {
		ID          int         `bson:"_id"`
		ClientNum   int         `bson:"clientnum"`
		RENNumber   string      `bson:"case,omitempty"`
		CMSID       string      `bson:"cms,omitempty"`
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
		Item    int
		Num     int
	}

	CommentRec struct {
		ID      int    `schema:"id"`
		Comment string `schema:"comment"`
		Commit  string `schema:"commit"`
	}
)

var (
	key     []byte
	decoder = schema.NewDecoder()

	ErrorDateCaseUsed  = errors.New("That case number is used elsewhere")
	ErrorDateCMSUsed   = errors.New("That CMS number is used elsewhere")
	errorCMSMissing    = errors.New("You must enter a CMS number.")
	errorReasonMissing = errors.New("You must enter a reason for the case.")
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

	app.Any("/caseadd/{clientnum:int}", createCase)

	app.Get("/caseedit/{casenum:int}", editCase)
	app.Post("/caseedit/{casenum:int}", editCase)

	app.Get("/casedelete/{casenum:int}", deleteCase)
	app.Get("/caseremove/{casenum:int}", removeCase)

	app.Post("/searchcase", searchCase)

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
