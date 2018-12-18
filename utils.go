// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// MakeHTTPPOSTRequest : Função que retorna uma request já com a URL e body prontos e
// o objeto do tipo error, caso haja erro na criação do objeto request.
func (ranchListener *RancherListener) MakeHTTPPOSTRequest(url string, body io.Reader) (r *http.Request, err error) {
	return ranchListener.HandleLogin(http.NewRequest("POST", url, body))
}

// MakeHTTPGETRequest : Função que retorna uma request já com a URL pronta
// junto com o objeto error.
func (ranchListener *RancherListener) MakeHTTPGETRequest(url string) (r *http.Request, err error) {
	return ranchListener.HandleLogin(http.NewRequest("GET", url, nil))
}

// MakeHTTPPUTRequest é uma função para efetuar uma request do tipo PUT
func (ranchListener *RancherListener) MakeHTTPPUTRequest(url string, data io.Reader) (r *http.Request, err error) {
	return ranchListener.HandleLogin(http.NewRequest("PUT", url, data))
}

// HandleLogin : Função que adiciona o usuário e senha à requisição
func (ranchListener *RancherListener) HandleLogin(r *http.Request, errParam error) (rWithCredentials *http.Request, err error) {
	r.PostFormValue(ranchListener.accessKey)
	r.PostFormValue(ranchListener.secretKey)

	return r, nil
}

// CheckErr : Função feita para checar os erros
func CheckErr(message string, err error) {
	if err != nil {
		log.Println(message)
		fmt.Printf("%+v", err)
		panic(err)
	}
}

// ConvertResponseToString : Converte uma série de dados do tipo ReadCloser
// para String
func ConvertResponseToString(response io.ReadCloser) (converted string) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(response)

	return buf.String()

	/*
		Outra solução seria:

		responseBytes, _ := ioutil.ReadAll(response)
		responseString := string(responseBytes)
		return responseString
	*/
}

// WriteOnFile é a função que escreve o que for passado por parâmetro
// dentro do arquivo passado no path
func WriteOnFile(path string, text string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	CheckErr("Erro ao abrir arquivo", err)

	_, err = f.WriteString(text)
	CheckErr("Erro ao escrever linha no arquivo", err)

	defer f.Close()
}
