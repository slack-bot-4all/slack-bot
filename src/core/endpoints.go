// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package core

import (
	"encoding/json"
	"net/http"
)

// Env usado para retornar o objeto env
type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var envs []Env

// GetEnvs mostra todas envs
func GetEnvs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	json.NewEncoder(w).Encode(envs)
}

// GetCommands retorna todos os comandos com todos seus atributos
func GetCommands(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	json.NewEncoder(w).Encode(Commands)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
