package mgo

import (
	"github.com/lstack-org/go-web-framework/pkg/req"
	"go.mongodb.org/mongo-driver/bson"
)

func NewSortBuilder(sort req.SortAble) *SortBuilder {
	var flag int
	if sort.IsAsc() {
		flag = 1
	} else {
		flag = -1
	}
	return &SortBuilder{
		field: sort.SortKey(),
		flag:  flag,
	}
}

func NewSortBuilderWithKey(key string, flag int) *SortBuilder {
	return &SortBuilder{
		field: key,
		flag:  flag,
	}
}

type SortBuilder struct {
	field string
	flag  int
}

func (s *SortBuilder) Do() bson.D {
	return bson.D{
		{
			Key: "$sort",
			Value: bson.M{
				s.field: s.flag,
			},
		},
	}
}
