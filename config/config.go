package config

type Config struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    
    Storage struct {
        Path string `yaml:"path"`
    } `yaml:"storage"`
    
    Subscription struct {
        UpdateInterval string `yaml:"update_interval"`
        TestInterval  string `yaml:"test_interval"`
        MaxConcurrent int    `yaml:"max_concurrent"`
    } `yaml:"subscription"`

    Filter struct {
        MaxLatency int     `yaml:"max_latency"` // 最大延迟(ms)
        MinSpeed   float64 `yaml:"min_speed"`   // 最小速度(MB/s)
    } `yaml:"filter"`
}

var GlobalConfig Config

func Init() error {
    // 设置默认配置
    GlobalConfig.Server.Port = 3355
    GlobalConfig.Server.Host = "localhost"
    GlobalConfig.Storage.Path = "./data"
    GlobalConfig.Subscription.MaxConcurrent = 5
    GlobalConfig.Subscription.UpdateInterval = "24h"
    GlobalConfig.Subscription.TestInterval = "4h"
    
    // 设置默认筛选配置
    GlobalConfig.Filter.MaxLatency = 400  // 默认最大延迟 400ms
    GlobalConfig.Filter.MinSpeed = 1.0    // 默认最小速度 1MB/s
    
    return nil
} 