// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
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

// ConfHaproxy é uma estrutura que é usada para enviar dados para conf Haproxy
type ConfHaproxy struct {
	cert                string //Optional 	- 	-
	certChain           string //Optional 	- 	-
	clusterSize         int    //Optional 	- 	3
	hostRegistrationURL string //Yes 	- 	-
	httpEnabled         bool   //Optional 	- 	true
	httpPort            int    //Optional 	- 	80
	httpsPort           int    //Optional 	- 	443
	key                 string //Optional 	- 	-
	ppHTTPPort          int    //Optional 	- 	81
	ppHTTPSPort         int    //Optional 	- 	444
	redisPort           int    //Optional 	- 	6379
	swarmEnabled        bool   //Optional 	- 	true
	swarmPort           int    //Optional 	- 	2376
	zookeeperClientPort int    //Optional 	- 	2181
	zookeeperLeaderPort int    //Optional 	- 	3888
	zookeeperQuorumPort int    //Optional 	- 	2888
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
				fmt.Println(string(msg))
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

// CreateConfigHaproxy Cria conf no lb
func (ranchListener *RancherListener) CreateConfigHaproxy(ID string) {
	//Código genérico para POST de criação das conf
	//Checar melhor a doc da API
	//https://rancher.com/docs/rancher/v1.6/en/api/v2-beta/api-resources/haConfig/#createscript

	client := &http.Client{}

	json := bytes.NewBuffer([]byte(`{
		"cert": "string",
		"certChain": "string",
		"clusterSize": 3,
		"hostRegistrationUrl": "string",
		"httpEnabled": true,
		"httpPort": 80,
		"httpsPort": 443,
		"key": "string",
		"ppHttpPort": 81,
		"ppHttpsPort": 444,
		"redisPort": 6379,
		"swarmEnabled": true,
		"swarmPort": 2376,
		"zookeeperClientPort": 2181,
		"zookeeperLeaderPort": 3888,
		"zookeeperQuorumPort": 2888
	}`))
	req, err := ranchListener.MakeHTTPPOSTRequest(fmt.Sprintf(ranchListener.baseURL+"/"+ranchListener.projectID+"/haConfigs/"+ID+"?action=createscript"), nil)
	CheckErr("Erro ao montar requisição", err)

	resp, err := client.Do(req)
	CheckErr("Erro ao enviar requisição", err)
	defer resp.Body.Close()

	log.Println(fmt.Sprintf("Configuração criada! %+v", resp))

}

// UpdateConfHaproxy Atualiza configurações existentes nas conf do lb
func (ranchListener *RancherListener) UpdateConfHaproxy(ID string) {
	// Criar request do tipo PUT
	var (
		y = "enabled"
		n = "disable"
	)
	body := fmt.Sprintf("{\"true\":%q, \"false\":%q}", y, n)
	ranchListener.MakeHTTPPUTRequest(fmt.Sprintf(ranchListener.baseURL+"/"+ranchListener.projectID+"/haConfigs/"+ID+"?action=createscript"), strings.NewReader(body))

}
