# terraform-provider-kubeflow-pipelines
Terraform provider for Kubeflow Pipelines API

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/datarootsio/terraform-provider-kubeflowpipelines)](https://registry.terraform.io/providers/datarootsio/kubeflowpipelines/latest) [![maintained by dataroots](https://img.shields.io/badge/maintained%20by-dataroots-%2300b189)](https://dataroots.io) [![tests](https://github.com/datarootsio/terraform-provider-kubeflowpipelines/workflows/tests/badge.svg?branch=master)](https://github.com/datarootsio/terraform-provider-kubeflowpipelines/actions?query=workflow%3Atests) [![Go Report Card](https://goreportcard.com/badge/github.com/datarootsio/terraform-provider-kubeflowpipelines)](https://goreportcard.com/report/github.com/datarootsio/terraform-provider-kubeflowpipelines) [![Coverage Status](https://coveralls.io/repos/github/datarootsio/terraform-provider-kubeflowpipelines/badge.svg)](https://coveralls.io/github/datarootsio/terraform-provider-kubeflowpipelines) ![GitHub](https://img.shields.io/github/license/datarootsio/terraform-provider-kubeflowpipelines)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14

## Usage

[Documentation available on the Terraform Registry](https://registry.terraform.io/providers/datarootsio/kubeflowpipelines/latest/docs)

Developed and validated against kubeflow 1.0.2

## Contributing

Contributions to this repository are very welcome! Found a bug or do you have a suggestion? Please open an issue. Do you know how to fix it? Pull requests are welcome as well! To get you started faster, a Makefile is provided.

Make sure to install [Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html), [Go](https://golang.org/doc/install) (for automated testing) and Make (optional, if you want to use the Makefile) on your computer. Install [tflint](https://github.com/terraform-linters/tflint) to be able to run the linting.

* Format your code: `make fmt`
* Run tests: `make test`
* Run acceptance tests: `make testacc`. This creates resources on your Kubernetes cluster, use with caution. We use [k3s](https://k3s.io/) in the CICD pipelines, to start from a fresh environment each time.

## License

MIT license. Please see [LICENSE](LICENSE) for details.