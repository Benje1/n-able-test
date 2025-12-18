package servive_monitor

import (
	"net/http"
	"testing"
)

func TestCallEndpoints(t *testing.T) {
	client := http.Client{}
	t.Run("Test gets 200", func(t *testing.T) {
		service := Service{
			Name:    "200",
			Url:     "https://htt.pavonz.com/200",
			Timeout: 5000,
		}

		res := CallEndpoints(service, client)

		if res.Status != Healthy {
			t.Error("status should be healthy")
		}

		if res.Error != nil {
			t.Error("should not have error")
		}
	})

	t.Run("test timeout with 200", func(t *testing.T) {
		service := Service{
			Name:    "200",
			Url:     "https://htt.pavonz.com/200",
			Timeout: 1,
		}

		// No response as would be nil
		res := CallEndpoints(service, client)

		if *res.Error != "Get \"https://htt.pavonz.com/200\": context deadline exceeded" {
			t.Error(*res.Error)
		}
	})

	t.Run("test get non 2xx code", func(t *testing.T) {
		service := Service{
			Name:    "500",
			Url:     "https://htt.pavonz.com/500",
			Timeout: 5000,
		}

		res := CallEndpoints(service, client)

		if *res.Error != "500 Internal Server Error" {
			t.Errorf("error should be internal error: got %s", *res.Error)
		}
	})
}
