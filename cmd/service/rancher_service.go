package service

import (
	"github.com/slack-bot-4all/slack-bot/cmd/model"
	"github.com/slack-bot-4all/slack-bot/cmd/repository"
)

// AddRancher : have a business rules to add a Rancher to db
func AddRancher(r *model.Rancher) error {
	var err error

	if r.Name != "" && r.URL != "" && r.AccessKey != "" && r.SecretKey != "" {
		err = repository.AddRancher(r)
	}

	if err != nil {
		return err
	}

	return nil
}
