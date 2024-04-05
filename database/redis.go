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

func (r *Redis) AddAd(ad model.Ad, id string) error {
	ctx := r.R.Context()
	cmd, err := r.ReadOnly.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		endAt := float64(ad.EndAt.Unix())

		expire := time.Until(ad.EndAt)
		adResponse := model.AdResponse{
			Title: ad.Title,
			EndAt: ad.EndAt,
		}
		jsonData, err := json.Marshal(adResponse)
		if err != nil {
			return err
		}
		pipe.Set(ctx, "ad:"+id, jsonData, expire)

		for _, country := range ad.Condition.Countries {
			log.Println("ad:country:" + country)
			pipe.SAdd(ctx, "ad:countries", country)
			pipe.ZAdd(ctx, "ad:country:"+country, &redis.Z{Score: endAt, Member: id})
		}
		for _, platform := range ad.Condition.Platform {
			pipe.ZAdd(ctx, "ad:platform:"+fmt.Sprint(platform), &redis.Z{Score: endAt, Member: id})
		}
		for _, gender := range ad.Condition.Gender {
			pipe.ZAdd(ctx, "ad:gender:"+fmt.Sprint(gender), &redis.Z{Score: endAt, Member: id})
		}
		for _, age := range ad.Condition.Age {
			pipe.ZAdd(ctx, "ad:age:"+fmt.Sprint(age), &redis.Z{Score: endAt, Member: id})
		}

		if ad.Condition.Countries == nil {
			log.Println("ad:country:NotSpecified")
			pipe.ZAdd(ctx, "ad:country:NotSpecified", &redis.Z{Score: endAt, Member: id})
		}
		if ad.Condition.Platform == nil {
			pipe.ZAdd(ctx, "ad:platform:NotSpecified", &redis.Z{Score: endAt, Member: id})
		}
		if ad.Condition.Gender == nil {
			pipe.ZAdd(ctx, "ad:gender:NotSpecified", &redis.Z{Score: endAt, Member: id})
		}
		if len(ad.Condition.Age) == 0 {
			pipe.ZAdd(ctx, "ad:age:NotSpecified", &redis.Z{Score: endAt, Member: id})
		}

		pipe.ZAdd(ctx, "ad:country:All", &redis.Z{Score: endAt, Member: id})
		pipe.ZAdd(ctx, "ad:platform:All", &redis.Z{Score: endAt, Member: id})
		pipe.ZAdd(ctx, "ad:gender:All", &redis.Z{Score: endAt, Member: id})
		pipe.ZAdd(ctx, "ad:age:All", &redis.Z{Score: endAt, Member: id})

		return nil
	})
	if err != nil {
		return err
	}
	for _, c := range cmd {
		if c.Err() != nil {
			return c.Err()
		}
	}
	return nil
}

func (r *Redis) AddAdToZSet(condition string, conditionContent string, adId string, endAt float64) error {
	key := "ad:" + condition + ":" + conditionContent
	err := r.R.ZAdd(r.R.Context(), key, &redis.Z{Score: endAt, Member: adId}).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetAdsFromCondition(cond api.Query) ([]model.AdResponse, error) {
	ctx := r.ReadOnly.Context()

	var key string
	if cond.Country == "" {
		cond.Country = "All"
	}
	if cond.Platform == "" {
		cond.Platform = "All"
	}
	if cond.Gender == "" {
		cond.Gender = "All"
	}
	if cond.Age == 0 {
		key = "inter:lv4:country:" + cond.Country + ":platform:" + cond.Platform + ":gender:" + cond.Gender + ":age:" + "All"
	} else {
		key = "inter:lv4:country:" + cond.Country + ":platform:" + cond.Platform + ":gender:" + cond.Gender + ":age:" + fmt.Sprint(cond.Age)
	}
	log.Println("key:", key)

	// TODO: remove pipeline
	pipe := r.ReadOnly.Pipeline()
	intersectionResult := r.R.ZRange(ctx, key, 0, -1)
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
		// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
		panic(err)
	}

	var result []model.AdResponse

	for _, adString := range ads {
		var ad model.AdResponse
		err = json.Unmarshal([]byte(adString.(*redis.StringCmd).Val()), &ad)
		if err != nil {
			// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
			log.Println("Error:", err)
		}
		result = append(result, ad)
	}

	return result, nil
}

func (r *Redis) GetCountries() []string {
	ctx := r.ReadOnly.Context()
	countries := r.R.SMembers(ctx, "ad:countries").Val()
	return countries
}

func (r *Redis) UpdateAdsIntersect() {
	ctx := r.R.Context()
	countries := r.GetCountries()
	countries = append(countries, "All")
	platforms := []string{"ios", "android", "web", "All"}
	genders := []string{"M", "F", "All"}

	pipe := r.R.Pipeline()

	for _, country := range countries {
		pipe.ZUnionStore(ctx, "union:ad:country:"+country, &redis.ZStore{
			Keys: []string{"ad:country:NotSpecified", "ad:country:" + country},
		})
	}

	for _, platform := range platforms {
		pipe.ZUnionStore(ctx, "union:ad:platform:"+platform, &redis.ZStore{
			Keys: []string{"ad:platform:NotSpecified", "ad:platform:" + platform},
		})
	}

	for _, gender := range genders {
		pipe.ZUnionStore(ctx, "union:ad:gender:"+gender, &redis.ZStore{
			Keys: []string{"ad:gender:NotSpecified", "ad:gender:" + gender},
		})
	}

	for age := 0; age < 100; age++ {
		pipe.ZUnionStore(ctx, "union:ad:age:"+fmt.Sprint(age), &redis.ZStore{
			Keys: []string{"ad:age:NotSpecified", "ad:age:" + fmt.Sprint(age)},
		})
	}
	pipe.ZUnionStore(ctx, "union:ad:age:"+"All", &redis.ZStore{
		Keys: []string{"ad:age:" + "All", "ad:age:" + "All"},
	})

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
				pipe.ZInterStore(ctx, "inter:lv4:country:"+country+":platform:"+platform+":gender:"+gender+":age:"+"All",
					&redis.ZStore{
						Keys: []string{
							"inter:lv3:country:" + country + ":platform:" + platform + ":gender:" + gender,
							"union:ad:age:" + "All",
						},
					},
				)
			}
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
		panic(err)
	}

}

func (r *Redis) ReplaceSet(key string, members []model.AdSet) error {
	ctx := r.R.Context()
	_, err := r.R.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, key)
		for _, member := range members {
			pipe.ZAdd(r.R.Context(), key, &redis.Z{Score: float64(member.EndAt.Unix()), Member: member.Id})
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error Redis ReplaceSet()" + err.Error())
		return err
	}
	return nil
}

func (r *Redis) ReplaceCountriesSet(countries []string) error {
	ctx := r.R.Context()
	pipe := r.R.Pipeline()
	pipe.Del(ctx, "ad:countries")
	for _, country := range countries {
		pipe.SAdd(ctx, "ad:countries", country)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error Redis ReplaceCountriesSet()" + err.Error())
		return err
	}
	return nil
}

// func (r *Redis) InitNeededSets(e *api.Env) {
// 	ctx := r.R.Context()
// 	pipe := r.R.Pipeline()
// 	for _, condition := range []string{"country", "platform", "gender", "age"} {

// 	}
// 	_, err := pipe.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("Error Redis ReplaceCountriesSet()" + err.Error())
// 		return err
// 	}
// }
