package main

import (
	"ad/api"
	"ad/database"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
)

func main() {
	adsUpdateDuration := 15 * time.Second

	version, isEnvSet := os.LookupEnv("APP_VERSION")
	if !isEnvSet {
		err := godotenv.Load(".env.dev")
		if err != nil {
			// TODO: handle error
			log.Fatalf("Error loading .env file")
		}
	} else {
		log.Printf("version: %s", version)
	}

	db, err := database.New()
	if err != nil {
		// TODO: handle error
		log.Fatalf("Error database.New()")
	}

	redisClient := database.NewRedis()

	env := &api.Env{
		DB: &database.MongoDB{
			DB:                    db,
			AdCollections:         db.Database("dcard_ads").Collection("ads"),
			CurrentAdsCollections: db.Database("dcard_ads").Collection("current_ads_0"),
		},
		Redis: &database.Redis{
			R:        redisClient,
			ReadOnly: database.NewRedisRead(),
		},
	}

	redisStore := persist.NewRedisStore(redisClient)

	go autoUpdateCurrentAds(env, adsUpdateDuration)

	r := gin.Default()
	r.RedirectFixedPath = true

	r.POST("/api/v1/ad", env.CreateAd)
	r.GET("/api/v1/ad",
		cache.CacheByRequestURI(redisStore, 60*time.Second),
		env.GetAds)

	r.Run(":80")

}

func autoUpdateCurrentAds(e *api.Env, adsUpdateDuration time.Duration) {
	var currentCollection = 0
	for {
		go func() {
			err := e.DB.UpdateCurrentAds(currentCollection)
			if err != nil {
				log.Fatalf("Error db.UpdateCurrentAds()")
			}
			log.Println("Current Ads Updated")
		}()
		go func() {
			e.Redis.UpdateAdsIntersect()
		}()
		if currentCollection == 0 {
			currentCollection = 1
		} else {
			currentCollection = 0
		}
		time.Sleep(adsUpdateDuration)
	}
}
