package kubeflowpipelines

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceKubeflowPipelinesExperiment_basic(t *testing.T) {
	resourceName := "data.kubeflowpipelines_experiment.test"
	resourceWithName := "data.kubeflowpipelines_experiment.test_name"
	experimentName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubeflowPipelinesExperimentBasic(experimentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", experimentName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", experimentName)),
					resource.TestCheckResourceAttr(resourceWithName, "name", experimentName),
					resource.TestCheckResourceAttr(resourceWithName, "description", fmt.Sprintf("Description %s", experimentName)),
				),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesExperimentBasic(experimentName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_experiment" "test" {
  name        = "%s"
  description = "Description %s"
}

data "kubeflowpipelines_experiment" "test" {
  id = kubeflowpipelines_experiment.test.id
}

data "kubeflowpipelines_experiment" "test_name" {
	name = "%s"
}
`, experimentName, experimentName, experimentName)
}
