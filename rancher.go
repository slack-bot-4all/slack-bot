// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rgamba/evtwebsocket"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// RancherListener é uma estrutura onde ficam armazenados os dados de acesso ao Rancher API
type RancherListener struct {
	accessKey string
	secretKey string
	baseURL   string
	projectID string
}

// Container é uma estrutura que é usada para mostrar informações ao usuário
type Container struct {
	id        string
	imageUUID string
	name      string
}

// LbConfig é uma estrutura que é usada para enviar dados para conf Haproxy
type LbConfig struct {
	Config string `json:"config"`
}

// LoadBalancerServices é uma estrutura que é usada para a construção do JSON de requisição
// quando se vai fazer o edit/upgrade de LB's
type LoadBalancerServices struct {
	LbConfig *LbConfig `json:"lbConfig"`
}

// LoadBalancer é a estrutura que tem como objetivo representar um LoadBalancer do Rancher
// de forma um pouco mais resumida (é usado na função GetLoadBalancers())
type LoadBalancer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RestartContainer : Função responsável por dar restart no container recebido por parâmetro
func (ranchListener *RancherListener) RestartContainer(containerID string) {
	url := fmt.Sprintf("%s/%s/containers/%s?action=restart", ranchListener.baseURL, ranchListener.projectID, containerID)
	resp := ranchListener.HTTPSendRancherRequest(url, PostHTTP, "")

	idValue := gjson.Get(resp, "id").String()

	log.Println("Container restartado! ID:", idValue)
}

