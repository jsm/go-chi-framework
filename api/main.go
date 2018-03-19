package main

import (
	"log"
	"net/http"

	"github.com/jsm/gode/api/application"
	"github.com/jsm/gode/services"
)

var version string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	instance := application.Initialize(version)
	defer instance.Close()
	services.InitializeAll(
		instance,
		instance.Machinery,
		instance.DB,
	)

	instance.StatsD.Increment("startup.api")

	instance.Log.Notice("Running API server")
	http.ListenAndServe(":8080", router())
}
