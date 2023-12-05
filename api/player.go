package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"singo/data"
	"singo/req"
	"singo/service"
)

// @Summary 玩家登录接口
// @Description 玩家登录接口
// @Tags 玩家
// @Accept json
// @Produce json
// @Param code path string true "优版权登录Code"
// @Success 200 {object} data.Response{data=data.PlayerLoginResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/{code} [get]
func PlayerLogin(c *gin.Context) {
	code := c.Param("code")
	if code != "" {
		res := service.PlayerLogin(code)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, data.NewErrorResponse(30001, "登录Code不得为空"))
	}
}

// @Summary 开始游戏接口
// @Description 开始游戏接口
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Success 200 {object} data.Response{data=data.StartGameResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/start [post]
func StartGame(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	res := service.StartGame(openId)
	c.JSON(http.StatusOK, res)
}

// @Summary 玩家能量查询接口
// @Description 玩家能量查询接口
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Success 200 {object} data.Response{data=data.EnergyResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/energy [get]
func MyEnergy(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	res := service.MyEnergy(openId)
	c.JSON(http.StatusOK, res)
}

// @Summary 玩家增加体力接口
// @Description 玩家增加体力接口
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Success 200 {object} data.Response "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/health [put]
func AddHealth(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	res := service.AddHealth(openId)
	c.JSON(http.StatusOK, res)
}

// @Summary 玩家购买道具接口
// @Description 玩家购买道具接口
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.AddItemReq true "请求参数"
// @Success 200 {object} data.Response "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/item [post]
func AddItem(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.AddItemReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.AddItem(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家分享拿双倍接口
// @Description 玩家分享拿双倍接口
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.ShareReq true "请求参数"
// @Success 200 {object} data.Response{data=data.ShareResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/share [post]
func Share(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.ShareReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.Share(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家队伍查询接口$
// @Description 玩家队伍查询接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request query req.TeamGetReq true "请求参数"
// @Success 200 {object} data.Response{data=data.GetTeamResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/team [get]
func MyTeam(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.TeamGetReq
	if err := c.ShouldBindQuery(&param); err == nil {
		res := service.GetTeam(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家对局结果查询接口$
// @Description 玩家对局结果查询接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request query req.ResultGetReq true "请求参数"
// @Success 200 {object} data.Response{data=data.GameResultResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/result [get]
func MyResult(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.ResultGetReq
	if err := c.ShouldBindQuery(&param); err == nil {
		res := service.GetResult(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家对局金币接口$
// @Description 玩家对局金币接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request query req.CoinGetReq true "请求参数"
// @Success 200 {object} data.Response{data=int} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/coin [get]
func MyCoin(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.CoinGetReq
	if err := c.ShouldBindQuery(&param); err == nil {
		res := service.GetCoin(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家保存对局金币接口$
// @Description 玩家保存对局金币接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request query req.CoinSetReq true "请求参数"
// @Success 200 {object} data.Response "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/coin [post]
func PostCoin(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.CoinSetReq
	if err := c.ShouldBindQuery(&param); err == nil {
		res := service.SaveCoin(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家保存队伍并对战接口$
// @Description 玩家保存队伍并对战接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.PlayerTeamReq true "请求参数"
// @Success 200 {object} data.Response{data=data.Player} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/team [post]
func PostTeam(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.PlayerTeamReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.SaveTeam(&param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家保存对局结果接口$
// @Description 玩家保存对局结果接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.PlayerResultReq true "请求参数"
// @Success 200 {object} data.Response{data=data.ResultResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/result [post]
func PostResult(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.PlayerResultReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.GameResult(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家查询个人信息接口$
// @Description 玩家查询个人信息接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Success 200 {object} data.Response{data=data.GameInfoResp} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/info [get]
func GetInfo(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	res := service.GetInfo(openId)
	c.JSON(http.StatusOK, res)
}

// @Summary 玩家查询能量恢复时间接口$
// @Description 玩家查询能量恢复时间接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Success 200 {object} data.Response{data=uint} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/rectime [get]
func GetRecTime(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	recTime := service.GetRectime()
	c.JSON(http.StatusOK, recTime)
}

// @Summary 玩家查询头像列表接口$
// @Description 玩家查询头像列表接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request query req.PageReq true "请求参数"
// @Success 200 {object} data.Response{data=data.Pagination{items=repo.RoleFlag}} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/avatars [get]
func GetAvatars(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.AvatarRoleListReq
	if err := c.ShouldBindQuery(&param); err == nil {
		res := service.GetAvatars(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家保存广告观看记录接口$
// @Description 玩家保存广告观看记录接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.LaAdReq true "请求参数"
// @Success 200 {object} data.Response "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/ad [post]
func BackpackAd(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.LaAdReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.BackpackAd(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家获取头像接口$
// @Description 玩家获取头像接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.PostAvatarReq true "请求参数"
// @Success 200 {object} data.Response "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/avatar [post]
func PostAvatar(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.PostAvatarReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.PostAvatar(openId, &param)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// @Summary 玩家查询是否首次修改昵称接口$
// @Description 玩家查询是否首次修改昵称接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Success 200 {object} data.Response{data=bool} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/nickname [get]
func CheckNickname(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	res := service.CheckNickname(openId)
	c.JSON(http.StatusOK, res)
}

// @Summary 玩家修改昵称接口$
// @Description 玩家修改昵称接口$
// @Tags 玩家
// @Accept json
// @Produce json
// @Param Authorization header string true "此接口需要token"
// @Param request body req.PutNickNameReq true "请求参数"
// @Success 200 {object} data.Response{data=bool} "成功返回"
// @Failure 400 {object} data.Response "失败返回"
// @Router /api/v1/player/nickname [put]
func PutNickname(c *gin.Context) {
	openId := c.GetString("username")
	if openId == "" {
		c.JSON(http.StatusOK, data.NewErrorResponse(30002, "登录状态异常"))
	}
	var param req.PutNickNameReq
	if err := c.ShouldBind(&param); err == nil {
		res := service.UpdateNickname(openId, param.NickName)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}
