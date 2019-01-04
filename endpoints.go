// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"encoding/json"
	"net/http"
)

// Env usado para retornar o objeto env
type Env struct {
	RancherAccessKey          string `json:"RANCHER_ACCESS_KEY"`
	RancherSecretKey          string `json:"RANCHER_SECRET_KEY"`
	RancherBaseURL            string `json:"RANCHER_BASE_URL"`
	RancherProjectID          string `json:"RANCHER_PROJECT_ID"`
	SlackBotToken             string `json:"SLACK_BOT_TOKEN"`
	SlackBotID                string `json:"SLACK_BOT_ID"`
	SlackBotChannel           string `json:"SLACK_BOT_CHANNEL"`
	SlackBotVerificationToken string `json:"SLACK_BOT_VERIFICATION_TOKEN"`
	HTTPPort                  string `json:"HTTP_PORT"`
	SplunkUsername            string `json:"SPLUNK_USERNAME"`
	SplunkPassword            string `json:"SPLUNK_PASSWORD"`
	SplunkBaseURL             string `json:"SPLUNK_BASE_URL"`
}

var env []Env

// GetEnvs mostra todas envs
func GetEnvs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(env)
}
