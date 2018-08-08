package clients

import (
	"errors"
	// "fmt"
	"sort"
	"strconv"
	"strings"
	// "log"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/schema"
	"github.com/kataras/iris"
	// "github.com/globalsign/mgo"

	"ssafa/crypto"
	"ssafa/db"
)

type (
	NameRec struct {
		ID      int
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
	ErrorDateLowYear      = errors.New("The year is too low")
	ErrorDateHighDay      = errors.New("The day is too high")
	ErrorDateHighMonth    = errors.New("The month is too high")
	ErrorDateHighYear     = errors.New("The year is too high")
	errorRenfundsUsed     = errors.New("That RENfunds number is used elsewhere")
)

func SetRoutes(app *iris.Application) {
	app.Get("/client/{clientnum:int}", showClient)
	app.Get("/clients", ListClients)
	app.Any("/editclient/{clientnum:int}", editClientHandler)
	app.Get("/clients/{pagenum:int}", ListClients)
	app.Post("/searchclient", searchClients)
	app.Post("/searchclient/{pagenum:int}", searchClients)
	app.Post("/commentclient/{clientnum:int}", addComment)
	app.Get("/clientcomment/{clientnum:int}/{commentnum:int}", editComment)
	app.Post("/clientcomment/{clientnum:int}/{commentnum:int}", editComment)
	// app.Get("/client/{clientnum:int}", showClient)
	app.Get("/clientfix/{usernum:int}", clientFix)
	app.Post("/clientfix/{usernum:int}", clientFix)
}

func ReOrder() {
	// session := db.MongoSession.Copy()
	// defer session.Close()
	// clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	// num, err := clientColl.Find(nil).Count()

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
		orderMap[theClient.ID] = theClient.Order

		newName := NameRec{ID: theClient.ID}
		newName.First = strings.ToLower(theClient.First)
		newName.Surname = strings.ToLower(theClient.Surname)
		newName.ID = theClient.ID
		theList = append(theList, newName)
	}
	iter.Close()

	sort.Sort(ByClientName{theList})
	for i, item := range theList {
		offset, ok := orderMap[item.ID]
		if offset == i && ok {
			continue
		}
		set := bson.M{db.KFieldClientsOrder: i}
		clientColl.UpdateId(item.ID, bson.M{"$set": set})
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
