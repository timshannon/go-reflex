// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"reflex"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}

	http.Handle("/", reflex.Must(reflex.ParseFile("index.template.html")).Setup(func() *reflex.Page {
		data := &struct {
			Count int
		}{
			Count: 0,
		}

		return &reflex.Page{
			ElementID: "app",
			Data:      data,
			Events: reflex.EventFuncs{
				"increment": func(r *http.Request) {
					data.Count++
				},
			},
		}
	}))

	err := server.ListenAndServe()

	log.Fatalf("Error Starting server: %s", err)
}
