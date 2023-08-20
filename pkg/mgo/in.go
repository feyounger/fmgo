package mgo

import "go.mongodb.org/mongo-driver/bson"

func NewInBuilder(inData interface{}) *InBuilder {
	return &InBuilder{data: inData}
}

type InBuilder struct {
	data interface{}
}

func (i *InBuilder) Do() bson.D {
	return bson.D{
		{
			Key:   "$in",
			Value: i.data,
		},
	}
}
