<!-- 状态监控页面 -->
<template>
  <div class="status page-container">
    <!-- 节点状态 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>节点状态</span>
          </div>
        </div>
      </template>

      <div class="node-stats">
        <div class="stat-item">
          <div class="stat-value">{{ nodeStats.total }}</div>
          <div class="stat-label">总节点数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ nodeStats.current }}</div>
          <div class="stat-label">当前节点</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ nodeStats.slow }}</div>
          <div class="stat-label">慢速节点</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ nodeStats.failed }}</div>
          <div class="stat-label">故障节点</div>
        </div>
      </div>
    </el-card>

    <!-- 当前状态 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>当前状态</span>
          </div>
        </div>
      </template>

      <div class="subscription-files">
        <div class="file-item">
          <div class="file-label">最终订阅文件</div>
          <div class="file-link">
            <el-link 
              type="primary" 
              :href="currentStatus.finalSubscription"
              target="_blank"
            >
              {{ currentStatus.finalSubscription }}
            </el-link>
          </div>
        </div>
        <div class="file-item">
          <div class="file-label">订阅整合文件</div>
          <div class="file-link">
            <el-link 
              type="primary" 
              :href="currentStatus.mergedSubscription"
              target="_blank"
            >
              {{ currentStatus.mergedSubscription }}
            </el-link>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 历史记录 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>历史记录</span>
          </div>
        </div>
      </template>

      <div class="history-list">
        <div 
          v-for="item in history" 
          :key="item.time" 
          class="history-item"
        >
          <div class="history-time">{{ formatTime(item.time) }}</div>
          <div class="history-url">
            <el-link 
              type="primary" 
              :href="item.url"
              target="_blank"
            >
              {{ item.url }}
            </el-link>
          </div>
          <div class="history-actions">
            <el-button 
              :icon="CopyDocument"
              circle
              @click="copyUrl(item.url)"
            />
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'
import { useStore } from '@/stores'
import type { SystemStatus } from '@/types'

// 状态管理
const store = useStore()

// 节点统计
const nodeStats = ref({
  total: 0,
  current: 0,
  slow: 0,
  failed: 0
})

// 当前状态
const currentStatus = ref({
  finalSubscription: '',
  mergedSubscription: ''
})

// 历史记录
const history = ref<SystemStatus['history']>([])

// 获取系统状态
const fetchStatus = async () => {
  try {
    await store.fetchSystemStatus()
    const status = store.systemStatus
    if (status) {
      nodeStats.value = status.nodeStats
      currentStatus.value = status.currentStatus
      history.value = status.history
    }
  } catch (error) {
    console.error('获取系统状态失败:', error)
  }
}

// 复制链接
const copyUrl = async (url: string) => {
  try {
    await navigator.clipboard.writeText(url)
    ElMessage.success('链接已复制')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 格式化时间
const formatTime = (time: string) => {
  return time.replace(/[T]/g, ' ').slice(5, 16)
}

// 页面加载时获取状态
onMounted(() => {
  fetchStatus()
})
</script>

<style scoped lang="scss">
@use '@/styles/common.scss';

.status {
  .node-stats {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    padding: 10px 0;

    @media screen and (max-width: 768px) {
      grid-template-columns: repeat(2, 1fr);
    }

    .stat-item {
      text-align: center;
      padding: 20px;
      background-color: var(--el-bg-color-page);
      border-radius: 4px;

      .stat-value {
        font-size: 24px;
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin-bottom: 8px;
      }

      .stat-label {
        font-size: 14px;
        color: var(--el-text-color-secondary);
      }
    }
  }

  .subscription-files {
    .file-item {
      margin-bottom: 15px;

      &:last-child {
        margin-bottom: 0;
      }

      .file-label {
        font-size: 14px;
        color: var(--el-text-color-secondary);
        margin-bottom: 5px;
      }

      .file-link {
        word-break: break-all;
      }
    }
  }

  .history-list {
    .history-item {
      display: flex;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid var(--el-border-color-lighter);

      &:last-child {
        border-bottom: none;
      }

      .history-time {
        width: 100px;
        flex-shrink: 0;
        color: var(--el-text-color-secondary);
      }

      .history-url {
        flex: 1;
        margin: 0 15px;
        word-break: break-all;
      }

      .history-actions {
        flex-shrink: 0;
      }
    }
  }
}
</style> 