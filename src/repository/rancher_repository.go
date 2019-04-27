package repository

import (
	"github.com/slack-bot-4all/slack-bot/src/config"
	"github.com/slack-bot-4all/slack-bot/src/model"
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

// FindRancherByName : consults the db with the name
func FindRancherByName(r *model.Rancher) (err error) {
	if err := config.DB.Where("name = ?", r.Name).First(r).Error; err != nil {
		return err
	}

	return nil
}
