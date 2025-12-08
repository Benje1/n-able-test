package ServiceMonitor

import (
	"fmt"
	"os"
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

func (sm ServiceMonitor) GetServiceStatus() (ServerResponse, error) {
	reses, err := sm.CallServices()
	start := time.Now().UTC()
	if err != nil {
		return ServerResponse{}, err
	}
	return ServerResponse{
		Status:    checkHealth(reses),
		Timestamp: start.Format(time.RFC3339),
		Service:   reses,
	}, nil

}

func (sm ServiceMonitor) CallServices() ([]Response, error) {
	var errs []error
	var responses []Response
	for _, server := range sm.Services {
		// Would have been nice to paralalise this so the total does not take too long
		res, err := CallEndpoints(server)
		if err != nil {
			errs = append(errs, err)
		} else {
			res = res.UpdateFileds()
			responses = append(responses, res)
		}
	}

	if len(errs) != 0 {
		return nil, fmt.Errorf("there was an issue reaching the services, check internet connection and try again")
	}
	return responses, nil
}

func checkHealth(responses []Response) ServerStatus {
	checker := make(map[string]bool)
	for _, res := range responses {
		checker[string(res.Status)] = true
	}
	if len(checker) > 1 {
		return ServerDegraded
	}
	if _, ok := checker[string(Healthy)]; ok {
		return ServerHealthy
	}
	return ServerDown
}
