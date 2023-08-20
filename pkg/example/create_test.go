package example

import (
	"context"
	"fmgo/pkg/mgo"
	"testing"
)

var UserDB *mgo.Mgo

func Init() {
	mgo.Connect()
	UserDB = mgo.NewMgo("user")
}

func TestCreate(t *testing.T) {
	Init()
	type User struct {
		Name string `bson:"name"`
		Age  int    `bson:"age"`
	}
	user := User{
		Name: "test",
		Age:  18,
	}
	_, err := UserDB.Creates(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
}
