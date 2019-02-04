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

package generator_test

import (
	"testing"

	"github.com/lodge93/drone-fpm/pkg/generator"
	"github.com/lodge93/drone-fpm/pkg/parser"
	"io/ioutil"
	"os"
)

const ExpectedOutput = `#!/bin/bash

set -e

COMMAND="fpm "

if [[ -n "$PLUGIN_FOO" ]]; then
    COMMAND="$COMMAND --foo"
fi

if [[ -n "$PLUGIN_FOO" ]]; then
    COMMAND="$COMMAND --bar $PLUGIN_FOO"
fi

echo $COMMAND`

func TestNewGenerator(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "generated-entrypoint")
	if err != nil {
		t.Fatalf("error creating temp file - %s", err)
	}
	defer os.Remove(tmpFile.Name())

	g := generator.NewGenerator()
	err = g.GenerateEntrypoint("testdata/entrypoint.sh", tmpFile.Name(), []parser.ParsedFlag{
		{
			Option: "--foo",
			EnvVar: "PLUGIN_FOO",
			HasInput: false,
		},
		{
			Option: "--bar",
			EnvVar: "PLUGIN_FOO",
			HasInput: true,
		},
	})

	if err != nil {
		t.Fatalf("could not generate entrypoint - %s", err)
	}

	b, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Fatalf("could not read generated entrypoint - %s", err)
	}

	actual := string(b)
	if actual != ExpectedOutput {
		t.Errorf("acutal '%s' does not match expected '%s'", actual, ExpectedOutput)
	}
}
