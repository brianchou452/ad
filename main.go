package main

import (
	"ad/api"
	"ad/database"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	adsUpdateDuration := 15 * time.Second

	version, isEnvSet := os.LookupEnv("APP_VERSION")
	if !isEnvSet {
		err := godotenv.Load(".env.dev")
		if err != nil {
			// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
			log.Fatalf("Error loading .env file")
		}
	} else {
		log.Printf("version: %s", version)
	}

	db, err := database.New()
	if err != nil {
		// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
		log.Fatalf("Error database.New()")
	}

	redisClient := database.NewRedis()

	env := &api.Env{
		DB: &database.MongoDB{
			DB:                    db,
			AdCollections:         db.Database("dcard_ads").Collection("ads"),
			CurrentAdsCollections: db.Database("dcard_ads").Collection("current_ads"),
		},
		Redis: &database.Redis{
			R:        redisClient,
			ReadOnly: database.NewRedisRead(),
		},
	}

	go env.Redis.ReplaceCountriesSet([]string{"TW", "JP", "HK", "VN"})
	go autoUpdateCurrentAds(env, adsUpdateDuration)

	r := gin.Default()
	r.RedirectFixedPath = true

	r.POST("/api/v1/ad", env.CreateAd)
	r.GET("/api/v1/ad", env.GetAds)

	r.Run(":80")

}

func autoUpdateCurrentAds(e *api.Env, adsUpdateDuration time.Duration) {
	for {
		go func() {
			err := e.DB.UpdateCurrentAds()
			if err != nil {
				log.Fatalf("Error db.UpdateCurrentAds()")
			}

			ReplaceCurrentAdsSet(e)
			e.Redis.UpdateAdsIntersect()

			log.Println("Current Ads Updated")
		}()
		time.Sleep(adsUpdateDuration)
	}
}

func ReplaceCurrentAdsSet(e *api.Env) {
	countries := e.Redis.GetCountries()
	platforms := []string{"ios", "android", "web"}
	genders := []string{"M", "F"}

	for _, country := range countries {
		newCountrySet, err := e.DB.GetAdIDsBySingleCondition("countries", country)
		if err != nil {
			log.Fatalf("Error db.GetAdIDsBySingleCondition()")
		}
		e.Redis.ReplaceSet("ad:country:"+country, newCountrySet)
	}

	for _, platform := range platforms {
		newPlatformSet, err := e.DB.GetAdIDsBySingleCondition("platform", platform)
		if err != nil {
			log.Fatalf("Error db.GetAdIDsBySingleCondition()")
		}
		e.Redis.ReplaceSet("ad:platform:"+platform, newPlatformSet)
	}

	for _, gender := range genders {
		newGenderSet, err := e.DB.GetAdIDsBySingleCondition("gender", gender)
		if err != nil {
			log.Fatalf("Error db.GetAdIDsBySingleCondition()")
		}
		e.Redis.ReplaceSet("ad:gender:"+gender, newGenderSet)
	}

	for age := 0; age < 100; age++ {
		newAgeSet, err := e.DB.GetAdIDsBySingleCondition("age", age)
		if err != nil {
			log.Fatalf("Error db.GetAdIDsBySingleCondition()")
		}
		e.Redis.ReplaceSet("ad:age:"+fmt.Sprint(age), newAgeSet)
	}

	for _, condition := range []string{"country", "platform", "gender", "age"} {
		set, err := e.DB.GetAdIDsBySingleCondition(condition, nil)
		if err != nil {
			log.Fatalf("Error db.GetAdIDsBySingleCondition()")
		}
		e.Redis.ReplaceSet("ad:"+condition+":NotSpecified", set)
	}

	currentAds, err := e.DB.GetAllCurrentAds()
	if err != nil {
		log.Fatalf("Error db.GetAllCurrentAds()")
	}
	for _, condition := range []string{"country", "platform", "gender", "age"} {
		e.Redis.ReplaceSet("ad:"+condition+":All", currentAds)
	}

}
