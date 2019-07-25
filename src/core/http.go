package core

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

	// DeleteHTTP é a constante usada para requisições de verbo DELETE
	DeleteHTTP = "DELETE"
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
	case "DELETE":
		req, err = http.NewRequest(method, url, nil)
	default:
		log.Println("[INFO] Not possible create request, method not found.")
	}
	CheckErr("Error to create request", err)

	rancherListener.RancherAuthAdd(req)

	resp, err := client.Do(req)
	CheckErr("Error to send request", err)

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
