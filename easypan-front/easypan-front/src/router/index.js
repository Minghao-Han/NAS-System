import { createRouter, createWebHistory } from 'vue-router'
import VueCookies from 'vue-cookies'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: '登录',
      component: () => import("@/views/Login.vue")
    },

    {
      path: "/",
      component: () => import("@/views/Framework.vue"),
      children: [
        {
          path: '/',
          redirect: "/login",
          //component: () => import("@/views/main/Main.vue")
        },
        {
          path: 'main/:category',
          name: '首页',
          meta: {
            needLogin: true,
            menuCode: "main"
          },
          component: () => import("@/views/main/Main.vue")
        },
        {
          path: '/myshare',
          name: '我的分享',
          meta: {
            needLogin: true,
            menuCode: "share"
          },
          component: () => import("@/views/share/Share.vue")
        },
        {
          path: '/recycle',
          name: '回收站',
          meta: {
            needLogin: true,
            menuCode: "recycle"
          },
          component: () => import("@/views/recycle/Recycle.vue")
        },
        {
          path: '/settings/sysSetting',
          name: '系统设置',
          meta: {
            needLogin: true,
            menuCode: "settings"
          },
          component: () => import("@/views/admin/SysSettings.vue")
        },
        {
          path: '/settings/userList',
          name: '用户管理',
          meta: {
            needLogin: true,
            menuCode: "settings"
          },
          component: () => import("@/views/admin/UserList.vue")
        },
        {
          path: '/settings/fileList',
          name: '用户文件',
          meta: {
            needLogin: true,
            menuCode: "settings"
          },
          component: () => import("@/views/admin/FileList.vue")
        },
      ]
    },
    {
      path: '/shareCheck/:shareId',
      name: '分享校验',
      component: () => import("@/views/webshare/ShareCheck.vue")
    },
    {
      path: '/share/:shareId',
      name: '分享',
      component: () => import("@/views/webshare/Share.vue")
    },
    // {
    //   path: '/qqlogincalback',
    //   name: "qq登录回调",
    //   component: () => import('@/views/QqLoginCallback.vue'),
    // }
  ]
})

router.beforeEach((to, from, next) => {
  if(to.path==='/login') return next()
  // const tokenstr = localStorage.getItem("token");
  // if (!tokenstr) return next('/login')
  next()
  //next('/login')
})

export default router
