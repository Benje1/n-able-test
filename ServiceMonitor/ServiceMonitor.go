package servicemonitor_test

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ServiceMonitor struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Timeout uint   `yaml:"timeout_ms"`
}

func SetupServiceMonitor() (ServiceMonitor, error) {
	data, err := os.ReadFile("../services.yaml")
	if err != nil {
		return ServiceMonitor{}, err
	}

	var service ServiceMonitor
	if err := yaml.Unmarshal(data, &service); err != nil {
		return ServiceMonitor{}, err
	}

	return service, nil
}
