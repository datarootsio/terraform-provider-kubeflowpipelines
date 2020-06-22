package main

import (
	"github.com/datarootsio/terraform-provider-kubeflow/kubeflow"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kubeflow.Provider})
}
