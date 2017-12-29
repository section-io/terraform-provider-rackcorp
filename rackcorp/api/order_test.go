package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestOrderConfirm(t *testing.T) {
	defer gock.Off()

	const orderId = "432"
	const responseBody = `{"contractID":[543],"code":"OK","message":"Order confirmed"}`

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	order, err := client.OrderConfirm(orderId)
	assert.Nil(t, err, "OrderConfirm error")

	assert.Contains(t, order.ContractIds, "543", "ContractIds")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}

func TestOrderCreate(t *testing.T) {
	defer gock.Off()

	const productCode = "SERVER_VIRTUAL_PERFORMANCE_AU"
	const customerId = "456"
	productDetails := ProductDetails{
		CpuCount: 1,
		Install: Install{
			OperatingSystem: "UBUNTU14.04_64",
		},
	}

	const responseBody = `{"orderId":123,"changeTxt":"Add NEW SUPPORT: SUPPORTSTD ($0.00)\nAdd NEW IPV6: 16 ($0.00)\n","code":"OK","message":"Order created"}`

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	order, err := client.OrderCreate(productCode, customerId, productDetails)
	assert.Nil(t, err, "OrderCreate error")

	assert.Equal(t, "123", order.OrderId, "OrderId")
	assert.Contains(t, order.ChangeText, "Add NEW", "ChangeText")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}

func TestOrderGet(t *testing.T) {
	defer gock.Off()

	const orderId = "123"
	const responseBody = `{"order":{"orderId":"123","customerId":"456","status":"ACCEPTED","contractId":"789"},"code":"OK","message":"Order lookup successful"}`

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	order, err := client.OrderGet(orderId)
	assert.Nil(t, err, "OrderGet error")

	assert.Equal(t, "123", order.OrderId, "OrderId")
	assert.Equal(t, "456", order.CustomerId, "CustomerId")
	assert.Equal(t, "789", order.ContractId, "ContractId")
	assert.Equal(t, "ACCEPTED", order.Status, "Status")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}
