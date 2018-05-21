package vouchers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"

	// "ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
	"ssafa/users"
	"ssafa/utils"
)

type (
	voucherShow struct {
		ID            int
		Client        string
		ClientID      int
		RENNumber     string
		CaseWorker    string
		CaseID        int
		Establishment string
		Amount        string
		Invoice       string
		Remains       string
		Date          string
		Updated       string
	}
)

func listVouchers(ctx iris.Context) {
	var (
		pageNum    int
		theVoucher db.Voucher
		theCase    db.Case
		details    []voucherShow
		header     types.HeaderRecord
		navButtons types.NavButtonRecord
		err        error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	pageNum, err = ctx.Params().GetInt("pagenum")
	if err != nil {
		pageNum = 0
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: PVs"

	skip := pageNum * types.KListLimit

	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	iter := voucherColl.Find(nil).Sort(db.KFieldVoucherDate).Skip(skip).Limit(types.KListLimit + 1).Iter()
	for iter.Next(&theVoucher) {
		caseColl.FindId(theVoucher.CaseID).One(&theCase)
		newVoucher := voucherShow{}
		newVoucher.ID = theVoucher.ID
		newVoucher.CaseID = theVoucher.CaseID
		newVoucher.Establishment = theVoucher.Establishment
		newVoucher.RENNumber = theCase.RENNumber
		newVoucher.Client = getClientName(theVoucher.ClientID)
		newVoucher.ClientID = theVoucher.ClientID
		newVoucher.Date = utils.GetDateAndTime(theVoucher.Issued, false)
		newVoucher.Amount = "£" + utils.IntToString(theVoucher.Amount, 2)
		newVoucher.CaseWorker = users.GetUserName(theVoucher.UserIssuing)

		if theVoucher.InvoiceReceived {
			newVoucher.Invoice = "Yes"
		} else {
			newVoucher.Invoice = "No"
		}
		newVoucher.Remains = "£" + utils.IntToString(theVoucher.Remaining, 2)
		details = append(details, newVoucher)
	}
	iter.Close()

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("vouchers/list.html")
}

func showVoucher(ctx iris.Context) {
	var (
		voucherNum int
		theVoucher db.Voucher
		theCase    db.Case
		details    voucherShow
		header     types.HeaderRecord
		err        error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	voucherNum, err = ctx.Params().GetInt("vouchernum")
	if err != nil {
		voucherNum = 0
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	voucherColl.FindId(voucherNum).One(&theVoucher)
	caseColl.FindId(theVoucher.CaseID).One(&theCase)

	details.ID = theVoucher.ID
	details.CaseID = theVoucher.CaseID
	details.Establishment = theVoucher.Establishment
	details.RENNumber = theCase.RENNumber
	details.Client = getClientName(theVoucher.ClientID)
	details.ClientID = theVoucher.ClientID
	details.Date = utils.GetDateAndTime(theVoucher.Issued, false)
	details.Updated = utils.GetDateAndTime(theVoucher.Updated, false)
	details.CaseWorker = users.GetUserName(theVoucher.UserIssuing)

	details.Amount = "£" + utils.IntToString(theVoucher.Amount, 2)
	if theVoucher.InvoiceReceived {
		details.Invoice = "Yes"
	} else {
		details.Invoice = "No"
	}
	details.Remains = "£" + utils.IntToString(theVoucher.Remaining, 2)

	header.Title = "RF: PV " + details.RENNumber

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("vouchers/show.html")
}

func addVoucher(ctx iris.Context) {
	var (
		header       types.HeaderRecord
		details      createVoucher
		theCase      db.Case
		errorMessage string
		caseNum      int
		err          error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	caseNum, err = ctx.Params().GetInt("casenum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	details.ID = caseNum

	session := db.MongoSession.Copy()
	defer session.Close()

	// voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)
	// clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	caseColl.FindId(caseNum).One(&theCase)
	details.RENNumber = theCase.RENNumber

	// clientColl.FindId(theCase.ClientNum).One(&theClient)
	details.Name = getClientName(theCase.ClientNum)
	details.Client = theCase.ClientNum

	switch ctx.Method() {
	case http.MethodGet:

	case http.MethodPost:
		ctx.FormValues()
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode add voucher", err)
		}
		err = details.checkAdd()
		if err == nil {
			details.saveAdd(theSession.(users.Session).UserNumber)
			theURL := fmt.Sprintf("/case/%d", details.ID)
			ctx.Redirect(theURL, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: Add voucher"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("vouchers/add.html")

}

func editCaseVoucher(ctx iris.Context) {
	var (
		header       types.HeaderRecord
		details      editVoucher
		theCase      db.Case
		theVoucher   db.Voucher
		errorMessage string
		voucherNum   int
		err          error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	voucherNum, err = ctx.Params().GetInt("vouchernum")
	if err != nil {
		ctx.Redirect("/cases", http.StatusFound)
		return
	}

	details.ID = voucherNum

	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	voucherColl.FindId(voucherNum).One(&theVoucher)
	caseColl.FindId(theVoucher.CaseID).One(&theCase)
	details.CaseID = theVoucher.CaseID

	details.Name = getClientName(theCase.ClientNum)
	details.Date = utils.GetDateAndTime(theVoucher.Issued, false)

	switch ctx.Method() {
	case http.MethodGet:
		details.Amount = utils.IntToString(theVoucher.Amount, 2)
		details.Remain = utils.IntToString(theVoucher.Remaining, 2)
		details.Establishment = theVoucher.Establishment
		details.Invoice = theVoucher.InvoiceReceived

	case http.MethodPost:
		ctx.FormValues()
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode add voucher", err)
		}
		err = details.checkEdit()
		if err == nil {
			details.saveEdit(theSession.(users.Session).UserNumber, &theVoucher)
			theURL := fmt.Sprintf("/case/%d", details.CaseID)
			ctx.Redirect(theURL, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}
	details.Action = fmt.Sprintf("/vouchereditcase/%d", details.ID)

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: Add voucher"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("vouchers/edit.html")
}

func voucherEdit(ctx iris.Context) {
	var (
		header       types.HeaderRecord
		details      editVoucher
		theCase      db.Case
		theVoucher   db.Voucher
		errorMessage string
		voucherNum   int
		err          error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}
	voucherNum, err = ctx.Params().GetInt("vouchernum")
	if err != nil {
		ctx.Redirect("/vouchers", http.StatusFound)
		return
	}

	details.ID = voucherNum

	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)

	voucherColl.FindId(voucherNum).One(&theVoucher)
	caseColl.FindId(theVoucher.CaseID).One(&theCase)
	details.CaseID = theVoucher.CaseID

	details.Name = getClientName(theCase.ClientNum)
	details.Date = utils.GetDateAndTime(theVoucher.Issued, false)

	switch ctx.Method() {
	case http.MethodGet:
		details.Amount = utils.IntToString(theVoucher.Amount, 2)
		details.Remain = utils.IntToString(theVoucher.Remaining, 2)
		details.Establishment = theVoucher.Establishment
		details.Invoice = theVoucher.InvoiceReceived

	case http.MethodPost:
		ctx.FormValues()
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode add voucher", err)
		}
		err = details.checkEdit()
		if err == nil {
			details.saveEdit(theSession.(users.Session).UserNumber, &theVoucher)
			theURL := fmt.Sprintf("/voucher/%d", voucherNum)
			ctx.Redirect(theURL, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}
	details.Action = fmt.Sprintf("/voucheredit/%d", details.ID)

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: Edit voucher"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("vouchers/edit.html")
}

func getVouchers(pageNum int) ([]voucherShow, bool) {
	var (
		theVoucher db.Voucher
		theCase    db.Case
		theClient  db.Client
		// theUser    db.User
		theList []voucherShow
		err     error
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	caseColl := session.DB(db.MainDB).C(db.CollectionCases)
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)
	// userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	skip := pageNum * types.KListLimit
	limit := types.KListLimit + 1

	// found := 0
	iter := voucherColl.Find(nil).Skip(skip).Limit(limit).Sort(db.KFieldClientsOrder).Iter()
	for iter.Next(&theVoucher) {

		if len(theList) < types.KListLimit {
			newVoucher := voucherShow{ID: theVoucher.ID}
			newVoucher.Establishment = theVoucher.Establishment
			newVoucher.Date = utils.GetDateAndTime(theVoucher.Issued, false)
			err = caseColl.FindId(theVoucher.CaseID).One(&theCase)
			if err == nil {
				newVoucher.RENNumber = theCase.RENNumber
				err = clientColl.FindId(theCase.ClientNum).One(&theClient)
				if err == nil {
				}
			}
			// err = userColl.FindId(theVoucher.UserIssuing).One(&theUser)
			// if err == nil {
			// 	newVoucher.
			// }

			// Client        string
			// Case          string
			// Establishment string
			// Amount        string
			// Remains       string
			// Issued          string
			theList = append(theList, newVoucher)
		} else {
			return theList, true
		}
	}
	iter.Close()

	return theList, false

}

func (cv *createVoucher) checkAdd() error {
	if len(cv.Establishment) == 0 {
		return errNoEstablish
	}
	val, err := utils.StringToInt(cv.Amount)
	if err != nil {
		return err
	}
	cv.Value = val
	return nil
}

func (cv *createVoucher) saveAdd(userNumber int) {
	var (
		theVoucher db.Voucher
	)

	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)

	theVoucher.ID = db.GetNextSequence(db.CollectionVouchers)
	theVoucher.ClientID = cv.Client
	theVoucher.CaseID = cv.ID
	theVoucher.Amount = cv.Value
	theVoucher.Remaining = cv.Value
	theVoucher.Establishment = cv.Establishment
	theVoucher.Issued = time.Now().Unix()
	theVoucher.Updated = theVoucher.Issued
	theVoucher.UserIssuing = userNumber

	err := voucherColl.Insert(&theVoucher)
	if err != nil {
		log.Println("Error: Voucher insert", err)
	}
}

func (ev *editVoucher) checkEdit() error {
	var (
		err error
	)
	if len(ev.Amount) == 0 {
		return errNoAmount
	}
	if len(ev.Establishment) == 0 {
		return errNoEstablish
	}
	if len(ev.Remain) == 0 {
		return errNoEstablish
	}
	ev.AmountVal, err = utils.StringToInt(ev.Amount)
	if err != nil {
		return err
	}
	ev.RemainVal, err = utils.StringToInt(ev.Remain)
	if err != nil {
		return err
	}
	if ev.RemainVal > ev.AmountVal {
		return errRemainMore
	}
	if ev.InvoiceVal == 1 {
		ev.Invoice = true
	}
	return nil
}

func (ev *editVoucher) saveEdit(userNum int, theVoucher *db.Voucher) {

	session := db.MongoSession.Copy()
	defer session.Close()

	voucherColl := session.DB(db.MainDB).C(db.CollectionVouchers)
	sets := bson.M{}
	if ev.AmountVal != theVoucher.Amount {
		sets[db.KFieldVoucherAmount] = ev.AmountVal
	}

	if ev.RemainVal != theVoucher.Remaining {
		sets[db.KFieldVoucherRemains] = ev.RemainVal
	}

	if ev.Establishment != theVoucher.Establishment {
		sets[db.KFieldVoucherEstablishment] = ev.Establishment
	}

	if ev.Invoice != theVoucher.InvoiceReceived {
		sets[db.KFieldVoucherInvoiceReceived] = ev.Invoice
	}

	sets[db.KFieldVoucherUpdated] = time.Now().Unix()

	err := voucherColl.UpdateId(theVoucher.ID, bson.M{"$set": sets})
	if err != nil {
		log.Println("Error: update voucher", err)
	}
}
