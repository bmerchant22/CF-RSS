package worker

import (
	"fmt"
	"github.com/bmerchant22/project/pkg/cfapi"
	"github.com/bmerchant22/project/pkg/models"
	"github.com/bmerchant22/project/pkg/store"
	"go.uber.org/zap"
	"time"
)

func PerformWork() {

	for {
		mongoStore := store.MongoStore{}

		mongoStore.ConnectToDatabase()

		obj := cfapi.NewCodeforcesClient()
		RecentActions, err := obj.RecentActions(100)
		if err != nil {
			fmt.Println("error occurred")
			return
		}

		maxTimeStamp, err := mongoStore.GetMaxTimeStamp()
		if err != nil {
			zap.S().Errorf("Error while getting maxTimeStamp: %v", maxTimeStamp)
		}

		zap.S().Info("Got maxTimeStamp successfully")

		var NewData []models.RecentAction

		for i := 0; i < len(RecentActions); i++ {
			if RecentActions[i].TimeSeconds > maxTimeStamp {
				NewData = append(NewData, RecentActions[i])
			}
		}

		zap.S().Info("RecentActions stored in NewData successfully ")

		err = mongoStore.StoreRecentActionsInTheDatabase(NewData)
		if err != nil {
			zap.S().Errorf("Error occurred while storing data : %v", err)
			return
		}

		var temp []models.RecentAction

		NewData = temp

		zap.S().Info("The worker will sleep for 5 min now.")
		time.Sleep(5 * time.Minute)
	}
}
