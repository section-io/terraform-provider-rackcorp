package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestClient(t *testing.T) Client {
	const uuid = "dummy-uuid"
	const secret = "dummy-secret"
	client, err := NewClient(uuid, secret)
	assert.Nil(t, err, "NewClient error")
	return client
}
