package security

import (
	"fmt"
	"testing"
)

func Test_HashPassword(t *testing.T) {
	password := "Messi123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(hashedPassword)
	match, err := ComparePassword(hashedPassword, password)
	if err != nil {
		t.Error(err)
	}

	if !match {
		t.Error("passwords do not match")
	}
}
