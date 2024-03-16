package api

import "ad/model"

type GormDatabase interface {
	CreateAd(ad *model.Ad) error
	GetAdByConditions(cond Query) ([]model.Ad, error)
}

type AdminEnv struct {
	DB GormDatabase
}
