package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/common/util"
	"github.com/kuangcp/gobase/mybook/app/param"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/kuangcp/gobase/mybook/app/vo"
	"github.com/wonderivan/logger"
)

// curl -s -F typeId=4 -F accountId=3 -F categoryId=102 -F amount=2000 -F date='2020-02-03' localhost:10006/record | pretty-json
func CreateRecord(c *gin.Context) {
	// typeId 含义为 categoryTypeId
	typeId := c.PostForm("typeId")
	accountId := c.PostForm("accountId")
	targetAccountId := c.PostForm("targetAccountId")
	categoryId := c.PostForm("categoryId")
	amount := c.PostForm("amount")
	date := c.PostForm("date")
	comment := c.PostForm("comment")

	recordVO := param.CreateRecordParam{TypeId: typeId, AccountId: accountId, CategoryId: categoryId,
		Amount: amount, Date: date, Comment: comment, TargetAccountId: targetAccountId}

	logger.Debug("createRecord: ", util.Json(recordVO))

	record := service.CreateMultipleTypeRecord(recordVO)
	if record != nil {
		logger.Debug("createRecord result: ", util.Json(record))
	}

	vo.FillResult(c, record)
}

func ListRecord(c *gin.Context) {
	accountId := c.Query("accountId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	typeId := c.Query("typeId")

	query := param.QueryRecordParam{AccountId: accountId, StartDate: startDate, EndDate: endDate, TypeId: typeId}
	result := service.FindRecord(query)
	vo.FillResult(c, result)
}

func CategoryRecord(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	typeId := c.Query("typeId")

	result := service.CategoryRecord(startDate, endDate, typeId)
	vo.FillResult(c, result)
}

func CategoryDetailRecord(c *gin.Context) {
	result := service.FindRecord(buildCategoryQueryParam(c))
	vo.FillResult(c, result)
}

func WeekCategoryDetailRecord(c *gin.Context) {
	result := service.WeekCategoryRecord(buildCategoryQueryParam(c))
	vo.FillResult(c, result)
}

func MonthCategoryDetailRecord(c *gin.Context) {
	result := service.MonthCategoryRecord(buildCategoryQueryParam(c))
	vo.FillResult(c, result)
}

func buildCategoryQueryParam(c *gin.Context) param.QueryRecordParam {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	categoryId := c.Query("categoryId")
	typeId := c.Query("typeId")

	return param.QueryRecordParam{StartDate: startDate, EndDate: endDate, CategoryId: categoryId, TypeId: typeId}
}