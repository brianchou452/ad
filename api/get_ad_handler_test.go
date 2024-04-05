package api_test

import (
	"ad/model"
	"ad/model/gender"
	"ad/model/platform"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *ApplicationSuite) Test_GetAd() {
	s.clearDB()

	adData := model.Ad{
		Title:   "Test_GetAd",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(time.Hour * 24),
		Condition: model.Conditions{
			AgeStart:  18,
			AgeEnd:    19,
			Gender:    gender.Genders{gender.F, gender.M},
			Countries: []string{"TW", "JP"},
			Platform:  platform.Platforms{platform.IOS, platform.ANDROID},
			Age:       []int{18, 19},
		},
	}

	result, _ := s.db.CreateAd(&adData)
	id, _ := result.InsertedID.(primitive.ObjectID)
	s.redis.AddAd(adData, id.Hex())
	s.redis.UpdateAdsIntersect()

	query := "offset=0&limit=10&age=19&gender=M&country=TW&platform=ios"
	s.withJSON(
		"GET",
		"/api/v1/ad"+"?"+query,
		nil,
	)

	s.a.GetAds(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.recorder.Body.String(), "items")
	assert.Contains(s.T(), s.recorder.Body.String(), "Test_GetAd")
}

func (s *ApplicationSuite) Test_SomeFieldsNotSpecified() {
	s.clearDB()

	adData1 := model.Ad{
		Title:   "Test_SomeFieldsNotSpecified 1",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(time.Hour * 24),
		Condition: model.Conditions{
			AgeStart: 18,
			AgeEnd:   19,
			Age:      []int{18, 19},
		},
	}
	result, _ := s.db.CreateAd(&adData1)
	id, _ := result.InsertedID.(primitive.ObjectID)
	s.redis.AddAd(adData1, id.Hex())

	adData2 := model.Ad{
		Title:   "Test_SomeFieldsNotSpecified 2",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(time.Hour * 24),
		Condition: model.Conditions{
			AgeStart: 20,
			AgeEnd:   21,
			Age:      []int{20, 21},
		},
	}
	result, _ = s.db.CreateAd(&adData2)
	id, _ = result.InsertedID.(primitive.ObjectID)
	s.redis.AddAd(adData2, id.Hex())

	s.redis.UpdateAdsIntersect()

	query := "offset=0&limit=10&gender=M&country=TW&platform=ios"
	s.withJSON(
		"GET",
		"/api/v1/ad"+"?"+query,
		nil,
	)

	s.a.GetAds(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.recorder.Body.String(), "items")
	assert.Contains(s.T(), s.recorder.Body.String(), "Test_SomeFieldsNotSpecified 1")
	assert.Contains(s.T(), s.recorder.Body.String(), "Test_SomeFieldsNotSpecified 2")
}

func (s *ApplicationSuite) Test_LimitFieldsNotSpecified() {
	s.clearDB()

	for i := 0; i < 10; i++ {
		adData := model.Ad{
			Title:   "Test_LimitFieldsNotSpecified " + fmt.Sprint(i),
			StartAt: time.Now(),
			EndAt:   time.Now().Add(time.Hour * 24),
			Condition: model.Conditions{
				Countries: []string{"TW", "JP"},
				Platform:  platform.Platforms{platform.IOS, platform.ANDROID},
			},
		}
		result, _ := s.db.CreateAd(&adData)
		id, _ := result.InsertedID.(primitive.ObjectID)
		s.redis.AddAd(adData, id.Hex())
	}

	s.redis.UpdateAdsIntersect()

	query := "gender=F&country=TW&platform=android"
	s.withJSON(
		"GET",
		"/api/v1/ad"+"?"+query,
		nil,
	)

	s.a.GetAds(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.recorder.Body.String(), "items")
	for i := 0; i < 5; i++ {
		assert.Contains(s.T(), s.recorder.Body.String(), "Test_LimitFieldsNotSpecified "+fmt.Sprint(i))
	}
}
