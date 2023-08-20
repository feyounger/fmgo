package mgo

import (
	"fmt"
	"strings"

	"github.com/lstack-org/go-web-framework/pkg/req"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMatchBuilder() *MatchBuilder {
	return &MatchBuilder{
		data: primitive.M{},
	}
}

type MatchBuilder struct {
	err  error
	data bson.M
}

func (m *MatchBuilder) EqWithKV(key string, value interface{}) *MatchBuilder {
	m.data[key] = value
	return m
}

func (m *MatchBuilder) Eq(query interface{}) *MatchBuilder {
	bsonMap, err := ConvertToBsonM(query)
	if err != nil {
		m.err = err
		return m
	}
	for key, value := range bsonMap {
		m.data[key] = value
	}
	return m
}

func (m *MatchBuilder) Or(or MatchOr, query interface{}) *MatchBuilder {
	value := or.OrValue()
	if matchIn, ok := query.(MatchIn); ok && matchIn.CanMatchIn() {
		value = append(value, map[string]interface{}{
			matchIn.InKey(): NewInBuilder(matchIn.InValue()).Do(),
		})
	}
	m.data["$or"] = value
	return m
}

func (m *MatchBuilder) In(in MatchIn) *MatchBuilder {
	m.data[in.InKey()] = NewInBuilder(in.InValue()).Do()
	return m
}

func (m *MatchBuilder) Regex(search req.SearchAble) *MatchBuilder {
	m.data[search.SKey()] = bson.M{
		"$regex": primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", CheckRegexString(search)),
			Options: "i",
		},
	}
	return m
}

func (m *MatchBuilder) Query(query interface{}) *MatchBuilder {
	builder := m.Eq(query)
	if matchOr, ok := query.(MatchOr); ok && matchOr.CanMatchOr() {
		builder.Or(matchOr, query)
	} else {
		if matchIn, ok := query.(MatchIn); ok && matchIn.CanMatchIn() {
			builder.In(matchIn)
		}
	}
	if search, ok := query.(req.SearchAble); ok && search.CanSearch() {
		return builder.Regex(search)
	}

	return builder
}

func (m *MatchBuilder) Do() (bson.D, error) {
	if m.err != nil {
		return nil, m.err
	}
	return bson.D{
		{
			Key:   "$match",
			Value: m.data,
		},
	}, nil
}

func (m *MatchBuilder) ElemMatchDo() (bson.D, error) {
	if m.err != nil {
		return nil, m.err
	}
	return bson.D{
		{
			Key:   "$elemMatch",
			Value: m.data,
		},
	}, nil
}

func CheckRegexString(search req.SearchAble) string {
	str := fmt.Sprintf("%v", search.SValue())
	//正则匹配出现的特殊字符串
	fbsArr := []string{"\\", "$", "(", ")", "*", "+", ".", "[", "]", "?", "^", "{", "}", "|"}
	for _, ch := range fbsArr {
		if StrContainers := strings.Contains(str, ch); StrContainers {
			str = strings.Replace(str, ch, "\\"+ch, -1)
		}
	}
	return str
}
