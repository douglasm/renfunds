package users

import (
	// "fmt"
	// "strings"
	// "github.com/gorilla/schema"
	// "github.com/kataras/iris"

	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "ssafa/crypto"

	"ssafa/db"
	// "ssafa/types"
)

func (sr *Session) ValidCookie(value string) bool {
	session := db.MongoSession.Copy()
	defer session.Close()

	sessionCollection := session.DB(db.MainDB).C(db.CollectionSessions)
	err := sessionCollection.FindId(value).One(sr)
	if err != nil {
		return false
	}

	return true
}
