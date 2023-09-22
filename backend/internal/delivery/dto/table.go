package dto

type TableInitializeResponse struct {
	NumberOfTables int `json:"number_of_tables"`
}

func (r *TableInitializeResponse) ToTableInitializeResponse(tableCount *int) {
	r.NumberOfTables = *tableCount
}
