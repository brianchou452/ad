package database

import (
	"ad/api"
	"ad/model"
)

func (d *GormDatabase) CreateAd(ad *model.Ad) error {
	return d.DB.Create(ad).Error
}

func (d *GormDatabase) GetAdByConditions(cond api.Query) ([]model.Ad, error) {
	var ads []model.Ad

	dbqurey := d.DB.Select("title", "end_at")
	dbqurey = dbqurey.Where("end_at >= NOW() AND start_at <= NOW()")
	if cond.Age != 0 {
		dbqurey = dbqurey.Where("age_start <= ? AND age_end >= ?", cond.Age, cond.Age).Or("age_start = NULL AND age_end = NULL")
	}
	if cond.Gender != "" {
		dbqurey = dbqurey.Where("FIND_IN_SET(?, gender)", cond.Gender).Or("gender = NULL")
	}
	if cond.Country != "" {
		dbqurey = dbqurey.Where("FIND_IN_SET(?, country)", cond.Country).Or("country = NULL")
	}
	if cond.Platform != "" {
		dbqurey = dbqurey.Where("FIND_IN_SET(?, platform)", cond.Platform).Or("platform = NULL")
	}
	dbqurey.Order("end_at asc").Limit(int(cond.Limit)).Offset(int(cond.Offset)).Find(&ads)

	return ads, nil
}
