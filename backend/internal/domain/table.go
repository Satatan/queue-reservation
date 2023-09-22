package domain

import "queue_reservation/internal/models"

type TableUsecaseInterface interface {
	InitializeTables(branch models.Branch) (*models.TableCount, error)
}

type TableRepositoryInterface interface {
	InitializeTables(branch models.Branch) (*models.TableCount, error)
}
