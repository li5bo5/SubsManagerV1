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
    
    return nil
} 