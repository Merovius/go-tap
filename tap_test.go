package tap_test

import (
	"github.com/Merovius/go-tap"
	. "github.com/robertkrimen/terst"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	Terst(t)

	r := strings.NewReader(`TAP version 13
1..2
ok 1
not ok 2`)
	p, e := tap.NewParser(r)
	Is(e, nil, "No error parsing preamble")

	s, e := p.Suite()
	Is(e, nil, "No error parsing input")

	Is(len(s.Tests), 2, "Right number of tests")

	Is(s.Tests[0].Ok, true, "First test ok")
	Is(s.Tests[1].Ok, false, "Second test not ok")
}
