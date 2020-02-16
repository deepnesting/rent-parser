package tomita

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Tomita struct {
	logger *log.Logger

	bin     string
	config  string
	verbose bool
}

func NewTomita(bin, config string, verbose bool) *Tomita {
	return &Tomita{
		logger:  log.New(os.Stdout, "tomita ", log.LstdFlags|log.Lshortfile),
		bin:     bin,
		config:  config,
		verbose: verbose,
	}
}

func (tomita Tomita) Parse(text string) (string, error) {
	command := exec.Command(tomita.bin, tomita.config)
	var Stdout bytes.Buffer
	command.Stdin = strings.NewReader(text)
	command.Stdout = &Stdout

	var stdErr bytes.Buffer

	if tomita.verbose {
		command.Stderr = &stdErr
	} else {
		command.Stderr = ioutil.Discard
	}

	err := command.Run()
	if tomita.verbose {
		tomita.logger.Printf("command %s %s done\nstderr: %s\nstdout: %s", tomita.bin, tomita.config, stdErr.String, Stdout.String)
	}
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, stdErr.String())
	}

	return Stdout.String(), nil
}
