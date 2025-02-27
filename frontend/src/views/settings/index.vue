<template>
  <div class="settings page-container">
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>系统设置</span>
          </div>
          <div class="header-actions">
            <el-button 
              type="primary" 
              :loading="saving"
              @click="handleSave"
            >
              保存设置
            </el-button>
          </div>
        </div>
      </template>

      <div class="content-section">
        <el-form 
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="180px"
        >
          <!-- 基本设置 -->
          <div class="setting-group">
            <div class="group-title">基本设置</div>
            <el-form-item 
              label="自动更新间隔(分钟)" 
              prop="updateInterval"
            >
              <el-input-number 
                v-model="form.updateInterval"
                :min="1"
                :max="1440"
                :step="1"
                class="form-item half-width"
              />
            </el-form-item>

            <el-form-item 
              label="延迟测试超时(秒)" 
              prop="testTimeout"
            >
              <el-input-number 
                v-model="form.testTimeout"
                :min="1"
                :max="60"
                :step="1"
                class="form-item half-width"
              />
            </el-form-item>

            <el-form-item 
              label="并发测试数量" 
              prop="concurrentTests"
            >
              <el-input-number 
                v-model="form.concurrentTests"
                :min="1"
                :max="50"
                :step="1"
                class="form-item half-width"
              />
            </el-form-item>
          </div>

          <!-- 代理设置 -->
          <div class="setting-group">
            <div class="group-title">代理设置</div>
            <el-form-item 
              label="本地监听地址" 
              prop="listenAddress"
            >
              <el-input 
                v-model="form.listenAddress"
                placeholder="127.0.0.1"
                class="form-item half-width"
              />
            </el-form-item>

            <el-form-item 
              label="本地监听端口" 
              prop="listenPort"
            >
              <el-input-number 
                v-model="form.listenPort"
                :min="1"
                :max="65535"
                :step="1"
                class="form-item half-width"
              />
            </el-form-item>

            <el-form-item 
              label="启用系统代理" 
              prop="enableSystemProxy"
            >
              <el-switch 
                v-model="form.enableSystemProxy"
              />
            </el-form-item>
          </div>

          <!-- 订阅设置 -->
          <div class="setting-group">
            <div class="group-title">订阅设置</div>
            <el-form-item 
              label="订阅更新间隔(小时)" 
              prop="subscriptionUpdateInterval"
            >
              <el-input-number 
                v-model="form.subscriptionUpdateInterval"
                :min="1"
                :max="168"
                :step="1"
                class="form-item half-width"
              />
            </el-form-item>

            <el-form-item 
              label="自动更新订阅" 
              prop="autoUpdateSubscription"
            >
              <el-switch 
                v-model="form.autoUpdateSubscription"
              />
            </el-form-item>
          </div>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useStore } from '@/stores'
import type { FormInstance } from 'element-plus'

// 状态管理
const store = useStore()

// 表单引用
const formRef = ref<FormInstance>()

// 加载状态
const saving = ref(false)

// 表单数据
const form = ref({
  enabled: false,
  subscriptionUpdate: {
    type: 'daily',
    weekday: '1',
    time: new Date(2000, 0, 1, 3, 0) // 默认凌晨3点
  },
  nodeCheck: {
    interval: '4'
  }
})

// 表单校验规则
const rules = {
  subscriptionUpdate: [
    { required: true, message: '请配置订阅更新时间', trigger: 'change' }
  ],
  nodeCheck: [
    { required: true, message: '请选择节点检测间隔', trigger: 'change' }
  ]
}

// 保存设置
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    
    saving.value = true
    const settings = {
      scheduledTasks: {
        enabled: form.value.enabled,
        subscriptionUpdate: {
          type: form.value.subscriptionUpdate.type,
          weekday: form.value.subscriptionUpdate.type === 'weekly' ? Number(form.value.subscriptionUpdate.weekday) : undefined,
          time: form.value.subscriptionUpdate.time ? {
            hour: form.value.subscriptionUpdate.time.getHours(),
            minute: form.value.subscriptionUpdate.time.getMinutes()
          } : undefined
        },
        nodeCheck: {
          interval: Number(form.value.nodeCheck.interval)
        }
      }
    }

    await store.updateSettings(settings)
    ElMessage.success('设置保存成功')
  } catch (error: any) {
    if (error?.errors) {
      ElMessage.error('请完善必填项')
    } else {
      ElMessage.error(error.message || '保存设置失败')
    }
  } finally {
    saving.value = false
  }
}

// 加载设置
const loadSettings = () => {
  const settings = store.settings
  if (!settings?.scheduledTasks) return

  const { enabled, subscriptionUpdate, nodeCheck } = settings.scheduledTasks
  
  form.value = {
    enabled,
    subscriptionUpdate: {
      type: subscriptionUpdate.type,
      weekday: String(subscriptionUpdate.weekday || '1'),
      time: subscriptionUpdate.time ? new Date(
        2000, 
        0, 
        1, 
        subscriptionUpdate.time.hour, 
        subscriptionUpdate.time.minute
      ) : new Date(2000, 0, 1, 3, 0)
    },
    nodeCheck: {
      interval: String(nodeCheck.interval)
    }
  }
}

// 页面加载时获取设置
onMounted(async () => {
  await store.fetchSettings()
  loadSettings()
})
</script>

<style scoped lang="scss">
@use '@/styles/common.scss';

.settings {
  .setting-group {
    margin-bottom: 30px;

    &:last-child {
      margin-bottom: 0;
    }

    .group-title {
      font-size: 16px;
      font-weight: 500;
      color: #303133;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 1px solid #EBEEF5;
    }
  }

  :deep(.el-form-item) {
    margin-bottom: 22px;

    @media screen and (max-width: 768px) {
      margin-bottom: 18px;
    }

    .el-form-item__label {
      font-weight: normal;
    }
  }

  :deep(.el-input-number) {
    width: 180px;

    @media screen and (max-width: 576px) {
      width: 100%;
    }
  }
}
</style> 