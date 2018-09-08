// Package tap implements a basic parser for the Test Anything Protocol
package tap

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	planRE        = regexp.MustCompile(`^1..(\d+)$`)
	versionRE     = regexp.MustCompile(`^TAP version (\d+)$`)
	testlineRE    = regexp.MustCompile(`^(not )?ok(?:\s+(\d+)\s*)?(?:\s*([\D\S][^#]+?)\s*)?(?i:#\s+(todo|skip)(?:\s+(.*))?)?$`)
	diagnosticsRE = regexp.MustCompile(`^\s*#\s+(.*)$`)
	yamlStartRE   = regexp.MustCompile(`\s*---\s*$`)
	yamlEndRE     = regexp.MustCompile(`\s*...\s*$`)
)

// A TAP-Directive (Todo/Skip)
//
// Deprecated: Project is unmaintained.
type Directive int

const (
	None Directive = iota // No directive given
	Todo                  // Testpoint is a TODO
	Skip                  // Testpoint was skipped
)

// Deprecated: Project is unmaintained.
func (d Directive) String() string {
	switch d {
	case None:
		return "None"
	case Todo:
		return "TODO"
	case Skip:
		return "SKIP"
	}
	return ""
}

// A single TAP-Testline
//
// Deprecated: Project is unmaintained.
type Testline struct {
	Ok          bool      // Whether the Testpoint executed ok
	Num         uint      // The number of the test
	Description string    // A short description
	Directive   Directive // Whether the test was skipped or is a todo
	Explanation string    // A short explanation why the test was skipped/is a todo
	Diagnostic  string    // A more detailed diagnostic message about the failed test
	Yaml        []byte    // The inline Yaml-document, if given
}

// The outcome of a Testsuite
//
// Deprecated: Project is unmaintained.
type Testsuite struct {
	Ok    bool        // Whether the Testsuite as a whole succeded
	Tests []*Testline // Description of all Testlines
	plan  int         // Number of tests intended to run
}

// Parses TAP
//
// Deprecated: Project is unmaintained.
type Parser struct {
	r     *bufio.Reader
	line  string
	suite Testsuite
}

func (p *Parser) parseLine(line string) (*Testline, error) {
	var err error

	matches := testlineRE.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("Does not match Testline: \"%s\"", line)
	}

	t := new(Testline)

	t.Ok = (len(matches[1]) == 0)

	if len(matches[2]) > 0 {
		var i int
		i, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("Could not parse Testnumber \"%s\"", matches[2])
		}
		t.Num = uint(i)
	}

	t.Description = matches[3]

	switch strings.ToLower(matches[4]) {
	case "":
		t.Directive = None
	case "todo":
		t.Directive = Todo
	case "skip":
		t.Directive = Skip
	}

	t.Explanation = matches[5]

	return t, nil
}

// Create a new TAP-Parser from the given reader
//
// Deprecated: Project is unmaintained.
func NewParser(r io.Reader) (*Parser, error) {
	p := &Parser{bufio.NewReader(r), "", Testsuite{true, nil, -1}}

	line, err := p.r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)

	if versionRE.MatchString(line) {
		line, err = p.r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
	}

	var matches []string
	if matches = planRE.FindStringSubmatch(line); matches != nil {
		i, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("Could not parse plan \"%s\": %s", matches[0], err)
		}
		p.suite.plan = i

		line, err = p.r.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = strings.TrimSpace(line)
	}

	p.line = line

	return p, nil
}

// Get the next Testline
//
// Deprecated: Project is unmaintained.
func (p *Parser) Next() (*Testline, error) {
	if len(p.line) == 0 {
		return nil, io.EOF
	}
	t, err := p.parseLine(p.line)
	if err != nil {
		return nil, err
	}
	p.line = ""

	var line string
	for {
		line, err = p.r.ReadString('\n')
		switch err {
		case nil:
		case io.EOF:
			if len(line) == 0 {
				p.suite.Tests = append(p.suite.Tests, t)
				return t, nil
			}
		default:
			return nil, err
		}
		line = strings.TrimSpace(line)

		var matches []string
		if matches = diagnosticsRE.FindStringSubmatch(line); matches != nil {
			if len(t.Diagnostic) == 0 {
				t.Diagnostic = matches[1]
			} else {
				t.Diagnostic = t.Diagnostic + "\n" + matches[1]
			}
			continue
		}

		if yamlStartRE.MatchString(line) {
			for {
				yaml, err := p.r.ReadBytes('\n')
				switch err {
				case nil:
				case io.EOF:
					if len(line) == 0 {
						p.suite.Tests = append(p.suite.Tests, t)
						return t, nil
					}
				default:
					return nil, err
				}
				if yamlEndRE.Match(yaml) {
					break
				}
				buf := make([]byte, len(t.Yaml)+len(yaml))
				copy(buf[:len(t.Yaml)], t.Yaml)
				copy(buf[len(t.Yaml):], yaml)
			}
			continue
		}

		if matches = planRE.FindStringSubmatch(line); matches != nil {
			if p.suite.plan != -1 {
				return nil, fmt.Errorf("Double plan")
			}
			i, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("Could not parse plan \"%s\": %s", matches[0], err)
			}
			p.suite.plan = i
			p.suite.Tests = append(p.suite.Tests, t)
			return t, nil
		}
		break
	}

	p.line = line
	p.suite.Tests = append(p.suite.Tests, t)

	return t, nil
}

// Get the whole Testsuite.
// This will block until the underlying reader reaches EOF or there is an error.
//
// Deprecated: Project is unmaintained.
func (p *Parser) Suite() (*Testsuite, error) {
	for {
		t, err := p.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !t.Ok {
			p.suite.Ok = false
		}
	}
	if p.suite.plan == 0 {
		p.suite.Ok = false
		return &p.suite, nil
	}
	if len(p.suite.Tests) != p.suite.plan {
		p.suite.Ok = false
		return &p.suite, nil
	}

	return &p.suite, nil
}

// Summarizes the Testline into a short String.
//
// The string will be formatted as followed: First a status (ok/fail/todo/skip)
// in […], then the description of the test, last the diagnostic-message, if
// the test failed and had such a message attached to it, or an explanation, if
// the test had a TODO/SKIP directive with an explanation attached.
//
// Deprecated: Project is unmaintained.
func (t *Testline) String() string {
	s := "["
	switch t.Directive {
	case None:
		if t.Ok {
			s += " ok "
		} else {
			s += "fail"
		}
	case Todo:
		s += "todo"
	case Skip:
		s += "skip"
	}
	s += "] "
	s += t.Description

	if t.Directive != None {
		if len(t.Explanation) > 0 {
			s += " # " + t.Explanation
		}
	} else if !t.Ok {
		if len(t.Diagnostic) > 0 {
			s += " # " + t.Diagnostic
		}
	}

	return s
}
