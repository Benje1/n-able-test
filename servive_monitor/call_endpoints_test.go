package servive_monitor

import "testing"

func TestCallEndpoints(t *testing.T) {
	t.Run("Test gets 200", func(t *testing.T) {
		service := Service{
			Name:    "200",
			Url:     "https://htt.pavonz.com/200",
			Timeout: 5000,
		}

		res, err := CallEndpoints(service)
		if err != nil {
			t.Errorf("no error expected, got: %s", err.Error())
		}

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
		_, err := CallEndpoints(service)
		if err == nil {
			t.Errorf("timeout error expected")
		}
	})

	t.Run("test get non 2xx code", func(t *testing.T) {
		service := Service{
			Name:    "500",
			Url:     "https://htt.pavonz.com/500",
			Timeout: 5000,
		}

		res, err := CallEndpoints(service)
		if err != nil {
			t.Errorf("no error expected, got: %s", err.Error())
		}

		if *res.Error != "500 Internal Server Error" {
			t.Errorf("error should be internal error: got %s", *res.Error)
		}
	})
}
