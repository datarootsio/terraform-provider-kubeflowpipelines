package kubeflowpipelines

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceKubeflowPipelinesExperiment_basic(t *testing.T) {
	resourceName := "data.kubeflowpipelines_experiment.test"
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
`, experimentName, experimentName)
}

func TestAccDataSourceKubeflowPipelinesExperiment_name(t *testing.T) {
	resourceName := "data.kubeflowpipelines_experiment.test_name"
	experimentName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubeflowPipelinesExperimentName(experimentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", experimentName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", experimentName)),
				),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesExperimentName(experimentName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_experiment" "test_name" {
  name        = "%s"
  description = "Description %s"
}

data "kubeflowpipelines_experiment" "test_name" {
  name = kubeflowpipelines_experiment.test_name.name
}
`, experimentName, experimentName)
}

func TestAccDataSourceKubeflowPipelinesExperiment_missing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceKubeflowPipelinesExperimentMissingName(),
				ExpectError: regexp.MustCompile("errors during refresh: unable to get experiment: non_existant"),
			},
			{
				Config:      testAccDataSourceKubeflowPipelinesExperimentMissingID(),
				ExpectError: regexp.MustCompile("errors during refresh: unable to get experiment: id-abc"),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesExperimentMissingID() string {
	return fmt.Sprintf(`
data "kubeflowpipelines_experiment" "test_missing" {
  id = "id-abc"
}
`)
}

func testAccDataSourceKubeflowPipelinesExperimentMissingName() string {
	return fmt.Sprintf(`
data "kubeflowpipelines_experiment" "test_missing_with_name" {
  name = "non_existant"
}
`)
}
