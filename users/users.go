package users

import (
	// "fmt"
	// "html"
	"log"
	"strings"
	// "github.com/gorilla/schema"
	"github.com/kataras/iris"

	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "ssafa/crypto"

	"ssafa/crypto"
	"ssafa/db"
	// "ssafa/types"
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
