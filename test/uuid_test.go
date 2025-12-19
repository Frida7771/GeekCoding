package test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestGenerateUUID(t *testing.T) {
	s := uuid.NewV4().String()
	print(s)

}
