package admin

import (
	// "errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	// "github.com/gorilla/schema"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	// "github.com/globalsign/mgo"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/mail"
	"ssafa/types"
	"ssafa/users"
	"ssafa/utils"
)

func adminMain(ctx iris.Context) {
	var (
		header types.HeaderRecord
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/clients", http.StatusFound)
	}

	header.Title = "RF: Admin"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.View("admin.html")
}

func adminPerson(ctx iris.Context) {
	var (
		header        types.HeaderRecord
		result        db.User
		allPersonnel  []managePerson
		columnHeaders []types.SortItem
		colNames      = []string{"First", "Surname", "Role", "Admin", "Based", "User name", "Active"}
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/clients", http.StatusFound)
	}

	sortNum, err := ctx.Params().GetInt("sortnum")
	if err != nil {
		sortNum = 2
	}

	if sortNum == 0 {
		sortNum = 2
	}

	for i, item := range colNames {
		newHeader := types.SortItem{}
		newHeader.Title = item

		newHeader.Link = fmt.Sprintf("/adminperson/%d", i+1)
		if sortNum == i+1 {
			newHeader.Sortable = false
		} else {
			newHeader.Sortable = true
		}
		columnHeaders = append(columnHeaders, newHeader)
	}
	session := db.MongoSession.Copy()
	defer session.Close()
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	iter := userColl.Find(nil).Sort("db.KFieldUserInactive").Iter()
	if iter.Err() != nil {
		log.Println("Error: adminPerson iter", iter.Err())
	}
	for iter.Next(&result) {
		newPerson := managePerson{}
		newPerson.ID = result.ID
		newPerson.First = crypto.Decrypt(result.First)
		newPerson.Surname = crypto.Decrypt(result.Surname)
		newPerson.Role = result.Position
		newPerson.Based = result.Based
		newPerson.UserName = result.Username
		if result.Admin {
			newPerson.AdminStr = "Yes"
		} else {
			newPerson.AdminStr = "No"
		}
		newPerson.AdminLink = fmt.Sprintf("/adminswap/%d/%d", newPerson.ID, sortNum)
		if result.InActive {
			newPerson.ActiveStr = "No"
		} else {
			newPerson.ActiveStr = "Yes"
		}
		newPerson.ActiveLink = fmt.Sprintf("/activeswap/%d/%d", newPerson.ID, sortNum)
		allPersonnel = append(allPersonnel, newPerson)
	}

	switch sortNum {
	case 1:
		sort.Sort(ByFirst{allPersonnel})
	case 2:
		sort.Sort(BySurname{allPersonnel})
	case 3:
		sort.Sort(ByRole{allPersonnel})
	case 4:
		sort.Sort(ByAdmin{allPersonnel})
	case 5:
		sort.Sort(ByBased{allPersonnel})
	case 6:
		sort.Sort(ByActive{allPersonnel})
	}

	header.Title = "RF: Admin"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("ColHeaders", columnHeaders)
	ctx.ViewData("Details", allPersonnel)
	ctx.View("adminmanage.html")
}

