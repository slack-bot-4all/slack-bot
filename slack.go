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

	// Parando a função caso a mensagem não traga o prefixo mencionando o BOT
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Tirando a menção ao BOT da mensagem e guardando em uma variável
	message := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1]

	if strings.Contains(ev.Msg.Text, "ajuda") {
		s.SlackCommandHelper(ev, message)
		return nil
	}

	// Fazendo as verificações de mensagens e jogando
	// para as devidas funções
	if strings.HasPrefix(message, restartContainer) {
		s.SlackRestartContainer(ev)
	} else if strings.HasPrefix(message, logsContainer) {
		s.SlackLogsContainer(ev)
	} else if strings.HasPrefix(message, canaryUpdate) {
		s.SlackUpdateCanary(ev)
	} else if strings.HasPrefix(message, haproxyList) {
		s.SlackListLoadBalancers(ev)
	} else if strings.HasPrefix(message, getServiceInfo) {
		s.SlackServiceInfo(ev)
	} else if strings.HasPrefix(message, listService) {
		s.SlackServicesList(ev)
	} else if strings.HasPrefix(message, upgradeService) {
		s.SlackServiceUpgrade(ev)
	} else if strings.HasPrefix(message, canaryDisable) {
		s.SlackCanaryDisable(ev)
	} else if strings.HasPrefix(message, canaryActivate) {
		s.SlackCanaryEnable(ev)
	} else if strings.HasPrefix(message, canaryInfo) {
		s.SlackCanaryInfo(ev)
	} else if strings.HasPrefix(message, comandos) {
		s.SlackHelper(ev)
	}

	return nil
}

