package parser

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/deepnesting/rent-parser/src/tomita"
)

const FLAT_1 = 1
const FLAT_2 = 2
const FLAT_3 = 3
const FLAT_4 = 4
const STUDIO = 5
const ROOM = 6
const WRONG = 7

const (
	BuildTypeApart = 1
	BuildTypeRoom  = 2
	BuildTypeHouse = 3
)

type Type struct {
	Type     int
	Position int
	Sequence int
	Value    string
}

type XmlType struct {
	XMLName xml.Name `xml:"Type"`
	Value   string   `xml:"val,attr"`
}

type XmlWrong struct {
	XMLName xml.Name `xml:"Wrong"`
	Value   string   `xml:"val,attr"`
}

type XmlError struct {
	XMLName xml.Name `xml:"Error"`
	Value   string   `xml:"val,attr"`
}

type XmlFactRealty struct {
	XMLName   xml.Name  `xml:"FactRealty"`
	TypeList  []XmlType `xml:"Type"`
	FirstWord int       `xml:"fw,attr"`
	LastWord  int       `xml:"lw,attr"`
	Sequence  int       `xml:"sn,attr"`
}

type XmlFactRent struct {
	XMLName   xml.Name  `xml:"FactRent"`
	TypeList  []XmlType `xml:"Type"`
	FirstWord int       `xml:"fw,attr"`
	LastWord  int       `xml:"lw,attr"`
	Sequence  int       `xml:"sn,attr"`
}

type XmlFactNeighbor struct {
	XMLName   xml.Name  `xml:"FactNeighbor"`
	TypeList  []XmlType `xml:"Type"`
	FirstWord int       `xml:"fw,attr"`
	LastWord  int       `xml:"lw,attr"`
	Sequence  int       `xml:"sn,attr"`
}

type XmlFactWrong struct {
	XMLName   xml.Name   `xml:"FactWrong"`
	WrongList []XmlWrong `xml:"Wrong"`
	FirstWord int        `xml:"fw,attr"`
	LastWord  int        `xml:"lw,attr"`
	Sequence  int        `xml:"sn,attr"`
}

type XmlFactError struct {
	XMLName   xml.Name   `xml:"FactError"`
	ErrorList []XmlError `xml:"Error"`
}

type XmlFacts struct {
	XMLName          xml.Name          `xml:"facts"`
	FactRealtyList   []XmlFactRealty   `xml:"FactRealty"`
	FactRentList     []XmlFactRent     `xml:"FactRent"`
	FactNeighborList []XmlFactNeighbor `xml:"FactNeighbor"`
	FactWrongList    []XmlFactWrong    `xml:"FactWrong"`
	FactErrorList    []XmlFactError    `xml:"FactError"`
}

type XmlDocument struct {
	XMLName  xml.Name `xml:"document"`
	XMLFacts XmlFacts `xml:"facts"`
}

type XmlFdoObject struct {
	XMLName  xml.Name    `xml:"fdo_objects"`
	Document XmlDocument `xml:"document"`
}

func Parse(tomitaBin, confPath, text string) (int, int, string, error) {
	tom := tomita.NewTomita(tomitaBin, confPath)
	text = normalize(text)
	xml, err := tom.Parse(text)
	if err != nil {
		return 0, 0, "", err
	}
	t, r, value, err := getByXML(xml)
	return t, r, value, err
}

func ParseFacts(tomitaBin, confPath, text string) (*Facts, error) {
	tom := tomita.NewTomita(tomitaBin, confPath)
	text = normalize(text)
	xml, err := tom.Parse(text)
	if err != nil {
		return nil, err
	}
	return GetFacts(xml)
}

var (
	flats  = []string{"кв", "ком", "одн", "дву", "тр(ё|e)", "студи", "сосед", "подсел"}
	reFlat = regexp.MustCompile(fmt.Sprintf(`(?i).*(%s).*`, strings.Join(flats, "|")))
	//searchRe = regexp.MustCompile(`(?i)(сним|ищ)(у|ем)[\w]*[^\w](кварт|комн|одн(ок|ушк)|дву(хк|шк)|тр[её](хк|шк))`)
	searchRe = regexp.MustCompile(`(?i)(сним|ищ)(у|ем)[\w]*[^\w](кварт|комн|одн(ок|ушк)|дву(хк|шк)|тр[её](хк|шк))`)
)

func PreValid(rawText string) error {
	normalizedText := normalize(rawText)
	normalizedByteText := []byte(normalizedText)

	if !reFlat.Match(normalizedByteText) {
		return fmt.Errorf("%w: flat not found in text(norm: %s)", ErrNotValid, normalizedText)
	}

	return nil
}

