package result

import (
	"os"
)

type Stage struct {
	Status int
	Name   string
	Jobs   []Job
}

func (stage Stage) Report() {
	Colors[stage.Status]("  Stage %s", Bold(stage.Name))
	Colors[stage.Status](" %s:\n    ", Labels[stage.Status])

	for index, job := range stage.Jobs {
		job.Report()

		if index < len(stage.Jobs)-1 {
			os.Stdout.WriteString(", ")
		} else {
			os.Stdout.WriteString("\n")
		}
	}
}

func CreateStageResult(name string, status int) *Stage {
	return &Stage{Status: status, Name: name, Jobs: make([]Job, 0)}
}
