package kubeflowpipelines

import (
	"fmt"
	"regexp"
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
resource "kubeflowpipelines_pipeline" "test" {
  name        = "%s"
  description = "Description %s"
  file_base64 = "YXBpVmVyc2lvbjogYXJnb3Byb2ouaW8vdjFhbHBoYTEKa2luZDogV29ya2Zsb3cKbWV0YWRhdGE6CiAgZ2VuZXJhdGVOYW1lOiBrdWJlZmxvdy1waXBlbGluZS1mb3ItdGVzdC0KICBhbm5vdGF0aW9uczoge3BpcGVsaW5lcy5rdWJlZmxvdy5vcmcva2ZwX3Nka192ZXJzaW9uOiAwLjUuMSwgcGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZV9jb21waWxhdGlvbl90aW1lOiAnMjAyMC0wNy0wNlQxMzozMzoxMi41MTM4NzInLAogICAgcGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZV9zcGVjOiAneyJkZXNjcmlwdGlvbiI6ICJUaGlzIHBpcGVsaW5lIGlzIGEgdGVzdCIsICJpbnB1dHMiOiBbeyJuYW1lIjogImtleSJ9LAogICAgICB7Im5hbWUiOiAiZm9vIn1dLCAibmFtZSI6ICJLdWJlZmxvdyBwaXBlbGluZSBmb3IgdGVzdCJ9J30KICBsYWJlbHM6IHtwaXBlbGluZXMua3ViZWZsb3cub3JnL2tmcF9zZGtfdmVyc2lvbjogMC41LjF9CnNwZWM6CiAgZW50cnlwb2ludDoga3ViZWZsb3ctcGlwZWxpbmUtZm9yLXRlc3QKICB0ZW1wbGF0ZXM6CiAgLSBuYW1lOiBrdWJlZmxvdy1waXBlbGluZS1mb3ItdGVzdAogICAgaW5wdXRzOgogICAgICBwYXJhbWV0ZXJzOgogICAgICAtIHtuYW1lOiBrZXl9CiAgICAgIC0ge25hbWU6IGZvb30KICAgIGRhZzoKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogc3RlcC0xCiAgICAgICAgdGVtcGxhdGU6IHN0ZXAtMQogICAgICAgIGFyZ3VtZW50czoKICAgICAgICAgIHBhcmFtZXRlcnM6CiAgICAgICAgICAtIHtuYW1lOiBrZXksIHZhbHVlOiAne3tpbnB1dHMucGFyYW1ldGVycy5rZXl9fSd9CiAgICAgICAgICAtIHtuYW1lOiBmb28sIHZhbHVlOiAne3tpbnB1dHMucGFyYW1ldGVycy5mb299fSd9CiAgICAgIC0gbmFtZTogc3RlcC0xLTIKICAgICAgICB0ZW1wbGF0ZTogc3RlcC0xLTIKICAgICAgICBkZXBlbmRlbmNpZXM6IFtzdGVwLTFdCiAgICAgICAgYXJndW1lbnRzOgogICAgICAgICAgcGFyYW1ldGVyczoKICAgICAgICAgIC0ge25hbWU6IGtleSwgdmFsdWU6ICd7e2lucHV0cy5wYXJhbWV0ZXJzLmtleX19J30KICAtIG5hbWU6IHN0ZXAtMQogICAgY29udGFpbmVyOgogICAgICBhcmdzOiBbLS1rZXksICd7e2lucHV0cy5wYXJhbWV0ZXJzLmtleX19JywgLS1mb28sICd7e2lucHV0cy5wYXJhbWV0ZXJzLmZvb319J10KICAgICAgY29tbWFuZDogW2VjaG8sICcnJ3BpcGVsaW5lJycnXQogICAgICBlbnY6CiAgICAgIC0ge25hbWU6IEVOVklST05FTUVOVCwgdmFsdWU6IHN1cGVyZW52fQogICAgICBpbWFnZTogaGVsbG8td29ybGQKICAgICAgaW1hZ2VQdWxsUG9saWN5OiBBbHdheXMKICAgIGlucHV0czoKICAgICAgcGFyYW1ldGVyczoKICAgICAgLSB7bmFtZToga2V5fQogICAgICAtIHtuYW1lOiBmb299CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxzOiB7cGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZS1zZGstdHlwZToga2ZwfQogIGFyZ3VtZW50czoKICAgIHBhcmFtZXRlcnM6CiAgICAtIHtuYW1lOiBrZXl9CiAgICAtIHtuYW1lOiBmb299CiAgc2VydmljZUFjY291bnROYW1lOiBwaXBlbGluZS1ydW5uZXIK"
  file_format = "yaml"
  version     = "v0.0.1"
}

data "kubeflowpipelines_pipeline" "test" {
  id = kubeflowpipelines_pipeline.test.id
}
`, pipelineName, pipelineName)
}

func TestAccDataSourceKubeflowPipelinesPipeline_name(t *testing.T) {
	resourceName := "data.kubeflowpipelines_pipeline.test_name"
	pipelineName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubeflowPipelinesPipelineName(pipelineName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", pipelineName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", pipelineName)),
				),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesPipelineName(pipelineName string) string {
	return fmt.Sprintf(`
resource "kubeflowpipelines_pipeline" "test_name" {
  name        = "%s"
  description = "Description %s"
  file_base64 = "YXBpVmVyc2lvbjogYXJnb3Byb2ouaW8vdjFhbHBoYTEKa2luZDogV29ya2Zsb3cKbWV0YWRhdGE6CiAgZ2VuZXJhdGVOYW1lOiBrdWJlZmxvdy1waXBlbGluZS1mb3ItdGVzdC0KICBhbm5vdGF0aW9uczoge3BpcGVsaW5lcy5rdWJlZmxvdy5vcmcva2ZwX3Nka192ZXJzaW9uOiAwLjUuMSwgcGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZV9jb21waWxhdGlvbl90aW1lOiAnMjAyMC0wNy0wNlQxMzozMzoxMi41MTM4NzInLAogICAgcGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZV9zcGVjOiAneyJkZXNjcmlwdGlvbiI6ICJUaGlzIHBpcGVsaW5lIGlzIGEgdGVzdCIsICJpbnB1dHMiOiBbeyJuYW1lIjogImtleSJ9LAogICAgICB7Im5hbWUiOiAiZm9vIn1dLCAibmFtZSI6ICJLdWJlZmxvdyBwaXBlbGluZSBmb3IgdGVzdCJ9J30KICBsYWJlbHM6IHtwaXBlbGluZXMua3ViZWZsb3cub3JnL2tmcF9zZGtfdmVyc2lvbjogMC41LjF9CnNwZWM6CiAgZW50cnlwb2ludDoga3ViZWZsb3ctcGlwZWxpbmUtZm9yLXRlc3QKICB0ZW1wbGF0ZXM6CiAgLSBuYW1lOiBrdWJlZmxvdy1waXBlbGluZS1mb3ItdGVzdAogICAgaW5wdXRzOgogICAgICBwYXJhbWV0ZXJzOgogICAgICAtIHtuYW1lOiBrZXl9CiAgICAgIC0ge25hbWU6IGZvb30KICAgIGRhZzoKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogc3RlcC0xCiAgICAgICAgdGVtcGxhdGU6IHN0ZXAtMQogICAgICAgIGFyZ3VtZW50czoKICAgICAgICAgIHBhcmFtZXRlcnM6CiAgICAgICAgICAtIHtuYW1lOiBrZXksIHZhbHVlOiAne3tpbnB1dHMucGFyYW1ldGVycy5rZXl9fSd9CiAgICAgICAgICAtIHtuYW1lOiBmb28sIHZhbHVlOiAne3tpbnB1dHMucGFyYW1ldGVycy5mb299fSd9CiAgICAgIC0gbmFtZTogc3RlcC0xLTIKICAgICAgICB0ZW1wbGF0ZTogc3RlcC0xLTIKICAgICAgICBkZXBlbmRlbmNpZXM6IFtzdGVwLTFdCiAgICAgICAgYXJndW1lbnRzOgogICAgICAgICAgcGFyYW1ldGVyczoKICAgICAgICAgIC0ge25hbWU6IGtleSwgdmFsdWU6ICd7e2lucHV0cy5wYXJhbWV0ZXJzLmtleX19J30KICAtIG5hbWU6IHN0ZXAtMQogICAgY29udGFpbmVyOgogICAgICBhcmdzOiBbLS1rZXksICd7e2lucHV0cy5wYXJhbWV0ZXJzLmtleX19JywgLS1mb28sICd7e2lucHV0cy5wYXJhbWV0ZXJzLmZvb319J10KICAgICAgY29tbWFuZDogW2VjaG8sICcnJ3BpcGVsaW5lJycnXQogICAgICBlbnY6CiAgICAgIC0ge25hbWU6IEVOVklST05FTUVOVCwgdmFsdWU6IHN1cGVyZW52fQogICAgICBpbWFnZTogaGVsbG8td29ybGQKICAgICAgaW1hZ2VQdWxsUG9saWN5OiBBbHdheXMKICAgIGlucHV0czoKICAgICAgcGFyYW1ldGVyczoKICAgICAgLSB7bmFtZToga2V5fQogICAgICAtIHtuYW1lOiBmb299CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxzOiB7cGlwZWxpbmVzLmt1YmVmbG93Lm9yZy9waXBlbGluZS1zZGstdHlwZToga2ZwfQogIGFyZ3VtZW50czoKICAgIHBhcmFtZXRlcnM6CiAgICAtIHtuYW1lOiBrZXl9CiAgICAtIHtuYW1lOiBmb299CiAgc2VydmljZUFjY291bnROYW1lOiBwaXBlbGluZS1ydW5uZXIK"
  file_format = "yaml"
  version     = "v0.0.1"
}

data "kubeflowpipelines_pipeline" "test_name" {
  name = kubeflowpipelines_pipeline.test_name.name
}
`, pipelineName, pipelineName)
}

func TestAccDataSourceKubeflowPipelinesPipeline_missingid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceKubeflowPipelinesPipelineMissingID(),
				ExpectError: regexp.MustCompile("errors during refresh: unable to get pipeline: id-abc"),
			},
		},
	})
}

func TestAccDataSourceKubeflowPipelinesPipeline_missingname(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceKubeflowPipelinesPipelineMissingName(),
				ExpectError: regexp.MustCompile("errors during refresh: unable to get pipeline: non_existent"),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesPipelineMissingID() string {
	return `
data "kubeflowpipelines_pipeline" "test_missing" {
  id = "id-abc"
}
`
}

func testAccDataSourceKubeflowPipelinesPipelineMissingName() string {
	return `
data "kubeflowpipelines_pipeline" "test_missing_with_name" {
  name = "non_existent"
}
`
}
