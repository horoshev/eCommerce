package application

import (
	"context"
	"eCommerce/registry/docs"
	"eCommerce/registry/internal/api"
	"eCommerce/registry/internal/consumers"
	"eCommerce/registry/internal/core"
	"eCommerce/registry/internal/data"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const ServiceName = `registry`

// App is used for configuration and running an application. Implements graceful shutdown.
type App struct {
	cfg       *Config
	ctx       context.Context
	cancelCtx context.CancelFunc
	log       *zap.SugaredLogger
	router    *chi.Router
	resources *RegistryResources

	OrderCoordinator   *core.OrderCoordinator
	PurchaseController *core.Purchaser
	RegistryController *core.RequestRegistry
}

func New() *App {
	logger, _ := zap.NewProduction()

	s := new(App)
	s.log = logger.Sugar()
	s.cfg = NewConfig()
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())

	return s
}

func (a *App) Build() {
	a.log.Info("Configuring server and initializing resources...")
	a.resources = NewRegistryResources(a.ctx, a.log, a.cfg)

	repository := data.NewMongoRegistryRepository(a.resources.Database)
	coordinator := core.NewOrderCoordinator(a.log, repository, a.resources.Producer)
	a.OrderCoordinator = coordinator
	a.PurchaseController = core.NewPurchaser(a.log, a.resources.Database, a.resources.Producer, coordinator)
	a.RegistryController = core.NewRequestRegistry(a.log, a.resources.Database, a.resources.Producer)

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
	a.OrderCoordinator.Run(consumers.BrokerReaderConfig(a.cfg.KafkaConnectionUrl))

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
		a.resources.Release(shutdownCtx)
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
