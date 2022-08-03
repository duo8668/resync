package tests

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetNumSlice(t *testing.T) {
	Convey("given i have num slice", t, func() {

		begin := 1
		for _, val := range GetNumSlice() {
			So(val, ShouldEqual, begin)
			begin++

		}
	})
}
