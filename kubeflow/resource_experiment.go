package kubeflow

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client/experiment_service"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_model"
)

func resourceKubeflowExperiment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeflowExperimentCreate,
		Read:   resourceKubeflowExperimentRead,
		Delete: resourceKubeflowExperimentDelete,

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
				Type:         schema.TypeString,
				Computed:     true,
				ForceNew:     true,
			},
			"storage_state": {
				Type:         schema.TypeString,
				Computed:     true,
				ForceNew:     true,
			},
			"resources_reference": {
				Type:         schema.TypeString,
				Computed:     true,
				ForceNew:     true,
			},
		},
	}
}

func resourceKubeflowExperimentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	apiExperiment := experiment_model.APIExperiment{Name: name}

	if description != "" {
		apiExperiment.Description = description
	}

	experimentParams := experiment_service.CreateExperimentParams{
		Body: &apiExperiment,
		Context: context,
	}

	resp, err := client.ExperimentService.CreateExperiment(&experimentParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create experiment: %s", err)
	}

	d.Set("name", resp.Payload.Name)
	d.Set("description", resp.Payload.Description)
	d.Set("created_at", resp.Payload.CreatedAt)
	d.Set("resource_references", resp.Payload.ResourceReferences)
	d.Set("storage_state", resp.Payload.StorageState)
	d.SetId(resp.Payload.ID)

	return resourceKubeflowExperimentRead(d, meta)
}

func resourceKubeflowExperimentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	id := d.Id()

	experimentParams := experiment_service.GetExperimentParams{
		ID: id,
		Context: context,
	}

	resp, err := client.ExperimentService.GetExperiment(&experimentParams, nil)
	if err != nil {
		return fmt.Errorf("unable to get experiment")
	}

	d.SetId(resp.Payload.ID)
	d.Set("name", resp.Payload.Name)
	d.Set("description", resp.Payload.Description)
	d.Set("created_at", resp.Payload.CreatedAt)
	d.Set("resource_references", resp.Payload.ResourceReferences)
	d.Set("storage_state", resp.Payload.StorageState)

	return nil
}

func resourceKubeflowExperimentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	id := d.Id()

	experimentParams := experiment_service.DeleteExperimentParams{
		ID: id,
		Context: context,
	}

	_, err := client.ExperimentService.DeleteExperiment(&experimentParams, nil)
	if err != nil {
		return fmt.Errorf("unable to delete experiment: %s", err)
	}

	d.SetId("")

	return nil
}