func adminReset(ctx iris.Context) {
	var (
		header        types.HeaderRecord
		result        db.User
		allPersonnel  []managePerson
		columnHeaders []types.SortItem
		colNames      = []string{"First", "Surname", "Code", "E Mail", "Link", "User name"}
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/clients", http.StatusFound)
	}

	sortNum, err := ctx.Params().GetInt("sortnum")
	if err != nil {
		sortNum = 2
	}

	if sortNum == 0 {
		sortNum = 2
	}

	for i, item := range colNames {
		newHeader := types.SortItem{}
		newHeader.Title = item

		newHeader.Link = fmt.Sprintf("/adminperson/%d", i+1)
		if sortNum == i+1 {
			newHeader.Sortable = false
		} else {
			newHeader.Sortable = true
		}
		columnHeaders = append(columnHeaders, newHeader)
	}
	session := db.MongoSession.Copy()
	defer session.Close()
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	query := bson.M{db.KFieldUserActivateCode: bson.M{"$exists": 1}}
	iter := userColl.Find(query).Iter()
	if iter.Err() != nil {
		log.Println("Error: adminPerson iter", iter.Err())
	}
	for iter.Next(&result) {
		if result.InActive {
			continue
		}
		newPerson := managePerson{}
		newPerson.ID = result.ID
		newPerson.First = crypto.Decrypt(result.First)
		newPerson.Surname = crypto.Decrypt(result.Surname)
		newPerson.Role = result.ActivateCode
		newPerson.Based = crypto.Decrypt(result.EMail)
		newPerson.UserName = result.Username
		newPerson.AdminLink = fmt.Sprintf("/adminswap/%d/%d", newPerson.ID, sortNum)
		if result.InActive {
			newPerson.ActiveStr = "No"
		} else {
			newPerson.ActiveStr = "Yes"
		}
		newPerson.ActiveLink = fmt.Sprintf("/activeswap/%d/%d", newPerson.ID, sortNum)
		allPersonnel = append(allPersonnel, newPerson)
	}

	switch sortNum {
	case 1:
		sort.Sort(ByFirst{allPersonnel})
	case 2:
		sort.Sort(BySurname{allPersonnel})
	case 3:
		sort.Sort(ByRole{allPersonnel})
	case 4:
		sort.Sort(ByAdmin{allPersonnel})
	case 5:
		sort.Sort(ByBased{allPersonnel})
	case 6:
		sort.Sort(ByActive{allPersonnel})
	}

	header.Title = "RF: Admin"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("ColHeaders", columnHeaders)
	ctx.ViewData("Details", allPersonnel)
	ctx.View("adminreset.html")
}

func adminSwap(ctx iris.Context) {
	var (
		theUser db.User
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/clients", http.StatusFound)
	}

	sortNum, err := ctx.Params().GetInt("sortnum")
	if err != nil {
		ctx.Redirect("/adminperson", http.StatusFound)
	}

	theUrl := fmt.Sprintf("/adminperson/%d", sortNum)

	userNum, err := ctx.Params().GetInt("usernum")
	if err != nil {
		ctx.Redirect(theUrl, http.StatusFound)
	}

	session := db.MongoSession.Copy()
	defer session.Close()
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	err = userColl.FindId(userNum).One(&theUser)
	if err != nil {
		ctx.Redirect(theUrl, http.StatusFound)
	}

	set := bson.M{}
	if theUser.Admin {
		set[db.KFieldUserAdmin] = false
	} else {
		set[db.KFieldUserAdmin] = true
	}
	userColl.UpdateId(userNum, bson.M{"$set": set})
	ctx.Redirect(theUrl, http.StatusFound)
}

func activeSwap(ctx iris.Context) {
	var (
		theUser db.User
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/clients", http.StatusFound)
	}

	sortNum, err := ctx.Params().GetInt("sortnum")
	if err != nil {
		ctx.Redirect("/adminperson", http.StatusFound)
	}

	theUrl := fmt.Sprintf("/adminperson/%d", sortNum)

	userNum, err := ctx.Params().GetInt("usernum")
	if err != nil {
		ctx.Redirect(theUrl, http.StatusFound)
	}

	session := db.MongoSession.Copy()
	defer session.Close()
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	err = userColl.FindId(userNum).One(&theUser)
	if err != nil {
		ctx.Redirect(theUrl, http.StatusFound)
	}

	set := bson.M{}
	if theUser.InActive {
		set[db.KFieldUserInactive] = false
	} else {
		set[db.KFieldUserInactive] = true
	}
	userColl.UpdateId(userNum, bson.M{"$set": set})
	ctx.Redirect(theUrl, http.StatusFound)
}

