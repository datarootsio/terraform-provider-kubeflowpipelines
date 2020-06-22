package kubeflow

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKubeflowExperiment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubeflowExperimentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceKubeflowExperimentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment

	name := d.Get("name").(string)

	resp, err := client.ExperimentService.ListExperiment(nil, nil)
	if err != nil {
		return fmt.Errorf("unable to get list of experiments: %s", name)
	}

	experimentFound := false

	for _, item := range resp.Payload.Experiments {
		if item.Name == name {
			d.SetId(item.ID)
			experimentFound = true
			break
		}
	}

	if !experimentFound {
		return fmt.Errorf("unable to get experiment: %s", name)
	}

	return nil
}
