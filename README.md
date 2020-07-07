# terraform-provider-kubeflow-pipelines
Terraform provider for Kubeflow Pipelines API

[![tests](https://github.com/datarootsio/terraform-provider-kubeflowpipelines/workflows/tests/badge.svg?branch=master)](https://github.com/datarootsio/terraform-provider-kubeflowpipelines/actions?query=workflow%3Atests)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14

## Usage

```hcl
provider "kubeflowpipelines" {
  host = "http://localhost:8080"
}

data "kubeflowpipelines_experiment" "my_experiment" {
  name = "MyExperiment"
}

resource "kubeflowpipelines_pipeline" "my_pipeline" {
  name        = "MyPipeline"
  description = "This is an example pipeline"
  file_base64 = base64encode(file("${path.module}/pipeline.yaml"))
  file_format = "yaml"
  version     = "v0.0.4"
}

resource "kubeflowpipelines_job" "my_job" {
  name            = "MyJob"
  description     = "An example job"
  service_account = "pipeline-runner"
  enabled         = true
  max_concurrency = 2
  no_catchup      = true
  experiment_id = kubeflowpipelines_experiment.my_experiment.id
  pipeline_spec {
    pipeline_id = kubeflowpipelines_pipeline.my_pipeline.id
    parameters = {
      "foo" = "bar"
      "bar" = "baz"
    }
  }
  trigger {
    cron_schedule {
      start_time = "2020-06-23T00:00:00Z"
      end_time   = "2030-06-23T00:00:00Z"
      cron       = "0 10 * * *"
    }
  }
}

resource "kubeflowpipelines_run" "my_run" {
  name            = "MyRun"
  description     = "An example run"
  service_account = "pipeline-runner"
  experiment_id = kubeflowpipelines_experiment.my_experiment.id
  pipeline_spec {
    pipeline_version_id = kubeflowpipelines_pipeline.my_pipeline.version_id
    parameters = {
      "foo" = "bar"
      "bar" = "baz"
    }
  }
}
```
