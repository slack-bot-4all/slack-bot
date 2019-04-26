package service

import (
	"github.com/slack-bot-4all/slack-bot/cmd/model"
	"github.com/slack-bot-4all/slack-bot/cmd/repository"
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
