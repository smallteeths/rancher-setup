package driver

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rancher-setup/internal/config"
	constants "github.com/rancher-setup/internal/constants/config"
	"github.com/spf13/viper"
)

type ViperConfig struct {
	viper  *viper.Viper
	option config.Options
}

var _ constants.DriverInterface = (*ViperConfig)(nil)

var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

func New() *ViperConfig {
	return &ViperConfig{}
}

// Apply 创建实例
func (v *ViperConfig) Apply(option config.Options) error {
	v.option = option
	viperConfig := viper.New()
	viperConfig.AddConfigPath(option.BasePath + "/config")
	if len(option.Filename) == 0 {
		viperConfig.SetConfigName("config")
	} else {
		viperConfig.SetConfigName(option.Filename)
	}
	viperConfig.SetConfigType(option.Cate)
	if err := viperConfig.ReadInConfig(); err != nil {
		return err
	}
	v.viper = viperConfig
	return nil
}

// Listen 监听文件变化
func (v *ViperConfig) Listen() {
	v.viper.OnConfigChange(func(in fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if in.Op.String() == "WRITE" {
				// 清除cache内的配置
				v.option.Cache.FuzzyDelete(v.option.CachePrefix)
				lastChangeTime = time.Now()
			}
		}
	})
}

func (v *ViperConfig) Get(key string) any {
	return v.viper.Get(key)
}

func (v *ViperConfig) Set(key string, value any) bool {
	v.viper.Set(key, value)
	return true
}

func (v *ViperConfig) Has(key string) bool {
	return v.viper.IsSet(key)
}
