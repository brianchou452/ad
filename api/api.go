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

type RedisStore interface {
	AddAdToCache(adId string, data *model.Ad) error
	GetAdFromCache(adId string) (model.Ad, error)
	AddAdToZSet(condition string, conditionContent string, adId string, endAt float64) error
	GetAdIdFromCondition(cond Query) ([]model.Ad, error)
}

type Env struct {
	DB    MongoDB
	Redis RedisStore
}
