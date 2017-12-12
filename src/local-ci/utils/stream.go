package utils

import (
	"github.com/docker/docker/pkg/term"
	"os"
)

type Out struct {
	stream *os.File
}

func (out *Out) FD() uintptr {
	return out.stream.Fd()
}

func (out *Out) IsTerminal() bool {
	return term.IsTerminal(out.FD())
}

func (out *Out) Write(p []byte) (int, error) {
	return out.stream.Write(p)
}

func MakeOutStream(file *os.File) *Out {
	return &Out{stream: file}
}
