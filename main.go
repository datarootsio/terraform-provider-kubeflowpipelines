package main

import (
	"github.com/datarootsio/terraform-provider-kubeflowpipelines/kubeflowpipelines"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kubeflowpipelines.Provider})
}
