package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLookUpBuilder() *LookUpBuilder {
	return &LookUpBuilder{
		data: make([]bson.M, 0),
	}
}

type LookUpBuilder struct {
	data []bson.M
}

func (l *LookUpBuilder) Add(lookUps ...LookUp) *LookUpBuilder {
	for _, lookUp := range lookUps {
		bsonM, _ := ConvertToBsonM(lookUp)
		l.data = append(l.data, bsonM)
	}
	return l
}

func (l *LookUpBuilder) Do() (pipeline mongo.Pipeline) {
	for _, bsonM := range l.data {
		pipeline = append(pipeline, bson.D{
			{
				Key:   "$lookup",
				Value: bsonM,
			},
		})
	}
	return
}
