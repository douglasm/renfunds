package users

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
)

type (
	editRec struct {
		Id          int           `schema:"id"`
		First       string        `schema:"first"`
		Surname     string        `schema:"surname"`
		Name        string        `schema:"_"`
		Position    string        `schema:"posit"`
		EMail       string        `schema:"email"`
		Address     string        `schema:"address"`
		AddressHTML template.HTML `schema:"-"`
		PostCode    string        `schema:"postcode"`
		Phone       string        `schema:"phone"`
		Mobile      string        `schema:"mobile"`
		Admin       bool          `schema:"-"`
		AdminStr    string        `schema:"admin"`
		Commit      string        `schema:"commit"`
	}
)

func GetUserName(userNumber int) string {
	var (
		theUser db.User
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	err := usersCollection.FindId(userNumber).One(&theUser)
	if err != nil {
		log.Println("Error: Reading user name error:", err)
		return "Unknown personnel"
	}
	return crypto.Decrypt(theUser.First) + " " + crypto.Decrypt(theUser.Surname)
}

func getUserList(ctx iris.Context) {
	var (
		ur       userRequest
		userList []userChoice
		theUser  db.User
	)
	decoder.Decode(&ur, ctx.FormValues())

	searchTerm := strings.ToLower(ur.SearchTerm)

	// theCookie := ctx.Values().Get("cookie")

	session := db.MongoSession.Copy()
	defer session.Close()

	userCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	count := 0
	// iter := userCollection.Find(bson.M{db.KFieldUserInactive: false}).Iter()
	iter := userCollection.Find(nil).Iter()
	for iter.Next(&theUser) {
		if count > 10 {
			break
		}
		if theUser.InActive {
			continue
		}

		theName := crypto.Decrypt(theUser.First) + " " + crypto.Decrypt(theUser.Surname)
		// fmt.Println(theName, searchTerm)
		if strings.Contains(strings.ToLower(theName), searchTerm) {
			newUser := userChoice{Id: theUser.Id, Name: theName}
			userList = append(userList, newUser)
			count++
		}
	}
	iter.Close()
	// fmt.Println(userList)

	ctx.JSON(userList)
}

func editUser(ctx iris.Context) {
	var (
		theUser      db.User
		header       types.HeaderRecord
		details      editRec
		userNum      int
		errorMessage string
		err          error
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	if !theSession.(Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	userNum, err = ctx.Params().GetInt("usernum")
	if err != nil {
		ctx.Redirect("/adminperson", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	switch ctx.Method() {
	case http.MethodGet:
		err = usersCollection.FindId(userNum).One(&theUser)
		if err != nil {
			log.Println("Error: editUser find error", err)
			ctx.Redirect("/adminperson", http.StatusFound)
			return
		}

		details.Id = userNum
		details.First = crypto.Decrypt(theUser.First)
		details.Surname = crypto.Decrypt(theUser.Surname)
		details.Position = theUser.Position
		details.EMail = crypto.Decrypt(theUser.EMail)
		details.Address = crypto.Decrypt(theUser.Address)
		details.PostCode = crypto.Decrypt(theUser.PostCode)
		details.Phone = crypto.Decrypt(theUser.Phone)
		details.Mobile = crypto.Decrypt(theUser.Mobile)
		details.Admin = theUser.Admin

	case http.MethodPost:
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode user edit", err)
		}
		details.save()
	}

	details.Name = strings.TrimSpace(details.First + " " + details.Surname)

	header.Title = "Edit user " + details.Name
	header.Admin = theSession.(Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("useredit.html")
}

func showUser(ctx iris.Context) {
	var (
		theUser db.User
		header  types.HeaderRecord
		details editRec
		userNum int
		err     error
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	if !theSession.(Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	userNum, err = ctx.Params().GetInt("usernum")
	if err != nil {
		ctx.Redirect("/adminperson", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	err = usersCollection.FindId(userNum).One(&theUser)
	if err != nil {
		log.Println("Error: editUser find error", err)
		ctx.Redirect("/adminperson", http.StatusFound)
		return
	}

	details.Id = userNum
	details.First = crypto.Decrypt(theUser.First)
	details.Surname = crypto.Decrypt(theUser.Surname)
	details.Position = theUser.Position
	details.EMail = crypto.Decrypt(theUser.EMail)
	details.Address = crypto.Decrypt(theUser.Address)
	details.AddressHTML = template.HTML(strings.Replace(details.Address, "\r", "<br />", -1))
	details.PostCode = crypto.Decrypt(theUser.PostCode)
	details.Phone = crypto.Decrypt(theUser.Phone)
	details.Mobile = crypto.Decrypt(theUser.Mobile)
	details.Admin = theUser.Admin

	details.Name = strings.TrimSpace(details.First + " " + details.Surname)

	header.Title = "User " + details.Name
	header.Admin = theSession.(Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("usershow.html")
}

func (er editRec) save() {
	var (
		theUser  db.User
		tempBool bool
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
	err := usersCollection.FindId(er.Id).One(&theUser)

	sets := bson.M{}

	if crypto.Decrypt(theUser.First) != er.First {
		sets[db.KFieldUserFirst] = crypto.Encrypt(er.First)
	}

	if crypto.Decrypt(theUser.Surname) != er.Surname {
		sets[db.KFieldUserSurname] = crypto.Encrypt(er.Surname)
	}

	if crypto.Decrypt(theUser.Address) != er.Address {
		sets[db.KFieldUserAddress] = crypto.Encrypt(er.Address)
	}

	if crypto.Decrypt(theUser.PostCode) != er.PostCode {
		sets[db.KFieldUserPostCode] = crypto.Encrypt(er.PostCode)
	}

	if crypto.Decrypt(theUser.EMail) != er.EMail {
		sets[db.KFieldUserEMail] = crypto.Encrypt(er.EMail)
	}

	if crypto.Decrypt(theUser.Phone) != er.Phone {
		sets[db.KFieldUserPhone] = crypto.Encrypt(er.Phone)
	}

	if crypto.Decrypt(theUser.Mobile) != er.Mobile {
		sets[db.KFieldUserMobile] = crypto.Encrypt(er.Mobile)
	}

	if theUser.Position != er.Position {
		sets[db.KFieldUserPosition] = er.Position
	}

	if er.AdminStr == "yes" {
		tempBool = true
	}

	if theUser.Admin != tempBool {
		sets[db.KFieldUserAdmin] = tempBool
	}

	if len(sets) == 0 {
		return
	}

	err = usersCollection.UpdateId(er.Id, bson.M{"$set": sets})
	if err != nil {
		log.Println("Error: update user", err)
	}
	// KFieldUserBased = "based"
	// KFieldUserArea = "area"

}
