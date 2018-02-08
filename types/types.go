package types

import (
	"fmt"
	"html/template"
)

const (
	ViewMenu   = "templates/menu.html"
	ViewHeader = "templates/header.html"
	// ViewSidebar      = "templates/sidebar.html"
	ViewErr          = "templates/error.html"
	ViewNavbar       = "templates/navbar.html"
	ViewMenuConstant = "templates/navbar_constant.html"
	ViewMenuUs       = "templates/navbarus.html"
	ViewNavButtons   = "templates/navbuttons.html"

	// ViewTextArea     = "templates/textarea_input.html"
	// ViewErr          = "templates/error_display.html"
	// ViewReplyButtons = "templates/replybuttons.html"
	// ViewHeader       = "templates/header.html"
	// ViewFooter       = "templates/footer.html"
	// ViewModalForm    = "templates/modal_form.html"
	KListLimit = 20

	KSortIdAscend  = "-_id"
	KSortIdDescend = "_id"
)

const (
	KFieldNavCount = 1 + iota
	KFieldNavPage
	KFieldNavLink
	KFieldNavNext
)

type (
	HeaderRecord struct {
		Title      string
		Loggedin   bool
		Admin      bool
		Angular    bool
		JSEditor   bool
		FooTable   bool
		DatePicker bool
		DocReady   bool
		NewStyle   bool
		AddUser    bool
		Scripts    []string
		// Load           string
		Load template.HTML
	}

	NavButtonRecord struct {
		HasNav   bool
		HasNext  bool
		HasPrev  bool
		PrevLink string
		NextLink string
	}

	MenuRecord struct {
		Current NavItem
		Items   []NavItem
	}

	NavItem struct {
		Text  string
		Title string
		Link  string
	}

	SortItem struct {
		Title    string
		Link     string
		Sortable bool
	}

	TableRow struct {
		First  string
		Second string
		Third  string
		Fourth string
		Fifth  string
		Sixth  string
	}

	Point struct {
		X float64
		Y float64
	}

	DrawItem struct {
		Strings map[string]string
		Floats  map[string]float64
	}

	LineRecord struct {
		FirstPoint   Point
		SecondPoint  Point
		ThirdPoint   Point
		LastPoint    Point
		ItemType     int
		NumPoints    int
		StrokeFill   int
		ColourNumber int
		RuleWeight   float64
		TextWidth    float64
		TheText      string
		Font         int
		TextAlign    int
		TextStyle    int
		TextSize     int
	}

	RowItem struct {
		Title string
		Value template.HTML
		Link  string
	}

	SearchRecord struct {
		Term string `schema:"search"`
		Type string `schema:"stype"`
	}

	M  map[int]interface{}
	MS map[string]interface{}
	S  []MS
)

var (
	GeneralMenu = []NavItem{
		{"Index", "/", "Go to the Agamik home page"},
		{"Programs", "/barcoder", "Information about our barcode programs, download the latest release"},
		{"Fonts", "/fonts", "Information about barcode fonts, download the latest release"},
		{"Creation", "/create", "We can supply your barcodes as files"},
		{"Downloads", "/download", "Download working versions and demos of our products"},
		{"Types", "/symbols", "Information about barcode types and how to identify different types"},
		{"Explained", "/explain", "Answers to the common questions. Information about barcoding. What types to use for which jobs"},
		{"Buying", "/buying", "Information about buying our products. How to buy and how to pay"},
		{"Contact", "/contact", "Contact information: e-mail, phone, mail, Skype, MSN Messenger and Yahoo Messenger details"},
	}
	AdminMenu = []NavItem{
		{"Cases", "/cases", "Go to the Agamik home page"},
		{"Clients", "/clients", "Information about our barcode programs, download the latest release"},
		{"Fonts", "/fonts", "Information about barcode fonts, download the latest release"},
		{"Creation", "/create", "We can supply your barcodes as files"},
		{"Downloads", "/download", "Download working versions and demos of our products"},
		{"Types", "/symbols", "Information about barcode types and how to identify different types"},
		{"Explained", "/explain", "Answers to the common questions. Information about barcoding. What types to use for which jobs"},
		{"Buying", "/buying", "Information about buying our products. How to buy and how to pay"},
		{"Contact", "/contact", "Contact information: e-mail, phone, mail, Skype, MSN Messenger and Yahoo Messenger details"},
	}
)

func GetUserMenu(isGod bool) []NavItem {
	if isGod {
		return GeneralMenu
	}
	return GeneralMenu
}

func GetGeneralItem(ignoreItem string, inUK bool) []NavItem {
	var (
		theMenu []NavItem
	)
	for _, elem := range GeneralMenu {
		if elem.Text != ignoreItem {
			if elem.Text == "Create" {
				if inUK {
					theMenu = append(theMenu, elem)
				}
			} else {
				theMenu = append(theMenu, elem)
			}
		}
	}
	return theMenu
}

func (nr *NavButtonRecord) SetNavButtons(data M) {
	nr.HasNav = false
	nr.HasNext = false
	nr.HasPrev = false

	// count, ok := data[KFieldNavCount].(int)
	// if !ok {
	// 	return
	// }

	pageNumber, ok := data[KFieldNavPage].(int)
	if !ok {
		return
	}

	link, ok := data[KFieldNavLink].(string)
	if !ok {
		return
	}

	navNext, ok := data[KFieldNavNext].(bool)
	if ok {
		nr.HasNav = true
		nr.HasNext = navNext
		nr.NextLink = fmt.Sprintf("%s/%d", link, pageNumber+1)
	}

	if pageNumber > 0 {
		nr.HasNav = true
		nr.HasPrev = true
		if pageNumber == 1 {
			nr.PrevLink = link
		} else {
			nr.PrevLink = fmt.Sprintf("%s/%d", link, pageNumber-1)
		}
	}
}
