package resources

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"

	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
	"ssafa/users"
	"ssafa/utils"
)

func listResources(ctx iris.Context) {
	type (
		resourceItem struct {
			ID      int    `bson:"_id"`
			Name    string `bson:"name,omitempty"`
			Phone   string `bson:"phone,omitempty"`
			URL     string `bson:"url,omitempty"`
			Email   string `bson:"email,omitempty"`
			Contact string `bson:"contact,omitempty"`
			Lower   string `bson:"lower,omitempty"`
		}
	)
	var (
		pageNum     int
		theResource resourceItem
		details     []listItem
		header      types.HeaderRecord
		// result      bson.M
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
	header.Title = "RF: Resources"

	theSkip := pageNum * types.KListLimit

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)

	project := bson.M{"$project": bson.M{"name": 1, "contact": 1, "email": 1, "phone": 1, "lower": bson.M{"$toLower": "$name"}}}
	skip := bson.M{"$skip": theSkip}
	limit := bson.M{"$limit": types.KListLimit + 1}
	sort := bson.M{"$sort": bson.M{"lower": 1}}

	pipeline := []bson.M{project, sort, skip, limit}

	iter := resourceColl.Pipe(pipeline).Iter()
	for iter.Next(&theResource) {
		if len(details) < types.KListLimit {
			newResource := listItem{}
			newResource.ID = theResource.ID
			newResource.Name = theResource.Name
			newResource.Contact = crypto.Decrypt(theResource.Contact)
			newResource.Phone = crypto.Decrypt(theResource.Phone)
			newResource.EMail = crypto.Decrypt(theResource.Email)
			newResource.EMail = crypto.Decrypt(theResource.Email)
			details = append(details, newResource)
		} else {
			navButtons.HasNav = true
			navButtons.HasNext = true
			navButtons.NextLink = fmt.Sprintf("resources/%d", pageNum+1)
		}
	}
	iter.Close()

	if pageNum > 0 {
		navButtons.HasNav = true
		navButtons.HasPrev = true
		if pageNum > 1 {
			navButtons.PrevLink = fmt.Sprintf("resources/%d", pageNum-1)
		} else {
			navButtons.PrevLink = "resources"
		}
	}
	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("NavButtons", navButtons)
	ctx.View("resources/list.html")
}

