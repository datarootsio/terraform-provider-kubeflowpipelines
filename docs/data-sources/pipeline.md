# kubeflowpipelines_pipeline

The `kubeflowpipelines_pipeline` fetches an existing kubeflow pipelines pipeline.

## Example Usage

```hcl-terraform
data "kubeflowpipelines_pipeline" "example" {
  name = "example-pipeline"
}
```

## Argument Reference

The following arguments are supported:

* `id` &mdash; (Optional) ID of the pipeline (mutually exclusive with `name`).
* `name` &mdash; (Optional) Name of the pipeline (mutually exclusive with `id`).

## Attributes Reference

The following attribute is exported:

* `description` &mdash; The pipeline's description.