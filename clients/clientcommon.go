package clients

import (
	"fmt"
	"sort"
	"strings"
	// "log"

	"github.com/gorilla/schema"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2"
	// "ssafa/crypto"

	"ssafa/db"
)

type (
	NameRec struct {
		Id      int
		First   string
		Surname string
	}

	NameList     []NameRec
	ByClientName struct{ NameList }
)

var (
	key     []byte
	decoder = schema.NewDecoder()
)

// func init() {
// 	dayNum := 2001*1000 + 1*50 + 16

// 	d, m, y := dateToDMY(dayNum)
// 	if d != 16 {
// 		fmt.Println("day bad", d)
// 	}
// 	if m != 1 {
// 		fmt.Println("day bad", m)
// 	}
// 	if y != 2001 {
// 		fmt.Println("year bad", y)
// 	}
// }

func SetRoutes(app *iris.Application) {
	app.Get("/client/{clientnum:int}", showClient)
	app.Get("/clients", ListClients)
	app.Get("/clients/{pagenum:int}", ListClients)
	app.Post("/searchclient", searchClients)
	app.Post("/searchclient/{pagenum:int}", searchClients)
	// app.Get("/client/{clientnum:int}", showClient)
}

func ReOrder() {
	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	num, err := clientColl.Find(nil).Count()

	fmt.Println("There are", num, err)
}

func orderNames() {
	var (
		theClient db.Client
		theList   NameList
		orderMap  = map[int]int{}
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	iter := clientColl.Find(nil).Sort(db.KFieldClientsOrder).Iter()
	for iter.Next(&theClient) {
		orderMap[theClient.Id] = theClient.Order

		newName := NameRec{Id: theClient.Id}
		newName.First = strings.ToLower(theClient.First)
		newName.Surname = strings.ToLower(theClient.Surname)
		newName.Id = theClient.Id
		theList = append(theList, newName)
	}
	iter.Close()

	sort.Sort(ByClientName{theList})
	for i, item := range theList {
		offset, ok := orderMap[item.Id]
		if offset == i && ok {
			continue
		}
		set := bson.M{db.KFieldClientsOrder: i}
		clientColl.UpdateId(item.Id, bson.M{"$set": set})
		// fmt.Println(num, item)
	}
}

func (s NameList) Len() int {
	return len(s)
}

func (s NameList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByClientName) Less(i, j int) bool {
	if s.NameList[i].Surname != s.NameList[j].Surname {
		return s.NameList[i].Surname < s.NameList[j].Surname
	}
	return s.NameList[i].First < s.NameList[j].First
}

func SetKey(theKey []byte) {
	key = theKey
}

func dateToDMY(theDate int) (day, month, year int) {
	day = theDate % 50
	theDate -= day
	month = theDate % 1000
	theDate -= month
	month /= 50

	year = theDate / 1000
	if year < 100 {
		year += 2000
	}
	return
}
