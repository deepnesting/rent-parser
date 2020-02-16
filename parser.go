package rentparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

type BuildType int

const (
	BuildTypeUnknown BuildType = iota
	Apartament
	Room
	House
)

type RoomCount int

const (
	RoomCountUnknown RoomCount = iota
	RoomCountOne
	RoomCountTwo
	RoomCountThree
	RoomCountFour
	RoomCountFive
	RoomCountStudio
)

type SearchType int

const (
	SearchTypeUnknown = iota
	SearchNester
	SearchNest
	SearchNeighbour
)

var (
	simpleNeighbourRe    = regexp.MustCompile(`ищу.*соседа`)
	simpleSearchNesterRe = regexp.MustCompile(`(сдам|сдаю)`)
	simpleSearchNestRe   = regexp.MustCompile(`(сниму)`)
)

func (p *Parser) ParseSearchType(text string) (SearchType, string, error) {
	if err := parsetype.PreValid(text); err != nil {
		return 0, "", fmt.Errorf("not valid: %w", err)
	}

	var reSearchType SearchType
	var value string

	if simpleNeighbourRe.MatchString(strings.ToLower(text)) {
		reSearchType = SearchNeighbour
		value = "neighbour re"
	}
	if simpleSearchNestRe.MatchString(strings.ToLower(text)) {
		reSearchType = SearchNest
		value = "nest re"
	}
	if simpleSearchNesterRe.MatchString(strings.ToLower(text)) {
		reSearchType = SearchNester
		value = "nester re"
	}

	typ, _, v, err := parsetype.Parse(p.tomitaBin, p.typeConf, text)
	if typ == 0 && reSearchType != SearchTypeUnknown {
		return reSearchType, value, nil
	}
	if err != nil {
		return 0, "", err
	}
	return SearchType(typ), v, nil
}

func (p *Parser) ParseRoomCount(text string) (RoomCount, error) {
	if err := parsetype.PreValid(text); err != nil {
		return 0, fmt.Errorf("not valid: %w", err)
	}
	_, roomCount, _, err := parsetype.Parse(p.tomitaBin, p.typeConf, text)
	if err != nil {
		return 0, err
	}
	return RoomCount(roomCount), nil
}

func (p *Parser) ParseFacts(text string) (*parsetype.Facts, error) {
	if err := parsetype.PreValid(text); err != nil {
		return nil, fmt.Errorf("not valid: %w", err)
	}
	return parsetype.ParseFacts(p.tomitaBin, p.typeConf, text)
}

func (p *Parser) ParsePrice(text string) (int, error) {
	if err := parsetype.PreValid(text); err != nil {
		return 0, fmt.Errorf("not valid: %w", err)
	}
	return price.Parse(p.tomitaBin, p.typeConf, text)
}

// fixFile change internal file config paths for tomita binary can access this files
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
		fullDir = filepath.Join(dir, filepath.Dir(configPath))
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
