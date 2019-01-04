// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
