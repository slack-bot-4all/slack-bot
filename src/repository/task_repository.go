package repository

import (
	"github.com/slack-bot-4all/slack-bot/src/config"
	"github.com/slack-bot-4all/slack-bot/src/model"
)

// AddTask : add a Task to database
func AddTask(t *model.Task) (err error) {
	if err := config.DB.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// ListTask :
func ListTask(t *[]model.Task) (err error) {
	if err = config.DB.Find(t).Error; err != nil {
		return err
	}

	return nil
}
