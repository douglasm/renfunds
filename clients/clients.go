package clients

import (
	// "log"
	"fmt"
	"html"
	"html/template"
	"net/http"
	// "sort"
	"strings"

	// "github.com/gorilla/schema"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"

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
		DOB     string
	}

	ClientList []ClientShow

	ClientDisplay struct {
		Id         int
		First      string
		Surname    string
		Phone      string
		Mobile     string
		EMail      string
		DOB        string
		NINum      string
		ServiceNum string
		Unit       string
		Address    template.HTML
		PostCode   string
		Comments   []CommentDisplay
		Reports    []CommentDisplay
		Cases      []cases.CaseList
	}

	ClientEdit struct {
		Id         int              `schema:"id"`
		First      string           `schema:"first"`
		Surname    string           `schema:"surname"`
		DOB        string           `schema:"dob"`
		NINum      string           `schema:"ninum"`
		ServiceNum string           `schema:"servicenum"`
		Unit       string           `schema:"unit"`
		Phone      string           `schema:"phone"`
		Mobile     string           `schema:"mobile"`
		EMail      string           `schema:"email"`
		Address    string           `schema:"address"`
		PostCode   string           `schema:"postcode"`
		Comments   []CommentDisplay `schema:"comments"`
		Checkfield string           `schema:"checkfield"`
		Commit     string           `schema:"commit"`
	}

	CommentDisplay struct {
		Date    string
		Comment template.HTML
		Name    template.HTML
	}
)

func ListClients(ctx iris.Context) {
	var (
		header     types.HeaderRecord
		pageNum    int
		nextPage   bool
		navButtons types.NavButtonRecord
		searchRec  types.SearchRecord
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
	ctx.ViewData("Search", searchRec)
	ctx.View("clients.html")
}

func editClientHandler(ctx iris.Context) {
	var (
		theClient    db.Client
		details      ClientEdit
		header       types.HeaderRecord
		errorMessage string
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	clientNum, err := ctx.Params().GetInt("clientnum")
	if err != nil {
		clientNum = 0
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	switch ctx.Method() {
	case http.MethodGet:
		if clientNum > 0 {
			clientColl := session.DB(db.MainDB).C(db.CollectionClients)
			err = clientColl.FindId(clientNum).One(&theClient)
			details.fillEdit(theClient)
		}

	case http.MethodPost:
		fmt.Println("Got edit post")
		fmt.Println(ctx.FormValues())
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			fmt.Println(err)
		}
		err = details.checkClient()
		if err == nil {
			fmt.Println("Good customer")
			if clientNum == 0 {
				clientNum = details.saveClient()
			} else {
				details.updateClient()
			}
			ctx.Redirect(fmt.Sprintf("/client/%d", clientNum), http.StatusFound)
			return
		} else {
			errorMessage = err.Error()
		}
	}

	header.Title = "RF: Edit client"
	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)

	ctx.View("clientedit.html")
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
	clientList, nextPage := GetList(theSearch.Type, theSearch.Term, pageNum)

	navData := types.M{}
	navData[types.KFieldNavPage] = pageNum
	navData[types.KFieldNavLink] = "clients"
	navData[types.KFieldNavNext] = nextPage

	navButtons.SetNavButtons(navData)

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", clientList)
	ctx.ViewData("NavButtons", navButtons)
	ctx.ViewData("Search", theSearch)
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
	if err != nil {
		fmt.Println("Client read error:", err)
	}
	// fmt.Println(theClient.Id)
	header.Title = "RF: Client"
	details.Id = theClient.Id
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
	details.Mobile = theClient.Mobile
	details.EMail = theClient.EMail
	details.NINum = theClient.NINum
	d, m, y := dateToDMY(theClient.DOB)
	if d != 0 && m != 0 {
		details.DOB = fmt.Sprintf("%02d/%02d/%04d", d, m, y)
	}
	details.ServiceNum = theClient.ServiceNum
	details.Unit = theClient.Unit
	details.Comments = getComments(theClient.Comments)
	details.Reports = getComments(theClient.Reports)

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
			d, m, y := dateToDMY(theClient.DOB)
			newClient.DOB = fmt.Sprintf("%02d/%02d/%04d", d, m, y)
			theList = append(theList, newClient)
		} else {
			return theList, true
		}
	}
	iter.Close()

	return theList, false
}

