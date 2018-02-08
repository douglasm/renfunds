package admin

import (
	// "errors"
	// "html/template"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	// "github.com/gorilla/schema"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2"

	"ssafa/crypto"
	"ssafa/db"
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
		colNames      = []string{"First", "Surname", "Role", "Admin", "Based", "Active"}
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
		newPerson.Id = result.Id
		newPerson.First = crypto.Decrypt(result.First)
		newPerson.Surname = crypto.Decrypt(result.Surname)
		newPerson.Role = result.Position
		newPerson.Based = result.Based
		if result.Admin {
			newPerson.AdminStr = "Yes"
		} else {
			newPerson.AdminStr = "No"
		}
		newPerson.AdminLink = fmt.Sprintf("/adminswap/%d/%d", newPerson.Id, sortNum)
		if result.InActive {
			newPerson.ActiveStr = "No"
		} else {
			newPerson.ActiveStr = "Yes"
		}
		newPerson.ActiveLink = fmt.Sprintf("/activeswap/%d/%d", newPerson.Id, sortNum)
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
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("ColHeaders", columnHeaders)
	ctx.ViewData("Details", allPersonnel)
	ctx.View("adminmanage.html")
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
			fmt.Println("Error: Decode adminadd:", err)
		}
		fmt.Printf("%+v\n", details)
		err = details.check()
		if err == nil {
			ctx.Redirect("/adminPerson", http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Title = "RF: Admin"
	header.Admin = theSession.(users.Session).Admin

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("adminaddperson.html")
}

func (pa *personAdd) check() error {
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

	return nil
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
