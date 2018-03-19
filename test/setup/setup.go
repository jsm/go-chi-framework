package setup

import (
	"log"

	"github.com/jsm/gode/api/application"
	"github.com/jsm/gode/services"
)

func Setup() {
	log.Println("Test Setup")
	instance := application.Initialize("test")
	services.InitializeAll(
		instance,
		instance.Machinery,
		instance.DB,
	)
}

func Teardown() {
	application.Instance.Close()
}
