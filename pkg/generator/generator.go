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

// Package generator provides mechanisms to generate an entrypoint using go template.
package generator

import (
	"text/template"
	"github.com/lodge93/drone-fpm/pkg/parser"
	"os"
)

// Generator is used to generate the entrypoint for the drone plugin.
type Generator struct {}

// NewGenerator creates a new, initialized generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateEntrypoint generates an entrypoint script using the source go template and destination file using the
// provided flags.
func (g *Generator) GenerateEntrypoint(src string, dst string, flags []parser.ParsedFlag) error {
	templ, err := template.ParseFiles(src)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	return templ.Execute(out, flags)
}