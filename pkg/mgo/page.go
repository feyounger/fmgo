package mgo

import (
	"github.com/lstack-org/go-web-framework/pkg/req"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPageBuilder(paging req.PageAble) *PageBuilder {
	return &PageBuilder{
		page:     paging.PNumber(),
		pageSize: paging.PSize(),
	}
}

type PageBuilder struct {
	page     int
	pageSize int
}

func (p *PageBuilder) Do() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{
			{
				Key:   "$skip",
				Value: (p.page - 1) * p.pageSize,
			},
		},
		bson.D{
			{
				Key:   "$limit",
				Value: p.pageSize,
			},
		},
	}
}
