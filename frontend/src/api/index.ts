import axios from 'axios'

const baseURL = 'http://localhost:3355'

const api = axios.create({
  baseURL,
  timeout: 30000
})

// 状态监控
export const getSystemStatus = () => api.get('/api/status')

// 订阅管理
export const importSubscription = (data: { name: string; url: string }) => api.post('/api/subscription/import', data)
export const getSubscriptionList = () => api.get('/api/subscription/list')
export const mergeSubscriptions = (ids: string[]) => api.post('/api/subscription/merge', { ids })
export const deleteSubscriptions = (ids: string[]) => api.delete('/api/subscription/delete', { data: { ids } })

// 节点测试
export const importNodes = () => api.post('/api/node/import')
export const testNodes = () => api.post('/api/node/test')
export const getNodeList = () => api.get('/api/node/list')

// 优选节点
export const filterNodes = (data: { maxDelay: number; minSpeed: number }) => api.post('/api/node/filter', data)
export const generateSubscription = () => api.post('/api/node/generate-subscription')

// 系统设置
export const getSystemSettings = () => api.get('/api/settings')
export const updateSystemSettings = (data: any) => api.post('/api/settings', data)

// 运行日志
export const getLogs = (params: { page: number; pageSize: number }) => 
  api.get('/api/logs', { params })

export const clearLogs = () => 
  api.delete('/api/logs')

export default api 