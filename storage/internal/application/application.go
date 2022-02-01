package application

import (
	"context"
	"eCommerce/storage/internal/consumers"
	"eCommerce/storage/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg *Config
	ctx context.Context
	log *zap.SugaredLogger

	resources       *StorageResources
	storageConsumer *consumers.StorageConsumer
}

func New() *App {
	logger, _ := zap.NewProduction()

	app := new(App)
	app.log = logger.Sugar()
	app.cfg = NewConfig()
	app.ctx = context.Background()

	return app
}

func (a *App) Build() {
	a.resources = NewStorageResources(a.ctx, a.cfg, a.log).Initialize()
	storage := core.NewStorage(a.log, a.resources.Database, a.resources.Producer)
	a.storageConsumer = consumers.NewStorageConsumer(context.Background(), a.log, a.cfg.KafkaConnectionUrl, storage)
}

func (a *App) InitTestProducts() {
	a.log.Info("generating test products...")

	products := GenerateProducts(5)
	collection := a.resources.Database.Collection(`products`)
	documents, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		a.log.Error(err)
		return
	}

	if documents > 0 {
		a.log.Info("products collection contains documents")
		return
	}

	a.log.Info("inserting test products to the database...")
	_, err = collection.InsertMany(context.Background(), AsInterfaces(products))
	if err != nil {
		a.log.Error(err)
		return
	}

	a.log.Info("successfully complete test products initialization...")
}

func (a *App) Run() {
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Error(err)
		}
	}(a.log)

	a.log.Info("Starting the storage service...")
	a.storageConsumer.Start()

	interrupt := make(chan os.Signal, 1)
	interruptSignals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	signal.Notify(interrupt, interruptSignals...)

	select {
	case x := <-interrupt:
		a.log.Infow("Received a signal.", "signal", x.String())
	}

	a.log.Info("Stopping the app...")

	if err := a.storageConsumer.Stop(); err != nil {
		a.log.Error("Got an error while stopping the storage consumer.", "err", err)
	}

	timeout, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := a.resources.Release(timeout)
	if err != nil {
		a.log.Error("Got an error while releasing resources.", "err", err)
	}

	a.log.Info("The app is calling the last defers and will be stopped.")
}
