package clients

import (
	// "log"
	// "fmt"
	"html"
	"html/template"
	"net/http"
	// "sort"
	"strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"

	"ssafa/cases"
	"ssafa/db"
	"ssafa/types"
	"ssafa/users"
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
)

func listClients(ctx iris.Context) {
	var (
		header     types.HeaderRecord
		pageNum    int
		nextPage   bool
		navButtons types.NavButtonRecord
		err        error
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	// fmt.Println("We are logged in")
	header.Title = "RF: Clients"
	clientList, nextPage := GetList("", "", pageNum)

	navData := types.M{}
	navData[types.KFieldNavPage] = pageNum
	navData[types.KFieldNavLink] = "clients"
	navData[types.KFieldNavNext] = nextPage

	navButtons.SetNavButtons(navData)

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", clientList)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("clients.html")
}

func searchClients(ctx iris.Context) {
	var (
		theSearch  types.SearchRecord
		header     types.HeaderRecord
		pageNum    int
		nextPage   bool
		navButtons types.NavButtonRecord
		err        error
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	data := ctx.FormValues()
	err = decoder.Decode(&theSearch, data)

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	// fmt.Println("We are logged in")
	header.Title = "RF: Clients"
	clientList, nextPage := GetList(theSearch.SearchType, theSearch.SearchTerm, pageNum)

	navData := types.M{}
	navData[types.KFieldNavPage] = pageNum
	navData[types.KFieldNavLink] = "clients"
	navData[types.KFieldNavNext] = nextPage

	navButtons.SetNavButtons(navData)

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", clientList)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("clients.html")
}

func displayClients(header types.HeaderRecord, clientList []ClientShow, navButtons types.NavButtonRecord, ctx iris.Context) {
	ctx.ViewData("Header", header)
	ctx.ViewData("Details", clientList)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("clients.html")
}

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

func GetList(searchCategory, searchTerm string, pageNum int) ([]ClientShow, bool) {
	var (
		theClient db.Client
		theList   []ClientShow
		skip      int
		limit     int
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	if searchTerm == "" {
		skip = pageNum * types.KListLimit
		limit = types.KListLimit + 1
	} else {
		limit = 100000
	}
	found := 0
	iter := clientColl.Find(nil).Skip(skip).Limit(limit).Sort(db.KFieldClientsOrder).Iter()
	for iter.Next(&theClient) {
		if searchTerm != "" {
			isValid := false
			if strings.Contains(strings.ToLower(theClient.First), searchTerm) {
				isValid = true
			}
			if strings.Contains(strings.ToLower(theClient.Surname), searchTerm) {
				isValid = true
			}
			if !isValid {
				continue
			}
			if found < pageNum*types.KListLimit {
				found++
				continue
			}
		}
		if len(theList) < types.KListLimit {
			newClient := ClientShow{Id: theClient.Id}
			newClient.First = theClient.First
			newClient.Surname = theClient.Surname
			newClient.NiNum = theClient.NINum
			theList = append(theList, newClient)
		} else {
			return theList, true
		}
	}
	iter.Close()

	return theList, false
}
