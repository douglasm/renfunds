package cases

import (
	// "log"
	"fmt"
	// "sort"
	// "strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// "ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
)

func showCase(ctx iris.Context) {
	var (
		details CaseDisplay
		theCase db.Case
		header  types.HeaderRecord
		caseNum int
		err     error
	)

	caseNum, err = ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/clients", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	err = caseColl.FindId(caseNum).One(&theCase)

	header.Title = "RF: Case"
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

func GetCases(clientNumber int) []CaseList {
	var (
		theCase  db.Case
		allCases []CaseList
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	iter := caseColl.Find(bson.M{db.KFieldCaseClientNum: clientNumber}).Sort("_id").Iter()
	for iter.Next(&theCase) {
		fmt.Printf("%+v\n", theCase)
		newCase := CaseList{CaseNumber: theCase.CaseNumber, CaseWorker: theCase.CaseFirst + " " + theCase.CaseSurname}
		newCase.Id = theCase.Id
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
