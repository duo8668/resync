package main

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/duo8668/resync/biz/otp"
	"github.com/gofiber/fiber/v2"
	"github.com/kpango/glg"
)

func main() {

	osSignal := waitCancellationSignal()
	defer signal.Stop(osSignal)

	_, stop := startServer()

	// wait for cancellation signal
	<-osSignal

	// call stop func
	stop()
}

func waitCancellationSignal() (osSignal chan os.Signal) {
	// listen to signal
	osSignal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	return
}

func startServer() (app *fiber.App, stopFunc func()) {

	// create server
	app = createServer()

	// start the server
	go func() {
		glg.Infof("start listen...")
		err := app.Listen(":3000")
		glg.Infof("started, the err:  %s", err.Error())
	}()

	// return create func
	stopFunc = func() {
		glg.Infof("app server is shutting down...")
		err := app.Shutdown()
		glg.Infof("done shutdown, err: %s", err)
	}

	return
}

func createServer() (app *fiber.App) {

	app = fiber.New(fiber.Config{DisableStartupMessage: false})

	//
	otpPeriod := 30
	bizOtp := otp.NewBizOtp(otpPeriod)

	// gen
	app.Get("/helloworld", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	// gen
	app.Post("/totp/gen", func(c *fiber.Ctx) error {
		payload := struct {
			PhoneNum string `json:"PhoneNum"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		payload.PhoneNum = strings.TrimSpace(payload.PhoneNum)
		//
		if payload.PhoneNum == "" {
			return errors.New("phone num cannot be empty")
		}

		otpRes, err := bizOtp.Generate(payload.PhoneNum, otpPeriod)
		if err != nil {
			return err
		}

		return c.SendString(otpRes)
	})

	//check
	app.Post("/totp/verify", func(c *fiber.Ctx) error {

		// read the OTP from the client and trim it
		// payload from POST
		payload := struct {
			PhoneNum string `json:"PhoneNum"`
			Otp      string `json:"otp"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		// phone num
		payload.PhoneNum = strings.TrimSpace(payload.PhoneNum)
		//
		if payload.PhoneNum == "" {
			return errors.New("phone num cannot be empty")
		}

		payload.Otp = strings.TrimSpace(payload.Otp)
		//
		if payload.Otp == "" {
			return errors.New("otp cannot be empty")
		}

		// validate now
		isOk, err := bizOtp.Verify(payload.PhoneNum, otpPeriod, payload.Otp)

		// some error with the OTP, need return err
		if err != nil {
			return err
		}

		// if result is ok
		if isOk {
			return c.SendString("OK")
		}

		// return fail here
		return c.SendString("FAIL")
	})

	return
}
