package users

import (
	// "fmt"
	"log"
	"net/http"
	"time"

	// "github.com/gorilla/schema"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"ssafa/cookie"
	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
)

type (
	LoginRecord struct {
		Username   string `schema:"username"`
		Password   string `schema:"password"`
		Remember   bool   `schema:"remember"`
		Checkfield string `schema:"checkfield"`
		Commit     string `schema:"commit"`
	}
)

func login(ctx iris.Context) {
	var (
		theSession   Session
		details      LoginRecord
		cookieString string
		theUser      db.User
		err          error
	)

	theSession = ctx.Values().Get("session").(Session)
	header := types.HeaderRecord{Title: "Renfunds login"}
	header.Scripts = append(header.Scripts, "passwordtoggle")

	if ctx.Method() == "POST" {
		data := ctx.FormValues()
		err = decoder.Decode(&details, data)
		if err == nil {
			session := db.MongoSession.Copy()
			defer session.Close()

			usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
			sessionCollection := session.DB(db.MainDB).C(db.CollectionSessions)
			err = usersCollection.Find(bson.M{db.KFieldUserUserName: details.Username}).One(&theUser)
			if err == nil {
				thePassword, err := crypto.GetHash(details.Password, theUser.Salt)
				if err == nil {
					if len(thePassword) == len(theUser.Password) {
						failed := false
						for i := 0; i < len(thePassword); i++ {
							if thePassword[i] != theUser.Password[i] {
								failed = true
							}
						}
						if theUser.InActive {
							failed = true
						}

						if !failed {
							for true {
								cookieString = crypto.RandomChars(16)
								err = sessionCollection.FindId(cookieString).One(&theSession)
								if err == mgo.ErrNotFound {
									break
								}
							}
							theCookie := http.Cookie{Name: "session", Value: cookieString}
							if details.Remember {
								theCookie.Expires = time.Now().Add(31 * 24 * time.Hour)
							}
							theSession.Id = cookieString
							theSession.UserNumber = theUser.Id
							theSession.Admin = theUser.Admin
							theSession.LoggedIn = true
							err = sessionCollection.Insert(&theSession)

							ctx.SetCookie(&theCookie)
							ctx.Redirect("/", http.StatusFound)
							cookie.MakeCookie(ctx)
							return
						}
					}
				}
			}
		} else {
			log.Println(err)
		}

	}
	details.Checkfield = crypto.MakeNonce(KLoginForm, "noone")
	ctx.ViewData("Title", "Hi Page")
	ctx.ViewData("Name", "iris")
	// ctx.ViewData("User", user)
	ctx.ViewData("Details", details)
	ctx.ViewData("Header", header)
	// ctx.ViewData("", myCcustomStruct{})
	ctx.View("login.html")
}
