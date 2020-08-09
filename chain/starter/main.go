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

	createdOrderTask := tasks.Signature{
		Name:       "createOrder",
		RoutingKey: "order",
	}

	decreaseStockTask := tasks.Signature{
		Name:       "decreaseStock",
		RoutingKey: "inventory",
	}

	// withdrawTask := tasks.Signature{
	// 	Name:       "withdraw",
	// 	RoutingKey: "wallet",
	// }

	ctx := context.Background()
	chain, _ := tasks.NewChain(&createdOrderTask, &decreaseStockTask)
	chainAsyncResult, err := server.SendChainWithContext(ctx, chain)
	if err != nil {
		// failed to send the chain
		// do something with the error
		log.Err(err).Warn("send chain failed")
		return
	}

	results, err := chainAsyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		log.Err(err).Warn("run chain failed")
		return
	}
	log.Infof("chain result: %s", tasks.HumanReadableResults(results))
}
