package cookie

import (
	"encoding/hex"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type (
	Cookie struct {
		User  int    `json:"num"`
		Admin bool   `json:"admin"`
		God   bool   `json:"god"`
		Auth  []byte `json:"auth"`
	}
)

var (
	hashKey  []byte
	blockKey []byte
	sc       *securecookie.SecureCookie

	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func SetVars(hk, bk string) {
	hashKey, _ = hex.DecodeString(hk)

	blockKey, _ = hex.DecodeString(bk)
	sc = securecookie.New(hashKey, blockKey)
}

func MakeCookie(ctx iris.Context) {
	session := sess.Start(ctx)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Set("authenticated", true)
	session.Set("user", 21)
}
