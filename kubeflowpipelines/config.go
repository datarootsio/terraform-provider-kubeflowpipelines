package kubeflowpipelines

import (
	"context"
	"fmt"
	"net/url"

	//	"github.com/Azure/go-autorest/autorest"
	//	"github.com/hashicorp/terraform/httpclient"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/job_client"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_client"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/pipeline_upload_client"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/run_client"
)

const TerraformProviderUserAgent = "terraform-provider-kubeflowpipelines"

type Config struct {
	Token            string
	Host             string
	terraformVersion string
}

type Meta struct {
	Experiment     *experiment_client.Experiment
	Pipeline       *pipeline_client.Pipeline
	PipelineUpload *pipeline_upload_client.PipelineUpload
	Job            *job_client.Job
	Run            *run_client.Run
	Context        context.Context
}

func (c *Config) Client() (*Meta, error) {
	u, err := url.Parse(c.Host)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %s", err)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return c.createClients(u.Host, u.Scheme)
}

func (c *Config) createClients(host string, scheme string) (*Meta, error) {
	meta := Meta{}

	experimentTransport := experiment_client.DefaultTransportConfig()
	experimentTransport.Host = host
	experimentTransport.Schemes = []string{scheme}

	meta.Experiment = experiment_client.NewHTTPClientWithConfig(nil, experimentTransport)

	pipelineTransport := pipeline_client.DefaultTransportConfig()
	pipelineTransport.Host = host
	pipelineTransport.Schemes = []string{scheme}

	meta.Pipeline = pipeline_client.NewHTTPClientWithConfig(nil, pipelineTransport)

	pipelineUploadTransport := pipeline_upload_client.DefaultTransportConfig()
	pipelineUploadTransport.Host = host
	pipelineUploadTransport.Schemes = []string{scheme}

	meta.PipelineUpload = pipeline_upload_client.NewHTTPClientWithConfig(nil, pipelineUploadTransport)

	jobTransport := job_client.DefaultTransportConfig()
	jobTransport.Host = host
	jobTransport.Schemes = []string{scheme}

	meta.Job = job_client.NewHTTPClientWithConfig(nil, jobTransport)

	runTransport := run_client.DefaultTransportConfig()
	runTransport.Host = host
	runTransport.Schemes = []string{scheme}

	meta.Run = run_client.NewHTTPClientWithConfig(nil, runTransport)

	return &meta, nil
}
