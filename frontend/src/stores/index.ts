import { defineStore } from 'pinia'
import type { SystemStatus, Subscription, Node, SystemSettings, Log } from '@/types'
import * as api from '@/api'

export const useStore = defineStore('main', {
  state: () => ({
    systemStatus: null as SystemStatus | null,
    subscriptions: [] as Subscription[],
    nodes: [] as Node[],
    settings: null as SystemSettings | null,
    logs: [] as Log[],
    loading: false
  }),

  actions: {
    // 状态监控
    async fetchSystemStatus() {
      try {
        const { data } = await api.getSystemStatus()
        this.systemStatus = data
      } catch (error) {
        console.error('获取系统状态失败:', error)
      }
    },

    // 订阅管理
    async fetchSubscriptions() {
      try {
        const { data } = await api.getSubscriptionList()
        this.subscriptions = data
      } catch (error) {
        console.error('获取订阅列表失败:', error)
      }
    },

    async importSubscription(name: string, url: string) {
      try {
        await api.importSubscription({ name, url })
        await this.fetchSubscriptions()
      } catch (error) {
        console.error('导入订阅失败:', error)
        throw error
      }
    },

    async mergeSubscriptions(ids: string[]) {
      try {
        await api.mergeSubscriptions(ids)
        await this.fetchSystemStatus()
      } catch (error) {
        console.error('合并订阅失败:', error)
        throw error
      }
    },

    async deleteSubscriptions(ids: string[]) {
      try {
        await api.deleteSubscriptions(ids)
        await this.fetchSubscriptions()
      } catch (error) {
        console.error('删除订阅失败:', error)
        throw error
      }
    },

    // 节点测试
    async importNodes() {
      try {
        await api.importNodes()
        const { data } = await api.getNodeList()
        this.nodes = data
      } catch (error) {
        console.error('导入节点失败:', error)
        throw error
      }
    },

    async testNodes() {
      try {
        this.loading = true
        await api.testNodes()
        const { data } = await api.getNodeList()
        this.nodes = data
      } catch (error) {
        console.error('测试节点失败:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // 节点筛选
    async filterNodes(maxDelay: number, minSpeed: number) {
      try {
        const { data } = await api.filterNodes({ maxDelay, minSpeed })
        return {
          nodes: data.nodes,
          totalNodes: data.total_nodes
        }
      } catch (error) {
        console.error('筛选节点失败:', error)
        throw error
      }
    },

    // 生成优选订阅
    async generateSubscription() {
      try {
        const { data } = await api.generateSubscription()
        return {
          url: data.sub_url,
          nodeCount: data.node_count
        }
      } catch (error) {
        console.error('生成订阅失败:', error)
        throw error
      }
    },

    // 系统设置
    async fetchSettings() {
      try {
        const { data } = await api.getSystemSettings()
        this.settings = data
      } catch (error) {
        console.error('获取系统设置失败:', error)
      }
    },

    async updateSettings(settings: SystemSettings) {
      try {
        await api.updateSystemSettings(settings)
        this.settings = settings
      } catch (error) {
        console.error('更新系统设置失败:', error)
        throw error
      }
    },

    // 运行日志
    async fetchLogs(params: { page: number; pageSize: number }) {
      try {
        const { data } = await api.getLogs(params)
        return {
          logs: data.logs,
          total: data.total
        }
      } catch (error) {
        console.error('获取运行日志失败:', error)
        throw error
      }
    },

    async clearLogs() {
      try {
        await api.clearLogs()
      } catch (error) {
        console.error('清空日志失败:', error)
        throw error
      }
    }
  }
}) 