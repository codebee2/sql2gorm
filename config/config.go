package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	SQLFile     string `yaml:"sql_file"`     // SQL文件路径
	OutputFile  string `yaml:"output_file"`  // 输出Go文件路径
	TableName   string `yaml:"table_name"`   // 表名
	PackageName string `yaml:"package_name"` // 包名
	DSN         string `yaml:"dsn"`          // 数据库连接字符串
}

// LoadConfig 从文件加载配置
func LoadConfig(configFile string) (*Config, error) {
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{
			PackageName: "models",
		}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	if config.PackageName == "" {
		config.PackageName = "models"
	}

	return &config, nil
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig(configFile string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}
