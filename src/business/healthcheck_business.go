package business

import "solace-events-consumer/src/models"

func HealthStatus() models.HealthCheck {
	return models.HealthCheck{
		Status: "The app is healthy",
	}
}
