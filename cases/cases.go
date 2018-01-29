package cases

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "sort"
	// "strings"

	"github.com/kataras/iris"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
	"ssafa/users"
	"ssafa/utils"
)

type (
	CaseDisplay struct {
		Id         int
		Open       bool
		ClientName types.RowItem
		CaseNumber string
		CaseWorker types.RowItem
		CMSNumber  types.RowItem
		Opened     types.RowItem
		State      types.RowItem
		Updated    types.RowItem
		ClientNum  int
		Reports    []CommentDisplay
	}

	CaseEdit struct {
		Id             int           `schema:"id"`
		CaseNumber     string        `schema:"casenumber"`
		CaseWorker     template.HTML `schema:"-"`
		CaseWorkerName string        `schema:"cwname"`
		CaseWorkerNum  int           `schema:"cwnum"`
		CMSNumber      string        `schema:"cms"`
		ClientNum      int           `schema:"client"`
		ClientName     string        `schema:"clientname"`
		Checkfield     int           `schema:"checkfield"`
		Commit         string        `schema:"commit"`
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

	theSession := ctx.Values().Get("session")

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

	details.Id = caseNum
	header.Admin = theSession.(users.Session).Admin

	header.Title = "RF: Case " + theCase.CaseNumber
	details.CaseNumber = theCase.CaseNumber

	details.ClientName.Title = "Client"
	tempStr := crypto.Decrypt(theClient.First) + " " + crypto.Decrypt(theClient.Surname)
	details.ClientName.Value = template.HTML(tempStr)
	details.ClientName.Link = fmt.Sprintf("/client/%d", theCase.ClientNum)

	details.CaseWorker.Title = "Case worker"
	if theCase.CaseWorkerNum == 0 {
		details.CaseWorker.Value = template.HTML(crypto.Decrypt(theCase.CaseWorker))
	} else {
		details.CaseWorker.Value = template.HTML(users.GetUserName(theCase.CaseWorkerNum))
	}

	details.CMSNumber.Title = "CMS number"
	details.CMSNumber.Value = template.HTML(theCase.CMSId)

	details.Opened.Title = "Opened"
	details.Opened.Value = template.HTML(utils.DateToString(theCase.Created))

	details.Updated.Title = "Updated"
	details.Updated.Value = template.HTML(utils.DateToString(theCase.Updated))

	details.State.Title = "Closed"
	if theCase.Closed {
		details.State.Value = template.HTML("Yes " + utils.DateToString(theCase.DateClosed))
	} else {
		details.Open = true
		details.State.Value = template.HTML("No")
	}

	details.ClientNum = theClient.Id
	for _, item := range theCase.Comments {
		newComment := CommentDisplay{}
		newComment.GetCommentDisplay(item)
		details.Reports = append(details.Reports, newComment)
	}

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("case.html")
}

func listCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header  types.HeaderRecord
		pageNum int
		err     error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	match := bson.M{"$match": bson.M{db.KFieldCaseClientNum: bson.M{"$ne": 0}}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, _ = GetCases(pageNum, match, sort)

	header.Title = "RF: Unassigned Cases"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("cases.html")
}

func openCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header  types.HeaderRecord
		pageNum int
		err     error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	match := bson.M{"$match": bson.M{db.KFieldCaseClosed: false}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, _ = GetCases(pageNum, match, sort)

	header.Title = "RF: Unassigned Cases"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("casesopen.html")
}

func unassignedCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header  types.HeaderRecord
		pageNum int
		err     error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	match := bson.M{"$match": bson.M{db.KFieldCaseWorkerNum: 0, db.KFieldCaseClosed: false}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, _ = GetCases(pageNum, match, sort)

	header.Title = "RF: Unassigned Cases"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("casesunassign.html")
}

func inactiveCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header  types.HeaderRecord
		pageNum int
		err     error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	match := bson.M{"$match": bson.M{db.KFieldCaseWorkerNum: 0, db.KFieldCaseClosed: false}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseUpdated: 1}}

	details, _ = GetCases(pageNum, match, sort)

	header.Title = "RF: Unassigned Cases"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("casesinactive.html")
}

func closeCase(ctx iris.Context) {
	caseNum, err := ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	theUpdate := bson.M{db.KFieldCaseClosed: true}
	updateCase(theUpdate, caseNum)

	theUrl := fmt.Sprintf("/case/%d", caseNum)
	ctx.Redirect(theUrl, http.StatusFound)
}

