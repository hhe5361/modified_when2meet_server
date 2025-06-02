package helper

import "testing"

func TestPassword(t *testing.T) {
	key := "test"
	result, err := EncryptPassword(key)
	if err != nil {
		t.Errorf("EncryptPassword not working")
	}
	if !CheckPasswordHash(key, result) {
		t.Errorf("result is not expected")
	}
}
