// 系统状态
export interface SystemStatus {
  nodeStats: {
    total: number
    current: number
    slow: number
    failed: number
  }
  currentStatus: {
    finalSubscription: string
    mergedSubscription: string
  }
  history: Array<{
    time: string
    url: string
  }>
}

// 订阅
export interface Subscription {
  id: string
  name: string
  type: string
  url: string
  nodeCount: number
  importTime: string
}

// 节点
export interface Node {
  id: string
  type: string
  alias: string
  address: string
  port: number
  protocol: string
  group: string
  delay?: number
  speed?: number
}

// 系统设置
export interface SystemSettings {
  scheduledTasks: {
    enabled: boolean
    subscriptionUpdate: {
      type: 'daily' | 'weekly'
      weekday?: number
      time?: {
        hour: number
        minute: number
      }
    }
    nodeCheck: {
      interval: number
    }
  }
}

// 日志
export interface Log {
  time: string
  operation: string
}

// API响应
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
} 