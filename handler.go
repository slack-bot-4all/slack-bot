package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/nlopes/slack"
)

// interactionHandler handles interactive message response.
type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
}

const (
	actionSelect           = "select"
	actionConfirm          = "confirm"
	actionCancel           = "cancel"
	actionRestartContainer = "restart-container"
	actionLogsContainer    = "logs-container"
	actionHaproxyCfgUpdate = "update-haproxy"
)

func (h interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := message.Actions[0]
	switch action.Name {
	case actionSelect:
		switch message.CallbackID {
		case "restart-container":
			actionRestartContainerFunction(message, w)
		case "logs-container":
			actionLogsContainerFunction(message, w)
		default:
			return
		}
	case actionCancel:
		title := fmt.Sprintf(":x: @%s cancelou a requisição", message.User.Name)
		responseMessage(w, message.OriginalMessage, title, "")
		getAPIConnection().client.DeleteMessage(message.Channel.ID, message.MessageTs)
	default:
		log.Printf("[ERROR] Ação inválida: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func actionRestartContainerFunction(message slack.AttachmentActionCallback, w http.ResponseWriter) {
	value := message.Actions[0].SelectedOptions[0].Value
	rancherListener.RestartContainer(value)

	title := fmt.Sprintf("Container de ID %s restartado por @%s com sucesso! :sunglasses:\n\n", value, message.User.Name)
	sendMessage(title)

	getAPIConnection().client.DeleteMessage(message.Channel.ID, message.MessageTs)
}

func actionLogsContainerFunction(message slack.AttachmentActionCallback, w http.ResponseWriter) {
	value := message.Actions[0].SelectedOptions[0].Value
	fileName := rancherListener.LogsContainer(value)

	time.Sleep(2 * time.Second)

	api := getAPIConnection()

	log.Println(fileName)

	file, err := api.client.UploadFile(slack.FileUploadParameters{
		File:     fileName,
		Filetype: "text",
		Channels: []string{
			api.channelID,
		},
	})
	CheckErr("Erro ao fazer upload de arquivo de logs de container", err)

	originalMessage := message.OriginalMessage
	originalMessage.Files = []slack.File{
		{
			ID:       file.ID,
			Title:    fmt.Sprintf("Logs do container: %s", value),
			Filetype: "text",
		},
	}
	originalMessage.Attachments = []slack.Attachment{}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&originalMessage)

	getAPIConnection().client.DeleteMessage(message.Channel.ID, message.MessageTs)
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}

func sendMessage(message string) {
	conn := getAPIConnection()

	conn.client.PostMessage(conn.channelID, slack.MsgOptionAttachments(slack.Attachment{
		Text:  message,
		Color: "#0C648A",
	}))
}

func getAPIConnection() *SlackListener {
	c := slack.New(SlackBotToken)

	s := &SlackListener{
		client:    c,
		botID:     SlackBotID,
		channelID: SlackBotChannel,
	}

	return s
}
