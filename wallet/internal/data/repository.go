package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletRepository interface {
	FindWallet()
	CreateWallet()
	FindTransaction()
	UpdateTransaction(id primitive.ObjectID, updates []primitive.E)
}

type MongoWalletRepository struct {
}

func (m MongoWalletRepository) FindWallet() {
	panic("implement me")
}

func (m MongoWalletRepository) CreateWallet() {
	panic("implement me")
}

func (m MongoWalletRepository) FindTransaction() {
	panic("implement me")
}

func (m MongoWalletRepository) UpdateTransaction(id primitive.ObjectID, updates []primitive.E) {
	panic("implement me")
}
