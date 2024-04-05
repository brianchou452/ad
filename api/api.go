package api

import (
	"ad/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB interface {
	CreateAd(ad *model.Ad) (*mongo.InsertOneResult, error)
	GetAdByConditions(cond Query) ([]primitive.M, error)
	UpdateCurrentAds() error
	GetAdIDsBySingleCondition(field string, content interface{}) ([]model.AdSet, error)
	GetAllCountries() ([]model.Conditions, error)
	GetAllCurrentAds() ([]model.AdSet, error)
}

type RedisStore interface {
	AddAdToCache(adId string, data *model.Ad) error
	GetAdFromCache(adId string) (model.Ad, error)
	AddAdToZSet(condition string, conditionContent string, adId string, endAt float64) error
	GetCountries() []string
	GetAdsFromCondition(cond Query) ([]model.AdResponse, error)
	UpdateAdsIntersect()
	AddAd(ad model.Ad, id string) error
	ReplaceSet(key string, members []model.AdSet) error
	ReplaceCountriesSet(countries []string) error
}

type Env struct {
	DB    MongoDB
	Redis RedisStore
}
