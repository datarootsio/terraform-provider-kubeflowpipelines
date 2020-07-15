package kubeflowpipelines

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_client/pipeline_service"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_model"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_upload_client/pipeline_upload_service"
)

func resourceKubeflowPipelinesPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeflowPipelinesPipelineCreate,
		Read:   resourceKubeflowPipelinesPipelineRead,
		Update: resourceKubeflowPipelinesPipelineUpdate,
		Delete: resourceKubeflowPipelinesPipelineDelete,
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf("version_id", func(d *schema.ResourceDiff, meta interface{}) bool {
				return d.HasChange("version") || d.HasChange("file_base64") || d.HasChange("url")
			}),
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ConflictsWith: []string{"url"},
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file_base64", "url"},
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"file_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zip", "tar.gz", "yaml"}, true),
				RequiredWith: []string{"file_base64"},
			},
			"file_base64": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file_base64", "url"},
				ValidateFunc: validation.StringIsBase64,
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func decodeFile(b64 string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "tf-kubeflow-pipeline-")
	if err != nil {
		return "", err
	}
	decoder, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	ioutil.WriteFile(file.Name(), decoder, 0644)
	if err != nil {
		return "", err
	}

	return file.Name(), err
}

func resourceKubeflowPipelinesPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	file := d.Get("file_base64").(string)

	if file != "" {
		return kubeflowCreatePipelineFromFile(d, meta)
	} else {
		return kubeflowCreatePipelineFromUrl(d, meta)
	}
}

func resourceKubeflowPipelinesPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	file := d.Get("file_base64").(string)
	if file != "" {
		return kubeflowCreatePipelineVersionFromFile(d, meta)
	} else {
		return kubeflowCreatePipelineVersionFromUrl(d, meta)
	}
}

func kubeflowCreatePipelineFromFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).PipelineUpload
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	file := d.Get("file_base64").(string)
	fileFormat := d.Get("file_format").(string)
	version := d.Get("version").(string)

	fileName, err := decodeFile(file)
	if err != nil {
		return fmt.Errorf("unable to decode file: %s", err)
	}
	defer os.Remove(fileName)
	decodedFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}

	pipelineParams := pipeline_upload_service.UploadPipelineParams{
		Context:    context,
		Uploadfile: runtime.NamedReader(fileName+"."+fileFormat, decodedFile),
	}

	if name != "" {
		pipelineParams.Name = &name
	}

	if description != "" {
		pipelineParams.Description = &description
	}

	resp, err := client.PipelineUploadService.UploadPipeline(&pipelineParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create pipeline: %s", err)
	}

	d.SetId(resp.Payload.ID)

	if version != "" {
		return kubeflowCreatePipelineVersionFromFile(d, meta)
	} else {
		return resourceKubeflowPipelinesPipelineRead(d, meta)
	}
}

func kubeflowCreatePipelineVersionFromFile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).PipelineUpload
	context := meta.(*Meta).Context

	pipeline_id := d.Id()
	name := d.Get("name").(string)
	file := d.Get("file_base64").(string)
	fileFormat := d.Get("file_format").(string)
	version := d.Get("version").(string)

	computedName := fmt.Sprintf("%s - %s", name, version)

	fileName, err := decodeFile(file)
	if err != nil {
		return fmt.Errorf("unable to decode file: %s", err)
	}
	defer os.Remove(fileName)
	decodedFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}

	pipelineParams := pipeline_upload_service.UploadPipelineVersionParams{
		Context:    context,
		Pipelineid: &pipeline_id,
		Name:       &computedName,
		Uploadfile: runtime.NamedReader(fmt.Sprintf("%s.%s", fileName, fileFormat), decodedFile),
	}

	resp, err := client.PipelineUploadService.UploadPipelineVersion(&pipelineParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create pipeline version: %s", err)
	}

	d.Set("version", version)
	d.Set("version_id", resp.Payload.ID)

	return resourceKubeflowPipelinesPipelineRead(d, meta)
}

func kubeflowCreatePipelineFromUrl(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Pipeline
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	version := d.Get("version").(string)
	url := d.Get("url").(string)

	apiPipeline := pipeline_model.APIPipeline{
		URL: &pipeline_model.APIURL{
			PipelineURL: url,
		},
	}

	if name != "" {
		apiPipeline.Name = name
	}

	if description != "" {
		apiPipeline.Description = description
	}

	pipelineParams := pipeline_service.CreatePipelineParams{
		Body:    &apiPipeline,
		Context: context,
	}

	resp, err := client.PipelineService.CreatePipeline(&pipelineParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create pipeline: %s", err)
	}

	d.SetId(resp.Payload.ID)

	if version != "" {
		return kubeflowCreatePipelineVersionFromUrl(d, meta)
	} else {
		return resourceKubeflowPipelinesPipelineRead(d, meta)
	}
}

func kubeflowCreatePipelineVersionFromUrl(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Pipeline
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	url := d.Get("url").(string)
	pipeline_id := d.Id()

	computedName := fmt.Sprintf("%s - %s", name, version)

	resourceReferences := pipeline_model.APIResourceReference{
		Key: &pipeline_model.APIResourceKey{
			ID:   pipeline_id,
			Type: "PIPELINE",
		},
		Relationship: "OWNER",
	}

	apiPipeline := pipeline_model.APIPipelineVersion{
		CodeSourceURL: url,
		Name:          computedName,
		PackageURL: &pipeline_model.APIURL{
			PipelineURL: url,
		},
		ResourceReferences: []*pipeline_model.APIResourceReference{&resourceReferences},
	}

	pipelineParams := pipeline_service.CreatePipelineVersionParams{
		Body:    &apiPipeline,
		Context: context,
	}

	resp, err := client.PipelineService.CreatePipelineVersion(&pipelineParams, nil)
	if err != nil {
		return fmt.Errorf("unable to create pipeline: %s", err)
	}

	d.Set("version", version)
	d.Set("version_id", resp.Payload.ID)

	return resourceKubeflowPipelinesPipelineRead(d, meta)
}

func resourceKubeflowPipelinesPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Pipeline
	context := meta.(*Meta).Context

	version := d.Get("version").(string)

	id := d.Id()

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

	versionId := resp.Payload.DefaultVersion.ID

	pipelineVersionParams := pipeline_service.GetPipelineVersionParams{
		VersionID: versionId,
		Context:   context,
	}

	respVersion, err := client.PipelineService.GetPipelineVersion(&pipelineVersionParams, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.Set("version_id", "")
		} else {
			return fmt.Errorf("unable to get pipeline version: %s", err)
		}
	} else {
		if strings.Contains(respVersion.Payload.Name, version) {
			d.Set("version_id", respVersion.Payload.ID)
		} else {
			d.Set("version_id", "")
		}
	}

	d.SetId(resp.Payload.ID)
	d.Set("name", resp.Payload.Name)
	d.Set("description", resp.Payload.Description)

	return nil
}

func resourceKubeflowPipelinesPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Pipeline
	context := meta.(*Meta).Context

	id := d.Id()

	pipelineParams := pipeline_service.DeletePipelineParams{
		ID:      id,
		Context: context,
	}

	_, err := client.PipelineService.DeletePipeline(&pipelineParams, nil)
	if err != nil {
		return fmt.Errorf("unable to delete pipeline: %s", err)
	}

	d.SetId("")

	return nil
}
