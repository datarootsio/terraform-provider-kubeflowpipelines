package kubeflowpipelines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_client/pipeline_service"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/run_client/run_service"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/run_model"
)

func resourceKubeflowPipelinesRun() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeflowPipelinesRunCreate,
		Read:   resourceKubeflowPipelinesRunRead,
		Delete: resourceKubeflowPipelinesRunDelete,

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
				Default:      true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"service_account": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"experiment_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"pipeline_spec": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_version_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceKubeflowPipelinesRunCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Run
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	serviceAccount := d.Get("service_account").(string)
	experimentId := d.Get("experiment_id").(string)
	pipelineVersionID, pipelineSpec, err := runExpandPipelineSpec(meta, d.Get("pipeline_spec").([]interface{}))

	if err != nil {
		return fmt.Errorf("unable to get pipeline: %s", err)
	}

	apiJob := run_model.APIRun{
		Name:         name,
		PipelineSpec: pipelineSpec,
	}

	if description != "" {
		apiJob.Description = description
	}

	if serviceAccount != "" {
		apiJob.ServiceAccount = serviceAccount
	}

	apiJob.ResourceReferences = []*run_model.APIResourceReference{
		&run_model.APIResourceReference{
			Key: &run_model.APIResourceKey{
				ID:   pipelineVersionID,
				Type: "PIPELINE_VERSION",
			},
			Relationship: "CREATOR",
		},
	}

	if experimentId != "" {
		apiJob.ResourceReferences = append(apiJob.ResourceReferences, &run_model.APIResourceReference{
			Key: &run_model.APIResourceKey{
				ID:   experimentId,
				Type: "EXPERIMENT",
			},
			Relationship: "OWNER",
		})
	}

	runParams := run_service.CreateRunParams{
		Body:    &apiJob,
		Context: context,
	}

	resp, err := client.RunService.CreateRun(&runParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create run: %s", err)
	}

	d.Set("name", resp.Payload.Run.Name)
	d.Set("description", resp.Payload.Run.Description)
	d.Set("created_at", resp.Payload.Run.CreatedAt)
	d.Set("resource_references", resp.Payload.Run.ResourceReferences)
	d.Set("pipeline_spec", pipelineSpec)
	d.Set("experiment_id", experimentId)
	d.Set("service_account", serviceAccount)

	d.SetId(resp.Payload.Run.ID)

	return resourceKubeflowPipelinesRunRead(d, meta)
}

func resourceKubeflowPipelinesRunRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Run
	context := meta.(*Meta).Context

	id := d.Id()

	runParams := run_service.GetRunParams{
		RunID:   id,
		Context: context,
	}

	resp, err := client.RunService.GetRun(&runParams, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to get run: %s", err)
	}

	d.SetId(resp.Payload.Run.ID)
	d.Set("name", resp.Payload.Run.Name)
	d.Set("description", resp.Payload.Run.Description)
	d.Set("created_at", resp.Payload.Run.CreatedAt)

	return nil
}

func resourceKubeflowPipelinesRunDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Run
	context := meta.(*Meta).Context

	id := d.Id()

	runParams := run_service.ArchiveRunParams{
		ID:      id,
		Context: context,
	}

	_, err := client.RunService.ArchiveRun(&runParams, nil)
	if err != nil {
		return fmt.Errorf("unable to delete job: %s", err)
	}

	d.SetId("")

	return nil
}

func runExpandPipelineSpec(meta interface{}, input []interface{}) (string, *run_model.APIPipelineSpec, error) {
	if len(input) == 0 {
		return "", nil, nil
	}

	values := input[0].(map[string]interface{})

	result := run_model.APIPipelineSpec{}

	pipelineVersionID := ""

	if v, ok := values["pipeline_version_id"]; ok {
		pipelineVersionID = string(v.(string))
	}

	pipelineID, err := getParentPipeline(meta, pipelineVersionID)

	if err != nil {
		return "", nil, fmt.Errorf("unable to get pipeline version: %s", err)
	}

	result.PipelineID = pipelineID

	if v, ok := values["parameters"]; ok {

		parameters := []*run_model.APIParameter{}

		for name, value := range v.(map[string]interface{}) {
			param := run_model.APIParameter{
				Name:  name,
				Value: value.(string),
			}

			parameters = append(parameters, &param)
		}
		result.Parameters = parameters
	}

	return pipelineVersionID, &result, nil
}

func getParentPipeline(meta interface{}, pipelineVersionID string) (string, error) {
	client := meta.(*Meta).Pipeline
	context := meta.(*Meta).Context

	pipelineParams := pipeline_service.GetPipelineVersionParams{
		VersionID: pipelineVersionID,
		Context:   context,
	}

	resp, err := client.PipelineService.GetPipelineVersion(&pipelineParams, nil)
	if err != nil {
		return "", fmt.Errorf("unable to get pipeline version: %s", err)
	}

	for reference := range resp.Payload.ResourceReferences {
		if resp.Payload.ResourceReferences[reference].Relationship == "OWNER" && resp.Payload.ResourceReferences[reference].Key.Type == "PIPELINE" {
			return resp.Payload.ResourceReferences[reference].Key.ID, nil
		}
	}

	return "", fmt.Errorf("unable to get parent pipeline ID: %s", err)

}
