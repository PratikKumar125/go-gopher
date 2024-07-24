package common_repository

import "go.mongodb.org/mongo-driver/bson"

func GetPaginationAggregation(initialStage bson.D, page int, limit int) bson.A {
	skip := 0
	if page > 1 {
		skip = ((page - 1) * limit)
	}
	return bson.A{
		initialStage,
		bson.D{{Key: "$facet", Value: bson.D{
			{Key: "count", Value: bson.A{
				bson.D{{Key: "$count", Value: "total"}},
			}},
			{Key: "docs", Value: bson.A{
				bson.D{{Key: "$skip", Value: skip}},
				bson.D{{Key: "$limit", Value: limit}},
			}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "total", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$count", 0}},
			}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "count", Value: "$total.total"},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "meta", Value: bson.D{
				{Key: "total", Value: "$total.total"},
				{Key: "hasNextPage", Value: bson.D{
					{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{
							{Key: "$gt", Value: bson.A{"$count", bson.D{
								{Key: "$multiply", Value: bson.A{limit, 1}},
							}}},
						}},
						{Key: "then", Value: true},
						{Key: "else", Value: false},
					}},
				}},
				{Key: "hasPrevPage", Value: bson.D{
					{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{
							{Key: "$gt", Value: bson.A{skip, 0}},
						}},
						{Key: "then", Value: true},
						{Key: "else", Value: false},
					}},
				}},
				{Key: "totalPages", Value: bson.D{
					{Key: "$ceil", Value: bson.D{
						{Key: "$divide", Value: bson.A{"$total.total", limit}},
					}},
				}},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "total", Value: 0},
			{Key: "count", Value: 0},
		}}},
	}
}
