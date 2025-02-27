<!-- 订阅管理页面 -->
<template>
  <div class="subscriptions page-container">
    <!-- 订阅列表 -->
    <el-card class="card-section">
      <template #header>
        <div class="section-header">
          <div class="header-title">
            <span>订阅列表</span>
            <el-tag 
              v-if="subscriptions.length"
              type="info" 
              class="count-tag"
            >
              共 {{ subscriptions.length }} 个
            </el-tag>
          </div>
          <div class="header-actions">
            <el-button 
              type="primary" 
              :icon="Plus"
              @click="showAddDialog = true"
            >
              添加订阅
            </el-button>
          </div>
        </div>
      </template>

      <div 
        v-loading="loading"
        class="content-section responsive-table"
        :class="{ empty: !subscriptions.length }"
      >
        <el-table 
          v-if="subscriptions.length"
          :data="subscriptions"
          style="width: 100%"
        >
          <el-table-column 
            prop="name" 
            label="名称"
            min-width="120"
          />
          <el-table-column 
            prop="url" 
            label="地址"
            min-width="200"
            show-overflow-tooltip
          />
          <el-table-column 
            prop="nodeCount" 
            label="节点数量"
            width="100"
            align="center"
          />
          <el-table-column 
            prop="lastUpdate" 
            label="最后更新"
            width="180"
            align="center"
          >
            <template #default="{ row }">
              {{ formatTime(row.lastUpdate) }}
            </template>
          </el-table-column>
          <el-table-column 
            label="操作"
            width="200"
            fixed="right"
          >
            <template #default="{ row }">
              <el-button-group>
                <el-button 
                  :loading="row.updating"
                  @click="handleUpdate(row)"
                >
                  更新
                </el-button>
                <el-button 
                  type="danger"
                  @click="handleDelete(row)"
                >
                  删除
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>

        <el-empty 
          v-else 
          description="暂无订阅"
        />
      </div>
    </el-card>

    <!-- 添加订阅对话框 -->
    <el-dialog
      v-model="showAddDialog"
      title="添加订阅"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item 
          label="名称" 
          prop="name"
        >
          <el-input 
            v-model="form.name"
            placeholder="请输入订阅名称"
          />
        </el-form-item>
        <el-form-item 
          label="地址" 
          prop="url"
        >
          <el-input 
            v-model="form.url"
            placeholder="请输入订阅地址"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button 
          type="primary" 
          :loading="adding"
          @click="handleAdd"
        >
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
@use '@/styles/common.scss';

.subscriptions {
  :deep(.el-dialog) {
    @media screen and (max-width: 576px) {
      width: 90% !important;
      margin: 0 auto;
    }
  }

  :deep(.el-button-group) {
    display: flex;
    gap: 8px;
  }
}
</style> 