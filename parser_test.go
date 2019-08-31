package rentparser

import (
	"testing"

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

	typ, err := p.ParseType(text)
	if err != nil {
		t.Errorf("type err=%s", err)
	}
	if typ != 2 {
		t.Errorf("bad type=%d", typ)
	}
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
				So(got,ShouldEqual,exp)
	})
}
