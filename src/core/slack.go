// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cayohollanda/runner"
	"github.com/slack-bot-4all/slack-bot/src/model"
	"github.com/slack-bot-4all/slack-bot/src/repository"
	"github.com/slack-bot-4all/slack-bot/src/service"

	"github.com/nlopes/slack"
	"github.com/tidwall/gjson"
)

const (
	canaryUpdate        = "canary-update"
	canaryDisable       = "canary-disable"
	canaryActivate      = "canary-enable"
	canaryInfo          = "canary-info"
	canaryUpTen         = "canary-up"
	haproxyList         = "lb-list"
	logsContainer       = "container-logs"
	restartContainer    = "container-restart"
	getServiceInfo      = "service-info"
	upgradeService      = "service-upgrade"
	listService         = "service-list"
	startService        = "service-start"
	stopService         = "service-stop"
	checkServiceHealth  = "task-add"
	removeServiceCheck  = "task-stop"
	listAllRunningTasks = "task-list"
	listAllEnvironments = "env-list"
	selectEnvironment   = "env-set"
	selectRancher       = "rancher-set"
	listRancher         = "rancher-list"
	commands            = "commands"
)

// SlackListener é a struct que armazena dados do BOT
type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

var (
	rancherListener *RancherListener
	tasks           []*runner.Task
)

// StartBot é a função que inicia o BOT e o prepara para receber eventos de mensagens
func (s *SlackListener) StartBot(rList *RancherListener) {
	log.Println("[INFO] Initializating BOT...")

	rancherListener = rList

	rtm := s.client.NewRTM()
	go rtm.ManageConnection()

	log.Println("[INFO] BOT connection successful!")

	// task := runner.Go(func(shouldStop runner.S) error {
	// 	defer func() {}()

	// 	for {
	// 		s.executeTasks()
	// 		time.Sleep(time.Minute * 2)
	// 	}
	// })
	// task.Running()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			s.client.PostMessage(s.channelID, slack.MsgOptionText("Hey brow, I'm here! Cry your tears :sob:", false))
			log.Println("[INFO] BOT started successfully!")
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
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s>", s.botID)) && !isReminder {
		return nil
	}

	var message string
	messageSlice := strings.Split(ev.Msg.Text, " ") // Tirando a menção ao BOT da mensagem e guardando em uma variável
	if len(messageSlice) <= 1 {
		if ev.Msg.Text != fmt.Sprintf("<@%s>", s.botID) {
			return nil
		}
	} else {
		message = messageSlice[1]
	}

	if strings.Contains(ev.Msg.Text, "help") {
		s.slackCommandHelper(ev, message)
		return nil
	}

	log.Printf("[INFO] New received message: %s", message)

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
	} else if strings.HasPrefix(message, startService) {
		s.slackStartService(ev)
	} else if strings.HasPrefix(message, stopService) {
		s.slackStopService(ev)
	} else if strings.HasPrefix(message, checkServiceHealth) {
		s.slackCheckServiceHealth(ev)
	} else if strings.HasPrefix(message, removeServiceCheck) {
		s.stopServiceCheck(ev)
	} else if strings.HasPrefix(message, listAllRunningTasks) {
		s.listAllRunningTasks(ev)
	} else if strings.HasPrefix(message, selectRancher) {
		s.selectRancher(ev)
	} else if strings.HasPrefix(message, listAllEnvironments) {
		s.listAllEnvironments(ev)
	} else if strings.HasPrefix(message, listRancher) {
		s.listAllRanchers(ev)
	} else if strings.HasPrefix(message, selectEnvironment) {
		s.selectEnvironment(ev)
	} else if strings.HasPrefix(message, commands) {
		s.slackHelper(ev)
	} else if strings.HasPrefix(message, canaryUpTen) {
		s.slackCanaryUpTen(ev)
	} else {
		s.interactiveMessage(ev)
	}

	return nil
}

