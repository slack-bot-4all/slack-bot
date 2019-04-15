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
	"regexp"
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
	ID        string `json:"id"`
	ImageUUID string `json:"imageUuid"`
	Name      string `json:"name"`
	State     string `json:"state"`
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

	log.Println("[INFO] Container restarted! ID:", idValue)
}

// StartService : Função responsável por dar start no container recebido por parâmetro
func (ranchListener *RancherListener) StartService(ID string) {
	url := fmt.Sprintf("%s/%s/services/%s?action=activate", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, PostHTTP, "")

	idValue := gjson.Get(resp, "id").String()

	log.Println("[INFO] Service started! ID:", idValue)
}

// StopService : Função responsável por dar stop no container recebido por parâmetro
func (ranchListener *RancherListener) StopService(ID string) {
	url := fmt.Sprintf("%s/%s/services/%s?action=deactivate", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, PostHTTP, "")

	idValue := gjson.Get(resp, "id").String()

	log.Println("[INFO] Service stopped! ID:", idValue)

}

// ListContainers é uma função que retornará uma lista de todos os containers de um projeto/environment
func (ranchListener *RancherListener) ListContainers() string {
	url := fmt.Sprintf("%s/%s/containers", ranchListener.baseURL, ranchListener.projectID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	return resp
}

// GetInstances é uma função que retornará uma lista de todas as instâncias de um serviço
func (ranchListener *RancherListener) GetInstances(serviceID string) string {
	url := fmt.Sprintf("%s/%s/services/%s/instances", ranchListener.baseURL, ranchListener.projectID, serviceID)
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

// GetServiceStack é uma função que retorna o JSON de uma requisição que busca
// informações da stack de um serviço em específico
func (ranchListener *RancherListener) GetServiceStack(ID string) string {
	url := fmt.Sprintf("%s/%s/services/%s/stack", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	return resp
}

// GetStacks é uma função que retorna o JSON de uma requisição que busca
// todas as stacks do environment
func (ranchListener *RancherListener) GetStacks() string {
	url := fmt.Sprintf("%s/%s/stacks", ranchListener.baseURL, ranchListener.projectID)
	resp := ranchListener.HTTPSendRancherRequest(url, GetHTTP, "")

	return resp
}

// GetServicesFromStack é uma função que retorna o JSON de uma requisição que busca
// todos os serviços de uma stack especificada
func (ranchListener *RancherListener) GetServicesFromStack(ID string) string {
	url := fmt.Sprintf("%s/%s/stacks/%s/services", ranchListener.baseURL, ranchListener.projectID, ID)
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
	CheckErr("Error to set value of new variable of service JSON", err)

	launchConfigOri := gjson.Get(originalServiceCfg, "launchConfig").String()

	json.Unmarshal([]byte(launchConfigOri), &data)

	jsonRequest, err = sjson.Set(jsonRequest, "inServiceStrategy.launchConfig", data)
	CheckErr("Error to set value of new variable of service JSON", err)

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
	CheckErr("Logs file creation error", err)

	socketConnectionLogsContainer(urlAndToken, f.Name())

	return f.Name()
}

func socketConnectionLogsContainer(urlAndToken string, fileName string) {
	conn := &evtwebsocket.Conn{
		OnConnected: func(w *evtwebsocket.Conn) {
			log.Println("[INFO] Connectec on WebSocket!")
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

// DisableCanary é a função que envia a requisição para a API do
// Rancher com a intenção de comentar todas as linhas do haproxy.cfg
func (ranchListener *RancherListener) DisableCanary(ID string) string {
	responseString := ranchListener.GetHaproxyCfg(ID)
	actualLbConfig := gjson.Get(responseString, "lbConfig.config").String()

	if actualLbConfig == "" {
		return "error"
	}

	scanner := bufio.NewScanner(strings.NewReader(actualLbConfig))

	newLbConfig := ""

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "#") {
			newLbConfig += fmt.Sprintf("%s\n", line)
		} else {
			newLbConfig += fmt.Sprintf("#%s\n", line)
		}
	}

	responseString, err := sjson.Set(responseString, "lbConfig.config", newLbConfig)
	CheckErr("Error to set new Custom haproxy.cfg on JSON", err)

	url := fmt.Sprintf("%s/%s/loadBalancerServices/%s", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, PutHTTP, responseString)

	return gjson.Get(resp, "lbConfig.config").String()
}

// EnableCanary é a função que retira os "#" de todo o haproxy.cfg
// depois envia como PUT para a API do Rancher
func (ranchListener *RancherListener) EnableCanary(ID string) string {
	responseString := ranchListener.GetHaproxyCfg(ID)
	actualLbConfig := gjson.Get(responseString, "lbConfig.config").String()

	if actualLbConfig == "" {
		return "error"
	}

	actualLbConfig = strings.Replace(actualLbConfig, "#", "", -1)

	responseString, err := sjson.Set(responseString, "lbConfig.config", actualLbConfig)
	CheckErr("Error to set new Custom haproxy.cfg on JSON", err)

	url := fmt.Sprintf("%s/%s/loadBalancerServices/%s", ranchListener.baseURL, ranchListener.projectID, ID)
	resp := ranchListener.HTTPSendRancherRequest(url, PutHTTP, responseString)

	return gjson.Get(resp, "lbConfig.config").String()
}

// UpdateCustomHaproxyCfg Edita o lbConfig.config do LB
func (ranchListener *RancherListener) UpdateCustomHaproxyCfg(ID string, newPercent string, oldPercent string) string {
	newPercentToInteger, _ := strconv.Atoi(newPercent)
	oldPercentToInteger, _ := strconv.Atoi(oldPercent)

	if (newPercentToInteger + oldPercentToInteger) != 100 {
		return "error"
	}

	ranchListener.EnableCanary(ID)
	responseString := ranchListener.GetHaproxyCfg(ID)
	actualLbConfig := gjson.Get(responseString, "lbConfig.config").String()

	if actualLbConfig == "" {
		return "error"
	}

	scanner := bufio.NewScanner(strings.NewReader(actualLbConfig))

	var newLbConfig string
	var lines []string

	for scanner.Scan() {
		lines = strings.Split(scanner.Text(), "\n")
		if strings.HasPrefix(lines[0], "server") {
			var n string
			var o string
			new := regexp.MustCompile(".+new.+(\\d{2})")
			old := regexp.MustCompile(".+old.+(\\d{2})")
			l := new.FindStringSubmatch(scanner.Text())
			z := old.FindStringSubmatch(scanner.Text())
			if len(l) >= 2 {
				n = strings.Replace(scanner.Text(), fmt.Sprintf("weight %s", l[1]), fmt.Sprintf("weight %s", newPercent), 1)
				if newLbConfig == "" {
					newLbConfig = strings.Replace(actualLbConfig, scanner.Text(), n, 1)
				} else {
					newLbConfig = strings.Replace(newLbConfig, scanner.Text(), n, 1)
				}
			}

			if len(z) >= 2 {
				o = strings.Replace(scanner.Text(), fmt.Sprintf("weight %s", z[1]), fmt.Sprintf("weight %s", oldPercent), 1)
				if newLbConfig == "" {
					newLbConfig = strings.Replace(actualLbConfig, scanner.Text(), o, 1)
				} else {
					newLbConfig = strings.Replace(newLbConfig, scanner.Text(), o, 1)
				}
			}

		}

	}

	fmt.Println("DEBUG [newLbConfig] =>", newLbConfig)

	responseString, err := sjson.Set(responseString, "lbConfig.config", newLbConfig)
	CheckErr("Error to set new Custom haproxy.cfg on JSON", err)

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
