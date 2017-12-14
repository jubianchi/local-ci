package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"local-ci/app"
	"local-ci/configuration"
	"local-ci/utils"
	"os"
	"path/filepath"
)

func main() {
	application := app.CreateApp()

	application.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a GitLab CI pipeline",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: ".gitlab-ci.yml",
				},
			},
			Action: func(c *cli.Context) error {
				gitlabCiYml, err := filepath.Abs(c.String("file"))
				utils.CheckError(err)

				workingDirectory := filepath.Dir(gitlabCiYml)

				fmt.Printf("Working directory: %s\n", workingDirectory)
				fmt.Printf("Using YAML: %s\n", gitlabCiYml)

				ciYml := utils.ReadYaml(gitlabCiYml)
				pipeline := *configuration.CreatePipelineFromYaml(ciYml)
				pipeline.DefineVariable("CI_CONFIG_PATH", filepath.Base(gitlabCiYml))
				pipeline.DefineVariable("CI_PROJECT_NAME", filepath.Base(workingDirectory))

				delete(ciYml, "before_script")
				delete(ciYml, "after_script")
				delete(ciYml, "variables")
				delete(ciYml, "services")
				delete(ciYml, "image")
				delete(ciYml, "cache")
				delete(ciYml, "stages")

				for name, jobYml := range ciYml {
					job := configuration.CreateJobFromYaml(name.(string), jobYml.(map[interface{}]interface{}), pipeline.CreateJob)
					job.AppendToPipeline(pipeline)
				}

				docker, err := client.NewEnvClient()
				utils.CheckError(err)

				pipelineResult := pipeline.Run(docker, workingDirectory)
				pipelineResult.Report()

				return nil
			},
		},
	}

	application.Run(os.Args)
}
