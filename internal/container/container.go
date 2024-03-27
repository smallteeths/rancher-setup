package container

import (
	"github.com/rancher-setup/internal/constants/config"
	"strings"
	"sync"
)

var containerMap sync.Map

type Container struct {
}

var _ config.CacheInterface = (*Container)(nil)

func CreateContainerFactory() *Container {
	return &Container{}
}

func (c *Container) Delete(key string) {
	containerMap.Delete(key)
}

func (c *Container) Get(key string) any {
	if value, exist := containerMap.Load(key); exist {
		return value
	}
	return nil
}

func (c *Container) Set(key string, value any) bool {
	containerMap.Store(key, value)
	return true
}

func (c *Container) Has(key string) bool {
	_, ok := containerMap.Load(key)
	return ok
}

func (c *Container) FuzzyDelete(keyPre string) {
	containerMap.Range(func(key, value interface{}) bool {
		if key, ok := key.(string); ok {
			if strings.HasPrefix(key, keyPre) {
				containerMap.Delete(key)
			}
		}
		return true
	})
}