// ListContainers é uma função que retornará uma lista de todos os containers de um projeto/environment
func (ranchListener *RancherListener) ListContainers() string {
	url := fmt.Sprintf("%s/%s/containers", ranchListener.baseURL, ranchListener.projectID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	return resp
}

// GetService é uma função que retorna o JSON de uma requisição que busca
// informações de um único serviço
func (ranchListener *RancherListener) GetService(ID string) string {
	url := fmt.Sprintf("%s/%s/services/%s", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	return resp
}

// UpgradeService é a função que faz o upgrade da imagem do serviço, recebendo
// como parâmetro o ID do serviço e o nome da nova imagem do serviço
func (ranchListener *RancherListener) UpgradeService(ID string, newImage string) string {
	var data interface{}
	var jsonRequest string

	originalServiceCfg := ranchListener.GetService(ID)

	originalServiceCfg, err := sjson.Set(originalServiceCfg, "launchConfig.imageUuid", newImage)
	CheckErr("Erro ao setar valor de nova variável no JSON do serviço", err)

	launchConfigOri := gjson.Get(originalServiceCfg, "launchConfig").String()

	json.Unmarshal([]byte(launchConfigOri), &data)

	jsonRequest, err = sjson.Set(jsonRequest, "inServiceStrategy.launchConfig", data)
	CheckErr("Erro ao setar valor de nova variável no JSON do serviço", err)

	url := fmt.Sprintf("%s/%s/services/%s?action=upgrade", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, PostHTTP, jsonRequest)

	return gjson.Get(resp, "launchConfig.imageUuid").String()
}

// ListServices é uma função que retorna o JSON (em string) de uma requisição que tem como
// objetivo buscar todos os serviços do Environment
func (ranchListener *RancherListener) ListServices() string {
	url := fmt.Sprintf("%s/%s/services", ranchListener.baseURL, ranchListener.projectID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	return resp
}

// LogsContainer : Função responsável retornar os logs do container
func (ranchListener *RancherListener) LogsContainer(containerID string) string {
	data := &url.Values{}
	data.Add("follow", "true")
	data.Add("lines", "50")

	url := fmt.Sprintf("%s/%s/containers/%s?action=logs", ranchListener.baseURL, ranchListener.projectID, containerID)

	resp := ranchListener.HTTPSendRancherRequest(url, PostHTTP, `{"follow": true, "lines": 50}`)

	tokenValue := gjson.Get(resp, "token").String()
	urlValue := gjson.Get(resp, "url").String()

	urlAndToken := fmt.Sprintf("%s?token=%s", urlValue, tokenValue)

	t := time.Now()

	f, err := os.Create(fmt.Sprintf("/tmp/logs-container-%d%d%d%02d%02d%02d.log", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second()))
	CheckErr("Erro na criação do arquivo de logs", err)

	SocketConnectionLogsContainer(urlAndToken, f.Name())

	return f.Name()
}

// SocketConnectionLogsContainer é uma função que retornará uma conexão com o ws da URL passada como parâmetro
func SocketConnectionLogsContainer(urlAndToken string, fileName string) {
	conn := &evtwebsocket.Conn{
		OnConnected: func(w *evtwebsocket.Conn) {
			log.Println("[INFO] Conectado no WebSocket!")
		},

		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			verifyExists, err := os.Open(fileName)

			if err == nil {
				WriteOnFile(fileName, string(msg))
			}

			defer verifyExists.Close()
		},

		OnError: func(err error) {
			log.Printf("[ERROR] %s\n", err.Error())
			os.Exit(1)
		},
	}

	conn.Dial(urlAndToken, "")
}

// UpdateCustomHaproxyCfg Edita o lbConfig.config do LB
func (ranchListener *RancherListener) UpdateCustomHaproxyCfg(ID string, newPercent string, oldPercent string) string {
	newPercentToInteger, _ := strconv.Atoi(newPercent)
	oldPercentToInteger, _ := strconv.Atoi(oldPercent)

	if (newPercentToInteger + oldPercentToInteger) != 100 {
		return "error"
	}

	responseString := ranchListener.GetHaproxyCfg(ID)
	actualLbConfig := gjson.Get(responseString, "lbConfig.config").String()

	if actualLbConfig == "" {
		return "error"
	}

	scanner := bufio.NewScanner(strings.NewReader(actualLbConfig))

	var firstWeight string
	var secondWeight string
	var newLbConfig string

	for scanner.Scan() {
		if line := strings.Split(scanner.Text(), "weight "); len(line) >= 2 {
			if firstWeight == "" {
				firstWeight = line[1]
				newLbConfig = strings.Replace(actualLbConfig, fmt.Sprintf("weight %s", firstWeight), fmt.Sprintf("weight %s", newPercent), 1)
			} else {
				secondWeight = line[1]
				newLbConfig = strings.Replace(newLbConfig, fmt.Sprintf("weight %s", secondWeight), fmt.Sprintf("weight %s", oldPercent), 1)
			}
		}
	}

	responseString, err := sjson.Set(responseString, "lbConfig.config", newLbConfig)
	CheckErr("Erro ao setar novo Custom haproxy.cfg no JSON", err)

	url := fmt.Sprintf("%s/%s/loadBalancerServices/%s", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, PutHTTP, responseString)

	return gjson.Get(resp, "lbConfig.config").String()
}

// GetHaproxyCfg Busca a Custom haproxy.cfg do LoadBalancer enviado como parâmetro
func (ranchListener *RancherListener) GetHaproxyCfg(containerID string) string {
	url := fmt.Sprintf(ranchListener.baseURL + "/" + ranchListener.projectID + "/loadBalancerServices/" + containerID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	if gjson.Get(resp, "id").String() != containerID {
		return ""
	}

	return resp
}

// GetLoadBalancers é a função responsável por trazer um slice
// de LoadBalancer, que pode ser usado para selects na interface
// do BOT do Slack
func (ranchListener *RancherListener) GetLoadBalancers() []*LoadBalancer {
	url := fmt.Sprintf("%s/%s/loadBalancerServices", ranchListener.baseURL, ranchListener.projectID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	loadBalancersSlice := []*LoadBalancer{}

	data := gjson.Get(resp, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		lb := new(LoadBalancer)
		lb.ID = value.Get("id").String()
		lb.Name = value.Get("name").String()
		loadBalancersSlice = append(loadBalancersSlice, lb)

		return true
	})

	return loadBalancersSlice
}
