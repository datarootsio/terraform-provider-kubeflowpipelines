# kubeflowpipelines_pipeline

The `kubeflowpipelines_pipeline` creates a kubeflow pipelines pipeline.

To avoid deleting old related objects (runs, jobs, etc), changes in the pipeline version will upload a new version, not recreate a new pipeline object.

## Example Usage

```hcl-terraform
resource "kubeflowpipelines_pipeline" "example" {
	name        = "example"
	description = "Description"
	file_base64 = base64encode(file("${path.module}/pipeline.yaml"))
	file_format = "yaml"
	version     = "v0.0.1"
}

resource "kubeflowpipelines_pipeline" "example_url" {
	name    = "example-with-url"
	url     = "https://raw.githubusercontent.com/datarootsio/terraform-provider-kubeflowpipelines/master/tests/kubeflow_setup/pipeline.yaml"
	version = "v0.0.1"
}
```

## Argument Reference

The following arguments are supported:
### Required arguments
* `name` &mdash; (Required).

### Optional arguments

* `description` &mdash; (Optional, conflicts with `url`). A description for the pipeline. Descriptions are ignored by the API for URL based pipelines so it's currently conflicting
* `url` &mdash; (Optional, conflicts with `file_base64`). The URL containing a pipeline definition
* `file_base64` &mdash; (Optional, conflicts with `url`). A base64 encoded pipeline.
* `file_format` &mdash; (Optional, required with `file_base64`). The format of the pipeline. One of `zip`,`tar.gz`,`yaml`
* `version` &mdash; (Optional). The version of the pipeline.

## Attributes Reference

The following attribute is exported:

* `version_id` &mdash; The pipeline's version ID.