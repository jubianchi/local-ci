package configuration

import (
	"io"
	"local-ci/utils"
)

type Script struct {
	commands []string
}

func CreateScriptFromYaml(yml interface{}) *Script {
	script := Script{}

	if command, ok := yml.(string); ok {
		script.commands = []string{command}
	}

	if commands, ok := yml.([]interface{}); ok {
		commands, _ := utils.InterfaceArrayToStringArray(commands)

		script.commands = commands
	}

	return &script
}

func (script Script) Write(writer io.Writer) (int, error) {
	str := ""

	for _, line := range script.commands {
		str += "echo '\033[33m> " + line + "'\033[39m\n"
		str += line + "\n"
	}

	return writer.Write([]byte(str))
}
