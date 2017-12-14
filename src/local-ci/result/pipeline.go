package result

type Pipeline struct {
	Status int
	Stages []Stage
}

func (pipeline Pipeline) Report() {
	Colors[pipeline.Status]("Pipeline %s:\n", Labels[pipeline.Status])

	for _, stage := range pipeline.Stages {
		stage.Report()
	}
}

func CreatePipelineResult(status int) *Pipeline {
	return &Pipeline{Status: status, Stages: make([]Stage, 0)}
}
