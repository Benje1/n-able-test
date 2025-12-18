package main

import (
	"encoding/json"
	"fmt"
	"n-able-test/servive_monitor"
	"net/http"
)

var server servive_monitor.ServiceMonitor

// Main sets up the service by reading the yaml config file
// Each call to the endpoint calls the same existing server
func main() {
	var err error
	server, err = servive_monitor.SetupServiceMonitor()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/health/aggregate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := server.GetServiceStatus()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	})

	fmt.Println("Server listening on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