func (ce *ClientEdit) fillEdit(theClient db.Client) {
	ce.Id = theClient.Id
	ce.First = theClient.First
	ce.Surname = theClient.Surname
	ce.Phone = theClient.Phone
	ce.Mobile = theClient.Mobile
	ce.EMail = theClient.EMail
	ce.Address = theClient.Address
	ce.PostCode = theClient.PostCode
	ce.NINum = theClient.NINum
	ce.ServiceNum = theClient.ServiceNum
	ce.Unit = theClient.Unit
}

func getComments(comments []db.Comment) (retVal []CommentDisplay) {
	for _, item := range comments {
		newComment := CommentDisplay{}
		if item.Date != 0 {
			d := item.Date % 50
			m := (item.Date - d) % 1000
			m /= 50
			y := item.Date / 1000
			newComment.Date = fmt.Sprintf("%d/%02d/%04d ", d, m, y)
		}
		newComment.Comment = template.HTML(html.EscapeString(item.Comment))
		if len(item.Name) > 0 {
			newComment.Name = template.HTML(html.EscapeString(item.Name))
		}
		retVal = append(retVal, newComment)
	}
	return
}

func (ce *ClientEdit) checkClient() error {
	if len(ce.First) == 0 {
		return ErrorNoFirst
	}
	if len(ce.Surname) == 0 {
		return ErrorNoSurname
	}

	return nil
}

func (ce *ClientEdit) saveClient() int {
	var (
		theClient db.Client
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	theClient.Id = db.GetNextSequence(db.CollectionClients)
	theClient.First = ce.First
	theClient.Surname = ce.Surname
	theClient.Phone = ce.Phone
	theClient.Mobile = ce.Mobile
	theClient.EMail = ce.EMail
	theClient.Address = ce.Address
	theClient.PostCode = ce.PostCode
	theClient.NINum = ce.NINum
	theClient.ServiceNum = ce.ServiceNum
	theClient.Unit = ce.Unit

	err := clientColl.Insert(&theClient)
	if err != nil {
		println("Client save error:", err)
	}

	return theClient.Id
}

func (ce *ClientEdit) updateClient() {
	var (
		theClient db.Client
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	err := clientColl.FindId(ce.Id).One(&theClient)
	sets := bson.M{}
	unsets := bson.M{}

	if ce.First != theClient.First {
		sets[db.KFieldClientsFirst] = ce.First
	}
	if ce.Surname != theClient.Surname {
		sets[db.KFieldClientsSurname] = ce.Surname
	}
	if ce.Address != theClient.Address {
		if len(ce.Address) > 0 {
			sets[db.KFieldClientsAddress] = ce.Address
		} else {
			unsets[db.KFieldClientsAddress] = 1
		}
	}

	if ce.Phone != theClient.Phone {
		if len(ce.Phone) > 0 {
			sets[db.KFieldClientsPhone] = ce.Phone
		} else {
			unsets[db.KFieldClientsPhone] = 1
		}
	}

	if ce.Mobile != theClient.Mobile {
		if len(ce.Mobile) > 0 {
			sets[db.KFieldClientsMobile] = ce.Mobile
		} else {
			unsets[db.KFieldClientsMobile] = 1
		}
	}

	if ce.EMail != theClient.EMail {
		if len(ce.EMail) > 0 {
			sets[db.KFieldClientsEMail] = ce.EMail
		} else {
			unsets[db.KFieldClientsEMail] = 1
		}
	}

	if ce.NINum != theClient.NINum {
		if len(ce.NINum) > 0 {
			sets[db.KFieldClientsNINum] = ce.NINum
		} else {
			unsets[db.KFieldClientsNINum] = 1
		}
	}

	if ce.ServiceNum != theClient.ServiceNum {
		if len(ce.ServiceNum) > 0 {
			sets[db.KFieldClientsServiceNum] = ce.ServiceNum
		} else {
			unsets[db.KFieldClientsServiceNum] = 1
		}
	}

	if ce.Unit != theClient.Unit {
		if len(ce.Unit) > 0 {
			sets[db.KFieldClientsUnit] = ce.Unit
		} else {
			unsets[db.KFieldClientsUnit] = 1
		}
	}

	update := bson.M{}
	if len(sets) > 0 {
		if len(unsets) > 0 {
			update = bson.M{"$set": sets, "$unset": unsets}
		} else {
			update = bson.M{"$set": sets}
		}

	} else {
		if len(unsets) == 0 {
			return
		}
		update = bson.M{"$unset": unsets}
	}
	err = clientColl.UpdateId(ce.Id, update)
	if err != nil {
		fmt.Println("Update error:", err)
	}
}
