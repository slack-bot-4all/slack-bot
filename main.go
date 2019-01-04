// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
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
)

func main() {
	File := os.Getenv("FILE")

	FileOpen, err := os.Open(File)
	CheckErr("Erro ao abrir o arquivo de environments", err)

	scanner := bufio.NewScanner(FileOpen)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	CheckErr("Erro ao scannear arquivo de environments", scanner.Err())

	for _, line := range lines {
		chave := strings.Split(line, "=")[0]
		valor := strings.Split(line, "=")[1]

		switch chave {
		case "RANCHER_ACCESS_KEY":
			RancherAccessKey = valor
		case "RANCHER_SECRET_KEY":
			RancherSecretKey = valor
		case "RANCHER_BASE_URL":
			RancherBaseURL = valor
		case "RANCHER_PROJECT_ID":
			RancherProjectID = valor
		case "SLACK_BOT_TOKEN":
			SlackBotToken = valor
		case "SLACK_BOT_ID":
			SlackBotID = valor
		case "SLACK_BOT_CHANNEL":
			SlackBotChannel = valor
		case "SLACK_BOT_VERIFICATION_TOKEN":
			SlackBotVerificationToken = valor
		case "HTTP_PORT":
			Port = valor
		case "SPLUNK_USERNAME":
			SplunkUsername = valor
		case "SPLUNK_PASSWORD":
			SplunkPassword = valor
		case "SPLUNK_BASE_URL":
			SplunkBaseURL = valor
		}

		envs = append(envs, Env{Key: chave, Value: valor})
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

	client := slack.New(
		SlackBotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(mw, "SLfR: ", log.Lshortfile|log.LstdFlags)),
	)

	slackListener := &SlackListener{
		client:    client,
		botID:     SlackBotID,
		channelID: SlackBotChannel,
	}

	rancherListener := &RancherListener{
		accessKey: RancherAccessKey,
		secretKey: RancherSecretKey,
		baseURL:   RancherBaseURL,
		projectID: RancherProjectID,
	}

	go slackListener.StartBot(rancherListener)

	router := mux.NewRouter()

	router.HandleFunc("/env", GetEnvs).Methods("GET")
	router.Handle("/interaction", interactionHandler{
		verificationToken: SlackBotVerificationToken,
	})

	log.Printf("[INFO] Servidor rodando na porta: %s", Port)
	if err := http.ListenAndServe(":"+Port, router); err != nil {
		log.Printf("[ERROR] %s", err)
		os.Exit(1)
	}
}
