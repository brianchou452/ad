package api_test

import (
	"ad/api"
	"ad/model"
	conutry "ad/model/country"
	"ad/model/gender"
	"ad/model/platform"

	"github.com/stretchr/testify/assert"
)

func (s *ApplicationSuite) Test_ensureApplicationHasCorrectJsonRepresentation() {

	s.withJSON(
		"POST",
		"/api/v1/ad",
		&api.PostData{
			Title:   "name",
			StratAt: "2021-01-01T00:00:00Z",
			EndAt:   "2021-01-01T00:00:00Z",
			Conditons: model.Conditon{
				AgeStart: 18,
				AgeEnd:   65,
				Gender:   gender.Genders{gender.F, gender.M},
				Country:  conutry.Countrys{conutry.TW, conutry.JP},
				Platform: platform.Platforms{platform.WEB, platform.IOS},
			},
		},
	)

	s.a.CreateAd(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.recorder.Body.String(), "id")
}
