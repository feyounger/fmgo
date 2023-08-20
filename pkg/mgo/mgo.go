package mgo

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"time"

	"k8s.io/klog/v2"

	"github.com/lstack-org/go-web-framework/pkg/req"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DefaultMongoClient   *mongo.Client
	DefaultMongoDatabase *mongo.Database
)

// Connect 数据库初始化，需要填写MongoAddr和MongoDatabase
func Connect() {
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()
	// 填写MongoAddr
	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI("MongoAddr"))
	if err != nil {
		panic(err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	// 填写MongoDatabase
	DefaultMongoDatabase = client.Database("MongoDatabase")
	DefaultMongoClient = client
}

func NewMgo(collection string) *Mgo {
	return &Mgo{DefaultMongoDatabase.Collection(collection)}
}

type Mgo struct {
	*mongo.Collection
}

// Creates 用于插入多条数据，语法：InsertMany
func (m *Mgo) Creates(ctx context.Context, arr interface{}) (*mongo.InsertManyResult, error) {
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Kind() != reflect.Slice && arrValue.Kind() != reflect.Array {
		panic("type must be Array or Slice")
	}

	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	for i := 0; i < arrValue.Len(); i++ {
		dataValue.Set(reflect.Append(dataValue, arrValue.Index(i)))
	}

	//当插入数组为空时，报错：must provide at least one element in input slice
	if len(data) == 0 {
		return nil, nil
	}

	return m.InsertMany(ctx, data)
}

// GetQuery 用于查询单条数据，语法：FindOne
func (m *Mgo) GetQuery(ctx context.Context, query, result interface{}) error {
	bsonMap, err := ConvertToBsonM(query)
	if err != nil {
		return err
	}
	singleResult := m.FindOne(ctx, bsonMap)
	if err := singleResult.Err(); err != nil {
		return err
	}
	return singleResult.Decode(result)
}

// Deletes 用于删除多条数据，语法：DeleteMany
func (m *Mgo) Deletes(ctx context.Context, eqFilter interface{}, key string, inFilter interface{}) error {
	bsonMap, err := ConvertToBsonM(eqFilter)
	if err != nil {
		return err
	}

	if key != "" && inFilter != nil {
		bsonMap[key] = inFilter
	}
	_, err = m.DeleteMany(ctx, bsonMap)
	return err
}

// UpdateOne 用于更新单条数据，语法：FindOneAndUpdate
func (m *Mgo) UpdateOne(ctx context.Context, data UpdateAble) error {
	filter, err := ConvertToBsonM(data.Filter())
	if err != nil {
		return err
	}
	update, err := NewSetBuilder().Add(data.SettAble()).Do()
	if err != nil {
		return err
	}
	err = m.FindOneAndUpdate(ctx, filter, update).Err()
	return err
}

// ListQuery 用于高级关联查询列表，语法：Aggregate
func (m *Mgo) ListQuery(ctx context.Context, query, data interface{}) (total int, err error) {
	var listRes ListRes
	err = m.RunQuery(ctx, query, &listRes)
	if err != nil {
		return
	}
	err = listRes.Convert(data)
	if err != nil {
		return
	}
	return listRes.GetTotal(), nil
}

// RunQuery 用于高级关联查询列表，语法：Aggregate
func (m *Mgo) RunQuery(ctx context.Context, query, result interface{}) error {
	var pipeline interface{}
	switch t := query.(type) {
	case []bson.D, bson.M, mongo.Pipeline:
		pipeline = t
	case bson.D:
		pipeline = mongo.Pipeline{t}
	default:
		p, err := m.pipelineGen(query)
		if err != nil {
			return err
		}

		facetBuilder := NewFacetBuilder().
			Add("total", mongo.Pipeline{NewCountBuilder().Do()})

		if paging, ok := query.(req.PageAble); ok && paging.CanPage() {
			facetBuilder.Add("items", NewPageBuilder(paging).Do())
		} else {
			facetBuilder.Add("items", p)
		}
		p = append(p, facetBuilder.Do())

		pipeline = p
	}

	bytes, _ := json.Marshal(pipeline)
	klog.V(7).Infof("mongo args: %s", string(bytes))
	cursor, err := m.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, result)
	return err
}

func (m *Mgo) pipelineGen(query interface{}) (pipeline mongo.Pipeline, _ error) {
	//往pipeline中增加模糊查询，精准匹配功能
	if bsonD, err := NewMatchBuilder().Query(query).Do(); err != nil {
		return nil, err
	} else {
		pipeline = append(pipeline, bsonD)
	}

	//增加排序功能
	if sort, ok := query.(req.SortAble); ok {
		if sort.CanSort() {
			bsonD := NewSortBuilder(sort).Do()
			pipeline = append(pipeline, bsonD)
		} else {
			//未指定排序时，默认按创建时间倒序
			bsonD := NewSortBuilderWithKey("createTime", -1).Do()
			pipeline = append(pipeline, bsonD)
		}
	}

	//增加多表关联查询功能
	if lookUp, ok := query.(LookUpAble); ok && lookUp.CanLookUp() {
		lookUpPipeline := NewLookUpBuilder().Add(lookUp.GetLookUps()...).Do()
		pipeline = append(pipeline, lookUpPipeline...)
	}

	return pipeline, nil
}
