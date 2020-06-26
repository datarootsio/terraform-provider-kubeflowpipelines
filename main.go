package main

import (
	"github.com/datarootsio/terraform-provider-kubeflow-pipelines/kubeflowpipelines"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kubeflowpipelines.Provider})
}
