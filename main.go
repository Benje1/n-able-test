package main

import (
	"encoding/json"
	"n-able-test/ServiceMonitor"
	"net/http"
)

var server ServiceMonitor.ServiceMonitor

// Main sets up the service by reading the yaml config file
// Each call to the endpoint calls the same existing server
func main() {
	var err error
	server, err = ServiceMonitor.SetupServiceMonitor()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/health/aggregate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status, err := server.GetServiceStatus()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	})
}
