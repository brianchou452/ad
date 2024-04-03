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
	ctx := r.ReadOnly.Context()

	pipe := r.ReadOnly.Pipeline()
	intersectionResult := pipe.ZRange(ctx,
		"inter:lv4:country:"+cond.Country+":platform:"+cond.Platform+":gender:"+cond.Gender+":age:"+fmt.Sprint(cond.Age),
		0, -1)
	pipe.Exec(ctx)

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

func (r *Redis) UpdateAdsIntersect() {
	ctx := r.R.Context()
	countries := []string{"TW", "JP", "HK", "VN"}
	platforms := []string{"ios", "android", "web"}
	genders := []string{"M", "F"}

	pipe := r.R.Pipeline()

	for _, country := range countries {
		pipe.ZUnionStore(ctx, "union:ad:country:"+country, &redis.ZStore{
			Keys: []string{"ad:country:All", "ad:country:" + country},
		})
	}

	for _, platform := range platforms {
		pipe.ZUnionStore(ctx, "union:ad:platform:"+platform, &redis.ZStore{
			Keys: []string{"ad:platform:All", "ad:platform:" + platform},
		})
	}

	for _, gender := range genders {
		pipe.ZUnionStore(ctx, "union:ad:gender:"+gender, &redis.ZStore{
			Keys: []string{"ad:gender:All", "ad:gender:" + gender},
		})
	}

	for age := 0; age < 100; age++ {
		pipe.ZUnionStore(ctx, "union:ad:age:"+fmt.Sprint(age), &redis.ZStore{
			Keys: []string{"ad:age:All", "ad:age:" + fmt.Sprint(age)},
		})
	}

	for _, country := range countries {
		for _, platform := range platforms {
			pipe.ZInterStore(ctx, "inter:lv2:country:"+country+":platform:"+platform,
				&redis.ZStore{
					Keys: []string{
						"union:ad:country:" + country,
						"union:ad:platform:" + platform,
					},
				},
			)
		}
	}

	for _, country := range countries {
		for _, platform := range platforms {
			for _, gender := range genders {
				pipe.ZInterStore(ctx, "inter:lv3:country:"+country+":platform:"+platform+":gender:"+gender,
					&redis.ZStore{
						Keys: []string{
							"inter:lv2:country:" + country + ":platform:" + platform,
							"union:ad:gender:" + gender,
						},
					},
				)
				for age := 0; age < 100; age++ {

				}
			}
		}
	}

	for _, country := range countries {
		for _, platform := range platforms {
			for _, gender := range genders {
				for age := 0; age < 100; age++ {
					pipe.ZInterStore(ctx, "inter:lv4:country:"+country+":platform:"+platform+":gender:"+gender+":age:"+fmt.Sprint(age),
						&redis.ZStore{
							Keys: []string{
								"inter:lv3:country:" + country + ":platform:" + platform + ":gender:" + gender,
								"union:ad:age:" + fmt.Sprint(age),
							},
						},
					)
				}
			}
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		// TODO: handle error
		panic(err)
	}

}
