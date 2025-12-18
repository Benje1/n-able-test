package servive_monitor

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type ServiceMonitor struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Timeout uint   `yaml:"timeout_ms"`
}

type ServerResponse struct {
	Status    ServerStatus `json:"status"`
	Timestamp string       `json:"timestamp"`
	Service   []Response   `json:"services"`
}

// Reducing typo possibility
type ServerStatus string

const (
	ServerHealthy  ServerStatus = "healthy"
	ServerDegraded ServerStatus = "degraded"
	ServerDown     ServerStatus = "down"
)

func SetupServiceMonitor() (ServiceMonitor, error) {
	data, err := os.ReadFile("services.yaml")
	if err != nil {
		return ServiceMonitor{}, fmt.Errorf("could not read yaml file, file is either missing or not named 'services.yaml': %w", err)
	}

	var service ServiceMonitor
	if err := yaml.Unmarshal(data, &service); err != nil {
		return ServiceMonitor{}, fmt.Errorf("could not read yaml file, file mignt not be in correct format: %w", err)
	}

	return service, nil
}

func (sm ServiceMonitor) GetServiceStatus() ServerResponse {
	reses := sm.CallServices()
	start := time.Now().UTC()

	return ServerResponse{
		Status:    checkHealth(reses),
		Timestamp: start.Format(time.RFC3339),
		Service:   reses,
	}
}

func (sm ServiceMonitor) CallServices() []Response {
	client := http.Client{}

	results := make(chan Response)
	var wg sync.WaitGroup
	for _, server := range sm.Services {
		wg.Add(1)

		go func(srv Service) {
			defer wg.Done()

			res := CallEndpoints(srv, client)
			results <- res
		}(server)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var responses []Response
	for res := range results {
		responses = append(responses, res)
	}

	return responses
}

func checkHealth(responses []Response) ServerStatus {
	checker := make(map[string]struct{})
	for _, res := range responses {
		checker[string(res.Status)] = struct{}{}
	}
	// If map is larger than 1 then we have mixed statuses
	if len(checker) > 1 {
		return ServerDegraded
	}
	// Otherwise check if healthy
	if _, ok := checker[string(Healthy)]; ok {
		return ServerHealthy
	}
	// Default down
	return ServerDown
}
