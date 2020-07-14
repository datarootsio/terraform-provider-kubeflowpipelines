package kubeflowpipelines

import (
	"fmt"
	"strings"
	"time"

	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client/experiment_service"
)

func readExperiment(meta interface{}, id string, name string) (string, string, string, string, error) {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	if id == "" {
		resp, err := client.ExperimentService.ListExperiment(nil, nil)
		if err != nil {
			return "", "", "", "", fmt.Errorf("unable to get list of experiments: %s", name)
		}

		experimentFound := false

		for _, item := range resp.Payload.Experiments {
			if item.Name == name {
				return item.ID, item.Name, item.Description, "", nil
			}
		}

		if !experimentFound {
			return "", "", "", "", fmt.Errorf("unable to get experiment: %s", name)
		}
	} else {
		experimentParams := experiment_service.GetExperimentParams{
			ID:      id,
			Context: context,
		}

		resp, err := client.ExperimentService.GetExperiment(&experimentParams, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return "", "", "", "", nil

			}
			return "", "", "", "", fmt.Errorf("unable to get experiment: %s", id)
		}

		return resp.Payload.ID, resp.Payload.Name, resp.Payload.Description, time.Time(resp.Payload.CreatedAt).Format(time.RFC3339), nil
	}

	return "", "", "", "", fmt.Errorf("unable to get experiment")

}
