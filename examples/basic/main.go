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
		data := &myData{}

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

type myData struct {
	Count int
}

// Double example of how you'd handle a "computed" property
func (m *myData) Double() int {
	return m.Count * 2
}
