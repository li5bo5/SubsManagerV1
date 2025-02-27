<!-- 优选节点页面 -->
<template>
  <div class="node-select page-container">
    <!-- 筛选条件 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>筛选条件</span>
          </div>
        </div>
      </template>

      <div class="form-group">
        <el-form :model="form" label-width="120px">
          <el-form-item label="延迟上限(ms)">
            <el-input-number 
              v-model="form.maxDelay" 
              :min="0" 
              :max="5000"
              :step="50"
              placeholder="请输入延迟上限"
            />
          </el-form-item>
          <el-form-item label="速度下限(M/s)">
            <el-input-number 
              v-model="form.minSpeed" 
              :min="0" 
              :max="100"
              :step="0.5"
              placeholder="请输入速度下限"
            />
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="filtering"
              :disabled="!form.maxDelay || !form.minSpeed"
              @click="handleFilter"
            >
              开始筛选
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 优选结果 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>优选结果</span>
            <div class="list-stats">
              <el-tag>总数: {{ totalNodes }}</el-tag>
              <el-tag type="success">可用: {{ availableNodes }}</el-tag>
              <el-tag type="danger">超时: {{ timeoutNodes }}</el-tag>
            </div>
          </div>
        </div>
      </template>

      <div class="content-section responsive-table">
        <el-table
          v-loading="loading"
          :data="currentPageNodes"
          style="width: 100%"
          :default-sort="{ prop: 'delay', order: 'ascending' }"
        >
          <el-table-column type="index" width="60" align="center" />
          <el-table-column prop="type" label="类型" width="100" align="center">
            <template #default="{ row }">
              <el-tag size="small" :type="getNodeTypeTag(row.type)">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="alias" label="别名" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="node-info">
                <span>{{ row.alias }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="address" label="地址" min-width="150" show-overflow-tooltip />
          <el-table-column prop="port" label="端口" width="100" align="center" />
          <el-table-column prop="protocol" label="传输协议" width="120" align="center">
            <template #default="{ row }">
              <el-tag size="small" effect="plain">{{ row.protocol }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="group" label="订阅分组" width="120" align="center" show-overflow-tooltip>
            <template #default="{ row }">
              <el-tag size="small" effect="plain" type="info">{{ row.group }}</el-tag>
            </template>
          </el-table-column>
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
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <span class="total-text">节点总数: {{ totalNodes }}</span>
          <el-select v-model="pageSize" class="page-size-select">
            <el-option
              v-for="size in [10, 20, 50, 100]"
              :key="size"
              :label="`${size}条/页`"
              :value="size"
            />
          </el-select>
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="totalNodes"
            :page-sizes="[10, 20, 50, 100]"
            layout="prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <!-- 优选订阅 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>优选订阅</span>
          </div>
        </div>
      </template>

      <div class="content-section">
        <el-button 
          type="primary" 
          :disabled="!filteredNodes.length"
          :loading="generating"
          @click="handleGenerateSubscription"
        >
          生成订阅
        </el-button>
        <div v-if="subscriptionInfo" class="subscription-info">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="订阅地址">
              <el-link type="primary" :href="subscriptionInfo.url" target="_blank">
                {{ subscriptionInfo.url }}
              </el-link>
            </el-descriptions-item>
            <el-descriptions-item label="节点数量">
              {{ subscriptionInfo.nodeCount }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useStore } from '@/stores'
import type { Node } from '@/types'

// 状态管理
const store = useStore()

// 表单数据
const form = ref({
  maxDelay: 200,
  minSpeed: 1
})

// 加载状态
const loading = ref(false)
const filtering = ref(false)
const generating = ref(false)

// 节点数据
const filteredNodes = ref<Node[]>([])

// 分页
const currentPage = ref(1)
const pageSize = ref(20)

// 订阅信息
const subscriptionInfo = ref<{
  url: string
  nodeCount: number
} | null>(null)

// 计算属性
const totalNodes = computed(() => filteredNodes.value.length)
const availableNodes = computed(() => filteredNodes.value.filter(node => node.delay && node.delay < 1000).length)
const timeoutNodes = computed(() => filteredNodes.value.filter(node => !node.delay || node.delay >= 1000).length)

// 当前页数据
const currentPageNodes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredNodes.value.slice(start, end)
})

// 处理筛选
const handleFilter = async () => {
  try {
    filtering.value = true
    const { maxDelay, minSpeed } = form.value
    const response = await store.filterNodes(maxDelay, minSpeed * 1024) // 转换为KB/s
    filteredNodes.value = response.nodes
    currentPage.value = 1 // 重置页码
    ElMessage.success('节点筛选完成')
  } catch (error: any) {
    ElMessage.error(error.message || '节点筛选失败')
  } finally {
    filtering.value = false
  }
}

// 处理生成订阅
const handleGenerateSubscription = async () => {
  try {
    generating.value = true
    const response = await store.generateSubscription()
    subscriptionInfo.value = {
      url: response.url,
      nodeCount: response.nodeCount
    }
    ElMessage.success('订阅生成成功')
  } catch (error: any) {
    ElMessage.error(error.message || '订阅生成失败')
  } finally {
    generating.value = false
  }
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

// 格式化函数
const formatDelay = (delay?: number) => {
  if (!delay) return '超时'
  return `${delay}ms`
}

const formatSpeed = (speed?: number) => {
  if (!speed) return '未测试'
  if (speed < 1024) return `${speed.toFixed(1)}KB/s`
  return `${(speed / 1024).toFixed(1)}MB/s`
}

// 样式函数
const getNodeTypeTag = (type: string) => {
  const typeMap: Record<string, string> = {
    vmess: '',
    ss: 'success',
    hysteria2: 'warning',
    trojan: 'danger'
  }
  return typeMap[type.toLowerCase()] || 'info'
}

const getDelayTag = (delay?: number) => {
  if (!delay) return 'danger'
  if (delay < 100) return 'success'
  if (delay < 300) return ''
  if (delay < 500) return 'warning'
  return 'danger'
}

const getSpeedTag = (speed?: number) => {
  if (!speed) return 'info'
  if (speed > 10 * 1024) return 'success' // 10MB/s
  if (speed > 5 * 1024) return '' // 5MB/s
  if (speed > 1024) return 'warning' // 1MB/s
  return 'danger'
}
</script>

<style scoped lang="scss">
@use '@/styles/common.scss';

.node-select {
  .list-stats {
    display: flex;
    gap: 10px;
    margin-left: 15px;
    flex-wrap: wrap;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 20px;
    flex-wrap: wrap;

    @media screen and (max-width: 576px) {
      justify-content: center;
    }

    .total-text {
      color: #606266;
      font-size: 14px;
    }

    .page-size-select {
      width: 120px;
    }
  }

  .subscription-info {
    margin-top: 20px;
  }

  :deep(.el-tag) {
    text-transform: none;
  }
}
</style> 