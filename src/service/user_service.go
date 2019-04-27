package service

import (
	"github.com/slack-bot-4all/slack-bot/src/model"
	"github.com/slack-bot-4all/slack-bot/src/repository"
)

// AddUser : have a business rules to add a User to db
func AddUser(u *model.User) error {
	var err error

	if u.Username != "" && u.Password != "" {
		err = repository.AddUser(u)
	}

	if err != nil {
		return err
	}

	return nil
}
