package auth

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

func TestRegisterEmail(t *testing.T) {
	if _, _, err := services.Auth.RegisterEmail("testregister@email.com", "password"); err != nil {
		t.Error(err)
		return
	}
}

func TestLoginUnregisteredEmail(t *testing.T) {
	if _, _, err := services.Auth.LoginEmail("testloginunregistered@email.com", "password"); err == nil {
		t.Error("Register should have returned an error")
		return
	}
}
