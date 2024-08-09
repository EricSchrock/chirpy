package main

import (
	"log"
	"net/http"

	"github.com/EricSchrock/chirpy/internal/api"
)

var port string = "8080"
var home string = "/app"
var assets string = home + "/assets"

func main() {
	log.Println("Starting server...")

	var apiCfg api.APIConfig
	mux := http.NewServeMux()

	// Front-end website
	mux.Handle(home+"/*", http.StripPrefix(home+"/", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))

	// Back-end APIs (health)
	mux.HandleFunc("GET "+api.HealthAPI, api.HealthHandler)

	// Back-end APIs (metrics)
	mux.HandleFunc("GET "+api.MetricsAPI, apiCfg.MetricsHandler)
	mux.HandleFunc(api.ResetAPI, apiCfg.ResetHandler)

	// Back-end APIs (chirps)
	mux.HandleFunc("POST "+api.ChirpAPI, api.PostChirpHandler)
	mux.HandleFunc("GET "+api.ChirpAPI, api.GetChirpHandler)

	corsMux := middlewareCors(mux)
	server := &http.Server{Addr: ":" + port, Handler: corsMux}
	err := server.ListenAndServe()
	log.Fatal(err)
}
