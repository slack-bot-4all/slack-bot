// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"github.com/tidwall/gjson"
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
			s.client.PostMessage(s.channelID, slack.MsgOptionText("Estou online! :nerd_face:", false))
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

	// Fazendo as verificações de mensagens e jogando
	// para as devidas funções
	switch message {
	case "restart-container":
		s.SlackRestartContainer(ev)
	case "logs-container":
		s.SlackLogsContainer(ev)
	case "haproxy-create":
		s.SlackUpdateHaproxy(ev)
	case "splunk":
		s.SlackSplunk(ev)
	}

	return nil
}

// SlackUpdateHaproxy é a função que busca a função em rancher.go para
// fazer a alteração dos pesos do canary deployment no haproxy.cfg
// dentro do Rancher
func (s *SlackListener) SlackUpdateHaproxy(ev *slack.MessageEvent) {
	// rancherListener.UpdateCustomHaproxyCfg("1s30", "40", "60")
	var attachment = slack.Attachment{
		Text:       "Selecione o Load Balancer, percentual da nova versão e da antiga versão, respectivamente. :grinning:",
		Color:      "#0C648A",
		CallbackID: "update-haproxy",
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Text:    "Load Balancer",
				Options: getLbOptions(),
			},
			{
				Name:    "select",
				Type:    "select",
				Text:    "Porc. Nova versão",
				Options: percentOptions(),
			},
			{
				Name:    "select",
				Type:    "select",
				Text:    "Porc. Antiga versão",
				Options: percentOptions(),
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
		CallbackID: "logs-container",
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

// Função responsável por fazer o reinício de um container dentro do Rancher
func (s *SlackListener) SlackRestartContainer(ev *slack.MessageEvent) {
	// Criando attachment e setando options como a lista de opcoes que foi
	// iterada acima
	var attachment = slack.Attachment{
		Text:       "Qual container deseja reiniciar? :yum:",
		Color:      "#0C648A",
		CallbackID: "restart-container",
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

func percentOptions() []slack.AttachmentActionOption {
	opcoes := []slack.AttachmentActionOption{}
	for i := 0; i < 100; i++ {
		opcoes = append(opcoes, slack.AttachmentActionOption{
			Text:  fmt.Sprintf("%d%%", i),
			Value: strconv.Itoa(i),
		})
	}

	return opcoes
}
