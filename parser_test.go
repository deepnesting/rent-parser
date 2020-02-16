package rentparser

import (
	"testing"

	parsetype "github.com/deepnesting/rent-parser/src/parser/type"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {

	t.Skip("skipped, change path to bin and configs")

	p := New("/Users/god/uapi/bin/tomita-mac",
		"/Users/god/uapi/docs/tomita/type/config.proto",
		"/Users/god/uapi/docs/tomita/price/config.proto")
	text := "сдаю двушку за 30 тыс в месяц"
	price, err := p.ParsePrice(text)
	if err != nil {
		t.Errorf("err=%s", err)
	}
	if price != 30000 {
		t.Errorf("price=%d", price)
	}

	Convey("test parsing", t, func() {
		for tc, c := range map[string]struct {
			In    string
			Facts parsetype.Facts
		}{
			"однушка":      {"сдам однокомнатную квартиру", parsetype.Facts{RealtyFacts: []parsetype.Fact{{"1 КВАРТИРА"}}, RentFacts: []parsetype.Fact{{"1 КВАРТИРА"}}}},
			"сдам однушку": {"сдам однушку", parsetype.Facts{RealtyFacts: []parsetype.Fact{{"1 КВАРТИРА"}}, RentFacts: []parsetype.Fact{{"1 КВАРТИРА"}}}},
			"ищу соседа":   {"ищу соседа в комнату", parsetype.Facts{RealtyFacts: []parsetype.Fact{{"КОМНАТА"}}, NeighbourFacts: []parsetype.Fact{{"КОМНАТА"}}}},
		} {
			Println(tc)
			facts, err := p.ParseFacts(c.In)
			So(err, ShouldBeNil)
			So(*facts, ShouldResemble, c.Facts)
		}
	})
}

func TestFixConfigPath(t *testing.T) {
	Convey("test fix", t, func() {
		inp := `TTextMinerConfig {

			Dictionary = "/Users/god/uapi/docs/tomita/price/dict.gzt";

			Output = {`

		got := fixFile(inp, "/allo/privet/dict.gzt")

		exp := `TTextMinerConfig {

			Dictionary = "/allo/privet/dict.gzt";

			Output = {`
		So(got, ShouldEqual, exp)
	})
}
