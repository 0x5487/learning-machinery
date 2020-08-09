package main

import (
	"sync"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
)

type Wallet struct {
	mu     sync.Mutex
	Amount int64
}

func (w *Wallet) Withdraw(amount int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Amount -= amount
}

func (w *Wallet) Deposit(amount int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Amount += amount
}

var (
	_wallet *Wallet
)

func main() {
	// set up log target
	log.
		Str("app_id", "wallet").
		SaveToDefault()

	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	// declared
	_wallet = &Wallet{
		Amount: 1000,
	}

	var cnf = config.Config{
		DefaultQueue:  "wallet",
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	server.RegisterTask("withdraw", withdraw)
	server.RegisterTask("deposit", deposit)

	worker := server.NewWorker("wallet-1", 10)
	err = worker.Launch()
	if err != nil {
		panic(err)
	}

}

func withdraw(orderID string, amount int64) error {
	_wallet.Withdraw(amount)
	log.Str("order_id", orderID).Infof("wallet amount withdraw %d, total: %d", amount, _wallet.Amount)
	return nil
}

func deposit(orderID string, amount int64) error {
	_wallet.Deposit(amount)
	log.Str("order_id", orderID).Infof("wallet amount deposit %d, total: %d", amount, _wallet.Amount)
	return nil
}
