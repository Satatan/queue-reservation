package dto

import "queue_reservation/internal/models"

type BookingRequest struct {
	NumberOfCustomers int `json:"number_of_customers" validate:"required"`
}

type BookingWithTableResponse struct {
	BookingId                 uint64 `json:"booking_id,omitempty"`
	NumberOfNewReservedTables int    `json:"number_of_reserved_table,omitempty"`
	NumberOfRemainingTables   *int   `json:"number_of_remaining_table,omitempty"`
}

func (r *BookingWithTableResponse) ToBookingWithTableResponse(id *uint64, tableCount *models.TableCount) {
	r.BookingId = *id
	r.NumberOfNewReservedTables = tableCount.NumberOfNewReservedTables
	r.NumberOfRemainingTables = &tableCount.NumberOfRemainingTables
}

func (r *BookingWithTableResponse) ToCancelBookingResponse(tableCount *models.TableCount) {
	r.NumberOfRemainingTables = &tableCount.NumberOfRemainingTables
}
