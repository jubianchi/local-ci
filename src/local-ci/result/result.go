package result

import "github.com/fatih/color"

const (
	PASSED int = iota
	FAILED
	ALLOWED_TO_FAIL
	SKIPPED
)

var Bold = color.New(color.Bold).SprintFunc()

var Colors = map[int]func(string, ...interface{}){
	0: color.New(color.FgGreen).PrintfFunc(),
	1: color.New(color.FgRed).PrintfFunc(),
	2: color.New(color.FgYellow).PrintfFunc(),
	3: color.New(color.FgWhite).PrintfFunc(),
}

var Labels = map[int]string{
	0: "passed",
	1: "failed",
	2: "passed with allowed failures",
	3: "was skipped",
}

type Result interface {
	Report()
}
