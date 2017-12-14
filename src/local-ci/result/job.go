package result

type Job struct {
	Status int
	Name   string
}

func (job Job) Report() {
	Colors[job.Status]("%s", Bold(job.Name))
	Colors[job.Status](" %s", Labels[job.Status])
}

func CreateJobResult(name string, status int) *Job {
	return &Job{Status: status, Name: name}
}
