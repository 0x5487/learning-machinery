package main

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/jasonsoft/learning-machinery/helloworld"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
)

func main() {
	// set up log target
	log.
		Str("app_id", "worker").
		SaveToDefault()

	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	var cnf = config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	server.RegisterTask("Say", helloworld.Say)

	worker := server.NewWorker("worker-1", 10)
	err = worker.Launch()
	if err != nil {
		panic(err)
	}

}
