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

	id, name, description, _, err := readExperiment(meta, d.Get("name").(string), d.Get("id").(string))

	if err != nil {
		return fmt.Errorf("unable to get experiment: %s", err)
	}

	d.SetId(id)
	d.Set("name", name)
	d.Set("description", description)

	return nil
}
