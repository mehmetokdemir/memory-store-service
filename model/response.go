package model

type Store struct {
	Key   string `json:"key" extensions:"x-order=1" example:"foo"`   // Key of the store
	Value string `json:"value" extensions:"x-order=1" example:"bar"` // Value of the store
}

type DescriptionEnum string

const (
	DescriptionEnumSuccess         DescriptionEnum = "Success"
	DescriptionEnumBodyError       DescriptionEnum = "Request body or parameters wrong"
	DescriptionEnumBodyDecodeError DescriptionEnum = "Can not decode data"
	DescriptionEnumServerError     DescriptionEnum = "Server error"

	DescriptionEnumCannotGetCurrencies DescriptionEnum = "Can not get currencies from coin market cap"
)

type ApiResponse struct {
	StatusCode  int             `json:"status_code"`           // Status code of the response
	Description DescriptionEnum `json:"description,omitempty"` // Description of the response
	Data        interface{}     `json:"data,omitempty"`        // Data of the response
}

func GenerateResponse(status int, description DescriptionEnum, data interface{}) ApiResponse {
	return ApiResponse{
		StatusCode:  status,
		Description: description,
		Data:        data,
	}
}
