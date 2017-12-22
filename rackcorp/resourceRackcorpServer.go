package rackcorp

import (
	"log"
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRackcorpServer() *schema.Resource {
	return &schema.Resource{
			Create: resourceRackcorpServerCreate,
			Delete:   resourceRackcorpServerRead,
			Read:   resourceRackcorpServerRead,
			Schema: map[string]*schema.Schema{
					"country": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"server_class": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
			},
		}
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	orderRequest := OrderRequest {
		ApiUuid: config.ApiUuid, 
		ApiSecret: config.ApiSecret, 
		Command: "order.create", 
		CustomerId: config.CustomerId,
		ProductCode: "SERVER_VIRTUAL_" + d.Get("server_class").(string) + "_" + d.Get("country").(string),
	}

	orderRequestJson, err := json.Marshal(orderRequest)
    if err != nil {
        panic(err)
    }

	order, err := http.Post(config.ApiAddress, "application/json", bytes.NewBuffer(orderRequestJson))

	if err != nil {
        panic(err)
	}

	var orderResponse OrderResponse
	json.NewDecoder(order.Body).Decode(orderResponse)

	log.Printf("%+v\n", orderResponse)

	// catch the order Id, and make the next api call.
	

	defer order.Body.Close()

	return resourceRackcorpServerRead(d, meta)
}

func resourceRackcorpServerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

type OrderRequest struct {
	ApiUuid string `json:"APIUUID"`
	ApiSecret string `json:"APISECRET"`
	Command string `json:"cmd"`
	ProductCode string `json:"productCode"`
	CustomerId string `json:"customerId"`
}

type OrderResponse struct {
	Code string `json:"code"`
	OrderId string `json:"orderId"`
}