package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"singo/logger"
	"sync"
)

type UbqConfig struct {
	AppId               string `mapstructure:"appId"`
	Secret              string `mapstructure:"secret"`
	PermitMetaProductNo string `mapstructure:"permitMetaProductNo"`
}

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Redis    *RedisConfig
	Ubq      *UbqConfig
	Game     *GameConfig
}

type GameConfig struct {
	EnergyTime uint
	MaxEnergy  int
}
type ServerConfig struct {
	Port   int    `mapstructure:"port"`
	Secret string `mapstructure:"secret"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db"`
}

// 定义配置结构体
var config *Config
var configOnce sync.Once

func GetConfig() *Config {
	configOnce.Do(func() {
		config = loadConfig()
	})
	return config
}

// Init 初始化配置项
func loadConfig() *Config {
	// 设置 Viper
	viper.SetConfigFile("./config.yaml") // 替换为你的配置文件路径
	viper.SetConfigType("yaml")          // 如果是 YAML 格式的配置文件，这里设置为 "yaml"

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		logger.Panic("读取配置文件出错 %s\n", err)
		return nil
	}

	var cfg Config
	// 解析配置文件到结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Panic("无法解码 %s\n", err)
		return nil
	}

	// 输出初始配置
	logger.Debug("配置文件初始化成功")

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Debug("配置文件变更:", e.Name)

		// 重新解析配置文件到结构体
		if err := viper.Unmarshal(&cfg); err != nil {
			logger.Panic("Unable to decode into struct: %s\n", err)
			return
		}
		// 输出更新后的配置
		logger.Debug("配置更新成功")
	})

	return &cfg
}
