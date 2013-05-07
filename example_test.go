package tap_test

import (
	"fmt"
	"github.com/Merovius/go-tap"
	"strings"
)

func ExampleNewParser() {
	r := strings.NewReader(`TAP version 13
1..3
ok 1 Squiggle
not ok 2 Wiggle
# Doesn't wiggle
not ok 3 Fliggle # TODO Fliggling not implemented yet`)
	p, err := tap.NewParser(r)
	if err != nil {
		panic(err)
	}

	suite, err := p.Suite()
	if err != nil {
		panic(err)
	}
	if suite.Ok {
		fmt.Println("Everything ok")
		return
	}

	for _, t := range suite.Tests {
		fmt.Println(t)
	}

	// Output:
	// [ ok ] Squiggle
	// [fail] Wiggle # Doesn't wiggle
	// [todo] Fliggle # Fliggling not implemented yet
}
