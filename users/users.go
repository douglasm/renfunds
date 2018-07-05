package users

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	// "github.com/globalsign/mgo"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
	"ssafa/utils"
)

type (
	activateRec struct {
		ID         int    `schema:"id"`
		Code       string `schema:"code"`
		Username   string `schema:"username"`
		Password1  string `schema:"pass1"`
		Password2  string `schema:"pass2"`
		Checkfield string `schema:"checkfield"`
		Commit     string `schema:"commit"`
	}

	editRec struct {
		ID          int           `schema:"id"`
		First       string        `schema:"first"`
		Surname     string        `schema:"surname"`
		Name        string        `schema:"_"`
		Position    string        `schema:"posit"`
		Based       string        `schema:"based"`
		EMail       string        `schema:"email"`
		Address     string        `schema:"address"`
		AddressHTML template.HTML `schema:"-"`
		PostCode    string        `schema:"postcode"`
		Phone       string        `schema:"phone"`
		Mobile      string        `schema:"mobile"`
		Admin       bool          `schema:"-"`
		AdminStr    string        `schema:"admin"`
		Commit      string        `schema:"commit"`
	}
)

func GetUserName(userNumber int) string {
	var (
		theUser db.User
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	err := usersCollection.FindId(userNumber).One(&theUser)
	if err != nil {
		log.Println("Error: Reading user name error:", err)
		return "Unknown personnel"
	}
	return crypto.Decrypt(theUser.First) + " " + crypto.Decrypt(theUser.Surname)
}

func getUserList(ctx iris.Context) {
	var (
		ur       userRequest
		userList []userChoice
		theUser  db.User
	)
	decoder.Decode(&ur, ctx.FormValues())

	searchTerm := strings.ToLower(ur.SearchTerm)

	// theCookie := ctx.Values().Get("cookie")

	session := db.MongoSession.Copy()
	defer session.Close()

	userCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	count := 0
	// iter := userCollection.Find(bson.M{db.KFieldUserInactive: false}).Iter()
	iter := userCollection.Find(nil).Iter()
	for iter.Next(&theUser) {
		if count > 10 {
			break
		}
		if theUser.InActive {
			continue
		}

		if !strings.Contains(strings.ToLower(theUser.Position), "caseworker") {
			continue
		}

		theName := crypto.Decrypt(theUser.First) + " " + crypto.Decrypt(theUser.Surname)
		if strings.Contains(strings.ToLower(theName), searchTerm) {
			newUser := userChoice{ID: theUser.ID, Name: theName}
			userList = append(userList, newUser)
			count++
		}
	}
	iter.Close()

	ctx.JSON(userList)
}

func editUser(ctx iris.Context) {
	var (
		theUser      db.User
		header       types.HeaderRecord
		details      editRec
		userNum      int
		errorMessage string
		err          error
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	if !theSession.(Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	userNum, err = ctx.Params().GetInt("usernum")
	if err != nil {
		ctx.Redirect("/adminperson", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	switch ctx.Method() {
	case http.MethodGet:
		err = usersCollection.FindId(userNum).One(&theUser)
		if err != nil {
			log.Println("Error: editUser find error", err)
			ctx.Redirect("/adminperson", http.StatusFound)
			return
		}

		details.ID = userNum
		details.First = crypto.Decrypt(theUser.First)
		details.Surname = crypto.Decrypt(theUser.Surname)
		details.Position = theUser.Position
		details.Based = theUser.Based
		details.EMail = crypto.Decrypt(theUser.EMail)
		details.Address = crypto.Decrypt(theUser.Address)
		details.PostCode = crypto.Decrypt(theUser.PostCode)
		details.Phone = crypto.Decrypt(theUser.Phone)
		details.Mobile = crypto.Decrypt(theUser.Mobile)
		details.Admin = theUser.Admin

	case http.MethodPost:
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode user edit", err)
		}
		details.save()
		tempStr := fmt.Sprintf("/user/%d", userNum)
		ctx.Redirect(tempStr, http.StatusFound)
		return
	}

	details.Name = strings.TrimSpace(details.First + " " + details.Surname)

	header.Title = "Edit user " + details.Name
	header.Loggedin = theSession.(Session).LoggedIn
	header.Admin = theSession.(Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("useredit.html")
}

func showUser(ctx iris.Context) {
	var (
		theUser db.User
		header  types.HeaderRecord
		details editRec
		userNum int
		err     error
	)
	theSession := ctx.Values().Get("session")
	if !theSession.(Session).LoggedIn {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	if !theSession.(Session).Admin {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	userNum, err = ctx.Params().GetInt("usernum")
	if err != nil {
		ctx.Redirect("/adminperson", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	err = usersCollection.FindId(userNum).One(&theUser)
	if err != nil {
		log.Println("Error: editUser find error", err)
		ctx.Redirect("/adminperson", http.StatusFound)
		return
	}

	details.ID = userNum
	details.First = crypto.Decrypt(theUser.First)
	details.Surname = crypto.Decrypt(theUser.Surname)
	details.Position = theUser.Position
	details.Based = theUser.Based
	details.EMail = crypto.Decrypt(theUser.EMail)
	details.Address = crypto.Decrypt(theUser.Address)
	details.AddressHTML = template.HTML(strings.Replace(details.Address, "\r", "<br />", -1))
	details.PostCode = crypto.Decrypt(theUser.PostCode)
	details.Phone = crypto.Decrypt(theUser.Phone)
	details.Mobile = crypto.Decrypt(theUser.Mobile)
	details.Admin = theUser.Admin

	details.Name = strings.TrimSpace(details.First + " " + details.Surname)

	header.Title = "User " + details.Name
	header.Loggedin = theSession.(Session).LoggedIn
	header.Admin = theSession.(Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("usershow.html")
}

func activateUser(ctx iris.Context) {
	var (
		theUser      db.User
		header       types.HeaderRecord
		details      activateRec
		errorMessage string
		ok           bool
		err          error
	)

	theCode := ctx.Params().Get("code")

	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)

	theTemplate := "activatecode.html"
	switch ctx.Method() {
	case http.MethodGet:
		if len(theCode) != 0 {
			err = usersCollection.Find(bson.M{db.KFieldUserActivateCode: theCode}).One(&theUser)
			if err == nil {
				theTemplate = "activate.html"
				details.Code = theCode
				details.ID = theUser.ID
			}
		}

	case http.MethodPost:
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode user activate", err)
		}

		err = usersCollection.Find(bson.M{db.KFieldUserActivateCode: details.Code}).One(&theUser)
		if err != nil {
			errorMessage = "The code and user name do not match"
			break
		}
		if theUser.Username != details.Username {
			errorMessage = "The code and user name do not match"
			break
		}

		ok, err = CheckPassword(details.Password1, details.Password2)
		if err == nil {
			if ok {
				theSalt := crypto.RandomChars(20)
				sets := bson.M{}
				unSets := bson.M{}
				sets[db.KFieldUserInactive] = false
				sets[db.KFieldUserSalt] = theSalt
				sets[db.KFieldUserPassword], err = crypto.GetHash(details.Password1, theSalt)

				unSets[db.KFieldUserActivateCode] = 1
				unSets[db.KFieldUserActivateTime] = 1

				err = usersCollection.UpdateId(theUser.ID, bson.M{"$set": sets, "$unset": unSets})
				if err != nil {
					log.Println("Error: users activate update", err)
				}

				ctx.Redirect("/login", http.StatusFound)
				return
			}
		} else {
			errorMessage = err.Error()
		}
		if len(theCode) > 0 {
			theTemplate = "activate.html"
		}
	}

	header.Title = "RF: Activate"
	header.Loggedin = false
	header.Scripts = append(header.Scripts, "passwordtoggle")

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View(theTemplate)
}

func (er editRec) save() {
	var (
		theUser  db.User
		tempBool bool
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	usersCollection := session.DB(db.MainDB).C(db.CollectionUsers)
	err := usersCollection.FindId(er.ID).One(&theUser)

	sets := bson.M{}

	if crypto.Decrypt(theUser.First) != er.First {
		sets[db.KFieldUserFirst] = crypto.Encrypt(er.First)
	}

	if crypto.Decrypt(theUser.Surname) != er.Surname {
		sets[db.KFieldUserSurname] = crypto.Encrypt(er.Surname)
	}

	if crypto.Decrypt(theUser.Address) != er.Address {
		sets[db.KFieldUserAddress] = crypto.Encrypt(er.Address)
	}

	if crypto.Decrypt(theUser.PostCode) != er.PostCode {
		sets[db.KFieldUserPostCode] = crypto.Encrypt(er.PostCode)
	}

	if crypto.Decrypt(theUser.EMail) != er.EMail {
		sets[db.KFieldUserEMail] = crypto.Encrypt(er.EMail)
	}

	if crypto.Decrypt(theUser.Phone) != er.Phone {
		sets[db.KFieldUserPhone] = crypto.Encrypt(er.Phone)
	}

	if crypto.Decrypt(theUser.Mobile) != er.Mobile {
		sets[db.KFieldUserMobile] = crypto.Encrypt(er.Mobile)
	}

	if theUser.Position != er.Position {
		sets[db.KFieldUserPosition] = er.Position
	}

	if theUser.Based != er.Based {
		sets[db.KFieldUserBased] = er.Based
	}

	if er.AdminStr == "yes" {
		tempBool = true
	}

	sets[db.KFieldUserAdmin] = tempBool

	if len(sets) == 0 {
		return
	}

	err = usersCollection.UpdateId(er.ID, bson.M{"$set": sets})
	if err != nil {
		log.Println("Error: update user", err)
	}
	// KFieldUserBased = "based"
	// KFieldUserArea = "area"

}

func CheckPassword(password, password2 string) (bool, error) {
	if password != password2 {
		return false, errPassMismatch
	}

	if utils.ValidateEmail(password) {
		return false, errEMailUsed
	}

	if !checkPasswordLen(password) {
		return false, errPassShort
	}
	return CheckPasswordSafe(password)
}

func CheckPasswordSafe(password string) (bool, error) {
	var (
		resp *http.Response
	)
	h := sha1.New()
	io.WriteString(h, password)
	theStr := fmt.Sprintf("%X", h.Sum(nil))

	theUrl := "https://api.pwnedpasswords.com/range/" + theStr[:5]
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodGet, theUrl, nil)
	if err != nil {
		log.Println("NewRequest Error:", err)
		return false, err
	}

	for i := 0; i < 3; i++ {
		resp, err = client.Do(r)
		if err == nil {
			break
		}
		i++
	}
	if err != nil {
		log.Println("Error: Pass check client", err)
		return false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), theStr[5:]) {
		return false, errPassUsed
	}

	return true, nil
}

func checkPasswordLen(password string) bool {
	var (
		ascii   int
		unicode int
	)
	for _, ch := range password {
		if int(ch) > 255 {
			unicode++
		} else {
			ascii++
		}
	}
	if unicode >= 6 {
		return true
	}
	if ascii >= 10 {
		return true
	}
	switch unicode {
	case 1:
		if ascii >= 8 {
			return true
		}
	case 2:
		if ascii >= 6 {
			return true
		}
	case 3:
		if ascii >= 4 {
			return true
		}
	case 4:
		if ascii >= 3 {
			return true
		}
	case 5:
		if ascii >= 2 {
			return true
		}
	}
	return false
}
