package cookie

import (
	"encoding/hex"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"

	"ssafa/crypto"
)

const (
	hkHash = "92534bcaae82b0ceb7962faa18a8f84f3d984785946338502cc44eddc9298ae7a97decfd6aeff743ad2e950d89502704199b121e747f2d97748dc5545fd6dcdb08c2234ba3f3c338d07774cfd60af3d0"
	bkHash = "6401eeb5afcd4f3201aa075d1d9baa190e7c5c4971707d6dc93c68c39643ab6fdd439561eccb602deb2a16a83aae7eb721c6c21cf714c5a05f7a1c9e4edd8e616e01de7fdea9e02ed53e0bcbb666e919"
)

type (
	Cookie struct {
		User             int    `json:"num"`
		Admin            bool   `json:"admin,omitempty"`
		God              bool   `json:"god,omitempty"`
		Auth             []byte `json:"auth,omitempty"`
		ClientSearchTerm string `json:"clst,omitempty"`
		ClientSearchType string `json:"clsc,omitempty"`
		CaseSearchTerm   string `json:"cast,omitempty"`
		CaseSearchType   string `json:"casc,omitempty"`
	}
)

var (
	hashKey  []byte
	blockKey []byte
	sc       *securecookie.SecureCookie

	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func SetVars() {
	hashKey, _ = hex.DecodeString(crypto.Decrypt(hkHash))

	blockKey, _ = hex.DecodeString(crypto.Decrypt(bkHash))
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
