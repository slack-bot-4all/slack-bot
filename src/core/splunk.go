// Slack BOT for Rancher API
// Created by: https://github.com/magnonta and https://github.com/cayohollanda

package core

import (
	"time"
	"log"
	"encoding/json"

	splunk "github.com/cayohollanda/go_splunk"
)

// SplunkListener é uma struct que armazena as credenciais e
// URL do Splunk
type SplunkListener struct {
	Username string
	Password string
	APIURL   string
}

// ResultSearch ::
type ResultSearch struct {
	Timestamp  string `json:"@timestamp"`
	Version    int    `json:"@version"`
	Message    string `json:"message"`
	LoggerName string `json:"logger_name"`
	ThreadName string `json:"thread_name"`
	Level      string `json:"level"`
	LevelValue int    `json:"level_value"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	AppProfile string `json:"app_profile"`
	Trace      struct {
		UUIDRequest         string    `json:"uuidRequest"`
		InsertedOnDate      time.Time `json:"insertedOnDate"`
		DurationMillis      int       `json:"durationMillis"`
		Host                string    `json:"host"`
		Verb                string    `json:"verb"`
		URL                 string    `json:"url"`
		Pattern             string    `json:"pattern"`
		Method              string    `json:"method"`
		ResultStatus        int       `json:"resultStatus"`
		App                 string    `json:"app"`
		IDEmissor           int       `json:"idEmissor"`
		Base                string    `json:"base"`
		ReceivedFromAddress string    `json:"receivedFromAddress"`
		StackTrace          struct {
			Clazz   string `json:"clazz"`
			Message string `json:"message"`
			Stack   string `json:"stack"`
		} `json:"stackTrace"`
		Logs []struct {
			Ts      time.Time `json:"ts"`
			Level   string    `json:"level"`
			Logger  string    `json:"logger"`
			Thread  string    `json:"thread"`
			Content string    `json:"content"`
		} `json:"logs"`
	} `json:"trace"`
}

// ConnectSplunk é uma função feita com objetivo fazer uma consulta no Splunk
func (s *SplunkListener) ConnectSplunk(query string) ResultSearch {
	conn := &splunk.SplunkConnection{
		Username: s.Username,
		Password: s.Password,
		APIURL:   s.APIURL,
	}

	searchResults, err := conn.GetSearchResults(query)
	if err != nil {
		log.Println("erro")
	}

	var rs ResultSearch
	_ = json.Unmarshal([]byte(searchResults[0].Result.Raw), &rs)

	return rs
}
