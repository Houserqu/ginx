// 配置 viper
// 支持从环境变量 CONFIG_PATH 指定配置文件路径，默认为 ./config.yaml
// 自动监听文件变化并重新加载配置

package core

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		log.Println("Load config file from: ", configPath)
		viper.SetConfigFile(configPath)
	} else {
		log.Println("Load config file from: ./config.yaml")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	// 设置默认值
	viper.SetDefault("server.addr", "0.0.0.0:8080")

	// 加载配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}
