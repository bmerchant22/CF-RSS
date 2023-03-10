package main

import (
	"github.com/bmerchant22/project/pkg/store"
	"github.com/bmerchant22/project/pkg/web"
	"github.com/bmerchant22/project/pkg/worker"
	"go.uber.org/zap"
	"log"
	"sync"
)

func main() {
	var logger *zap.Logger
	var loggerErr error

	if logger, loggerErr = zap.NewDevelopment(); loggerErr != nil {
		log.Fatalln(loggerErr)
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	mongoStore := store.MongoStore{}

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		worker.PerformWork(&mongoStore)
	}(wg)

	func(wg *sync.WaitGroup) {
		defer wg.Done()
		srv := web.CreateWebServer(&mongoStore)
		if err := srv.StartListeningForRequests(":8080"); err != nil {
			zap.S().Errorf("Error while starting server from main : %v", err)
		}
	}(wg)

	wg.Wait()

	return
}
