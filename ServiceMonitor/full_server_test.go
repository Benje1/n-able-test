package ServiceMonitor

import "testing"

// Happy path test only due to lack of time
func TestServiceMonitor(t *testing.T) {
	t.Run("Test happy path server", func(t *testing.T) {

		server := ServiceMonitor{
			Services: []Service{
				Service{
					Name:    "200",
					Url:     "https://htt.pavonz.com/200",
					Timeout: 500,
				},
			},
		}

		res, err := server.GetServiceStatus()
		if err != nil {
			t.Errorf("should not have error: %s", err.Error())
		}

		if res.Status != ServerHealthy {
			t.Errorf("server should be healty: %s", res.Status)
		}

		if res.Service[0].Status != Healthy {
			t.Errorf("server's service should be healthy: %s", res.Service[0].Status)
		}
	})
}
