package code

import (
	"github.com/lstack-org/go-web-framework/pkg/code"
	"net/http"
)

var (
	DuplicateValueErr = code.ServiceCode{
		HttpCode:     http.StatusOK,
		BusinessCode: 11000,
		EnglishMsg:   "duplicate key error",
		ChineseMsg:   "值重复",
	}

	ResourceNotFound = code.ServiceCode{
		HttpCode:     http.StatusOK,
		BusinessCode: 11001,
		EnglishMsg:   "resource not found",
		ChineseMsg:   "资源不存在",
	}
)
