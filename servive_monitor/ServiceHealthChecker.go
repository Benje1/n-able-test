package servive_monitor

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type ServiceResponse struct {
	Name         string        `json:"name"`
	Status       ServiceStatus `json:"status"`
	ResponseTime *uint         `json:"response_time_ms,omitempty"`
	Error        *string       `json:"error,omitempty"`
}

type ServiceStatus string

const (
	Healthy ServiceStatus = "healthy"
	Down    ServiceStatus = "down"
)

// There must be a better way of doing this but I was short on time
// Only returning the properties that are wanted
func (res ServiceResponse) SanitiseFields() ServiceResponse {
	if res.Status == Down {
		return ServiceResponse{
			Name:   res.Name,
			Status: res.Status,
			Error:  res.Error,
		}
	}
	return ServiceResponse{
		Name:         res.Name,
		Status:       res.Status,
		ResponseTime: res.ResponseTime,
	}
}

func CallEndpoints(service Service, client *http.Client) ServiceResponse {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(service.Timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, service.Url, nil)

	if err != nil {
		errStr := "could not make requst"
		return ServiceResponse{
			Name:   service.Name,
			Status: "could not make request, error in creating request",
			Error:  &errStr,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		errStr := err.Error()
		return ServiceResponse{
			Name:   service.Name,
			Status: "could not complete requst",
			Error:  &errStr,
		}
	}

	defer resp.Body.Close()

	duration := uint(time.Since(start).Milliseconds())

	respo := ServiceResponse{
		Name:   service.Name,
		Status: getStatusFromCode(resp.StatusCode),
		// Should have a function that modifies this error
		Error:        &resp.Status,
		ResponseTime: &duration,
	}

	return respo.SanitiseFields()
}

func getStatusFromCode(code int) ServiceStatus {
	if code >= 200 && code < 300 {
		return Healthy
	}

	return Down
}

func worker(
	client *http.Client,
	jobs <-chan Service,
	results chan<- ServiceResponse,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for srv := range jobs {
		results <- CallEndpoints(srv, client)
	}
}
