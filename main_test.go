package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kpango/glg"
	. "github.com/smartystreets/goconvey/convey"
)

var fiberClient = fiber.Client{}

const host = "http://127.0.0.1:3000/"

func init() {

}

func TestServerReply(t *testing.T) {

	// start server
	go func() {
		startServer()
	}()

	// let the server fully start
	time.Sleep(100 * time.Millisecond)

	//
	glg.Get().SetLineTraceMode(glg.TraceLineShort)

	Convey("given i have server started", t, func() {

		// just simple test for API can be called
		Convey("the server should be able to call directly", func() {
			agent := fiberClient.Get(fmt.Sprintf("%s%s", host, "helloworld"))

			code, body, errs := agent.String()

			So(code, ShouldEqual, 200)
			So(body, ShouldEqual, "Hello World!")
			So(errs, ShouldBeEmpty)
		})

		//
		Convey("i should be able to get one otp using my phone num", func() {
			//
			phoneNum := "87870001"
			// should be able to get token
			genArgs := fiber.AcquireArgs()
			genArgs.Set("PhoneNum", phoneNum)
			genAgent := fiberClient.Post(fmt.Sprintf("%s%s", host, "/totp/gen")).Form(genArgs)

			code, body, errs := genAgent.String()

			So(code, ShouldEqual, 200)
			So(len(body), ShouldEqual, 6)
			So(errs, ShouldBeEmpty)

			myCode := body

			Convey("i should be able to verify the received otp using my phone num", func() {
				// should be true for the token
				args := genArgs
				args.Set("otp", myCode)
				verifyAgent := fiberClient.Post(fmt.Sprintf("%s%s", host, "/totp/verify")).Form(args)

				code, body, errs = verifyAgent.String()
				So(code, ShouldEqual, 200)
				So(body, ShouldEqual, "OK")
				So(errs, ShouldBeEmpty)
			})

		})

	})

}
