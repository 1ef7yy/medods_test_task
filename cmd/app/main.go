package main

import (
	"net/http"
	"os"

	"github.com/1ef7yy/medods_test_task/internal/routes"
	"github.com/1ef7yy/medods_test_task/internal/view"
	"github.com/1ef7yy/medods_test_task/pkg/logger"
)

func main() {
	logger := logger.NewLogger()

	logger.Info("starting auth service...")

	view, err := view.NewView(logger)
	if err != nil {
		logger.Fatalf("could not initialize the view layer: %s", err.Error())
	}

	mux := routes.InitRouter(view)

	logger.Info("initialized router...")

	serverAddr, ok := os.LookupEnv("SERVER_ADDRESS")

	if !ok {
		serverAddr = "localhost:8080"
		logger.Warnf("could not resolve SERVER_ADDRESS from environment, reverting to default: %s", serverAddr)
	}

	logger.Infof("starting server on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		logger.Errorf("error starting server: %s", err.Error())
	}
}
