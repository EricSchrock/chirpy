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

func requestTest(t *testing.T, method string, path string, requestBody string, responseStatus int, responseBody string, exactMatch bool) {
	req, err := http.NewRequest(method, host+":"+port+path, bytes.NewReader([]byte(requestBody)))
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
	requestTest(t, http.MethodGet, home, "", http.StatusOK, "Welcome to Chirpy", false)
}

func TestLogo(t *testing.T) {
	requestTest(t, http.MethodGet, assets, "", http.StatusOK, `<a href="logo.png">logo.png</a>`, false)
}

func TestHealth(t *testing.T) {
	requestTest(t, http.MethodGet, api.HealthAPI, "", http.StatusOK, "OK", true)
}

func TestMetrics(t *testing.T) {
	requestTest(t, http.MethodGet, api.ResetAPI, "", http.StatusOK, "", true)
	requestTest(t, http.MethodGet, api.MetricsAPI, "", http.StatusOK, "0 times", false)
	requestTest(t, http.MethodGet, home, "", http.StatusOK, "", false)
	requestTest(t, http.MethodGet, api.MetricsAPI, "", http.StatusOK, "1 times", false)
	requestTest(t, http.MethodGet, home, "", http.StatusOK, "", false)
	requestTest(t, http.MethodGet, api.MetricsAPI, "", http.StatusOK, "2 times", false)
}

func TestChirps(t *testing.T) {
	requestTest(t, http.MethodPost, api.ChirpAPI, `{"body": "hello"}`, http.StatusCreated, `{"id":1,"body":"hello"}`, true)
	requestTest(t, http.MethodPost, api.ChirpAPI, `{"body": "world"}`, http.StatusCreated, `{"id":2,"body":"world"}`, true)
	requestTest(t, http.MethodGet, api.ChirpAPI, "", http.StatusOK, `[{"id":1,"body":"hello"},{"id":2,"body":"world"}]`, true)
}

func TestChirpLengthLimit(t *testing.T) {
	requestTest(t, http.MethodPost, api.ChirpAPI, `{"body": "`+strings.Repeat("hello", (api.ChirpLengthLimit/len("hello"))+1)+`"}`, http.StatusBadRequest, `{"error":"Chirp is too long"}`, true)
}

func TestChirpProfanityFilter(t *testing.T) {
	for _, profanity := range api.Profanities {
		t.Run(profanity, func(t *testing.T) {
			requestTest(t, http.MethodPost, api.ChirpAPI, `{"body": "abc `+profanity+` 123"}`, http.StatusCreated, `"body":"abc **** 123"`, false)
		})
	}
}
