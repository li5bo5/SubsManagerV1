<template>
  <div class="status page-container">
    <!-- 系统状态 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>系统状态</span>
          </div>
          <div class="header-actions">
            <el-tooltip 
              content="刷新状态" 
              placement="top"
            >
              <el-button 
                :icon="Refresh"
                circle
                @click="refreshStatus"
              />
            </el-tooltip>
          </div>
        </div>
      </template>

      <div 
        v-loading="loading"
        class="content-section"
      >
        <el-descriptions :column="2" border>
          <el-descriptions-item label="运行状态">
            <el-tag 
              :type="status.running ? 'success' : 'danger'"
            >
              {{ status.running ? '运行中' : '已停止' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="运行时间">
            {{ formatUptime(status.uptime) }}
          </el-descriptions-item>
          <el-descriptions-item label="CPU 使用率">
            {{ status.cpuUsage }}%
          </el-descriptions-item>
          <el-descriptions-item label="内存使用">
            {{ formatMemory(status.memoryUsage) }}
          </el-descriptions-item>
          <el-descriptions-item label="订阅数量">
            {{ status.subscriptionCount }}
          </el-descriptions-item>
          <el-descriptions-item label="节点数量">
            {{ status.nodeCount }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>

    <!-- 代理状态 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>代理状态</span>
          </div>
        </div>
      </template>

      <div class="content-section">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="系统代理">
            <el-tag 
              :type="proxy.enabled ? 'success' : 'info'"
            >
              {{ proxy.enabled ? '已启用' : '未启用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="监听地址">
            {{ proxy.address }}:{{ proxy.port }}
          </el-descriptions-item>
          <el-descriptions-item label="当前节点">
            {{ proxy.currentNode || '未选择' }}
          </el-descriptions-item>
          <el-descriptions-item label="连接数">
            {{ proxy.connections }}
          </el-descriptions-item>
          <el-descriptions-item label="上传流量">
            {{ formatTraffic(proxy.uploadTraffic) }}
          </el-descriptions-item>
          <el-descriptions-item label="下载流量">
            {{ formatTraffic(proxy.downloadTraffic) }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 流量图表 -->
        <div class="traffic-chart">
          <div class="chart-title">实时流量</div>
          <div class="chart-container">
            <!-- 这里放置流量图表组件 -->
          </div>
        </div>
      </div>
    </el-card>

    <!-- 定时任务 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>定时任务</span>
          </div>
        </div>
      </template>

      <div class="content-section responsive-table">
        <el-table :data="tasks" style="width: 100%">
          <el-table-column 
            prop="name" 
            label="任务名称"
            min-width="120"
          />
          <el-table-column 
            prop="schedule" 
            label="执行计划"
            min-width="150"
          />
          <el-table-column 
            prop="lastRun" 
            label="上次执行"
            min-width="180"
          >
            <template #default="{ row }">
              {{ formatTime(row.lastRun) }}
            </template>
          </el-table-column>
          <el-table-column 
            prop="nextRun" 
            label="下次执行"
            min-width="180"
          >
            <template #default="{ row }">
              {{ formatTime(row.nextRun) }}
            </template>
          </el-table-column>
          <el-table-column 
            prop="status" 
            label="状态"
            width="100"
            align="center"
          >
            <template #default="{ row }">
              <el-tag :type="getTaskStatusType(row.status)">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { useStore } from '@/stores'

// 状态管理
const store = useStore()

// 加载状态
const loading = ref(false)

// 系统状态数据
const status = ref({
  running: false,
  uptime: 0,
  cpuUsage: 0,
  memoryUsage: 0,
  subscriptionCount: 0,
  nodeCount: 0
})

// 代理状态数据
const proxy = ref({
  enabled: false,
  address: '',
  port: 0,
  currentNode: '',
  connections: 0,
  uploadTraffic: 0,
  downloadTraffic: 0
})

// 定时任务数据
const tasks = ref([])

// 刷新状态
const refreshStatus = async () => {
  try {
    loading.value = true
    await store.fetchSystemStatus()
    // 更新状态数据
    // ...
  } catch (error) {
    console.error('刷新状态失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化函数
const formatUptime = (seconds: number) => {
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${days}天 ${hours}小时 ${minutes}分钟`
}

const formatMemory = (bytes: number) => {
  const mb = bytes / 1024 / 1024
  return `${mb.toFixed(1)} MB`
}

const formatTraffic = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`
  return `${(bytes / 1024 / 1024 / 1024).toFixed(1)} GB`
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

const getTaskStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    running: 'primary',
    success: 'success',
    failed: 'danger',
    pending: 'info'
  }
  return typeMap[status.toLowerCase()] || ''
}

// 页面加载时获取状态
onMounted(() => {
  refreshStatus()
})
</script>

<style scoped lang="scss">
@use '@/styles/common.scss';

.status {
  :deep(.el-descriptions) {
    margin-bottom: 30px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .traffic-chart {
    margin-top: 30px;

    .chart-title {
      font-size: 16px;
      font-weight: 500;
      color: #303133;
      margin-bottom: 20px;
    }

    .chart-container {
      height: 300px;
      background-color: #f5f7fa;
      border-radius: 4px;
      display: flex;
      justify-content: center;
      align-items: center;
      color: #909399;
    }
  }

  :deep(.el-tag) {
    text-transform: none;
  }
}
</style> 