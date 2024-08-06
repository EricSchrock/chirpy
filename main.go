package main

import (
	"log"
	"net/http"
)

var port string = "8080"
var home string = "/app"
var assets string = home + "/assets"

func main() {
	log.Println("Starting server...")

	var apiCfg apiConfig
	mux := http.NewServeMux()

	// Front-end website
	mux.Handle(home+"/*", http.StripPrefix(home+"/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	// Back-end APIs (health)
	mux.HandleFunc("GET "+healthAPI, healthHandler)

	// Back-end APIs (metrics)
	mux.HandleFunc("GET "+metricsAPI, apiCfg.metricsHandler)
	mux.HandleFunc(resetAPI, apiCfg.resetHandler)

	// Back-end APIs (chirps)
	mux.HandleFunc("POST "+chirpAPI, validateChirpHandler)

	corsMux := middlewareCors(mux)
	server := &http.Server{Addr: ":" + port, Handler: corsMux}
	err := server.ListenAndServe()
	log.Fatal(err)
}
