package config

import "github.com/rancher-setup/internal/constants/container"

type DriverInterface interface {
	container.ContainerInterface
	Listen()
}
