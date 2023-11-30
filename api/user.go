package api

import (
	"net/http"
	"singo/req"
	"singo/service"

	"github.com/gin-gonic/gin"
)

// @Summary 用户注册接口
// @Description 用户注册接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.UserRegisterReq true "请求参数"
// @Success 200 {object} data.Response{data=data.UserReq} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/user/register [post]
func UserRegister(c *gin.Context) {
	var param service.UserRegisterReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.Register(&param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 用户登录接口
// @Description 用户登录接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.UserLoginReq true "请求参数"
// @Success 200 {object} data.Response{data=data.UserReq} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/user/login [post]
func UserLogin(c *gin.Context) {
	var param service.UserLoginReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.Login(&param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 用户详情接口
// @Description 用户详情接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} data.Response{data=model.User} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/user/info [get]
func UserMe(c *gin.Context) {
	user := service.Me(c.GetString("username"))
	c.JSON(http.StatusOK, user)
}

// @Summary 用户列表接口
// @Description 用户列表接口
// @Tags 用户
// @Accept x-www-form-urlencoded
// @Produce json
// @Param request query req.PageUserReq true "请求参数"
// @Param Authorization header string true "token"
// @Success 200 {object} data.Response{data=data.Pagination{items=[]model.User}} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/user/list [get]
func Get(c *gin.Context) {
	var param req.PageUserReq
	if err := c.ShouldBindQuery(&param); err == nil {
		res := service.GetAllUsers(&param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}
