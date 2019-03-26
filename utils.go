// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

type Kanye struct{
	Quote string `json:"quote"`
} 

// CheckErr : Função feita para checar os erros
func CheckErr(message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s\n%s", message, err)
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

// RemoveLastCharacter é a função que remove o último caracter de uma string
func RemoveLastCharacter(s string) string {
	sz := len(s)

	if sz > 0 && s[sz-1] == '.' {
		s = s[:sz-1]
	}

	return s
}
