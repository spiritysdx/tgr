# 使用说明

拖拽插件压缩包进入插件安装页面进行安装再进行后续操作

## 具体说明

### server

查看 ```server/initialize/plugin.go``` 文件中是否已注册插件，如若未注册，在

```
func InstallPlugin
```

函数中插入

```
  // 8881 为普通子用户ID，可自行更改替换注册的角色
  PluginInit(PublicGroup, register.CreateRegisterPlug("8881"))
```

### web

#### 1

查看 ```web/src/permission.js``` 并参照

```
const whiteList = ['Admin', 'Login', 'Init', 'Register']
```

添加对应的路由名字

#### 2

查看 ```web/src/modules/user.js``` 并参照

```
  const Register = async (loginInfo) => {
    loadingInstance.value = ElLoading.service({
        fullscreen: true,
        text: "注册中，请稍候...",
    });
    try {
        const res = await userRegister(loginInfo);
        if (res.code === 0) {
            setUserInfo(res.data.user);
            setToken(res.data.token);
            const routerStore = useRouterStore();
            await routerStore.SetAsyncRouter();
            const asyncRouters = routerStore.asyncRouters;
            asyncRouters.forEach((asyncRouter) => {
                router.addRoute(asyncRouter);
            });
            router.push({name: userInfo.value.authority.defaultRouter});
            return true;
        }
    } catch (e) {
        loadingInstance.value.close();
    }
    loadingInstance.value.close();
  };

  return {
    Register,
    // 后续有若干原有配置的内容
  }
```

在对应位置添加上述内容，注意return只是多添加了一个返回，不是写多一个return

#### 3

查看 ```web/src/api/user.js``` 在文件最后添加

```
// web/src/api/user.js
// @Summary 用户注册
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /register [post]
export const userRegister = (data) => {
  return service({
    url: '/register',
    method: 'post',
    data: data,
  })
}
```

#### 4

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

进行更改，确保用户进入的首页是插件中定义的首页，而不是官方GVA定义的登录页

## 后言

所有更改成功后务必保存所有文件再重启前后端，才能正常使用
