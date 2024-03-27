package bootstrap

import (
	"log"

	"github.com/rancher-setup/internal/config"
	"github.com/rancher-setup/internal/config/driver"
	"github.com/rancher-setup/internal/logger"
	"github.com/rancher-setup/internal/variable"
	"github.com/rancher-setup/internal/variable/consts"
)

func init() {
	var err error
	if variable.Config, err = config.New(driver.New(), config.Options{
		BasePath: variable.BasePath,
	}); err != nil {
		log.Fatal("aaaaaaaaaaaa")
		log.Fatal(consts.ErrorInitConfig)
		log.Fatal(err.Error())
	}
	if variable.Log, err = logger.New(
		logger.WithDebug(true),
		logger.WithEncode("json"),
		logger.WithFilename(variable.BasePath+"/storage/logs/system.log"),
	); err != nil {
		log.Fatal(consts.ErrorInitLogger)
	}
}
