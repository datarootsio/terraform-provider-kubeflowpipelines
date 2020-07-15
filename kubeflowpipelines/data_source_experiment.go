package kubeflowpipelines

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKubeflowPipelinesExperiment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubeflowPipelinesExperimentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubeflowPipelinesExperimentRead(d *schema.ResourceData, meta interface{}) error {

	experiment, err := readExperiment(meta, d.Get("id").(string), d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("%s", err)
	}

	if experiment.ID == "" {
		return fmt.Errorf("unable to get experiment: %s", d.Get("id").(string))
	}

	d.SetId(experiment.ID)
	d.Set("name", experiment.Name)
	d.Set("description", experiment.Description)

	return nil
}
