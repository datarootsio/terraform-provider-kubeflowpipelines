package kubeflowpipelines

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/job_client/job_service"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/job_model"
)

func resourceKubeflowPipelinesJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeflowPipelinesJobCreate,
		Read:   resourceKubeflowPipelinesJobRead,
		Delete: resourceKubeflowPipelinesJobDelete,

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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
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
			"max_concurrency": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"no_catchup": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"pipeline_spec": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_id": {
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
			"trigger": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_schedule": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
									"cron": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"periodic_schedule": {
							Type:     schema.TypeList,
							MinItems: 1,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
									"interval_seconds": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceKubeflowPipelinesJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Job
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	serviceAccount := d.Get("service_account").(string)
	enabled := d.Get("enabled").(bool)
	noCatchup := d.Get("no_catchup").(bool)
	maxConcurrency := int64(d.Get("max_concurrency").(int))
	experimentId := d.Get("experiment_id").(string)

	pipelineSpec := jobExpandPipelineSpec(d.Get("pipeline_spec").([]interface{}))
	trigger := expandTrigger(d.Get("trigger").([]interface{}))

	apiJob := job_model.APIJob{
		Name:           name,
		MaxConcurrency: maxConcurrency,
		Enabled:        enabled,
		NoCatchup:      noCatchup,
		PipelineSpec:   pipelineSpec,
		Trigger:        trigger,
	}

	if description != "" {
		apiJob.Description = description
	}

	if serviceAccount != "" {
		apiJob.ServiceAccount = serviceAccount
	}

	if experimentId != "" {
		apiJob.ResourceReferences = []*job_model.APIResourceReference{
			&job_model.APIResourceReference{
				Key: &job_model.APIResourceKey{
					ID:   experimentId,
					Type: "EXPERIMENT",
				},
				Relationship: "OWNER",
			},
		}
	}

	jobParams := job_service.CreateJobParams{
		Body:    &apiJob,
		Context: context,
	}

	resp, err := client.JobService.CreateJob(&jobParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create job: %s", err)
	}

	d.Set("name", resp.Payload.Name)
	d.Set("description", resp.Payload.Description)
	d.Set("created_at", time.Time(resp.Payload.CreatedAt).Format(time.RFC3339))
	d.Set("resource_references", resp.Payload.ResourceReferences)
	d.Set("trigger", trigger)
	d.Set("pipeline_spec", pipelineSpec)
	d.Set("enabled", enabled)
	d.Set("no_catchup", noCatchup)
	d.Set("experiment_id", experimentId)
	d.Set("service_account", serviceAccount)
	d.Set("max_concurrency", maxConcurrency)

	d.SetId(resp.Payload.ID)

	return resourceKubeflowPipelinesJobRead(d, meta)
}

func resourceKubeflowPipelinesJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Job
	context := meta.(*Meta).Context

	id := d.Id()

	jobParams := job_service.GetJobParams{
		ID:      id,
		Context: context,
	}

	resp, err := client.JobService.GetJob(&jobParams, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to get job: %s", err)
	}

	d.SetId(resp.Payload.ID)
	d.Set("name", resp.Payload.Name)
	d.Set("description", resp.Payload.Description)
	d.Set("created_at", resp.Payload.CreatedAt)

	return nil
}

func resourceKubeflowPipelinesJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Job
	context := meta.(*Meta).Context

	id := d.Id()

	jobParams := job_service.DeleteJobParams{
		ID:      id,
		Context: context,
	}

	_, err := client.JobService.DeleteJob(&jobParams, nil)
	if err != nil {
		return fmt.Errorf("unable to delete job: %s", err)
	}

	d.SetId("")

	return nil
}

func jobExpandPipelineSpec(input []interface{}) *job_model.APIPipelineSpec {
	if len(input) == 0 {
		return nil
	}

	values := input[0].(map[string]interface{})

	result := job_model.APIPipelineSpec{}

	if v, ok := values["pipeline_id"]; ok {
		result.PipelineID = string(v.(string))
	}

	if v, ok := values["parameters"]; ok {

		parameters := []*job_model.APIParameter{}

		for name, value := range v.(map[string]interface{}) {
			param := job_model.APIParameter{
				Name:  name,
				Value: value.(string),
			}

			parameters = append(parameters, &param)
		}
		result.Parameters = parameters
	}

	return &result
}

func expandTrigger(input []interface{}) *job_model.APITrigger {
	if len(input) == 0 {
		return nil
	}

	values := input[0].(map[string]interface{})

	result := job_model.APITrigger{}

	if v, ok := values["cron_schedule"]; ok {
		vv := v.([]interface{})
		cron := vv[0].(map[string]interface{})

		cron_result := job_model.APICronSchedule{}

		if w, ok := cron["cron"]; ok {
			cron_result.Cron = string(w.(string))
		}

		if w, ok := cron["start_time"]; ok {
			date, _ := strfmt.ParseDateTime(string(w.(string)))
			cron_result.StartTime = date
		}

		if w, ok := cron["end_time"]; ok {
			date, _ := strfmt.ParseDateTime(string(w.(string)))
			cron_result.EndTime = date
		}

		result.CronSchedule = &cron_result

	} else if v, ok := values["periodic_schedule"]; ok {
		vv := v.([]interface{})
		periodic := vv[0].(map[string]interface{})

		periodic_result := job_model.APIPeriodicSchedule{}

		if w, ok := periodic["cron"]; ok {
			periodic_result.IntervalSecond = int64(w.(int64))
		}

		if w, ok := periodic["start_time"]; ok {
			date, _ := strfmt.ParseDateTime(string(w.(string)))
			periodic_result.StartTime = date
		}

		if w, ok := periodic["end_time"]; ok {
			date, _ := strfmt.ParseDateTime(string(w.(string)))
			periodic_result.EndTime = date
		}

		result.PeriodicSchedule = &periodic_result

	}

	return &result
}
