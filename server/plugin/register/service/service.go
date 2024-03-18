package service

import (
	"context"
	"errors"
	"fmt"
	gvaGlobal "github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	plugGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/register/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/telegram_bot/service"
	userService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/mojocn/base64Captcha"
	"time"
)

type RegisterService struct{}

func (e *RegisterService) Code(tgid string) (err error) {
	// 制作四位数code
	code := utils.RandomString(plugGlobal.GlobalConfig.CodeLength)
	// 发送code
	_, err = service.ServiceGroupApp.SendTgMessage(plugGlobal.GlobalConfig.TgBotToken, tgid,
		fmt.Sprintf("注册验证码：<code>%v</code>", code), "html")
	if err != nil {
		return errors.New(fmt.Sprintf("发送TG验证码错误：%v", err))
	}
	// 存储code
	ctx := context.Background()
	gvaGlobal.GVA_REDIS.Set(ctx, tgid, code, 5*time.Minute)
	return nil
}

func (e *RegisterService) Register(register model.RegisterReq) (res *system.SysUser, err error) {
	res = &system.SysUser{}
	// 检测tgcode是否正确
	ctx := context.Background()
	code, err := gvaGlobal.GVA_REDIS.Get(ctx, register.Tgid).Result()
	if err != nil {
		return res, errors.New(fmt.Sprintf("存储的TG验证码获取错误：%v", err))
	} else if register.Code != code {
		return res, errors.New(fmt.Sprintf("验证码填写错误：%v", register.Code))
	}
	// 检测用户是否在特定的频道中
	_, err = service.ServiceGroupApp.IsTgMember(plugGlobal.GlobalConfig.TgBotToken, register.Tgid,
		plugGlobal.GlobalConfig.ChannelId)
	if err != nil {
		return res, errors.New(fmt.Sprintf("检测是否在频道错误：%v", err))
	}
	// 获取注册的信息
	if err := utils.Verify(register, utils.LoginVerify); err != nil {
		return res, errors.New(fmt.Sprintf("获取登录状态错误：%v", err))
	}
	var (
		store = base64Captcha.DefaultMemStore
		user  system.SysUser
		us    *userService.UserService
	)
	if !store.Verify(register.CaptchaId, register.Captcha, true) {
		return res, errors.New(fmt.Sprintf("图片验证码错误"))
	}
	// 加密密码
	plaintext_password := register.Password
	register.Password = utils.BcryptHash(register.Password)
	// 创建用户需要传入的信息
	// 用 Phone 字段存用户的 TGID 了
	UUID, _ := uuid.NewV4()
	u := &system.SysUser{UUID: UUID, Username: register.Username, Password: register.Password, NickName: "注册用户", Phone: register.Tgid, AuthorityId: plugGlobal.GlobalConfig.AuthorityId}
	// 检测传入信息是否为空
	if u.Username == "" {
		return res, errors.New(fmt.Sprintf("用户名为空：%v", err))
	}
	if u.Password == "" {
		return res, errors.New(fmt.Sprintf("密码为空：%v", err))
	}
	// 检测账户是否存在
	err = gvaGlobal.GVA_DB.Where("phone = ?", u.Phone).First(&user).Error
	if err == nil {
		return res, errors.New(fmt.Sprintf("该TGID已注册，用户名为：%v", user.Username))
	}
	// 创建用户账户
	err = gvaGlobal.GVA_DB.Create(&u).Error
	if err != nil {
		return res, errors.New(fmt.Sprintf("注册账户失败：%v", err))
	}
	// 创建用户角色
	gvaGlobal.GVA_DB.Where("phone = ?", u.Phone).First(&user) // 查询对应tgid的用户id存到user里
	a := &system.SysUserAuthority{SysUserId: user.ID, SysAuthorityAuthorityId: plugGlobal.GlobalConfig.AuthorityId}
	err = gvaGlobal.GVA_DB.Create(&a).Error
	if err != nil {
		return res, errors.New(fmt.Sprintf("注册角色失败：%v", err))
	}
	// 创建完毕后密码需要改回明文密码再登录
	u.Password = plaintext_password
	if _, err := us.Login(u); err != nil {
		return res, errors.New(fmt.Sprintf("登录失败：%v", err))
	}
	return res, nil
}
