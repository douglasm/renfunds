package users

import (
	"errors"
	"log"
	// "fmt"
	"html/template"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/schema"
	"github.com/kataras/iris"

	"ssafa/crypto"
	"ssafa/db"
	// "ssafa/types"
)

const (
	kLoginForm = 150 + iota
	kResetForm
	kChangeForm
	kActivateForm

	KActivateLength = 12
	KResetLength    = 14

	KActivateTime = 86400
	KResetTime    = 7200
)

type (
	User struct {
		Name  string
		Towns []string
		Num   int
	}

	resetRec struct {
		Username   string `schema:"un"`
		Pass1      string `schema:"pwd0"`
		Pass2      string `schema:"pwd1"`
		Code       string `schema:"code"`
		HasCode    bool   `schema:"-"`
		Checkfield string `schema:"checkfield"`
		Commit     string `schema:"commit"`
	}

	changePass struct {
		Password   string `schema:"code"`
		Pass1      string `schema:"pwd0"`
		Pass2      string `schema:"pwd1"`
		Number     int    `schema:"num"`
		Checkfield string `schema:"checkfield"`
		Commit     string `schema:"commit"`
	}

	Session struct {
		ID         string `bson:"_id"`
		UserNumber int    `bson:"usernum"`
		LoggedIn   bool   `bson:"logged,omitempty"`
		Admin      bool   `bson:"admin,omitempty"`
	}

	userDetails struct {
		First    string
		Surname  string
		Address  template.HTML
		PostCode template.HTML
		Phone    template.HTML
		Mobile   template.HTML
		EMail    template.HTML
	}

	userChoice struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	userRequest struct {
		SearchTerm string `schema:"search"`
	}
)

var (
	decoder = schema.NewDecoder()

	errEMailUsed    = errors.New("You cannot use an e-mail address as a password")
	errPassMismatch = errors.New("The two passwords do not match")
	errPassUsed     = errors.New("That password has been revealed in a data breach")
	errPassShort    = errors.New("The password is not long enough")
	errPassBad      = errors.New("Incorrect password")
)

func SetRoutes(app *iris.Application) {
	app.Get("/login", login)
	app.Post("/login", login)
	app.Get("/logout", logout)

	app.Post("/usersget", getUserList)

	app.Get("/user/{usernum:int}", showUser)
	app.Get("/useredit/{usernum:int}", editUser)
	app.Post("/useredit/{usernum:int}", editUser)

	app.Get("/activate", activateUser)
	app.Post("/activate", activateUser)
	app.Get("/activate/{code}", activateUser)
	app.Post("/activate/{code}", activateUser)

	app.Get("/me", myDetails)

	app.Get("/changepassword", changePassword)
	app.Post("/changepassword", changePassword)

	app.Get("/resetpassword", resetPassword)
	app.Post("/resetpassword", resetPassword)

	app.Get("/resetsent", resetSent)

	app.Get("/reset", resetHandler)
	app.Get("/reset/{code}", resetHandler)
	app.Post("/reset", resetHandler)
	app.Post("/reset/{code}", resetHandler)
}

func setNewPassword(userID int, thePass string, setActive bool) {
	var (
		err error
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	theSalt := crypto.RandomChars(20)
	sets := bson.M{}
	unSets := bson.M{}

	if setActive {
		sets[db.KFieldUserInactive] = false
	}
	sets[db.KFieldUserSalt] = theSalt
	sets[db.KFieldUserPassword], err = crypto.GetHash(thePass, theSalt)

	unSets[db.KFieldUserActivateCode] = 1
	unSets[db.KFieldUserActivateTime] = 1

	err = usersCollection.UpdateId(userID, bson.M{"$set": sets, "$unset": unSets})
	if err != nil {
		log.Println("Error: setNewPassword", err)
	}
}
