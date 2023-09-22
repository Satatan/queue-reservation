package repository

import (
	"errors"
	"queue_reservation/internal/domain"
	"queue_reservation/internal/models"
	"queue_reservation/internal/models/enum"
	"queue_reservation/internal/repository/dbmodels"
	"queue_reservation/pkg/logx"
)

type tableRepository struct {
}

func NewTableRepository() domain.TableRepositoryInterface {
	return &tableRepository{}
}

func (r *tableRepository) InitializeTables(branch models.Branch) (*models.TableCount, error) {
	tableDB := dbmodels.TableDB
	branchDB := dbmodels.BranchDB
	tableRemainCount := 0

	for bidx := range branchDB {
		// check correct restaurant and branch
		if branchDB[bidx].ID == branch.ID && branchDB[bidx].RestaurantId == branch.RestaurantId {
			if !branchDB[bidx].IsInit {
				// set initialize first call
				branchDB[bidx].IsInit = true
				for tidx := range tableDB {
					if tableDB[tidx].BranchID == branch.ID {
						// set table to empty booking
						tableDB[tidx].BookingId = nil
						// count table remain
						tableRemainCount += 1
					}
				}
			} else {
				// Return error if this API is called again after initialization.
				err := errors.New(enum.ErrorTableInitialized)
				logx.GetLog().Errorf("InitializeTables Repo Error: %s", err.Error())
				return nil, err
			}
		}
	}

	logx.GetLog().Info("Initialize Tables Success")
	result := models.TableCount{
		NumberOfRemainingTables: tableRemainCount,
	}
	return &result, nil
}
