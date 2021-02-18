package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type daoMock struct {
}

func (d daoMock) Hello() string {
	return "Hello"
}

func TestSayHelo(t *testing.T) {
	// Mockeamos
	MockedDao = new(daoMock)

	s := NewService()
	assert.Equal(t, "Hello", s.SayHello())

	// Volvemos al original
	MockedDao = nil
}
