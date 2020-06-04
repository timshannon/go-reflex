# reflex-template
An experimental library for building reactive templates for web applications inspired by Vue.

## Goals
The goal of this project is to try to allow for building frontend code in Go similarily to how you would in a frontend
framework like [Vue](https://vuejs.org), i.e. responding to data changes and events automatically.


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
