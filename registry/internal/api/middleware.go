package api

import (
	"context"
	"eCommerce/registry/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type RequestRegistryMiddleware struct {
	rh *RegistryHandlers
}

// UserRegistry user registration
func (rr *RequestRegistryMiddleware) UserRegistry(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, `swagger`) {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := rr.rh.RegistryController.RegisterUser(r)
		if err != nil {
			return
		}

		w.Header().Set("X-UID", uid)
		id, _ := primitive.ObjectIDFromHex(uid)

		ctx := context.WithValue(r.Context(), `ident`, models.Identity{Id: id})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// RequestRegistry logging requests
func (rr *RequestRegistryMiddleware) RequestRegistry(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		if strings.Contains(r.URL.Path, `swagger`) {
			return
		}

		err := rr.rh.RegisterRequest(r)
		if err != nil {
			InternalErrorResponse(w, InternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn)
}

// DefaultContentType set content type
func DefaultContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		if !strings.Contains(r.URL.Path, `swagger`) {
			w.Header().Set(`Content-Type`, `application/json`)
		}
	}

	return http.HandlerFunc(fn)
}
