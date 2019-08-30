package rentparser

import (
	"fmt"

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

	return price.Parse(p.tomitaBin, p.typeConf, text)
}