func normalize(raw_text string) string {
	byte_text := []byte(strings.TrimSpace(raw_text))
	re := regexp.MustCompile(`\?\W{0,10}$`)
	if nil != re.Find(byte_text) {
		return ""
	}

	re2 := regexp.MustCompile(`(?i)(публиковать|варианты .* нашего сайта|правила темы|сайт)(.|\n)*`)
	byte_text = re2.ReplaceAll(byte_text, []byte(""))

	re3 := regexp.MustCompile(`(?i)http(s):(\w|\/|\.)*`)
	byte_text = re3.ReplaceAll(byte_text, []byte(""))

	byte_text = []byte(strings.Replace(string(byte_text), `\n`, "\n", -1))
	byte_text = []byte(strings.Replace(string(byte_text), "-", " ", -1))

	months := []string{"январ", "феврал", "март", "апрел", "ма(я|е|й)", "июн", "июл", "август", "сентябр", "октябр", "ноябр", "декабр"}

	tmp_re := regexp.MustCompile(fmt.Sprintf(`(?i)\d{0,2}[^0-9\.\!\?;]{0,3}(%s)[а-я]{0,4}`, strings.Join(months, "|")))

	byte_text = tmp_re.ReplaceAll(byte_text, []byte(``))

	flats := [...]string{
		`(?i)(\s|[0-9])кк(\s|\.)`,
		`(?i)(\s|[0-9])ккв(\s|\.)`,
		`(?i)(\s|[0-9])к\.\s{0,2}к(\s|\.)`,
		`(?i)(\s|[0-9])к\.\s{0,2}квартира(\s|\.)`,
		`(?i)(\s|[0-9])к\.\s{0,2}квартиру(\s|\.)`,
		`(?i)(\s|[0-9])к\.\s{0,2}кв(\s|\.)`,
		`(?i)(\s|[0-9])к\.\s{0,2}кварт(\s|\.)`,
		`(?i)(\s|[0-9])комн\.\s{0,2}кв(\s|\.)`,
		`(?i)(\s|[0-9])хкк(\s|\.)`,
		`(?i)(\s|[0-9])х\.\s{0,2}ком.*\b(\s|\.)`,
		`(?i)(\s|[0-9])хк\.\s{0,2}кв.*\b(\s|\.)`,
		`(?i)(\s|[0-9])х\.\s{0,2}к.*\b(\s|\.)кв.*\b(\s|\.)`,
		`(?i)(\s|[0-9])хк\.\s{0,2}к.*\b(\s|\.)кв.*\b(\s|\.)`,
	}

	for _, flat := range flats {
		tmp_re := regexp.MustCompile(flat)
		byte_text = tmp_re.ReplaceAll(byte_text, []byte("$1 комнатная квартира "))
	}

	re4 := regexp.MustCompile(`(?i)\d{1,3}\s{0,10}(кв(\.|\s){0,1}м(\.|\s){0,1}|м²|м(\.|\s))`)
	byte_text = re4.ReplaceAll(byte_text, []byte(" "))

	re5 := regexp.MustCompile(`(?i)([\d-=\+.\!?])([а-яеёa-z])`)
	byte_text = re5.ReplaceAll(byte_text, []byte("$1 $2"))

	re6 := regexp.MustCompile(`(?i)([а-яеёa-z])([\d-=\+.\!?])`)
	byte_text = re6.ReplaceAll(byte_text, []byte("$1 $2"))

	re7 := regexp.MustCompile(`(?i)\sквартир[а-яА-Яeё]*`)
	byte_text = re7.ReplaceAll(byte_text, []byte(" квартира "))

	re8 := regexp.MustCompile(`(?i)\sкомната[а-яА-Яeё]*`)
	byte_text = re8.ReplaceAll(byte_text, []byte(" комната "))

	text := []rune(string(byte_text))

	if len(text) > 500 {
		byte_text = []byte(string(text[:500]))
	}

	return string(byte_text)
}

var ErrNotValid = fmt.Errorf("not valid")

type Facts struct {
	RentFacts      []Fact
	RealtyFacts    []Fact
	NeighbourFacts []Fact
}

func (f Facts) SearchType() (int, error) {
	if len(f.NeighbourFacts) > 0 && len(f.RealtyFacts) == 0 && len(f.RentFacts) == 0 {
		return int(Neighbour), nil
	}
	if len(f.RentFacts) > 0 && len(f.RealtyFacts) == 0 && len(f.NeighbourFacts) == 0 {
		return int(Rent), nil
	}
	if len(f.RealtyFacts) > 0 && len(f.RentFacts) == 0 && len(f.NeighbourFacts) == 0 {
		return int(Realty), nil
	}
	return 0, ErrNotValid
}

