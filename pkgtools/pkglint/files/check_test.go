package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/check.v1"
	"netbsd.org/pkglint/textproc"
	"netbsd.org/pkglint/trace"
)

var equals = check.Equals
var deepEquals = check.DeepEquals

type Suite struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	tmpdir string
	checkC *check.C
}

// Init initializes the suite with the check.C instance for the actual
// test run. See https://github.com/go-check/check/issues/22
func (s *Suite) Init(c *check.C) {
	if s.checkC != nil {
		panic("Suite.Init must only be called once.")
	}
	s.checkC = c
}

func (s *Suite) c() *check.C {
	if s.checkC == nil {
		panic("Suite.Init must be called before accessing check.C.")
	}
	return s.checkC
}

func (s *Suite) Stdout() string {
	defer s.stdout.Reset()
	return s.stdout.String()
}

func (s *Suite) Stderr() string {
	defer s.stderr.Reset()
	return s.stderr.String()
}

// Returns and consumes the output from both stdout and stderr.
// The temporary directory is replaced with a tilde (~).
func (s *Suite) Output() string {
	output := s.Stdout() + s.Stderr()
	if s.tmpdir != "" {
		output = strings.Replace(output, s.tmpdir, "~", -1)
	}
	return output
}

func (s *Suite) CheckOutputEmpty() {
	s.c().Check(s.Output(), equals, "")
}

func (s *Suite) CheckOutputLines(expectedLines ...string) {
	expectedOutput := ""
	for _, expectedLine := range expectedLines {
		expectedOutput += expectedLine + "\n"
	}
	s.c().Check(s.Output(), equals, expectedOutput)
}

func (s *Suite) BeginDebugToStdout() {
	G.logOut = NewSeparatorWriter(os.Stdout)
	trace.Out = os.Stdout
	trace.Tracing = true
}

func (s *Suite) EndDebugToStdout() {
	G.logOut = NewSeparatorWriter(&s.stdout)
	trace.Out = &s.stdout
	trace.Tracing = false
}

// UseCommandLine simulates a command line for the remainder of the test.
// See Pkglint.ParseCommandLine
func (s *Suite) UseCommandLine(args ...string) {
	exitcode := new(Pkglint).ParseCommandLine(append([]string{"pkglint"}, args...))
	if exitcode != nil && *exitcode != 0 {
		s.CheckOutputEmpty()
		s.c().Fatalf("Cannot parse command line: %#v", args)
	}
	G.opts.LogVerbose = true // See SetUpTest
}

func (s *Suite) RegisterMasterSite(varname string, urls ...string) {
	name2url := &G.globalData.MasterSiteVarToURL
	url2name := &G.globalData.MasterSiteURLToVar
	if *name2url == nil {
		*name2url = make(map[string]string)
		*url2name = make(map[string]string)
	}
	(*name2url)[varname] = urls[0]
	for _, url := range urls {
		(*url2name)[url] = varname
	}
}

func (s *Suite) RegisterTool(tool *Tool) {
	reg := G.globalData.Tools

	if len(reg.byName) == 0 && len(reg.byVarname) == 0 {
		reg = NewToolRegistry()
		G.globalData.Tools = reg
	}
	if tool.Name != "" {
		reg.byName[tool.Name] = tool
	}
	if tool.Varname != "" {
		reg.byVarname[tool.Varname] = tool
	}
}

func (s *Suite) CreateTmpFile(relFname, content string) (absFname string) {
	c := s.c()
	absFname = s.TmpDir() + "/" + relFname
	err := os.MkdirAll(path.Dir(absFname), 0777)
	c.Assert(err, check.IsNil)

	err = ioutil.WriteFile(absFname, []byte(content), 0666)
	c.Check(err, check.IsNil)
	return
}

func (s *Suite) CreateTmpFileLines(relFname string, rawTexts ...string) (absFname string) {
	text := ""
	for _, rawText := range rawTexts {
		text += rawText + "\n"
	}
	return s.CreateTmpFile(relFname, text)
}

