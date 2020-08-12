# kubeflowpipelines_run

The `kubeflowpipelines_run` creates a kubeflow pipelines run

## Example Usage

```hcl-terraform
resource "kubeflowpipelines_run" "example" {
    name            = "example"
    description     = "Description"
    service_account = "pipeline-runner"
    experiment_id   = kubeflowpipelines_experiment.example.id
    pipeline_spec {
        pipeline_version_id = kubeflowpipelines_pipeline.example.version_id
        parameters = {
            "key" = "value"
        }
    }
}
```

## Argument Reference

The following arguments are supported:
### Required arguments
* `name` &mdash; (Required).
* `pipeline_spec` &mdash; (Required). The pipeline specification of this run.

### Optional arguments

* `description` &mdash; (Optional). A description for the run.
* `service_account` &mdash; (Optional). The service account to attach this run to.
* `experiment_id` &mdash; (Optional). The experiment that will store this run.

### Nested Blocks

#### pipeline_spec

* `pipeline_version_id` &mdash; (Required). The pipeline version ID associated to this run.
* `parameters` &mdash; (Optional). A map of key values of parameters.

## Attributes Reference

The following attribute is exported:

* `created_at` &mdash; The date and time of creation, formatted with RFC3339.