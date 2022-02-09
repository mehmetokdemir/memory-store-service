package model

type Store struct {
	Key   string `json:"key" extensions:"x-order=1" example:"foo"`   // Key of the store
	Value string `json:"value" extensions:"x-order=1" example:"bar"` // Value of the store
}

type DescriptionEnum string

const (
	DescriptionEnumSuccess          DescriptionEnum = "Success"
	DescriptionEnumBodyError        DescriptionEnum = "Request body or parameters wrong"
	DescriptionEnumBodyDecodeError  DescriptionEnum = "Can not decode data"
	DescriptionEnumBodyReadError    DescriptionEnum = "Can not read data"
	DescriptionEnumValueTypeError   DescriptionEnum = "Value type is not a string"
	DescriptionEnumKeyNotFoundError DescriptionEnum = "Key not found"
	DescriptionEnumInvalidKeyError  DescriptionEnum = "Invalid key"
	DescriptionEnumServerError      DescriptionEnum = "Server error"
)

func (e DescriptionEnum) String() string {
	return string(e)
}

type ApiResponse struct {
	StatusCode  int             `json:"status_code"`           // Status code of the response
	Description DescriptionEnum `json:"description,omitempty"` // Description of the response
	Data        interface{}     `json:"data,omitempty"`        // Data of the response
}

//GenerateResponse create own response
func GenerateResponse(status int, description DescriptionEnum, data interface{}) ApiResponse {
	return ApiResponse{
		StatusCode:  status,
		Description: description,
		Data:        data,
	}
}