func (s *SlackListener) executeTasks() {
	var stackName string
	var serviceName string
	var stackID string
	var serviceID string
	var serviceState string
	var containers []Container

	tasks, err := service.ListTask()
	if err != nil {
		log.Println("[ERROR] Error on execute task check, no response from database")
		return
	}

	for _, task := range tasks {
		rancherListener := &RancherListener{
			baseURL:   task.RancherURL,
			accessKey: task.RancherAccessKey,
			secretKey: task.RancherSecretKey,
			projectID: task.RancherProjectID,
		}

		argSplitted := strings.Split(task.Service, "/")

		if len(argSplitted) >= 2 {
			stackName = argSplitted[0]
			serviceName = argSplitted[1]
		} else {
			log.Println("Error! service name is not declared right. Right declaration example: stackName/serviceName")
			return
		}

		respAllStacks := rancherListener.GetStacks()

		dataStack := gjson.Get(respAllStacks, "data")
		dataStack.ForEach(func(key, value gjson.Result) bool {
			if value.Get("name").String() == stackName {
				stackID = value.Get("id").String()
			}
			return true
		})

		respAllServicesFromStack := rancherListener.GetServicesFromStack(stackID)

		dataService := gjson.Get(respAllServicesFromStack, "data")
		dataService.ForEach(func(key, value gjson.Result) bool {
			if value.Get("name").String() == serviceName {
				serviceID = value.Get("id").String()
				serviceState = value.Get("healthState").String()
			}
			return true
		})

		if stackID == "" || serviceID == "" {
			log.Println("Error! Check if you are passing correct argument, the correct is: @bot command stackName/serviceName")
			return
		}

		respAllInstances := rancherListener.GetInstances(serviceID)
		dataInstances := gjson.Get(respAllInstances, "data")
		dataInstances.ForEach(func(key, value gjson.Result) bool {
			var container Container
			container.ID = value.Get("id").String()
			container.Name = value.Get("name").String()
			container.State = value.Get("healthState").String()

			containers = append(containers, container)

			return true
		})

		if serviceState != "healthy" {
			var downContainers []Container
			var upContainers []Container

			var msg string

			for _, container := range containers {
				if container.State == "healthy" {
					upContainers = append(upContainers, container)
				} else {
					downContainers = append(downContainers, container)
				}
				msg += fmt.Sprintf("`%s` - `%s`\n", container.Name, container.State)
			}

			s.client.PostMessage(task.ChannelToSendAlert, slack.MsgOptionText(fmt.Sprintf("Please, check the containers health, the service `%s/%s` actually is `%s` with `%d` up containers and `%d` down containers\n\n%s", stackName, serviceName, serviceState, len(upContainers), len(downContainers), msg), true))
		}
	}
}

func (s *SlackListener) listAllRanchers(ev *slack.MessageEvent) {
	ranchers, err := service.ListRancher()
	if err != nil {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error, verify if database is active", false))
		return
	}

	msg := "*Registered Ranchers:*\n\n"
	for _, rancher := range ranchers {
		msg += fmt.Sprintf("Name: `%s`\nURL: `%s`\nAccess Key: `%s`\n\n", rancher.Name, rancher.URL, rancher.AccessKey)
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) selectEnvironment(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		var idEnv string
		var haveEnv bool

		environment := args[2]
		environment = strings.Replace(environment, "_", " ", -1)

		resp := rancherListener.GetAllEnvironmentsFromRancher()

		data := gjson.Get(resp, "data")
		data.ForEach(func(key, value gjson.Result) bool {
			if value.Get("name").String() == environment {
				idEnv = value.Get("id").String()
				haveEnv = true
			}

			return true
		})

		if haveEnv && idEnv != "" {
			rancherListener.projectID = idEnv
			s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Environment `%s` selected successfully!", environment), false))
			return
		}

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Error on select environment `%s`, check if it exists!", environment), false))
	}
}

