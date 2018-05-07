package users_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/users"
	"ssafa/utils"
)

const (
	usersCopyDB    = "userscopy"
	countersCopyDB = "counterscopy"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Suite")
}

func copyUsers() {
	var (
		theUser    db.User
		theCounter db.Counter
	)
	defer db.MongoSession.Close()
	db.MongoSession.SetMode(mgo.Monotonic, true)

	session := db.MongoSession.Copy()
	defer session.Close()
	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
	countersCollection := session.DB(db.MainDB).C(db.CollectionCounters)
	usersCopyCollection := session.DB(db.MainDB).C(usersCopyDB)
	countersCopyCollection := session.DB(db.MainDB).C(countersCopyDB)

	usersCopyCollection.DropCollection()
	countersCopyCollection.DropCollection()

	iter := usersCollection.Find(nil).Iter()
	for iter.Next(&theUser) {
		usersCopyCollection.Insert(&theUser)
	}
	iter.Close()

	iter = countersCollection.Find(nil).Iter()
	for iter.Next(&theCounter) {
		countersCopyCollection.Insert(&theCounter)
	}
	iter.Close()

}

func restoreUsers() {
	var (
		theUser    db.User
		theCounter db.Counter
	)
	defer db.MongoSession.Close()
	db.MongoSession.SetMode(mgo.Monotonic, true)

	session := db.MongoSession.Copy()
	defer session.Close()
	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
	countersCollection := session.DB(db.MainDB).C(db.CollectionCounters)
	usersCopyCollection := session.DB(db.MainDB).C(usersCopyDB)
	countersCopyCollection := session.DB(db.MainDB).C(countersCopyDB)

	usersCollection.DropCollection()
	countersCollection.DropCollection()

	iter := usersCopyCollection.Find(nil).Iter()
	for iter.Next(&theUser) {
		usersCollection.Insert(&theUser)
	}
	iter.Close()

	iter = countersCopyCollection.Find(nil).Iter()
	for iter.Next(&theCounter) {
		countersCollection.Insert(&theCounter)
	}
	iter.Close()
}

var _ = Describe("Password check", func() {
	err := errors.New("test")
	db.MongoSession, err = mgo.Dial(db.DialStr)
	if err != nil {
		fmt.Println("Bugger Mongo doesn't open", err)
	}

	copyUsers()

	Context("fred.smith@gmail.com", func() {
		thePassword := "fred.smith@gmail.com"
		It(thePassword+" is an e-mail", func() {
			Expect(utils.ValidateEmail(thePassword)).Should(BeTrue())
		})
		Specify("the total amount is 0.00", func() {})
	})

	Context("fred.smith@gmail.co.uk", func() {
		thePassword := "fred.smith@gmail.co.uk"
		It(thePassword+" is an e-mail", func() {
			Expect(utils.ValidateEmail(thePassword)).Should(BeTrue())
		})
		Specify("the total amount is 0.00", func() {})
	})

	Context("fred.smith@gmail", func() {
		thePassword := "fred.smith@gmail"
		It(thePassword+" is not an e-mail", func() {
			Expect(utils.ValidateEmail(thePassword)).Should(BeFalse())
		})
		// Specify("the total amount is 0.00", func() {})
	})

	Context("Length check", func() {
		Context("length check 9 ascii", func() {
			thePassword := "Hello all"
			ok, err := users.CheckPassword(thePassword)
			It(thePassword+" not long enough ascii", func() {
				Expect(ok).Should(BeFalse())
				Expect(err).ShouldNot(BeNil())
			})
		})

		Context("length check 10 ascii", func() {
			thePassword := "Hello alls"
			ok, err := users.CheckPassword(thePassword)
			It(thePassword+" long enough ascii", func() {
				Expect(ok).Should(BeTrue())
				Expect(err).Should(BeNil())
			})
		})

		Context("length check 5 non ascii", func() {
			thePassword := "首相も出席"
			ok, err := users.CheckPassword(thePassword)
			It(thePassword+" long enough non ascii", func() {
				Expect(ok).Should(BeFalse())
				Expect(err).ShouldNot(BeNil())
			})
		})

		Context("length check 6 non ascii", func() {
			thePassword := "う首相も出席"
			ok, err := users.CheckPassword(thePassword)
			It(thePassword+" not long enough non ascii", func() {
				Expect(ok).Should(BeTrue())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("Check password has been used before", func() {
		Context("known password", func() {
			thePassword := "password12"
			ok, err := users.CheckPasswordSafe(thePassword)
			It("has the password been used", func() {
				Expect(ok).Should(BeFalse())
				Expect(err).Should(BeNil())
			})
		})
		Context("unknown password", func() {
			thePassword := "there 6547 hdjkhdhd ouepeupoue pouepouepeou"
			ok, err := users.CheckPasswordSafe(thePassword)
			It("has the password been used", func() {
				Expect(ok).Should(BeTrue())
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("Check activation expired time", func() {
		theUser := db.User{}
		theUser.ActivateCode = crypto.RandomChars(users.KActivateLength)
		theUser.ActivateTime = time.Now().Unix() - 100000 + users.KActivateTime
		Context("check time has expired", func() {
			activateTime := time.Now().Unix()
			It("has the password been used", func() {
				Expect(activateTime).Should(BeNumerically(">", theUser.ActivateTime))
			})
		})
	})

	Context("Check activation not expired", func() {
		theUser := db.User{}
		theUser.ActivateCode = crypto.RandomChars(users.KActivateLength)
		theUser.ActivateTime = time.Now().Unix() - 1000 + users.KActivateTime
		Context("check time has not expired", func() {
			activateTime := time.Now().Unix()
			It("has the password been used", func() {
				Expect(activateTime).Should(BeNumerically("<", theUser.ActivateTime))
			})
		})
	})
})
