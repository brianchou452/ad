package api_test

import (
	"ad/api"
	"ad/database"
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

type ApplicationSuite struct {
	suite.Suite
	db       *database.MongoDB
	redis    *database.Redis
	a        *api.Env
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *ApplicationSuite) BeforeTest(suiteName, testName string) {

	err := godotenv.Load("../.env.dev")
	if err != nil {
		// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
		log.Fatalf("Error loading .env file")
	}

	s.recorder = httptest.NewRecorder()
	db, err := database.New()
	if err != nil {
		panic(err)
	}
	redisClient := database.NewRedis()

	s.db = &database.MongoDB{
		DB:                    db,
		AdCollections:         db.Database("dcard_ads").Collection("ads"),
		CurrentAdsCollections: db.Database("dcard_ads").Collection("current_ads"),
	}

	s.redis = &database.Redis{
		R:        redisClient,
		ReadOnly: database.NewRedisRead(),
	}

	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.a = &api.Env{
		DB: s.db,
		Redis: &database.Redis{
			R:        s.redis.R,
			ReadOnly: s.redis.ReadOnly,
		},
	}
}

func (s *ApplicationSuite) AfterTest(suiteName, testName string) {
	// db.Close()
}

func (s *ApplicationSuite) withJSON(method string, path string, value interface{}) {
	jsonVal, _ := json.Marshal(value)
	s.ctx.Request = httptest.NewRequest(method, path, bytes.NewBuffer(jsonVal))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
}

func (s *ApplicationSuite) clearDB() {
	s.db.DB.Database("dcard_ads").Drop(s.ctx)
	s.redis.R.FlushAll(s.ctx)
}
