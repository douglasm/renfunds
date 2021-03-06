package cases

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	// "sort"

	"github.com/kataras/iris"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
	"ssafa/users"
	"ssafa/utils"
)

type (
	caseCreate struct {
		ID         int    `schema:"id"`
		RENNumber  string `schema:"casenumber"`
		CMSNumber  string `schema:"cms"`
		ClientNum  int    `schema:"client"`
		ClientName string `schema:"_"`
		Text       string `schema:"text"`
		Checkfield int    `schema:"checkfield"`
		Commit     string `schema:"commit"`
	}

	caseDisplay struct {
		ID          int
		Open        bool
		Voucher     bool
		ClientName  types.RowItem
		RENNumber   string
		CaseWorker  types.RowItem
		CMSNumber   types.RowItem
		Opened      types.RowItem
		State       types.RowItem
		Updated     types.RowItem
		ClientNum   int
		HasVoucher  bool
		VoucherList []voucherItem
		Reports     []CommentDisplay
	}

	caseEdit struct {
		ID             int           `schema:"id"`
		RENNumber      string        `schema:"casenumber"`
		CaseWorker     template.HTML `schema:"-"`
		CaseWorkerName string        `schema:"cwname"`
		CaseWorkerNum  int           `schema:"cwnum"`
		CMSNumber      string        `schema:"cms"`
		ClientNum      int           `schema:"client"`
		ClientName     string        `schema:"clientname"`
		Checkfield     int           `schema:"checkfield"`
		Commit         string        `schema:"commit"`
	}

	voucherItem struct {
		ID            int
		Date          string
		Amount        string
		Establishment string
		Remaining     string
		Invoice       string
		Updated       string
	}
)

func showCase(ctx iris.Context) {
	var (
		details    caseDisplay
		theCase    db.Case
		theClient  db.Client
		theVoucher db.Voucher
		header     types.HeaderRecord
		caseNum    int
		err        error
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
	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)

	err = caseColl.FindId(caseNum).One(&theCase)
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	err = clientColl.FindId(theCase.ClientNum).One(&theClient)

	details.ID = caseNum
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	header.Title = "RF: Case " + theCase.RENNumber

	details.RENNumber = theCase.RENNumber

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
	details.CMSNumber.Value = template.HTML(theCase.CMSID)

	details.Opened.Title = "Opened"
	details.Opened.Value = template.HTML(utils.DateToString(theCase.Created))

	details.Updated.Title = "Updated"
	details.Updated.Value = template.HTML(utils.DateToString(theCase.Updated))

	details.State.Title = "Closed"
	if theCase.Closed {
		details.State.Value = template.HTML("Yes " + utils.DateToString(theCase.DateClosed))
	} else {
		details.Open = true
		details.Voucher = true
		details.State.Value = template.HTML("No")
	}

	iter := voucherColl.Find(bson.M{db.KFieldVoucherCase: theCase.ID}).Iter()
	for iter.Next(&theVoucher) {
		details.HasVoucher = true
		newVoucher := voucherItem{}
		newVoucher.ID = theVoucher.ID
		newVoucher.Date = utils.GetDateAndTime(theVoucher.Issued, false)
		newVoucher.Amount = "£" + utils.IntToString(theVoucher.Amount, 2)
		newVoucher.Establishment = theVoucher.Establishment
		newVoucher.Remaining = "£" + utils.IntToString(theVoucher.Remaining, 2)
		if theVoucher.InvoiceReceived {
			newVoucher.Invoice = "Yes"
		} else {
			newVoucher.Invoice = "No"
		}
		newVoucher.Updated = utils.GetDateAndTime(theVoucher.Updated, false)
		details.VoucherList = append(details.VoucherList, newVoucher)
	}
	iter.Close()

	details.ClientNum = theClient.ID

	theCase.Comments = checkComments(theCase.Comments, details.ID)
	for _, item := range theCase.Comments {
		newComment := CommentDisplay{}
		newComment.GetCommentDisplay(item, details.ID)
		details.Reports = append(details.Reports, newComment)
	}

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("/cases/show.html")
}

