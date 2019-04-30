package service

import (
	"github.com/slack-bot-4all/slack-bot/src/model"
	"github.com/slack-bot-4all/slack-bot/src/repository"
)

// AddTask : have a business rules to add a Task to db
func AddTask(t *model.Task) error {
	var err error

	if t.RancherURL != "" && t.RancherAccessKey != "" && t.RancherSecretKey != "" && t.Service != "" {
		err = repository.AddTask(t)
	}

	if err != nil {
		return err
	}

	return nil
}

// ListTask : list all ranchers
func ListTask() (tasksList []model.Task, err error) {
	var tasks []model.Task

	err = repository.ListTask(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// DeleteTask :
func DeleteTask(t model.Task) error {
	if err := repository.DeleteTask(&t); err != nil {
		return err
	}

	return nil
}
