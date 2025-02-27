# SubsManager - 多订阅管理工具

SubsManager 是一个轻量级的多订阅管理工具，通过Docker容器化的方式提供服务，支持节点优选功能。

## 项目结构

- `internal/`：包含项目的主要逻辑代码
  - `models/`：定义数据模型
  - `services/`：包含业务逻辑代码
  - `handlers/`：包含HTTP处理程序

## 参考

- [v2rayN](https://github.com/2dust/v2rayN)
- [Sub-Store](https://github.com/sub-store-org/Sub-Store)
- [subs-check](https://github.com/beck-8/subs-check)
- [subs-check](https://github.com/bestruirui/subs-check)
- [OpenClash](https://github.com/vernesong/OpenClash)

## 功能特点

- 订阅管理：导入、解析、编辑、删除和整合订阅
- 节点测试：对节点进行延迟和下载速度测试
- 优选节点：根据条件筛选节点并生成订阅
- 状态监控：实时监控节点状态和历史记录
- 系统设置：配置定时任务和节点检测周期
- 运行日志：记录操作日志和异常信息

## 快速开始

### 使用Docker运行

```bash
# 拉取镜像
docker pull li5bo5/subsmanager:latest

# 运行容器
docker run -d \
  --name subsmanager \
  -p 3355:3355 \
  -v /path/to/data:/app/data \
  li5bo5/subsmanager:latest
```

## 配置说明

### 以下配置是Cursor说的，我也不是很懂，先这样用着

### 默认配置
```yaml
server:
  port: 3355        # 服务端口
  host: "0.0.0.0"   # 监听地址

storage:
  path: "/app/data" # 数据存储目录

speedtest:
  max_concurrent: 5 # 最大并发测试数
  timeout: 10      # 测试超时时间(秒)
  test_url: "http://cachefly.cachefly.net/100mb.test" # 测速URL

schedule:
  subscription_update: "0 0 */24 * * *" # 订阅更新间隔（每24小时）
  node_test: "0 0 */4 * * *"           # 节点检测间隔（每4小时）

filter:
  max_latency: 400    # 延迟阈值(ms)
  min_speed: 2.0      # 最低速度(MB/s)

log:
  max_entries: 1000   # 内存中保留的最大日志条数
  level: "INFO"       # 日志级别(INFO/ERROR/WARNING)
```

### 修改配置

有两种方式可以修改默认配置：

1. 使用配置文件：
```bash
# 1. 创建配置文件
cat > config.yaml << EOF
server:
  port: 3366        # 修改端口为3366
  
speedtest:
  max_concurrent: 10 # 修改并发数为10
  
schedule:
  subscription_update: "0 0 */12 * * *" # 修改为每12小时更新一次
EOF

# 2. 挂载配置文件运行容器
docker run -d \
  --name subsmanager \
  -p 3366:3366 \
  -v /path/to/config.yaml:/app/config.yaml \
  -v /path/to/data:/app/data \
  li5bo5/subsmanager:latest
```

2. 使用环境变量：
```bash
docker run -d \
  --name subsmanager \
  -p 3366:3366 \
  -v /path/to/data:/app/data \
  -e SUBS_SERVER_PORT=3366 \
  -e SUBS_SPEEDTEST_MAX_CONCURRENT=10 \
  -e SUBS_SCHEDULE_SUBSCRIPTION_UPDATE="0 0 */12 * * *" \
  li5bo5/subsmanager:latest
```

环境变量命名规则：
- 使用 `SUBS_` 前缀
- 配置项用下划线连接
- 全部大写

例如：
- `server.port` → `SUBS_SERVER_PORT`
- `speedtest.max_concurrent` → `SUBS_SPEEDTEST_MAX_CONCURRENT`
- `schedule.subscription_update` → `SUBS_SCHEDULE_SUBSCRIPTION_UPDATE`

配置优先级：
1. 环境变量（最高）
2. 配置文件
3. 默认配置（最低）

## 许可证

MIT License 