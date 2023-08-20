package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewFacetBuilder() *FacetBuilder {
	return &FacetBuilder{
		data: primitive.M{},
	}
}

type FacetBuilder struct {
	data bson.M
}

func (m *FacetBuilder) Add(key string, value mongo.Pipeline) *FacetBuilder {
	m.data[key] = value
	return m
}

func (m *FacetBuilder) Do() bson.D {
	return bson.D{
		{
			Key:   "$facet",
			Value: m.data,
		},
	}
}
