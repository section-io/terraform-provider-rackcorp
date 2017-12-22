package rackcorp

import (
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
						"image": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
				},
			},
		}
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	orderRequest := OrderRequest{ApiUuid: config.ApiUuid, ApiSecret: config.ApiSecret, Command: "order.create"}

	orderRequestJson, err := json.Marshal(orderRequest)
    if err != nil {
        panic(err)
    }

	order, err := http.Post("https://requestb.in/txetp8tx", "application/json", bytes.NewBuffer(orderRequestJson))

	if err != nil {
		// handle error
	}
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