package kubeflowpipelines

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/job_client/job_service"
)

func TestAccResourceKubeflowPipelinesJob_basic(t *testing.T) {
	resourceName := "kubeflowpipelines_job.test"
	jobName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccResourceKubeflowPipelinesJobDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceKubeflowPipelinesJobBasic(jobName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", jobName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", jobName)),
				),
			},
		},
	})
}

func testAccResourceKubeflowPipelinesJobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubeflowpipelines_job" {
			continue
		}

		id := rs.Primary.Attributes["id"]

		client := testAccProvider.Meta().(*Meta).Job
		context := testAccProvider.Meta().(*Meta).Context

		jobParams := job_service.GetJobParams{
			ID:      id,
			Context: context,
		}

		_, err := client.JobService.GetJob(&jobParams, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil
			}
			return fmt.Errorf("unable to get job: %s", err)
		}
		return fmt.Errorf("Job still exists: %s", id)
	}
	return nil
}

func testAccResourceKubeflowPipelinesJobBasic(jobName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_experiment" "test" {
	name        = "%[1]s"
	description = "%[1]s"
}

resource "kubeflowpipelines_pipeline" "test" {
	name        = "%[1]s"
	description = "Description %[1]s"
	file_base64 = "YXBpVmVyc2lvbjogYXJnb3Byb2ouaW8vdjFhbHBoYTEKa2luZDogV29ya2Zsb3cKbWV0YWRhdGE6CiAgZ2VuZXJhdGVOYW1lOiBrdWJlZmxvdy1waXBlbGluZS1mb3ItdGVzdC0KICBhbm5vdGF0aW9uczoge3BpcGVsaW5lcy5rdWJlZmxvdy5vcmcva2ZwX3Nka192ZXJzaW9uOiAwLjUuMSwgcGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZV9jb21waWxhdGlvbl90aW1lOiAnMjAyMC0wNy0wNlQxMzozMzoxMi41MTM4NzInLAogICAgcGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZV9zcGVjOiAneyJkZXNjcmlwdGlvbiI6ICJUaGlzIHBpcGVsaW5lIGlzIGEgdGVzdCIsICJpbnB1dHMiOiBbeyJuYW1lIjogImtleSJ9LAogICAgICB7Im5hbWUiOiAiZm9vIn1dLCAibmFtZSI6ICJLdWJlZmxvdyBwaXBlbGluZSBmb3IgdGVzdCJ9J30KICBsYWJlbHM6IHtwaXBlbGluZXMua3ViZWZsb3cub3JnL2tmcF9zZGtfdmVyc2lvbjogMC41LjF9CnNwZWM6CiAgZW50cnlwb2ludDoga3ViZWZsb3ctcGlwZWxpbmUtZm9yLXRlc3QKICB0ZW1wbGF0ZXM6CiAgLSBuYW1lOiBrdWJlZmxvdy1waXBlbGluZS1mb3ItdGVzdAogICAgaW5wdXRzOgogICAgICBwYXJhbWV0ZXJzOgogICAgICAtIHtuYW1lOiBrZXl9CiAgICAgIC0ge25hbWU6IGZvb30KICAgIGRhZzoKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogc3RlcC0xCiAgICAgICAgdGVtcGxhdGU6IHN0ZXAtMQogICAgICAgIGFyZ3VtZW50czoKICAgICAgICAgIHBhcmFtZXRlcnM6CiAgICAgICAgICAtIHtuYW1lOiBrZXksIHZhbHVlOiAne3tpbnB1dHMucGFyYW1ldGVycy5rZXl9fSd9CiAgICAgICAgICAtIHtuYW1lOiBmb28sIHZhbHVlOiAne3tpbnB1dHMucGFyYW1ldGVycy5mb299fSd9CiAgICAgIC0gbmFtZTogc3RlcC0xLTIKICAgICAgICB0ZW1wbGF0ZTogc3RlcC0xLTIKICAgICAgICBkZXBlbmRlbmNpZXM6IFtzdGVwLTFdCiAgICAgICAgYXJndW1lbnRzOgogICAgICAgICAgcGFyYW1ldGVyczoKICAgICAgICAgIC0ge25hbWU6IGtleSwgdmFsdWU6ICd7e2lucHV0cy5wYXJhbWV0ZXJzLmtleX19J30KICAtIG5hbWU6IHN0ZXAtMQogICAgY29udGFpbmVyOgogICAgICBhcmdzOiBbLS1rZXksICd7e2lucHV0cy5wYXJhbWV0ZXJzLmtleX19JywgLS1mb28sICd7e2lucHV0cy5wYXJhbWV0ZXJzLmZvb319J10KICAgICAgY29tbWFuZDogW2VjaG8sICcnJ3BpcGVsaW5lJycnXQogICAgICBlbnY6CiAgICAgIC0ge25hbWU6IEVOVklST05FTUVOVCwgdmFsdWU6IHN1cGVyZW52fQogICAgICBpbWFnZTogaGVsbG8td29ybGQKICAgICAgaW1hZ2VQdWxsUG9saWN5OiBBbHdheXMKICAgIGlucHV0czoKICAgICAgcGFyYW1ldGVyczoKICAgICAgLSB7bmFtZToga2V5fQogICAgICAtIHtuYW1lOiBmb299CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxzOiB7cGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZS1zZGstdHlwZToga2ZwfQogIGFyZ3VtZW50czoKICAgIHBhcmFtZXRlcnM6CiAgICAtIHtuYW1lOiBrZXl9CiAgICAtIHtuYW1lOiBmb299CiAgc2VydmljZUFjY291bnROYW1lOiBwaXBlbGluZS1ydW5uZXIK"
	file_format = "yaml"
	version     = "v0.0.1"
}


resource "kubeflowpipelines_job" "test" {
	name            = "%[1]s"
	description     = "Description %[1]s"
	service_account = "pipeline-runner"
	enabled         = true
	max_concurrency = 2
	no_catchup      = true
	experiment_id = kubeflowpipelines_experiment.test.id
	pipeline_spec {
		pipeline_version_id = kubeflowpipelines_pipeline.test.version_id
		parameters = {
			"key" = "value"
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
`, jobName)
}
