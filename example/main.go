// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"html/template"
	"log"
	"net/http"
	"reflex"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}

	http.Handle("/", reflex.Setup(template.Must(template.ParseFiles("index.template.html")), func() *reflex.Page {
		data := struct {
			Count int
		}{
			Count: 0,
		}

		return &reflex.Page{
			Data: data,
			Events: map[string]reflex.EventFunc{
				"increment": func(e *reflex.Event) {
					data.Count++
				},
			},
		}
	}))

	err := server.ListenAndServe()

	log.Fatalf("Error Starting server: %s", err)
}
