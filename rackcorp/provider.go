package rackcorp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

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
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackcorp_server": resourceRackcorpServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ApiUuid: d.Get("api_uuid").(string),
		ApiSecret: d.Get("api_secret").(string),
	}

	return config, nil
}

type Config struct {
	ApiUuid string
	ApiSecret string
}