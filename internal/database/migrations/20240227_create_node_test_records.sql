-- 创建节点测试记录表
CREATE TABLE IF NOT EXISTS node_test_records (
    id VARCHAR(36) PRIMARY KEY,
    node_id VARCHAR(36) NOT NULL,
    latency INTEGER,           -- 延迟(ms)
    speed DECIMAL(10,2),      -- 下载速度(MB/s)
    test_time TIMESTAMP NOT NULL,
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_node_id (node_id),
    INDEX idx_test_time (test_time)
); 