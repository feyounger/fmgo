package mgo

import (
	"fmgo/pkg/code"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	duplicateKey = "E11000 duplicate key error"
	noResult     = "mongo: no documents in result"
)

// CheckIndexConflict 用于检测是否是索引冲突错误，若无错误，返回nil
// 若是索引冲突，返回code.DingNotifyError，key表示冲突的字段值
func CheckIndexConflict(err error, key string) error {
	if err == nil {
		return nil
	}
	if IsIndexConflictError(err) {
		return code.DuplicateValueErr.MergeObj(key)
	}
	return err
}

func IsIndexConflictError(err error) bool {
	switch we := err.(type) {
	case mongo.WriteException:
		for _, e := range we.WriteErrors {
			if e.Code == 11000 && strings.Contains(err.Error(), duplicateKey) {
				return true
			}
		}
	case mongo.CommandError:
		if we.Code == 11000 && strings.Contains(we.Message, duplicateKey) {
			return true
		}
	}

	return false
}

// CheckNoDocumentsInResult 用于在更新单个资源或获取单个资源时，检测对应的资源是否存在
func CheckNoDocumentsInResult(err error, key string) error {
	if err == nil {
		return nil
	}
	if IsNoDocumentsInResult(err) {
		return code.ResourceNotFound.MergeObj(key)
	}
	return err
}

func IsNoDocumentsInResult(err error) bool {
	return err.Error() == noResult
}

func ErrCheck(err error, id, name string) error {
	if err == nil {
		return nil
	}
	if IsIndexConflictError(err) {
		return code.DuplicateValueErr.MergeObj(name)
	}
	if IsNoDocumentsInResult(err) {
		return code.ResourceNotFound.MergeObj(id)
	}
	return err
}
