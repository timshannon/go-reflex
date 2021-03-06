# go-reflex
An experimental library for building reactive templates for web applications inspired by Vue.

## Goals
The goal of this project is to try to allow for building frontend code in Go similarily to how you would in a frontend
framework like [Vue](https://vuejs.org), i.e. responding to data changes and events automatically.

* Build dynamic and reactive web pages without having to write any javascript
* Use standard library http handle so it can be used with any library that works with the go http package
* Use standard `html/template` without requiring any special syntax

## Overview
When a user connects to an endpoint, we'll create a websocket connection.  The template and any template changes (diffs)
will be sent on that websocket connection.  Client events will be sent back on that same websocket to initiate changes in
the template.  


## Use Case
Right now this is an experiment.  I have no idea if this is going to work well, if at all.  I know that there is no
way the reactive performance of this library will every compare to client js libraries.  If you are chasing for the
peak of client responsiveness, this library probably isn't for you. If you want to build responsive web UI's without
having to touch **any** javascript, then maybe this will be useful to you.  

So to summarize: 

**No** for high performance ui applications.  
**Yes** for for *"I need a web ui in my Go project and don't 
want to install NPM"*.

## Preliminary API Example

`index.template.html`
```html
<!doctype html>
<html lang="en">
    <head>
        <title>Example</title>
    </head>
    <body>
        <button onclick="{{increment 2}}" class="btn btn-primary">
            Count is: {{ .Count }}
        </button>
        {{client}} <!-- inject js client -->
    </body>
</html>
```

`Go Code`
```go
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

```