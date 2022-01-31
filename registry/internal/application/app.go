package application

import (
	"context"
	"eCommerce/registry/docs"
	"eCommerce/registry/internal/api"
	"eCommerce/registry/internal/consumers/topics"
	"eCommerce/registry/internal/core"
	"eCommerce/registry/internal/data"
	"eCommerce/registry/internal/resources"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// App implements graceful shutdown
type App struct {
	cfg       *Config
	ctx       context.Context
	cancelCtx context.CancelFunc
	log       *zap.SugaredLogger
	router    *chi.Router
	resources *resources.Resources

	OrderCoordinator   *core.OrderCoordinator
	PurchaseController *core.Purchaser
	RegistryController *core.RequestRegistry
}

func New() *App {
	logger, _ := zap.NewProduction()

	s := new(App)
	s.log = logger.Sugar()
	s.cfg = s.configFromEnv()
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())

	return s
}

func (a *App) configFromEnv() *Config {
	cfg := new(Config)
	err := envconfig.Process("", cfg)

	if err != nil {
		a.log.Fatalf("can't process the config: %s", err)
	}

	return cfg
}

func (a *App) Build() {
	a.log.Info("Configuring server and initializing resources...")
	a.resources = resources.New(a.ctx, a.log, &resources.ConnectionConfig{
		MongoConnectionUrl: a.cfg.MongoConnectionUrl,
		KafkaConnectionUrl: a.cfg.KafkaConnectionUrl,
	})
	a.resources.InitMongoDB()
	a.resources.InitKafka(&resources.TopicConfig{
		UsersTopic:  topics.WalletCreateTopic,
		OrdersTopic: topics.StorageReserveOrderTopic,
	})

	repository := data.NewMongoRegistryRepository(a.resources.Mongo)
	coordinator := core.NewOrderCoordinator(a.log, repository, a.resources.Producer)
	a.OrderCoordinator = coordinator
	a.PurchaseController = core.NewPurchaser(a.resources, coordinator)
	a.RegistryController = core.NewRequestRegistry(a.log, a.resources)

	a.router = api.NewRouter(a.PurchaseController, a.RegistryController, &api.RouterConfig{
		Host: a.cfg.ApplicationHost,
	})
}

func (a *App) Run() {
	docs.SwaggerInfo.Host = a.cfg.ApplicationHost
	docs.SwaggerInfo.BasePath = "/"

	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Error(err)
		}
	}(a.log)

	a.log.Info("Starting the server...")
	a.OrderCoordinator.Run(kafka.ReaderConfig{
		Brokers:  []string{a.cfg.KafkaConnectionUrl},
		MinBytes: 10e1,
		MaxBytes: 10e3,
		Logger:   log.Default(),
	})

	addr := ":" + strconv.Itoa(a.cfg.ApplicationPort)
	server := &http.Server{Addr: addr, Handler: *a.router}

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancelShutdownCtx := context.WithTimeout(a.ctx, 30*time.Second)
		defer cancelShutdownCtx()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				a.log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		a.resources.Release()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			a.log.Fatal(err)
		}
		a.cancelCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.log.Fatal(err)
	}

	<-a.ctx.Done()
}
