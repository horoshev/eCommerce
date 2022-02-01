package api

import (
	"eCommerce/registry/internal/core"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	swag "github.com/swaggo/http-swagger"
	"time"
)

type RouterConfig struct {
	Host string
}

func NewRouter(pc core.PurchaseController, rc core.RegistryController, cfg *RouterConfig) *chi.Router {
	ph := OrderHandlers{pc}
	rh := RegistryHandlers{rc}
	rr := RequestRegistryMiddleware{&rh}

	swagIndex := fmt.Sprintf("http://%s/swagger/index.html", cfg.Host)
	swagDoc := fmt.Sprintf("http://%s/swagger/doc.json", cfg.Host)

	var r chi.Router = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(DefaultContentType)
	r.Use(rr.UserRegistry)
	r.Use(rr.RequestRegistry)

	r.Get("/", Index(swagIndex))
	r.Get("/requests", rh.ListRequestsHandler)
	r.Get("/requests/{id}", rh.ListUserRequestsHandler)
	r.Get("/orders", ph.ListOrdersHandler)
	r.Get("/orders/{id}", ph.ListUserOrdersHandler)
	r.Post("/order", ph.OrderHandler)

	r.Get("/swagger/*", swag.Handler(swag.URL(swagDoc)))

	return &r
}
