package application

import (
	"context"
	"eCommerce/wallet/internal/consumers"
	"eCommerce/wallet/internal/core"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg       *Config
	ctx       context.Context
	cancelCtx context.CancelFunc
	log       *zap.SugaredLogger

	resource      *WalletResources
	userConsumer  *consumers.UserConsumer
	orderConsumer *consumers.OrderConsumer
}

func New() *App {
	logger, _ := zap.NewProduction()

	app := new(App)
	app.cfg = NewConfig()
	app.log = logger.Sugar()
	app.ctx, app.cancelCtx = context.WithCancel(context.Background())

	return app
}

func (a *App) Build() {
	a.resource = NewWalletResources(a.ctx, a.log, a.cfg).Initialize()
	controller := core.NewWalletController(a.ctx, a.log, a.resource.Database, a.resource.KafkaProducer)
	a.userConsumer = consumers.NewUserConsumer(a.ctx, a.log, a.cfg.KafkaConnectionUrl, controller)
	a.orderConsumer = consumers.NewOrderConsumer(a.ctx, a.log, a.cfg.KafkaConnectionUrl, controller)
}

func (a *App) Run() {
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Error(err)
		}
	}(a.log)

	a.log.Info("Starting the wallet service...")
	a.userConsumer.Start()
	a.orderConsumer.Start()

	interrupt := make(chan os.Signal, 1)
	interruptSignals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	signal.Notify(interrupt, interruptSignals...)

	select {
	case x := <-interrupt:
		a.log.Infow("Received a signal.", "signal", x.String())
	}

	a.log.Info("Stopping the app...")

	if err := a.userConsumer.Stop(); err != nil {
		a.log.Error("Got an error while stopping the business logic server.", "err", err)
	}

	if err := a.orderConsumer.Stop(); err != nil {
		a.log.Error("Got an error while stopping the business logic server.", "err", err)
	}

	timeout, _ := context.WithTimeout(context.Background(), 10*time.Second)
	a.resource.Release(timeout)

	a.log.Info("The app is calling the last defers and will be stopped.")
}
