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

	http.Handle("/", reflex.ParseFiles("index.template.html").Setup(func() *reflex.Page {
		data := struct {
			Count int
		}{
			Count: 0,
		}

		return &reflex.Page{
			Data: data,
			Events: reflex.EventFuncs{
				"increment": func(i int) {
					data.Count += i
				},
			},
		}
	}))

	err := server.ListenAndServe()

	log.Fatalf("Error Starting server: %s", err)
}
