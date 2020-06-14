// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package reflex

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"reflex/client"

	"github.com/gorilla/websocket"
)

// Page defines the Data and events used to run the reflex Template
// if Upgrader is not set, defaults to buffer sizes of 512
// if ErrorHandler is not set, defaults to http.Error 500 and logs error to stdout
type Page struct {
	Events       EventFuncs
	Data         interface{}
	Upgrader     websocket.Upgrader
	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func defaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Print(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// EventFuncs can have Event, *http.Request, or websocket.Conn parameters
type EventFuncs map[string]interface{}

// SetupFunc is the function that builds and returns the Page elements for use with the reflex template
type SetupFunc func() *Page

// Template defines a reflex template that update and respond to DOM events
type Template struct {
	text string
}

// Parse creates a reflex template from the passed in text
func Parse(text string) *Template {
	return &Template{text: text}
}

// ParseFile creates a reflex template from the passed in file location
func ParseFile(file string) (*Template, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}
	return &Template{text: string(b)}, nil
}

// Must is similar to core Go template.Must, panics if an error is thrown during parsing
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}

	return t
}

// Setup sets up the reflex template for use as a standard http.Handler
func (t *Template) Setup(setup SetupFunc) http.Handler {
	// compare sending entire new template vs diff w/ https://github.com/sergi/go-diff

	pg := setup()
	if pg.ErrorHandler == nil {
		pg.ErrorHandler = defaultErrorHandler
	}

	funcs := map[string]interface{}{
		"client": func() template.HTML {
			return template.HTML(client.Inject)
		},
	}

	for name := range pg.Events {
		funcs[name] = func() template.JS {
			// TODO: Event arguments
			return template.JS(`reflex.event(event, '` + name + `');`)
		}
	}

	tmpl := template.Must(template.New("reflex-template").Funcs(funcs).Parse(t.text))

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
	if websocket.IsWebSocketUpgrade(r) {
		p.handleWebsocket(w, r)
		return
	}
	p.handleTemplate(w, r)
}

type eventCall struct {
	Name  string
	Args  []reflect.Value
	Event Event
}

func (p *page) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := p.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.ErrorHandler(w, r, err)
		return
	}

	defer func() {
		ws.Close()
	}()

	for {
		e := &eventCall{}
		err = websocket.ReadJSON(ws, e)
		if err != nil {
			// TODO: Handle websocket disconnects
			// preserve template state and try to reconnect?
			p.ErrorHandler(w, r, err)
			return
		}
		fn, ok := p.Events[e.Name]
		if !ok {
			continue
		}

		fnVal := reflect.ValueOf(fn)
		fnType := reflect.TypeOf(fn)

		args := make([]reflect.Value, 0, fnType.NumIn())
		for i := 0; i < cap(args); i++ {
			if fnType.In(i) == reflect.TypeOf(e.Event) {
				args = append(args, reflect.ValueOf(e.Event))
			} else if fnType.In(i) == reflect.TypeOf(r) {
				args = append(args, reflect.ValueOf(r))
			}
		}

		args = append(args, e.Args...)

		out := fnVal.Call(args)
		if len(out) > 0 {
			err, ok := out[0].Interface().(error)
			if ok {
				p.ErrorHandler(w, r, err)
				return
			}
		}
		// TODO: ignore all other returns?

		var b bytes.Buffer
		p.template.Execute(&b, p.Data)
		ws.WriteJSON(b.String())
	}
}

func (p *page) handleTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var b bytes.Buffer
	err := p.template.Execute(&b, p.Data)

	if err != nil {
		p.ErrorHandler(w, r, err)
		return
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		log.Printf("Error Copying template data to template writer: %s", err)
	}
}
