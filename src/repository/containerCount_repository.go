package repository

import (
	"github.com/slack-bot-4all/slack-bot/src/config"
	"github.com/slack-bot-4all/slack-bot/src/model"
)

// ChangeToZeroCounter ::
func ChangeToZeroCounter(counter *model.ContainerCount) error {
	counter.Count = 0

	if err := config.DB.Save(counter).Error; err != nil {
		return err
	}

	return nil
}

// CreateCounterToContainer ::
func CreateCounterToContainer(counter *model.ContainerCount) error {
	if err := config.DB.Create(counter).Error; err != nil {
		return err
	}

	return nil
}

// CreateCounterToService ::
func CreateCounterToService(counter *model.ContainerCount) error {
	if err := config.DB.Create(counter).Error; err != nil {
		return err
	}

	return nil
}

// GetCounterByContainerID ::
func GetCounterByContainerID(counter *model.ContainerCount, containerID string) error {

	if err := config.DB.Where("container_id = ?", containerID).Find(&counter).Error; err != nil {
		return err
	}

	return nil
}

// IncrementCounterByContainerID ::
func IncrementCounterByContainerID(containerID string) error {

	var counter model.ContainerCount

	if err := config.DB.Where("container_id = ?", containerID).Find(&counter).Error; err != nil {
		return err
	}

	counter.Count = counter.Count + 1

	if err := config.DB.Save(counter).Error; err != nil {
		return err
	}

	return nil

}
