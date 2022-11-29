package models

import "github.com/kataras/iris/v12"

// 错误
var (
	ErrorPageNotFound = iris.Map{
		"errmsg": "操作失败",
		"code":   404,
		"result": "地址未找到",
	}
)

func SuccessResponse(result interface{}) iris.Map {
	return iris.Map{
		"errmsg": "操作成功",
		"code":   200,
		"result": result,
	}
}

func ErrorResponse(result interface{}) iris.Map {
	return iris.Map{
		"errmsg": "操作失败",
		"code":   400,
		"result": result,
	}
}