func adminAddPerson(ctx iris.Context) {
	var (
		details      personAdd
		header       types.HeaderRecord
		errorMessage string
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/clients", http.StatusFound)
	}

	switch ctx.Method() {
	case http.MethodPost:
		err := decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: Decode adminadd:", err)
		}

		err = details.check()
		if err == nil {
			details.save()
			ctx.Redirect("/adminperson", http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Title = "RF: Admin"
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("adminaddperson.html")
}

func (pa *personAdd) check() error {
	var (
		theUser db.User
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	if pa.AdminStr == "yes" {
		pa.Admin = true
	} else {
		pa.Admin = false
	}

	if len(pa.First) == 0 {
		return errorNoFirst
	}

	if len(pa.Surname) == 0 {
		return errorNoSurname
	}

	if len(pa.UserName) == 0 {
		return errorNoUsername
	}

	if len(pa.UserName) < 4 {
		return errorShortUsername
	}

	if len(pa.EMail) == 0 {
		return errorNoEMail
	}

	if !utils.ValidateEmail(pa.EMail) {
		return errorBadEMail
	}

	err := userColl.Find(bson.M{db.KFieldUserUserName: pa.UserName}).One(&theUser)
	if err == nil {
		return errorUsernameUsed
	}

	return nil
}

func (pa *personAdd) save() {
	var (
		theUser db.User
		tempStr string
	)

	session := db.MongoSession.Copy()
	defer session.Close()
	userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	theUser.First = crypto.Encrypt(pa.First)
	theUser.Surname = crypto.Encrypt(pa.Surname)
	if len(pa.First) > 0 {
		if len(pa.Surname) > 0 {
			tempStr = pa.First + " " + pa.Surname
		} else {
			tempStr = pa.First
		}
	} else {
		tempStr = pa.Surname
	}
	theUser.Name = crypto.Encrypt(tempStr)
	theUser.Username = pa.UserName
	theUser.Address = crypto.Encrypt(pa.Address)
	theUser.PostCode = crypto.Encrypt(pa.Postcode)
	theUser.Phone = crypto.Encrypt(pa.Phone)
	theUser.Mobile = crypto.Encrypt(pa.Mobile)
	theUser.EMail = crypto.Encrypt(pa.EMail)
	theUser.Position = pa.Role
	theUser.Based = pa.Based
	theUser.ActivateCode = crypto.RandomLower(users.KActivateLength)
	theUser.ActivateTime = time.Now().Unix() + users.KActivateTime

	theUser.ID = db.GetNextSequence(db.CollectionUsers)
	userColl.Insert(&theUser)

	mail.SendActivate(pa.EMail, theUser.ActivateCode)
}

func (s manageList) Len() int {
	return len(s)
}

func (s manageList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByFirst) Less(i, j int) bool {
	name1 := strings.ToLower(s.manageList[i].First)
	name2 := strings.ToLower(s.manageList[j].First)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].Surname)
	name2 = strings.ToLower(s.manageList[j].Surname)
	return name1 < name2
}

func (s BySurname) Less(i, j int) bool {
	name1 := strings.ToLower(s.manageList[i].Surname)
	name2 := strings.ToLower(s.manageList[j].Surname)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].First)
	name2 = strings.ToLower(s.manageList[j].First)
	return name1 < name2
}

func (s ByRole) Less(i, j int) bool {
	name1 := strings.ToLower(s.manageList[i].Role)
	name2 := strings.ToLower(s.manageList[j].Role)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].Surname)
	name2 = strings.ToLower(s.manageList[j].Surname)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].First)
	name2 = strings.ToLower(s.manageList[j].First)
	return name1 < name2
}

func (s ByAdmin) Less(i, j int) bool {
	name1 := strings.ToLower(s.manageList[i].AdminStr)
	name2 := strings.ToLower(s.manageList[j].AdminStr)
	if name1 != name2 {
		return name1 > name2
	}
	name1 = strings.ToLower(s.manageList[i].Surname)
	name2 = strings.ToLower(s.manageList[j].Surname)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].First)
	name2 = strings.ToLower(s.manageList[j].First)
	return name1 < name2
}

func (s ByBased) Less(i, j int) bool {
	name1 := strings.ToLower(s.manageList[i].Based)
	name2 := strings.ToLower(s.manageList[j].Based)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].Surname)
	name2 = strings.ToLower(s.manageList[j].Surname)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].First)
	name2 = strings.ToLower(s.manageList[j].First)
	return name1 < name2
}

func (s ByActive) Less(i, j int) bool {
	name1 := strings.ToLower(s.manageList[i].ActiveStr)
	name2 := strings.ToLower(s.manageList[j].ActiveStr)
	if name1 != name2 {
		return name1 > name2
	}
	name1 = strings.ToLower(s.manageList[i].Surname)
	name2 = strings.ToLower(s.manageList[j].Surname)
	if name1 != name2 {
		return name1 < name2
	}
	name1 = strings.ToLower(s.manageList[i].First)
	name2 = strings.ToLower(s.manageList[j].First)
	return name1 < name2
}
