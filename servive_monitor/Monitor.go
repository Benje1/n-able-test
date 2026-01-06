package servive_monitor

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Monitor struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Timeout uint   `yaml:"timeout_ms"`
}

type MonitorResponse struct {
	Status    MonitorStatus     `json:"status"`
	Timestamp string            `json:"timestamp"`
	Service   []ServiceResponse `json:"services"`
}

// Reducing typo possibility
type MonitorStatus string

const (
	ServerHealthy  MonitorStatus = "healthy"
	ServerDegraded MonitorStatus = "degraded"
	ServerDown     MonitorStatus = "down"
)

func SetupServiceMonitor() (Monitor, error) {
	data, err := os.ReadFile("services.yaml")
	if err != nil {
		return Monitor{}, fmt.Errorf("could not read yaml file, file is either missing or not named 'services.yaml': %w", err)
	}

	var service Monitor
	if err := yaml.Unmarshal(data, &service); err != nil {
		return Monitor{}, fmt.Errorf("could not read yaml file, file mignt not be in correct format: %w", err)
	}

	return service, nil
}

func (sm Monitor) GetServiceStatus() MonitorResponse {
	reses := sm.CallServices()
	start := time.Now().UTC()

	return MonitorResponse{
		Status:    checkHealth(reses),
		Timestamp: start.Format(time.RFC3339),
		Service:   reses,
	}
}

func (sm Monitor) CallServices() []ServiceResponse {
	client := http.Client{}

	results := make(chan ServiceResponse)
	jobs := make(chan Service)

	// This could/should be an env var
	const maxWorkers int = 5
	var wg sync.WaitGroup

	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go worker(&client, jobs, results, &wg)
	}

	go func() {
		for _, server := range sm.Services {
			jobs <- server
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var responses []ServiceResponse
	for res := range results {
		responses = append(responses, res)
	}

	return responses
}

func checkHealth(responses []ServiceResponse) MonitorStatus {
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
