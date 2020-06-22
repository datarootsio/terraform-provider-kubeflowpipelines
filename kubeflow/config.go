package kubeflow

import (
	"context"
	"fmt"
	"net/url"

	//	"github.com/Azure/go-autorest/autorest"
	//	"github.com/hashicorp/terraform/httpclient"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client"
)

const TerraformProviderUserAgent = "terraform-provider-kubeflow"

type Config struct {
	Token            string
	Host             string
	terraformVersion string
}

type Meta struct {
	Experiment *experiment_client.Experiment
	Context    context.Context
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

	return &meta, nil
}