func (s *SlackListener) listAllEnvironments(ev *slack.MessageEvent) {
	resp := rancherListener.GetAllEnvironmentsFromRancher()

	msg := "*Environments from this Rancher:*\n\n"

	data := gjson.Get(resp, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		msg += fmt.Sprintf("`%s`\n", value.Get("name").String())

		return true
	})

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) selectRancher(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		rancherInstance := args[2]

		var rancher model.Rancher
		rancher.Name = rancherInstance

		err := repository.FindRancherByName(&rancher)
		if err != nil {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Error on select Rancher `%s`, make sure it is registered!", rancherInstance), false))
			return
		}

		rancherListener.ID = rancher.ID
		rancherListener.baseURL = rancher.URL
		rancherListener.accessKey = rancher.AccessKey
		rancherListener.secretKey = rancher.SecretKey

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Rancher `%s` selected successfully!", rancherInstance), false))
	}
}

func (s *SlackListener) listAllRunningTasks(ev *slack.MessageEvent) {

	msg := "*Running Tasks List:* \n\n"

	var tasks []model.Task
	err := repository.ListTask(&tasks)
	if err != nil {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on check running tasks. Verify if the BOT have connection with database", false))
		return
	}

	for _, task := range tasks {
		var envName string
		resp := rancherListener.GetAllEnvironmentsFromRancher()

		data := gjson.Get(resp, "data")
		data.ForEach(func(key, value gjson.Result) bool {
			if value.Get("id").String() == task.RancherProjectID {
				envName = value.Get("name").String()
			}

			return true
		})
		if string(task.ID) != "" {
			msg += fmt.Sprintf("*%d* / %s - Environment `%s`\n", task.ID, task.Service, envName)
		}
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) stopServiceCheck(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		taskIDToStop := args[2]

		var tasks []model.Task
		err := repository.ListTask(&tasks)

		if err != nil {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Failed to stop task `%s`, check if this task is already running or if the database is running", taskIDToStop), false))
			return
		}

		var taskToStop model.Task
		for _, task := range tasks {
			if fmt.Sprintf("%d", task.ID) == taskIDToStop {
				taskToStop = task
			}
		}

		if taskToStop.Service == "" {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Failed to stop task `%s`. Verify if this task is running", taskIDToStop), false))
			return
		}

		err = service.DeleteTask(taskToStop)
		if err != nil {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Failed to stop task `%s`", taskIDToStop), false))
			return
		}

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Task *%s*/`%s` stopped successfully!", taskIDToStop, taskToStop.Service), false))
	}
}

func (s *SlackListener) slackCheckServiceHealth(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")
	if len(args) == 4 {
		task := &model.Task{
			Service:            args[2],
			ChannelToSendAlert: args[3],
			RancherURL:         rancherListener.baseURL,
			RancherAccessKey:   rancherListener.accessKey,
			RancherSecretKey:   rancherListener.secretKey,
			RancherProjectID:   rancherListener.projectID,
		}

		err := service.AddTask(task)
		if err != nil {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on register task, verify if BOT haves connection with database", false))
		} else {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Task added successfully!", false))
		}
	}
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

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("ConfigHaprox:\n\n\n%s\n", msg), true))
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
	msg := "*Commands:*\n\n"

	for _, cmd := range Commands {
		msg += fmt.Sprintf("`%s` -> %s\n", cmd.Cmd, cmd.Description)
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
	var channelToSendMessage string

	args := strings.Split(ev.Msg.Text, " ")

	if len(args) < 5 {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Command call error, correct syntax: @name-of-bot %s LB-id new-version-weight old-version-weight", canaryUpdate), false))
		return
	}

	lb := args[2]
	newVersionPercent := args[3]
	oldVersionPercent := args[4]

	if len(args) == 6 {
		channelToSendMessage = args[5]
	}

	resp := rancherListener.UpdateCustomHaproxyCfg(lb, newVersionPercent, oldVersionPercent)

	if resp == "error" {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on update haproxy.cfg, check if ID param is right, the body of haproxy.cfg is empty or if weights not sum 100", false))
		return
	}
	//v := strconv.FormatBool(resp)
	s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("File 'haproxy.cfg' updated successfuly!\n```%s```", resp), false))

	if channelToSendMessage != "" {
		resp := rancherListener.GetService(lb)

		serviceName := gjson.Get(resp, "name").String()

		s.client.PostMessage(channelToSendMessage, slack.MsgOptionText(fmt.Sprintf("Canary of `%s` has been updated.\nNew version: `%s`\nOld version: `%s`", serviceName, newVersionPercent, oldVersionPercent), false))
	}
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

	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		id := args[2]

		rancherListener.RestartContainer(id)

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Container restarted"), true))
	} else {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Parameters is required"), true))
	}

	// s.createAndSendAttachment(
	// 	ev,
	// 	"Which container you need restart? :yum:",
	// 	restartContainer,
	// 	getContainers(),
	// 	nil,
	// )
}

