import { useUserStore } from '@/pinia/modules/user'
import { ElLoading } from 'element-plus'
import { userRegister } from '@/plugin/register/api/api'

// 获取原始的 useUserStore 模块
const originalUserStore = useUserStore()

// 扩展原始的 useUserStore 模块
export const extendedUseUserStore = originalUserStore.extend({
  Register: async (loginInfo) => {
    const loadingInstance = ElLoading.service({
        fullscreen: true,
        text: "注册中，请稍候...",
    });
    try {
        const res = await userRegister(loginInfo);
        if (res.code === 0) {
            originalUserStore.userInfo = res.data.user;
            originalUserStore.token = res.data.token;
            const routerStore = originalUserStore.useRouterStore();
            await routerStore.SetAsyncRouter();
            const asyncRouters = routerStore.asyncRouters;
            asyncRouters.forEach((asyncRouter) => {
                router.addRoute(asyncRouter);
            });
            router.push({name: originalUserStore.userInfo.authority.defaultRouter});
            return true;
        }
    } catch (e) {
        loadingInstance.close();
    }
    loadingInstance.close();
  }
})
