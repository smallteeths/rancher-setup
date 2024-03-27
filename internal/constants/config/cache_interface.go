package config

import (
	"github.com/rancher-setup/internal/constants/container"
)

// CacheInterface config cache
type CacheInterface interface {
	container.ContainerInterface
	FuzzyDelete(key string)
}
