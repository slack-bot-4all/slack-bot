// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package main

import (
	splunk "github.com/drewrm/splunk-golang"
)

// SplunkListener é uma struct que armazena as credenciais e
// URL do Splunk
type SplunkListener struct {
	Username string
	Password string
	BaseURL  string
}

// ConnectSplunk é uma função feita com objetivo de fazer a conexão ao Splunk
func (s *SplunkListener) ConnectSplunk() splunk.SessionKey {
	conn := splunk.SplunkConnection{
		Username: s.Username,
		Password: s.Password,
		BaseURL:  s.BaseURL,
	}

	key, err := conn.Login()
	CheckErr("Não foi possível fazer login no Splunk", err)

	return key
}
