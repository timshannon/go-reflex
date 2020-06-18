// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package reflex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"reflex/client"

	"github.com/gorilla/websocket"
	"golang.org/x/net/html"
)

// Page defines the Data and events used to run the reflex Template
// if Upgrader is not set, defaults to buffer sizes of 512
// if ErrorHandler is not set, defaults to http.Error 500 and logs error to stdout
type Page struct {
	ElementID    string
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

var requestType = reflect.TypeOf(&http.Request{})

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

	pg := setup()

	if reflect.TypeOf(pg.Data).Kind() != reflect.Ptr {
		panic("Page Data must be a pointer")
	}

	if pg.ElementID == "" {
		panic("Element ID must be set")
	}

	funcs := map[string]interface{}{
		"client": func() template.HTML {
			return template.HTML(client.Inject)
		},
	}

	for name := range pg.Events {
		eventName := name
		funcs[eventName] = func(in ...interface{}) (template.JS, error) {
			js := "reflex.event(event, '" + eventName + "'"
			if len(in) > 0 {
				args := make([]interface{}, 0, len(in))
				for i := range in {
					argType := reflect.TypeOf(in[i])
					if argType != eventType && argType != requestType {
						args = append(args, in[i])
					}

				}
				param, err := json.Marshal(args)
				if err != nil {
					return "", err
				}
				js += "," + string(param)
			}
			js += ");"
			return template.JS(js), nil
		}
	}

	tmpl := template.Must(template.New("reflex-template").Funcs(funcs).Parse(t.text))

	return &handler{
		setup:    setup,
		template: tmpl,
	}
}

type handler struct {
	setup    SetupFunc
	template *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := h.setup()
	if p.ErrorHandler == nil {
		p.ErrorHandler = defaultErrorHandler
	}
	if websocket.IsWebSocketUpgrade(r) {
		h.handleWebsocket(p, w, r)
		return
	}
	h.handleTemplate(p, w, r)
}

type eventCall struct {
	Name  string        `json:"name"`
	Args  []interface{} `json:"args"`
	Event Event         `json:"event"`
}

func (h *handler) handleWebsocket(p *Page, w http.ResponseWriter, r *http.Request) {
	ws, err := p.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.ErrorHandler(w, r, err)
		return
	}

	defer func() {
		ws.Close()
	}()

	// send inital page data
	ws.WriteJSON(struct {
		ElementID string `json:"elementID"`
	}{
		ElementID: p.ElementID,
	})

	for {
		e := &eventCall{}
		err = websocket.ReadJSON(ws, e)
		if err != nil {
			if err == websocket.ErrCloseSent {
				return
			}
			// TODO: Handle websocket disconnects
			// preserve template state and try to reconnect?
			// send a UUID on first connect, and if reconnect with same UUID load template state from memory?
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
			inType := fnType.In(i)
			if inType == eventType {
				args = append(args, reflect.ValueOf(e.Event))
			} else if inType == requestType {
				args = append(args, reflect.ValueOf(r))
			} else {
				for j := range e.Args {
					val := reflect.ValueOf(e.Args[j])
					if val.Type() == inType {
						args = append(args, val)
						e.Args = append(e.Args[:j], e.Args[j+1:]...)
						break
					} else if val.Type().ConvertibleTo(inType) {
						args = append(args, val.Convert(inType))
						e.Args = append(e.Args[:j], e.Args[j+1:]...)
						break
					}
				}
			}
		}

		out := fnVal.Call(args)
		if len(out) > 0 {
			err, ok := out[0].Interface().(error)
			if ok {
				p.ErrorHandler(w, r, err)
				return
			}
		}
		// TODO: ignore all other func returns? Error?

		var b bytes.Buffer
		h.template.Execute(&b, p.Data)
		doc, err := html.Parse(&b)
		if err != nil {
			p.ErrorHandler(w, r, err)
			return
		}

		b.Reset()

		el := findElement(p.ElementID, doc)
		if el == nil {
			p.ErrorHandler(w, r, fmt.Errorf("No element found with an id of %s in the template", p.ElementID))
			return
		}

		err = html.Render(&b, el)
		if err != nil {
			p.ErrorHandler(w, r, err)
			return
		}

		ws.WriteMessage(websocket.TextMessage, b.Bytes())
	}
}

func (h *handler) handleTemplate(p *Page, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var b bytes.Buffer
	err := h.template.Execute(&b, p.Data)

	if err != nil {
		p.ErrorHandler(w, r, err)
		return
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		log.Printf("Error Copying template data to template writer: %s", err)
	}
}

func findElement(id string, parent *html.Node) *html.Node {
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && len(c.Attr) > 0 {
			for i := range c.Attr {
				if c.Attr[i].Key == "id" && c.Attr[i].Val == id {
					return c
				}
			}
		}

		el := findElement(id, c)
		if el != nil {
			return el
		}
	}
	return nil
}
