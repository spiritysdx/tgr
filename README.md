# tgr

gva: 2.6.1

版本：v0.0.6

结合 Telegram 进行验证的注册/登录/找回密码插件

- [x] 基本的注册页面和注册逻辑
- [x] 用户需要使用TG进行验证再进行注册
- [x] 用户需要关注特定频道才能进行注册
- [x] 支持注册成功后自动登录
- [x] 支持用户通过TGID找回密码
- [x] 用户登录需要确保仍然关注了特定频道(测试用户和管理员无需校验)
- ~~查询用户是否为指定注册期限长度的账户，避免频繁删创号刷注册~~ 官方没有公布该API，仅提供应用程序版本的查询，暂无GO版本适配库

具体部署说明见 server/plugin/register 下的 README.md 说明

由于有TG发信和获取频道用户信息的需求，所以需要事先安装GVA插件市场中的```灰机消息发送插件```(也即是 https://github.com/spiritysdx/tgm )

用户需要事先私聊过机器人，否则机器人无法主动私聊用户，且注册用户在频道中不应该为管理员身份，仅为普通订阅者

界面展示

![screenshot-1710854022227](https://github.com/spiritysdx/tgr/assets/97792170/a90b138f-6d70-4485-84db-b5e87187c40a)
![screenshot-1710854046018](https://github.com/spiritysdx/tgr/assets/97792170/a1b6a231-c49e-43c5-9f83-145d0a05a208)
![screenshot-1710854032749](https://github.com/spiritysdx/tgr/assets/97792170/d24dc87f-7397-4f11-9f47-6c9405b3bdcc)
