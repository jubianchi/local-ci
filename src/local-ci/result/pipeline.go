package result

type Pipeline struct {
	Status int
	Stages []Stage
}

func CreatePipelineResult(status int) *Pipeline {
	return &Pipeline{Status: status, Stages: make([]Stage, 0)}
}
