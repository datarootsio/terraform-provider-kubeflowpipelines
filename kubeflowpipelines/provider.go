package kubeflowpipelines

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KUBEFLOW_HOST", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kubeflow_pipelines_experiment": resourceKubeflowPipelinesExperiment(),
			"kubeflow_pipelines_pipeline":   resourceKubeflowPipelinesPipeline(),
			"kubeflow_pipelines_job":        resourceKubeflowPipelinesJob(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"kubeflow_pipelines_experiment": dataSourceKubeflowPipelinesExperiment(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			Host:             d.Get("host").(string),
			terraformVersion: p.TerraformVersion,
		}

		client, err := config.Client()
		if err != nil {
			return nil, err
		}

		client.Context = p.StopContext()

		return client, nil
	}
}
