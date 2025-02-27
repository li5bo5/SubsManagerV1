<!-- 节点测试页面 -->
<template>
  <div class="node-test">
    <!-- 操作区域 -->
    <el-card class="operation-section">
      <template #header>
        <div class="section-header">
          <span>节点操作</span>
          <div class="header-actions">
            <el-button 
              type="primary" 
              :loading="importing" 
              :disabled="testing"
              @click="handleImport"
            >
              导入节点
            </el-button>
            <el-button 
              type="success" 
              :loading="testing" 
              :disabled="importing || !nodes.length"
              @click="handleTest"
            >
              开始测速
            </el-button>
          </div>
        </div>
      </template>

      <!-- 测试进度 -->
      <div v-if="testing" class="progress-section">
        <el-progress 
          :percentage="testProgress" 
          :status="testProgress === 100 ? 'success' : ''"
        >
          <template #default="{ percentage }">
            <span class="progress-label">
              {{ percentage === 100 ? '测试完成' : `正在测试: ${percentage}%` }}
            </span>
          </template>
        </el-progress>
        <div class="progress-stats">
          <el-tag>总节点: {{ totalNodes }}</el-tag>
          <el-tag type="success">已测试: {{ testedNodes }}</el-tag>
          <el-tag type="warning">待测试: {{ remainingNodes }}</el-tag>
        </div>
      </div>
    </el-card>

    <!-- 节点列表 -->
    <el-card class="node-list-section">
      <template #header>
        <div class="section-header">
          <span>节点列表</span>
          <div class="list-stats">
            <el-tag>总数: {{ nodes.length }}</el-tag>
            <el-tag type="success">可用: {{ availableNodes }}</el-tag>
            <el-tag type="danger">超时: {{ timeoutNodes }}</el-tag>
          </div>
        </div>
      </template>

      <el-table
        v-loading="importing"
        :data="nodes"
        style="width: 100%"
        :default-sort="{ prop: 'delay', order: 'ascending' }"
      >
        <el-table-column prop="alias" label="节点名称" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="node-info">
              <el-tag size="small" :type="getNodeTypeTag(row.type)">{{ row.type }}</el-tag>
              <span>{{ row.alias }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="服务器" min-width="150" show-overflow-tooltip />
        <el-table-column prop="port" label="端口" width="100" align="center" />
        <el-table-column prop="delay" label="延迟" width="120" align="center" sortable>
          <template #default="{ row }">
            <el-tag 
              :type="getDelayTag(row.delay)" 
              :effect="row.delay ? 'light' : 'plain'"
            >
              {{ formatDelay(row.delay) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="speed" label="速度" width="120" align="center" sortable>
          <template #default="{ row }">
            <el-tag 
              :type="getSpeedTag(row.speed)"
              :effect="row.speed ? 'light' : 'plain'"
            >
              {{ formatSpeed(row.speed) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="protocol" label="协议" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ row.protocol }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useStore } from '@/stores'
import type { Node } from '@/types'

// 状态管理
const store = useStore()

// 节点列表
const nodes = ref<Node[]>([])

// 加载状态
const importing = ref(false)
const testing = ref(false)

// 测试进度
const totalNodes = ref(0)
const testedNodes = ref(0)
const testProgress = computed(() => {
  if (!totalNodes.value) return 0
  return Math.round((testedNodes.value / totalNodes.value) * 100)
})
const remainingNodes = computed(() => totalNodes.value - testedNodes.value)

// 节点统计
const availableNodes = computed(() => nodes.value.filter(node => node.delay && node.delay < 1000).length)
const timeoutNodes = computed(() => nodes.value.filter(node => !node.delay || node.delay >= 1000).length)

// 导入节点
const handleImport = async () => {
  try {
    importing.value = true
    await store.importNodes()
    await fetchNodes()
    ElMessage.success('节点导入成功')
  } catch (error: any) {
    ElMessage.error(error.message || '节点导入失败')
  } finally {
    importing.value = false
  }
}

// 开始测速
const handleTest = async () => {
  try {
    testing.value = true
    totalNodes.value = nodes.value.length
    testedNodes.value = 0

    // 启动测速
    await store.testNodes()
    
    // 更新节点列表
    await fetchNodes()
    ElMessage.success('节点测试完成')
  } catch (error: any) {
    ElMessage.error(error.message || '节点测试失败')
  } finally {
    testing.value = false
    totalNodes.value = 0
    testedNodes.value = 0
  }
}

// 获取节点列表
const fetchNodes = async () => {
  try {
    const { data } = await store.getNodeList()
    nodes.value = data
  } catch (error: any) {
    ElMessage.error(error.message || '获取节点列表失败')
  }
}

// 格式化延迟
const formatDelay = (delay?: number) => {
  if (!delay) return '超时'
  return `${delay}ms`
}

// 格式化速度
const formatSpeed = (speed?: number) => {
  if (!speed) return '未测试'
  if (speed < 1024) return `${speed.toFixed(1)}KB/s`
  return `${(speed / 1024).toFixed(1)}MB/s`
}

// 获取节点类型标签样式
const getNodeTypeTag = (type: string) => {
  const typeMap: Record<string, string> = {
    vmess: '',
    shadowsocks: 'success',
    hysteria2: 'warning',
    trojan: 'danger'
  }
  return typeMap[type.toLowerCase()] || 'info'
}

// 获取延迟标签样式
const getDelayTag = (delay?: number) => {
  if (!delay) return 'danger'
  if (delay < 100) return 'success'
  if (delay < 300) return ''
  if (delay < 500) return 'warning'
  return 'danger'
}

// 获取速度标签样式
const getSpeedTag = (speed?: number) => {
  if (!speed) return 'info'
  if (speed > 10 * 1024) return 'success' // 10MB/s
  if (speed > 5 * 1024) return '' // 5MB/s
  if (speed > 1024) return 'warning' // 1MB/s
  return 'danger'
}

// 页面加载时获取节点列表
onMounted(() => {
  fetchNodes()
})
</script>

<style scoped lang="scss">
.node-test {
  .operation-section,
  .node-list-section {
    margin-bottom: 20px;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 500;

    .header-actions {
      display: flex;
      gap: 10px;
    }

    .list-stats {
      display: flex;
      gap: 10px;
    }
  }

  .progress-section {
    padding: 20px 0;

    .progress-label {
      font-size: 14px;
      color: #606266;
    }

    .progress-stats {
      margin-top: 10px;
      display: flex;
      gap: 10px;
      justify-content: center;
    }
  }

  .node-info {
    display: flex;
    align-items: center;
    gap: 8px;

    .el-tag {
      text-transform: uppercase;
    }
  }

  :deep(.el-progress-bar__inner) {
    transition: width 0.3s ease;
  }

  :deep(.el-tag) {
    text-transform: none;
  }
}
</style> 