func (s *Suite) LoadTmpFile(relFname string) (absFname string) {
	c := s.c()
	bytes, err := ioutil.ReadFile(s.TmpDir() + "/" + relFname)
	c.Assert(err, check.IsNil)
	return string(bytes)
}

func (s *Suite) TmpDir() string {
	if s.tmpdir == "" {
		s.tmpdir = filepath.ToSlash(s.c().MkDir())
	}
	return s.tmpdir
}

func (s *Suite) ExpectFatalError(action func()) {
	if r := recover(); r != nil {
		if _, ok := r.(pkglintFatal); ok {
			action()
			return
		}
		panic(r)
	}
}

func (s *Suite) SetUpTest(c *check.C) {
	G = GlobalVars{Testing: true}
	textproc.Testing = true
	G.logOut, G.logErr, trace.Out = NewSeparatorWriter(&s.stdout), NewSeparatorWriter(&s.stderr), &s.stdout
	s.checkC = c
	s.UseCommandLine( /* no arguments */ )
	s.checkC = nil
	G.opts.LogVerbose = true // To detect duplicate work being done
}

func (s *Suite) TearDownTest(c *check.C) {
	G = GlobalVars{}
	textproc.Testing = false
	if out := s.Output(); out != "" {
		fmt.Fprintf(os.Stderr, "Unchecked output in %q; check with: c.Check(s.Output(), equals, %q)", c.TestName(), out)
	}
	s.tmpdir = ""
}

var _ = check.Suite(new(Suite))

func Test(t *testing.T) { check.TestingT(t) }

var T TestObjectCreator

type TestObjectCreator struct{}

// Arguments are either (lineno, orignl) or (lineno, orignl, textnl).
func (t TestObjectCreator) NewRawLines(args ...interface{}) []*RawLine {
	rawlines := make([]*RawLine, len(args)/2)
	j := 0
	for i := 0; i < len(args); i += 2 {
		lineno := args[i].(int)
		orignl := args[i+1].(string)
		textnl := orignl
		if i+2 < len(args) {
			if s, ok := args[i+2].(string); ok {
				textnl = s
				i++
			}
		}
		rawlines[j] = &RawLine{lineno, orignl, textnl}
		j++
	}
	return rawlines[:j]
}

func (t TestObjectCreator) NewLine(filename string, lineno int, text string) Line {
	textnl := text + "\n"
	rawLine := RawLine{lineno, textnl, textnl}
	return NewLine(filename, lineno, text, []*RawLine{&rawLine})
}

func (t TestObjectCreator) NewMkLine(fileName string, lineno int, text string) MkLine {
	return NewMkLine(t.NewLine(fileName, lineno, text))
}

func (t TestObjectCreator) NewShellLine(fileName string, lineno int, text string) *ShellLine {
	return NewShellLine(t.NewMkLine(fileName, lineno, text))
}

// NewLines generates a slice of simple lines,
// i.e. each logical line has exactly one physical line.
// To work with line continuations like in Makefiles,
// use Suite.CreateTmpFileLines together with Suite.LoadExistingLines.
func (t TestObjectCreator) NewLines(fname string, texts ...string) []Line {
	return t.NewLinesAt(fname, 1, texts...)
}

func (t TestObjectCreator) NewMkLines(fname string, lines ...string) *MkLines {
	return NewMkLines(t.NewLines(fname, lines...))
}

// NewLinesAt generates a slice of simple lines,
// i.e. each logical line has exactly one physical line.
// To work with line continuations like in Makefiles,
// use Suite.CreateTmpFileLines together with Suite.LoadExistingLines.
func (t TestObjectCreator) NewLinesAt(fname string, firstLine int, texts ...string) []Line {
	result := make([]Line, len(texts))
	for i, text := range texts {
		textnl := text + "\n"
		result[i] = NewLine(fname, i+firstLine, text, t.NewRawLines(i+firstLine, textnl))
	}
	return result
}
