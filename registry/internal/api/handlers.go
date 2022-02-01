package api

import (
	"eCommerce/registry/internal/api/requests"
	"eCommerce/registry/internal/core"
	"eCommerce/registry/internal/models"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type OrderHandlers struct {
	core.PurchaseController
}

type RegistryHandlers struct {
	core.RegistryController
}

// Index godoc
// @Summary 	Welcome message from microservice.
// @Description	Returns name, version and link to the swagger.
// @Tags        general
// @Accept      json
// @Produce     json
// @Param   	uid	query	string	false	"auth User ID"
// @Success 	200 {object} ServerResponse
// @Router 		/ [get]
func Index(swagURI string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		OkResponse(w, ServerResponse{
			Name:    `registry`,
			Version: `1.0.0`,
			Swagger: swagURI,
		})
	}
}

// OrderHandler godoc
// @Summary 	Creates new order.
// @Description	Processing user order request. When `uid` param missed - creates new user and wallet as well.
// @Description	Parameter `uid` is stored in response headers under `x-uid` name.
// @Tags        orders
// @Accept      json
// @Produce     json
// @Param   	uid		query	string					false	"auth User ID"
// @Param   	order	body	requests.OrderRequest	true	"Order"
// @Success 	200 {object} models.Order
// @Success 	400 {object} api.Response
// @Failure 	500 {object} api.Response
// @Router 		/order [post]
func (c *OrderHandlers) OrderHandler(w http.ResponseWriter, r *http.Request) {
	req := new(requests.OrderRequest)
	err := json.NewDecoder(r.Body).Decode(req)

	identity, err := core.Identity(r)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	result, err := c.PurchaseController.Order(identity.Id, req)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	OkResponse(w, result)
}

// ListOrdersHandler godoc
// @Summary 	Returns list of created orders for all users.
// @Description Find and return created orders of the user using paging.
// @Tags        orders
// @Accept      json
// @Produce     json
// @Param   	uid  query string false "auth User ID"
// @Param   	page query string false "Page number"
// @Param   	size query string false "Page size"
// @Success 	200 {object} []models.Order
// @Failure 	500 {object} api.Response
// @Router 		/orders [get]
func (c *OrderHandlers) ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	request := requests.ParsePageRequest(r)
	result, err := c.PurchaseController.ListOrders(request)

	if err != nil {
		BadRequestResponse(w, RequestBodyParseError)
		return
	}

	OkResponse(w, result)
}

// ListUserOrdersHandler godoc
// @Summary 	Returns list of created orders for user.
// @Description Find and return created orders of the user using paging.
// @Tags        orders
// @Accept      json
// @Produce     json
// @Param   	id	path string true "User ID to filter orders"
// @Param   	uid  query string false "auth User ID"
// @Success 	200 {object} []models.Order
// @Failure 	500 {object} api.Response
// @Router 		/orders/{id} [get]
func (c *OrderHandlers) ListUserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	request := requests.ParsePageRequest(r)
	id, err := GetPathParameter(r.URL.Path, 1)
	if err != nil {
		BadRequestResponse(w, BadPathParameterError)
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		BadRequestResponse(w, BadPathParameterError)
	}

	result, err := c.PurchaseController.ListUserOrders(userId, request)
	if err != nil {
		BadRequestResponse(w, RequestBodyParseError)
		return
	}

	OkResponse(w, result)
}

// ListRequestsHandler godoc
// @Summary 	Returns list of created requests for all users.
// @Description Find and return created requests of the users using paging.
// @Tags        requests
// @Accept      json
// @Produce     json
// @Param   	uid	 query	string false "auth User ID"
// @Param   	page query string false "Page number"
// @Param   	size query string false "Page size"
// @Success 	200 {object} []models.UserRequest
// @Failure 	500 {object} api.Response
// @Router 		/requests [get]
func (c *RegistryHandlers) ListRequestsHandler(w http.ResponseWriter, r *http.Request) {
	request := requests.ParsePageRequest(r)
	result, err := c.RegistryController.ListRequests(request)

	if err != nil {
		BadRequestResponse(w, RequestBodyParseError)
		return
	}

	OkResponse(w, result)
}

// ListUserRequestsHandler godoc
// @Summary 	Returns list of created requests for user.
// @Description Find and return created requests of the user using paging.
// @Tags        requests
// @Accept      json
// @Produce     json
// @Param   	id   path 	string true  "User ID to filter requests"
// @Param   	uid	 query	string false "auth User ID"
// @Param   	page query 	string false "Page number"
// @Param   	size query 	string false "Page size"
// @Success 	200 {object} []models.UserRequest
// @Failure 	500 {object} api.Response
// @Router 		/requests/{id} [get]
func (c *RegistryHandlers) ListUserRequestsHandler(w http.ResponseWriter, r *http.Request) {
	request := requests.ParsePageRequest(r)
	id, err := GetPathParameter(r.URL.Path, 1)
	if err != nil {
		BadRequestResponse(w, BadPathParameterError)
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		BadRequestResponse(w, BadPathParameterError)
	}

	if err != nil {
		UnauthorizedResponse(w, RequestBodyParseError)
		return
	}

	result, err := c.RegistryController.ListUserRequests(&models.Identity{Id: userId}, request)
	if err != nil {
		BadRequestResponse(w, RequestBodyParseError)
		return
	}

	OkResponse(w, result)
}

func GetPathParameter(path string, index int) (string, error) {
	path = strings.Trim(path, "/")
	split := strings.Split(path, "/")

	if len(split) < index {
		return "", errors.New(`out of range`)
	}

	return split[index], nil
}
