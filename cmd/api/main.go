package main

import (
	"github.com/SlavaShagalov/proxy-server/internal/pkg/config"
	"github.com/SlavaShagalov/proxy-server/internal/pkg/db"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	pLog "github.com/SlavaShagalov/proxy-server/internal/pkg/log/zap"

	repeatDel "github.com/SlavaShagalov/proxy-server/internal/repeat/delivery/http"
	requestsDel "github.com/SlavaShagalov/proxy-server/internal/requests/delivery/http"
	requestsRepository "github.com/SlavaShagalov/proxy-server/internal/requests/repository/mongo"
	requestsUsecase "github.com/SlavaShagalov/proxy-server/internal/requests/usecase"
	scanDel "github.com/SlavaShagalov/proxy-server/internal/scan/delivery/http"
)

const (
	apiAddress = ":8000"
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

	router := mux.NewRouter()

	requestsRep := requestsRepository.New(mongoDB.Collection("requests"), logger)
	requestsUC := requestsUsecase.New(requestsRep)
	requestsDel.RegisterHandlers(router, requestsUC)
	repeatDel.RegisterHandlers(router, requestsRep, logger)
	scanDel.RegisterHandlers(router, requestsRep, logger)

	server := http.Server{
		Addr:    apiAddress,
		Handler: router,
	}

	logger.Info("Starting API server", zap.String("address", apiAddress))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("API server stopped %v", zap.Error(err))
	}
}
