package tap_test

import (
	"fmt"
	"strings"
)

func ExampleNewParser() {
	r := strings.NewReader(`TAP version 14
1..3
ok 1 Squiggle
not ok 2 Wiggle
# Doesn't wiggle
not ok 3 Fliggle # TODO Fliggling not implemented yet
`)
	p, err := NewParser(r)
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

	for _, t := range(suite.Tests) {
		if t.Ok {
			fmt.Println("Test", t.Num, "ok")
		} else {
			fmt.Println("Test", t.Num, "not ok:", t.Diagnostic)
		}
	}

	//Output:
	//Test 1 ok
	//Test 2 not ok: Doesn't wiggle
	//
	//Test 3 not ok:
}
