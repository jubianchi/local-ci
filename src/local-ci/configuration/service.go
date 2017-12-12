package configuration

type Service struct {
	Name       string
	Alias      string
	Entrypoint []string
	Command    []string
}
