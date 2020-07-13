package kubeflowpipelines

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceKubeflowPipelinesExperiment_basic(t *testing.T) {
	resourceName := "kubeflowpipelines_experiment.test"
	experimentName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func testAccResourceKubeflowPipelinesExperimentBasic(experimentName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_experiment" "test" {
  name        = "%s"
  description = "Description %s"
}

data "kubeflowpipelines_experiment" "test" {
  id = kubeflowpipelines_experiment.test.id
}
`, experimentName, experimentName)
}

func TestAccResourceKubeflowPipelinesExperiment_name(t *testing.T) {
	resourceName := "data.kubeflowpipelines_experiment.test_name"
	experimentName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceKubeflowPipelinesExperimentName(experimentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", experimentName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", experimentName)),
				),
			},
		},
	})
}

func testAccResourceKubeflowPipelinesExperimentName(experimentName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_experiment" "test_name" {
  name        = "%s"
  description = "Description %s"
}
`, experimentName, experimentName)
}

func testAccResourceKubeflowPipelinesExperimentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubeflowpipelines_experiment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		description := rs.Primary.Attributes["description"]

		client := testAccProvider.Meta().(*Meta).Experiment
		ctx := testAccProvider.Meta().(*Meta).StopContext

		resp, err := client.ExperimentService.ListExperiment(nil, nil)
		if err != nil {
			return fmt.Errorf("unable to get list of experiments: %s", name)
		}

		experimentFound := false

		for _, item := range resp.Payload.Experiments {
			if item.Name == name {
				d.SetId(item.ID)
				d.Set("name", item.Name)
				d.Set("description", item.Description)
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