func openCase(ctx iris.Context) {
	caseNum, err := ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	theUpdate := bson.M{db.KFieldCaseClosed: false}
	updateCase(theUpdate, caseNum)

	theUrl := fmt.Sprintf("/case/%d", caseNum)
	ctx.Redirect(theUrl, http.StatusFound)
}

func addCase(ctx iris.Context) {
	var (
		theCase db.Case
	)

	clientNum, err := ctx.Params().GetInt("clientnum")
	if err != nil {
		ctx.Redirect("/clients", http.StatusFound)
		return
	}

	// theSession := ctx.Values().Get("session")

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	theCase.Id = db.GetNextSequence(db.CollectionCases)
	theCase.ClientNum = clientNum
	theCase.Created = utils.CurrentDate()
	theCase.Updated = theCase.Created

	caseColl.Insert(&theCase)

	theUrl := fmt.Sprintf("/case/%d", theCase.Id)
	ctx.Redirect(theUrl, http.StatusFound)
}

func editCase(ctx iris.Context) {
	var (
		theCase      db.Case
		details      CaseEdit
		header       types.HeaderRecord
		errorMessage string
	)

	caseNum, err := ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	theSession := ctx.Values().Get("session")

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	switch ctx.Method() {
	case http.MethodGet:
		err = caseColl.FindId(caseNum).One(&theCase)
		if err != nil {
			log.Println("Error: editcase get", err)
		}
		details.Id = caseNum
		details.CaseNumber = theCase.CaseNumber
		details.CMSNumber = theCase.CMSId
		details.ClientNum = theCase.ClientNum
		details.ClientName = getClientName(theCase.ClientNum)
		if theCase.CaseWorkerNum == 0 {
			details.CaseWorker = template.HTML(crypto.Decrypt(theCase.CaseWorker))
		} else {
			details.CaseWorker = template.HTML(users.GetUserName(theCase.CaseWorkerNum))
		}

	case http.MethodPost:
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode editcase", err)
		}
		details.CaseWorker = template.HTML(details.CaseWorkerName)
		err = details.Save()
		if err == nil {
			theUrl := fmt.Sprintf("/case/%d", details.Id)
			ctx.Redirect(theUrl, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Admin = theSession.(users.Session).Admin
	header.Title = "Edit case"
	header.Scripts = append(header.Scripts, "getusers")

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("caseedit.html")
}

func updateCase(theUpdate bson.M, caseId int) {
	theUpdate[db.KFieldCaseUpdated] = utils.CurrentDate()

	session := db.MongoSession.Copy()
	defer session.Close()

	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	err := caseColl.UpdateId(caseId, bson.M{"$set": theUpdate})
	if err != nil {
		log.Println("Error: Case update:", err)
	}
}

func addComment(ctx iris.Context) {
	var (
		newComment CommentRec
		caseNum    int
		err        error
	)

	caseNum, err = ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	decoder.Decode(&newComment, ctx.FormValues())

	if len(newComment.Comment) < 1 {
		ctx.Redirect(fmt.Sprintf("/case/%d", caseNum), http.StatusFound)
		return
	}

	theSession := ctx.Values().Get("session")

	tempStr := crypto.Encrypt(newComment.Comment)

	y := utils.CurrentDate()

	theSet := bson.M{db.KFieldCaseUpdated: y}
	comment := bson.M{"comment": tempStr, "user": theSession.(users.Session).UserNumber, "date": y}
	each := []bson.M{comment}

	commentEach := bson.M{"$each": each, "$position": 0}

	thePush := bson.M{"$push": bson.M{"comments": commentEach}}
	// db.clients.update({"_id": 1375}, { $push: { comments: { $each: [ {date: 2018067, name: "Fred Smith"} ], $position: 0 } }})

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	err = caseColl.UpdateId(caseNum, thePush)
	if err != nil {
		log.Println("Error: Case update push:", err)
	}
	err = caseColl.UpdateId(caseNum, bson.M{"$set": theSet})
	if err != nil {
		log.Println("Error: Case update set:", err)
	}

	ctx.Redirect(fmt.Sprintf("/case/%d", caseNum), http.StatusFound)
}

func GetCases(pageNum int, match, sort bson.M) ([]CaseList, bool) {
	var (
		// theCase  db.Case
		allCases []CaseList
		result   CaseRec
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	theSkip := pageNum * types.KListLimit

	// match := bson.M{"$match": bson.M{db.KFieldCaseWorkerNum: 0, db.KFieldCaseClosed: false}}
	projectFields := bson.M{}
	projectFields[db.KFieldCaseClientNum] = 1
	projectFields[db.KFieldCaseNum] = 1
	projectFields[db.KFieldCaseCMS] = 1
	projectFields[db.KFieldCaseCreated] = 1
	projectFields[db.KFieldCaseClosed] = 1
	projectFields[db.KFieldCaseUpdated] = 1
	projectFields[db.KFieldCaseWorkerNum] = 1
	projectFields[db.KFieldCaseFirst] = 1
	projectFields[db.KFieldCaseSurname] = 1

	project := bson.M{"$project": projectFields}
	skip := bson.M{"$skip": theSkip}
	limit := bson.M{"$limit": types.KListLimit + 1}
	lookupClient := bson.M{"$lookup": bson.M{"from": db.CollectionClients, "localField": db.KFieldCaseClientNum, "foreignField": "_id", "as": "client"}}
	lookupUser := bson.M{"$lookup": bson.M{"from": db.CollectionUsers, "localField": db.KFieldCaseWorkerNum, "foreignField": "_id", "as": "user"}}

	// pipeline := []bson.M{match, sort, skip, limit, lookup}
	// pipeline := []bson.M{match, sort, project, skip, limit}
	pipeline := []bson.M{match, sort, project, skip, limit, lookupClient, lookupUser}
	// fmt.Println(project)
	// fmt.Println(lookup)

	// iter := caseColl.Find(nil).Skip(theSkip).Limit(types.KListLimit + 1).Sort(db.KFieldCaseUpdated).Iter()
	iter := caseColl.Pipe(pipeline).Iter()
	if iter.Err() != nil {
		log.Println("Error: GetCases", iter.Err())
	}
	count := 0
	for iter.Next(&result) {
		count++
		if count > types.KListLimit {
			return allCases, true
		}
		// fmt.Println(result)
		newCase := CaseList{Id: result.Id, CaseNumber: result.CaseNumber, CMSNumber: result.CMSId}
		if len(newCase.CaseNumber) == 0 {
			newCase.CaseNumber = "None"
		}
		newCase.ClientId = result.Client[0].Id
		newCase.Name = crypto.Decrypt(result.Client[0].First) + " " + crypto.Decrypt(result.Client[0].Surname)

		if len(result.User) > 0 {
			newCase.CaseWorker = crypto.Decrypt(result.User[0].First) + " " + crypto.Decrypt(result.User[0].Surname)
		} else {
			newCase.CaseWorker = crypto.Decrypt(result.CaseFirst) + " " + crypto.Decrypt(result.CaseSurname)
		}

		if result.Closed {
			newCase.State = "Closed"
		} else {
			newCase.State = "Open"
		}
		if result.Created != 0 {
			newCase.Opened = utils.DateToString(result.Created)
		}

		if result.Updated != 0 {
			newCase.Updated = utils.DateToString(result.Updated)
		}

		allCases = append(allCases, newCase)
	}
	iter.Close()
	return allCases, false
}

func (ce *CaseEdit) Save() error {
	var (
		theCase db.Case
	)
	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)
	// KFieldCaseNum        = "case"
	// KFieldCaseCMS        = "cms"

	gotOne := false
	iter := caseColl.Find(bson.M{db.KFieldCaseNum: ce.CaseNumber}).Iter()
	for iter.Next(&theCase) {
		if theCase.Id != ce.Id {
			gotOne = true
			break
		}
	}
	iter.Close()

	if gotOne {
		return ErrorDateCaseUsed
	}

	iter = caseColl.Find(bson.M{db.KFieldCaseCMS: ce.CMSNumber}).Iter()
	for iter.Next(&theCase) {
		if theCase.Id != ce.Id {
			gotOne = true
			break
		}
	}
	iter.Close()

	if gotOne {
		return ErrorDateCMSUsed
	}

	sets := bson.M{}
	if ce.CMSNumber != theCase.CMSId {
		sets[db.KFieldCaseCMS] = ce.CMSNumber
	}

	if ce.CaseNumber != theCase.CaseNumber {
		sets[db.KFieldCaseNum] = ce.CaseNumber
	}

	if ce.CaseWorkerNum != theCase.CaseWorkerNum {
		sets[db.KFieldCaseWorkerNum] = ce.CaseWorkerNum
		if ce.CaseWorkerNum != 0 {
			if ce.CaseWorkerName != theCase.CaseWorker {
				sets[db.KFieldCaseWorker] = ""
			}
		}
	}

	updateCase(sets, ce.Id)

	return nil
}

// db.cases.aggregate([{$match: {"cwnum": {$ne: 0}}}, {$lookup: {from: "clients", localField: "clientnum", foreignField: "_id", as: "cd"}}])
