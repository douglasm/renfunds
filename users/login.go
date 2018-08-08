package users

import (
	// "fmt"
	"log"
	"net/http"
	"strings"
	"time"

	// "github.com/gorilla/schema"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"

	"ssafa/cookie"
	"ssafa/crypto"
	"ssafa/db"
	"ssafa/mail"
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
		// theSession Session
		details LoginRecord
		// cookieString string
		theUser db.User
		err     error
	)

	// theSession = ctx.Values().Get("session").(Session)
	header := types.HeaderRecord{Title: "Renfunds login"}
	header.Scripts = append(header.Scripts, "passwordtoggle")

	if ctx.Method() == http.MethodPost {
		data := ctx.FormValues()
		err = decoder.Decode(&details, data)
		if err == nil {
			session := db.MongoSession.Copy()
			defer session.Close()

			usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
			// sessionCollection := session.DB(db.MainDB).C(db.CollectionSessions)
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
							theCookie := createCookie(theUser.ID, theUser.Admin, details.Remember)
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
	details.Checkfield = crypto.MakeNonce(kLoginForm, "noone")
	ctx.ViewData("Details", details)
	ctx.ViewData("Header", header)
	ctx.View("login.html")
}

func logout(ctx iris.Context) {

	cookieName := ctx.GetCookie("session")
	ctx.RemoveCookie(cookieName)

	session := db.MongoSession.Copy()
	defer session.Close()

	sessionCollection := session.DB(db.MainDB).C(db.CollectionSessions)
	err := sessionCollection.RemoveId(cookieName)
	if err != nil {
		log.Println("Error: logout fail", cookieName, err)
	}

	ctx.Redirect("/", http.StatusFound)
}

func resetPassword(ctx iris.Context) {
	var (
		// theSession   Session
		details resetRec
		// cookieString string
		// err          error
	)

	if ctx.Method() == http.MethodPost {
		data := ctx.FormValues()
		err := decoder.Decode(&details, data)
		if err == nil {
			theEMail, theCode := getEMailAddress(details.Username)
			if len(theEMail) > 0 {
				mail.SendReset(theEMail, theCode)
			}

			ctx.Redirect("/resetsent", http.StatusFound)
			return
		}
	}

	header := types.HeaderRecord{Title: "Renfunds reset password"}

	details.Checkfield = crypto.MakeNonce(kLoginForm, "noone")
	ctx.ViewData("Details", details)
	ctx.ViewData("Header", header)
	ctx.View("resetemail.html")
}

func resetSent(ctx iris.Context) {
	header := types.HeaderRecord{Title: "Renfunds reset password e-mail sent"}

	ctx.ViewData("Header", header)
	ctx.View("resetsent.html")
}

func resetHandler(ctx iris.Context) {
	var (
		details      resetRec
		theUser      db.User
		errorMessage string
	)

	theCode := ctx.Params().Get("code")
	if len(theCode) > 0 {
		details.HasCode = true
	}

	switch ctx.Method() {
	case http.MethodGet:
		if details.HasCode {
			details.Code = theCode
		}

	case http.MethodPost:
		data := ctx.FormValues()
		err := decoder.Decode(&details, data)
		if err != nil {
			log.Println("Error: resetHandler decode", err)
		}
		if details.HasCode {
			details.Code = theCode
		}

		session := db.MongoSession.Copy()
		defer session.Close()

		usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

		err = usersCollection.Find(bson.M{db.KFieldUserActivateCode: details.Code}).One(&theUser)
		if err != nil {
			errorMessage = "The code and user name do not match"
			break
		}
		if len(theUser.Username) < 3 {
			errorMessage = "Please enter your user name"
			break
		}

		if theUser.Username != details.Username {
			errorMessage = "The code and user name do not match"
			break
		}

		if theUser.InActive {
			errorMessage = "The code and user name do not match"
			break
		}

		_, err = CheckPassword(details.Pass1, details.Pass2)
		if err == nil {
			setNewPassword(theUser.ID, details.Pass1, false)
			theCookie := createCookie(theUser.ID, theUser.Admin, false)
			ctx.SetCookie(&theCookie)
			ctx.Redirect("/", http.StatusFound)
			return
		}

		errorMessage = err.Error()

	}

	header := types.HeaderRecord{Title: "Renfunds reset password"}
	header.Scripts = append(header.Scripts, "passwordtoggle")

	details.Checkfield = crypto.MakeNonce(kResetForm, "noone")

	ctx.ViewData("Details", details)
	ctx.ViewData("Header", header)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("resetpass.html")
}

func getEMailAddress(name string) (string, string) {
	var (
		theUser db.User
		retVal  string
		theCode string
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	theCode = crypto.RandomLower(12)
	theUpdate := bson.M{}
	theUpdate[db.KFieldUserActivateTime] = time.Now().Unix() + 108000
	theUpdate[db.KFieldUserActivateCode] = theCode

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
	err := usersCollection.Find(bson.M{db.KFieldUserUserName: name}).One(&theUser)
	if err == nil {
		if theUser.InActive {
			return "", ""
		}

		err = usersCollection.UpdateId(theUser.ID, bson.M{"$set": theUpdate})
		return crypto.Decrypt(theUser.EMail), theCode
	}
	iter := usersCollection.Find(nil).Iter()

	name = strings.ToLower(name)
	for iter.Next(&theUser) {
		theStr := crypto.Decrypt(theUser.EMail)
		theStr = strings.ToLower(theStr)
		if theStr == name {
			retVal = name
			err = usersCollection.UpdateId(theUser.ID, bson.M{"$set": theUpdate})
			break
		}
	}
	iter.Close()

	return retVal, theCode
}

func createCookie(userNum int, admin bool, remember bool) http.Cookie {
	var (
		theSession   Session
		cookieString string
		err          error
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	sessionCollection := session.DB(db.MainDB).C(db.CollectionSessions)

	for true {
		cookieString = crypto.RandomChars(16)
		err = sessionCollection.FindId(cookieString).One(&theSession)
		if err == mgo.ErrNotFound {
			break
		}
	}
	theCookie := http.Cookie{Name: "session", Value: cookieString}
	if remember {
		theCookie.Expires = time.Now().Add(31 * 24 * time.Hour)
	}
	theSession.ID = cookieString
	theSession.UserNumber = userNum
	theSession.Admin = admin
	theSession.LoggedIn = true
	err = sessionCollection.Insert(&theSession)
	return theCookie
}

func verifyPassword(currentPassword string, theUser db.User) error {
	thePassword, err := crypto.GetHash(currentPassword, theUser.Salt)
	if err != nil {
		return err
	}
	if len(thePassword) == len(theUser.Password) {
		for i := 0; i < len(thePassword); i++ {
			if thePassword[i] != theUser.Password[i] {
				return errPassBad
			}
		}
		if theUser.InActive {
			return errPassBad
		}
	}
	return nil
}
