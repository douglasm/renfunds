package clients

import (
	"testing"

	// "github.com/gorilla/schema"

	// "github.com/kataras/iris"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Dates(t *testing.T) {
	t.Parallel()

	Convey("Testing Date bad format", t, func() {
		ce := ClientEdit{First: "Alan", Surname: "Bennett"}
		ce.DOB = "hdhddh"
		So(ce.checkClient(), ShouldEqual, ErrorDateBadFormat)

		ce.DOB = "uu/uu/uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateBadDay)

		ce.DOB = "12/uu/uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateBadMonth)

		ce.DOB = "12/12/uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateBadYear)

		ce.DOB = "0/uu/uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateLowDay)

		ce.DOB = "1/0/uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateLowMonth)

		ce.DOB = "1/13/uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateHighMonth)

		ce.DOB = "1/1/0"
		So(ce.checkClient(), ShouldEqual, ErrorDateLowYear)

		ce.DOB = "1//uu"
		So(ce.checkClient(), ShouldEqual, ErrorDateBadMonth)

		ce.DOB = "1/1/"
		So(ce.checkClient(), ShouldEqual, ErrorDateBadYear)

		ce.DOB = "1/1/1900"
		So(ce.checkClient(), ShouldEqual, ErrorDateLowYear)

		ce.DOB = "1/1/2100"
		So(ce.checkClient(), ShouldEqual, ErrorDateHighYear)

		ce.DOB = "32/1/1945"
		So(ce.checkClient(), ShouldEqual, ErrorDateHighDay)

		ce.DOB = "31/4/1945"
		So(ce.checkClient(), ShouldEqual, ErrorDateHighDay)

		ce.DOB = "29/2/1945"
		So(ce.checkClient(), ShouldEqual, ErrorDateHighDay)

		ce.DOB = "30/2/1944"
		So(ce.checkClient(), ShouldEqual, ErrorDateHighDay)

		ce.DOB = "28/2/1944"
		So(ce.checkClient(), ShouldBeNil)
	})
}

func Test_NINumbers(t *testing.T) {
	t.Parallel()

	Convey("Testing NINum bad format", t, func() {
		ce := ClientEdit{First: "Alan", Surname: "Bennett"}

		ce.NINum = "hdhddh"
		So(ce.checkClient(), ShouldEqual, ErrorNINumWrongLength)

		ce.NINum = "1dhddhqwe"
		So(ce.checkClient(), ShouldEqual, ErrorNINumBadFormat)

		ce.NINum = "d1hddhqwe"
		So(ce.checkClient(), ShouldEqual, ErrorNINumBadFormat)

		ce.NINum = "dd1234564"
		So(ce.checkClient(), ShouldEqual, ErrorNINumBadFormat)

		ce.NINum = "dd12345xx"
		So(ce.checkClient(), ShouldEqual, ErrorNINumBadFormat)

		ce.NINum = "dd123x56x"
		So(ce.checkClient(), ShouldEqual, ErrorNINumBadFormat)

		ce.NINum = "we123456x"
		So(ce.checkClient(), ShouldBeNil)
	})
}

// ErrorDateBadDay       = errors.New("The day value is wrong")
// ErrorDateBadMonth     = errors.New("The month value is wrong")
// ErrorDateBadYear      = errors.New("The year value is wrong")
// ErrorDateLowDay       = errors.New("The day is too low")
// ErrorDateLowMonth     = errors.New("The month is too low")
// ErrorDateLowYear      = errors.New("The aeary is too low")
