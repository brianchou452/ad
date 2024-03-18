package api

import "ad/model"

type GormDatabase interface {
	CreateAd(ad *model.Ad, age *model.Age, country *[]model.Country, platform *[]model.Platform, gender *[]model.Gender) error
	GetAdByConditions(cond Query) ([]model.Ad, error)
}

type AdminEnv struct {
	DB GormDatabase
}
