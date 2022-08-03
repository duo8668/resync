package otp

import (
	"encoding/base32"

	"github.com/jltorresm/otpgo"
)

type BizOtp struct {
	OtpTool otpgo.TOTP
	enc     *base32.Encoding
}

// initiate a new struct with a default period
func NewBizOtp(period int) *BizOtp {

	return &BizOtp{
		OtpTool: otpgo.TOTP{Period: period},
		enc:     base32.StdEncoding.WithPadding(base32.NoPadding),
	}
}

// generate token
func (o *BizOtp) Generate(identifier string, period int) (string, error) {

	totp := otpgo.TOTP{Key: o.enc.EncodeToString([]byte(identifier)), Period: period}

	return totp.Generate()
}

// verify if otp is valid
func (o *BizOtp) Verify(identifier string, period int, otp string) (bool, error) {

	totp := otpgo.TOTP{Key: o.enc.EncodeToString([]byte(identifier)), Period: period}

	return totp.Validate(otp)
}
