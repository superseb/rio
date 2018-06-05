package create

import (
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func populateHealthCheck(c *Create, service *client.Service) error {
	var err error

	hc := &client.HealthConfig{
		HealthyThreshold: int64(c.HealthRetries),
	}

	hc.InitialDelaySeconds, err = ParseSeconds(c.HealthStartPeriod, "--health-start-period")
	if err != nil {
		return err
	}

	hc.IntervalSeconds, err = ParseSeconds(c.HealthInterval, "--health-interval")
	if err != nil {
		return err
	}

	if len(c.HealthCmd) > 0 {
		hc.Test = []string{"CMD-SHELL", c.HealthCmd}
	}

	hc.TimeoutSeconds, err = ParseSeconds(c.HealthTimeout, "--health-timeout")

	if len(c.HealthCmd) > 0 {
		service.Healthcheck = hc
	}

	return err
}
