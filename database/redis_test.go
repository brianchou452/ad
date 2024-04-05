package database_test

import (
	"ad/api"
	"ad/database"
	"testing"

	"github.com/joho/godotenv"
)

var env *database.Redis

func init() {
	err := godotenv.Load("../.env.dev")
	if err != nil {
		panic("Error loading .env file")
	}
	env = &database.Redis{
		R:        database.NewRedis(),
		ReadOnly: database.NewRedisRead(),
	}
}

func Benchmark_GetAdIdFromCondition(b *testing.B) {
	// b.ReportAllocs()
	// b.ReportMetric(1, "ns/op")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		env.GetAdsFromCondition(
			api.Query{
				Offset:   0,
				Limit:    10,
				Country:  "US",
				Gender:   "M",
				Age:      20,
				Platform: "ios",
			},
		)
	}

}
