package result

type Stage struct {
	Status int
	Name   string
	Jobs   []Job
}

func CreateStageResult(name string, status int) *Stage {
	return &Stage{Status: status, Name: name, Jobs: make([]Job, 0)}
}
