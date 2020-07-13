package kubeflowpipelines

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceKubeflowPipelinesExperiment_basic(t *testing.T) {
	resourceName := "kubeflowpipelines_experiment.test"
	experimentName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccResourceKubeflowPipelinesExperimentDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceKubeflowPipelinesExperimentBasic(experimentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", experimentName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", experimentName)),
				),
			},
		},
	})
}

func testAccResourceKubeflowPipelinesExperimentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubeflowpipelines_experiment" {
			continue
		}

		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*Meta).Experiment

		resp, err := client.ExperimentService.ListExperiment(nil, nil)
		if err != nil {
			return fmt.Errorf("unable to get list of experiments: %s", name)
		}

		experimentFound := false

		for _, item := range resp.Payload.Experiments {
			if item.Name == name {
				experimentFound = true
				break
			}
		}

		if experimentFound {
			return fmt.Errorf("experiment %s still exists", name)
		}
	}

	return nil
}

func testAccResourceKubeflowPipelinesExperimentBasic(experimentName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_experiment" "test" {
  name        = "%s"
  description = "Description %s"
}
`, experimentName, experimentName)
}
