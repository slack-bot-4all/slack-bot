// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/tidwall/gjson"
)

const (
	canaryUpdate     = "update-canary"
	canaryDisable    = "disable-canary"
	canaryActivate   = "enable-canary"
	canaryInfo       = "info-canary"
	haproxyList      = "list-lb"
	logsContainer    = "logs-container"
	restartContainer = "restart-container"
	getServiceInfo   = "info-service"
	upgradeService   = "upgrade-service"
	listService      = "list-service"
	commands         = "commands"
)

// SlackListener é a struct que armazena dados do BOT
type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

var rancherListener *RancherListener

// StartBot é a função que inicia o BOT e o prepara para receber eventos de mensagens
func (s *SlackListener) StartBot(rList *RancherListener) {
	log.Println("[INFO] Iniciando o BOT...")

	rancherListener = rList

	rtm := s.client.NewRTM()
	go rtm.ManageConnection()

	log.Println("[INFO] Conexão com o BOT feita com sucesso!")

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			s.client.PostMessage(s.channelID, slack.MsgOptionText("Hey brow, I'm here! Cry your tears :sob:", false))
			log.Println("[INFO] BOT iniciado com sucesso!")
		case *slack.MessageEvent:
			s.handleMessageEvent(ev)
		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Parando a função caso a msg não venha do mesmo canal que o BOT está
	if ev.Channel != s.channelID {
		return nil
	}

	// Parando a função caso a msg tenha vindo do BOT
	if ev.User == s.botID {
		log.Println("[INFO] Mensagem não recebida. Mensagem vinda do BOT")
		return nil
	}

	var isReminder bool
	if strings.Contains(ev.Msg.Text, fmt.Sprintf("Reminder: <@%s", s.botID)) {
		ev.Msg.Text = strings.Replace(ev.Msg.Text, "Reminder: ", "", 1)
		ev.Msg.Text = RemoveLastCharacter(ev.Msg.Text)
		isReminder = true
	}

	// Parando a função caso a mensagem não traga o prefixo mencionando o BOT
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) && !isReminder {
		return nil
	}

	// Tirando a menção ao BOT da mensagem e guardando em uma variável
	message := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1]

	if strings.Contains(ev.Msg.Text, "help") {
		s.slackCommandHelper(ev, message)
		return nil
	}

	// Fazendo as verificações de mensagens e jogando
	// para as devidas funções
	if strings.HasPrefix(message, restartContainer) {
		s.slackRestartContainer(ev)
	} else if strings.HasPrefix(message, logsContainer) {
		s.slackLogsContainer(ev)
	} else if strings.HasPrefix(message, canaryUpdate) {
		s.slackUpdateCanary(ev)
	} else if strings.HasPrefix(message, haproxyList) {
		s.slackListLoadBalancers(ev)
	} else if strings.HasPrefix(message, getServiceInfo) {
		s.slackServiceInfo(ev)
	} else if strings.HasPrefix(message, listService) {
		s.slackServicesList(ev)
	} else if strings.HasPrefix(message, upgradeService) {
		s.slackServiceUpgrade(ev)
	} else if strings.HasPrefix(message, canaryDisable) {
		s.slackCanaryDisable(ev)
	} else if strings.HasPrefix(message, canaryActivate) {
		s.slackCanaryEnable(ev)
	} else if strings.HasPrefix(message, canaryInfo) {
		s.slackCanaryInfo(ev)
	} else if strings.HasPrefix(message, commands) {
		s.slackHelper(ev)
	} else {
		s.interactiveMessage(ev)
	}

	return nil
}

func (s *SlackListener) slackCanaryInfo(ev *slack.MessageEvent) {

	args := strings.Split(ev.Msg.Text, " ")
	if len(args) == 3 {
		lbid := args[2]

		resp := rancherListener.GetHaproxyCfg(lbid)
		lbConfig := gjson.Get(resp, "lbConfig.config").String()

		msg := fmt.Sprintf("haproxy.cfg file of Load Balancer `%s`.\n```%s```",
			lbid, lbConfig)

		if resp == "error" {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error", false))
			return
		}
		fmt.Println(msg)
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("ConfigHaprox:\n\n\"%s\"", msg), true))
	}
}

