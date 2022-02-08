package handler

import (
	// Go imports
	"errors"
	"testing"
)

func TestGetValue(t *testing.T) {
	var key str = "key"
	expectedValue := "value"

	val, err := key.getValue()
	if err != nil || val != expectedValue {
		t.Errorf("error:%s\n got %s, want %s", err.Error(), val, expectedValue)
	}
}

type str string

func (s str) getValue() (string, error) {

	if s == "" {
		return "", errors.New("key is empty")
	}

	value := "value"
	if s != "key" {
		return "", errors.New("invalid key")
	}

	return value, nil
}
