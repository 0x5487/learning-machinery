package main

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
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

	server.RegisterTask("job1", job1)
	server.RegisterTask("job2", job2)

	worker := server.NewWorker("worker-1", 10)
	err = worker.Launch()
	if err != nil {
		panic(err)
	}

}

func job1(name string) (string, error) {
	log.Info("job1 is calling")
	return "Hello " + name + "!", nil
}

func job2(name string) (string, error) {
	return "Hello " + name + "!", nil
}
