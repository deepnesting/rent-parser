package price

import (
	"testing"
)

func TestNormilize(t *testing.T) {
	var texts = map[string]struct {
		In  string
		Out string
	}{
		"1a": {
			In:  "1а",
			Out: "1 а",
		},
		"127000к плюс ку": {
			In:  "127000к плюс ку",
			Out: "127000 к плюс ку",
		},
		"13 000": {
			In:  "13 000",
			Out: "13000",
		},
	}
	for name, tt := range texts {
		if norm := splitNumbers(tt.In); norm != tt.Out {
			t.Errorf("%s: exp:%s got: %s", name, tt.Out, norm)
		}
	}

}
