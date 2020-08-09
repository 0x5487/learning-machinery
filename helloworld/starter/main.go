package main

import (
	"context"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
)

func main() {
	// set up log target
	log.
		Str("app_id", "starter").
		SaveToDefault()

	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)
	defer log.Flush() // flush log buffer

	var cnf = config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	sayTask := tasks.Signature{
		Name: "Say",
		Args: []tasks.Arg{
			{Type: "string", Value: "Jason"},
		},
	}

	ctx := context.Background()
	asyncResult, err := server.SendTaskWithContext(ctx, &sayTask)
	if err != nil {
		panic(err)
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		log.Infof("result: %s", result.String())
	}
	//log.Infof("%s\n", tasks.HumanReadableResults(results))
}
