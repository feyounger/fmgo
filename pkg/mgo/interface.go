package mgo

type UpdateAble interface {
	Filter() interface{}
	SettAble() interface{}
}

type LookUpAble interface {
	CanLookUp() bool
	GetLookUps() []LookUp
}

type MatchIn interface {
	InKey() string
	InValue() interface{}
	CanMatchIn() bool
}

type MatchOr interface {
	OrValue() []interface{}
	CanMatchOr() bool
}
