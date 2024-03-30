package main

import (
	"ad/api"
	"ad/database"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/go-redis/redis/v8"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
)

func main() {
	adsUpdateDuration := 10 * time.Second

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

	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	}))

	env := &api.AdminEnv{DB: &database.MongoDB{
		DB:                    db,
		AdCollections:         db.Database("dcard_ads").Collection("ads"),
		CurrentAdsCollections: db.Database("dcard_ads").Collection("current_ads"),
	}}

	go autoUpdateCurrentAds(env.DB, adsUpdateDuration)

	r := gin.Default()
	r.RedirectFixedPath = true

	r.POST("/api/v1/ad", env.PostAdminAPIController)
	r.GET("/api/v1/ad",
		cache.CacheByRequestURI(redisStore, 20*time.Second),
		env.GetAdController)

	r.Run(":80")
}

func autoUpdateCurrentAds(db api.MongoDB, adsUpdateDuration time.Duration) {
	for {
		err := db.UpdateCurrentAds()
		if err != nil {
			log.Fatalf("Error db.UpdateCurrentAds()")
		}
		log.Println("Current Ads Updated")
		time.Sleep(adsUpdateDuration)
	}
}
