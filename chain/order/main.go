package main

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/google/uuid"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
)

func main() {
	// set up log target
	log.
		Str("app_id", "order").
		SaveToDefault()

	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	var cnf = config.Config{
		DefaultQueue:  "order",
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	server.RegisterTask("createOrder", createdOrder)

	worker := server.NewWorker("order-1", 10)
	err = worker.Launch()
	if err != nil {
		panic(err)
	}

}

func createdOrder() (string, error) {
	orderID := uuid.New().String()
	log.Infof("order %s was created", orderID)
	return orderID, nil
}
