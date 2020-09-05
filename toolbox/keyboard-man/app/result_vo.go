package app

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS = 0
	FAILED  = 1
)

type (
	ResultVO struct {
		Data    interface{}
		Code    int8
		Success bool
		Msg     string
	}
)

var success = ResultVO{Code: SUCCESS, Success: true}
var failed = ResultVO{Code: FAILED, Success: false}

func SuccessWith(data interface{}) ResultVO {
	return ResultVO{Data: data, Code: SUCCESS, Success: true}
}

func Success() ResultVO {
	return success
}

func Failed() ResultVO {
	return failed
}

func FailedWithMsg(msg string) ResultVO {
	return ResultVO{Msg: msg, Code: FAILED, Success: false}
}

func (result ResultVO) IsSuccess() bool {
	return result.Code == SUCCESS
}
func (result ResultVO) IsFailed() bool {
	return !result.IsSuccess()
}

// web util
func GinFailed(c *gin.Context) {
	c.JSON(http.StatusOK, failed)
}

func GinSuccessWith(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessWith(data))
}

// https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
func GinResult(c *gin.Context, result interface{}) {
	if result == nil || (reflect.ValueOf(result).Kind() == reflect.Ptr && reflect.ValueOf(result).IsNil()) {
		GinFailed(c)
	} else {
		GinSuccessWith(c, result)
	}
}
