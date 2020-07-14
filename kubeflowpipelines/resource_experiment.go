package kubeflowpipelines

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client/experiment_service"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_model"
)

func resourceKubeflowPipelinesExperiment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeflowPipelinesExperimentCreate,
		Read:   resourceKubeflowPipelinesExperimentRead,
		Delete: resourceKubeflowPipelinesExperimentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKubeflowPipelinesExperimentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	apiExperiment := experiment_model.APIExperiment{Name: name}

	if description != "" {
		apiExperiment.Description = description
	}

	experimentParams := experiment_service.CreateExperimentParams{
		Body:    &apiExperiment,
		Context: context,
	}

	resp, err := client.ExperimentService.CreateExperiment(&experimentParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create experiment: %s", err)
	}

	d.Set("name", resp.Payload.Name)
	d.Set("description", resp.Payload.Description)
	d.Set("created_at", time.Time(resp.Payload.CreatedAt).Format(time.RFC3339))
	d.SetId(resp.Payload.ID)

	return resourceKubeflowPipelinesExperimentRead(d, meta)
}

func resourceKubeflowPipelinesExperimentRead(d *schema.ResourceData, meta interface{}) error {
	id, name, description, createdAt, err := readExperiment(meta, d.Id(), "")

	if err != nil {
		return fmt.Errorf("%s", err)
	}

	d.SetId(id)
	d.Set("name", name)
	d.Set("description", description)
	d.Set("created_at", createdAt)

	return nil
}

func resourceKubeflowPipelinesExperimentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	id := d.Id()

	experimentParams := experiment_service.DeleteExperimentParams{
		ID:      id,
		Context: context,
	}

	_, err := client.ExperimentService.DeleteExperiment(&experimentParams, nil)
	if err != nil {
		return fmt.Errorf("unable to delete experiment: %s", err)
	}

	d.SetId("")

	return nil
}
