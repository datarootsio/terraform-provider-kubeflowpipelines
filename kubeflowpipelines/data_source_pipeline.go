package kubeflowpipelines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_client/pipeline_service"
)

func dataSourceKubeflowPipelinesPipeline() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubeflowPipelinesPipelineRead,

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

func dataSourceKubeflowPipelinesPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Pipeline
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" {
		resp, err := client.PipelineService.ListPipelines(nil, nil)
		if err != nil {
			return fmt.Errorf("unable to get list of pipelines: %s", name)
		}

		pipelineFound := false

		for _, item := range resp.Payload.Pipelines {
			if item.Name == name {
				d.SetId(item.ID)
				d.Set("name", item.Name)
				d.Set("description", item.Description)
				pipelineFound = true
				break
			}
		}

		if !pipelineFound {
			return fmt.Errorf("unable to get pipeline: %s", name)
		}
	} else {
		pipelineParams := pipeline_service.GetPipelineParams{
			ID:      id,
			Context: context,
		}

		resp, err := client.PipelineService.GetPipeline(&pipelineParams, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("unable to get pipeline: %s", err)
		}

		d.SetId(resp.Payload.ID)
		d.Set("name", resp.Payload.Name)
		d.Set("description", resp.Payload.Description)
	}

	return nil
}
