package servive_monitor

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

		res := server.GetServiceStatus()

		if res.Status != ServerHealthy {
			t.Errorf("server should be healty: %s", res.Status)
		}

		if res.Service[0].Status != Healthy {
			t.Errorf("server's service should be healthy: %s", res.Service[0].Status)
		}
	})

	t.Run("Test mixed response", func(t *testing.T) {
		server := ServiceMonitor{
			Services: []Service{
				Service{
					Name:    "200",
					Url:     "https://htt.pavonz.com/200",
					Timeout: 500,
				},
				Service{
					Name:    "Teapot",
					Url:     "https://htt.pavonz.com/418",
					Timeout: 500,
				},
			},
		}

		res := server.GetServiceStatus()

		if res.Status != ServerDegraded {
			t.Errorf("server should be healty: %s", res.Status)
		}

		for _, rep := range res.Service {
			if rep.Name == "200" && rep.Status != Healthy {
				t.Errorf("server's service should be healthy: %s", res.Service[0].Status)
			}
			if rep.Name == "Teapot" && rep.Status != Down {
				t.Errorf("server's service should be healthy: %s", res.Service[0].Status)
			}
		}
	})

	t.Run("Test service down", func(t *testing.T) {
		server := ServiceMonitor{
			Services: []Service{
				Service{
					Name:    "500",
					Url:     "https://htt.pavonz.com/500",
					Timeout: 500,
				},
			},
		}

		res := server.GetServiceStatus()

		if res.Status != ServerDown {
			t.Errorf("server should be healty: %s", res.Status)
		}

		if res.Service[0].Status != Down {
			t.Errorf("server's service should be healthy: %s", res.Service[0].Status)
		}
	})
}
