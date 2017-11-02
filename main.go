package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"gopkg.in/mgo.v2"

	"ssafa/cases"
	"ssafa/clients"
	"ssafa/cookie"
	"ssafa/crypto"
	"ssafa/db"
	"ssafa/types"
	"ssafa/users"
)

type (
	Config struct {
		Port int
		// Templates string
		// Posts     string
		// Public    string
		// Admin     string
		// Metadata  string
		// Index     string
	}
)

// type (
// 	uUser struct {
// 		Name  string
// 		Towns []string
// 		Num   int
// 	}

// 	uLoginRecord struct {
// 		EMail      string
// 		Remember   bool
// 		Checkfield string
// 	}
// )

var (
	key string

	config = ReadConfig()

	excludeExtensions = [...]string{
		".js",
		".css",
		".jpg",
		".png",
		".ico",
		".svg",
	}
)

func main() {
	// var (
	// 	err error
	// )
	parts := strings.Split(key, "|")
	crypto.SetKey([]byte(parts[0]))
	clients.SetKey([]byte(parts[0]))
	cookie.SetVars(parts[2], parts[3])

	f, _ := os.OpenFile("./logs/ssafa.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)
	log.Println("Application starting")

	db.MongoSession, _ = mgo.Dial(db.DialStr)
	// if err != nil {
	// 	log.Fatal("Bugger Mongo doesn't open", err)
	// }

	defer db.MongoSession.Close()
	db.MongoSession.SetMode(mgo.Monotonic, true)

	// generateKeys()
	// f, _ := os.OpenFile("./logs/ssafa.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// defer f.Close()

	// log.SetOutput(f)
	// log.Println("Application starting")

	app := iris.New() // defaults to these

	r, close := newRequestLogger()
	defer close()

	app.Use(r)
	// customLogger := logger.New(logger.Config{
	// 	// Status displays status code
	// 	Status: true,
	// 	// IP displays request's remote address
	// 	IP: true,
	// 	// Method displays the http method
	// 	Method: true,
	// 	// Path displays the request path
	// 	Path: true,

	// 	//Columns: true,

	// 	// if !empty then its contents derives from `ctx.Values().Get("logger_message")
	// 	// will be added to the logs.
	// 	MessageContextKey: "logger_message",
	// })

	// app.Use(customLogger)
	// - standard html  | iris.HTML(...)
	app.Use(authCheck)
	// - standard html  | iris.HTML(...)
	// - django         | iris.Django(...)
	// - pug(jade)      | iris.Pug(...)
	// - handlebars     | iris.Handlebars(...)
	// - amber          | iris.Amber(...)

	tmpl := iris.HTML("./templates", ".html")
	tmpl.Reload(true) // reload templates on each request (development mode)
	// default template funcs are:
	//
	// - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
	// - {{ render "header.html" }}
	// - {{ render_r "header.html" }} // partial relative path to current page
	// - {{ yield }}
	// - {{ current }}
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})
	app.RegisterView(tmpl)

	app.StaticWeb("/css", "./css")
	app.StaticWeb("/js", "./js")

	app.Get("/", hi)
	app.Post("/", hi)

	users.SetRoutes(app)
	clients.SetRoutes(app)
	cases.SetRoutes(app)

	// http://localhost:9039
	thePort := fmt.Sprintf(":%d", config.Port)
	app.Run(iris.Addr(thePort), iris.WithCharset("UTF-8")) // defaults to that but you can change it.
}

func hi(ctx iris.Context) {
	var (
		header types.HeaderRecord
	)
	// userNumber, _ := ctx.Values().GetInt("user")
	// fmt.Println(ctx.Values())
	// fmt.Println(ctx.Values().Get("logged"))
	// fmt.Println(ctx.Values().Get("user"))
	// fmt.Println(ctx.Values().Get("admin"))
	// fmt.Println(ctx.Values().Get("session"))
	theSession := ctx.Values().Get("session")
	header.Loggedin = theSession.(users.Session).LoggedIn
	header.Admin = theSession.(users.Session).Admin
	// theSession := sessionVal.(users.Session)
	// fmt.Println(header.Loggedin, "vs", userNumber)
	if !header.Loggedin {
		user := users.User{Name: "Albert Einstein", Towns: []string{"London", "Manchester", "Paris"}, Num: 65}
		// fmt.Println("No session set")
		// }
		// if !theSession.LoggedIn {
		header.Title = "RF: Login"
		// details := users.LoginRecord{Username: "mccallum_ir", Password: "bowpa1132"}
		details := users.LoginRecord{}
		ctx.ViewData("User", user)
		ctx.ViewData("Details", details)
		ctx.ViewData("Header", header)
		ctx.View("login.html")
		return
	}

	// fmt.Println("We are logged in")
	header.Title = "RF: Clients"
	clientList, _ := clients.GetList("", "", 0)
	ctx.ViewData("Header", header)
	ctx.ViewData("Details", clientList)
	ctx.View("main.html")
}

func authCheck(ctx iris.Context) {
	var (
		theSession users.Session
	)
	// make auth check
	path := ctx.Path()
	for _, ext := range excludeExtensions {
		if strings.HasSuffix(path, ext) {
			ctx.Next()
			return
		}
	}

	// fmt.Println("In authCheck", path)

	theSession.ValidCookie(ctx.GetCookie("session"))
	ctx.Values().Set("logged", theSession.LoggedIn)
	ctx.Values().Set("admin", theSession.Admin)
	ctx.Values().Set("user", theSession.UserNumber)
	ctx.Values().Set("session", theSession)

	if ctx.Method() == "POST" {
		// fmt.Println("checking a post")
		// nonceString := ctx.FormValue("checkfield")
		// fmt.Println(nonceString)
		// if !crypto.CheckNonce(nonceString) {
		// 	fmt.Println("Failed nonce")
		// }
	}

	ctx.Values().Set("fred", "Hello there")
	ctx.Next()
}

// func newLogFile() *os.File {
// 	filename := "logs/ssafa2.log"
// 	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return f
// }

func newRequestLogger() (h iris.Handler, close func() error) {
	close = func() error { return nil }

	c := logger.Config{
		Status: true,
		IP:     true,
		Method: true,
		Path:   true,
		// Columns: true,
	}

	// logFile := newLogFile()
	// close = func() error {
	// 	err := logFile.Close()
	// 	return err
	// }

	c.LogFunc = func(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}) {
		// output := logger.Columnize(now.Format("2006/01/02 - 15:04:05"), latency, status, ip, method, path, message)
		f := float64(latency)
		f /= 1000000.0
		output := fmt.Sprintf("%s %.3fms %s %s %s %s", now.Format("2006/01/02 - 15:04:05"), f, ip, status, method, path)
		//, latency, status, ip, method, path, message
		// ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())
		log.Println(output)
		// logFile.Write([]byte(output))
	}

	//	we don't want to use the logger
	// to log requests to assets and etc
	c.AddSkipper(func(ctx iris.Context) bool {
		path := ctx.Path()
		for _, ext := range excludeExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		return false
	})

	h = logger.New(c)

	return
}

func ReadConfig() Config {
	var (
		configfile = "ssafa.cfg"
		config     Config
	)
	_, err := os.Stat(configfile)
	if err != nil {
		config.Port = 9039
		log.Println("Config file is missing: ", configfile)
		return config
	}

	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		config.Port = 9039
		log.Println(err)
	}
	return config
}

func generateKeys() {
	theBytes := securecookie.GenerateRandomKey(32)
	if theBytes == nil {
		// fmt.Println("Bugger")
		return
	}
	// fmt.Printf("%x\n", theBytes)
}
