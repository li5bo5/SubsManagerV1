import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 布局组件
import MainLayout from '@/layouts/MainLayout.vue'

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Status',
        component: () => import('@/views/status/index.vue'),
        meta: {
          title: '状态监控',
          keepAlive: true
        }
      },
      {
        path: 'subscription',
        name: 'Subscription',
        component: () => import('@/views/subscription/index.vue'),
        meta: {
          title: '订阅管理',
          keepAlive: true
        }
      },
      {
        path: 'node-test',
        name: 'NodeTest',
        component: () => import('@/views/node-test/index.vue'),
        meta: {
          title: '节点测试',
          keepAlive: true
        }
      },
      {
        path: 'node-select',
        name: 'NodeSelect',
        component: () => import('@/views/node-select/index.vue'),
        meta: {
          title: '优选节点',
          keepAlive: true
        }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/index.vue'),
        meta: {
          title: '系统设置',
          keepAlive: true
        }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('@/views/logs/index.vue'),
        meta: {
          title: '运行日志',
          keepAlive: true
        }
      }
    ]
  },
  // 404页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue')
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  // 平滑滚动
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || '未知页面'} - SubsManager`
  next()
})

export default router 