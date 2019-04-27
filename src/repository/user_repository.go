package repository

import (
	"github.com/slack-bot-4all/slack-bot/src/config"
	"github.com/slack-bot-4all/slack-bot/src/model"
)

// AddUser : add a User to database
func AddUser(u *model.User) (err error) {
	var hash config.Hash
	if u.Password, err = hash.Generate(u.Password); err != nil {
		return err
	}

	if err := config.DB.Create(u).Error; err != nil {
		return err
	}

	return nil
}

// FindUserByUsername : consults the db with the username
func FindUserByUsername(u *model.User) (err error) {
	if err := config.DB.Where("username = ?", u.Username).First(u).Error; err != nil {
		return err
	}

	return nil
}
