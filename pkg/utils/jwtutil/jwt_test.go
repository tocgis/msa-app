package jwtutil

import (
	"testing"
)

func TestCreateJwtToken(t *testing.T) {
	jwtToken, err := CreateJwtToken("Cheng",2)
	if err != nil {
		t.Error(err)
	}
	t.Log(jwtToken)
	jwtInfo, err := ParseToken(jwtToken)
	if err != nil {
		t.Error(err)
	}
	t.Log(jwtInfo)
}