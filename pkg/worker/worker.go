package worker

import (
	"fmt"
	"github.com/bmerchant22/project/pkg/cfapi"
	"github.com/bmerchant22/project/pkg/models"
	"github.com/bmerchant22/project/pkg/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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

		e := echo.New()
		e.GET("/recent-actions", func(c echo.Context) error {

			after, err := strconv.ParseInt(c.QueryParam("after"), 10, 64)
			if err != nil {
				zap.S().Errorf("Error while converting after string to int")
				c.String(http.StatusBadRequest, "Enter valid query params.")
			}

			zap.S().Infof("After converted to int successfully %v", after)
			recentActions, err := mongoStore.QueryRecentActions(after)
			if err != nil {
				zap.S().Errorf("Error occurred while calling QueryRecentActions: %v", err)
				return nil
			}
			return c.JSON(http.StatusOK, recentActions)
		})
		e.Logger.Fatal(e.Start(":4000"))

		//zap.S().Info("The worker will sleep for 5 min now.")
		//time.Sleep(5 * time.Minute)
	}
}
