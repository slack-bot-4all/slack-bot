// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rgamba/evtwebsocket"
	"github.com/tidwall/gjson"
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
	client := &http.Client{}

	req, err := ranchListener.MakeHTTPPOSTRequest(fmt.Sprintf(ranchListener.baseURL+"/"+ranchListener.projectID+"/containers/"+containerID+"?action=restart"), nil)
	CheckErr("Erro ao montar requisição", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição", err)
	defer resp.Body.Close()

	responseString := ConvertResponseToString(resp.Body)
	idValue := gjson.Get(responseString, "id")

	log.Println(fmt.Sprintf("Container restartado! ID: %+v", idValue))
}

// ListContainers é uma função que retornará uma lista de todos os containers de um projeto/environment
func (ranchListener *RancherListener) ListContainers() string {
	client := &http.Client{}

	req, err := ranchListener.MakeHTTPGETRequest(fmt.Sprintf(ranchListener.baseURL + "/" + ranchListener.projectID + "/containers/"))
	CheckErr("Erro ao montar requisição", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição", err)
	defer resp.Body.Close()

	responseString := ConvertResponseToString(resp.Body)

	return responseString
}

// LogsContainer : Função responsável retornar os logs do container
func (ranchListener *RancherListener) LogsContainer(containerID string) string {
	client := &http.Client{}

	var jsonStr = []byte(`{"follow": true, "lines": 50}`)

	req, err := ranchListener.MakeHTTPPOSTRequest(fmt.Sprintf(ranchListener.baseURL+"/"+ranchListener.projectID+"/containers/"+containerID+"?action=logs"), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	CheckErr("Erro ao montar requisição", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição", err)
	defer resp.Body.Close()

	responseString := ConvertResponseToString(resp.Body)
	tokenValue := gjson.Get(responseString, "token")
	urlValue := gjson.Get(responseString, "url")

	urlAndToken := fmt.Sprintf("%s?token=%s", urlValue.String(), tokenValue.String())

	t := time.Now()

	f, err := os.Create(fmt.Sprintf("/tmp/logs-container-%d%d%d%02d%02d%02d.log", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second()))

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
func (ranchListener *RancherListener) UpdateCustomHaproxyCfg(ID string, newPercent string, oldPercent string) bool {
	newPercentToInteger, _ := strconv.Atoi(newPercent)
	oldPercentToInteger, _ := strconv.Atoi(oldPercent)

	if (newPercentToInteger + oldPercentToInteger) != 100 {
		return false
	}

	client := &http.Client{}

	actualLbConfig := ranchListener.GetHaproxyCfg(ID)

	if actualLbConfig == "" {
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(actualLbConfig))

	var firstWeight string
	var secondWeight string

	for scanner.Scan() {
		if line := strings.Split(scanner.Text(), "weight "); len(line) >= 2 {
			if firstWeight == "" {
				firstWeight = line[1]
			}
			if secondWeight == "" {
				secondWeight = line[1]
			}
		}
	}

	newLbConfig := strings.Replace(actualLbConfig, fmt.Sprintf("weight %s", firstWeight), fmt.Sprintf("weight %s", newPercent), -1)
	newLbConfig = strings.Replace(newLbConfig, fmt.Sprintf("weight %s", secondWeight), fmt.Sprintf("weight %s", oldPercent), -1)

	lbConfig := &LoadBalancerServices{
		LbConfig: &LbConfig{
			Config: newLbConfig,
		},
	}

	payload, err := json.Marshal(lbConfig)
	CheckErr("Erro ao fazer conversão de struct para JSON", err)

	payloadReader := bytes.NewReader(payload)

	req, err := ranchListener.MakeHTTPPUTRequest(fmt.Sprintf(ranchListener.baseURL+"/"+ranchListener.projectID+"/loadBalancerServices/"+ID), payloadReader)
	CheckErr("Erro ao montar requisição", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição", err)
	defer resp.Body.Close()

	return true
}

// GetHaproxyCfg Busca a Custom haproxy.cfg do LoadBalancer enviado como parâmetro
func (ranchListener *RancherListener) GetHaproxyCfg(containerID string) string {
	client := &http.Client{}

	req, err := ranchListener.MakeHTTPGETRequest(fmt.Sprintf(ranchListener.baseURL + "/" + ranchListener.projectID + "/loadBalancerServices/" + containerID))
	CheckErr("Erro ao montar requisição de haproxy.cfg", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição", err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	responseString := ConvertResponseToString(resp.Body)
	lbConfig := gjson.Get(responseString, "lbConfig.config").String()

	return lbConfig
}

// GetLoadBalancers é a função responsável por trazer um slice
// de LoadBalancer, que pode ser usado para selects na interface
// do BOT do Slack
func (ranchListener *RancherListener) GetLoadBalancers() []*LoadBalancer {
	client := &http.Client{}

	req, err := ranchListener.MakeHTTPGETRequest(fmt.Sprintf(ranchListener.baseURL + "/" + ranchListener.projectID + "/loadBalancerServices"))
	CheckErr("Erro ao montar requisição em GetLoadBalancers()", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição em GetLoadBalancers()", err)
	defer resp.Body.Close()

	respString := ConvertResponseToString(resp.Body)

	loadBalancersSlice := []*LoadBalancer{}

	data := gjson.Get(respString, "data")
	data.ForEach(func(key, value gjson.Result) bool {
		lb := new(LoadBalancer)
		lb.ID = value.Get("id").String()
		lb.Name = value.Get("name").String()
		loadBalancersSlice = append(loadBalancersSlice, lb)

		return true
	})

	return loadBalancersSlice
}
