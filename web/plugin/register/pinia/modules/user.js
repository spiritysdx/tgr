import { useUserStore } from '@/pinia/modules/user'
import { ElLoading, ElMessage } from 'element-plus'
import { userRegister } from '@/plugin/register/api/api'
import { defineStore } from 'pinia'

export const extendedUseUserStore = defineStore({
  id: 'extendedUser',
  state: () => ({
    originalUserStore: useUserStore(),
    loadingInstance: null
  }),
  actions: {
    async Register(loginInfo) {
      this.loadingInstance = ElLoading.service({
        fullscreen: true,
        text: "注册中，请稍候...",
      });
      try {
        const res = await userRegister(loginInfo);
        if (res.code === 0) {
          this.originalUserStore.userInfo = res.data.user;
          this.originalUserStore.token = res.data.token;
          const routerStore = this.originalUserStore.useRouterStore();
          await routerStore.SetAsyncRouter();
          const asyncRouters = routerStore.asyncRouters;
          asyncRouters.forEach((asyncRouter) => {
            router.addRoute(asyncRouter);
          });
          router.push({ name: this.originalUserStore.userInfo.authority.defaultRouter });
          return true;
        }
      } catch (e) {
        this.loadingInstance.close();
      }
      this.loadingInstance.close();
    }
  }
})
