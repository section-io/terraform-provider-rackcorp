package rackcorp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/section-io/rackcorp-sdk-go/api"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKCORP_API_UUID", nil),
				Description: "The API UUID provided by Rackcorp.",
			},
			"api_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKCORP_API_SECRET", nil),
				Description: "The API secret provided by Rackcorp.",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKCORP_CUSTOMER_ID", nil),
				Description: "Your Rackcorp Customer ID.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackcorp_server": resourceRackcorpServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client, err := api.NewClient(d.Get("api_uuid").(string), d.Get("api_secret").(string))
	if err != nil {
		return nil, err
	}

	config := providerConfig{
		Client:     client,
		CustomerID: d.Get("customer_id").(string),
	}

	return config, nil
}

type providerConfig struct {
	Client     api.Client
	CustomerID string
}
