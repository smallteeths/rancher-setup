package config

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/rancher-setup/internal/constants/config"
	"github.com/rancher-setup/internal/container"
	"github.com/spf13/cast"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Driver config.DriverInterface
	cache  config.CacheInterface
	mu     *sync.Mutex
}

var _ config.ConfigInterface = (*Config)(nil)

type Options struct {
	Filename    string
	BasePath    string
	Cate        string
	CachePrefix string
	Cache       config.CacheInterface
}

// New config
func New(config config.DriverInterface, option Options) (provider *Config, err error) {
	if option.Cache == nil {
		option.Cache = container.CreateContainerFactory()
	}
	if option.CachePrefix == "" {
		option.CachePrefix = "config"
	}
	if option.Cate == "" {
		option.Cate = "yaml"
	}
	if d, ok := config.(interface{ Apply(Options) error }); ok {
		if err = d.Apply(option); err != nil {
			return
		}
	}
	if d, ok := config.(interface{ Listen() }); ok {
		d.Listen()
	}
	return &Config{
		Driver: config,
		cache:  option.Cache,
		mu:     new(sync.Mutex),
	}, nil
}

func (c *Config) Cache(key string, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Set(key, value)
}

func (c *Config) Get(key string) any {
	if c.cache.Has(key) {
		return c.cache.Get(key)
	}
	val := c.Driver.Get(key)
	c.Cache(key, val)
	return val
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

func (c *Config) GetInt32(key string) int32 {
	return cast.ToInt32(c.Get(key))
}

func (c *Config) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

func (c *Config) GetClientset() (*kubernetes.Clientset, error) {
	var kubeconfig string
	kubeconfig = filepath.Join("/root", ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
