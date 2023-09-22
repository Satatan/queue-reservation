package usecase

import (
	"queue_reservation/internal/domain"
	"queue_reservation/internal/models"
)

type tableUsecase struct {
	tableRepo domain.TableRepositoryInterface
}

func NewTableUsecase(tableRepository domain.TableRepositoryInterface) domain.TableUsecaseInterface {
	return &tableUsecase{
		tableRepo: tableRepository,
	}
}

func (u *tableUsecase) InitializeTables(branch models.Branch) (*models.TableCount, error) {
	tableCount, err := u.tableRepo.InitializeTables(branch)
	if err != nil {
		return nil, err
	}
	return tableCount, nil
}
