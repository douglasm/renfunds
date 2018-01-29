package clients

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	// "github.com/gorilla/schema"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"

	"ssafa/cases"
	"ssafa/crypto"
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
		Comments   []cases.CommentDisplay
		Cases      []cases.CaseList
	}

	ClientEdit struct {
		Id         int                    `schema:"id"`
		First      string                 `schema:"first"`
		Surname    string                 `schema:"surname"`
		DOB        string                 `schema:"dob"`
		NINum      string                 `schema:"ninum"`
		ServiceNum string                 `schema:"servicenum"`
		Unit       string                 `schema:"unit"`
		Phone      string                 `schema:"phone"`
		Mobile     string                 `schema:"mobile"`
		EMail      string                 `schema:"email"`
		Address    string                 `schema:"address"`
		PostCode   string                 `schema:"postcode"`
		Comments   []cases.CommentDisplay `schema:"comments"`
		Checkfield string                 `schema:"checkfield"`
		Commit     string                 `schema:"commit"`
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
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println(err)
		}
		err = details.checkClient()
		if err == nil {
			if clientNum == 0 {
				clientNum = details.saveClient()
				OrderClients()
			} else {
				if details.updateClient() {
					OrderClients()
				}
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
	theSearch.Term = strings.ToLower(theSearch.Term)

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
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
		log.Println("Error: Client read error:", err)
	}
	header.Title = "RF: Client"

	decryptClient(&theClient)
	details.Id = theClient.Id
	details.First = theClient.First
	details.Surname = theClient.Surname
	tempStr := html.EscapeString(theClient.Address)
	tempStr = strings.TrimSpace(tempStr)
	tempStr = strings.Replace(tempStr, "\r", "<br />", -1)
	if len(tempStr) > 0 {
		tempStr += "<br />"
	}
	if len(theClient.PostCode) > 0 {
		tempStr += theClient.PostCode + "<br />"
	}
	details.Address = template.HTML(tempStr)
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

	allCases := getCases(theClient.Id)

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("Cases", allCases)
	ctx.View("client.html")
}

func addComment(ctx iris.Context) {
	var (
		newComment cases.CommentRec
		clientNum  int
		err        error
	)

	clientNum, err = ctx.Params().GetInt("clientnum")
	if err != nil {
		ctx.Redirect("/clients", http.StatusFound)
		return
	}

	decoder.Decode(&newComment, ctx.FormValues())

	if len(newComment.Comment) < 1 {
		ctx.Redirect(fmt.Sprintf("/client/%d", clientNum), http.StatusFound)
		return
	}

	theSession := ctx.Values().Get("session")

	tempStr := crypto.Encrypt(newComment.Comment)

	theTime := time.Now()
	y := theTime.Year() * 1000
	y += (int(theTime.Month()) * 50)
	y += theTime.Day()

	comment := bson.M{"comment": tempStr, "user": theSession.(users.Session).UserNumber, "date": y}
	each := []bson.M{comment}

	commentEach := bson.M{"$each": each, "$position": 0}

	thePush := bson.M{"$push": bson.M{"comments": commentEach}}
	// db.clients.update({"_id": 1375}, { $push: { comments: { $each: [ {date: 2018067, name: "Fred Smith"} ], $position: 0 } }})

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	err = clientColl.UpdateId(clientNum, thePush)
	if err != nil {
		log.Println("Error: Client addComment Push error:", err)
	}

	ctx.Redirect(fmt.Sprintf("/client/%d", clientNum), http.StatusFound)
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
		decryptClient(&theClient)
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
	decryptClient(&theClient)
	ce.Id = theClient.Id
	ce.First = theClient.First
	ce.Surname = theClient.Surname
	d, m, y := dateToDMY(theClient.DOB)
	if d != 0 && m != 0 {
		ce.DOB = fmt.Sprintf("%02d/%02d/%04d", d, m, y)
	}
	ce.Phone = theClient.Phone
	ce.Mobile = theClient.Mobile
	ce.EMail = theClient.EMail
	ce.Address = theClient.Address
	ce.PostCode = theClient.PostCode
	ce.NINum = theClient.NINum
	ce.ServiceNum = theClient.ServiceNum
	ce.Unit = theClient.Unit
}

