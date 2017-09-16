package clients

import (
	// "log"
	"strings"

	// "github.com/kataras/iris"
	// "gopkg.in/mgo.v2"

	// "ssafa/crypto"
	"ssafa/db"
)

type (
	Client struct {
		Id          int    `bson:"_id"`
		CaseNumber  string `bson:"case,omitempty"`
		CMSId       string `bson:"cms,omitempty"`
		VisitNumber string `bson:"visit,omitempty"`
		DOB         int    `bson:"dob"`
		Title       string `bson:"title,omitempty"`
		First       string `bson:"first,omitempty"`
		Surname     string `bson:"surname,omitempty"`
		Letters     string `bson:"case,letters"`
		Address     string `bson:"address,omitempty"`
		// `address2` varchar(255) NOT NULL DEFAULT '',
		// `address3` varchar(255) NOT NULL DEFAULT '',
		// `address4` varchar(255) NOT NULL DEFAULT '',
		PostCode string `bson:"postcode,omitempty"`
		Phone    string `bson:"phone,omitempty"`
		Mobile   string `bson:"mobile,omitempty"`
		EMail    string `bson:"email,omitempty"`
		NINum    string `bson:"ninum,omitempty"`
		// `ni2` smallint(2) unsigned zerofill NOT NULL DEFAULT '00',
		// `ni3` smallint(2) unsigned zerofill NOT NULL DEFAULT '00',
		// `ni4` smallint(2) unsigned zerofill NOT NULL DEFAULT '00',
		// `ni5` char(1) DEFAULT NULL,
		ServiceNo string `bson:"serviceno,omitempty"`
		Service   string `bson:"services,omitempty"`
		Annuity   bool   `bson:"annuity,omitempty"`
		Comments  string `bson:"comments,omitempty"`
		Report    string `bson:"report,omitempty"`
		CaseFirst string `bson:"casefirst,omitempty"`
		CaseSurn  string `bson:"casesurn,omitempty"`
		Based     string `bson:"based,omitempty"`
		Closed    bool   `bson:"closed,omitempty"`
		Shredded  bool   `bson:"shredded,omitempty"`
	}

	ClientShow struct {
		Id      int
		Case    string
		First   string
		Surname string
		Service string
	}

	ClientList []ClientShow

	ByClientName struct{ ClientShow }
)

var (
	key []byte
)

func (s ClientList) Len() int {
	return len(s)
}

func (s ClientList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByClientName) Less(i, j int) bool {
	if strings.ToLower(s.ClientShow[i].Surname) != strings.ToLower(s.ClientShow[j].Surname) {
		return strings.ToLower(s.ClientShow[i].Surname) < strings.ToLower(s.ClientShow[j].Surname)
	}
	return s.Organs[i].Weight < s.Organs[j].Weight
}

func GetList(searchCategory, searchTerm string, offset int) ClientList {
	var (
		theClient Client
		theList   ClientList
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	clientColl := session.DB(db.MainDB).C(db.CollectionClients)

	iter := clientColl.Find(nil).Sort("_id").Iter()
	for iter.Next(&theClient) {
		newClient := ClientShow{Id: theClient.Id, Case: theClient.CaseNumber}
		newClient.First = theClient.First
		newClient.Surname = theClient.Surname
		newClient.Service = theClient.Service
		theList = append(theList, newClient)
	}
	iter.Close()

	return theList
}

func SetKey(theKey []byte) {
	key = theKey
}
