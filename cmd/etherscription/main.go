package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ziollek/etherscription/internal/storage/memory"
	"github.com/ziollek/etherscription/pkg/api"
	"github.com/ziollek/etherscription/pkg/config"
	"github.com/ziollek/etherscription/pkg/etherum"
	"github.com/ziollek/etherscription/pkg/logging"
	"github.com/ziollek/etherscription/pkg/model"
	"github.com/ziollek/etherscription/pkg/parser"
)

const txBufferSize = 1000

var (
	node            string
	cfgPath         string
	port            int
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func init() {
	flag.StringVar(&node, "node", "", "Ethereum node URL")
	flag.StringVar(&cfgPath, "config", "configuration.yaml", "Ethereum node URL")
	flag.IntVar(&port, "port", 8888, "HTTP server port")
	flag.Parse()
}

func main() {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		logging.Logger().Err(err).Msg("Error while loading configuration")
		os.Exit(1)
	}
	ctx, done := signal.NotifyContext(context.Background(), shutdownSignals...)
	defer done()
	txChan := make(chan model.Transaction, txBufferSize)
	blocksChan := make(chan int)
	subscribersStorage := memory.NewKVStorage[string]()
	transactionsStorage := memory.NewListStorage[model.Transaction]()
	stateStorage := memory.NewKVStorage[int]()
	broker := parser.NewBroker(
		txChan,
		blocksChan,
		parser.NewConsumerService(transactionsStorage, subscribersStorage, cfg.Storage),
		parser.NewStateConsumerService(stateStorage),
	)
	go broker.Start(ctx)
	fetcher := etherum.NewFetcher(
		cfg.RPC.Interval,
		etherum.NewRPCClient(node, cfg.RPC),
		txChan,
		blocksChan,
	)
	go func(ctx context.Context) {
		err := fetcher.Start(ctx)
		if err != nil {
			logging.Logger().Err(err).Msg("Error while starting fetcher")
			panic(err)
		}
	}(ctx)
	cleaner := memory.NewCleaner[model.Transaction](transactionsStorage, cfg.Storage.CleanInterval)
	go cleaner.Start(ctx)

	subscriptionService := parser.NewSubscriptionService(transactionsStorage, subscribersStorage, stateStorage)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: api.ConfigureRouting(api.NewHandler(subscriptionService)),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logging.Logger().Err(err).Msg("Error while starting HTTP server")
		}
	}()
	<-ctx.Done()
	shutdownCtx, shutdownDone := context.WithTimeout(context.Background(), time.Second)
	defer shutdownDone()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logging.Logger().Err(err).Msg("Error while shutting down HTTP server")
	}
}
