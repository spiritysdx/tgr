import { ElLoading } from 'element-plus'
import { userRegister } from '@/api/user'
import router from '@/router'

const Register = async(loginInfo) => {
  loadingInstance.value = ElLoading.service({
    fullscreen: true,
    text: '注册中，请稍候...',
  })
  try {
    const res = await userRegister(loginInfo)
    if (res.code === 0) {
      setUserInfo(res.data.user)
      setToken(res.data.token)
      const routerStore = useRouterStore()
      await routerStore.SetAsyncRouter()
      const asyncRouters = routerStore.asyncRouters
      asyncRouters.forEach((asyncRouter) => {
        router.addRoute(asyncRouter)
      })
      router.push({ name: userInfo.value.authority.defaultRouter })
      return true
    }
  } catch (e) {
    loadingInstance.value.close()
  }
  loadingInstance.value.close()
}