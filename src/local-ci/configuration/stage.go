package configuration

import (
	"github.com/docker/docker/client"
	"local-ci/result"
)

type Stage struct {
	name string
	jobs []Job
}

func (stage Stage) Empty() bool {
	return len(stage.jobs) == 0
}

func (stage *Stage) AppendJob(job Job) {
	stage.jobs = append(stage.jobs, job)
}

func (stage Stage) Run(client *client.Client, workingDirectory string) *result.Stage {
	stageResult := result.CreateStageResult(stage.name, result.PASSED)

	for _, job := range stage.jobs {
		jobResult := job.Run(client, workingDirectory)

		if jobResult.Status == result.FAILED {
			stageResult.Status = result.FAILED
		}

		stageResult.Jobs = append(stageResult.Jobs, *jobResult)
	}

	return stageResult
}

func (stage Stage) Skip() *result.Stage {
	stageResult := result.CreateStageResult(stage.name, result.SKIPPED)

	for _, job := range stage.jobs {
		stageResult.Jobs = append(stageResult.Jobs, *job.Skip())
	}

	return stageResult
}

func CreateStage(name string) *Stage {
	return &Stage{
		name: name,
	}
}
