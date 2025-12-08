package ServiceMonitor

import (
	"fmt"
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
	data, err := os.ReadFile("services.yaml")
	if err != nil {
		return ServiceMonitor{}, fmt.Errorf("could not read yaml file, file is either missing or not named 'services.yaml': %w", err)
	}

	var service ServiceMonitor
	if err := yaml.Unmarshal(data, &service); err != nil {
		return ServiceMonitor{}, fmt.Errorf("could not read yaml file, file mignt not be in correct format", err)
	}

	return service, nil
}

func (sm ServiceMonitor) CallServices() (any, error) {
	var errs []error
	for _, server := range sm.Services {
		// Would have been nice to paralalise this so the total does not take too long
		res, err := CallEndpoints(server)
		if err != nil {
			errs = append(errs, err)
		}
		fmt.Println(res.Name, res.Status)
	}

	if len(errs) != 0 {
		return nil, fmt.Errorf("There was an issue reaching the services, check internet connection and try again")
	}
	return nil, nil
}
