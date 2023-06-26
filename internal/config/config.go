package config

import (
	"github.com/SideSwapIN/Analystic/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		DNS string `mapstructure:"dns"`
	} `mapstructure:"database"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

var cfg Config

func Init(cfgPath string) error {
	// 读取配置文件
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Failed to read config file: %v", err)
		return err
	}

	// 将配置文件解析到结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Errorf("Failed to unmarshal config: %v", err)
		return err
	}
	return nil
}

// GetConfig 返回全局配置结构体
func GetConfig() *Config {
	return &cfg
}
