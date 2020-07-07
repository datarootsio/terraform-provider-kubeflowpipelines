package kubeflowpipelines

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceKubeflowPipelinesPipeline_basic(t *testing.T) {
	resourceName := "data.kubeflowpipelines_pipeline.test"
	pipelineName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubeflowPipelinesPipelineBasic(pipelineName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", pipelineName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", pipelineName)),
				),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesPipelineBasic(pipelineName string) string {
	return fmt.Sprintf(`
data "local_file" "pipeline_yaml" {
	filename = "${path.module}/pipeline.yaml"
}

resource "kubeflowpipelines_pipeline" "test" {
  name          = "%s"
  description   = "Description %s"
  file_base64   = data.local_file.pipeline_yaml.content_base64
}

data "kubeflowpipelines_pipeline" "test" {
  id = kubeflowpipelines_pipeline.test.id
}
`, pipelineName, pipelineName)
}
