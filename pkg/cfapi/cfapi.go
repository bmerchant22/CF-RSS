package cfapi

import "github.com/bmerchant2253/project/pkg/models"

type CodeforcesAPI interface {
	RecentActions(maxCount int) ([]models.RecentAction, error)
}
