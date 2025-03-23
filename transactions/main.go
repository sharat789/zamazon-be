package main

import (
	"github.com/sharat789/zamazon-be-ms/transactions/configs"
	api "github.com/sharat789/zamazon-be-ms/transactions/internal/api"
	"log"
)

func main() {

	cfg, err := configs.EnvSetup()

	if err != nil {
		log.Fatalf("config file is not loaded properly %v\n", err)
	}
	api.StartServer(cfg)
}
