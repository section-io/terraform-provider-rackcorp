package main // import "github.com/section-io/terraform-provider-rackcorp"

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/section-io/terraform-provider-rackcorp/rackcorp"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return rackcorp.Provider()
		},
	})
}
