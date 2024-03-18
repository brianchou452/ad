package database

import (
	"ad/api"
	"ad/model"
	"fmt"
	"log"
	"sort"

	"gorm.io/gorm"
)

func (d *GormDatabase) CreateAd(ad *model.Ad, age *model.Age, country *[]model.Country, platform *[]model.Platform, gender *[]model.Gender) error {
	err := d.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(ad).Error; err != nil {
			return err
		}

		adID := ad.ID

		if age != nil {
			age.AdID = adID
			if err := tx.Create(age).Error; err != nil {
				return err
			}
		}
		if len(*country) != 0 {
			for i := range *country {
				(*country)[i].AdID = adID
			}
			if err := tx.Create(country).Error; err != nil {
				return err
			}
		}
		if len(*platform) != 0 {
			for i := range *platform {
				(*platform)[i].AdID = adID
			}
			if err := tx.Create(platform).Error; err != nil {
				return err
			}
		}
		if len(*gender) != 0 {
			for i := range *gender {
				(*gender)[i].AdID = adID
			}
			if err := tx.Create(gender).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

type AdID struct {
	ID int64
}

func (d *GormDatabase) GetAdByConditions(cond api.Query) ([]model.Ad, error) {
	var ads []model.Ad

	// var adIDs []int64
	var prevAdIDs []int64

	err := d.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.Ad{}).Select("id").Where("end_at >= NOW() AND start_at <= NOW()").Order("id asc").Find(&prevAdIDs)

		if cond.Gender != "" {
			var genderIDs []int64
			tx.Model(&model.Gender{}).Select("ad_id").
				Where("ad_id IN ?", prevAdIDs).
				Where("gender = ?", cond.Gender).
				Or("gender = NULL").
				Order("ad_id asc").Find(&genderIDs)
			log.Println("gender " + cond.Gender)
			prevAdIDs = intersectInt64s(prevAdIDs, genderIDs)
			// prevAdIDs = genderIDs
			// log.Println(prevAdIDs)
		}
		if cond.Country != "" {
			var countryIDs []int64
			tx.Model(&model.Country{}).Select("ad_id").
				Where("ad_id IN ?", prevAdIDs).
				Where("country = ?", cond.Country).
				Or("country = NULL").
				Order("ad_id asc").Find(&countryIDs)
			log.Println("country " + cond.Country)
			prevAdIDs = intersectInt64s(prevAdIDs, countryIDs)
			// log.Println(prevAdIDs)
		}
		if cond.Platform != "" {
			var platformIDs []int64
			tx.Model(&model.Platform{}).Select("ad_id").
				Where("ad_id IN ?", prevAdIDs).
				Where("platform = ?", cond.Platform).
				Or("platform = NULL").
				Order("ad_id asc").Find(&platformIDs)
			log.Println("platform " + cond.Platform)
			prevAdIDs = intersectInt64s(prevAdIDs, platformIDs)
			// log.Println(prevAdIDs)
		}
		if cond.Age != 0 {
			var ageIDs []int64
			tx.Model(&model.Age{}).Select("ad_id").
				Where("ad_id IN ?", prevAdIDs).
				Where("age_start <= ? AND age_end >= ?", cond.Age, cond.Age).
				Or("age_start = NULL AND age_end = NULL").
				Order("ad_id asc").Find(&ageIDs)
			log.Println("age " + fmt.Sprint(cond.Age))
			prevAdIDs = intersectInt64s(prevAdIDs, ageIDs)
			// log.Println(prevAdIDs)
		}

		log.Println(prevAdIDs)

		if len(prevAdIDs) == 0 {
			return nil
		}

		tx = tx.Model(&model.Ad{}).Select("title", "end_at").Where("id IN ?", prevAdIDs)
		tx = tx.Order("end_at asc").Limit(int(cond.Limit)).Offset(int(cond.Offset))
		err := tx.Find(&ads).Error
		return err
	})

	return ads, err
}

// TODO: 檢查正確性
// Sorted has complexity: O(n * log(n)), a needs to be sorted
func SortedGeneric[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		idx := sort.Search(len(b), func(i int) bool {
			return b[i] == v
		})
		if idx < len(b) && b[idx] == v {
			set = append(set, v)
		}
	}

	return set
}

func intersectInt64s(a []int64, b []int64) []int64 {
	// sort.Slice(a, func(i, j int) bool {
	// 	return a[i] < a[j]
	// })
	// sort.Slice(b, func(i, j int) bool {
	// 	return b[i] < b[j]
	// })

	var result []int64
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			result = append(result, a[i])
			i++
			j++
		}
	}

	return result
}
