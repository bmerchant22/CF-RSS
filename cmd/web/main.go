package main

import (
	"encoding/json"
	"fmt"
	"github.com/bmerchant22/project/pkg/cfapi"
	"go.uber.org/zap"
	"log"
)

func main() {
	environment := "development"
	var logger *zap.Logger
	var loggerErr error

	if environment == "development" {
		if logger, loggerErr = zap.NewDevelopment(); loggerErr != nil {
			log.Fatalln(loggerErr)
		}
	} else {
		if logger, loggerErr = zap.NewProduction(); loggerErr != nil {
			log.Fatalln(loggerErr)
		}
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	obj := cfapi.NewCodeforcesClient()
	//obj.RecentActions(1)
	recentActions, err := obj.RecentActions(1)
	if err != nil {
		fmt.Println("error occured")
		return
	}
	//zap.S().Info(recentActions)
	data, err1 := json.MarshalIndent(recentActions, "", " ")
	if err1 != nil {
		fmt.Println("error occurred")
		return
	}
	zap.S().Info(string(data))
}
