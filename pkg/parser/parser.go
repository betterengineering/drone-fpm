// Copyright 2019 Mark Spicer
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package parser parses fpm help output into a usable object model.
package parser

import (
	"os/exec"
	"strings"
)

// PluginPrefix is the prefix drone will add to all options passed into the plugin as environment variables.
const PluginPrefix = "PLUGIN_"

// Parser is a struct containing the lines to parse.
type Parser struct {
	// Lines is a slice of strings containing the fpm help output.
	Lines []string
}

// ParsedFlag contains a flag parsed from the fpm help output.
type ParsedFlag struct {
	// Option is the command line option for fpm.
	Option string

	// EnvVar is the environment variable that drone will use to inject the variable into the container.
	EnvVar string

	// HasInput is used to determine if the flag contains input or not.
	HasInput bool
}

// ParsedOptionLine is a structured representation of an option line for fpm's help message.
type ParsedOptionLine struct {
	// Options is a comma separated list of potential options for the flag in question.
	Options  string

	// Docs is the doc string for the command line flag.
	Docs     string

	// HasInput is used to determine if the flag contains input or not.
	HasInput bool
}

// NewInitialisedParser creates a new parser populated with output from fpm -h.
func NewInitialisedParser() (*Parser, error) {
	out, err := exec.Command("fpm", "-h").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")

	return &Parser{
		Lines: lines,
	}, nil
}

// Parse will return the output of fpm -h as a slice of structured flag objects.
func (p *Parser) Parse() ([]ParsedFlag, error) {
	parsedFlags := []ParsedFlag{}

	optionsIndex := 0
	for index, line := range p.Lines {
		if line == "Options:" {
			optionsIndex = index
		}
	}

	options := p.Lines[optionsIndex+1 : len(p.Lines)-1]
	for _, option := range options {
		parsedOptionLine := p.parseOptionLine(option)
		parsedFlagsFromOptions := p.convertParsedOptionLine(parsedOptionLine)
		parsedFlags = append(parsedFlags, parsedFlagsFromOptions...)
	}

	return parsedFlags, nil
}

func (p *Parser) parseOptionLine(line string) ParsedOptionLine {
	// Remove proceeding and trailing whitespace.
	lineTrimmed := strings.TrimSpace(line)

	// Split string based on a space as a separator.
	portions := strings.Split(lineTrimmed, " ")

	options := ""
	endOfCommand := false
	hasInput := false
	inputPassed := false
	doc := ""

	for _, portion := range portions {
		portion = strings.Replace(portion, " ", "", -1)
		if portion == "" {
			continue
		}

		// A option has been found, but there are more to be found.
		if strings.HasPrefix(portion, "-") && strings.HasSuffix(portion, ",") && !endOfCommand {
			options += portion
			continue
		}

		// The last option in the option line has been found.
		if strings.HasPrefix(portion, "--") && !strings.HasSuffix(portion, ",") && !endOfCommand {
			options += portion
			endOfCommand = true
			continue
		}

		// Check for the input variable, but make sure that we have not already passed that.
		if strings.ToUpper(portion) == portion && !inputPassed {
			hasInput = true
			inputPassed = true
			continue
		}
		inputPassed = true

		// Add the rest of the string to the doc.
		if doc == "" {
			doc = portion
			continue
		}
		doc = doc + " " + portion
	}

	return ParsedOptionLine{
		Options:  options,
		Docs:     doc,
		HasInput: hasInput,
	}
}

func (p *Parser) convertParsedOptionLine(parsedOptionLine ParsedOptionLine) []ParsedFlag {
	parsedFlags := []ParsedFlag{}
	for _, option := range strings.Split(parsedOptionLine.Options, ",") {
		// Skip shorthand options
		if !strings.HasPrefix(option, "--") {
			continue
		}

		// Add no prefixed options for both with and without no
		if strings.Contains(option, "--[no-]") {
			postive := "--" + strings.TrimPrefix(option, "--[no-]")
			parsedFlag := ParsedFlag{
				Option: postive,
				EnvVar: p.generateEnvVarFromOption(postive),
				HasInput: parsedOptionLine.HasInput,
			}
			parsedFlags = append(parsedFlags, parsedFlag)

			negative := "--no-" + strings.TrimPrefix(option, "--[no-]")
			parsedFlag = ParsedFlag{
				Option: negative,
				EnvVar: p.generateEnvVarFromOption(negative),
				HasInput: parsedOptionLine.HasInput,
			}
			parsedFlags = append(parsedFlags, parsedFlag)
			continue
		}

		// Add regular option.
		parsedFlag := ParsedFlag{
			Option:   option,
			EnvVar:   p.generateEnvVarFromOption(option),
			HasInput: parsedOptionLine.HasInput,
		}
		parsedFlags = append(parsedFlags, parsedFlag)
	}

	return parsedFlags
}

func (p *Parser) generateEnvVarFromOption(option string) string {
	optionName := strings.TrimLeft(option, "-")
	underscores := strings.Replace(optionName, "-", "_", -1)
	uppercase := strings.ToUpper(underscores)
	return PluginPrefix + uppercase
}
