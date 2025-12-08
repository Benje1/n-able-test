package ServiceMonitor

import (
	"net/http"
	"time"
)

type Response struct {
	Name         string  `json:"name"`
	Status       Status  `json:"status"`
	ResponseTime *uint   `json:"response_time_ms,omitempty"`
	Error        *string `json:"error,omitempty"`
}

type Status string

const (
	Healthy Status = "healthy"
	Down    Status = "down"
)

// There must be a better way of doing this but I was short on time
// Only returning the properties that are wanted
func (res Response) UpdateFileds() Response {
	if res.Status == Down {
		return Response{
			Name:   res.Name,
			Status: res.Status,
			Error:  res.Error,
		}
	}
	return Response{
		Name:         res.Name,
		Status:       res.Status,
		ResponseTime: res.ResponseTime,
	}
}

func CallEndpoints(service Service) (Response, error) {
	start := time.Now()
	resp, err := http.Get(service.Url)
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	duration := uint(time.Since(start).Milliseconds())
	respo := Response{
		Name:         service.Name,
		Status:       getStatusFromCode(service, resp.StatusCode, duration),
		Error:        &resp.Status,
		ResponseTime: &duration,
	}

	return respo, err
}

func getStatusFromCode(service Service, code int, duration uint) Status {
	if code >= 200 && code < 300 || service.Timeout > duration {
		return Healthy
	}
	return Down
}