func deleteCase(ctx iris.Context) {
	var (
		details   caseDisplay
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
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	if !header.Admin {
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

	details.ID = caseNum

	header.Title = "RF: Delete Case " + theCase.RENNumber

	details.RENNumber = theCase.RENNumber

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
	details.CMSNumber.Value = template.HTML(theCase.CMSID)

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

	details.ClientNum = theClient.ID
	for _, item := range theCase.Comments {
		newComment := CommentDisplay{}
		newComment.GetCommentDisplay(item, details.ID)
		details.Reports = append(details.Reports, newComment)
	}

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("cases/delete.html")
}

func removeCase(ctx iris.Context) {
	var (
		caseNum int
		err     error
	)

	caseNum, err = ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	theSession := ctx.Values().Get("session")

	if !theSession.(users.Session).Admin {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	err = caseColl.RemoveId(caseNum)
	if err != nil {
		log.Println("Error: Remove case", err)
	}
	ctx.Redirect("/cases", http.StatusFound)
}

func listCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header     types.HeaderRecord
		navButtons types.NavButtonRecord
		pageNum    int
		err        error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	theSession := ctx.Values().Get("session")

	match := bson.M{"$match": bson.M{db.KFieldCaseClientNum: bson.M{"$ne": 0}}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, navButtons.HasNext = GetCases(pageNum, match, sort)

	navButtons.FillNav(navButtons.HasNext, pageNum, "cases")

	header.Title = "RF: Cases"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("cases/list.html")
}

func openCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header     types.HeaderRecord
		navButtons types.NavButtonRecord
		pageNum    int
		err        error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	theSession := ctx.Values().Get("session")
	match := bson.M{"$match": bson.M{db.KFieldCaseClosed: false}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, navButtons.HasNext = GetCases(pageNum, match, sort)

	navButtons.FillNav(navButtons.HasNext, pageNum, "casesopen")

	header.Title = "RF: Open Cases"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("cases/open.html")
}

func unassignedCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header     types.HeaderRecord
		navButtons types.NavButtonRecord
		pageNum    int
		err        error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	theSession := ctx.Values().Get("session")
	match := bson.M{"$match": bson.M{db.KFieldCaseWorkerNum: 0, db.KFieldCaseClosed: false}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, navButtons.HasNext = GetCases(pageNum, match, sort)

	navButtons.FillNav(navButtons.HasNext, pageNum, "casesunassign")

	header.Title = "RF: Unassigned Cases"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("cases/unassign.html")
}

func inactiveCases(ctx iris.Context) {
	var (
		details []CaseList
		// theCase db.Case
		// theClient db.Client
		header     types.HeaderRecord
		navButtons types.NavButtonRecord
		pageNum    int
		err        error
	)

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	theSession := ctx.Values().Get("session")
	match := bson.M{"$match": bson.M{db.KFieldCaseWorkerNum: 0, db.KFieldCaseClosed: false}}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseUpdated: 1}}

	details, navButtons.HasNext = GetCases(pageNum, match, sort)

	navButtons.FillNav(navButtons.HasNext, pageNum, "casesinactive")

	header.Title = "RF: Inactive Cases"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("cases/inactive.html")
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

func createCase(ctx iris.Context) {
	var (
		theCase      db.Case
		theClient    db.Client
		details      caseCreate
		errorMessage string
		header       types.HeaderRecord
		err          error
	)

	clientNum, err := ctx.Params().GetInt("clientnum")
	if err != nil {
		ctx.Redirect("/clients", http.StatusFound)
		return
	}

	theSession := ctx.Values().Get("session")

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	clientColl.FindId(clientNum).One(&theClient)

	switch ctx.Method() {
	case http.MethodGet:

	case http.MethodPost:
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode casecreate", err)
			errorMessage = err.Error()
			break
		}
		err = details.check()
		if err == nil {
			theCase.ID = db.GetNextSequence(db.CollectionCases)
			theCase.ClientNum = clientNum
			theCase.Created = utils.CurrentDate()
			theCase.Updated = theCase.Created
			theCase.CMSID = details.CMSNumber
			theCase.RENNumber = details.RENNumber

			newComment := db.Comment{}
			newComment.Comment = crypto.Encrypt(details.Text)
			newComment.Date = theCase.Created
			newComment.User = theSession.(users.Session).UserNumber

			theCase.Comments = append(theCase.Comments, newComment)

			caseColl.Insert(&theCase)
			tempStr := fmt.Sprintf("/case/%d", theCase.ID)
			ctx.Redirect(tempStr, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	details.ClientNum = clientNum
	details.ClientName = crypto.Decrypt(theClient.First) + " " + crypto.Decrypt(theClient.Surname)

	header.Admin = theSession.(users.Session).Admin
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Title = "Add new case"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("cases/create.html")
}

func editCase(ctx iris.Context) {
	var (
		theCase      db.Case
		details      caseEdit
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
		details.ID = caseNum
		details.RENNumber = theCase.RENNumber
		details.CMSNumber = theCase.CMSID
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
		err = details.save()
		if err == nil {
			theUrl := fmt.Sprintf("/case/%d", details.ID)
			ctx.Redirect(theUrl, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Admin = theSession.(users.Session).Admin
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Title = "Edit case"
	header.Scripts = append(header.Scripts, "getusers")

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("cases/edit.html")
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

func searchCase(ctx iris.Context) {
	var (
		details []CaseList
		header  types.HeaderRecord
	)

	theSession := ctx.Values().Get("session")

	searchStr := ctx.FormValue("search")

	reg := bson.M{"$regex": searchStr, "$options": "-i"}

	theQuery := bson.M{"$or": types.S{{db.KFieldCaseNum: reg}, {db.KFieldCaseCMS: reg}}}

	match := bson.M{"$match": theQuery}
	sort := bson.M{"$sort": bson.M{db.KFieldCaseCreated: 1}}

	details, _ = GetCases(0, match, sort)

	header.Title = "RF: Cases"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("cases/list.html")
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

func editComment(ctx iris.Context) {
	var (
		theCase      db.Case
		details      commentEdit
		header       types.HeaderRecord
		caseNum      int
		commentNum   int
		gotOne       bool
		errorMessage string
		err          error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	caseNum, err = ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	commentNum, err = ctx.Params().GetInt("commentnum")
	if err != nil {
		theUrl := fmt.Sprintf("/case/%d", caseNum)
		ctx.Redirect(theUrl, http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	details.Message = "Edit case comment"
	details.Link = fmt.Sprintf("/casecomment/%d/%d", caseNum, commentNum)
	details.ID = caseNum
	details.Num = commentNum

	switch ctx.Method() {
	case http.MethodGet:
		err = caseColl.FindId(caseNum).One(&theCase)
		if err != nil {
			theUrl := fmt.Sprintf("/case/%d", caseNum)
			ctx.Redirect(theUrl, http.StatusFound)
			return
		}

		for _, item := range theCase.Comments {
			if item.Num == commentNum {
				details.Comment = decodeComment(item.Comment)
				// details.Comment = crypto.Decrypt(item.Comment)
				gotOne = true
				break
			}
		}
		if !gotOne {
			theUrl := fmt.Sprintf("/case/%d", caseNum)
			ctx.Redirect(theUrl, http.StatusFound)
			return
		}

	case http.MethodPost:
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: edit case comment decode", err)
		}

		err = caseColl.FindId(caseNum).One(&theCase)
		if err != nil {
			theUrl := fmt.Sprintf("/case/%d", caseNum)
			ctx.Redirect(theUrl, http.StatusFound)
			return
		}

		if len(details.Comment) == 0 {
			theUrl := fmt.Sprintf("/case/%d", caseNum)
			ctx.Redirect(theUrl, http.StatusFound)
			return
		}

		newComments := []db.Comment{}
		for _, item := range theCase.Comments {
			if item.Num == commentNum {
				item.Comment = crypto.Encrypt(details.Comment)
				gotOne = true
			}
			newComments = append(newComments, item)
		}
		if gotOne {
			err = caseColl.UpdateId(caseNum, bson.M{"$set": bson.M{db.FieldCaseComments: newComments}})
			if err != nil {
				log.Println("Update err:", err)
			}
		}

		theUrl := fmt.Sprintf("/case/%d", caseNum)
		ctx.Redirect(theUrl, http.StatusFound)
		return
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: Edit case comment"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("commentedit.html")
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
		newCase := CaseList{ID: result.ID, RENNumber: result.RENNumber, CMSNumber: result.CMSID}
		if len(newCase.RENNumber) == 0 {
			newCase.RENNumber = "None"
		}
		newCase.ClientID = result.Client[0].ID
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

func (ce *caseEdit) save() error {
	var (
		theCase db.Case
	)
	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)
	// KFieldCaseNum        = "case"
	// KFieldCaseCMS        = "cms"

	gotOne := false
	ce.RENNumber = strings.TrimSpace(ce.RENNumber)
	if len(ce.RENNumber) > 0 {
		ce.RENNumber = strings.ToUpper(ce.RENNumber)
		iter := caseColl.Find(bson.M{db.KFieldCaseNum: ce.RENNumber}).Iter()
		for iter.Next(&theCase) {
			if theCase.ID != ce.ID {
				gotOne = true
				break
			}
		}
		iter.Close()

		if gotOne {
			return ErrorDateCaseUsed
		}
	}

	ce.CMSNumber = strings.TrimSpace(ce.CMSNumber)
	if len(ce.CMSNumber) < 5 {
		if ce.CaseWorkerNum != 0 {
			return errorCMSMissing
		}
	}

	iter := caseColl.Find(bson.M{db.KFieldCaseCMS: ce.CMSNumber}).Iter()
	for iter.Next(&theCase) {
		if theCase.ID != ce.ID {
			gotOne = true
			break
		}
	}
	iter.Close()

	if gotOne {
		return ErrorDateCMSUsed
	}

	err := caseColl.FindId(ce.ID).One(&theCase)
	if err != nil {
		log.Println("Error: case save read", err)
	}
	sets := bson.M{}
	if ce.CMSNumber != theCase.CMSID {
		sets[db.KFieldCaseCMS] = ce.CMSNumber
	}

	if ce.RENNumber != theCase.RENNumber {
		sets[db.KFieldCaseNum] = ce.RENNumber
	}

	if len(ce.CaseWorkerName) != 0 {
		if ce.CaseWorkerNum != theCase.CaseWorkerNum {
			sets[db.KFieldCaseWorkerNum] = ce.CaseWorkerNum
			if ce.CaseWorkerNum != 0 {
				if ce.CaseWorkerName != theCase.CaseWorker {
					sets[db.KFieldCaseWorker] = ""
				}
			}
		}
	}

	updateCase(sets, ce.ID)

	return nil
}

func (cc *caseCreate) check() error {
	var (
		theCase db.Case
		err     error
	)

	cc.CMSNumber = strings.TrimSpace(cc.CMSNumber)
	cc.RENNumber = strings.TrimSpace(cc.RENNumber)
	cc.Text = strings.TrimSpace(cc.Text)

	session := db.MongoSession.Copy()
	defer session.Close()
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	err = caseColl.Find(bson.M{db.KFieldCaseCMS: cc.CMSNumber}).One(&theCase)
	if err == nil {
		return ErrorDateCMSUsed
	}

	if len(cc.RENNumber) > 0 {
		cc.RENNumber = strings.ToUpper(cc.RENNumber)
		err = caseColl.Find(bson.M{db.KFieldCaseNum: cc.RENNumber}).One(&theCase)
		if err == nil {
			return ErrorDateCaseUsed
		}
	}

	if len(cc.Text) == 0 {
		return errorReasonMissing
	}
	return nil
}

func checkComments(comments []db.Comment, caseNum int) []db.Comment {
	var (
		newComments  []db.Comment
		tempComments []db.Comment
	)
	for _, item := range comments {
		theStr := item.Comment
		if len(item.Comment) == 0 {
			continue
		}
		looping := true
		for looping {
			_, err := hex.DecodeString(theStr)
			if err != nil {
				looping = false
			} else {
				theStr = crypto.Decrypt(theStr)
			}

			item.Comment = theStr
		}
		tempComments = append(tempComments, item)
	}
	for i, item := range tempComments {
		item.Num = i + 1
		if len(item.Comment) > 0 {
			item.Comment = crypto.Encrypt(item.Comment)
		}
		newComments = append(newComments, item)
	}
	if len(newComments) != 0 {
		session := db.MongoSession.Copy()
		defer session.Close()
		caseColl := session.DB(db.MainDB).C(db.CollectionCases)
		update := bson.M{"$set": bson.M{db.FieldCaseComments: newComments}}
		caseColl.UpdateId(caseNum, update)
	}
	return newComments
}
