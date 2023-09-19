package main

import (
	"github.com/SlavaShagalov/proxy-server/internal/pkg/config"
	pLog "github.com/SlavaShagalov/proxy-server/internal/pkg/log/zap"
	"github.com/SlavaShagalov/proxy-server/internal/proxy"
	requestsRepository "github.com/SlavaShagalov/proxy-server/internal/requests/repository/mongo"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/SlavaShagalov/proxy-server/internal/pkg/db"
)

const (
	proxyAddress = ":8080"
)

func main() {
	// Logger
	logger := pLog.NewProdLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	config.SetDefaultMongoConfig()

	// Storage
	mongoDB, err := db.NewMongoDB(logger)
	if err != nil {
		os.Exit(1)
	}

	requestsRep := requestsRepository.New(mongoDB.Collection("requests"), logger)
	handler := proxy.New(requestsRep, logger)

	logger.Info("Starting proxy server", zap.String("proxy_addr", proxyAddress))
	err = http.ListenAndServe(proxyAddress, handler)
	if err != nil {
		log.Fatalf("Proxy server stopped %s", err.Error())
	}
}
