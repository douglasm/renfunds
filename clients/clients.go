package clients

import (
	// "log"
	// "fmt"
	"html"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"

	"ssafa/cases"
	"ssafa/db"
	"ssafa/types"
)

type (
	ClientShow struct {
		Id      int
		Case    string
		First   string
		Surname string
		NiNum   string
	}

	ClientList []ClientShow

	ClientDisplay struct {
		Id       int
		First    string
		Surname  string
		Phone    string
		NiNum    string
		Address  template.HTML
		PostCode string
		Comments []string
		Reports  []string
		Cases    []cases.CaseList
	}

	ByClientName struct{ ClientList }
)

var (
	key []byte
)

func showClient(ctx iris.Context) {
	var (
		details   ClientDisplay
		theClient db.Client
		header    types.HeaderRecord
		// theUser      db.User
		clientNum int
		err       error
	)

	clientNum, err = ctx.Params().GetInt("clientnum")
	if err != nil {
		ctx.Redirect("/clients", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	err = clientColl.FindId(clientNum).One(&theClient)

	header.Title = "RF: Client"
	details.First = theClient.First
	details.Surname = theClient.Surname
	tempStr := html.EscapeString(theClient.Address)
	tempStr = strings.TrimSpace(tempStr)
	tempStr = strings.Replace(theClient.Address, "\r", "<br />", -1)
	if len(tempStr) > 0 {
		tempStr += "<br />"
	}
	if len(theClient.PostCode) > 0 {
		tempStr += theClient.PostCode + "<br />"
	}
	details.Address = template.HTML(tempStr)
	details.PostCode = theClient.PostCode
	details.Phone = theClient.Phone
	details.Comments = theClient.Comments
	details.Reports = theClient.Reports

	details.Cases = cases.GetCases(theClient.Id)

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("client.html")
}

func (s ClientList) Len() int {
	return len(s)
}

func (s ClientList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByClientName) Less(i, j int) bool {
	if strings.ToLower(s.ClientList[i].Surname) != strings.ToLower(s.ClientList[j].Surname) {
		return strings.ToLower(s.ClientList[i].Surname) < strings.ToLower(s.ClientList[j].Surname)
	}
	return strings.ToLower(s.ClientList[i].First) < strings.ToLower(s.ClientList[j].First)
}

func GetList(searchCategory, searchTerm string, offset int) []ClientShow {
	var (
		theClient db.Client
		theList   []ClientShow
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	iter := clientColl.Find(nil).Sort("_id").Iter()
	for iter.Next(&theClient) {
		newClient := ClientShow{Id: theClient.Id}
		newClient.First = theClient.First
		newClient.Surname = theClient.Surname
		newClient.NiNum = theClient.NINum
		theList = append(theList, newClient)
	}
	iter.Close()

	sort.Sort(ByClientName{theList})

	return theList
}

func SetKey(theKey []byte) {
	key = theKey
}
