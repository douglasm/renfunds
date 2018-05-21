package vouchers

import (
	"errors"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/schema"
	"github.com/kataras/iris"

	"ssafa/crypto"
	"ssafa/db"
)

type (
	createVoucher struct {
		ID            int    `schema:"id"`
		Name          string `schema:"-"`
		Amount        string `schema:"amount"`
		Value         int    `schema:"-"`
		Establishment string `schema:"establish"`
		RENNumber     string `schema:"-"`
		Client        int    `schema:"-"`
		Commit        string `schema:"commit"`
	}

	editVoucher struct {
		ID            int    `schema:"id"`
		Name          string `schema:"-"`
		Amount        string `schema:"amount"`
		AmountVal     int    `schema:"-"`
		Establishment string `schema:"establish"`
		CMSNumber     string `schema:"-"`
		RENNumber     string `schema:"-"`
		Date          string `schema:"-"`
		Remain        string `schema:"remain"`
		RemainVal     int    `schema:"-"`
		Invoice       bool   `schema:"-"`
		InvoiceVal    int    `schema:"invoice"`
		CaseID        int    `schema:"caseid"`
		Action        string `schema:"-"`
		Commit        string `schema:"commit"`
	}
)

var (
	decoder = schema.NewDecoder()

	errNoEstablish = errors.New("You must enter an establishment")
	errNoAmount    = errors.New("You must enter an amount")
	errRemainMore  = errors.New("The amount remaining is more that the original amount")
)

func SetRoutes(app *iris.Application) {
	app.Get("/vouchers", listVouchers)
	app.Get("/vouchers/{pagenum:int}", listVouchers)
	app.Get("/voucher/{vouchernum:int}", showVoucher)
	app.Get("/addvoucher/{casenum:int}", addVoucher)
	app.Post("/addvoucher/{casenum:int}", addVoucher)
	app.Get("/vouchereditcase/{vouchernum:int}", editCaseVoucher)
	app.Post("/vouchereditcase/{vouchernum:int}", editCaseVoucher)
	app.Get("/voucheredit/{vouchernum:int}", voucherEdit)
	app.Post("/voucheredit/{vouchernum:int}", voucherEdit)
}

func getClientName(clientNum int) string {
	var (
		theClient db.Client
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	theSelect := bson.M{db.KFieldClientsFirst: 1, db.KFieldClientsSurname: 1}
	clientColl.FindId(clientNum).Select(theSelect).One(&theClient)
	theName := crypto.Decrypt(theClient.First) + " " + crypto.Decrypt(theClient.Surname)
	return strings.TrimSpace(theName)
}
