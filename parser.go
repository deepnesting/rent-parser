package rentparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/deepnesting/rent-parser/src/parser/price"
	parsetype "github.com/deepnesting/rent-parser/src/parser/type"
)

type Parser struct {
	tomitaBin string
	typeConf  string
	priceConf string
}

func New(bin, typeConf, priceConf string) *Parser {
	return &Parser{
		tomitaBin: bin,
		typeConf:  typeConf,
		priceConf: priceConf,
	}
}

func (p *Parser) ParseType(text string) (int, error) {
	if !parsetype.PreValid(text) {
		return -1, fmt.Errorf("not valid")
	}

	return parsetype.Parse(p.tomitaBin, p.typeConf, text)
}

func (p *Parser) ParsePrice(text string) (int, error) {
	if !parsetype.PreValid(text) {
		return -1, fmt.Errorf("not valid")
	}

	return price.Parse(p.tomitaBin, p.priceConf, text)
}

func fixFile(inp, fullFilePath string) string {
	lines := strings.Split(inp, "\n")
	var result []string
	for _, l := range lines {
		if !strings.Contains(l, "Dictionary") {
			result = append(result, l)
			continue
		}
		fixed := ""
		for _, letter := range []rune(l) {
			if letter != '"' {
				fixed += string(letter)
			} else {
				fixed += `"` + fullFilePath + `";`
				break
			}
		}
		result = append(result, fixed)
	}
	return strings.Join(result, "\n")
}

func FixConfigPathToGzt(configPath string) error {
	var fullDir string
	if filepath.IsAbs(configPath) {
		fullDir = filepath.Dir(configPath)
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fullDir = filepath.Join(dir, configPath)
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	fullGztPath := filepath.Join(fullDir, "dict.gzt")
	res := fixFile(string(data), fullGztPath)

	err = ioutil.WriteFile(configPath, []byte(res), os.ModePerm)
	return err
}
