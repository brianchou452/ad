package api

import (
	"ad/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB interface {
	CreateAd(ad *model.Ad) (*mongo.InsertOneResult, error)
	GetAdByConditions(cond Query) ([]primitive.M, error)
	UpdateCurrentAds(currentCollection int) error
}

type AdminEnv struct {
	DB MongoDB
}
