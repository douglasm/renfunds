package clients

import (
	"sort"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"ssafa/crypto"
	"ssafa/db"
)

type (
	sliceRec struct {
		Id      int
		First   string
		Surname string
	}

	ByName []sliceRec
)

func OrderClients() {
	var (
		allPeople []sliceRec
		theClient db.Client
	)
	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	iter := clientColl.Find(nil).Iter()
	for iter.Next(&theClient) {
		newRec := sliceRec{}
		newRec.Id = theClient.Id
		newRec.First = strings.ToLower(crypto.Decrypt(theClient.First))
		newRec.Surname = strings.ToLower(crypto.Decrypt(theClient.Surname))
		allPeople = append(allPeople, newRec)
	}
	iter.Close()

	sort.Sort(ByName(allPeople))

	posn := 10
	for _, item := range allPeople {
		clientColl.UpdateId(item.Id, bson.M{"$set": bson.M{db.KFieldClientsOrder: posn}})
		posn += 10
	}
}

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	if a[i].Surname < a[j].Surname {
		return true
	}
	if a[i].Surname > a[j].Surname {
		return false
	}
	return a[i].First < a[j].First
}
