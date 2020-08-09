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

	ctx := context.Background()

	// job1
	job1Task := tasks.Signature{
		UUID: "task_e6e300f5-edc6-4865-81a7-e0c71ce7ae2e",
		Name: "job1",
		Args: []tasks.Arg{
			{Type: "string", Value: "job1"},
		},
		Immutable: true,
	}

	asyncResult, err := server.SendTaskWithContext(ctx, &job1Task)
	if err != nil {
		panic(err)
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		log.Infof("job1 result: %s", result.String())
	}

	// job2
	job2Task := tasks.Signature{
		Name: "job2",
		Args: []tasks.Arg{
			{Type: "string", Value: "job2"},
		},
	}

	asyncResult, err = server.SendTaskWithContext(ctx, &job2Task)
	if err != nil {
		panic(err)
	}

	results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		log.Infof("job2 result: %s", result.String())
	}

}
