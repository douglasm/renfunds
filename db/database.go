package db

import (
	// "fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	MainDB                   = "ssafa"
	CollectionCounters       = "counters"
	CollectionCases          = "cases"
	CollectionClients        = "clients"
	CollectionCollections    = "collections"
	CollectionEstablishments = "establishments"
	CollectionAgents         = "agents"
	CollectionFunds          = "funds"
	CollectionLogins         = "logins"
	CollectionNumbers        = "numbers"
	CollectionPeople         = "people"
	CollectionResources      = "resources"
	CollectionSessions       = "sessions"
	CollectionServices       = "services"
	CollectionTransactions   = "transactions"
	CollectionUsers          = "users"
	CollectionVouchers       = "vouchers"

	KFieldNumCustomer = "cust"
	KFieldNumDecimal  = "decimal"
	KFieldNumBinary   = "binary"
	KFieldNumType     = "type"
	KFieldNumDeleted  = "deleted"

	KFieldUserUserName = "username"
	KFieldNameFirst    = "first"
	KFieldNameSurname  = "surname"

	DialStr = "127.0.0.1"
)

const (
	KFieldId = "_id"
)

const (
	KFieldCaseClientNum  = "clientnum"
	KFieldCaseNum        = "case"
	KFieldCaseCMS        = "cms"
	KFieldCaseFirst      = "casefirst"
	KFieldCaseSurname    = "casesurn"
	KFieldCaseWorker     = "caseworker"
	KFieldCaseCreated    = "create"
	KFieldCaseUpdated    = "updated"
	KFieldCaseWorkerNum  = "cwnum"
	KFieldCaseClosed     = "closed"
	KFieldCaseShredded   = "shredded"
	KFieldCaseAdded      = "added"
	KFieldCaseDateClosed = "dateclosed"
	FieldCaseComments    = "comments"
)

const (
	// KFieldClientsRenNumber  = "ren"
	KFieldClientsComments   = "comments"
	KFieldClientsReports    = "report"
	KFieldClientsOrder      = "order"
	KFieldClientsDOB        = "dob"
	KFieldClientsTitle      = "title"
	KFieldClientsFirst      = "first"
	KFieldClientsSurname    = "surname"
	KFieldClientsLetters    = "letters"
	KFieldClientsAddress    = "address"
	KFieldClientsPostCode   = "postcode"
	KFieldClientsPhone      = "phone"
	KFieldClientsMobile     = "mobile"
	KFieldClientsEMail      = "email"
	KFieldClientsNINum      = "ninum"
	KFieldClientsServiceNum = "serviceno"
	KFieldClientsUnit       = "services"
	// Annuity     bool      `bson:"annuity,omitempty"`
	// Comments    []Comment `bson:"comments,omitempty"`
	// Reports     []Comment `bson:"report,omitempty"`
	// UserIssuing int       `bson:"usernum,omitempty"`
	// Based       string    `bson:"based,omitempty"`
)

const (
	KFieldUserAdmin        = "admin_access"
	KFieldUserInactive     = "inactive"
	KFieldUserPosition     = "position"
	KFieldUserFirst        = "first"
	KFieldUserSurname      = "surname"
	KFieldUserAddress      = "address"
	KFieldUserPostCode     = "postcode"
	KFieldUserEMail        = "email"
	KFieldUserPhone        = "telephone"
	KFieldUserMobile       = "mobile"
	KFieldUserBased        = "based"
	KFieldUserArea         = "area"
	KFieldUserResetCode    = "resetcode"
	KFieldUserResetTime    = "resettime"
	KFieldUserActivateCode = "activate"
	KFieldUserActivateTime = "activatetime"
	KFieldUserPassword     = "password"
	KFieldUserSalt         = "salt"
)

const (
	FieldResourceName     = "name"
	FieldResourceContact  = "contact"
	FieldResourceAddress  = "address"
	FieldResourcePhone    = "phone"
	FieldResourceEmail    = "email"
	FieldResourceURL      = "url"
	FieldResourceComments = "comments"
	FieldResourceUpdated  = "updated"
)

