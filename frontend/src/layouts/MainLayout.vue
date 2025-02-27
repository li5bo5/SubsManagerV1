<!-- MainLayout.vue -->
<template>
  <el-container class="layout-container">
    <!-- 头部区域 -->
    <el-header class="header">
      <div class="logo">
        <img src="@/assets/logo.svg" alt="Logo" class="logo-img" />
        <h1>多订阅管理工具（SubsManager）</h1>
      </div>
      <div class="github">
        <a href="https://github.com/li5bo5" target="_blank" rel="noopener noreferrer">
          <el-icon><svg-icon name="github" /></el-icon>
          PowerBy: Li5Bo5
        </a>
      </div>
    </el-header>

    <el-container class="main-container">
      <!-- 侧边栏 -->
      <el-aside width="200px" class="sidebar">
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          router
          :collapse="isCollapse"
        >
          <el-menu-item index="/">
            <el-icon><Monitor /></el-icon>
            <template #title>状态监控</template>
          </el-menu-item>
          <el-menu-item index="/subscription">
            <el-icon><Document /></el-icon>
            <template #title>订阅管理</template>
          </el-menu-item>
          <el-menu-item index="/node-test">
            <el-icon><Connection /></el-icon>
            <template #title>节点测试</template>
          </el-menu-item>
          <el-menu-item index="/node-select">
            <el-icon><Select /></el-icon>
            <template #title>优选节点</template>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <template #title>系统设置</template>
          </el-menu-item>
          <el-menu-item index="/logs">
            <el-icon><List /></el-icon>
            <template #title>运行日志</template>
          </el-menu-item>
        </el-menu>
        <div class="version">SubsManager-V1</div>
      </el-aside>

      <!-- 主要内容区域 -->
      <el-main class="main-content">
        <!-- 面包屑导航 -->
        <el-breadcrumb class="breadcrumb">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>{{ currentMenuTitle }}</el-breadcrumb-item>
        </el-breadcrumb>

        <!-- 路由视图 -->
        <div class="content-wrapper">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Monitor,
  Document,
  Connection,
  Select,
  Setting,
  List
} from '@element-plus/icons-vue'

const route = useRoute()
const isCollapse = ref(false)
const activeMenu = computed(() => route.path)

// 计算当前菜单标题
const currentMenuTitle = computed(() => {
  const menuMap = {
    '/': '状态监控',
    '/subscription': '订阅管理',
    '/node-test': '节点测试',
    '/node-select': '优选节点',
    '/settings': '系统设置',
    '/logs': '运行日志'
  }
  return menuMap[route.path] || '未知页面'
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  background-color: #f5f7fa;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
  height: 60px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.logo {
  display: flex;
  align-items: center;
}

.logo-img {
  height: 32px;
  margin-right: 10px;
}

.logo h1 {
  margin: 0;
  font-size: 18px;
  color: #303133;
  font-weight: 600;
}

.github {
  display: flex;
  align-items: center;
}

.github a {
  display: flex;
  align-items: center;
  color: #606266;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

.github a:hover {
  color: #409eff;
}

.github .el-icon {
  margin-right: 5px;
  font-size: 18px;
}

.main-container {
  height: calc(100vh - 60px);
}

.sidebar {
  background-color: #fff;
  border-right: 1px solid #dcdfe6;
  position: relative;
}

.el-menu-vertical {
  border-right: none;
}

.version {
  position: absolute;
  bottom: 20px;
  left: 0;
  width: 100%;
  text-align: center;
  color: #909399;
  font-size: 12px;
}

.main-content {
  padding: 20px;
  overflow-y: auto;
}

.breadcrumb {
  margin-bottom: 20px;
  padding: 8px 0;
}

.content-wrapper {
  background-color: #fff;
  border-radius: 4px;
  padding: 20px;
  min-height: calc(100vh - 180px);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 0 10px;
  }

  .logo h1 {
    font-size: 16px;
  }

  .main-content {
    padding: 10px;
  }

  .content-wrapper {
    padding: 15px;
  }
}
</style> 