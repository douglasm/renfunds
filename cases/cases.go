package cases

import (
	// "log"
	"fmt"
	"net/http"
	// "sort"
	// "strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// "ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
)

type (
	CaseDisplay struct {
		ClientName string
		ClientNum  int
		Reports    []string
		Comments   []string
	}
)

func showCase(ctx iris.Context) {
	var (
		details   CaseDisplay
		theCase   db.Case
		theClient db.Client
		header    types.HeaderRecord
		caseNum   int
		err       error
	)

	caseNum, err = ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	clientColl := session.DB(db.MainDB).C(db.CollectionClients)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	err = caseColl.FindId(caseNum).One(&theCase)
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}
	err = clientColl.FindId(theCase.ClientNum).One(&theClient)

	header.Title = "RF: Case " + theCase.CaseNumber
	details.ClientName = theClient.First + " " + theClient.Surname + " " + theCase.CaseNumber
	// tempStr := html.EscapeString(theClient.Address)
	// tempStr = strings.TrimSpace(tempStr)
	// tempStr = strings.Replace(theClient.Address, "\r", "<br />", -1)
	// if len(tempStr) > 0 {
	// 	tempStr += "<br />"
	// }
	// if len(theClient.PostCode) > 0 {
	// 	tempStr += theClient.PostCode + "<br />"
	// }
	// details.Address = template.HTML(tempStr)
	// details.PostCode = theClient.PostCode
	// details.Phone = theClient.Phone
	details.ClientNum = theClient.Id
	for _, item := range theClient.Comments {
		theDate := ""
		if item.Date != 0 {
			d := item.Date % 50
			m := (item.Date - d) % 1000
			m /= 50
			y := item.Date / 1000
			theDate = fmt.Sprintf("%d/%02d/%04d ", d, m, y)
		}
		theDate += item.Comment
		if len(item.Name) > 0 {
			theDate += " - " + item.Name
		}
		details.Comments = append(details.Comments, theDate)
	}
	for _, item := range theClient.Reports {
		details.Reports = append(details.Reports, item.Comment)
	}

	// details.Cases = cases.GetCases(theClient.Id)

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("case.html")
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
		newCase := CaseList{CaseNumber: theCase.CaseNumber, CaseWorker: theCase.CaseFirst + " " + theCase.CaseSurname}
		newCase.Id = theCase.Id
		newCase.CMSNumber = theCase.CMSId
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
