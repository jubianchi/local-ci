package main

import (
	//"context"
	"os"
	"path/filepath"

	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/pkg/stdcopy"
	"github.com/urfave/cli"

	"local-ci/app"
	"local-ci/configuration"
	//"./container"
	"fmt"
	"local-ci/utils"
	//"github.com/docker/docker/pkg/stdcopy"
	//"context"
	//"github.com/docker/docker/api/types/container"
	//"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	//"github.com/docker/docker/pkg/jsonmessage"
	//"github.com/docker/docker/pkg/stdcopy"
	"local-ci/result"
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

				switch pipelineResult.Status {
				case result.PASSED:
					fmt.Print("\033[1;32mPipeline passed:\033[0;39m\n")
					break
				case result.ALLOWED_TO_FAIL:
					fmt.Print("\033[1;33mPipeline passed with allowed failures:\033[0;39m\n")
					break
				case result.FAILED:
					fmt.Print("\n\033[1;31mPipeline failed:\033[0;39m\n")
					break
				}

				for _, stageResult := range pipelineResult.Stages {
					switch stageResult.Status {
					case result.PASSED:
						fmt.Printf("  \033[32mStage \033[1m%s\033[0;32m passed:\033[0;39m\n    ", stageResult.Name)
						break
					case result.ALLOWED_TO_FAIL:
						fmt.Printf("  \033[33mStage \033[1m%s\033[0;33m passed with allowed failures:\033[0;39m\n    ", stageResult.Name)
						break
					case result.FAILED:
						fmt.Printf("  \033[31mStage \033[1m%s\033[0;31m failed:\033[0;39m\n    ", stageResult.Name)
						break
					case result.SKIPPED:
						fmt.Printf("  \033[37mStage \033[1m%s\033[0;37m was skipped:\033[0;39m\n    ", stageResult.Name)
						break
					}

					for index, jobResult := range stageResult.Jobs {
						switch jobResult.Status {
						case result.PASSED:
							fmt.Printf("\033[1;32m%s (passed)\033[0;39m", jobResult.Name)
							break
						case result.ALLOWED_TO_FAIL:
							fmt.Printf("\033[1;33m%s (passed with allowed failures)\033[0;39m", jobResult.Name)
							break
						case result.FAILED:
							fmt.Printf("\033[1;31m%s (failed)\033[0;39m", jobResult.Name)
							break
						case result.SKIPPED:
							fmt.Printf("\033[1;37m%s (skipped)\033[0;39m", jobResult.Name)
							break
						}

						if index < len(stageResult.Jobs)-1 {
							fmt.Print(", ")
						} else {
							fmt.Print("\n")
						}
					}
				}

				return nil
			},
		},
	}

	application.Run(os.Args)
}