func (s *SlackListener) slackCanaryEnable(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		lb := args[2]

		resp := rancherListener.EnableCanary(lb)

		if resp == "error" {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on update haproxy.cfg, check if ID param is right or the body of haproxy.cfg is empty", false))
			return
		}

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("File 'haproxy.cfg' updated success! *Canary Deployment* enabled.\n```%s```", resp), false))
	} else {
		s.createAndSendAttachment(
			ev,
			"Which Load Balancer you need enable the Canary?",
			canaryActivate,
			getLbOptions(),
			&slack.ConfirmationField{
				Title:       "Are you sure?",
				Text:        "You sure to enable Canary? :thinking_face:",
				OkText:      "Yes",
				DismissText: "No",
			},
		)
	}

}

func (s *SlackListener) slackCanaryDisable(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		lb := args[2]

		resp := rancherListener.DisableCanary(lb)

		if resp == "error" {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on update haproxy.cfg, check if ID param is right or the body of haproxy.cfg is empty", false))
			return
		}

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("File 'haproxy.cfg' updated success! *Canary Deployment* disabled.\n```%s```", resp), false))
	} else {
		s.createAndSendAttachment(
			ev,
			"Which Load Balancer you need disable the Canary?",
			canaryDisable,
			getLbOptions(),
			&slack.ConfirmationField{
				Title:       "Are you sure?",
				Text:        "You sure to disable Canary? :scream:",
				OkText:      "Yes",
				DismissText: "No",
			},
		)
	}

}

func (s *SlackListener) slackServiceUpgrade(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) != 4 {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Command call error, correct syntax: @name-of-bot %s service-id new-image", upgradeService), false))
		return
	}

	serviceID := args[2]
	newServiceImage := args[3]

	if !strings.HasPrefix(newServiceImage, "docker:") {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Image name needs to start with 'docker:'. Ex.: docker:ubuntu:14.04", false))
		return
	}

	resp := rancherListener.UpgradeService(serviceID, newServiceImage)

	if resp == "" {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Service upgrade error. Check:\n*- If service ID really exists*\n*- If service is not in upgrading state*", false))
		return
	}

	msg := fmt.Sprintf("Service updated successfuly! New image of the service `%s` is `%s`", serviceID, resp)

	log.Printf("[INFO] Service %s updated by %s\n", serviceID, ev.Msg.User)
	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackServicesList(ev *slack.MessageEvent) {
	resp := rancherListener.ListServices()

	msg := "*Service List:* \n\n"

	data := gjson.Get(resp, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		msg += fmt.Sprintf("`%s | %s`\n", value.Get("id").String(), value.Get("name").String())
		return true
	})

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackServiceInfo(ev *slack.MessageEvent) {
	s.createAndSendAttachment(
		ev,
		"Which service you need informations? :sunglasses:",
		getServiceInfo,
		getServices(),
		nil,
	)
}