func getComments(comments []db.Comment) (retVal []cases.CommentDisplay) {
	for _, item := range comments {
		newComment := cases.CommentDisplay{}
		newComment.GetCommentDisplay(item)
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

	if len(ce.NINum) > 0 {
		err := checkNiNum(ce.NINum)
		if err != nil {
			return err
		}
	}

	if len(ce.DOB) > 0 {
		parts := strings.Split(ce.DOB, "/")
		if len(parts) != 3 {
			return ErrorDateBadFormat
		}

		d, err := strconv.Atoi(parts[0])
		if err != nil {
			return ErrorDateBadDay
		}
		if d < 1 {
			return ErrorDateLowDay
		}

		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return ErrorDateBadMonth
		}
		if m < 1 {
			return ErrorDateLowMonth
		}
		if m > 12 {
			return ErrorDateHighMonth
		}

		y, err := strconv.Atoi(parts[2])
		if err != nil {
			return ErrorDateBadYear
		}

		theTime := time.Now()
		if y < theTime.Year()-112 {
			return ErrorDateLowYear
		}

		if y > theTime.Year()-18 {
			return ErrorDateHighYear
		}

		if d > 31 {
			return ErrorDateHighDay
		}

		// Check for the number of days in a month
		switch m {
		case 2:
			if y%4 != 0 {
				if d > 28 {
					return ErrorDateHighDay
				}
			} else {
				if d > 29 {
					return ErrorDateHighDay
				}
			}
		case 4, 6, 9, 11:
			if d > 30 {
				return ErrorDateHighDay
			}
		}
	}

	ce.PostCode = strings.ToUpper(ce.PostCode)
	ce.NINum = strings.ToUpper(ce.NINum)

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
	theClient.First = crypto.Encrypt(ce.First)
	theClient.Surname = crypto.Encrypt(ce.Surname)
	theClient.Phone = crypto.Encrypt(ce.Phone)
	theClient.Mobile = crypto.Encrypt(ce.Mobile)
	theClient.EMail = crypto.Encrypt(ce.EMail)
	theClient.Address = crypto.Encrypt(ce.Address)
	theClient.PostCode = crypto.Encrypt(ce.PostCode)
	theClient.NINum = crypto.Encrypt(ce.NINum)
	theClient.ServiceNum = crypto.Encrypt(ce.ServiceNum)
	theClient.Unit = ce.Unit

	theClient.DOB = parseDateString(ce.DOB)

	theClient.Created = db.GetCurrentDate()
	theClient.Changed = theClient.Created

	err := clientColl.Insert(theClient)
	if err != nil {
		log.Println("Error: Client save error:", err)
	}

	return theClient.Id
}

