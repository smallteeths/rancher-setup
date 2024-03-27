package main

import (
	_ "github.com/rancher-setup/internal/bootstrap"
	"github.com/rancher-setup/internal/server"
	"github.com/rancher-setup/internal/variable"
	"github.com/rancher-setup/router"
)

func main() {
	port := variable.Config.GetString("HttpServer.Port")
	mode := variable.Config.GetString("HttpServer.Mode")

	http := server.New(
		server.WithPort(port),
		server.WithMode(mode),
		server.WithLogger(variable.Log),
	)
	http.SetRouters(router.New(http)).Run()
}
