package enum

type ApiMessageResponse string

const (
	MessageSuccess ApiMessageResponse = "success"
	MessageError   ApiMessageResponse = "error"
)
