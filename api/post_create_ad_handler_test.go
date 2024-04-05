package api_test

import (
	"ad/api"
	"ad/model"
	"ad/model/gender"
	"ad/model/platform"

	"github.com/stretchr/testify/assert"
)

func (s *ApplicationSuite) Test_insserAdWithEmptyCondition() {

	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "Test_insserAdWithEmptyCondition",
			StratAt: "2021-01-01T00:00:00Z",
			EndAt:   "2021-01-01T00:00:00Z",
		},
	)
	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.recorder.Body.String(), "id")
}

func (s *ApplicationSuite) Test_startAtAfterEndAt() {
	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "Test_startAtAfterEndAt",
			StratAt: "2021-01-02T00:00:00Z",
			EndAt:   "2021-01-01T00:00:00Z",
			Conditons: model.Conditions{
				AgeStart:  18,
				AgeEnd:    65,
				Gender:    gender.Genders{gender.F, gender.M},
				Countries: []string{"TW", "JP"},
				Platform:  platform.Platforms{platform.WEB, platform.IOS},
			},
		},
	)
	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *ApplicationSuite) Test_conditionAgeStartBiggerThenAgeEnd() {
	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "Test_conditionAgeStartBiggerThenAgeEnd",
			StratAt: "2021-01-01T00:00:00Z",
			EndAt:   "2021-01-02T00:00:00Z",
			Conditons: model.Conditions{
				AgeStart: 65,
				AgeEnd:   18,
			},
		},
	)
	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *ApplicationSuite) Test_conditionAgeStartSameAsAgeEnd() {
	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "Test_conditionAgeStartSameAsAgeEnd",
			StratAt: "2021-01-01T00:00:00Z",
			EndAt:   "2021-01-02T00:00:00Z",
			Conditons: model.Conditions{
				AgeStart: 18,
				AgeEnd:   18,
			},
		},
	)
	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *ApplicationSuite) Test_sentJsonWithMissingField() {
	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "Test_sentJsonWithMissingField",
			StratAt: "2021-01-01T00:00:00Z",
		},
	)
	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *ApplicationSuite) Test_sentJsonWithWrongDateFormat() {
	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "Test_sentJsonWithWrongDateFormat",
			StratAt: "2021/01/01",
			EndAt:   "2021/01/02",
		},
	)
	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}
