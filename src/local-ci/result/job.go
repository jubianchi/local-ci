package result

type Job struct {
	Status int
	Name   string
}

func CreateJobResult(name string, status int) *Job {
	return &Job{Status: status, Name: name}
}