func (ce *ClientEdit) updateClient() bool {
	var (
		theClient db.Client
		retVal    bool
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	err := clientColl.FindId(ce.Id).One(&theClient)
	decryptClient(&theClient)
	sets := bson.M{}
	unsets := bson.M{}

	if ce.First != theClient.First {
		sets[db.KFieldClientsFirst] = crypto.Encrypt(ce.First)
		retVal = true
	}
	if ce.Surname != theClient.Surname {
		sets[db.KFieldClientsSurname] = crypto.Encrypt(ce.Surname)
		retVal = true
	}

	if ce.Address != theClient.Address {
		if len(ce.Address) > 0 {
			sets[db.KFieldClientsAddress] = crypto.Encrypt(ce.Address)
		} else {
			unsets[db.KFieldClientsAddress] = 1
		}
	}

	if ce.PostCode != theClient.PostCode {
		if len(ce.PostCode) > 0 {
			sets[db.KFieldClientsPostCode] = crypto.Encrypt(ce.PostCode)
		} else {
			unsets[db.KFieldClientsPostCode] = 1
		}
	}

	if ce.Phone != theClient.Phone {
		if len(ce.Phone) > 0 {
			sets[db.KFieldClientsPhone] = crypto.Encrypt(ce.Phone)
		} else {
			unsets[db.KFieldClientsPhone] = 1
		}
	}

	if ce.Mobile != theClient.Mobile {
		if len(ce.Mobile) > 0 {
			sets[db.KFieldClientsMobile] = crypto.Encrypt(ce.Mobile)
		} else {
			unsets[db.KFieldClientsMobile] = 1
		}
	}

	if ce.EMail != theClient.EMail {
		if len(ce.EMail) > 0 {
			sets[db.KFieldClientsEMail] = crypto.Encrypt(ce.EMail)
		} else {
			unsets[db.KFieldClientsEMail] = 1
		}
	}

	if ce.NINum != theClient.NINum {
		if len(ce.NINum) > 0 {
			sets[db.KFieldClientsNINum] = crypto.Encrypt(ce.NINum)
		} else {
			unsets[db.KFieldClientsNINum] = 1
		}
	}

	if ce.ServiceNum != theClient.ServiceNum {
		if len(ce.ServiceNum) > 0 {
			sets[db.KFieldClientsServiceNum] = crypto.Encrypt(ce.ServiceNum)
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

	val := parseDateString(ce.DOB)
	if val != theClient.DOB {
		if val != 0 {
			sets[db.KFieldClientsDOB] = val
		} else {
			unsets[db.KFieldClientsDOB] = 1
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
			return false
		}
		update = bson.M{"$unset": unsets}
	}

	theClient.Changed = db.GetCurrentDate()

	err = clientColl.UpdateId(ce.Id, update)
	if err != nil {
		log.Println("Update error:", err)
	}
	return retVal
}

func checkNiNum(nINum string) error {
	if len(nINum) != 9 {
		return ErrorNINumWrongLength
	}
	nINum = strings.ToUpper(nINum)
	for i, ch := range nINum {
		switch i {
		case 0, 1, 8:
			if ch < 'A' {
				return ErrorNINumBadFormat
			}
			if ch > 'Z' {
				return ErrorNINumBadFormat
			}

		case 2, 3, 4, 5, 6, 7:
			if ch < '0' {
				return ErrorNINumBadFormat
			}
			if ch > '9' {
				return ErrorNINumBadFormat
			}
		}
	}
	return nil
}

func getCases(clientNum int) []cases.CaseList {
	var (
		allCases []cases.CaseList
		theCase  db.Case
		theUser  db.User
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	iter := caseColl.Find(bson.M{db.KFieldCaseClientNum: clientNum}).Sort(db.KFieldCaseClosed, db.KFieldCaseUpdated).Iter()
	for iter.Next(&theCase) {
		newCase := cases.CaseList{Id: theCase.Id}
		if len(theCase.CaseNumber) > 0 {
			newCase.CaseNumber = theCase.CaseNumber
		} else {
			newCase.CaseNumber = "None"
		}
		newCase.CMSNumber = theCase.CMSId
		if theCase.CaseWorkerNum == 0 {
			newCase.CaseWorker = crypto.Decrypt(theCase.CaseWorker)
		} else {
			err := userColl.FindId(theCase.CaseWorkerNum).One(&theUser)
			if err == nil {
				newCase.CaseWorker = crypto.Decrypt(theUser.Name)
			}
		}
		if theCase.Closed {
			newCase.State = "Closed"
		} else {
			newCase.State = "Open"
		}

		allCases = append(allCases, newCase)
	}
	iter.Close()
	return allCases
}

// db.clients.update({"_id": 1375}, { $push: { comments: { $each: [ {date: 2018067, name: "Fred Smith"} ], $position: 0 } }})
// db.clients.update({"_id": 1375}, { $pop: { comments: -1 }})
