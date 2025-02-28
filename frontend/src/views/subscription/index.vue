<template>
  <div class="subscription">
    <!-- 订阅导入区域 -->
    <el-card class="subscription-section">
      <template #header>
        <div class="section-header">
          <span>订阅导入</span>
        </div>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入订阅名称" :disabled="loading" />
        </el-form-item>
        <el-form-item label="链接" prop="url">
          <el-input v-model="form.url" placeholder="请输入订阅链接" :disabled="loading" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleImport">导入订阅</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 订阅列表区域 -->
    <el-card class="subscription-section">
      <template #header>
        <div class="section-header">
          <span>订阅列表</span>
          <div class="header-actions">
            <el-button 
              type="primary" 
              :disabled="!selectedSubscriptions.length || loading" 
              @click="handleMerge"
            >
              合并选中订阅
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="subscriptions"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="订阅名称" min-width="120" />
        <el-table-column prop="url" label="订阅链接" min-width="200" show-overflow-tooltip />
        <el-table-column prop="nodeCount" label="节点数量" width="100" align="center" />
        <el-table-column prop="updateTime" label="更新时间" width="180" align="center">
          <template #default="{ row }">
            {{ formatDate(row.updateTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" align="center">
          <template #default="{ row }">
            <el-button 
              type="danger" 
              link
              :disabled="loading"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useStore } from '@/stores'
import type { FormInstance } from 'element-plus'
import type { Subscription } from '@/types'

// 状态管理
const store = useStore()

// 表单相关
const formRef = ref<FormInstance>()
const form = ref({
  name: '',
  url: ''
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入订阅名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入订阅链接', trigger: 'blur' },
    { type: 'url', message: '请输入正确的URL格式', trigger: 'blur' }
  ]
}

// 加载状态
const loading = ref(false)

// 订阅列表
const subscriptions = ref<Subscription[]>([])
const selectedSubscriptions = ref<Subscription[]>([])

// 选择变化处理
const handleSelectionChange = (selection: Subscription[]) => {
  selectedSubscriptions.value = selection
}

// 导入订阅
const handleImport = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    await store.importSubscription(form.value.name, form.value.url)
    ElMessage.success('订阅导入成功')
    
    // 重置表单
    formRef.value.resetFields()
    // 刷新列表
    await fetchSubscriptions()
  } catch (error: any) {
    ElMessage.error(error.message || '订阅导入失败')
  } finally {
    loading.value = false
  }
}

// 删除订阅
const handleDelete = async (subscription: Subscription) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除订阅 "${subscription.name}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    loading.value = true
    await store.deleteSubscriptions([subscription.id])
    ElMessage.success('删除成功')
    await fetchSubscriptions()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  } finally {
    loading.value = false
  }
}

// 合并订阅
const handleMerge = async () => {
  if (!selectedSubscriptions.value.length) {
    ElMessage.warning('请选择要合并的订阅')
    return
  }

  try {
    loading.value = true
    const ids = selectedSubscriptions.value.map(sub => sub.id)
    await store.mergeSubscriptions(ids)
    ElMessage.success('订阅合并成功')
  } catch (error: any) {
    ElMessage.error(error.message || '订阅合并失败')
  } finally {
    loading.value = false
  }
}

// 获取订阅列表
const fetchSubscriptions = async () => {
  try {
    loading.value = true
    await store.fetchSubscriptions()
    subscriptions.value = store.subscriptions
  } catch (error: any) {
    ElMessage.error(error.message || '获取订阅列表失败')
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (timestamp: number) => {
  return new Date(timestamp).toLocaleString()
}

// 页面加载时获取订阅列表
onMounted(() => {
  fetchSubscriptions()
})
</script>

<style scoped lang="scss">
.subscription {
  .subscription-section {
    margin-bottom: 20px;

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: 500;

      .header-actions {
        display: flex;
        gap: 10px;
      }
    }
  }
}
</style> 