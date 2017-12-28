package main // import "github.com/section-io/terraform-provider-rackcorp"

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/section-io/terraform-provider-rackcorp/rackcorp"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rackcorp.Provider})
}
