package main

import (
	"fmt"
	"log"

	"github.com/slack-bot-4all/slack-bot/cmd/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/slack-bot-4all/slack-bot/cmd/config"
	"github.com/slack-bot-4all/slack-bot/cmd/core"
)

func main() {
	core.PrintLogoOnConsole()

	err := initializeDB()
	if err != nil {
		log.Fatalf("[ERROR] Error to connect on database\n%s", err.Error())
	}
	core.Start()
}

func initializeDB() error {
	var err error
	config.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", core.DatabaseUsername, core.DatabasePassword, core.DatabaseURL, core.DatabaseSchema))

	// u := model.User{
	// 	Username: "admin",
	// 	Password: "admin",
	// }
	// err = repository.AddUser(&u)

	if err != nil {
		return err
	}

	log.Println("[INFO] Connected to database")

	config.DB.AutoMigrate(&model.Rancher{}, &model.User{})

	return nil
}
