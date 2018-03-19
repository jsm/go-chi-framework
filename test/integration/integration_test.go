package integration

import (
	"os"
	"testing"

	"github.com/jsm/gode/services"
	"github.com/jsm/gode/test/setup"
)

func TestMain(m *testing.M) {
	setup.Setup()
	code := m.Run()
	setup.Teardown()
	os.Exit(code)
}

func TestRegisterAndLoginSuccess(t *testing.T) {
	email := "testregisterandloginsuccess@email.com"
	password := "password"

	if _, _, err := services.Auth.RegisterEmail(email, password); err != nil {
		t.Error(err)
		return
	}

	if _, _, err := services.Auth.LoginEmail(email, password); err != nil {
		t.Error(err)
		return
	}
}

func TestRegisterAndLoginWrongPassword(t *testing.T) {
	email := "testregisterandloginwrongpassword@email.com"
	password := "password"

	if _, _, err := services.Auth.RegisterEmail(email, password); err != nil {
		t.Error(err)
		return
	}

	if _, _, err := services.Auth.LoginEmail(email, "wrongpassword"); err == nil {
		t.Error("Login should have returned error")
		return
	}
}
