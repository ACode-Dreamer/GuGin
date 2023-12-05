package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
	"net/http"
	"singo/data"
	"singo/logger"
)

// @Summary 状态检查
// @Description 状态检查
// @Tags 系统
// @Accept json
// @Produce json
// @NewSuccessResponse 200 {object} data.Response "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, data.NewSuccessResponse())
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) *data.Response {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range ve {
			field := fmt.Sprintf("Field.%s", e.Field)
			tag := fmt.Sprintf("Tag.Valid.%s", e.Tag)
			logger.Error("字段错误", err)
			return data.ParamErr(
				fmt.Sprintf("%s%s", field, tag),
			)
		}
	}
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		logger.Error("JSON类型不匹配", err)
		return data.ParamErr("JSON类型不匹配")
	}
	logger.Error("参数错误", err)
	return data.ParamErr("参数错误")
}
