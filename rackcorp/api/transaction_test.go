package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestTransactionCreate(t *testing.T) {
	defer gock.Off()

	const objectId = "879"
	responseBody := getTestDataString(t, "rctransaction.create.responseBody.json")

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	transaction, err := client.TransactionCreate(
		TransactionTypeStartup,
		TransactionObjectTypeDevice,
		objectId,
		false)
	assert.Nil(t, err, "TransactionCreate error")

	assert.Equal(t, "141414", transaction.TransactionId, "TransactionId")
	assert.Equal(t, "STARTUP", transaction.Type, "Type")
	assert.Equal(t, "DEVICE", transaction.ObjectType, "ObjectType")
	assert.Equal(t, objectId, transaction.ObjectId, "ObjectId")
	assert.False(t, transaction.ConfirmationRequired, "ConfirmationRequired")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}
