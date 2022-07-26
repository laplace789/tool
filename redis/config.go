package rdb

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Host         string `json:"host"      mapstructure:"host"`
	Port         string `json:"port"      mapstructure:"port"`
	UserName     string `json:"user_name" mapstructure:"user_name"`
	Password     string `json:"password" mapstructure:"password"`
	Db           int    `json:"db" mapstructure:"db"`
	PoolSize     int    `json:"pool_size" mapstructure:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns" mapstructure:"min_idle_conns"`
	MaxRetries   int    `json:"max_retries" mapstructure:"max_retries"`
	Open         bool   `json:"open" mapstructure:"open"`
}

func initRedisConf() *Config {
	path := "./conf"
	name := "redis"
	configType := "yml"
	conf := &Config{}
	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType(configType)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	v.Unmarshal(conf)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		v.Unmarshal(conf)
		log.Println(conf)
	})

	return conf
}
