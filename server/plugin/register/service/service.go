package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	gvaGloval "github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	plugGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/register/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/telegram_bot/service"
	userService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/mojocn/base64Captcha"
)

type RegisterService struct{}

var telegramServer = new(service.TelegramBotService)

func (e *RegisterService) Code(tgid string) (err error) {
	// 制作四位数code
	code := utils.RandomString(plugGlobal.GlobalConfig.CodeLength)
	// 发送code
	err = service.ServiceGroupApp.SendTgMessage(plugGlobal.GlobalConfig.TgBotToken, tgid, fmt.Sprintf("注册验证码：<code>%v</code>", code), "html")
	if err != nil {
		return err
	}
	// 存储code
	ctx := context.Background()
	gvaGloval.GVA_REDIS.Set(ctx, tgid, code, 5*time.Minute)
	return nil
}

func (e *RegisterService) Register(register model.RegisterReq) (err error) {
	// 检测tgcode是否正确
	ctx := context.Background()
	code, err := gvaGloval.GVA_REDIS.Get(ctx, register.Tgid)
	if register.Code != code {
		return errors.New("验证码错误")
	} else if err != nil {
		return errors.New(fmt.Sprintf("存储的TG验证码获取错误：%v", err))
	}
	// 获取注册的信息
	if err := utils.Verify(register, utils.LoginVerify); err != nil {
		return errors.New(fmt.Sprintf("获取登录状态错误：%v", err))
	}
	var (
		store = base64Captcha.DefaultMemStore
		user  system.SysUser
		us    *userService.UserService
	)
	if !store.Verify(register.CaptchaId, register.Captcha, true) {
		return errors.New(fmt.Sprintf("图片验证码错误"))
	}
	u := &system.SysUser{Username: register.Username, Password: register.Password, Tgid: register.Tgid}
	// 检测账户是否存在
	err = gvaGloval.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		return errors.New(fmt.Sprintf("用户名已注册：%v", err))
	}
	if user.Username != "" && user.Password != "" {
		return errors.New(fmt.Sprintf("用户名已注册：%v", err))
	}
	// 默认用户结构体
	var sysAuthority systemReq.Register
	sysAuthority.Username = u.Username
	sysAuthority.NickName = u.NickName
	sysAuthority.Password = u.Password
	sysAuthority.Phone = u.Phone
	sysAuthority.AuthorityId = plugGlobal.GlobalConfig.AuthorityId
	sysAuthority.AuthorityIds = append(sysAuthority.AuthorityIds, plugGlobal.GlobalConfig.AuthorityId)
	// 因为上面定义过，且得到了数据库默认的值，所以直接使用
	user.Password = u.Password
	user.UUID, _ = uuid.NewV4()
	user.Username = u.Username
	user.NickName = u.Username
	user.Phone = u.Tgid
	user.AuthorityId = plugGlobal.GlobalConfig.AuthorityId
	for _, v := range sysAuthority.AuthorityIds {
		user.Authorities = append(user.Authorities, system.SysAuthority{
			AuthorityId: v,
			// 系统注册的时候有这个参数 DefaultRouter 用户登录后默认的router设置为dashboard，如果注册的用户首页不是后台，需要自行更改
			DefaultRouter: "dashboard",
		})
	}
	if _, err := us.Register(*u); err != nil {
		return errors.New(fmt.Sprintf("注册失败：%v", err))
	}
	if _, err := us.Login(u); err != nil {
		return errors.New(fmt.Sprintf("登录失败：%v", err))
	}
	return nil
}
