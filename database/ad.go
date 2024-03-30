package database

import (
	"ad/api"
	"ad/model"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *MongoDB) CreateAd(ad *model.Ad) (*mongo.InsertOneResult, error) {
	result, err := d.AdCollections.InsertOne(context.TODO(), ad)
	return result, err
}

func (d *MongoDB) GetAdByConditions(cond api.Query) ([]primitive.M, error) {
	// if matchCondition == nil {
	// 	matchCondition = []bson.E{}
	// }

	// 建立一個空的map來儲存查詢條件
	matchCondition := make(map[string]interface{})

	if cond.Country != "" {
		matchCondition["conditions.country"] = bson.D{
			{Key: "$in", Value: bson.A{cond.Country, nil}},
		}
	}

	if cond.Gender != "" {
		matchCondition["conditions.gender"] = bson.D{
			{Key: "$in", Value: bson.A{cond.Gender, nil}},
		}
	}

	if cond.Platform != "" {
		matchCondition["conditions.platform"] = bson.D{
			{Key: "$in", Value: bson.A{cond.Platform, nil}},
		}
	}

	if cond.Age != 0 {
		matchCondition["conditions.age"] = bson.D{
			{Key: "$in", Value: bson.A{cond.Age, nil}},
		}
	}

	bsonMatchCondition := make(bson.D, 0, len(matchCondition))
	for key, value := range matchCondition {
		bsonMatchCondition = append(bsonMatchCondition, bson.E{Key: key, Value: value})
	}

	// https://www.mongodb.com/community/forums/t/mongodb-go-primative-e/168870
	findAd := bson.A{
		bson.D{
			{Key: "$match",
				Value: bsonMatchCondition,
			},
		},
		// TODO: 會讓query變慢，改在API內部處理
		// bson.D{
		// 	{Key: "$group",
		// 		Value: bson.D{
		// 			{Key: "_id", Value: "$adId"},
		// 			{Key: "title", Value: bson.D{{Key: "$first", Value: "$title"}}},
		// 			{Key: "endAt", Value: bson.D{{Key: "$first", Value: "$endAt"}}},
		// 		},
		// 	},
		// },
		// bson.D{
		// 	{Key: "$sort",
		// 		Value: bson.D{
		// 			{Key: "endAt", Value: 1},
		// 		},
		// 	},
		// },
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "title", Value: 1},
					{Key: "endAt", Value: 1},
				},
			},
		},
		bson.D{
			{Key: "$skip",
				Value: cond.Offset,
			},
		},
		bson.D{
			{Key: "$limit",
				Value: cond.Limit,
			},
		},
	}

	cursor, err := d.CurrentAdsCollections.Aggregate(context.TODO(), findAd)
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return results, err
}

func (d *MongoDB) UpdateCurrentAds(collectionId int) error {

	var coll *mongo.Collection

	if collectionId == 0 {
		coll = d.DB.Database("dcard_ads").Collection("current_ads_0")
	} else {
		coll = d.DB.Database("dcard_ads").Collection("current_ads_1")
	}

	_, err := d.AdCollections.Aggregate(context.TODO(), bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "startAt", Value: bson.D{{Key: "$lte", Value: time.Now()}}},
					{Key: "endAt", Value: bson.D{{Key: "$gte", Value: time.Now()}}},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$conditions.gender"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$conditions.country"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$conditions.platform"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{{Key: "$addFields", Value: bson.D{{Key: "adId", Value: "$_id"}}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}}}},
		bson.D{
			{Key: "$out",
				Value: bson.D{
					{Key: "db", Value: "dcard_ads"},
					{Key: "coll", Value: coll.Name()},
				},
			},
		},
	})
	if err != nil {
		// TODO: handle error
		log.Fatal(err)
		return err
	}

	d.CurrentAdsCollections = coll

	return nil
}