func (f Facts) BuildType() (int, error) {
	if f.IsStudio() {
		return BuildTypeApart, nil
	}
	if f.IsSRoom() {
		return BuildTypeRoom, nil
	}
	return 0, fmt.Errorf("invalid")
}

func (f Facts) IsStudio() bool {
	for _, fact := range f.NeighbourFacts {
		if getTypeByString(fact.Value) == STUDIO {
			return true
		}
	}
	for _, fact := range f.RentFacts {
		if getTypeByString(fact.Value) == STUDIO {
			return true
		}
	}
	for _, fact := range f.RealtyFacts {
		if getTypeByString(fact.Value) == STUDIO {
			return true
		}
	}

	return false
}

func (f Facts) IsSRoom() bool {
	for _, fact := range f.NeighbourFacts {
		if getTypeByString(fact.Value) == ROOM {
			return true
		}
	}
	for _, fact := range f.RentFacts {
		if getTypeByString(fact.Value) == ROOM {
			return true
		}
	}
	for _, fact := range f.RealtyFacts {
		if getTypeByString(fact.Value) == ROOM {
			return true
		}
	}

	return false
}

func (f Facts) RoomCount() (int, error) {
	var vals []int
	for _, fact := range f.NeighbourFacts {
		if n := getTypeByString(fact.Value); n != 0 && n != WRONG && n != ROOM {
			vals = append(vals, n)
		}
	}
	for _, fact := range f.RentFacts {
		if n := getTypeByString(fact.Value); n != 0 && n != WRONG && n != ROOM {
			vals = append(vals, n)
		}
	}
	for _, fact := range f.RealtyFacts {
		if n := getTypeByString(fact.Value); n != 0 && n != WRONG && n != ROOM {
			vals = append(vals, n)
		}
	}

	if len(vals) > 0 {
		return vals[0], nil
	}

	return 0, ErrNotValid
}

type Fact struct {
	Value string
}

func GetFacts(input string) (*Facts, error) {
	var document XmlFdoObject
	err := xml.Unmarshal([]byte(input), &document)
	if err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}
	if len(document.Document.XMLFacts.FactErrorList) > 0 {
		return nil, fmt.Errorf("fact errors (%v): %w", document.Document.XMLFacts.FactErrorList, ErrNotValid)
	}

	var result = new(Facts)

	// rentlist
	for _, fact := range document.Document.XMLFacts.FactRentList {
		if len(fact.TypeList) == 0 {
			continue
		}
		result.RentFacts = append(result.RentFacts, Fact{
			Value: fact.TypeList[0].Value,
		})
	}
	for _, fact := range document.Document.XMLFacts.FactRealtyList {
		if len(fact.TypeList) == 0 {
			continue
		}
		result.RealtyFacts = append(result.RealtyFacts, Fact{
			Value: fact.TypeList[0].Value,
		})
	}
	for _, fact := range document.Document.XMLFacts.FactNeighborList {
		if len(fact.TypeList) == 0 {
			continue
		}
		result.NeighbourFacts = append(result.NeighbourFacts, Fact{
			Value: fact.TypeList[0].Value,
		})
	}
	return result, nil
}

