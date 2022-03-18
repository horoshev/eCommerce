package core

import (
	"bytes"
	"context"
	"eCommerce/registry/internal/api/requests"
	"eCommerce/registry/internal/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type RegistryController interface {
	RegisterUser(r *http.Request) (string, error)
	RegisterRequest(r *http.Request) error
	ListRequests(r *requests.PageRequest) ([]models.UserRequest, error)
	ListUserRequests(identity *models.Identity, r *requests.PageRequest) ([]models.UserRequest, error)
}

type RequestRegistry struct {
	log      *zap.SugaredLogger
	Requests *mongo.Collection
	Users    *mongo.Collection
	Producer *kafka.Writer
}

func NewRequestRegistry(log *zap.SugaredLogger, db *mongo.Database, producer *kafka.Writer) *RequestRegistry {
	p := new(RequestRegistry)
	p.log = log
	p.Producer = producer
	p.Users = db.Collection("users")
	p.Requests = db.Collection("requests")

	return p
}

func (rr *RequestRegistry) RegisterUser(r *http.Request) (string, error) {
	identity, err := Identity(r)

	if err == nil {
		query := bson.M{"_id": identity.Id}
		result := rr.Users.FindOne(context.Background(), query)

		if result.Err() == nil {
			return identity.Id.Hex(), nil
		}
	}

	user, err := rr.CreateUser()
	if err != nil {
		return "", err
	}

	return user.Id.Hex(), nil
}

func (rr *RequestRegistry) CreateUser() (*models.User, error) {
	user := models.User{}
	result, err := rr.Users.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)

	// Create event
	go func(rr RequestRegistry, u models.User) {
		payload, err := json.Marshal(u)
		if err != nil {
			rr.log.Error(err)
			return
		}

		err = rr.Producer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(u.Id.Hex()),
			Value: payload,
			Topic: models.WalletCreateTopic,
		})
		if err != nil {
			rr.log.Error(err)
		}
	}(*rr, user)

	return &user, nil
}

func (rr *RequestRegistry) RegisterRequest(r *http.Request) error {
	identity, err := Identity(r)
	if err != nil {
		return err
	}

	request := models.UserRequest{
		UserId:    identity.Id,
		Type:      r.Method,
		Path:      r.URL.Path,
		Timestamp: time.Now().UTC(),
	}

	_, err = rr.Requests.InsertOne(context.Background(), request)
	if err != nil {
		return err
	}

	return nil
}

func (rr *RequestRegistry) ListRequests(r *requests.PageRequest) ([]models.UserRequest, error) {
	opt := options.Find()
	opt.SetSort(bson.D{{"_id", -1}})
	opt.SetSkip(int64(r.Page * r.Size))
	opt.SetLimit(int64(r.Size))

	records, err := rr.Requests.Find(context.Background(), bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	list := make([]models.UserRequest, 0, r.Size)
	for records.Next(context.Background()) {
		var record models.UserRequest
		if err = records.Decode(&record); err != nil {
			return nil, err
		}
		list = append(list, record)
	}

	return list, nil
}

func (rr *RequestRegistry) ListUserRequests(identity *models.Identity, r *requests.PageRequest) ([]models.UserRequest, error) {
	opt := options.Find()
	opt.SetSort(bson.D{{"_id", -1}})
	opt.SetSkip(int64(r.Page * r.Size))
	opt.SetLimit(int64(r.Size))
	filter := bson.D{{"user_id", identity.Id}}

	records, err := rr.Requests.Find(context.Background(), filter, opt)
	if err != nil {
		return nil, err
	}

	list := make([]models.UserRequest, 0, r.Size)
	for records.Next(context.Background()) {
		var record models.UserRequest
		if err = records.Decode(&record); err != nil {
			return nil, err
		}
		list = append(list, record)
	}

	return list, nil
}

func Identity(r *http.Request) (*models.Identity, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if identity, err := IdentityFromParams(r); err == nil {
		return identity, nil
	}

	if identity, err := IdentityFromBody(body); err == nil {
		return identity, nil
	}

	if identity, err := IdentityFromCookies(r); err == nil {
		return identity, nil
	}

	if identity, err := IdentityFromContext(r.Context()); err == nil {
		return identity, nil
	}

	// А можешь пояснить, зачем столько вариантов?
	return nil, errors.New(`no identity`)
}

func NewIdentity(uid string) (*models.Identity, error) {
	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return nil, errors.New(`bad id format`)
	}

	return &models.Identity{
		Id: id,
	}, nil
}

func IdentityFromBody(body []byte) (*models.Identity, error) {
	identity := new(models.Identity)
	if err := json.Unmarshal(body, identity); err != nil {
		return nil, err
	}

	return identity, nil
}

func IdentityFromCookies(r *http.Request) (*models.Identity, error) {
	cookie, err := r.Cookie(`uid`)
	if err != nil {
		return nil, err
	}

	return NewIdentity(cookie.Value)
}

func IdentityFromContext(ctx context.Context) (*models.Identity, error) {
	if identity, ok := ctx.Value(`ident`).(models.Identity); ok {
		return &identity, nil
	}

	return nil, errors.New(`no identity in context`)
}

func IdentityFromParams(r *http.Request) (*models.Identity, error) {
	param := r.URL.Query().Get(`uid`)
	if param == "" {
		return nil, errors.New(`no uid param`)
	}

	return NewIdentity(param)
}