// SlackCanaryInfo é a função que é responsável por trazer o
// haproxy.cfg do Load Balancer, com o propósito do usuário
// visualizar como está configurado o Canary
func (s *SlackListener) SlackCanaryInfo(ev *slack.MessageEvent) {
	var attachment = slack.Attachment{
		Text:       "Qual Load Balancer deseja buscar informações do Canary?",
		Color:      "#0C648A",
		CallbackID: canaryInfo,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: getLbOptions(),
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	// Mandando a mensagem pro Slack com o Attachment feito acima
	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}

// SlackCanaryEnable é a função que é responsável por descomentar todas
// as linhas do haproxy.cfg do Load Balancer que for recebido como
// parâmetro
func (s *SlackListener) SlackCanaryEnable(ev *slack.MessageEvent) {
	var attachment = slack.Attachment{
		Text:       "Qual Load Balancer deseja ativar o Canary?",
		Color:      "#0C648A",
		CallbackID: canaryActivate,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: getLbOptions(),
				Confirm: &slack.ConfirmationField{
					Title:       "Deseja mesmo ativar o Canary? :yum:",
					Text:        "Verifique se este é mesmo o Load Balancer que você quer ativar o Canary :smile:",
					OkText:      "Sim",
					DismissText: "Não",
				},
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	// Mandando a mensagem pro Slack com o Attachment feito acima
	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}

// SlackCanaryDisable é a função que é responsável por comentar todas
// as linhas do haproxy.cfg do Load Balancer que for recebido como
// parâmetro
func (s *SlackListener) SlackCanaryDisable(ev *slack.MessageEvent) {
	var attachment = slack.Attachment{
		Text:       "Qual Load Balancer deseja desativar o Canary?",
		Color:      "#0C648A",
		CallbackID: canaryDisable,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: getLbOptions(),
				Confirm: &slack.ConfirmationField{
					Title:       "Deseja mesmo desativar o Canary? :thinking_face:",
					Text:        "Verifique se este é mesmo o Load Balancer que você quer desativar o Canary :smile:",
					OkText:      "Sim",
					DismissText: "Não",
				},
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	// Mandando a mensagem pro Slack com o Attachment feito acima
	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}

// SlackServiceUpgrade é a função responsável por fazer o upgrade da
// imagem de um container que será recebido como parâmetro
func (s *SlackListener) SlackServiceUpgrade(ev *slack.MessageEvent) {
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

// SlackServicesList é a função que enviará uma mensagem no geral listando todos os
// serviços que existem no Environment
func (s *SlackListener) SlackServicesList(ev *slack.MessageEvent) {
	resp := rancherListener.ListServices()

	msg := "*Lista de serviços:* \n\n"

	data := gjson.Get(resp, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		msg += fmt.Sprintf("*ID:* `%s` | *Nome:* `%s`\n", value.Get("id").String(), value.Get("name").String())
		return true
	})

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

// SlackServiceInfo é a função que envia um Attachment para o Slack com
// a lista de serviços que tem no Environment, após isso, o usuário
// selecionará um container, que com isso, o BOT retornará informações sobre
// esse serviço
func (s *SlackListener) SlackServiceInfo(ev *slack.MessageEvent) {
	var attachment = slack.Attachment{
		Text:       "Qual serviço deseja obter informações? :sunglasses:",
		Color:      "#0C648A",
		CallbackID: getServiceInfo,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: getServices(),
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	// Mandando a mensagem pro Slack com o Attachment feito acima
	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}

// SlackCommandHelper é a função que retorna melhores informações
// sobre um comando específico
func (s *SlackListener) SlackCommandHelper(ev *slack.MessageEvent, message string) {
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

// SlackHelper é a função que retorna os comandos possíveis juntamente
// com breves resumos e formas de uso das mesmas
func (s *SlackListener) SlackHelper(ev *slack.MessageEvent) {
	msg := "*Comandos:* "

	for _, cmd := range Commands {
		msg += fmt.Sprintf("`%s` ", cmd.Cmd)
	}

	msg += "\n\n_*Obs.:* Caso queira informações mais detalhadas sobre um comando, você pode chamar este comando seguido de *ajuda*._\n_*Ex.:* @bot comando ajuda_"

	s.client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
}

// SlackListLoadBalancers é a função responsável por retornar para o usuário a lista
// de LB's
func (s *SlackListener) SlackListLoadBalancers(ev *slack.MessageEvent) {
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

// SlackUpdateCanary é a função que busca a função em rancher.go para
// fazer a alteração dos pesos do canary deployment no haproxy.cfg
// dentro do Rancher
func (s *SlackListener) SlackUpdateCanary(ev *slack.MessageEvent) {
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

// SlackSplunk é a função responsável por retornar informações sobre o Splunk
func (s *SlackListener) SlackSplunk(ev *slack.MessageEvent) {
	splunkListener := &SplunkListener{
		Username: SplunkUsername,
		Password: SplunkPassword,
		BaseURL:  SplunkBaseURL,
	}

	key := splunkListener.ConnectSplunk()

	fmt.Println(key)
}

// SlackLogsContainer é a função responsável por mandar o attachment
// com todos os containers, para o usuário selecionar um para recuperar
// os logs
func (s *SlackListener) SlackLogsContainer(ev *slack.MessageEvent) {
	var attachment = slack.Attachment{
		Text:       "Qual container deseja baixar os logs? :yum:",
		Color:      "#0C648A",
		CallbackID: logsContainer,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: getContainers(),
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}

// SlackRestartContainer Função responsável por fazer o reinício de um container dentro do Rancher
func (s *SlackListener) SlackRestartContainer(ev *slack.MessageEvent) {
	// Criando attachment e setando options como a lista de opcoes que foi
	// iterada acima
	var attachment = slack.Attachment{
		Text:       "Qual container deseja reiniciar? :yum:",
		Color:      "#0C648A",
		CallbackID: restartContainer,
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: getContainers(),
				Confirm: &slack.ConfirmationField{
					Title:       "Deseja reiniciar este container?",
					Text:        "Verifique se realmente é este container que você deseja reiniciar",
					OkText:      "Sim",
					DismissText: "Não",
				},
			},
			{
				Name:  "cancel",
				Text:  "Cancelar",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	// Mandando a mensagem pro Slack com o Attachment feito acima
	s.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
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
