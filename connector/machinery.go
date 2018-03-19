package connector

import (
	"log"

	"github.com/jsm/gode/worker/tasknames"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// ConnectMachinery initializes a machinery Server
func ConnectMachinery(broker string, backend string) *machinery.Server {
	cnf := config.Config{
		Broker:          broker,
		DefaultQueue:    "machinery_tasks",
		ResultBackend:   backend,
		ResultsExpireIn: 3600,
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		log.Fatalln(err)
	}

	testMachinery(server)

	return server
}

func testMachinery(server *machinery.Server) {
	signature := tasks.NewSignature(tasknames.Test, nil)
	_, err := server.SendTask(signature)
	if err != nil {
		log.Println("Failed to setup Machinery")
		log.Fatalln(err)
	}
}
