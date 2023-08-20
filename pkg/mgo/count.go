package mgo

import "go.mongodb.org/mongo-driver/bson"

func NewCountBuilder() *CountBuilder {
	return &CountBuilder{}
}

type CountBuilder struct {
}

func (c *CountBuilder) Do() bson.D {
	return bson.D{
		{
			Key:   "$count",
			Value: "total",
		},
	}
}
