package kubeflowpipelines

import (
	"fmt"
	"strings"
	"time"

	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client/experiment_service"
)

type ExperimentResult struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
}

func readExperiment(meta interface{}, id string, name string) (ExperimentResult, error) {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	if id == "" {
		resp, err := client.ExperimentService.ListExperiment(nil, nil)
		if err != nil {
			return ExperimentResult{}, fmt.Errorf("unable to get list of experiments: %s", name)
		}

		experimentFound := false

		for _, item := range resp.Payload.Experiments {
			if item.Name == name {

				return ExperimentResult{item.ID, item.Name, item.Description, ""}, nil
			}
		}

		if !experimentFound {
			return ExperimentResult{}, fmt.Errorf("unable to get experiment: %s", name)
		}
	} else {
		experimentParams := experiment_service.GetExperimentParams{
			ID:      id,
			Context: context,
		}

		resp, err := client.ExperimentService.GetExperiment(&experimentParams, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return ExperimentResult{}, nil

			}
			return ExperimentResult{}, fmt.Errorf("unable to get experiment: %s", id)
		}

		return ExperimentResult{resp.Payload.ID, resp.Payload.Name, resp.Payload.Description, time.Time(resp.Payload.CreatedAt).Format(time.RFC3339)}, nil
	}
	return ExperimentResult{}, fmt.Errorf("unable to get experiment")
}
