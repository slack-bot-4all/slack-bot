package go_splunk

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// SplunkConnection : guard data of connection with Splunk API
type SplunkConnection struct {

	// Username to connect on Splunk (usually is your login username)
	Username string

	// Password to connect on Splunk (usually is your logi password of your username)
	Password string

	// URL to connect with Splunk API. Ex.: https://localhost:8089/
	APIURL string
}

// SearchResult : have result of a search that has been sended
// to Splunk, having EventsCounter that is a counter to how many
// events has returned on request and Raw that is a []string
// of returned raws of events
type SearchResult struct {
	Preview bool  `json:"preview"`
	Offset  int64 `json:"offset"`
	Result  struct {
		Bkt          string   `json:"_bkt"`
		Cd           string   `json:"_cd"`
		IndexTime    string   `json:"_indextime"`
		Raw          string   `json:"_raw"`
		Serial       string   `json:"_serial"`
		Si           []string `json:"_si"`
		Sourcetype   string   `json:"_sourcetype"`
		Time         string   `json:"_time"`
		Host         string   `json:"host"`
		Index        string   `json:"index"`
		LineCount    string   `json:"linecount"`
		Source       string   `json:"source"`
		SourceType   string   `json:"sourcetype"`
		SplunkServer string   `json:"splunkserver"`
	} `json:"result"`
}

// GetSearchResults : is a function to return a SearchResult of one search
// sended to Splunk API. This function receives search param that is a string of
// search
func (conn SplunkConnection) GetSearchResults(search string) (results []SearchResult, err error) {
	var searchResults []SearchResult

	response, err := conn.HTTPGetRequest(fmt.Sprintf("%s/services/search/jobs/export?search=search%%20%s&output_mode=json", conn.APIURL, search), nil)
	if err != nil {
		return []SearchResult{}, err
	}

	responseToByteSlice, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []SearchResult{}, err
	}

	responseReader := bytes.NewReader(responseToByteSlice)
	responseScanner := bufio.NewScanner(responseReader)

	for responseScanner.Scan() {
		var resultOfSearch SearchResult
		err = json.Unmarshal(responseScanner.Bytes(), &resultOfSearch)

		searchResults = append(searchResults, resultOfSearch)
	}

	if err != nil {
		return []SearchResult{}, err
	}

	return searchResults, nil
}
