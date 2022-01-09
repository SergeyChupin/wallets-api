package main

import (
	"flag"

	"github.com/SergeyChupin/wallets-api/internal/app/httpserver"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/httpserver.yaml", "path to config file")
}

func main() {
	flag.Parse()
	httpserver.Run(configPath)
}