func (s *SlackListener) slackStartService(ev *slack.MessageEvent) {

	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		id := args[2]

		rancherListener.StartService(id)

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Service started"), true))
	} else {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Parameters is required"), true))
	}
}

func (s *SlackListener) slackStopService(ev *slack.MessageEvent) {

	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		id := args[2]

		rancherListener.StopService(id)

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Service stopped"), true))
	} else {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Parameters is required"), true))
	}
}

func (s *SlackListener) interactiveMessage(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) >= 0 {
		client := createHTTPClient()

		req, err := http.NewRequest("GET", "https://api.kanye.rest", nil)
		CheckErr("", err)

		resp, err := client.Do(req)
		CheckErr("", err)

		body, _ := ioutil.ReadAll(resp.Body)

		var kanye Kanye
		_ = json.Unmarshal(body, &kanye)

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Little Friend, what did you mean? I do not understand, use @jeremias help or @jeremias commands!\n\n So here's a message to make your day better:\n\n\"%s\"", kanye.Quote), true))
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
		container.ID = value.Get("id").String()
		container.ImageUUID = value.Get("imageUuid").String()
		container.Name = value.Get("name").String()
		containers = append(containers, container)

		return true
	})

	// Criando lista de opções, fazendo um ForEach na lista
	// de structs de containers, criando opcao dentro do ForEach
	// e adicionando à lista de opcoes
	opcoes := []slack.AttachmentActionOption{}
	for _, container := range containers {
		opcoes = append(opcoes, slack.AttachmentActionOption{
			Text:  fmt.Sprintf("%s | %s", container.ID, container.Name),
			Value: container.ID,
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
func (s *SlackListener) slackCanaryUpTen(ev *slack.MessageEvent) {
	var channelToSendMessage string

	args := strings.Split(ev.Msg.Text, " ")
	if len(args) < 3 {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Command call error, correct syntax: @name-of-bot %s canaryUpTen LB-id channel-to-send-alert (optional)", canaryUpTen), false))
		return
	}
	if len(args) == 4 {
		channelToSendMessage = args[3]
	}
	lb := args[2]

	new, old := rancherListener.SearchForLbPercent(lb)

	newToInt, _ := strconv.Atoi(new)
	oldToInt, _ := strconv.Atoi(old)

	newMoreTen := newToInt + 10
	oldLessTen := oldToInt - 10

	newToString := strconv.Itoa(newMoreTen)
	oldToString := strconv.Itoa(oldLessTen)

	resp := rancherListener.UpdateCustomHaproxyCfg(lb, newToString, oldToString)

	if resp == "error" {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Error on update haproxy.cfg, check if ID param is right, the body of haproxy.cfg is empty or if weights not sum 100", false))
		return
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("File 'haproxy.cfg' updated successfuly!\n```%s```", resp), false))

	if channelToSendMessage != "" {
		resp := rancherListener.GetService(lb)

		serviceName := gjson.Get(resp, "name").String()

		s.client.PostMessage(channelToSendMessage, slack.MsgOptionText(fmt.Sprintf("Canary of `%s` has been updated.\nNew version: `%s`\nOld version: `%s`", serviceName, newToString, oldToString), false))
	}
}
