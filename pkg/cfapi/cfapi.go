package cfapi

import (
	"encoding/json"
	"fmt"
	"github.com/bmerchant22/project/pkg/models"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type CodeforcesAPI interface {
	RecentActions(maxCount int) ([]models.RecentAction, error)
}

type CodeforcesClient struct {
	client http.Client
}

func (cfClient *CodeforcesClient) RecentActions(maxCount int) ([]models.RecentAction, error) {

	resp, err := cfClient.client.Get("https://codeforces.com/api/recentActions?maxCount=100")
	if err != nil {
		zap.S().Errorf("Error occured while calling cf api: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error occured while reading the resp body")
		return nil, err
	}

	//zap.S().Info(string(data))

	wrapper := struct {
		Status string
		Result []models.RecentAction
	}{}

	json.Unmarshal(data, &wrapper)

	return wrapper.Result, err
}

func NewCodeforcesClient() CodeforcesAPI {
	obj := new(CodeforcesClient)
	return obj
}