func getByXML(xml_row string) (int, int, string, error) {
	if xml_row == "" {
		return 0, 0, "", fmt.Errorf("empty string: %w", ErrNotValid)
	}

	var document XmlFdoObject

	err := xml.Unmarshal([]byte(xml_row), &document)
	if err != nil {
		return 0, 0, "", fmt.Errorf("unmarshal xml: %w", err)
	}

	if len(document.Document.XMLFacts.FactErrorList) > 0 {
		return 0, 0, "", fmt.Errorf("fact errors (%v): %w", document.Document.XMLFacts.FactErrorList, ErrNotValid)
	}

	rent := Type{Type: -1, Position: 99, Sequence: 99}
	for _, fact := range document.Document.XMLFacts.FactRentList {
		if len(fact.TypeList) == 0 {
			continue
		}
		position := fact.FirstWord
		sequence := fact.Sequence
		rtype := getTypeByString(fact.TypeList[0].Value)

		if STUDIO == rtype {
			rent.Position = position
			rent.Sequence = sequence
			rent.Type = rtype
			rent.Value = fact.TypeList[0].Value
			break
		}

		if sequence == rent.Sequence && position < rent.Position {
			rent.Position = position
			rent.Sequence = sequence
			rent.Type = rtype
			rent.Value = fact.TypeList[0].Value
			continue
		}

		if sequence < rent.Sequence {
			rent.Position = position
			rent.Sequence = sequence
			rent.Type = rtype
			rent.Value = fact.TypeList[0].Value
			continue
		}
	}

	neighbor := Type{Type: -1, Position: 99, Sequence: 99}
	for _, fact := range document.Document.XMLFacts.FactNeighborList {
		if len(fact.TypeList) == 0 {
			continue
		}

		position := fact.FirstWord
		sequence := fact.Sequence
		rtype := getTypeByString(fact.TypeList[0].Value)

		if STUDIO == rtype {
			neighbor.Position = position
			neighbor.Sequence = sequence
			neighbor.Type = rtype

			break
		}

		if sequence == neighbor.Sequence && position < neighbor.Position {
			neighbor.Position = position
			neighbor.Sequence = sequence
			neighbor.Type = rtype

			continue
		}

		if sequence < neighbor.Sequence {
			neighbor.Position = position
			neighbor.Sequence = sequence
			neighbor.Type = rtype

			continue
		}
	}

	realty := Type{Type: -1, Position: 99, Sequence: 99}
	for _, fact := range document.Document.XMLFacts.FactRealtyList {
		if len(fact.TypeList) == 0 {
			continue
		}

		position := fact.FirstWord
		sequence := fact.Sequence
		rtype := getTypeByString(fact.TypeList[0].Value)

		if STUDIO == rtype {
			realty.Position = position
			realty.Sequence = sequence
			realty.Type = rtype

			break
		}

		if sequence == realty.Sequence && position < realty.Position {
			realty.Position = position
			realty.Sequence = sequence
			realty.Type = rtype

			continue
		}

		if sequence < realty.Sequence {
			realty.Position = position
			realty.Sequence = sequence
			realty.Type = rtype

			continue
		}
	}

	wrong := Type{Type: -1, Position: 99, Sequence: 99}
	for _, fact := range document.Document.XMLFacts.FactWrongList {
		if len(fact.WrongList) == 0 {
			continue
		}

		position := fact.FirstWord
		sequence := fact.Sequence
		rtype := getTypeByString(fact.WrongList[0].Value)

		if sequence == wrong.Sequence && position < wrong.Position {
			wrong.Position = position
			wrong.Sequence = sequence
			wrong.Type = rtype
			wrong.Value = fact.WrongList[0].Value

			continue
		}

		if sequence < wrong.Sequence {
			wrong.Position = position
			wrong.Sequence = sequence
			wrong.Type = rtype
			wrong.Value = fact.WrongList[0].Value

			continue
		}
	}

	switch true {
	case -1 != wrong.Type &&
		(wrong.Sequence < rent.Sequence || (wrong.Sequence == rent.Sequence && wrong.Position < rent.Position)) &&
		(wrong.Sequence < realty.Sequence || (wrong.Sequence == realty.Sequence && wrong.Position < realty.Position)):
		return 0, 0, "", fmt.Errorf("some wrong (value: %s)", wrong.Value)
	case -1 != rent.Type:
		return int(Rent), rent.Type, rent.Value, nil
	case -1 != neighbor.Type:
		return int(Neighbour), neighbor.Type, "", nil
	case -1 != wrong.Type:
		return 0, 0, "", fmt.Errorf("some wrong 2")
	case -1 != realty.Type:
		return int(Realty), realty.Type, "", nil
	default:
		return 0, 0, "", fmt.Errorf("unexpected error: %w", ErrNotValid)
	}
}

type OfferType int

const (
	Realty OfferType = iota + 1
	Rent
	Neighbour
)

func getTypeByString(raw_text string) int {

	if "" == raw_text {
		return WRONG
	}

	text := strings.ToLower(raw_text)

	if -1 != strings.Index(text, "студи") {
		return STUDIO
	}

	re := regexp.MustCompile(`(^|\W)комнаты($|\W)`)

	if nil != re.Find([]byte(text)) {
		return ROOM
	}

	if -1 != strings.Index(text, "1") {
		return FLAT_1
	}

	if -1 != strings.Index(text, "2") {
		return FLAT_2
	}

	if -1 != strings.Index(text, "3") {
		return FLAT_3
	}

	re2 := regexp.MustCompile(`(([^\d,\.!?]|^)[4-9]\D{0,30}квартир|четыр\Sх|много)|(квартир\D{0,3}1\D.{0,10}комнатн)`)

	if nil != re2.Find([]byte(text)) {
		return FLAT_4
	}

	re3 := regexp.MustCompile(`(^|\W)квартир\W{1,4}($|\W)`)

	if nil != re3.Find([]byte(text)) {
		return FLAT_1
	}

	re4 := regexp.MustCompile(`(^|\W)комнат`)

	if nil != re4.Find([]byte(text)) {
		return ROOM
	}

	return WRONG
}
