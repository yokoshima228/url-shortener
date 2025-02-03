package main

import (
	"fmt"

	"github.com/yokoshima228/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.Env)
	fmt.Println(cfg.StoragePath)
	fmt.Println(cfg.HttpServer.Address)
}
