// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package reflex

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Page struct {
	TemplateText  string
	TemplateFiles []string
	TemplateGlob  string
	Events        EventFuncs
	Data          interface{}
}

// EventFuncs can have Event or *http.Request parameters
type EventFuncs map[string]interface{}

type Event struct {
	// Example?: https://godoc.org/honnef.co/go/js/dom#BasicEvent
}

type SetupFunc func() *Page

type Template struct {
	text    string
	files   []string
	pattern string
}

func Parse(text string) *Template {
	return &Template{text: text}
}

func ParseFiles(files ...string) *Template {
	return &Template{files: files}
}

func ParseGlob(pattern string) *Template {
	return &Template{pattern: pattern}
}

func (t *Template) Setup(setup SetupFunc) http.Handler {
	// compare sending entire new template vs diff w/ https://github.com/sergi/go-diff

	pg := setup()

	funcs := map[string]interface{}{
		"client": func() template.JS {
			return template.JS(`<script type="text/javascript">var reflex = {
					event: function(name) {
						console.log("Event: ", name);
					},
				};</script>`)
		},
	}

	for k := range pg.Events {
		funcs[k] = func() template.JS {
			return template.JS(`reflex.event('` + k + `');`)
		}
	}

	tmpl := template.New("reflex-template").Funcs(funcs)

	if t.files != nil {
		tmpl = template.Must(tmpl.ParseFiles(t.files...))
	} else if t.pattern != "" {
		tmpl = template.Must(tmpl.ParseGlob(t.pattern))
	} else if t.text != "" {
		tmpl = template.Must(tmpl.Parse(t.text))
	}

	return &page{
		Page:     pg,
		template: tmpl,
	}
}

type page struct {
	*Page
	template *template.Template
}

func (p *page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: make ws connection

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var b bytes.Buffer
	err := p.template.Execute(&b, p.Data)

	if err != nil {
		// TODO: error handling
		log.Printf("Error executing template: %s", err)
	} else {
		_, err = io.Copy(w, &b)
		if err != nil {
			log.Printf("Error Copying template data to template writer: %s", err)
		}
	}
}
