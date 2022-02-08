package model

type SetMemory struct {
	Key   *string `json:"key" extensions:"x-order=1" example:"foo"`   // Key of the store
	Value *string `json:"value" extensions:"x-order=2" example:"bar"` // Value of the store
}
