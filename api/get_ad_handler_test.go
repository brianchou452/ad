package api_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func (s *ApplicationSuite) Test_GetAd() {
	s.clearDB()

	s.withJSON(
		"GET",
		"/api/v1/ad",
		nil,
	)
	s.ctx.Params = gin.Params{
		{Key: "offset", Value: "0"},
		{Key: "limit", Value: "10"},
		{Key: "country", Value: "TW"},
		{Key: "gender", Value: "F"},
		{Key: "platform", Value: "web"},
		{Key: "age", Value: "18"},
	}

	s.a.GetAds(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
}
