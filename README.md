Equinix oAuth2 Go client
================
Go implementation of oAuth2 enabbled HTTP client for interactions with Equinix APIs.
Module implementes Equinix specific client credentials grant type with custom `TokenSource` from standard Go oauth2 module.

[![Build Status](https://travis-ci.com/equinix/oauth2-go.svg?branch=master)](https://travis-ci.com/github/equinix/oauth2-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/oauth2-go)](https://goreportcard.com/report/github.com/equinix/oauth2-go)
[![GoDoc](https://godoc.org/github.com/equinix/oauth2-go?status.svg)](https://godoc.org/github.com/equinix/oauth2-go)
![GitHub](https://img.shields.io/github/license/equinix/oauth2-go)

---

Requirements
----------------
* [Go](https://golang.org/doc/install) 1.14+ (to build provider plugin)

Usage
----------------
1. Import
    ```
    import "github.com/equinix/oauth2-go"
    ```

2. Prepare configuration and get http client 
    ```
	authConfig := oauth2.Config{
		ClientID:     "myClientId",
		ClientSecret: "myClientSecret"
		BaseURL:      "https://api.equinix.com"}

    //*http.Client is returned
	hc := authConfig.New(context.Background())
    ```

3. Use client

    `*http.Client` created by oAuth2 library will deal with token acquisition, refreshment and population of Authorization headers in subsequent requests. 

    Below example shows how to use oAuth2 client with [Resty REST client library](https://github.com/go-resty/resty)
    ```
	rc := resty.NewWithClient(hc)
    resp, err := rc.R().Get("https://api.equinix.com/ecx/v3/port/userport")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Status Code:", resp.StatusCode())
        fmt.Println("Body:\n", resp)
    }
    ```