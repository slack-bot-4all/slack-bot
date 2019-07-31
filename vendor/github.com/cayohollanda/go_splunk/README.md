# Go Splunk 
[![GoDoc](https://godoc.org/github.com/cayohollanda/go_splunk?status.svg)](https://godoc.org/github.com/cayohollanda/go_splunk) [![Go Report Card](https://goreportcard.com/badge/github.com/cayohollanda/go_splunk)](https://goreportcard.com/report/github.com/cayohollanda/go_splunk)

A library that provides functions to programmer to connect and get responses of Splunk API with more abstract form and no stress

# Installation
To install, you need only get the package
```
go get github.com/cayohollanda/go_splunk
```

# Usage
To call all methods to connect on Splunk, you need to have a declared variable with credentials of Splunk
```go
package main

import "github.com/cayohollanda/go_splunk"

func main() {
  conn := &go_splunk.SplunkConnection{
    APIURL:   "https://localhost:8089",
    Username: "splunk-username",
    Password: "splunk-password",
  }
}
```

# Get search results
A example of usage of ```GetSearchResults()```:
```go
package main

import "github.com/cayohollanda/go_splunk"

func main() {
  conn := &go_splunk.SplunkConnection{
    APIURL:   "https://localhost:8089",
    Username: "splunk-username",
    Password: "splunk-password",
  }
  
  results, err := conn.GetSearchResults("index=test")
  if err != nil {
    log.Fatalf("Error: %s", err.Error())
  }
  
  // Use the 'results' variable to execute what do you need
}
```
