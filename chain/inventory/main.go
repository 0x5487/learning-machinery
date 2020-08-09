package main

import (
	"sync"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
)

type Stock struct {
	mu     sync.Mutex
	Amount int64
}

func (s *Stock) Decrease() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Amount--
}

func (s *Stock) Increase() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Amount++
}

var (
	_stock *Stock
)

func main() {
	// set up log target
	log.
		Str("app_id", "inventory").
		SaveToDefault()

	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	// declared
	_stock = &Stock{
		Amount: 1000,
	}

	var cnf = config.Config{
		DefaultQueue:  "inventory",
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	server.RegisterTask("decreaseStock", decreaseStock)
	server.RegisterTask("increaseStock", increaseStock)

	worker := server.NewWorker("inventory-1", 10)
	err = worker.Launch()
	if err != nil {
		panic(err)
	}

}

func decreaseStock(orderID string) error {
	_stock.Decrease()
	logger := log.Str("order_id", orderID)
	logger.Info("decreaseStock is calling")

	time.Sleep(10 * time.Second)

	logger.Infof("stock amount decrease 1, total: %d", _stock.Amount)
	return nil
}

func increaseStock(orderID string) error {
	_stock.Increase()
	log.Str("order_id", orderID).Infof("stock amount increase 1, total: %d", _stock.Amount)
	return nil
}
