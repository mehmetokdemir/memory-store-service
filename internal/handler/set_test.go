package handler

import (
	// Go imports
	"errors"
	"testing"
)

type memoryTest struct {
	Key   string
	Value string
}

func TestSetValue(t *testing.T) {
	key := "key"
	val := "value"

	//
	// I didn't want to combine "if" commands to check all controls
	//

	memory, err := setValue(key, val)
	if err != nil  {
		t.Fatalf("error:%s", err.Error())
	}

	if memory == nil {
		t.Fatal("memory struct is empty")
	}

	if memory.Value != val || memory.Key != key {
		t.Fatalf("got key: %s, value: %s, want key: %s, value: %s", memory.Key, memory.Value, key, val)
	}
}

func setValue(key, val string) (*memoryTest, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	if val == "" {
		return nil, errors.New("value is empty")
	}

	return &memoryTest{
		Key:   key,
		Value: val,
	}, nil
}
