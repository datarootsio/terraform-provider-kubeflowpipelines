# kubeflowpipelines_experiment

The `kubeflowpipelines_experiment` fetches an existing kubeflow pipelines experiment.

## Example Usage

```hcl-terraform
data "kubeflowpipelines_pipeline" "example" {
  name = "example-experiment"
}
```

## Argument Reference

The following arguments are supported:

* `id` &mdash; (Optional) ID of the experiment (mutually exclusive with `name`).
* `name` &mdash; (Optional) Name of the experiment (mutually exclusive with `id`).

## Attributes Reference

The following attribute is exported:

* `description` &mdash; The experiment's description.