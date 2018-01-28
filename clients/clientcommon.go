package clients

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	// "log"

	"github.com/gorilla/schema"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2"

	"ssafa/crypto"
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
	decoder = schema.NewDecoder()
)

var (
	ErrorNoFirst          = errors.New("No first name")
	ErrorNoSurname        = errors.New("No surname")
	ErrorNINumWrongLength = errors.New("The NI Number should be 9 characters")
	ErrorNINumBadFormat   = errors.New("The NI Number should be 2 letters, 6 digits and 1 letter")
	ErrorDateBadFormat    = errors.New("The date format is wrong")
	ErrorDateBadDay       = errors.New("The day value is wrong")
	ErrorDateBadMonth     = errors.New("The month value is wrong")
	ErrorDateBadYear      = errors.New("The year value is wrong")
	ErrorDateLowDay       = errors.New("The day is too low")
	ErrorDateLowMonth     = errors.New("The month is too low")
	ErrorDateLowYear      = errors.New("The aeary is too low")
	ErrorDateHighDay      = errors.New("The day is too high")
	ErrorDateHighMonth    = errors.New("The month is too high")
	ErrorDateHighYear     = errors.New("The year is too high")
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
	app.Any("/editclient/{clientnum:int}", editClientHandler)
	app.Get("/clients/{pagenum:int}", ListClients)
	app.Post("/searchclient", searchClients)
	app.Post("/searchclient/{pagenum:int}", searchClients)
	app.Post("/commentclient/{clientnum:int}", addComment)
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

func decryptClient(theClient *db.Client) {
	theClient.First = crypto.Decrypt(theClient.First)
	theClient.Surname = crypto.Decrypt(theClient.Surname)
	theClient.Address = crypto.Decrypt(theClient.Address)
	theClient.PostCode = crypto.Decrypt(theClient.PostCode)
	theClient.Phone = crypto.Decrypt(theClient.Phone)
	theClient.Mobile = crypto.Decrypt(theClient.Mobile)
	theClient.EMail = crypto.Decrypt(theClient.EMail)
	theClient.NINum = crypto.Decrypt(theClient.NINum)
	theClient.ServiceNum = crypto.Decrypt(theClient.ServiceNum)
}

func parseDateString(theDate string) int {
	theDate = strings.TrimSpace(theDate)
	if len(theDate) == 0 {
		return 0
	}
	parts := strings.Split(theDate, "/")
	d, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	y, _ := strconv.Atoi(parts[2])
	return (y * 1000) + (m * 50) + d
}
