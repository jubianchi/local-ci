package configuration

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"local-ci/result"
	"local-ci/utils"
	"os"
)

type Job struct {
	name         string
	stage        string
	image        Image
	beforeScript Script
	afterScript  Script
	script       Script
	variables    Variables
	allowFailure bool
	retry        int
}

func (job Job) writeScript(writer io.Writer) {
	writer.Write([]byte("#!/bin/sh\n\n"))
	writer.Write([]byte("(\n"))
	writer.Write([]byte("set -e\n"))
	job.beforeScript.Write(writer)
	writer.Write([]byte("\n"))
	job.script.Write(writer)
	writer.Write([]byte(")\n"))
	writer.Write([]byte("STATUS=$?"))
	writer.Write([]byte("\n"))
	writer.Write([]byte("(\n"))
	writer.Write([]byte("set -e\n"))
	job.afterScript.Write(writer)
	writer.Write([]byte(")\n"))
	writer.Write([]byte("exit $STATUS\n"))
}

func (job Job) pullImage(client *client.Client) {
	reader, err := job.image.Pull(client)
	utils.CheckError(err)

	jsonmessage.DisplayJSONMessagesToStream(reader, utils.MakeOutStream(os.Stdout), nil)
	reader.Close()
}

func (job Job) configureContainer(containerConfiguration *container.Config) {
	job.image.ConfigureContainer(containerConfiguration)
	job.variables.ConfigureContainer(containerConfiguration)
}

func (job Job) AppendToPipeline(pipeline Pipeline) {
	pipeline.AppendJob(job.stage, job)
}

func (job Job) Run(client *client.Client, workingDirectory string) *result.Job {
	job.pullImage(client)

	os.MkdirAll(workingDirectory+"/.local-ci", 0755)

	entrypointPath := workingDirectory + "/.local-ci/entrypoint"
	writer, _ := os.Create(entrypointPath)

	job.writeScript(writer)

	containerConfiguration := container.Config{
		Cmd:        []string{"/bin/sh", "/runner"},
		WorkingDir: "/build",
	}
	job.configureContainer(&containerConfiguration)

	hostConfiguration := container.HostConfig{
		Binds: []string{
			entrypointPath + ":/runner",
			workingDirectory + ":/build",
		},
	}

	networkConfiguration := network.NetworkingConfig{}

	ct, createErr := client.ContainerCreate(context.Background(), &containerConfiguration, &hostConfiguration, &networkConfiguration, "")

	defer client.ContainerRemove(context.Background(), ct.ID, types.ContainerRemoveOptions{Force: true})

	utils.CheckError(createErr)

	attach, attachErr := client.ContainerAttach(
		context.Background(),
		ct.ID,
		types.ContainerAttachOptions{
			Stream: true,
			Stdin:  false,
			Stdout: true,
			Stderr: true,
		},
	)

	utils.CheckError(attachErr)

	startErr := client.ContainerStart(
		context.Background(),
		ct.ID,
		types.ContainerStartOptions{},
	)

	utils.CheckError(startErr)

	stdcopy.StdCopy(os.Stdout, os.Stderr, attach.Reader)

	bodyChan, errChan := client.ContainerWait(context.Background(), ct.ID, "")

	select {
	case body := <-bodyChan:
		if body.StatusCode == 0 {
			return result.CreateJobResult(job.name, result.PASSED)
		}

		if job.allowFailure {
			return result.CreateJobResult(job.name, result.ALLOWED_TO_FAIL)
		}

		return result.CreateJobResult(job.name, result.FAILED)
	case waitErr := <-errChan:
		panic(waitErr)
	}
}

func (job Job) Skip() *result.Job {
	return result.CreateJobResult(job.name, result.SKIPPED)
}

func CreateJob(name, stage string) *Job {
	return &Job{
		name:  name,
		stage: stage,
		variables: *CreateVariables(map[string]string{
			"CI_JOB_ID":     "1",
			"CI_JOB_MANUAL": "1",
			"CI_JOB_NAME":   name,
			"CI_JOB_STAGE":  stage,
			"CI_JOB_TOKEN":  "",
		}),
	}
}

func CreateJobWithDefault(name, stage string, image Image, beforeScript, afterScript Script, variables Variables) *Job {
	job := CreateJob(name, stage)

	job.image = image
	job.beforeScript = beforeScript
	job.afterScript = afterScript
	job.variables = *job.variables.MergeVariables(variables)

	return job
}

func CreateJobFromYaml(name string, yml map[interface{}]interface{}, factory func(string, string) *Job) *Job {
	stage := "test"

	if nil != yml["stage"] {
		stage = yml["stage"].(string)
	}

	if nil == factory {
		factory = CreateJob
	}

	job := factory(name, stage)

	if nil != yml["allow_failure"] {
		job.allowFailure = yml["allow_failure"].(bool)
	}

	if nil != yml["retry"] {
		job.retry = yml["retry"].(int)
	}

	if nil != yml["image"] {
		job.image = *CreateImageFromYaml(yml["image"])
	}

	if nil != yml["before_script"] {
		job.beforeScript = *CreateScriptFromYaml(yml["before_script"])
	}

	if nil != yml["after_script"] {
		job.afterScript = *CreateScriptFromYaml(yml["after_script"])
	}

	if nil != yml["script"] {
		job.script = *CreateScriptFromYaml(yml["script"])
	}

	if nil != yml["variables"] {
		job.variables = *job.variables.MergeVariables(*CreateVariablesFromYaml(yml["variables"].(map[interface{}]interface{})))
	}

	return job
}
