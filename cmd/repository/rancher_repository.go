package repository

import (
	"github.com/slack-bot-4all/slack-bot/cmd/config"
	"github.com/slack-bot-4all/slack-bot/cmd/model"
)

// AddRancher : add a Rancher to database
func AddRancher(r *model.Rancher) error {
	if err := config.DB.Create(r).Error; err != nil {
		return err
	}

	return nil
}

// ListRancher :
func ListRancher(r *[]model.Rancher) (err error) {
	if err = config.DB.Find(r).Error; err != nil {
		return err
	}

	return nil
}
