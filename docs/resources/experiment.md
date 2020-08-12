# kubeflowpipelines_experiment

The `kubeflowpipelines_experiment` creates a kubeflow pipelines experiment.

## Example Usage

```hcl-terraform
resource "kubeflowpipelines_experiment" "example" {
  name        = "example-experiment"
  description = "Some description"
}
```

## Argument Reference

The following arguments are supported:
### Required arguments
* `name` &mdash; (Required) Name of the experiment.
### Optional arguments
* `description` &mdash; (Optional) Description of the experiment.

## Attributes Reference

The following attribute is exported:

* `created_at` &mdash; The date and time of creation, formatted with RFC3339.