func (s *SlackListener) slackCommandHelper(ev *slack.MessageEvent, message string) {
	var msg string

	for _, cmd := range Commands {
		if cmd.Cmd == message {
			cmd.Usage = strings.Replace(cmd.Usage, "command", cmd.Cmd, 1)
			msg = fmt.Sprintf("*Command:* `%s`\n*Description:* _%s_\n*Usage:* _%s_\n*Lint:* _%s_", cmd.Cmd, cmd.Description, cmd.Usage, cmd.Lint)
		}
	}

	if msg == "" {
		msg = "Command not found."
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackHelper(ev *slack.MessageEvent) {
	msg := "*Commands:* "

	for _, cmd := range Commands {
		msg += fmt.Sprintf("`%s` ", cmd.Cmd)
	}

	msg += "\n\n_*PS.:* If you need detailed informations for a command, you can call command followed by *help*._\n_*Ex.:* @jeremias command help_"

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackListLoadBalancers(ev *slack.MessageEvent) {
	loadBalancers := rancherListener.GetLoadBalancers()

	var lines []string

	for _, lb := range loadBalancers {
		line := fmt.Sprintf("`%s | %s`", lb.ID, lb.Name)

		lines = append(lines, line)
	}

	msg := "*Load Balancers list:*"

	for _, line := range lines {
		msg += fmt.Sprintf("\n%s", line)
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackUpdateCanary(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) != 5 {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Command call error, correct syntax: @name-of-bot %s LB-id new-version-weight old-version-weight", canaryUpdate), false))
		return
	}

	lb := args[2]
	newVersionPercent := args[3]
	oldVersionPercent := args[4]

	resp := rancherListener.UpdateCustomHaproxyCfg(lb, newVersionPercent, oldVersionPercent)

	if resp == "error" {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on update haproxy.cfg, check if ID param is right, the body of haproxy.cfg is empty or if weights not sum 100", false))
		return
	}
	//v := strconv.FormatBool(resp)
	s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("File 'haproxy.cfg' updated successfuly!\n```%s```", resp), false))
}

func (s *SlackListener) slackLogsContainer(ev *slack.MessageEvent) {

	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		container := args[2]

		fileName := rancherListener.LogsContainer(container)

		time.Sleep(2 * time.Second)

		api := getAPIConnection()

		_, err := api.client.UploadFile(slack.FileUploadParameters{
			File:     fileName,
			Filename: fileName,
			Filetype: "text",
			Channels: []string{
				api.channelID,
			},
		})
		CheckErr("Upload logs container error", err)

		// FileSlack := []slack.File{
		// 	{
		// 		ID:       file.ID,
		// 		Title:    fmt.Sprintf("Logs of container: %s", container),
		// 		Filetype: "text",
		// 	},
		// }
		// s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on update haproxy.cfg, check if ID param is right or the body of haproxy.cfg is empty", true))
	} else {
		s.createAndSendAttachment(
			ev,
			"Which container you need to download logs? :yum:",
			logsContainer,
			getContainers(),
			nil,
		)

	}

}

func (s *SlackListener) slackRestartContainer(ev *slack.MessageEvent) {
	s.createAndSendAttachment(
		ev,
		"Which container you need restart? :yum:",
		restartContainer,
		getContainers(),
		nil,
	)
}

func (s *SlackListener) interactiveMessage(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) >= 0 {
		fmt.Println(args)
		client := createHTTPClient()

		req, err := http.NewRequest("GET", "https://api.kanye.rest", nil)
		CheckErr("", err)

		resp, err := client.Do(req)
		CheckErr("", err)

		body, _ := ioutil.ReadAll(resp.Body)

		var kanye Kanye
		_ = json.Unmarshal(body, &kanye)

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Friend, what did you mean? I do not understand, so here's a message to make your day better:\n\n\"%s\"", kanye.Quote), true))
	}

}

func (s *SlackListener) createAndSendAttachment(ev *slack.MessageEvent, text string, callbackID string, options []slack.AttachmentActionOption, confirmation *slack.ConfirmationField) {
	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(slack.Attachment{
		Text:       text,
		Color:      "#0C648A",
		CallbackID: callbackID,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: options,
				Confirm: confirmation,
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}))
}

func getContainers() []slack.AttachmentActionOption {
	// Pegando a lista de containers lá do rancher.go
	containersList := rancherListener.ListContainers()

	// Criando uma lista de estruturas
	containers := []*Container{}

	// Pegando a lista que veio da API do Rancher, convertendo pra String
	// atribuindo à uma struct e adicionando na lista de structs
	data := gjson.Get(containersList, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		container := new(Container)
		container.id = value.Get("id").String()
		container.imageUUID = value.Get("imageUuid").String()
		container.name = value.Get("name").String()
		containers = append(containers, container)

		return true
	})

	// Criando lista de opções, fazendo um ForEach na lista
	// de structs de containers, criando opcao dentro do ForEach
	// e adicionando à lista de opcoes
	opcoes := []slack.AttachmentActionOption{}
	for _, container := range containers {
		opcoes = append(opcoes, slack.AttachmentActionOption{
			Text:  fmt.Sprintf("%s | %s", container.id, container.name),
			Value: container.id,
		})
	}

	return opcoes
}

func getServices() []slack.AttachmentActionOption {
	servicesList := rancherListener.ListServices()

	opcoes := []slack.AttachmentActionOption{}

	data := gjson.Get(servicesList, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		serviceID := value.Get("id").String()
		serviceName := value.Get("name").String()
		opcoes = append(opcoes, slack.AttachmentActionOption{
			Text:  fmt.Sprintf("%s | %s", serviceID, serviceName),
			Value: serviceID,
		})

		return true
	})

	return opcoes
}

func getLbOptions() []slack.AttachmentActionOption {
	opcoes := []slack.AttachmentActionOption{}
	for _, lb := range rancherListener.GetLoadBalancers() {
		opcoes = append(opcoes, slack.AttachmentActionOption{
			Text:  fmt.Sprintf("%s | %s", lb.ID, lb.Name),
			Value: lb.ID,
		})
	}

	return opcoes
}
