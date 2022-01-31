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

	// TODO:
	// 1. Consumer for new orders
	// 2. Consumer for cancel order
	// 3. Producer for reservation result
	// 4. Database
	// 5. Kafka
	// 6. Create test products

	a.resources = NewStorageResources(a.ctx, a.cfg, a.log)
	a.resources.Initialize()

	storage := core.NewStorage(a.log, a.resources.Database, a.resources.Producer)
	a.storageConsumer = consumers.NewStorageConsumer(context.Background(), a.log, a.cfg.ToConsumerConfig(), storage)
}

func (a *App) InitTestProducts() {
	products := GenerateProducts(5)
	collection := a.resources.Database.Collection(`products`)

	documents, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		a.log.Error(err)
	}

	if documents > 0 {
		a.log.Info("products collection contains documents")
		return
	}

	_, err = collection.InsertMany(context.TODO(), AsInterfaces(products))
	if err != nil {
		a.log.Error(err)
	}
}

func (a *App) Run() {

	// TODO: Update Graceful shutdown
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
		a.log.Error("Got an error while stopping the business logic server.", "err", err)
	}

	timeout, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	err := a.resources.Release(timeout)
	if err != nil {
		a.log.Error("Got an error while releasing resources.", "err", err)
	}

	a.log.Info("The app is calling the last defers and will be stopped.")
}
