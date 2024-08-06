package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

var host string = "http://localhost"

func getRequestTest(t *testing.T, path string, responseStatus int, responseBody string, exactMatch bool) {
	resp, err := http.Get(host + ":" + port + path)
	if err != nil {
		t.Fatal(err.Error())
	} else if resp.StatusCode != responseStatus {
		t.Errorf("Unexpected status: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	} else if exactMatch && (string(body) != responseBody) {
		t.Errorf("Unexpected body: %v", string(body))
	} else if !strings.Contains(string(body), responseBody) {
		t.Errorf("Unexpected body: %v", string(body))
	}
}

func postRequestTest(t *testing.T, path string, requestBody string, responseStatus int, responseBody string) {
	req, err := http.NewRequest(http.MethodPost, host+":"+port+path, bytes.NewReader([]byte(requestBody)))
	if err != nil {
		t.Fatal(err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	} else if resp.StatusCode != responseStatus {
		t.Errorf("Unexpected status: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	} else if string(body) != responseBody {
		t.Errorf("Unexpected body: %v", string(body))
	}
}

func TestWelcome(t *testing.T) {
	getRequestTest(t, home, http.StatusOK, "Welcome to Chirpy", false)
}

func TestLogo(t *testing.T) {
	getRequestTest(t, assets, http.StatusOK, `<a href="logo.png">logo.png</a>`, false)
}

func TestHealth(t *testing.T) {
	getRequestTest(t, healthAPI, http.StatusOK, "OK", true)
}

func TestMetrics(t *testing.T) {
	getRequestTest(t, resetAPI, http.StatusOK, "", true)
	getRequestTest(t, metricsAPI, http.StatusOK, "0 times", false)
	getRequestTest(t, home, http.StatusOK, "", false)
	getRequestTest(t, metricsAPI, http.StatusOK, "1 times", false)
}

func TestChirp(t *testing.T) {
	postRequestTest(t, chirpAPI, `{"body": "hello"}`, http.StatusOK, `{"cleaned_body":"hello"}`)
}

func TestChirpLengthLimit(t *testing.T) {
	postRequestTest(t, chirpAPI, `{"body": "`+strings.Repeat("hello", (chirpLengthLimit/len("hello"))+1)+`"}`, http.StatusBadRequest, `{"error":"Chirp is too long"}`)
}

func TestChirpProfanityFilter(t *testing.T) {
	for _, profanity := range profanities {
		t.Run(profanity, func(t *testing.T) {
			postRequestTest(t, chirpAPI, `{"body": "abc `+profanity+` 123"}`, http.StatusOK, `{"cleaned_body":"abc **** 123"}`)
		})
	}
}
