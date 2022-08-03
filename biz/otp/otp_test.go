package otp

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOTP_Expiry(t *testing.T) {
	Convey("given i would like to test on expiry", t, func() {
		bizOtp := NewBizOtp(30)

		// generate 1 second expiry
		otpToken, err := bizOtp.Generate("88887777", 1)
		So(err, ShouldBeNil)
		So(otpToken, ShouldNotBeEmpty)

		// verify
		isOK, err := bizOtp.Verify("88887777", 1, otpToken)
		So(err, ShouldBeNil)
		So(isOK, ShouldBeTrue)

		// sleep for 2s
		time.Sleep(2 * time.Second)

		// should be false now
		isOK, err = bizOtp.Verify("88887777", 1, otpToken)
		So(err, ShouldBeNil)
		So(isOK, ShouldBeFalse)
	})
}
