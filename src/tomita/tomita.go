package tomita

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Tomita struct {
	bin    string
	config string
}

func NewTomita(bin string, config string) *Tomita {
	p := new(Tomita)

	p.bin = bin
	p.config = config

	return p
}

func (tomita Tomita) Parse(text string) (string, error) {
	command := exec.Command(tomita.bin, tomita.config)
	var Stdout bytes.Buffer
	command.Stdin = strings.NewReader(text)
	command.Stdout = &Stdout

	// сюда попадает всякая дичь, типо дебага
	command.Stderr = ioutil.Discard

	err := command.Run()
	if err != nil {
		return "", err
	}

	return Stdout.String(), nil
}
