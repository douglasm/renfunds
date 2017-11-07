package db

import (
	// "fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	KFieldClientsNiNum    = "ninum"
	KFieldClientsComments = "comments"
	KFieldClientsReports  = "report"
	KFieldClientsOrder    = "order"
)

const (
	KFieldCaseClientNum = "clientnum"
)

type (
	Counter struct {
		Id    string `bson:"_id"`
		Value int    `bson:"value"`
	}

	Case struct {
		Id          int      `bson:"_id"`
		ClientNum   int      `bson:"clientnum"`
		CaseNumber  string   `bson:"case,omitempty"`
		CMSId       string   `bson:"cms,omitempty"`
		VisitNumber string   `bson:"visit,omitempty"`
		Annuity     bool     `bson:"annuity,omitempty"`
		Comments    []string `bson:"comments,omitempty"`
		Reports     []string `bson:"report,omitempty"`
		UserIssuing int      `bson:"usernum,omitempty"`
		CaseId      string   `bson:"cwid,omitempty"`
		CaseFirst   string   `bson:"casefirst,omitempty"`
		CaseSurname string   `bson:"casesurn,omitempty"`
		Closed      bool     `bson:"closed,omitempty"`
		Shredded    bool     `bson:"shredded,omitempty"`
		Added       int      `bson:"added"`
		DateClosed  int      `bson:"dateclosed"`
		Updated     int      `bson:"updated"`
		CWNum       int      `bson:"cwnum"`
	}

	Client struct {
		Id          int      `bson:"_id"`
		DOB         int      `bson:"dob,omitempty"`
		Title       string   `bson:"title,omitempty"`
		First       string   `bson:"first,omitempty"`
		Surname     string   `bson:"surname,omitempty"`
		Letters     string   `bson:"letters,omitempty"`
		Address     string   `bson:"address,omitempty"`
		PostCode    string   `bson:"postcode,omitempty"`
		Phone       string   `bson:"phone,omitempty"`
		Mobile      string   `bson:"mobile,omitempty"`
		EMail       string   `bson:"email,omitempty"`
		NINum       string   `bson:"ninum,omitempty"`
		ServiceNo   string   `bson:"serviceno,omitempty"`
		Services    string   `bson:"services,omitempty"`
		Annuity     bool     `bson:"annuity,omitempty"`
		Comments    []string `bson:"comments,omitempty"`
		Reports     []string `bson:"report,omitempty"`
		UserIssuing int      `bson:"usernum,omitempty"`
		Based       string   `bson:"based,omitempty"`
		Order       int      `bson:"order"`
	}

	Collection struct {
		Id            int     `bson:"_id"`
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
		Comment string `bson:"comment,omitempty"`
		Date    int    `bson:"date,omitempty"`
		User    int    `bson:"user,omitempty"`
		Name    string `bson:"name,omitempty"`
	}

	Establishment struct {
		Id            int    `bson:"_id"`
		Establishment string `bson:"establishment,omitempty"`
		Address       string `bson:"address,omitempty"`
		PostCode      string `bson:"postcode,omitempty"`
		Phone         string `bson:"telephone,omitempty"`
		Title         string `bson:"title,omitempty"`
		First         string `bson:"first,omitempty"`
		Surname       string `bson:"surname,omitempty"`
	}

	EstateAgent struct {
		Id       int    `bson:"_id"`
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
		Id           int    `json:"_id"`                     // `id` smallint(3) unsigned zerofill NOT NULL AUTO_INCREMENT,
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
		Id          int      `bson:"_id"`
		CaseNumber  string   `bson:"case,omitempty"`
		CMSId       string   `bson:"cms,omitempty"`
		VisitNumber string   `bson:"visit,omitempty"`
		DOB         int      `bson:"dob,omitempty"`
		Title       string   `bson:"title,omitempty"`
		First       string   `bson:"first,omitempty"`
		Surname     string   `bson:"surname,omitempty"`
		Letters     string   `bson:"letters,omitempty"`
		Address     string   `bson:"address,omitempty"`
		PostCode    string   `bson:"postcode,omitempty"`
		Phone       string   `bson:"phone,omitempty"`
		Mobile      string   `bson:"mobile,omitempty"`
		EMail       string   `bson:"email,omitempty"`
		NINum       string   `bson:"ninum,omitempty"`
		ServiceNo   string   `bson:"serviceno,omitempty"`
		Services    string   `bson:"services,omitempty"`
		Annuity     bool     `bson:"annuity,omitempty"`
		Comments    []string `bson:"comments,omitempty"`
		Reports     []string `bson:"report,omitempty"`
		UserIssuing int      `bson:"usernum,omitempty"`
		CaseFirst   string   `bson:"casefirst,omitempty"`
		CaseSurn    string   `bson:"casesurn,omitempty"`
		Based       string   `bson:"based,omitempty"`
		Closed      bool     `bson:"closed,omitempty"`
		Shredded    bool     `bson:"shredded,omitempty"`
		Added       int      `bson:"added"`
		DateClosed  int      `bson:"dateclosed"`
		Updated     int      `bson:"updated"`
		CWNum       int      `bson:"cwnum"`
	}

	User struct {
		Id       int    `bson:"_id"`
		Area     string `bson:"area,omitempty"`
		Admin    bool   `bson:"admin_access,omitempty"`
		Title    string `bson:"title,omitempty"`
		First    string `bson:"first,omitempty"`
		Surname  string `bson:"surname,omitempty"`
		Letters  string `bson:"letters,omitempty"`
		Position string `bson:"position,omitempty"`
		Username string `bson:"username,omitempty"`
		Password []byte `bson:"password"`
		Salt     string `bson:"salt"`
		Address  string `bson:"address,omitempty"`
		PostCode string `bson:"postcode,omitempty"`
		Phone    string `bson:"telephone,omitempty"`
		Mobile   string `bson:"mobile,omitempty"`
		EMail    string `bson:"email,omitempty"`
		Based    string `bson:"based,omitempty"`
		Comments string `bson:"comments,omitempty"`
	}

	Voucher struct {
		Id              int     `bson:"_id"`
		Closed          bool    `bson:"closed,omitempty"`
		Title           string  `bson:"title,omitempty"`
		First           string  `bson:"first,omitempty"`
		Surname         string  `bson:"surname,omitempty"`
		CaseNumber      string  `bson:"casenumber,omitempty"`
		Amount          float64 `bson:"amount,omitempty"`
		Establishment   string  `bson:"establishment,omitempty"`
		DateIssued      int     `bson:"date_issued,omitempty"`
		UserIssuing     int     `bson:"usernum,omitempty"`
		IssuedFirst     string  `bson:"issued_by_first,omitempty"`
		IssuedSurname   string  `bson:"issued_by_surname,omitempty"`
		InvoiceReceived int     `bson:"invoice_received,omitempty"`
		AmountRemaining float64 `bson:"amount_remaining,omitempty"`
	}

	Resource struct {
		Id           int    `bson:"_id"`                     // `id` smallint(3) unsigned zerofill NOT NULL AUTO_INCREMENT,
		Resource     string `bson:"resource,omitempty"`      // `resource` varchar(255) NOT NULL DEFAULT '',
		Address      string `bson:"address,omitempty"`       // `address1` varchar(255) NOT NULL DEFAULT '',
		Postcode     string `bson:"postcode,omitempty"`      // `postcode` varchar(255) NOT NULL DEFAULT '',
		Phone        string `bson:"telephone,omitempty"`     // `telephone` varchar(255) NOT NULL DEFAULT '',
		Mobile       string `bson:"mobile,omitempty"`        // `mobile` varchar(255) NOT NULL DEFAULT '',
		Url          string `bson:"url,omitempty"`           // `url` varchar(255) NOT NULL DEFAULT '',
		Email        string `bson:"email,omitempty"`         // `email` varchar(255) NOT NULL DEFAULT '',
		Title        string `bson:"title,omitempty"`         // `title` varchar(255) NOT NULL DEFAULT '',
		First        string `bson:"first,omitempty"`         // `firstname` varchar(255) NOT NULL DEFAULT '',
		Surname      string `bson:"surname,omitempty"`       // `surname` varchar(255) NOT NULL DEFAULT '',
		Position     string `bson:"position,omitempty"`      // `position` varchar(255) NOT NULL DEFAULT '',
		ContactEMail string `bson:"contact_email,omitempty"` // `contact_email` varchar(255) NOT NULL DEFAULT '',
		Comments     string `bson:"comments,omitempty"`      // `comments` mediumblob NOT NULL,
		DateAdded    int    `bson:"date_added,omitempty"`    // `date_added` date NOT NULL DEFAULT '0000-00-00',
	}

	Service struct {
		Id              int     `bson:"_id"`                     // `id` smallint(4) unsigned zerofill NOT NULL AUTO_INCREMENT,
		Closed          bool    `bson:"closed,omitempty"`        // `closed` enum('C','O') NOT NULL DEFAULT 'O',
		Title           string  `bson:"title,omitempty"`         // `title` varchar(255) NOT NULL DEFAULT '',
		First           string  `bson:"first,omitempty"`         // `firstname` varchar(255) NOT NULL DEFAULT '',
		Surname         string  `bson:"surname,omitempty"`       // `surname` varchar(255) NOT NULL DEFAULT '',
		CaseNumber      string  `bson:"case_number,omitempty"`   // `case_number` varchar(255) NOT NULL DEFAULT '',
		Amount          float64 `bson:"amount,omitempty"`        // `amount` double(12,2) NOT NULL DEFAULT '0.00',
		Establishment   string  `bson:"establishment,omitempty"` // `establishment` varchar(255) NOT NULL DEFAULT '',
		DateIssued      int     `bson:"date_issued,omitempty"`   // `date_issued` date NOT NULL DEFAULT '0000-00-00',
		UserIssuing     int     `bson:"usernum,omitempty"`
		IssuedFirst     string  `bson:"issued_by_first,omitempty"`   // `issued_by_firstname` varchar(255) NOT NULL DEFAULT '',
		IssuedSurname   string  `bson:"issued_by_surname,omitempty"` // `issued_by_surname` varchar(255) NOT NULL DEFAULT '',
		InvoiceReceived int     `bson:"invoice_received,omitempty"`  // `invoice_received` date NOT NULL DEFAULT '0000-00-00',
		AmountRemaining float64 `bson:"amount_remaining,omitempty"`  // `amount_remaining` double NOT NULL DEFAULT '0',
	}
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
		theCounter.Id = name
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
