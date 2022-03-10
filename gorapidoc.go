package rapidoc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type Config struct {
	DocsPath    string
	SpecFile    string
	Title       string
	Description string
}

var (
	HTML string
 	JavaScript string
)

func (c Config) Body() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	tpl, err := template.New("rapidoc").Parse(HTML)
	if err != nil {
		return nil, err
	}

	if err = tpl.Execute(buf, map[string]string{
		"title":       c.Title,
		"url":         c.SpecFile,
		"description": c.Description,
	}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c Config) Handler() http.HandlerFunc {
	data, err := c.Body()
	if err != nil {
		panic(err)
	}

	specFile := c.SpecFile

	if specFile == "" {
		panic(errors.New("spec not found"))
	}

	spec, err := ioutil.ReadFile(specFile)
	if err != nil {
		panic(err)
	}

	docsPath := c.DocsPath
	return func(w http.ResponseWriter, req *http.Request) {
		method := strings.ToLower(req.Method)
		if method != "get" && method != "head" {
			return
		}

		if strings.HasSuffix(req.URL.Path, c.SpecFile) {
			w.WriteHeader(200)
			w.Header().Set("content-type", "application/json")
			w.Write(spec)
			return
		}

		if docsPath == "" || docsPath == req.URL.Path {
			w.WriteHeader(200)
			w.Header().Set("content-type", "text/html")
			w.Write(data)
		}
	}
}