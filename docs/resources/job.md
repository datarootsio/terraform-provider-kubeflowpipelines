# kubeflowpipelines_job

The `kubeflowpipelines_job` creates a kubeflow pipelines job

## Example Usage

```hcl-terraform
resource "kubeflowpipelines_job" "example" {
	name            = "example"
	description     = "Description"
	service_account = "pipeline-runner"
	enabled         = true
	max_concurrency = 2
	no_catchup      = true
	experiment_id = kubeflowpipelines_experiment.example.id
	pipeline_spec {
		pipeline_version_id = kubeflowpipelines_pipeline.example.version_id
		parameters = {
			"key" = "value"
		}
	}
	trigger {
		cron_schedule {
            start_time = "2020-06-23T00:00:00Z"
            end_time   = "2030-06-23T00:00:00Z"
            cron       = "0 0/10 * * * ?"

		}
	}
}
```

## Argument Reference

The following arguments are supported:

### Required arguments
* `name` &mdash; (Required).
* `max_concurrency` &mdash; (Required). Maximum number of concurrent runs for this job. Value must be between `1` and `10` included.
* `pipeline_spec` &mdash; (Required). The pipeline specifications of this job.
* `trigger` &mdash; (Required). The trigger rules for this job. Contains either `cron_schedule` or `periodic_schedule` block.

### Optional arguments
* `description` &mdash; (Optional). A description for the run.
* `enabled` &mdash; (Optional, defaults to `true`). Is the job enabled ?
* `service_account` &mdash; (Optional). The service account to attach this job to.
* `experiment_id` &mdash; (Optional). The experiment that will store this jon.
* `no_catchup` &mdash; (Optional, defaults to `true`). Whether the job should catch up if behind schedule. If true, the job will only schedule the latest interval if behind schedule. If false, the job will catch up on each past interval.

### Nested Blocks

#### pipeline_spec

* `pipeline_version_id` &mdash; (Required). The pipeline version ID associated to this run.
* `parameters` &mdash; (Optional). A map of key values of parameters.

#### cron_schedule

* `start_time` &mdash; (Required). The start time in RFC3339 format.
* `end_time` &mdash; (Required). The end time in RFC3339 format.
* `cron` &mdash; (Required). A schedule represented with the cron format, with the [quartz format](http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html). Must include seconds.

#### periodic_schedule

* `start_time` &mdash; (Required). The start time in RFC3339 format.
* `end_time` &mdash; (Required). The end time in RFC3339 format.
* `interval_seconds` &mdash; (Required). The interval between each run, in seconds.

## Attributes Reference

The following attribute is exported:

* `created_at` &mdash; The date and time of creation, formatted with RFC3339.