func showResource(ctx iris.Context) {
	var (
		resourceNum int
		theResource db.Resource
		details     resourceShow
		header      types.HeaderRecord
		err         error
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

	resourceNum, err = ctx.Params().GetInt("resourcenum")
	if err != nil {
		resourceNum = 0
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)
	err = resourceColl.FindId(resourceNum).One(&theResource)
	if err != nil {
		log.Println("Error: show resource", resourceNum, err)
		ctx.Redirect("/resources", http.StatusNotFound)
		return
	}

	details.ID = theResource.ID
	details.Name = theResource.Name
	details.Contact = crypto.Decrypt(theResource.Contact)
	details.Phone = crypto.Decrypt(theResource.Phone)
	details.EMail = crypto.Decrypt(theResource.EMail)
	if len(theResource.URL) > 0 {
		details.URL = theResource.URL
		details.Link = utils.ValidateURL(theResource.URL)
	}

	details.Comments = utils.EscapeString(theResource.Comments)
	details.Address = utils.EscapeString(theResource.Address)
	if theResource.Updated > 0 {
		details.Updated = utils.GetDateAndTime(theResource.Updated, false)
	}

	header.Title = "RF: Resource " + details.Name

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("resources/show.html")
}

func addResource(ctx iris.Context) {
	var (
		header       types.HeaderRecord
		details      editResource
		errorMessage string
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

	session := db.MongoSession.Copy()
	defer session.Close()

	// resourceColl := session.DB(db.MainDB).C(db.CollectionResources)

	switch ctx.Method() {
	case http.MethodGet:

	case http.MethodPost:
		ctx.FormValues()
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode add voucher", err)
		}
		err = details.checkEdit()
		if err == nil {
			details.saveAdd()
			theURL := fmt.Sprintf("/resource/%d", details.ID)
			ctx.Redirect(theURL, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: Add resource"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("resources/add.html")

}

func resourceEdit(ctx iris.Context) {
	var (
		header       types.HeaderRecord
		details      editResource
		theResource  db.Resource
		errorMessage string
		resourceNum  int
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
	resourceNum, err = ctx.Params().GetInt("resourcenum")
	if err != nil {
		ctx.Redirect("/resources", http.StatusNotFound)
		return
	}

	details.ID = resourceNum

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)

	resourceColl.FindId(resourceNum).One(&theResource)

	switch ctx.Method() {
	case http.MethodGet:
		details.Name = theResource.Name
		details.Contact = crypto.Decrypt(theResource.Contact)
		details.Phone = crypto.Decrypt(theResource.Phone)
		details.EMail = crypto.Decrypt(theResource.EMail)
		details.URL = theResource.URL
		details.Comment = theResource.Comments
		details.Address = theResource.Address

	case http.MethodPost:
		ctx.FormValues()
		err = decoder.Decode(&details, ctx.FormValues())
		if err != nil {
			log.Println("Error: decode add resource", err)
		}
		err = details.checkEdit()
		if err == nil {
			details.saveEdit(&theResource)
			theURL := fmt.Sprintf("/resource/%d", resourceNum)
			ctx.Redirect(theURL, http.StatusFound)
			return
		}
		errorMessage = err.Error()
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	header.Title = "RF: Edit resource"

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.ViewData("ErrorMessage", errorMessage)
	ctx.View("resources/edit.html")
}

func resourceDelete(ctx iris.Context) {
	var (
		resourceNum int
		theResource db.Resource
		details     resourceShow
		header      types.HeaderRecord
		err         error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	resourceNum, err = ctx.Params().GetInt("resourcenum")
	if err != nil {
		resourceNum = 0
	}

	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)
	err = resourceColl.FindId(resourceNum).One(&theResource)
	if err != nil {
		log.Println("Error: show resource", resourceNum, err)
		ctx.Redirect("/resources", http.StatusNotFound)
		return
	}

	details.ID = theResource.ID
	details.Name = theResource.Name

	header.Title = "RF: Delete resource " + details.Name

	ctx.ViewData("Header", header)
	ctx.ViewData("Details", details)
	ctx.View("resources/delete.html")
}

func resourceRemove(ctx iris.Context) {
	var (
		resourceNum int
		err         error
	)

	theSession := ctx.Values().Get("session")
	if !theSession.(users.Session).LoggedIn {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}
	if !theSession.(users.Session).Admin {
		ctx.Redirect("/resources", http.StatusNotFound)
		return
	}

	resourceNum, err = ctx.Params().GetInt("resourcenum")
	if err != nil {
		ctx.Redirect("/resources", http.StatusFound)
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)
	err = resourceColl.RemoveId(resourceNum)
	if err != nil {
		log.Println("Error: remove resource", resourceNum, err)
	}

	ctx.Redirect("/resources", http.StatusFound)
}

func getResources(pageNum int) ([]listItem, bool) {
	var (
		theResource db.Resource
		// theUser    db.User
		theList []listItem
		// err     error
	)
	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)
	// userColl := session.DB(db.MainDB).C(db.CollectionUsers)

	skip := pageNum * types.KListLimit
	limit := types.KListLimit + 1

	// found := 0
	iter := resourceColl.Find(nil).Skip(skip).Limit(limit).Sort(db.KFieldClientsOrder).Iter()
	for iter.Next(&theResource) {
		if len(theList) < types.KListLimit {
			newResource := listItem{ID: theResource.ID}
			newResource.Name = theResource.Name
			// err = userColl.FindId(theResource.UserIssuing).One(&theUser)
			// if err == nil {
			// 	newVoucher.
			// }

			theList = append(theList, newResource)
		} else {
			return theList, true
		}
	}
	iter.Close()

	return theList, false

}

func (cv *createResource) checkAdd() error {
	if len(cv.Name) == 0 {
		return errNoName
	}
	return nil
}

func (cv *createResource) saveAdd(userNumber int) {
	/*	var (
		theResource db.Resource
	)*/

	/*	session := db.MongoSession.Copy()
		defer session.Close()

		resourceColl := session.DB(db.MainDB).C(db.CollectionResources)*/

	/*	err := resourceColl.Insert(&theResource)
		if err != nil {
			log.Println("Error: Voucher insert", err)
		}*/
}

func (er *editResource) checkEdit() error {
	/*	var (
		err error
	)*/
	if len(er.Name) == 0 {
		return errNoName
	}
	return nil
}

func (ev *editResource) saveEdit(theResource *db.Resource) {
	sets := bson.M{}
	unsets := bson.M{}
	if theResource.Name != ev.Name {
		if len(ev.Name) > 0 {
			sets[db.FieldResourceName] = ev.Name
		} else {
			unsets[db.FieldResourceName] = 1
		}
	}

	if crypto.Decrypt(theResource.Contact) != ev.Contact {
		if len(ev.Contact) > 0 {
			sets[db.FieldResourceContact] = crypto.Encrypt(ev.Contact)
		} else {
			unsets[db.FieldResourceContact] = 1
		}
	}

	if crypto.Decrypt(theResource.Phone) != ev.Phone {
		if len(ev.Phone) > 0 {
			sets[db.FieldResourcePhone] = crypto.Encrypt(ev.Phone)
		} else {
			unsets[db.FieldResourcePhone] = 1
		}
	}

	if crypto.Decrypt(theResource.EMail) != ev.EMail {
		if len(ev.EMail) > 0 {
			sets[db.FieldResourceEmail] = crypto.Encrypt(ev.EMail)
		} else {
			unsets[db.FieldResourceEmail] = 1
		}
	}

	if theResource.URL != ev.URL {
		if len(ev.URL) > 0 {
			sets[db.FieldResourceURL] = ev.URL
		} else {
			unsets[db.FieldResourceURL] = 1
		}
	}

	if theResource.Address != ev.Address {
		if len(ev.Address) > 0 {
			sets[db.FieldResourceAddress] = ev.Address
		} else {
			unsets[db.FieldResourceAddress] = 1
		}
	}

	if theResource.Comments != ev.Comment {
		if len(ev.Comment) > 0 {
			sets[db.FieldResourceComments] = ev.Comment
		} else {
			unsets[db.FieldResourceComments] = 1
		}
	}

	gotOne := false
	if len(sets) > 0 {
		sets[db.FieldResourceUpdated] = time.Now().Unix()
		gotOne = true
	}
	if len(unsets) > 0 {
		sets[db.FieldResourceUpdated] = time.Now().Unix()
		gotOne = true
	}

	if !gotOne {
		return
	}

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)

	theUpdate := bson.M{}
	if len(sets) > 0 {
		theUpdate["$set"] = sets
	}
	if len(unsets) > 0 {
		theUpdate["$unset"] = unsets
	}

	err := resourceColl.UpdateId(theResource.ID, theUpdate)
	if err != nil {
		log.Println("Error: update resource", err)
	}
}

func (ev *editResource) saveAdd() {
	var (
		theResource db.Resource
		err         error
	)
	ev.ID = db.GetNextSequence(db.CollectionResources)
	theResource.ID = ev.ID
	theResource.Name = ev.Name
	theResource.Contact = crypto.Encrypt(ev.Contact)
	theResource.Phone = crypto.Encrypt(ev.Phone)
	theResource.EMail = crypto.Encrypt(ev.EMail)
	theResource.URL = ev.URL
	theResource.Address = ev.Address
	theResource.Comments = ev.Comment

	session := db.MongoSession.Copy()
	defer session.Close()

	resourceColl := session.DB(db.MainDB).C(db.CollectionResources)

	err = resourceColl.Insert(&theResource)
	if err != nil {
		log.Println("Error: add resource", err)
	}
}
