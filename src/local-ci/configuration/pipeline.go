package configuration

import (
	"fmt"
	"github.com/docker/docker/client"
	"local-ci/app"
	"local-ci/result"
	"local-ci/utils"
)

type Pipeline struct {
	image        Image
	variables    Variables
	services     []Service
	beforeScript Script
	afterScript  Script
	stages       []Stage
}

func (pipeline *Pipeline) AppendJob(stage string, job Job) error {
	for index, current := range pipeline.stages {
		if current.name == stage {
			pipeline.stages[index].AppendJob(job)

			return nil
		}
	}

	return fmt.Errorf("Stage '%s' does not exist", stage)
}

func (pipeline *Pipeline) DefineVariable(name string, value string) {
	variables := make(map[string]string)
	variables[name] = value

	pipeline.variables = *pipeline.variables.MergeVariables(*CreateVariables(variables))
}

func (pipeline Pipeline) CreateJob(name, stage string) *Job {
	return CreateJobWithDefault(name, stage, pipeline.image, pipeline.beforeScript, pipeline.afterScript, pipeline.variables)
}

func (pipeline Pipeline) Run(client *client.Client, workingDirectory string) *result.Pipeline {
	pipelineResult := result.CreatePipelineResult(result.PASSED)

	for _, stage := range pipeline.stages {
		if stage.Empty() {
			continue
		}

		if pipelineResult.Status == result.FAILED {
			pipelineResult.Stages = append(pipelineResult.Stages, *stage.Skip())

			continue
		}

		stageResult := stage.Run(client, workingDirectory)

		if stageResult.Status == result.FAILED {
			pipelineResult.Status = result.FAILED
		}

		pipelineResult.Stages = append(pipelineResult.Stages, *stageResult)
	}

	return pipelineResult
}

func CreatePipeline() *Pipeline {
	return &Pipeline{
		variables: *CreateVariables(map[string]string{
			"CI":                        "1",
			"CI_PIPELINE_ID":            "-1",
			"CI_PIPELINE_TRIGGERED":     "0",
			"CI_PIPELINE_SOURCE":        "external",
			"CI_DISPOSABLE_ENVIRONMENT": "1",
			"CI_RUNNER_DESCRIPTION":     app.Name + " version " + app.Version,
			"CI_RUNNER_ID":              "-1",
			"CI_RUNNER_TAGS":            app.Name,
		}),
	}
}

func CreatePipelineFromYaml(yml map[interface{}]interface{}) *Pipeline {
	pipeline := CreatePipeline()

	if nil != yml["before_script"] {
		pipeline.beforeScript = *CreateScriptFromYaml(yml["before_script"])
	}

	if nil != yml["after_script"] {
		pipeline.afterScript = *CreateScriptFromYaml(yml["after_script"])
	}

	if nil != yml["variables"] {
		pipeline.variables = *pipeline.variables.MergeVariables(*CreateVariablesFromYaml(yml["variables"].(map[interface{}]interface{})))
	}

	if nil != yml["image"] {
		pipeline.image = *CreateImageFromYaml(yml["image"])
	}

	if nil != yml["services"] {
		if services, ok := yml["services"].([]interface{}); ok {
			pipeline.services = make([]Service, len(services))

			for index, service := range services {
				if name, stringOk := service.(string); stringOk {
					pipeline.services[index] = Service{Name: name}
				}

				if desc, mapOk := service.(map[interface{}]interface{}); mapOk {
					pipeline.services[index] = Service{Name: desc["name"].(string)}

					if nil != desc["alias"] {
						pipeline.services[index].Alias = desc["alias"].(string)
					}

					if nil != desc["entrypoint"] {
						entrypoint, _ := utils.InterfaceArrayToStringArray(desc["entrypoint"].([]interface{}))

						pipeline.services[index].Entrypoint = entrypoint
					}

					if nil != desc["command"] {
						command, _ := utils.InterfaceArrayToStringArray(desc["command"].([]interface{}))

						pipeline.services[index].Command = command
					}
				}
			}
		}
	}

	if nil != yml["stages"] {
		stages, _ := utils.InterfaceArrayToStringArray(yml["stages"].([]interface{}))
		pipeline.stages = make([]Stage, len(stages))

		for index, stage := range stages {
			pipeline.stages[index] = *CreateStage(stage)
		}
	} else {
		pipeline.stages = []Stage{*CreateStage("test")}
	}

	return pipeline
}
