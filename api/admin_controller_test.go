package api_test

import (
	"ad/api"
	"ad/database"
	"ad/model"
	conutry "ad/model/country"
	"ad/model/gender"
	"ad/model/platform"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

type ApplicationSuite struct {
	suite.Suite
	db       *database.GormDatabase
	a        *api.AdminEnv
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *ApplicationSuite) BeforeTest(suiteName, testName string) {

	s.recorder = httptest.NewRecorder()
	db, err := database.NewTestDB()
	if err != nil {
		panic(err)
	}
	s.db = &database.GormDatabase{DB: db}
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.a = &api.AdminEnv{DB: s.db}
}

func (s *ApplicationSuite) AfterTest(suiteName, testName string) {
	// db.Close()
}

func (s *ApplicationSuite) Test_ensureApplicationHasCorrectJsonRepresentation() {

	s.withJSON(&api.PostData{
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
	})

	expectedJSONValue, _ := json.Marshal(&gin.H{
		"message": "success",
	})
	json, _ := strconv.Unquote(string(expectedJSONValue))

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Equal(s.T(), json, s.recorder.Body.String())
}

func (s *ApplicationSuite) withJSON(value interface{}) {
	jsonVal, _ := json.Marshal(value)
	s.ctx.Request = httptest.NewRequest("POST", "/application", bytes.NewBuffer(jsonVal))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
}
