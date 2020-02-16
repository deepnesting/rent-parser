package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPreValidRe(t *testing.T) {
	Convey("test prevalid", t, func() {
		for c, cc := range map[string]struct {
			matchFlat   bool
			matchSearch bool
		}{
			"ищу однушку":  {true, true},
			"сдам однушку": {true, true},
			"ищу соседа":   {true, true},
		} {
			Println(c)
			So(reFlat.MatchString(normalize(c)), ShouldEqual, cc.matchFlat)
		}
	})
}
