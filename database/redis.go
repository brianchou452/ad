package database

import (
	"ad/api"
	"ad/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	R        *redis.Client
	ReadOnly *redis.Client
}

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
}

func NewRedisRead() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    os.Getenv("REDIS_HOST_SLAVE") + ":" + os.Getenv("REDIS_PORT"),
	})
}

func (r *Redis) AddAdToCache(adId string, data *model.Ad) error {
	key := "ad:" + adId
	expire := time.Until(data.EndAt)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.R.Set(r.R.Context(), key, jsonData, expire).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetAdFromCache(adId string) (model.Ad, error) {
	key := "ad:" + adId

	result, err := r.R.Get(r.R.Context(), key).Result()
	if err != nil {
		return model.Ad{}, err
	}

	var ad model.Ad
	err = json.Unmarshal([]byte(result), &ad)
	if err != nil {
		return model.Ad{}, err
	}

	return ad, nil
}

func (r *Redis) AddAdToZSet(condition string, conditionContent string, adId string, endAt float64) error {
	key := "ad:" + condition + ":" + conditionContent
	err := r.R.ZAdd(r.R.Context(), key, &redis.Z{Score: endAt, Member: adId}).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetAdIdFromCondition(cond api.Query) ([]model.Ad, error) {
	ctx := r.R.Context()

	ageCondition := "ad:age:" + fmt.Sprint(cond.Age)
	ageResultName := "ad:age:All+" + ageCondition
	countryCondition := "ad:country:" + fmt.Sprint(cond.Country)
	countryResultName := "ad:country:All+" + countryCondition
	genderCondition := "ad:gender:" + fmt.Sprint(cond.Gender)
	genderResultName := "ad:gender:All+" + genderCondition
	platformCondition := "ad:platform:" + fmt.Sprint(cond.Platform)
	platformResultName := "ad:platform:All+" + platformCondition

	checkExistPipeline := r.ReadOnly.Pipeline()
	isAgeExist := checkExistPipeline.Exists(ctx, ageResultName)
	isCountryExist := checkExistPipeline.Exists(ctx, countryResultName)
	isGenderExist := checkExistPipeline.Exists(ctx, genderResultName)
	isPlatformExist := checkExistPipeline.Exists(ctx, platformResultName)
	_, err := checkExistPipeline.Exec(ctx)
	if err != nil {
		log.Println("Error:", err)
	}

	unionPipeline := r.R.Pipeline()
	if isAgeExist.Val() == 0 {
		unionPipeline.ZUnionStore(ctx, ageResultName, &redis.ZStore{
			Keys: []string{"ad:age:All", ageCondition},
		})
	}

	if isCountryExist.Val() == 0 {
		unionPipeline.ZUnionStore(ctx, countryResultName, &redis.ZStore{
			Keys: []string{"ad:country:All", countryCondition},
		})
	}

	if isGenderExist.Val() == 0 {
		unionPipeline.ZUnionStore(ctx, genderResultName, &redis.ZStore{
			Keys: []string{"ad:gender:All", genderCondition},
		})
	}

	if isPlatformExist.Val() == 0 {
		unionPipeline.ZUnionStore(ctx, platformResultName, &redis.ZStore{
			Keys: []string{"ad:platform:All", platformCondition},
		})
	}

	intersectionResult := unionPipeline.ZInter(ctx, &redis.ZStore{
		Keys: []string{ageResultName, countryResultName, genderResultName, platformResultName},
	})

	unionPipeline.Exec(ctx)

	intersectionResultArray := intersectionResult.Val()

	if cond.Offset+cond.Limit > int64(len(intersectionResultArray)) {
		cond.Limit = int64(len(intersectionResultArray)) - cond.Offset
	}
	if cond.Limit < 0 {
		cond.Limit = 0
		cond.Offset = 0
	}

	ads, err := r.ReadOnly.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, member := range intersectionResultArray[cond.Offset : cond.Offset+cond.Limit] {
			pipe.Get(ctx, "ad:"+member)
		}
		return nil
	})
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	var result []model.Ad

	for _, adString := range ads {
		var ad model.Ad
		err = json.Unmarshal([]byte(adString.(*redis.StringCmd).Val()), &ad)
		if err != nil {
			// TODO: handle error
			log.Println("Error:", err)
		}
		result = append(result, ad)
	}
	// log.Println(result)

	return result, nil
}
