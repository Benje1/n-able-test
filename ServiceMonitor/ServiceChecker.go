package ServiceMonitor

import (
	"net/http"
	"time"
)

type Response struct {
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	ResponseTime *uint   `json:"response_time_ms"`
	Error        *string `json:"error"`
}

// There must be a better way of doing this but I was short on time
func (res Response) GetFileds() Response {
	if res.Status == "down" {
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

	ducration := uint(time.Since(start))
	respo := Response{
		Name:         service.Name,
		Status:       getStatusFromCode(resp.StatusCode),
		Error:        &resp.Status,
		ResponseTime: &ducration,
	}

	return respo, err
}

func getStatusFromCode(code int) string {
	if code >= 200 && code < 300 {
		return "healthy"
	}
	return "down"

}
