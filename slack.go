// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"fmt"
	"log"
	"strings"

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
	comandos         = "comandos"
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
			s.client.PostMessage(s.channelID, slack.MsgOptionText("Fala mano, to aqui! :nerd_face:", false))
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

	log.Println(ev.User)

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

	if strings.Contains(ev.Msg.Text, "ajuda") {
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
	} else if strings.HasPrefix(message, comandos) {
		s.slackHelper(ev)
	}

	return nil
}

func (s *SlackListener) slackCanaryInfo(ev *slack.MessageEvent) {
	s.createAndSendAttachment(
		ev,
		"Qual Load Balancer deseja buscar informações do Canary?",
		canaryInfo,
		getLbOptions(),
		nil,
	)
}

func (s *SlackListener) slackCanaryEnable(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) == 3 {
		lb := args[2]

		resp := rancherListener.EnableCanary(lb)

		if resp == "error" {
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Erro ao fazer update no haproxy.cfg, verifique se o ID passado está correto ou se o conteúdo do haproxy.cfg atual está em branco", false))
			return
		}

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Arquivo 'haproxy.cfg' alterado com sucesso! *Canary Deployment* ativado.\n```%s```", resp), false))
	} else {
		s.createAndSendAttachment(
			ev,
			"Qual Load Balancer deseja ativar o Canary?",
			canaryActivate,
			getLbOptions(),
			&slack.ConfirmationField{
				Title:       "Tem certeza disso?",
				Text:        "Deseja mesmo ativar o Canary? :thinking_face:",
				OkText:      "Sim",
				DismissText: "Não",
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
			s.client.PostMessage(ev.Channel, slack.MsgOptionText("Erro ao fazer update no haproxy.cfg, verifique se o ID passado está correto ou se o conteúdo do haproxy.cfg atual está em branco", false))
			return
		}

		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Arquivo 'haproxy.cfg' alterado com sucesso! *Canary Deployment* desativado.\n```%s```", resp), false))
	} else {
		s.createAndSendAttachment(
			ev,
			"Qual Load Balancer deseja desativar o Canary?",
			canaryDisable,
			getLbOptions(),
			&slack.ConfirmationField{
				Title:       "Tem certeza disso?",
				Text:        "Deseja mesmo desativar o Canary? :scream:",
				OkText:      "Sim",
				DismissText: "Não",
			},
		)
	}

}

func (s *SlackListener) slackServiceUpgrade(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) != 4 {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Erro na chamada do comando, sintaxe correta: @nome-do-bot %s id-serviço nova-imagem", upgradeService), false))
		return
	}

	serviceID := args[2]
	newServiceImage := args[3]

	if !strings.HasPrefix(newServiceImage, "docker:") {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("O nome da imagem deve começar com 'docker:'. Ex.: docker:ubuntu:14.04", false))
		return
	}

	resp := rancherListener.UpgradeService(serviceID, newServiceImage)

	if resp == "" {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Erro no upgrade do serviço. Você pode verificar:\n*- Se o ID do serviço que foi passado realmente existe*\n*- Se o serviço já não está passando por um processo de Upgrade*", false))
		return
	}

	msg := fmt.Sprintf("Serviço atualizado com sucesso! A nova imagem do serviço `%s` é `%s`", serviceID, resp)

	log.Printf("[INFO] Serviço %s atualizado pelo usuário %s\n", serviceID, ev.Msg.User)
	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackServicesList(ev *slack.MessageEvent) {
	resp := rancherListener.ListServices()

	msg := "*Lista de serviços:* \n\n"

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
		"Qual serviço deseja obter informações? :sunglasses:",
		getServiceInfo,
		getServices(),
		nil,
	)
}

func (s *SlackListener) slackCommandHelper(ev *slack.MessageEvent, message string) {
	var msg string

	for _, cmd := range Commands {
		if cmd.Cmd == message {
			cmd.Usage = strings.Replace(cmd.Usage, "comando", cmd.Cmd, 1)
			msg = fmt.Sprintf("*Comando:* `%s`\n*Descrição:* _%s_\n*Uso:* _%s_\n*Dica:* _%s_", cmd.Cmd, cmd.Description, cmd.Usage, cmd.Lint)
		}
	}

	if msg == "" {
		msg = "Comando não encontrado."
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackHelper(ev *slack.MessageEvent) {
	msg := "*Comandos:* "

	for _, cmd := range Commands {
		msg += fmt.Sprintf("`%s` ", cmd.Cmd)
	}

	msg += "\n\n_*Obs.:* Caso queira informações mais detalhadas sobre um comando, você pode chamar este comando seguido de *ajuda*._\n_*Ex.:* @bot comando ajuda_"

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackListLoadBalancers(ev *slack.MessageEvent) {
	loadBalancers := rancherListener.GetLoadBalancers()

	var lines []string

	for _, lb := range loadBalancers {
		line := fmt.Sprintf("`%s | %s`", lb.ID, lb.Name)

		lines = append(lines, line)
	}

	msg := "*Lista de Load Balancers:*"

	for _, line := range lines {
		msg += fmt.Sprintf("\n%s", line)
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

func (s *SlackListener) slackUpdateCanary(ev *slack.MessageEvent) {
	args := strings.Split(ev.Msg.Text, " ")

	if len(args) != 5 {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Erro na chamada do comando, sintaxe correta: @nome-do-bot %s id-do-LB peso-nova-versao peso-antiga-versao", canaryUpdate), false))
		return
	}

	lb := args[2]
	newVersionPercent := args[3]
	oldVersionPercent := args[4]

	resp := rancherListener.UpdateCustomHaproxyCfg(lb, newVersionPercent, oldVersionPercent)

	if resp == "error" {
		s.client.PostMessage(ev.Channel, slack.MsgOptionText("Erro ao fazer update no haproxy.cfg, verifique se o ID passado está correto, se o conteúdo do haproxy.cfg atual está em branco ou se os pesos passados não somam 100", false))
		return
	}
	//v := strconv.FormatBool(resp)
	s.client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Arquivo 'haproxy.cfg' alterado com sucesso!\n```%s```", resp), false))
}

func (s *SlackListener) slackLogsContainer(ev *slack.MessageEvent) {
	s.createAndSendAttachment(
		ev,
		"Qual container deseja baixar os logs? :yum:",
		logsContainer,
		getContainers(),
		nil,
	)
}

func (s *SlackListener) slackRestartContainer(ev *slack.MessageEvent) {
	s.createAndSendAttachment(
		ev,
		"Qual container deseja reiniciar? :yum:",
		restartContainer,
		getContainers(),
		nil,
	)
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
