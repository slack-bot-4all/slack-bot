package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

const (
	// GetHTTP é a constante usada para requisições de verbo GET
	GetHTTP = "GET"

	// PostHTTP é a constante usada para requisições de verbo POST
	PostHTTP = "POST"

	// PutHTTP é a constante usada para requisições de verbo PUT
	PutHTTP = "PUT"
)

// CreateHTTPClient é a função responsável por retornar um client para que possam ser
// enviadas as requisições
func CreateHTTPClient() *http.Client {
	transp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transp}

	return client
}

// HTTPSendRancherRequest é a função que envia a requisição para a
// API do Rancher e retorna o body do response já convertido em
// String
func (rancherListener *RancherListener) HTTPSendRancherRequest(url string, method string, data string) string {
	client := CreateHTTPClient()

	var payload io.Reader
	if data != "" {
		payload = bytes.NewBufferString(data)
	}

	var req *http.Request
	var err error
	switch method {
	case "GET":
		req, err = http.NewRequest(method, url, nil)
	case "POST":
		req, err = http.NewRequest(method, url, payload)
	case "PUT":
		req, err = http.NewRequest(method, url, payload)
	default:
		log.Println("[INFO] Não possível criar requisição, método não encontrado.")
	}
	CheckErr("[ERROR] Erro ao criar requisição", err)

	rancherListener.RancherAuthAdd(req)

	resp, err := client.Do(req)
	CheckErr("[ERROR] Erro ao enviar requisição", err)

	return ConvertResponseToString(resp.Body)
}

// RancherAuthAdd é a função que adiciona as credenciais na requisição que será feita
// para a API do Rancher
func (rancherListener *RancherListener) RancherAuthAdd(request *http.Request) {
	if rancherListener.accessKey != "" && rancherListener.secretKey != "" {
		request.SetBasicAuth(rancherListener.accessKey, rancherListener.secretKey)
	}
}

func createHTTPClient() *http.Client {
    transp := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{Transport: transp}

    return client
}