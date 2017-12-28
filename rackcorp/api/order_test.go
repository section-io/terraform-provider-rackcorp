package api

import (
	"log"
	"testing"

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
	if err != nil {
		t.Error(err)
	}
	order, err := client.OrderGet(orderId)
	log.Printf("order:%#v, err:%#v\n", order, err)

	if !gock.IsDone() {
		t.Error("gock not done")
	}
}