const (
	KFieldVoucherAmount          = "amount"
	KFieldVoucherCase            = "case"
	KFieldVoucherClient          = "client"
	KFieldVoucherClosed          = "closed"
	KFieldVoucherEstablishment   = "establishment"
	KFieldVoucherDate            = "date"
	KFieldVoucherUpdated         = "update"
	KFieldVoucherUserIssuing     = "usernum"
	KFieldVoucherInvoiceReceived = "invoice"
	KFieldVoucherRemains         = "remaining"
)

type (
	Counter struct {
		ID    string `bson:"_id"`
		Value int    `bson:"value"`
	}

	Annuity struct {
		ID            int       `bson:"_id"`
		ClientNum     int       `bson:"clientnum"`
		RENNumber     string    `bson:"case,omitempty"`
		CMSID         string    `bson:"cms,omitempty"`
		Comments      []Comment `bson:"comments,omitempty"`
		Amount        int       `bson:"amount,omitempty"`
		Organisation  string    `bson:"organisation,omitempty"`
		CaseWorkerNum int       `bson:"cwnum"`
		Closed        bool      `bson:"closed"`
		DateClosed    int       `bson:"dateclosed"`
		Updated       int       `bson:"updated"`
		Created       int       `bson:"create,omitempty"`
	}

	Case struct {
		ID            int       `bson:"_id"`
		ClientNum     int       `bson:"clientnum"`
		RENNumber     string    `bson:"case,omitempty"`
		CMSID         string    `bson:"cms,omitempty"`
		VisitNumber   string    `bson:"visit,omitempty"`
		Annuity       bool      `bson:"annuity,omitempty"`
		Comments      []Comment `bson:"comments,omitempty"`
		UserIssuing   int       `bson:"usernum,omitempty"`
		CaseWorkerNum int       `bson:"cwnum"`
		CaseWorker    string    `bson:"caseworker,omitempty"`
		Closed        bool      `bson:"closed"`
		Shredded      bool      `bson:"shredded,omitempty"`
		Added         int       `bson:"added"`
		DateClosed    int       `bson:"dateclosed"`
		Updated       int       `bson:"updated"`
		Created       int       `bson:"create,omitempty"`
	}

	Client struct {
		ID  int `bson:"_id"`
		DOB int `bson:"dob,omitempty"`
		// RenNumber   string    `bson:"ren"`
		Title       string    `bson:"title,omitempty"`
		First       string    `bson:"first,omitempty"`
		Surname     string    `bson:"surname,omitempty"`
		Letters     string    `bson:"letters,omitempty"`
		Address     string    `bson:"address,omitempty"`
		PostCode    string    `bson:"postcode,omitempty"`
		Phone       string    `bson:"phone,omitempty"`
		Mobile      string    `bson:"mobile,omitempty"`
		EMail       string    `bson:"email,omitempty"`
		NINum       string    `bson:"ninum,omitempty"`
		ServiceNum  string    `bson:"serviceno,omitempty"`
		Unit        string    `bson:"services,omitempty"`
		Annuity     bool      `bson:"annuity,omitempty"`
		Comments    []Comment `bson:"comments,omitempty"`
		OldComments []Comment `bson:"oldcomm,omitempty"`
		UserIssuing int       `bson:"usernum,omitempty"`
		Based       string    `bson:"based,omitempty"`
		Alert       string    `bson:"alert,omitempty"`
		Order       int       `bson:"order"`
		Created     int       `bson:"create,omitempty"`
		Changed     int       `bson:"change,omitempty"`
	}

	Collection struct {
		ID            int     `bson:"_id"`
		Can           bool    `bson:"can,omitempty"`
		CanNumber     int     `bson:"can_number"`
		Establishment string  `bson:"establishment,omitempty"`
		Sum           float64 `bson:"sum,omitempty"`
		DateIssued    int     `bson:"date_issued,omitempty"`
		DateCollected int     `bson:"date_collected"`
		First         string  `bson:"first,omitempty"`
		Surname       string  `bson:"surname,omitempty"`
		LetterSent    bool    `bson:"letter_sent,omitempty"`
	}

	Comment struct {
		Num     int    `bson:"num"`
		Comment string `bson:"comment,omitempty"`
		Date    int    `bson:"date,omitempty"`
		User    int    `bson:"user,omitempty"`
		Name    string `bson:"name,omitempty"`
	}

	Establishment struct {
		ID            int    `bson:"_id"`
		Establishment string `bson:"establishment,omitempty"`
		Address       string `bson:"address,omitempty"`
		PostCode      string `bson:"postcode,omitempty"`
		Phone         string `bson:"telephone,omitempty"`
		Title         string `bson:"title,omitempty"`
		First         string `bson:"first,omitempty"`
		Surname       string `bson:"surname,omitempty"`
	}

	EstateAgent struct {
		ID       int    `bson:"_id"`
		Property string `bson:"property"`
		Address  string `bson:"name_and_address"`
		Title    string `bson:"title,omitempty"`
		First    string `bson:"first,omitempty"`
		Surname  string `bson:"surname,omitempty"`
		Letters  string `bson:"letters,omitempty"`
		Phone    string `bson:"telephone,omitempty"`
		URL      string `bson:"url,omitempty"`
		EMail    string `bson:"email,omitempty"`
	}

	Fund struct {
		ID           int    `json:"_id"`                     // `id` smallint(3) unsigned zerofill NOT NULL AUTO_INCREMENT,
		Source       string `json:"source,omitempty"`        // `source` varchar(255) NOT NULL DEFAULT '',
		Address      string `json:"address,omitempty"`       // `address1` varchar(255) NOT NULL DEFAULT '',
		Country      string `json:"country,omitempty"`       // `country` varchar(255) NOT NULL DEFAULT '',
		PostCode     string `json:"postcode,omitempty"`      // `postcode` varchar(255) NOT NULL DEFAULT '',
		Phone        string `json:"telephone,omitempty"`     // `telephone` varchar(255) NOT NULL DEFAULT '',
		Fax          string `json:"fax,omitempty"`           // `fax` varchar(255) NOT NULL DEFAULT '',
		EMail        string `json:"email,omitempty"`         // `email` varchar(255) NOT NULL DEFAULT '',
		FormA        bool   `json:"forma,omitempty"`         // `forma` enum('N','Y') NOT NULL DEFAULT 'N',
		URL          string `json:"url,omitempty"`           // `url` varchar(255) NOT NULL DEFAULT '',
		Title        string `json:"title,omitempty"`         // `title` varchar(255) NOT NULL DEFAULT '',
		First        string `json:"first,omitempty"`         // `firstname` varchar(255) NOT NULL DEFAULT '',
		Surname      string `json:"surname,omitempty"`       // `surname` varchar(255) NOT NULL DEFAULT '',
		Letters      string `json:"letters,omitempty"`       // `letters` varchar(255) NOT NULL DEFAULT '',
		Position     string `json:"position,omitempty"`      // `position` varchar(255) NOT NULL DEFAULT '',
		SsafaContact string `json:"ssafa_contact,omitempty"` // `ssafa_contact` varchar(255) NOT NULL DEFAULT '',
		ContactEMail string `json:"contact_email,omitempty"` // `contact_email` varchar(255) NOT NULL DEFAULT '',
		FormaEMail   string `json:"forma_email,omitempty"`   // `forma_email` varchar(255) NOT NULL DEFAULT '',
		DoesFund     string `json:"does_fund,omitempty"`     // `does_fund` mediumblob NOT NULL,
		DoesNotFund  string `json:"does_not_fund,omitempty"` // `does_not_fund` mediumblob NOT NULL,
		Comments     string `json:"comments,omitempty"`      // `comments` mediumblob NOT NULL,
		Added        int    `json:"date_added,omitempty"`    // `date_added` date NOT NULL DEFAULT '0000-00-00',
		Username     string `json:"username,omitempty"`      // `username` varchar(255) NOT NULL DEFAULT '',
	}

	People struct {
		ID          int       `bson:"_id"`
		RENNumber   string    `bson:"case,omitempty"`
		CMSID       string    `bson:"cms,omitempty"`
		VisitNumber string    `bson:"visit,omitempty"`
		DOB         int       `bson:"dob,omitempty"`
		Title       string    `bson:"title,omitempty"`
		First       string    `bson:"first,omitempty"`
		Surname     string    `bson:"surname,omitempty"`
		Letters     string    `bson:"letters,omitempty"`
		Address     string    `bson:"address,omitempty"`
		PostCode    string    `bson:"postcode,omitempty"`
		Phone       string    `bson:"phone,omitempty"`
		Mobile      string    `bson:"mobile,omitempty"`
		EMail       string    `bson:"email,omitempty"`
		NINum       string    `bson:"ninum,omitempty"`
		ServiceNo   string    `bson:"serviceno,omitempty"`
		Services    string    `bson:"services,omitempty"`
		Annuity     bool      `bson:"annuity,omitempty"`
		Comments    []Comment `bson:"comments,omitempty"`
		Reports     []Comment `bson:"reports,omitempty"`
		UserIssuing int       `bson:"usernum,omitempty"`
		CaseFirst   string    `bson:"casefirst,omitempty"`
		CaseSurn    string    `bson:"casesurn,omitempty"`
		Based       string    `bson:"based,omitempty"`
		Closed      bool      `bson:"closed,omitempty"`
		Shredded    bool      `bson:"shredded,omitempty"`
		Added       int       `bson:"added"`
		DateClosed  int       `bson:"dateclosed"`
		Updated     int       `bson:"updated"`
		CWNum       int       `bson:"cwnum"`
	}

	Resource struct {
		ID       int    `bson:"_id"`
		Name     string `bson:"name,omitempty"`
		Address  string `bson:"address,omitempty"`
		Phone    string `bson:"phone,omitempty"`
		URL      string `bson:"url,omitempty"`
		EMail    string `bson:"email,omitempty"`
		Contact  string `bson:"contact,omitempty"`
		Comments string `bson:"comments,omitempty"`
		Fund     bool   `bson:"fund,omitempty"`
		Updated  int64  `bson:"updated,omitempty"`
	}

	Servicee struct {
		ID       int    `bson:"_id"`
		Name     string `bson:"area,omitempty"`
		Contact  bool   `bson:"admin_access,omitempty"`
		Phone    string `bson:"title,omitempty"`
		EMail    string `bson:"title,omitempty"`
		URL      string `bson:"title,omitempty"`
		Comments string `bson:"title,omitempty"`
	}

	User struct {
		ID           int    `bson:"_id"`
		Area         string `bson:"area,omitempty"`
		Admin        bool   `bson:"admin_access,omitempty"`
		Title        string `bson:"title,omitempty"`
		First        string `bson:"first,omitempty"`
		Surname      string `bson:"surname,omitempty"`
		Name         string `bson:"name,omitempty"`
		Letters      string `bson:"letters,omitempty"`
		Position     string `bson:"position,omitempty"`
		Username     string `bson:"username,omitempty"`
		Password     []byte `bson:"password"`
		Salt         string `bson:"salt"`
		Address      string `bson:"address,omitempty"`
		PostCode     string `bson:"postcode,omitempty"`
		Phone        string `bson:"telephone,omitempty"`
		Mobile       string `bson:"mobile,omitempty"`
		EMail        string `bson:"email,omitempty"`
		Based        string `bson:"based,omitempty"`
		ResetCode    string `bson:"resetcode,omitempty"`
		ResetTime    int64  `bson:"resettime,omitempty"`
		ActivateCode string `bson:"activate,omitempty"`
		ActivateTime int64  `bson:"activatetime,omitempty"`
		InActive     bool   `bson:"inactive"`
		Hidden       bool   `bson:"hidden"`
		Comments     string `bson:"comments,omitempty"`
	}

	Voucher struct {
		ID              int    `bson:"_id"`
		Closed          bool   `bson:"closed,omitempty"`
		ClientID        int    `bson:"client,omitempty"`
		CaseID          int    `bson:"case,omitempty"`
		Amount          int    `bson:"amount,omitempty"`
		Establishment   string `bson:"establishment,omitempty"`
		Issued          int64  `bson:"date,omitempty"`
		Updated         int64  `bson:"update,omitempty"`
		UserIssuing     int    `bson:"usernum,omitempty"`
		InvoiceReceived bool   `bson:"invoice,omitempty"`
		Remaining       int    `bson:"remaining,omitempty"`
	}

	Proof struct {
		ID       int    `bson:"_id"`
		Closed   bool   `bson:"closed,omitempty"`
		CaseID   int    `bson:"case,omitempty"`
		ClientID int    `bson:"client,omitempty"`
		Name     string `bson:"name,omitempty"`
		Date     int    `bson:"date,omitempty"`
	}

	// Service struct {
	// 	ID              int     `bson:"_id"`                     // `id` smallint(4) unsigned zerofill NOT NULL AUTO_INCREMENT,
	// 	Closed          bool    `bson:"closed,omitempty"`        // `closed` enum('C','O') NOT NULL DEFAULT 'O',
	// 	Title           string  `bson:"title,omitempty"`         // `title` varchar(255) NOT NULL DEFAULT '',
	// 	First           string  `bson:"first,omitempty"`         // `firstname` varchar(255) NOT NULL DEFAULT '',
	// 	Surname         string  `bson:"surname,omitempty"`       // `surname` varchar(255) NOT NULL DEFAULT '',
	// 	RENNumber       string  `bson:"case_number,omitempty"`   // `case_number` varchar(255) NOT NULL DEFAULT '',
	// 	Amount          float64 `bson:"amount,omitempty"`        // `amount` double(12,2) NOT NULL DEFAULT '0.00',
	// 	Establishment   string  `bson:"establishment,omitempty"` // `establishment` varchar(255) NOT NULL DEFAULT '',
	// 	DateIssued      int     `bson:"date_issued,omitempty"`   // `date_issued` date NOT NULL DEFAULT '0000-00-00',
	// 	UserIssuing     int     `bson:"usernum,omitempty"`
	// 	IssuedFirst     string  `bson:"issued_by_first,omitempty"`   // `issued_by_firstname` varchar(255) NOT NULL DEFAULT '',
	// 	IssuedSurname   string  `bson:"issued_by_surname,omitempty"` // `issued_by_surname` varchar(255) NOT NULL DEFAULT '',
	// 	InvoiceReceived int     `bson:"invoice_received,omitempty"`  // `invoice_received` date NOT NULL DEFAULT '0000-00-00',
	// 	AmountRemaining float64 `bson:"amount_remaining,omitempty"`  // `amount_remaining` double NOT NULL DEFAULT '0',
	// }
)

var (
	MongoSession *mgo.Session
)

func GetNextSequence(name string) int {
	var (
		theCounter Counter
		err        error
	)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"value": 1}},
		ReturnNew: true,
	}
	session := MongoSession.Copy()
	countersCollection := session.DB(MainDB).C(CollectionCounters)

	_, err = countersCollection.FindId(name).Apply(change, &theCounter)
	if err != nil {
		theCounter.ID = name
		theCounter.Value = 1
		countersCollection.Insert(theCounter)
		return theCounter.Value
	}
	return theCounter.Value
}

func GetCurrentSequenceNumber(name string) int {
	var theCounter Counter

	session := MongoSession.Copy()
	countersCollection := session.DB(MainDB).C(CollectionCounters)

	err := countersCollection.FindId(name).One(&theCounter)
	if err == nil {
		return theCounter.Value
	}
	return 0
}

func GetCurrentDate() int {
	theTime := time.Now()

	theDate := theTime.Year() * 1000
	theDate += int(theTime.Month()) * 50
	theDate += theTime.Day()
	return theDate
}
