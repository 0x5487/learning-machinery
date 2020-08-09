package main

import (
	"context"
	"errors"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
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
		DefaultQueue:  "err_topic",
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	// try error
	errTask := tasks.Signature{
		Name: "ThrowErr",
		Args: []tasks.Arg{
			{Type: "string", Value: "Jason"},
		},
		RetryCount: 3,
		//RetryTimeout: 1, // Second
	}

	ctx := context.Background()
	asyncResult, err := server.SendTaskWithContext(ctx, &errTask)
	if err != nil {
		log.Err(err).Warn("send task failed")
		return
	}

	log.Str("taskUUID", asyncResult.Signature.UUID).Debug("wait for result")

	_, err = asyncResult.GetWithTimeout(time.Duration(time.Second*10), time.Duration(time.Millisecond*5))
	if err != nil {
		if errors.Is(err, result.ErrTimeoutReached) {
			log.Str("id", asyncResult.Signature.UUID).Warn("task is timeout")
			return
		}
	}

	state := asyncResult.GetState()
	log.Infof("state: %s", state.State) // FAILURE

	if err != nil {
		// handle error here.
		log.Err(err).Warn("run task failed")
		//return
	}

	// try panic
	panicTask := tasks.Signature{
		Name: "ThrowPanic",
		Args: []tasks.Arg{
			{Type: "string", Value: "Jason"},
		},
	}

	asyncResult, err = server.SendTaskWithContext(ctx, &panicTask)
	if err != nil {
		log.Err(err).Warn("send task failed")
		return
	}
	log.Str("taskUUID", asyncResult.Signature.UUID).Debug("wait for result")

	_, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
	state = asyncResult.GetState()
	log.Infof("state: %s", state.State) // FAILURE

	if err != nil {
		// handle error here.
		log.Err(err).Warn("run task failed")
		//return
	}

}
