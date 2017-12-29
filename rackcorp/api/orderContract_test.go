package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestOrderContractGet(t *testing.T) {
	defer gock.Off()

	const contractId = "543"
	const responseBody = `{"contract":{"contractId":"543","customerId":"456","serviceBillId":"135","contractInfo":"1_","currency":"6","created":1514337760,"modified":1514337781,"lastBilled":false,"referenceID":"2816","notes":null,"status":"ACTIVE","type":"VIRTUALSERVER","deviceID":"678"},"code":"OK","message":"Contract lookup successful"}`

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	contract, err := client.OrderContractGet(contractId)
	assert.Nil(t, err, "OrderContractGet error")

	assert.Equal(t, "543", contract.ContractId, "ContractId")
	assert.Equal(t, "456", contract.CustomerId, "CustomerId")
	assert.Equal(t, "678", contract.DeviceId, "DeviceId")
	assert.Equal(t, "ACTIVE", contract.Status, "Status")
	assert.Equal(t, "VIRTUALSERVER", contract.Type, "Type")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}
