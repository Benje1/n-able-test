package main

import (
	"n-able-test/ServiceMonitor"
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
}
