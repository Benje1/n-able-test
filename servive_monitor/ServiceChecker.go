package servive_monitor

import (
	"context"
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

func CallEndpoints(service Service, client http.Client) Response {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(service.Timeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, service.Url, nil)

	if err != nil {
		errStr := "could not make requst"
		return Response{
			Name:   service.Name,
			Status: "could not make request, error in creating request",
			Error:  &errStr,
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		errStr := err.Error()
		return Response{
			Name:   service.Name,
			Status: "could not complete requst",
			Error:  &errStr,
		}
	}

	defer resp.Body.Close()

	duration := uint(time.Since(start).Milliseconds())

	respo := Response{
		Name:   service.Name,
		Status: getStatusFromCode(resp.StatusCode),
		// Should have a function that modifies this error
		Error:        &resp.Status,
		ResponseTime: &duration,
	}

	return respo.UpdateFileds()
}

func getStatusFromCode(code int) Status {
	if code >= 200 && code < 300 {
		return Healthy
	}

	return Down
}
