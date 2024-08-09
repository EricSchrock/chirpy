package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/EricSchrock/chirpy/internal/api"
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

func postRequestTest(t *testing.T, path string, requestBody string, responseStatus int, responseBody string, exactMatch bool) {
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
	} else if exactMatch && string(body) != responseBody {
		t.Errorf("Unexpected body: %v", string(body))
	} else if !strings.Contains(string(body), responseBody) {
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
	getRequestTest(t, api.HealthAPI, http.StatusOK, "OK", true)
}

func TestMetrics(t *testing.T) {
	getRequestTest(t, api.ResetAPI, http.StatusOK, "", true)
	getRequestTest(t, api.MetricsAPI, http.StatusOK, "0 times", false)
	getRequestTest(t, home, http.StatusOK, "", false)
	getRequestTest(t, api.MetricsAPI, http.StatusOK, "1 times", false)
	getRequestTest(t, home, http.StatusOK, "", false)
	getRequestTest(t, api.MetricsAPI, http.StatusOK, "2 times", false)
}

func TestChirps(t *testing.T) {
	postRequestTest(t, api.ChirpAPI, `{"body": "hello"}`, http.StatusCreated, `{"id":1,"body":"hello"}`, true)
	postRequestTest(t, api.ChirpAPI, `{"body": "world"}`, http.StatusCreated, `{"id":2,"body":"world"}`, true)
	getRequestTest(t, api.ChirpAPI, http.StatusOK, `[{"id":1,"body":"hello"},{"id":2,"body":"world"}]`, true)
}

func TestChirpLengthLimit(t *testing.T) {
	postRequestTest(t, api.ChirpAPI, `{"body": "`+strings.Repeat("hello", (api.ChirpLengthLimit/len("hello"))+1)+`"}`, http.StatusBadRequest, `{"error":"Chirp is too long"}`, true)
}

func TestChirpProfanityFilter(t *testing.T) {
	for _, profanity := range api.Profanities {
		t.Run(profanity, func(t *testing.T) {
			postRequestTest(t, api.ChirpAPI, `{"body": "abc `+profanity+` 123"}`, http.StatusCreated, `"body":"abc **** 123"`, false)
		})
	}
}
