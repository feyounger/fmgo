package mgo

import (
	"encoding/json"
	"time"

	"github.com/lstack-org/go-web-framework/pkg/req"
	"go.mongodb.org/mongo-driver/bson"
)

func NewIdFilter(id, primaryUserId string) *IdFilter {
	return &IdFilter{
		Id:            id,
		PrimaryUserId: primaryUserId,
	}
}

type IdFilter struct {
	Id            string `bson:"id"`
	PrimaryUserId string `bson:"primaryUserId"`
}

func NewBase(initer req.IamToken) (base Base) {
	timeNow := time.Now().UTC()
	base.CreateTime = &timeNow
	base.CreatorUserId = initer.GetUserID()
	base.PrimaryUserId = initer.GetPrimaryAuthUserID()
	base.CreatorUserName = initer.GetUserName()
	base.PrimaryUserName = initer.GetPrimaryAuthUserName()
	return
}

type Base struct {
	CreateTime      *time.Time `json:"createTime,omitempty" bson:"createTime,omitempty"`
	CreatorUserId   string     `json:"creatorUserId" bson:"creatorUserId,omitempty"`
	CreatorUserName string     `json:"creatorUserName" bson:"creatorUserName,omitempty"`
	PrimaryUserId   string     `json:"primaryUserId" bson:"primaryUserId,omitempty"`
	PrimaryUserName string     `json:"primaryUserName" bson:"primaryUserName,omitempty"`
}

type ListRes []ListData

func (l ListRes) Convert(data interface{}) error {
	if len(l) == 0 {
		return nil
	}
	bytes, err := json.Marshal(l[0].Items)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, data)
}

func (l ListRes) GetTotal() int {
	if len(l) == 0 {
		return 0
	}
	return l[0].GetTotal()
}

type ListData struct {
	Items []bson.M `json:"items"`
	Total []Total  `json:"total"`
}

func (l *ListData) GetTotal() int {
	if l == nil || len(l.Total) == 0 {
		return 0
	}
	return l.Total[0].Total
}

type Total struct {
	Total int `json:"total"`
}

func ConvertToBsonM(data interface{}) (bson.M, error) {
	bytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	bsonMap := bson.M{}
	err = bson.Unmarshal(bytes, bsonMap)
	return bsonMap, err
}

type LookUp struct {
	From         string `bson:"from"`
	LocalField   string `bson:"localField"`
	ForeignField string `bson:"foreignField"`
	As           string `bson:"as"`
}

type LookUps []LookUp

func (l LookUps) CanLookUp() bool {
	return len(l) != 0
}

func (l LookUps) GetLookUps() []LookUp {
	return l
}

type InMatch struct {
	InMatchKey   string      `json:"inMatchKey"`
	InMatchValue interface{} `json:"inMatchValue"`
}

func (m InMatch) InKey() string {
	return m.InMatchKey
}

func (m InMatch) InValue() interface{} {
	return m.InMatchValue
}

func (m InMatch) CanMatchIn() bool {
	return m.InMatchKey != "" && m.InMatchValue != nil
}

func NewOrMatchEq(eq interface{}) OrMatch {
	m, _ := ConvertToBsonM(eq)

	var d []interface{}
	for k, v := range m {
		d = append(d, map[string]interface{}{
			k: v,
		})
	}
	return OrMatch{OrMatchValue: d}
}

type OrMatch struct {
	OrMatchValue []interface{} `json:"orMatchValue"`
}

func (o OrMatch) OrValue() []interface{} {
	return o.OrMatchValue
}

func (o OrMatch) CanMatchOr() bool {
	return o.OrMatchValue != nil
}

type IdName struct {
	Id      string `json:"id,omitempty" bson:"id,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	OrgId   string `json:"orgId,omitempty" bson:"orgId,omitempty"`
	Project string `json:"project,omitempty" bson:"project,omitempty"`
}

type OwnerFilter struct {
	Owner string `json:"owner" bson:"owner"`
}
