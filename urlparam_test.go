package neo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	param := &UrlParam{
		"Key": "Value",
	}

	assert.Equal(t, "Value", param.Get("Key"))
}

func TestExist(t *testing.T) {
	param := &UrlParam{
		"Key": "Value",
	}

	assert.True(t, param.Exist("Key"))
}
