package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestOrderGet(t *testing.T) {
	defer gock.Off()

	responseBody := `{"order":{"orderId":"123","customerId":"456","status":"ACCEPTED","contractId":"789"},"code":"OK","message":"Order lookup successful"}`

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	uuid := "dummy-uuid"
	secret := "dummy-secret"
	orderId := "123"
	client, err := NewClient(uuid, secret)
	assert.Nil(t, err, "NewClient error")

	order, err := client.OrderGet(orderId)
	assert.Nil(t, err, "OrderGet error")

	assert.Equal(t, "123", order.OrderId, "OrderId")
	assert.Equal(t, "456", order.CustomerId, "CustomerId")
	assert.Equal(t, "789", order.ContractId, "ContractId")
	assert.Equal(t, "ACCEPTED", order.Status, "Status")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}
