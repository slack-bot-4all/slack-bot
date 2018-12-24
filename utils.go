// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

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
