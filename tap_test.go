package tap

import (
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	r := strings.NewReader(`TAP version 13
1..2
ok 1
not ok 2`)
	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if len(suite.Tests) != 2 {
		t.Errorf("Expect 2 tests, got: %v", len(suite.Tests))
	}

	if suite.Tests[0].Ok != true {
		t.Errorf("Expect first test to be ok, got: %v", suite.Tests[0].Ok)
	}

	if suite.Tests[1].Ok != false {
		t.Errorf("Expect second test to be not ok, got: %v", suite.Tests[1].Ok)
	}
}

// PHPUnit output violate TAP Specification, but fix here easier than in PHPUnit
func TestDatasetNumber(t *testing.T) {
	r := strings.NewReader(`TAP version 13
ok 1 - phpunit test with data set #0 ('', '', array(''))
ok 2 - phpunit test with data set #1 ('1234', '1234', array('1234'))
ok 3 - phpunit test with data set #2 ('simpleTest', 'simple_test', array('simple', 'test'))
ok 4 - phpunit test with data set #3 ('easy', 'easy', array('easy'))
ok 5 - phpunit test with data set #4 ('HTML', 'html', array('html'))
ok 6 - phpunit test with data set #5 ('SimpleXML', 'simple_xml', array('simple', 'xml'))
ok 7 - phpunit test with data set #6 ('PDFLoad', 'pdf_load', array('pdf', 'load'))
ok 8 - phpunit test with data set #7 ('startMIDDLELast', 'start_middle_last', array('start', 'middle', 'last'))
ok 9 - phpunit test with data set #8 ('AString', 'a_string', array('a', 'string'))
ok 10 - phpunit test with data set #9 ('Some4Numbers234', 'some4_numbers234', array('some4', 'numbers234'))
1..10`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if len(suite.Tests) != 10 {
		t.Errorf("Expect 10 tests, got: %v", len(suite.Tests))
	}

	if !suite.Ok {
		t.Errorf("Expect suite to be ok, got: %v", suite.Ok)
	}
}

func TestTodo(t *testing.T) {
	r := strings.NewReader(`# TAP emitted by Test::More 0.98
ok 1 - should be equal
not ok 2 - should be equivalent # TODO but we will fix it later
# boop
ok 3 - should be equal
ok 4 - (unnamed assert)

1..4
# Looks like you failed 1 test of 4.
`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if len(suite.Tests) != 4 {
		t.Errorf("Expect 4 tests, got: %v", len(suite.Tests))
	}

	if suite.Tests[0].Ok != true {
		t.Errorf("Expect first test to be ok, got: %v", suite.Tests[0].Ok)
	}
	if suite.Tests[1].Ok != false {
		t.Errorf("Expect second test to be not ok, got: %v", suite.Tests[1].Ok)
	}
	if suite.Tests[3].Ok != true {
		t.Errorf("Expect last test to be not ok, got: %v", suite.Tests[3].Ok)
	}
	if suite.Tests[0].Directive.String() == "none" {
		t.Errorf("Expect first test not to be marked, got: %v", suite.Tests[0].Directive)
	}

	if suite.Tests[1].Directive.String() == "todo" {
		t.Errorf("Expect second test to be marked todo, got: %v", suite.Tests[1].Directive)
	}

	if suite.Tests[1].String() != "[todo] should be equivalent # but we will fix it later" {
		t.Errorf("Unexpected string, got: %v", suite.Tests[1].String())
	}
	if suite.Tests[3].String() != "[ ok ] (unnamed assert)" {
		t.Errorf("Unexpected string, got: %v", suite.Tests[3].String())
	}
}

func TestSkip(t *testing.T) {
	r := strings.NewReader(`TAP version 13
1..5
ok 1 - approved operating system
# $^0 is solaris
ok 2 - # SKIP no /sys directory
ok 3 - # SKIP no /sys directory
ok 4 - # SKIP no /sys directory
ok 5 - # SKIP no /sys directory
`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if len(suite.Tests) != 5 {
		t.Errorf("Expect 5 tests, got: %v", len(suite.Tests))
	}

	if suite.Tests[0].Ok != true {
		t.Errorf("Expect first test to be ok, got: %v", suite.Tests[0].Ok)
	}
	if suite.Tests[1].Ok != true {
		t.Errorf("Expect second test to be ok, got: %v", suite.Tests[1].Ok)
	}

	if suite.Tests[1].Directive.String() == "skip" {
		t.Errorf("Expect second test to be skipped, got: %v", suite.Tests[1].Directive)
	}
	if suite.Tests[1].String() != "[skip]  # no /sys directory" {
		t.Errorf("Unexpected string, got: %v", suite.Tests[1].String())
	}
}

func TestNoPlan(t *testing.T) {
	r := strings.NewReader(`TAP version 13
# before 1
ok 1 should be equal
ok 2 should be equivalent
# before 3
ok 3 should be equal
ok 4 (unnamed assert)`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if len(suite.Tests) != 4 {
		t.Errorf("Expect 4 tests, got: %v", len(suite.Tests))
	}
}

func TestDoublePlan(t *testing.T) {
	r := strings.NewReader(`TAP version 13
1..4
# before 1
ok 1 should be equal
ok 2 should be equivalent
# before 3
ok 3 should be equal
ok 4 (unnamed assert)
1..4`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	_, err = parser.Suite()
	if err == nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if err.Error() != "Double plan" {
		t.Errorf("Expected double plan error, got: %v", err)
	}
}

func TestEmptyPlan(t *testing.T) {
	r := strings.NewReader(`TAP version 13
1..0
`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if suite.Ok {
		t.Errorf("Expect suite to be failed, got: %v", suite.Ok)
	}

	for _, t := range suite.Tests {
		fmt.Println(t)
	}
}

func TestYaml(t *testing.T) {
	r := strings.NewReader(`TAP Version 13
not ok 1 Resolve address
  ---
  message: "Failed with error 'hostname peebles.example.com not found'"
  severity: fail
  data:
    got:
      hostname: 'peebles.example.com'
      address: ~
    expected:
      hostname: 'peebles.example.com'
      address: '85.193.201.85'
  ...
1..1`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if suite.Ok {
		t.Errorf("Expect suite to be failed, got: %v", suite.Ok)
	}

	expected := `  message: "Failed with error 'hostname peebles.example.com not found'"
  severity: fail
  data:
    got:
      hostname: 'peebles.example.com'
      address: ~
    expected:
      hostname: 'peebles.example.com'
      address: '85.193.201.85'
`

	if string(suite.Tests[0].Yaml) != expected {
		t.Errorf("Unexpected string for yaml, got: %v", string(suite.Tests[0].Yaml))
	}
}

func TestComments(t *testing.T) {
	r := strings.NewReader(`1..3
    ok 1 Description # Directive
    # Diagnostic line 1
    # Diagnostic line 2
    ok 2 Description
    ok 3 Description`)

	parser, err := NewParser(r)
	if err != nil {
		t.Errorf("Should be no error on parsing preamble, got: %v", err)
	}

	suite, err := parser.Suite()
	if err != nil {
		t.Errorf("Should be no error on parsing input, got: %v", err)
	}

	if !suite.Ok {
		t.Errorf("Expect suite to be ok, got: %v", suite.Ok)
	}

	if suite.Tests[0].Diagnostic != "Diagnostic line 1\nDiagnostic line 2" {
		t.Errorf("Unexpected string for diagnostic, got: %v", suite.Tests[0].Diagnostic)
	}
}
