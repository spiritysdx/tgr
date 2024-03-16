# 使用说明

拖拽插件压缩包进入插件安装页面进行安装再进行后续操作

由于有TG发信和获取频道用户信息的需求，所以需要事先安装GVA插件市场中的```灰机消息发送插件```

## 具体说明

### server

查看 ```server/initialize/plugin.go``` 文件中是否已注册插件，如若未注册，在

```
func InstallPlugin
```

函数中插入

```
  // 8881 为普通子用户ID，可自行更改替换注册的角色
  PluginInit(PublicGroup, register.CreateRegisterPlug("8881"), "tgbot的token", "验证码的长度", "频道的chat_id")
```

这里的bot务必已加入到频道中并给予了管理员权限

### web

#### 1

查看 ```web/src/permission.js``` 并参照

```
const whiteList = ['Admin', 'Login', 'Init', 'Register']
```

添加对应的路由名字

#### 2

查看 ```web/src/router/index.js``` 参照

```
const routes = [{
  path: '/',
  redirect: '/login'
},
{
  path: '/init',
  name: 'Init',
  component: () => import('@/view/init/index.vue')
},
{
  path: '/admin',
  name: 'Admin',
  component: () => import('@/view/login/index.vue')
},
{
  path: '/login',
  name: 'Login',
  component: () => import('@/plugin/register/view/index.vue')
},
{
  path: '/register',
  name: 'Register',
  component: () => import('@/plugin/register/view/index.vue')
},
{
  path: '/:catchAll(.*)',
  meta: {
    closeTab: true,
  },
  component: () => import('@/view/error/index.vue')
}
]
```

进行更改，定义用户进入的首页是插件中定义的首页，而不是官方GVA定义的初始化页面

## 后言

所有更改成功后务必保存所有文件再重启前后端，才能正常使用
