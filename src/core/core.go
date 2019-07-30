// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package core

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
	"github.com/slack-bot-4all/slack-bot/src/config"
	"github.com/slack-bot-4all/slack-bot/src/model"
	"github.com/slack-bot-4all/slack-bot/src/repository"
	"github.com/slack-bot-4all/slack-bot/src/routes"
)

var (
	// RancherAccessKey é a KEY de usuário para acessar na API do Rancher 1.6
	RancherAccessKey string

	// RancherSecretKey é a KEY de senha para acessar na API do Rancher 1.6
	RancherSecretKey string

	// RancherBaseURL é a URL base para acesso aos End-points da API do Rancher 1.6
	RancherBaseURL string

	// RancherProjectID é o ID do projeto base que será usado nas requisições
	RancherProjectID string

	// SlackBotToken é o token que será usado para ter acesso a aplicação do BOT no Slack API
	SlackBotToken string

	// SlackBotID é o ID do BOT que poderá ser usado para futuras comparações de mensagens
	SlackBotID string

	// SlackBotChannel é o canal padrão que o BOT irá escutar
	SlackBotChannel string

	// SlackBotVerificationToken é o Verification Token do BOT que será usado
	// no interactive
	SlackBotVerificationToken string

	// Port é a porta onde a API irá rodar
	Port string

	// SplunkUsername para login no Splunk
	SplunkUsername string

	// SplunkPassword para login no Splunk
	SplunkPassword string

	// SplunkBaseURL para login no Splunk
	SplunkBaseURL string

	// DatabaseUsername is a user to connect on MySQL database
	DatabaseUsername string

	// DatabasePassword is a password to connect on MySQL database
	DatabasePassword string

	// DatabaseURL is a URL to connect on MySQL database
	DatabaseURL string

	// DatabaseSchema is a schema to connect on MySQL database
	DatabaseSchema string

	// StatusCakeChannelID ::
	StatusCakeChannelID string

	RanchListener *RancherListener
)

func init() {
	flag.StringVar(&RancherAccessKey, "rancher_access_key", os.Getenv("RANCHER_ACCESS_KEY"), "Access key to connect on Rancher API")
	flag.StringVar(&RancherSecretKey, "rancher_secret_key", os.Getenv("RANCHER_SECRET_KEY"), "Secret key to connect on Rancher API")
	flag.StringVar(&RancherBaseURL, "rancher_base_url", os.Getenv("RANCHER_BASE_URL"), "Base URL of Rancher API")
	flag.StringVar(&RancherProjectID, "rancher_project_id", os.Getenv("RANCHER_PROJECT_ID"), "Project ID default to API requests")
	flag.StringVar(&SlackBotToken, "slack_bot_token", os.Getenv("SLACK_BOT_TOKEN"), "Slack Bot Token to connect on Slack API")
	flag.StringVar(&SlackBotID, "slack_bot_id", os.Getenv("SLACK_BOT_ID"), "Slack Bot ID to compare messages on channel's")
	flag.StringVar(&SlackBotChannel, "slack_bot_channel", os.Getenv("SLACK_BOT_CHANNEL"), "Channel where the BOT will listen")
	flag.StringVar(&Port, "http_port", os.Getenv("HTTP_PORT"), "HTTP Port where API's gonna run")
	flag.StringVar(&SlackBotVerificationToken, "slack_bot_verification_token", os.Getenv("SLACK_BOT_VERIFICATION_TOKEN"), "Verification token of BOT")
	flag.StringVar(&DatabaseUsername, "database_username", os.Getenv("DATABASE_USERNAME"), "Username of db")
	flag.StringVar(&DatabasePassword, "database_password", os.Getenv("DATABASE_PASSWORD"), "Password of db")
	flag.StringVar(&DatabaseURL, "database_url", os.Getenv("DATABASE_URL"), "URL of db")
	flag.StringVar(&DatabaseSchema, "database_schema", os.Getenv("DATABASE_SCHEMA"), "Schema of db")
	flag.StringVar(&StatusCakeChannelID, "statuscake_channel_id", os.Getenv("STATUSCAKE_CHANNEL_ID"), "Schema of db")
	flag.StringVar(&SplunkBaseURL, "splunk_base_url", os.Getenv("SPLUNK_BASE_URL"), "Schema of db")
	flag.StringVar(&SplunkUsername, "splunk_username", os.Getenv("SPLUNK_USERNAME"), "Schema of db")
	flag.StringVar(&SplunkPassword, "splunk_password", os.Getenv("SPLUNK_PASSWORD"), "Schema of db")
}

// Start : start all proccesses
func Start() {
	// parsing environmnets to variables
	flag.Parse()

	if SlackBotToken == "" || SlackBotID == "" || SlackBotChannel == "" || Port == "" || DatabaseURL == "" || DatabaseUsername == "" || DatabasePassword == "" || DatabaseSchema == "" {
		log.Fatal("[ERROR] To run the BOT, you need to set the environments, questions, see README")
	}

	err := initializeDB()
	if err != nil {
		log.Fatalf("[ERROR] Error to connect on database\n%s", err.Error())
	}

	t := time.Now()
	fileName := fmt.Sprintf("logs/logs-%d%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
	f, err := os.Create(fileName)
	defer f.Close()

	fileOpen, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer fileOpen.Close()

	mw := io.MultiWriter(os.Stdout, fileOpen)

	log.SetOutput(mw)

	log.Println("[INFO] Updating commands...")
	CreateCommands()
	log.Println("[INFO] Commands has been updated!")

	client := slack.New(
		SlackBotToken,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(mw, "SLfR: ", log.Lshortfile|log.LstdFlags)),
	)

	slackListener := &SlackListener{
		client:    client,
		botID:     SlackBotID,
		channelID: SlackBotChannel,
		statusCakeChannelID: StatusCakeChannelID,
	}

	RanchListener = &RancherListener{
		accessKey: RancherAccessKey,
		secretKey: RancherSecretKey,
		baseURL:   RancherBaseURL,
		projectID: RancherProjectID,
	}

	go slackListener.StartBot(RanchListener)

	router := routes.GetRoutes()

	router.Run(fmt.Sprintf(":%s", Port))
}

func initializeDB() error {
	var err error
	config.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", DatabaseUsername, DatabasePassword, DatabaseURL, DatabaseSchema))

	if err != nil {
		return err
	}

	log.Println("[INFO] Connected to database")

	config.DB.AutoMigrate(&model.Rancher{}, &model.User{}, &model.Task{}, &model.ContainerCount{})

	adminUser := model.User{
		Username: "admin",
		Password: "admin",
	}
	err = repository.FindUserByUsername(&adminUser)
	if err != nil {
		repository.AddUser(&adminUser)
	}

	return nil
}
