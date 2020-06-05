// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package reflex

import (
	"html/template"
	"net/http"
)

type Page struct {
	Events map[string]EventFunc
	Data   interface{}
}

type Event struct{}

type EventFunc func(*Event)

type SetupFunc func() *Page

func Setup(t *template.Template, setup SetupFunc) http.Handler {

	// TODO: Serve all the HTTPs
	return http.NotFoundHandler()
}
