import { useUserStore } from '@/pinia/modules/user'
import { ElLoading } from 'element-plus'
import { userRegister } from '@/plugin/register/api/api'
import { defineStore } from 'pinia'

// 获取原始的 useUserStore 模块
const originalUserStore = useUserStore()

// 扩展原始的 useUserStore 模块
export const extendedUseUserStore = defineStore({
  id: 'user',
  state: () => ({
    userInfo: null,
    token: null
  }),
  actions: {
    async Register(loginInfo) {
      const loadingInstance = ElLoading.service({
        fullscreen: true,
        text: "注册中，请稍候...",
      });
      try {
        const res = await userRegister(loginInfo);
        if (res.code === 0) {
          this.userInfo = res.data.user;
          this.token = res.data.token;
          const routerStore = this.useRouterStore();
          await routerStore.SetAsyncRouter();
          const asyncRouters = routerStore.asyncRouters;
          asyncRouters.forEach((asyncRouter) => {
            router.addRoute(asyncRouter);
          });
          router.push({ name: this.userInfo.authority.defaultRouter });
          return true;
        }
      } catch (e) {
        loadingInstance.close();
      }
      loadingInstance.close();
    }
  }
}, () => originalUserStore)