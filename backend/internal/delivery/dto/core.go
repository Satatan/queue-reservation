package dto

import "queue_reservation/internal/models/enum"

type CoreResponse struct {
	Message enum.ApiMessageResponse `json:"message"`
	Result  interface{}             `json:"result"`
}
