package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewSetBuilder() *SetBuilder {
	return &SetBuilder{
		data: primitive.M{},
	}
}

type SetBuilder struct {
	err  error
	data bson.M
}

func (s *SetBuilder) Add(query interface{}) *SetBuilder {
	bsonMap, err := ConvertToBsonM(query)
	if err != nil {
		s.err = err
		return s
	}
	for key, value := range bsonMap {
		s.data[key] = value
	}
	return s
}

func (s *SetBuilder) Do() (bson.D, error) {
	if s.err != nil {
		return nil, s.err
	}
	return bson.D{
		{
			Key:   "$set",
			Value: s.data,
		},
	}, nil
}
