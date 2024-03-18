package api

import (
	systemApi "github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RegisterApi struct{}

// @Tags code
// @Summary 发送code
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /register/code [post]
func (p *RegisterApi) Code(c *gin.Context) {
	var req model.CodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("获取tg_id失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.Code(req.TgId); err != nil {
		global.GVA_LOG.Error("发送Code失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed("", "发送Code成功", c)
	}
}

// @Tags Register
// @Summary 注册用户
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /register/register [post]
func (p *RegisterApi) Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if res, err := service.ServiceGroupApp.Register(req); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(*res, "注冊成功", c)
	}
}
